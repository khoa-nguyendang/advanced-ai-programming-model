build:
	docker-compose -f docker-compose.base.yaml \
			-f docker-compose.mongodb.yaml \
			-f docker-compose.redis.yaml \
			-f docker-compose.mysql.yaml \
			-f docker-compose.services.yaml build

# Assume all code has been run make prepare-all before
build-web:
	cd ./externals/web-ui/ & docker build -t web-ui .

build-engine:
	cd ./externals/engine-service/ & docker build -t aapi-engine-service .

build-faiss:
	cd ./externals/faiss-service/ & docker build -t aapi-faiss-service .

# useage: make build
build-app:
	docker build .  & docker-compose -f docker-compose.base.yaml \
			-f docker-compose.mongodb.yaml \
			-f docker-compose.redis.yaml \
			-f docker-compose.mysql.yaml \
			-f docker-compose.services.yaml build

# useage: make up
up:
	docker-compose -f docker-compose.base.yaml \
			-f docker-compose.mongodb.yaml \
			-f docker-compose.redis.yaml \
			-f docker-compose.mysql.yaml \
			-f docker-compose.services.yaml up -d



# useage: make down
down:
	docker-compose -f docker-compose.services.yaml \
		-f docker-compose.mysql.yaml \
		-f docker-compose.redis.yaml \
		-f docker-compose.mongodb.yaml \
		-f docker-compose.base.yaml down

# useage: make down
stop:
	docker-compose -f docker-compose.services.yaml \
		-f docker-compose.mysql.yaml \
		-f docker-compose.redis.yaml \
		-f docker-compose.mongodb.yaml \
		-f docker-compose.base.yaml stop

# useage: make logs service=user-service
logs:
	docker-compose -f docker-compose.services.yaml \
		-f docker-compose.mysql.yaml \
		-f docker-compose.redis.yaml \
		-f docker-compose.mongodb.yaml \
		-f docker-compose.base.yaml logs $(service) -f
		
# useage: make exec service=aapimongodb
exec:
	docker-compose -f docker-compose.services.yaml \
		-f docker-compose.mysql.yaml \
		-f docker-compose.redis.yaml \
		-f docker-compose.mongodb.yaml \
		-f docker-compose.base.yaml exec $(service) sh

# useage: make ps
ps:
	docker-compose -f docker-compose.services.yaml \
		-f docker-compose.mysql.yaml \
		-f docker-compose.redis.yaml \
		-f docker-compose.mongodb.yaml \
		-f docker-compose.base.yaml ps

# useage: make deploy token=your_gitlab_token pull=1
# if don't want to pull then remove "pull=1" in the end
deploy:
	git fetch & git pull
	make pull-source u=$(u) p=$(p)
	docker rmi $$(docker images -f "dangling=true" -q)
	make build
	make run

# useage: make all token=your_gitlab_token
# if don't want to pull then remove "pull=1" in the end
all:
	make prepare
	make clone-source u=$(u) p=$(p)
	make pull-source u=$(u) p=$(p)
	make build
	make run

# clear old unused images
clear-images:
	docker image rm $$(docker image ls -qf "dangling=true")

# clear old unused volume
clear-volumes:
	docker volume rm $$(docker volume ls -qf "dangling=true")

# clear old unused network
clear-networks:
	docker network rm $$(docker network ls -qf "dangling=true")

# useage: make prepare-all token=your_gitlab_token only first time
prepare-all:
	make prepare
	sudo mkdir -p ./externals
	sudo mkdir -p ./externals/aapi-faiss-service
	sudo chmod -R 777 ./externals
	sudo chmod -R 777 ./externals/aapi-faiss-service
	sudo mkdir -p ./externals/aapi-engine-service
	sudo chmod -R 777 ./externals/aapi-engine-service


# ==============================================================================
prepare:
	sudo mkdir -p ./volumes/mysqldb
	sudo mkdir -p ./volumes/mysqldb_rep1
	sudo mkdir -p ./volumes/mongodb_data
	sudo mkdir -p ./volumes/mongodb_config
	sudo mkdir -p ./volumes/mongodb_data_replicas_1
	sudo mkdir -p ./volumes/mongodb_config_replicas_1
	sudo mkdir -p ./volumes/mongodb_data_replicas_2
	sudo mkdir -p ./volumes/mongodb_config_replicas_2
	sudo mkdir -p ./volumes/minio
	sudo mkdir -p ./volumes/zookeeper
	sudo mkdir -p ./volumes/kafka
	sudo mkdir -p ./volumes/elasticsearch
	sudo mkdir -p ./volumes/logstash/pipeline
	sudo mkdir -p ./volumes/logstash/config/queries

	sudo chown 999:999 ./volumes/mysqldb
	sudo chown 999:999 ./volumes/mysqldb_rep1
	sudo chown 999:999 ./volumes/mongodb_data
	sudo chown 999:999 ./volumes/mongodb_config
	sudo chown 999:999 ./volumes/mongodb_data_replicas_1
	sudo chown 999:999 ./volumes/mongodb_config_replicas_1
	sudo chown 999:999 ./volumes/mongodb_data_replicas_2
	sudo chown 999:999 ./volumes/mongodb_config_replicas_2
	sudo chown 999:999 ./volumes/minio
	sudo chown 999:999 ./volumes/zookeeper
	sudo chown 999:999 ./volumes/kafka
	sudo chown 999:999 ./volumes/elasticsearch
	sudo chown 999:999 -R ./volumes/logstash
	sudo chown 999:999 ./docker/mongodb
	sudo chmod -R 777 ./volumes/mysqldb
	sudo chmod -R 777 ./volumes/mysqldb_rep1
	sudo chmod -R 777 ./volumes/mongodb_data
	sudo chmod -R 777 ./volumes/mongodb_config
	sudo chmod -R 777 ./volumes/mongodb_data_replicas_1
	sudo chmod -R 777 ./volumes/mongodb_config_replicas_1
	sudo chmod -R 777 ./volumes/mongodb_data_replicas_2
	sudo chmod -R 777 ./volumes/mongodb_config_replicas_2
	sudo chmod -R 777 ./volumes/minio
	sudo chmod -R 777 ./volumes/zookeeper
	sudo chmod -R 777 ./volumes/kafka
	sudo chmod -R 777 ./volumes/mosquitto
	sudo chmod -R 777 ./volumes/elasticsearch
	sudo chmod -R 777 ./volumes/logstash
	sudo chmod -R 777 ./docker/*
	sudo chmod 644 ./docker/mysql_installer/master/master.cnf
	sudo chmod 644 ./docker/mysql_installer/replicas/replica.cnf

clear-mongo:
	sudo rm -rfv ./volumes/mongodb_data/*
	sudo rm -rfv ./volumes/mongodb_config/*
	sudo rm -rfv ./volumes/mongodb_data_replicas_1/*
	sudo rm -rfv ./volumes/mongodb_config_replicas_1/*
	sudo rm -rfv ./volumes/mongodb_data_replicas_2/*
	sudo rm -rfv ./volumes/mongodb_config_replicas_2/*

clear-mysql:
	sudo rm -rfv ./volumes/mysqldb/*
	sudo rm -rfv ./volumes/mysqldb_rep1/*

# ===============================================================================
# Generate all proto
protos-gen:
	./services/user/proto-gen.sh services/user/protos/user/v1/user.proto
	./services/user/proto-gen.sh services/user/protos/faiss/v1/faiss.proto
	./services/user/proto-gen.sh services/user/protos/ai_engine/v1/ai_engine.proto
	./services/logging/proto-gen.sh services/logging/protos/log/v1/log.proto
	./proto-gen.sh protos/v1/
	./proto-gen-gateway.sh
