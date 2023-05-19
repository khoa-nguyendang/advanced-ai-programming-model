"""
Environment-based configuration

To be loaded using `app.config.from_object` inside `main.py`
"""
import os


class Config(object):
    DEBUG = True
    TESTING = False
    INDEX_FILE = '/tmp/faiss.index'
    INDEX_LOCK_FILE = '/tmp/faiss.lock'
    INDEX_STATUS_FILE = '/tmp/faiss-status.json'
    DIMENSIONS = 512
    NEW_FACE = "NEW-FACE"
    SERVER_PORT  = '50051'
    DATABASE_NAME = 'faiss'
    DATABASE_FACE_INFO = 'face_info'

    DATABASE_HOST = 'localhost'
    DATABASE_PORT = 27017
    DATABASE_USERNAME = 'app_User'
    DATABASE_PASSWORD = 'app_Password1234'
    DATABASE_FACE_INFO = 'face_info'

    MATCH_SCORE_THRESHOLD = 0.65 # model_vectorizer 220208


class ProductionConfig(Config):
    INDEX_FILE = '$INDEX_FILE'
    INDEX_LOCK_FILE = '$INDEX_LOCK_FILE'
    INDEX_STATUS_FILE = '$INDEX_STATUS_FILE'
    DATABASE_HOST = os.environ.get('DATABASE_HOST')
    DATABASE_PORT = int(os.environ.get('DATABASE_PORT'))
    DATABASE_USERNAME = os.environ.get('DATABASE_USERNAME')
    DATABASE_PASSWORD = os.environ.get('DATABASE_PASSWORD')
    DATABASE_FACE_INFO = os.environ.get('DATABASE_FACE_INFO')
    SERVER_PORT  = os.environ.get('SERVER_PORT')
    MATCH_SCORE_THRESHOLD = float(os.environ.get('MATCH_SCORE_THRESHOLD'))

class DevelopmentConfig(Config):
    INDEX_FILE = '$INDEX_FILE'
    INDEX_LOCK_FILE = '$INDEX_LOCK_FILE'
    INDEX_STATUS_FILE = '$INDEX_STATUS_FILE'
    DATABASE_HOST = os.environ.get('DATABASE_HOST')
    DATABASE_PORT = int(os.environ.get('DATABASE_PORT'))
    DATABASE_USERNAME = os.environ.get('DATABASE_USERNAME')
    DATABASE_PASSWORD = os.environ.get('DATABASE_PASSWORD')
    DATABASE_FACE_INFO = os.environ.get('DATABASE_FACE_INFO')
    SERVER_PORT  = os.environ.get('SERVER_PORT')
    MATCH_SCORE_THRESHOLD = float(os.environ.get('MATCH_SCORE_THRESHOLD'))

class LocalhostConfig(Config):
    INDEX_FILE = '$INDEX_FILE'
    INDEX_LOCK_FILE = '$INDEX_LOCK_FILE'
    INDEX_STATUS_FILE = '$INDEX_STATUS_FILE'
    DATABASE_HOST = "localhost"
    DATABASE_PORT = 27017
    DATABASE_USERNAME = 'app_User'
    DATABASE_PASSWORD = 'app_Password1234'
    DATABASE_NAME = 'faiss'
    DATABASE_FACE_INFO = 'face_info'
    SERVER_PORT  = '50051'
    MATCH_SCORE_THRESHOLD = 0.65

class TestConfig(Config):
    DEBUG = True
    TESTING = True
    DATABASE_HOST = 'mongomock://localhost'


def get_config():
    """
    Return the appropriate configuration object, given the environment
    """
    config_mode = os.environ.get('CLOUD_SERVICE_ENV') or 'Localhost'
    cls_name = '{0}Config'.format(config_mode.title())
    return globals()[cls_name]()
