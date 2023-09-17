#!/bin/sh

echo "### initialize start ####"
CMD_MYSQL="mysql -uuser -ppassword db"
$CMD_MYSQL -e "CREATE TABLE users ( id varchar(50) NOT NULL primary key, name varchar(50) NOT NULL);"
echo "### initialize end   ####"