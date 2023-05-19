"""
Environment-based configuration

To be loaded using `app.config.from_object` inside `main.py`
"""
import os


class Config(object):
    DEBUG = False
    TESTING = False
    APP_NAME = 'ai-engine-service'
    SERVER_PORT  = '50053'

    # faiss
    FAISS_SERVICE_HOST = '0.0.0.0'
    FAISS_SERVICE_PORT = 50051
    MAX_ENGINE         = '6'

    # sentry
    SENTRY_DSN = None

    ENGINE_HEALTH_TIMEOUT = 2.5

    FACE_MATCHING_CUDA = 1

    MINIMUM_ENROLLMENT_IMAGE = '1'
    GET_MODEL_FROM_CLOUD = 'DISABLE'

    MINIO_ENDPOINT = "localhost:9000"
    MINIO_ENGINE_MODEL_BUCKET = "engine-model"
    MINIO_ROOT_USER = "admin"
    MINIO_ROOT_PASSWORD = "app@2022"

    VERIFICATION_IMAGE_LOG="DISABLED"


class ProductionConfig(Config):
    # The deployment script works better with string for now
    SENTRY_DSN = os.environ.get('SENTRY_DSN')
    SERVER_PORT  = os.environ.get('SERVER_PORT')

    FAISS_SERVICE_HOST = os.environ.get('FAISS_SERVICE_HOST')
    FAISS_SERVICE_PORT = os.environ.get('FAISS_SERVICE_PORT')
    MAX_ENGINE = os.environ.get('MAX_ENGINE')
    MINIMUM_ENROLLMENT_IMAGE = os.environ.get('MINIMUM_ENROLLMENT_IMAGE')
    GET_MODEL_FROM_CLOUD = os.environ.get('GET_MODEL_FROM_CLOUD')
    MINIO_ENDPOINT = os.environ.get('MINIO_ENDPOINT')
    MINIO_ENGINE_MODEL_BUCKET = os.environ.get('MINIO_ENGINE_MODEL_BUCKET')
    MINIO_ROOT_USER = os.environ.get('MINIO_ROOT_USER')
    MINIO_ROOT_PASSWORD = os.environ.get('MINIO_ROOT_PASSWORD')
    VERIFICATION_IMAGE_LOG = os.environ.get('VERIFICATION_IMAGE_LOG')

class DevelopmentConfig(Config):
    DEBUG = True
    SERVER_PORT  = os.environ.get('SERVER_PORT')

    FAISS_SERVICE_HOST = os.environ.get('FAISS_SERVICE_HOST')
    FAISS_SERVICE_PORT = os.environ.get('FAISS_SERVICE_PORT')
    MAX_ENGINE = os.environ.get('MAX_ENGINE')
    MINIMUM_ENROLLMENT_IMAGE = os.environ.get('MINIMUM_ENROLLMENT_IMAGE')
    GET_MODEL_FROM_CLOUD = os.environ.get('GET_MODEL_FROM_CLOUD')
    MINIO_ENDPOINT = os.environ.get('MINIO_ENDPOINT')
    MINIO_ENGINE_MODEL_BUCKET = os.environ.get('MINIO_ENGINE_MODEL_BUCKET')
    MINIO_ROOT_USER = os.environ.get('MINIO_ROOT_USER')
    MINIO_ROOT_PASSWORD = os.environ.get('MINIO_ROOT_PASSWORD')
    VERIFICATION_IMAGE_LOG = os.environ.get('VERIFICATION_IMAGE_LOG')

class TestConfig(Config):
    TESTING = True
    SERVER_PORT  = '50053'

    # FAISS Service mock
    FAISS_SERVICE_HOST = 'localhost'
    FAISS_SERVICE_PORT = '50051'
    MAX_ENGINE = 6
    MINIMUM_ENROLLMENT_IMAGE = 1
    GET_MODEL_FROM_CLOUD = 'DISABLED'
    MINIO_ENDPOINT = 'localhost:9000'
    MINIO_ENGINE_MODEL_BUCKET = '=engine-model'
    MINIO_ROOT_USER = 'admin'
    MINIO_ROOT_PASSWORD = 'app@2022'
    VERIFICATION_IMAGE_LOG = 'ENABLED'

def get_config():
    """
    Return the appropriate configuration object, given the environment
    """
    config_mode = os.environ.get('CLOUD_SERVICE_ENV') or 'Test'
    cls_name = '{0}Config'.format(config_mode.title())
    return globals()[cls_name]()
