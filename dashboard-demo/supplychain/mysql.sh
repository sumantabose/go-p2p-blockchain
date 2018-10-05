#!/bin/bash

MYSQL_ROOT=""
MYSQL_PASS=""

option=$1

if [ "$option" == "-start" ]; then
	echo "Starting MySQL ..."
	sudo /usr/local/mysql/support-files/mysql.server start

elif [ "$option" == "-stop" ]; then
	echo "Stopping MySQL ..."
	sudo /usr/local/mysql/support-files/mysql.server stop

elif [ "$option" == "-restart" ]; then
	echo "Restarting MySQL ..."
	sudo /usr/local/mysql/support-files/mysql.server restart

elif [ "$option" == "-login" ]; then
	echo "Logging into MySQL ..."
	mysql -u "$MYSQL_ROOT" -p"$MYSQL_PASS"

elif [ "$option" == "-setup" ]; then
	echo "Setting up Blockchain database ..."
	mysql -u "$MYSQL_ROOT" -p"$MYSQL_PASS" < setup.sql

elif [ "$option" == "-destroy" ]; then
	echo "Destroying Blockchain database ..."
	mysql -u "$MYSQL_ROOT" -p"$MYSQL_PASS" < destroy.sql

elif [ "$option" == "-reset" ]; then
	echo "Destroying and Setting up Blockchain database ..."
	mysql -u "$MYSQL_ROOT" -p"$MYSQL_PASS" < destroy.sql
	mysql -u "$MYSQL_ROOT" -p"$MYSQL_PASS" < setup.sql

elif [ "$option" == "-help" ]; then
	echo -e "\nHelp guide to MySQL Blockchain database ..."
	echo "Usage: bash mysql.sh -[option]"
	echo "List of options:"
	echo -e "\t1. start: To start the MySQL database"
	echo -e "\t2. stop: To stop the MySQL database"
	echo -e "\t3. restart: To restart the MySQL database"
	echo -e "\t4. bcsetup: To setup the database"
	echo -e "\t5. bcdestroy: To destroy the database"
	echo -e "\t6. bcresetup: To destroy and setup the database"
	echo -e "\t7. help: To see help guide"

else
	echo "Invalid Option ..."

fi
