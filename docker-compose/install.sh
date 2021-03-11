#!/bin/bash

CONTAINER=sqle-1
MYSQL_CONTAINER=mysql-for-sqle

SQLE_RPM=sqle.standard.qa.el7.x86_64.rpm

# create default database for sqle forever
for i in $(seq 1 10)
do
  sleep 5
  docker exec -it ${MYSQL_CONTAINER} bash -c\
 "mysql -uroot -pmysqlpass -e \"CREATE DATABASE IF NOT EXISTS sqle default character set utf8mb4 collate utf8mb4_unicode_ci \""
  if [[ $? != 0 ]]; then
    echo "create database failed, retry..."
    echo
  else
    echo "create database success"
    break
  fi
done

# install sqle
docker cp ${SQLE_RPM} ${CONTAINER}:/opt
docker exec -it ${CONTAINER} bash -c "rpm -ivh --prefix=/opt/sqle /opt/${SQLE_RPM}"
# update config of sqle
docker cp ./docker-compose/sqled.yml ${CONTAINER}:/opt/sqle/etc/
# start sqled
echo "start sqled"
docker exec -it ${CONTAINER} bash -c "systemctl start sqled"


