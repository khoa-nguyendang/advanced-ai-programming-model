#!/bin/bash
BASE_PATH=$(dirname $0)

echo "Waiting for mysql to get up"
# Give 10 seconds for master and slave to come up
sleep 10

echo "* Create replication user"

mysql --host $MYSQL_REP1_HOST -uroot -p$MYSQL_SLAVE_PASSWORD  -AN -e 'STOP SLAVE;';
mysql --host $MYSQL_REP1_HOST -uroot -p$MYSQL_MASTER_PASSWORD  -AN -e 'RESET SLAVE ALL;';

mysql --host $APP_MYSQL_HOST -uroot -p$MYSQL_MASTER_PASSWORD  -AN -e "CREATE USER '$MYSQL_REPLICATION_USER'@'%' IDENTIFIED WITH mysql_native_password BY  '$MYSQL_REPLICATION_PASSWORD';"
mysql --host $APP_MYSQL_HOST -uroot -p$MYSQL_MASTER_PASSWORD  -AN -e "GRANT REPLICATION SLAVE ON *.* TO '$MYSQL_REPLICATION_USER'@'%';"
mysql --host $APP_MYSQL_HOST -uroot -p$MYSQL_MASTER_PASSWORD  -AN -e 'flush privileges;'


echo "* Set MySQL01 as master on MySQL02"

MYSQL01_Position=$(eval "mysql --host $APP_MYSQL_HOST -uroot -p$MYSQL_MASTER_PASSWORD  -e 'show master status \G' | grep Position | sed -n -e 's/^.*: //p'")
MYSQL01_File=$(eval "mysql --host $APP_MYSQL_HOST -uroot -p$MYSQL_MASTER_PASSWORD  -e 'show master status \G'     | grep File     | sed -n -e 's/^.*: //p'")
MASTER_IP=$(eval "getent hosts $APP_MYSQL_HOST|awk '{print \$1}'")
echo $MASTER_IP
mysql --host $MYSQL_REP1_HOST -uroot -p$MYSQL_SLAVE_PASSWORD  -AN -e "CHANGE MASTER TO master_host='$APP_MYSQL_HOST', master_port=3306, \
        master_user='$MYSQL_REPLICATION_USER', master_password='$MYSQL_REPLICATION_PASSWORD', master_log_file='$MYSQL01_File', \
        master_log_pos=$MYSQL01_Position;"

echo "* Set MySQL02 as master on MySQL01"

MYSQL02_Position=$(eval "mysql --host $MYSQL_REP1_HOST -uroot -p$MYSQL_SLAVE_PASSWORD  -e 'show master status \G' | grep Position | sed -n -e 's/^.*: //p'")
MYSQL02_File=$(eval "mysql --host $MYSQL_REP1_HOST -uroot -p$MYSQL_SLAVE_PASSWORD  -e 'show master status \G'     | grep File     | sed -n -e 's/^.*: //p'")

SLAVE_IP=$(eval "getent hosts $MYSQL_REP1_HOST|awk '{print \$1}'")
echo $SLAVE_IP

echo "* Start Slave on both Servers"
mysql --host $MYSQL_REP1_HOST -uroot -p$MYSQL_SLAVE_PASSWORD  -AN -e "start slave;"

echo "Increase the max_connections to 2000"
mysql --host $APP_MYSQL_HOST -uroot -p$MYSQL_MASTER_PASSWORD  -AN -e 'set GLOBAL max_connections=2000';
mysql --host $MYSQL_REP1_HOST  -uroot -p$MYSQL_SLAVE_PASSWORD  -AN -e 'set GLOBAL max_connections=2000';

mysql --host $MYSQL_REP1_HOST -uroot -p$MYSQL_MASTER_PASSWORD  -e "show slave status \G"

echo "MySQL servers created!"
echo "--------------------"
echo
echo Variables available fo you :-
echo
echo MYSQL01_IP       : $APP_MYSQL_HOST
echo MYSQL02_IP       : $MYSQL_REP1_HOST