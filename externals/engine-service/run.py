import logging
import app.main.service.image_vector_proxy as Server

logging.basicConfig(level=logging.INFO)

if __name__ == '__main__':
    Server.run()