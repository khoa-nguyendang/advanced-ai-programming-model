import logging
import random
import sys

import app.base_server
import numpy as np
from app.config import get_config
from app.index import index
from app.utils import database, make_response

logging.basicConfig(format='%(asctime)s,%(msecs)03d %(levelname)-8s [%(filename)s:%(lineno)d] %(message)s',
    datefmt='%Y-%m-%d:%H:%M:%S',
    level=logging.DEBUG)

logger = logging.getLogger(__name__)

class FaissServer(app.base_server.AbstractServer):

    def __init__(self):
        super().__init__()
        config = get_config()
        self.app.config.from_object(config)

        self.app.logger.addHandler(logging.StreamHandler(sys.stderr))
        self.app.logger.setLevel('INFO')
        self.app.before_first_request(self._preload_index)
        self.database = database.FaceInfoDatabase()

    def add_endpoint(self):
        self.app.add_url_rule('/health/', 'health', self._health, methods=['GET'])
        self.app.add_url_rule('/search/', 'search', self._search, methods=['POST'])
        self.app.add_url_rule('/register/', 'register', self._register, methods=['POST'])
        self.app.add_url_rule('/deregister/', 'deregister', self._deregister, methods=['POST'])

    def _preload_index(self):
        """
        Before the first request, preload index from disk if it is not already done
        """
        _ = index(database=self.database)

    def _health(self):
        """
        Health check endpoint

        If the index fails to load, this will fail immediately
        """
        if index() is not None:
            # Generate a vector to perform a search
            test_vector = [random.uniform(0, 1000) for _ in range(0, 512)]
            test_vector = np.array(test_vector, dtype=np.float32)[np.newaxis]
            # Reload the index if we see a newer update via the `force_reload` flag
            result = index().match(test_vector)
            result_index = np.argmax(result[0][0])
            score = result[0][0][result_index].item()

            if score is not None and isinstance(score, float):
                # A valid score is calculated
                return make_response({}, 200)
        # If we reach this point, the index is confirmed to not work
        return make_response({'reason': 'Index is not loaded correctly'}, 500)

    def _search(self):
        """
        Given a vector, perform a search query

        The input vector is a JSON object has the following format:
        ```
        {"vector": <vector>}
        ```
        where `<vector>` is simply an array of numbers
        """
        try:
            request_json = self.request.get_json()
            if request_json is None:
                return make_response({'reason': 'Malformed request body'}, 400)
            vector = request_json.get('vector')

            if not isinstance(vector, list):
                return make_response({'reason': '`vector` should be an array'}, 400)

            if len(vector) != self.app.config['DIMENSIONS']:
                return make_response({'reason': '`vector` should contain {0} entries'
                                      .format(self.app.config['DIMENSIONS'])}, 400)

            vector = np.array(vector, dtype=np.float32)[np.newaxis]
            result = index().match(vector)
            result_index = np.argmax(result[0][0])
            score = result[0][0][result_index].item()
            if score < self.app.config['MATCH_SCORE_THRESHOLD']:
                # This is the max score we've found - if the score does not pass the threshold,
                # we skip it anyway
                return make_response({
                    'reason': 'No matching item is found',
                    'score': score
                }, 404)  # Not found
            else:
                image_id = result[1][0][result_index].item()
                user_uid = self.database.find_user_uids_by_image_ids([image_id])
                user_uid = user_uid[0] if (user_uid is not None) and len(user_uid) else self.app.config['NEW_FACE']
                return make_response({
                    'faiss_id': image_id,
                    'user_uid': user_uid,
                    'score': score
                }, 200)  # Successful search
        except:
            import traceback
            traceback.print_exc()
            return make_response({'reason': 'Internal error encountered'}, 500)

    def _register(self):
        """
        Register a list of vectors into the index where each item in the list is actually a tuple (ID, vector)
        ```
        {"vectors": [[<vector ID>, <user ID>, <vector>]]}
        ```
        """
        try:
            request_json = self.request.get_json()
            if request_json is None:
                return make_response({'reason': 'Malformed request body'}, 400)
            vectors = request_json.get('vectors')

            # Perform a check on input vectors
            for _id, _user_id, vector in vectors:
                if not isinstance(vector, list):
                    return make_response({'reason': '`vector` at ID `{0}` should be an array'.format(_id)}, 400)

                if len(vector) != self.app.config['DIMENSIONS']:
                    return make_response({'reason': '`vector` at ID `{0}` should contain {1} entries'
                                          .format(_id, self.app.config['DIMENSIONS'])}, 400)

            ids = [x[0] for x in vectors]
            user_uids = [x[1] for x in vectors]
            vectors = [x[2] for x in vectors]
            # index().train(vectors)
            # Add them to the index
            index().update(vectors, ids)
            self.database.insert_many_new_faces(
                [dict(imageId=ids[i], userId=user_uids[i], embedding=vectors[i]) for i in range(len(ids))])
            # Then write to disk
            # write_index_to_disk()
            # Return a 200 to signal a successful creation
            return make_response({}, 200)  # Successful registration
        except:
            import traceback
            traceback.print_exc()
            return make_response({'reason': 'Internal error encountered'}, 500)

    def _deregister(self):
        """
        Deregister a list of vector IDs from the index
        ```
        {"ids": [<vector ID>]}
        ```
        """
        try:
            request_json = self.request.get_json()
            if request_json is None:
                return make_response({'reason': 'Malformed request body'}, 400)
            vector_ids = request_json.get('vector_ids')
            # ids = np.array(vector_ids)
            # Remove the specified IDs from the index
            index().remove(vector_ids)
            self.database.remove_many_image_ids(vector_ids)
            # Then write to disk
            # write_index_to_disk()
            return make_response({}, 200)  # Successful deregistration
        except:
            import traceback
            traceback.print_exc()
            return make_response({'reason': 'Internal error encountered'}, 500)
