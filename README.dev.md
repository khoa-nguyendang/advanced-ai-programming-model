# Troubleshooting guide lines

clear old data of mysql

```
sudo rm -R ./volumes/mysqldb/*

```

clear old data of mongodb
```
make clear-mongo
```

```
# access to mongodb primary
docker logs -n 20 aapi-golang-microservices_mongodb_installer_1
docker logs -n 20 aapimongodb
docker logs -n 20 mongo_rep1
mongo --host localhost --port 27017 -u app_User -p app_Password1234 --authenticationDatabase faiss
mongo --host localhost --port 27018 -u app_User -p app_Password1234 --authenticationDatabase faiss
# access to mongodb services
docker exec -it aapimongodb mongo
docker exec -it mongo-rep1 mongo

rs.initiate(
    {
      _id : 'rs0',
      members: [
        { _id : 0, host : "localhost:27017" },
        { _id : 1, host : "localhost:27018" },
        { _id : 2, host : "localhost:27019", arbiterOnly: true }
      ]
    }
  )
```

### update config of mongoDB

[Sync-Target](https://docs.mongodb.com/manual/tutorial/configure-replica-set-secondary-sync-target/)
```
rs.syncFrom("localhost:27017");
```



### First Deloyment tutorial

Step 1. Clear old data of mysql and mongodb
```
make clear-mongodb
make clear-mysql
```

Step 2. Prepare new necessary folder
```
make prepare
```

Step 3. Pull latest source
```
make pull-source
```

Step 4. Build docker-compose
```
make build
```

Step 5. Start docker-compose
```
make up
```

Step 6. View process
```
clear
docker ps
```

#### Stop current process
```
make down
```

#### display logs for specific service
```
make logs service=name_of_service_here
#example : make logs service=user-service
```

#### access to specific service
```
make exec service=name_of_service_here
#example : make exec service=user-service
```


test mongodb
```
rs.status();
docker logs -n 20 aapi-golang-microservices_mongodb_installer_1
docker logs -n 20 aapimongodb
docker logs -n 20 mongo-rep1
mongo --host localhost --port 27017 -u app_User -p app_Password1234 --authenticationDatabase faiss
mongo --host localhost --port 27018 -u app_User -p app_Password1234 --authenticationDatabase faiss

```

troubleshooting Mongodb cannot create user for faiss

step 1. install mongodb-shell.
```
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.0.list
sudo apt-get update
sudo apt-get install -y mongodb-org-shell
```

step 2. acccess master with mongodb client
```

```

#### Access to MySql
```
 docker-compose exec aapimysql mysql -uroot -papp_Password
 docker-compose exec aapimysql mysql -uroot -papp_Password -e "CREATE DATABASE test_replication;"
 docker-compose exec mysql-rep1 mysql -uroot -papp_Password
 docker-compose exec mysql-rep1 mysql -uroot -papp_Password -e "SHOW DATABASES;"



 make logs service=mysql_installer
 $(eval "mysql --host 0.0.0.0 --port 3307 -uroot -papp_Password -e 'show master status \G' | grep File | sed -n -e 's/^: //p'")
 
```



install necessary Go gateway packages
```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

```

example running gateway


example post request
```
curl -X 'POST' \
  'http://localhost:8081/api/v1/authentication/administrator-authorize' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "username": "admin",
  "password": "admin",
  "callbackUrl": "string"
}'
```

build local then add
```
 -buildvcs=false
 ```