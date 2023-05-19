import logging
import os
import sys
import time
from abc import ABCMeta, abstractmethod

import faiss
import numpy as np
from app.config import get_config
from app.utils import rwlock

config = get_config()
logging.basicConfig(format='%(asctime)s,%(msecs)03d %(levelname)-8s [%(filename)s:%(lineno)d] %(message)s',
    datefmt='%Y-%m-%d:%H:%M:%S',
    level=logging.DEBUG)

logger = logging.getLogger(__name__)

class FaissMatcher(object):
    '''
    Classify face id using Faiss
    '''

    def __init__(self, matcher_type='IndexFlatIP', *args, **kwargs):
        """ Construct a matcher
        
        Keyword Arguments:
            matcher_type {str} -- type of matcher (default: {'IndexFlatIP'})
        """
        self._labels_arr = None
        self._matcher_type = matcher_type
        self._emb_dimensions = config.DIMENSIONS
        if matcher_type in dir(faiss):
            self._classifier = getattr(faiss, matcher_type)(self._emb_dimensions, *args, **kwargs)
        else:
            self._classifier = faiss.index_factory(self._emb_dimensions, self._matcher_type)
        self._rwlock = rwlock.RWLock()

    def build(self, database):
        """ Build matcher data from database
        
        Arguments:
            database {Database Class}
        """
        labels_embs = database.get_labels_and_embs()
        if not labels_embs is None and len(labels_embs):
            labels, embs = labels_embs
            self.update(embs, labels)

    def update(self, new_embs_ls, new_labels_ls):
        """ Add new embs and labels to current matcher
        Arguments:
            new_embs_ls {np.array, list} -- list of embs
            new_labels_ls {np.array, list} -- [list of label (face id) for each emb
        """
        if len(new_embs_ls):
            try:
                new_embs_ls = np.array(new_embs_ls).reshape((-1, self._emb_dimensions)).astype(np.float32)
                new_labels_ls = np.array(new_labels_ls).reshape((len(new_embs_ls)))
                self._update_array(np.array(new_embs_ls), np.array(new_labels_ls))
            except Exception as err:
                logger.error("Got error update: {} -- error: {}".format(new_embs_ls, err))

    def _update_array(self, new_embs_arr, new_labels_arr):
        try:
            self._rwlock.writer_acquire()
            self._classifier.add(new_embs_arr)
            if self._labels_arr is None:
                self._labels_arr = new_labels_arr
            else:
                self._labels_arr = np.concatenate((self._labels_arr, new_labels_arr))
            self._rwlock.writer_release()
        except Exception as err:
            logger.error("Got error _update_array: {} -- error: {}".format(new_embs_arr, err))

    def match(self, embs, top_matches=5, return_dists=True):
        """ Search nearest vectors from vector pool
        
        Arguments:
            embs {np.array, list} -- embedding to search
        
        Keyword Arguments:
            top_matches {int} -- number of top matches (default: {5})
            return_dists {bool} -- return distances (default: {True})
        
        Returns:
            [type] -- nearest IDs and distances
        """
        if len(embs):
            embs = np.array(embs).reshape((-1, self._emb_dimensions)).astype(np.float32)
            return self._match_array(embs, top_matches, return_dists)

    def _match_array(self, embs, top_matches, return_dists):
        if self._labels_arr is not None:
            top_matches = min(top_matches, len(self._labels_arr))
            dists, inds = self._classifier.search(embs, k=top_matches)
            logger.info("labels is not nil , dits: {} --- inds: {}".format(dists, inds))
            top_match_ids = self._labels_arr[inds]
            
        else:
            logger.info("labels is nil")
            top_match_ids = np.full((embs.shape[0], top_matches), config.NEW_FACE)
            dists = np.full((embs.shape[0], top_matches), -1).astype(np.float32)

        if return_dists:
            return dists, top_match_ids
        return top_match_ids

    def remove(self, labels_arr):
        """ Remove ids from the indexer

        Arguments:
            labels_arr {np.array} -- [array of ids need to be removed]
        """
        if len(labels_arr):
            labels_arr = np.array(labels_arr).reshape(-1)
            self._remove_array(labels_arr)

    def _remove_array(self, labels_arr):
        if self._labels_arr is not None:
            _, match_indices, _ = np.intersect1d(self._labels_arr, labels_arr, return_indices=True)
            inverse_indices = np.setdiff1d(np.arange(len(self._labels_arr)), match_indices)

            self._rwlock.writer_acquire()
            self._labels_arr = self._labels_arr[inverse_indices]
            self._classifier.remove_ids(match_indices)
            self._rwlock.writer_release()