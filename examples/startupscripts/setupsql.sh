#!/bin/sh
echo " ====== setting up the my sql cient ===== "
sudo apt-get -y update
sudo apt-get -y install mariadb-client-10.6

mysql --version

echo " ====== setting up the cloud sql proxy ===== "
curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.4.0/cloud-sql-proxy.linux.amd64
chmod +x cloud-sql-proxy

./cloud-sql-proxy --version

echo " ====== Startup Script Execution Complete ===== "

#Creating database inside mysql instance using the private IP
mysql --host=${host_ip} --user=${default_username} --password=${default_password} -e "create database ${database_name};"
