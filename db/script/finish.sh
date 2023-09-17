#!/bin/sh

CMD_MYSQL="mysql -uuser -ppassword db"
$CMD_MYSQL -e "DROP TABLE users"