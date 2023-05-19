"""
FAISS index management
"""
import json
import os

# import faiss
from app import matcher
from app.utils import database
from delorean import Delorean
from filelock import FileLock
from flask import current_app
from pytz import utc

_INDEX = None
_INDEX_LAST_UPDATED = None


def index(force_reload=False, *args, **kwargs):
    """
    Index singleton

    If the `force_reload` is set to True, always reload the index
    """
    global _INDEX, _INDEX_LAST_UPDATED
    if _INDEX and not force_reload:
        return _INDEX
    else:
        # lock = FileLock(current_app.config['INDEX_LOCK_FILE'])
        # with lock:
        #     if os.path.isfile(current_app.config['INDEX_FILE']):
        #         _INDEX = faiss.read_index(current_app.config['INDEX_FILE'])
        #         _INDEX.nprobe = 16
        #         _INDEX_LAST_UPDATED = last_updated()
        #     else:
        #         # The index does not exist on disk so let's create it from scratch
        #         d = current_app.config['DIMENSIONS']
        #         quantizer = faiss.IndexFlatL2(d)
        #         _INDEX = faiss.IndexIVFFlat(quantizer, d, 1, faiss.METRIC_INNER_PRODUCT)
        #         faiss.write_index(_INDEX, current_app.config['INDEX_FILE'])
        #         _last_updated = record_last_updated()
        #         _INDEX_LAST_UPDATED = _last_updated
        _database = kwargs.get('database', database.FaceInfoDatabase())
        _INDEX = matcher.FaissMatcher("IndexFlatIP") #cosine similarity , that runs from [0-1], required to normalize vector before append
        _INDEX.build(_database)
        return _INDEX


# def write_index_to_disk():
#     """
#     Write the index to disk with a lock
#     """
#     lock = FileLock(current_app.config['INDEX_LOCK_FILE'])
#     with lock:
#         faiss.write_index(_INDEX, current_app.config['INDEX_FILE'])
#         record_last_updated()


# def search_faiss(vector):
#     """
#     Perform a search within the FAISS index with lock
#     """
#     # if _INDEX_LAST_UPDATED:
#     #     lock = FileLock(current_app.config['INDEX_LOCK_FILE'])
#     #     with lock:
#     #         force_reload = _INDEX_LAST_UPDATED < last_updated()
#     # result = index(force_reload).search(vector, 5)
#     result = index().match(vector, 5)
#     return result


# def record_last_updated():
#     """
#     Record the last update timestamp for the index and return it

#     It should always be called after the index is written to on disk
#     """
#     _last_updated = int(Delorean(timezone=utc).epoch)
#     status = {
#         'last_updated': _last_updated
#     }
#     with open(current_app.config['INDEX_STATUS_FILE'], 'w') as f:
#         f.write(json.dumps(status))

#     return _last_updated


# def last_updated():
#     """
#     Fetch the last updated timestamp from the status file
#     """
#     if not os.path.isfile(current_app.config['INDEX_STATUS_FILE']):
#         # If the file does not exist, create it
#         return record_last_updated()

#     with open(current_app.config['INDEX_STATUS_FILE'], 'r') as f:
#         status = json.load(f)
#         return int(status.get('last_updated'))
