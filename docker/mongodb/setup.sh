#!/bin/bash

#MONGODB1=`ping -c 1 aapimongodb | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`
#MONGODB2=`ping -c 1 mongo-rep1 | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`
#MONGODB3=`ping -c 1 mongo-rep2 | head -1  | cut -d "(" -f 2 | cut -d ")" -f 1`

MONGODB1=aapimongodb
MONGODB2=mongo-rep1
UserName=$MONGO_INITDB_ROOT_USERNAME
Password=$MONGO_INITDB_ROOT_PASSWORD
Database=$MONGO_INITDB_DATABASE

echo "********MONGODB1 = " ${MONGODB1}
echo "********MONGODB2 = " ${MONGODB2}
echo "********UserName = " ${UserName}
echo "********Password = " ${Password}
echo "********Database = " ${Database}

echo "********Waiting for startup..********"
until curl http://${MONGODB1}:27017/serverStatus\?text\=1 2>&1 | grep uptime | head -1; do
  printf '.'
done

echo curl http://${MONGODB1}:27017/serverStatus\?text\=1 2>&1 | grep uptime | head -1
echo "********Started..********"

echo SETUP.sh time now: `date +"%T" `
mongo --host ${MONGODB1} --port 27017 <<EOF
use faiss;
var cfg = {
    "_id": "rs0",
    "protocolVersion": 1,
    "version": 1,
    "members": [
        {
            "_id": 0,
            "host": "${MONGODB1}:27017",
            "priority": 1
        },
        {
            "_id": 1,
            "host": "${MONGODB2}:27017",
            "priority": 0,
        }
    ],settings: {chainingAllowed: true}
};
rs.initiate(cfg, { force: true });
rs.status();

EOF


mongo --host ${MONGODB1} --port 27017  <<EOF

use admin;
db.createUser( { user: 'app_User', pwd: 'app_Password1234', roles: [ { role: 'readWrite', db: 'admin' }]});
use faiss;
db.createUser( { user: 'app_User', pwd: 'app_Password1234', roles: [ { role: 'readWrite', db: 'faiss' }]});
EOF