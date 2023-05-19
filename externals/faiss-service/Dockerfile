FROM python:3.7

RUN pip install Click==7.0
RUN pip install delorean==1.0.0
RUN pip install numpy
RUN pip install faiss-gpu==1.6.0
RUN pip install filelock==3.0.12
RUN pip install Flask==1.1.1
RUN pip install gunicorn==19.9.0
RUN pip install itsdangerous==1.1.0
RUN pip install Jinja2==2.10.3
RUN pip install MarkupSafe==1.1.1
RUN pip install pytz==2019.3
RUN pip install Werkzeug==0.16.0
RUN pip install pymongo==3.10.1
RUN pip install grpcio-tools protobuf==3.20.*

COPY app/ /app
COPY app/main.py /run.py
WORKDIR /

ENV INDEX_FILE /faiss/faiss.index
ENV INDEX_LOCK_FILE /faiss/faiss.lock
ENV INDEX_STATUS_FILE /faiss/faiss-status.json
ENV NUMBER_OF_WORKERS 1
ENV DATABASE_HOST localhost

CMD python run.py

# build
# docker build -t  aapi/aapi-faiss-service:dev .
# docker run -p 5051:5051 --rm --net aapivsm-demo-net --name aapi-faiss-service aapi/aapi-faiss-service:dev
