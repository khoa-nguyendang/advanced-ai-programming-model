import logging
import time
from concurrent import futures

import app.base_server
import app.faiss_pb2_grpc
import faiss
import grpc
import numpy as np
from app.config import get_config
from app.faiss_pb2 import (StatusCode, VectorDeletionResponse,
                           VectorEnrollmentResponse, VectorSearchResponse)
from app.index import index
from app.utils import database, make_response
from google.protobuf import message

config = get_config()
logging.basicConfig(format='%(asctime)s,%(msecs)03d %(levelname)-8s [%(filename)s:%(lineno)d] %(message)s',
    datefmt='%Y-%m-%d:%H:%M:%S',
    level=logging.DEBUG)

logger = logging.getLogger(__name__)

class FAISSService(
    app.faiss_pb2_grpc.VectorAPIServicer
):
    def __init__(self):
        super().__init__()
        self.database = database.FaceInfoDatabase()

    def _preload_index(self):
        """
        Before the first request, preload index from disk if it is not already done
        """
        _ = index(database=self.database)

    def _search(self, vector):
        vector = np.array(vector, dtype=np.float32)[np.newaxis]
        result = index().match(vector)
        if len(result[0][0]) > 0:
            result_index = np.argmax(result[0][0])
            score = result[0][0][result_index].item()
            logger.info("FAISS Searching score : {}".format(score))
            logger.info("FAISS MATCH_SCORE_THRESHOLD =  : {}".format(config.MATCH_SCORE_THRESHOLD))
            if score < config.MATCH_SCORE_THRESHOLD:
                return None, None, score
            else:
                image_id = result[1][0][result_index].item()
                user_uids = self.database.find_user_uids_by_image_ids([image_id])
                logger.info("FAISS Searching user_uid for image_id : {}".format(image_id))
                if len(user_uids):
                    user_uid = user_uids[0]
                    return user_uid, image_id, score
                else:
                    logger.info("FAISS user_uid NOT_FOUND")
        return None, None, None

    def Search(self, request, context):
        if len(request.vectors) > 0:
            start = time.time()
            vector = request.vectors[0]
            if len(vector.data) != config.DIMENSIONS:
                    return VectorEnrollmentResponse(code = StatusCode.BAD_VECTOR, message = 'Vector should contain 512 entries')
            vector_data = np.array([vector.data]).astype(np.float32)
            faiss.normalize_L2(vector_data)
            user_uid, image_id, score = self._search(vector_data[0])
            logger.info("Search time {} ".format(time.time() - start))
            if user_uid:
                logger.info("FAISS return : user_uid = {} image_id = {}".format(user_uid, image_id))
                return VectorSearchResponse(code = StatusCode.FOUND, message = "Found" , image_id = image_id, user_uid = user_uid, score = score)

        logger.info("FAISS return NOT_FOUND")
        return VectorSearchResponse(code = StatusCode.NOT_FOUND, message = 'No matching item is found', score = score)

    def Enroll(self, request, context):
        try:
            user_uid = request.user_uid
            req_vectors = list(request.vectors)

            # Perform a check on input vectors
            new_vectors = []
            for vector in req_vectors:
                try:
                    if len(vector.data) != config.DIMENSIONS:
                        return VectorEnrollmentResponse(code = StatusCode.BAD_VECTOR, message = 'Vector should contain 512 entries')
                    vector_data = np.array([vector.data]).astype(np.float32)
                    logger.info("vector input: {}".format(vector_data))
                    faiss.normalize_L2(vector_data)
                    logger.info("vector after normalized: {}".format(vector_data))
                    existed_user_id, image_id, score = self._search(vector_data)

                    logger.info('existed_user_id: {} -- user_uid: {} -- score: {}'.format(existed_user_id, user_uid, score))
                    if existed_user_id and user_uid != existed_user_id:
                        return VectorEnrollmentResponse(code = StatusCode.IMAGE_REGISTERED_FOR_ANOTHER_USER, message = "Vector registered for user id {} with image id {}".format(existed_user_id, image_id))
                    new_vector = [vector.vector_id, vector_data[0].tolist()]
                    new_vectors.append(new_vector)
                except Exception as err:
                    logger.error("Got error for vector input: {} -- error: {}".format(vector.data, err))
                    return VectorEnrollmentResponse(code = StatusCode.BAD_VECTOR, message = "Something wrong {}".format(err))


            ids = [x[0] for x in new_vectors]
            vectors = [x[1] for x in new_vectors]
            logger.info("FAISS vectors : {}".format(vectors))
            # index().train(vectors)
            # faiss.normalize_L2(vectors)
            # Add them to the index
            index().update(vectors, ids)
            logger.info("FAISS vectors size : {}".format(len(ids)))
            self.database.insert_many_new_faces([dict(imageId=ids[i], userId=user_uid, embedding=vectors[i]) for i in range(len(ids))])
            # Then write to disk
            # write_index_to_disk()
            return VectorEnrollmentResponse(code = StatusCode.SUCCESS, message = "Vectors created")  # Successful registration
        
        except Exception as err:
            logger.error("Got error for vector input: {} -- error: {}".format(vector.data, err))
            return VectorEnrollmentResponse(code = StatusCode.BAD_VECTOR, message = "Something wrong {}".format(err))

    def Delete(self, request, context):
        image_ids = list(request.image_ids)
        # ids = np.array(vector_ids)
        # Remove the specified IDs from the index
        if request.user_uid != "" :
            image_ids.extend(self.database.find_image_ids_by_user_id(request.user_uid))
        index().remove(image_ids)
        self.database.remove_many_image_ids(image_ids)
        # Then write to disk
        # write_index_to_disk()
        return VectorDeletionResponse(code = StatusCode.SUCCESS)  # Successful deregistration

def run():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=500))
    app.faiss_pb2_grpc.add_VectorAPIServicer_to_server(
        FAISSService(), server
    )
    print('config.SERVER_PORT' + config.SERVER_PORT)
    server.add_insecure_port("[::]:" + config.SERVER_PORT)
    server.start()
    server.wait_for_termination()