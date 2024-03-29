#!/bin/sh
# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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