import logging
import os
import time
from abc import ABCMeta, abstractmethod
from collections import Counter
from datetime import datetime

import numpy as np
from app.config import get_config
from app.utils import catch_error_for_all_methods
from pymongo import ASCENDING, HASHED, MongoClient, errors

# from bson.objectid import ObjectId


config = get_config()

class AbstractDatabase(metaclass=ABCMeta):
    def __init__(self):
        '''
        For communication with the database
        '''
        logging.info('Connect to db host: {} - username: {} - password: {} - port'.format(config.DATABASE_HOST, config.DATABASE_USERNAME, config.DATABASE_PASSWORD, config.DATABASE_PORT))
        if not config.TESTING:
            self.mongodb_client = MongoClient(config.DATABASE_HOST, config.DATABASE_PORT, username=config.DATABASE_USERNAME, password=config.DATABASE_PASSWORD)
            self.setup()

    @abstractmethod
    def setup(self):
        pass


@catch_error_for_all_methods
class FaceInfoDatabase(AbstractDatabase):
    """Database Class stores UserID - ImageID relation
    """
    def setup(self):
        self.mongodb_db = self.mongodb_client[config.DATABASE_NAME]
        self.mongodb_faceinfo = self.mongodb_db[config.DATABASE_FACE_INFO]
        self.mongodb_faceinfo.create_index([("imageId", HASHED)])
        self.mongodb_faceinfo.create_index([("userId", ASCENDING)])

    def find_user_uids_by_image_ids(self, image_ids):
        """Given a image_ids list, return user_uids, irrespectively

        Arguments:
            image_ids {list} -- image ids want to search

        Returns:
            [list] -- user ids
        """
        count = Counter(image_ids)
        user_uids = []
        cursor = self.mongodb_faceinfo.find({'imageId': { '$in': list(count.keys()) }})
        for key in cursor:
            for _ in range(count[key['imageId']]):
                user_uids.append(key['userId'])
        return user_uids

    def find_image_ids_by_user_id(self, user_uid):
        """Given a user_uid, return image_ids, irrespectively

        Arguments:
            user_uid {string} -- user id want to search

        Returns:
            [list] -- image ids
        """
        count = Counter([user_uid])
        image_ids = []
        cursor = self.mongodb_faceinfo.find({'userId': { '$in': list(count.keys()) }})
        for key in cursor:
            for _ in range(count[key['userId']]):
                image_ids.append(key['imageId'])
        return image_ids

    def insert_many_new_faces(self, faces):
        """Insert FaceInfo to Database
        
        Arguments:
            faces {list of dictionaries}
        """
        if len(faces):
            image_ids = [item['imageId'] for item in faces]
            self.mongodb_faceinfo.remove({'imageId': {'$in': image_ids}})
            self.mongodb_faceinfo.insert_many(faces)

    def remove_many_image_ids(self, image_ids):
        """Remove image IDs
        
        Arguments:
            image_ids {list}
        """
        if len(image_ids):
            self.mongodb_faceinfo.remove({'imageId': {'$in': image_ids}})

    def get_labels_and_embs(self, user_uids=None):
        '''
        Get  images and labels for building matcher
        '''
        if user_uids is None or len(user_uids) == 0:
            reg_image_face_dict = self.mongodb_faceinfo.find(projection={'_id': False})
        else:
            reg_image_face_dict = self.mongodb_faceinfo.find(
                {'userId': { '$in': user_uids } },
                projection={'_id': False}
            )
        labels = []
        embs = []

        for cursor in reg_image_face_dict:
            emb = np.array(cursor['embedding'])
            embs.append(emb)
            labels.append(cursor['imageId'])
        return np.array(labels), np.array(embs).squeeze()
