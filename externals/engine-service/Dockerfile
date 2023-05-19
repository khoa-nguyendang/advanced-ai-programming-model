
FROM nvcr.io/nvidia/tensorrt:22.02-py3

RUN pip install torch torchvision opencv-python
RUN pip install opencv-contrib-python-headless scikit-image
RUN pip install grpcio-tools protobuf==3.20.*
RUN pip install minio==7.1.1

COPY ./ /aapi-engine-service
WORKDIR /aapi-engine-service

ENV FAISS_SERVICE_HOST 0.0.0.0
ENV FAISS_SERVICE_PORT 50051

CMD python run.py
#docker build -t aapi/aapi-engine-service:dev .
#docker run --gpus all --rm --net host --name aapi-engine-service aapi/aapi-engine-service:dev