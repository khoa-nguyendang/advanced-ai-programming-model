import io
import logging


def upload_file(minio_client, value_as_bytes, file_path) :
    try:
        value_as_a_stream = io.BytesIO(value_as_bytes)
        minio_client.put_object("user-service-verification", file_path, value_as_a_stream , length=len(value_as_bytes))
    except Exception as e:
        logging.error(e)

