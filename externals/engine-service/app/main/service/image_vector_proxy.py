import io
import json
import logging
import os.path
import time
from concurrent import futures

import app.main.grpc.ai_engine_pb2_grpc
import app.main.grpc.faiss_pb2_grpc
import cv2
import grpc
import numpy as np
from app.main.grpc.ai_engine_pb2 import (ImageEnrollmentRequest,
                                         ImageEnrollmentResponse,
                                         ImageSearchResponse, StatusCode)
from app.main.grpc.faiss_pb2 import StatusCode as Faiss_StatusCode
from app.main.grpc.faiss_pb2 import (Vector, VectorEnrollmentRequest,
                                     VectorSearchRequest)
from app.main.grpc.faiss_pb2_grpc import VectorAPIStub
from app.main.service import engine
from app.main.service.config import get_config
from app.main.service.utils import upload_file
from google.protobuf import message
from minio import Minio

config = get_config()

metadata_json_string = """
{
    "match_score": ""
}
"""

logging.basicConfig(format='%(asctime)s,%(msecs)03d %(levelname)-8s [%(filename)s:%(lineno)d] %(message)s',
    datefmt='%Y-%m-%d:%H:%M:%S',
    level=logging.DEBUG)

logger = logging.getLogger(__name__)


class EngineService(app.main.grpc.ai_engine_pb2_grpc.ImageAPIServicer):
    def __init__(self, faiss_url, engine_instance_number, minimum_enrollment_image, minio_client):
        self.faiss_grpc_channel = grpc.insecure_channel(faiss_url)
        self.faiss_grpc_client = VectorAPIStub(self.faiss_grpc_channel)
        self.minio_client = minio_client
        engine.test_engine()

    def Enroll(self, request: ImageEnrollmentRequest, context):
        logger.info("Enroll")
        vectors = []
        images = list(request.images)
        logger.info("Enroll with {} images".format(len(images)))

        user_uid = request.user_uid
        for image in list(request.images):
            try:
                '''from binary image, that pass though grpc, decode it and pass to prodet'''
                npimg = np.frombuffer(image.data, dtype=np.int8)
                img = cv2.imdecode(npimg, cv2.IMREAD_COLOR)
                try:
                    feature = engine.extract_feature(img)
                    if feature is None:
                        return ImageEnrollmentResponse(code = StatusCode.BAD_IMAGE, message = "Must verify with at least one face in an image")
                    if len(feature) == 0:
                        logger.error("Unable to extract feature for image {}".format(image.image_id))
                    else:
                        vectors.append(Vector(vector_id = image.image_id, data = feature))
                except Exception as e:
                    logger.error(e)
                    return ImageEnrollmentResponse(code = StatusCode.ENGINE_ERROR)
            except Exception as e:
                logger.error(e)
                return ImageEnrollmentResponse(code = StatusCode.BAD_IMAGE, message = "Cannot decode image")

        if len(vectors) == 0:
            logger.error("Unable to extract any feature for current request")
            return ImageEnrollmentResponse(code = StatusCode.ENGINE_ERROR)

        # Send to FAISS
        response = self.faiss_grpc_client.Enroll(VectorEnrollmentRequest(user_uid=user_uid, vectors = vectors))
        logger.info(response)
        if response.code == Faiss_StatusCode.SUCCESS:
            return ImageEnrollmentResponse(code = StatusCode.SUCCESS)
        elif response.code == Faiss_StatusCode.IMAGE_REGISTERED_FOR_ANOTHER_USER:
                return ImageEnrollmentResponse(code = StatusCode.IMAGE_REGISTERED_FOR_ANOTHER_USER, message = response.message)
        else :
            return ImageEnrollmentResponse(code = StatusCode.ENGINE_ERROR)

    def Search(self, request, context):
        start = time.time()
        logger.info("Search user")
        if len(request.images) > 0:
            try:
                npimg = np.frombuffer(request.images[0].data, dtype=np.int8)
                img = cv2.imdecode(npimg, cv2.IMREAD_COLOR)
                try:
                    feature = engine.extract_feature(img)
                    if feature is None:
                        return ImageEnrollmentResponse(code = StatusCode.BAD_IMAGE, message = "Must verify with at least one face in an image")
                    if len(feature) == 0:
                        logger.error("Unable to extract feature for image ")
                        return ImageSearchResponse(code = StatusCode.ENGINE_ERROR)

                    logger.info("Start Searching user via faiss")
                    response = self.faiss_grpc_client.Search(VectorSearchRequest(vectors = [Vector(data = feature)]))
                    logger.info('received data {}'.format(response))

                    if config.VERIFICATION_IMAGE_LOG == "ENABLED":
                        metadata = json.loads(metadata_json_string)
                        metadata["match_score"] = response.score
                        logger.info(metadata)
                        value_as_bytes = json.dumps(metadata).encode('utf-8')
                        upload_file(self.minio_client, value_as_bytes,  "/metadata.json")

                    if response.code == Faiss_StatusCode.FOUND:
                        logger.info("Running time {} ".format(time.time() - start))
                        if config.VERIFICATION_IMAGE_LOG == "ENABLED":
                            upload_file(self.minio_client, request.images[0].data, "/image_matched.jpg")
                        return ImageSearchResponse(code = StatusCode.FOUND, message = "Face image found", image_id = response.image_id, user_uid = response.user_uid, score = response.score)
                    elif response.code == Faiss_StatusCode.NOT_FOUND:
                        if config.VERIFICATION_IMAGE_LOG == "ENABLED":
                            upload_file(self.minio_client, request.images[0].data, "/image_unmatched.jpg")
                        return ImageSearchResponse(code = StatusCode.NOT_FOUND, score = response.score)
                    else :
                        return ImageSearchResponse(code = StatusCode.ENGINE_ERROR)
                except Exception as e:
                    logger.error('Search exception 1 {}', e)
                    return ImageSearchResponse(code = StatusCode.ENGINE_ERROR)
            except Exception as e:
                logger.error('Search exception 2 {}', e)
                return ImageSearchResponse(code = StatusCode.BAD_IMAGE, message = "Cannot decode image")
        else:
            return ImageSearchResponse(code = StatusCode.NOT_FOUND)

def run():
    minio_client = Minio(
        config.MINIO_ENDPOINT,
        access_key = config.MINIO_ROOT_USER,
        secret_key = config.MINIO_ROOT_PASSWORD,
        secure = False,
    )


    faiss_url = config.FAISS_SERVICE_HOST + ":" + config.FAISS_SERVICE_PORT
    minimum_enrollment_image = int(config.MINIMUM_ENROLLMENT_IMAGE)

    max_worker_and_engine = int(config.MAX_ENGINE)

    srv = EngineService(faiss_url, max_worker_and_engine, minimum_enrollment_image, minio_client)
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=100))

    app.main.grpc.ai_engine_pb2_grpc.add_ImageAPIServicer_to_server(
        srv, server
    )
    server.add_insecure_port("[::]:" + config.SERVER_PORT)
    server.start()
    server.wait_for_termination()