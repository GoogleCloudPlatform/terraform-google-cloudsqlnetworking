#!/bin/sh
# Copyright 2023-2024 Google LLC
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

sudo sysctl -w net.ipv4.ip_forward=1

sudo iptables -t nat -A POSTROUTING -o ens4 -j MASQUERADE

sudo apt-get update

sudo apt-get -y install sshpass

#Setting nameserver resolution
sshpass -p 'basebase' ssh -o StrictHostKeyChecking=no base@${bastion_ip} "echo 'basebase' | sudo -S sh -c 'echo nameserver 8.8.8.8 >> /etc/resolv.conf'"

sshpass -p 'basebase' ssh -o StrictHostKeyChecking=no base@${bastion_ip} "echo 'basebase' | sudo -S sh -c 'echo nameserver 8.8.4.4 >> /etc/resolv.conf'"

#Installing mysql client
sshpass -p 'basebase' ssh -o StrictHostKeyChecking=no base@${bastion_ip} "echo 'basebase' | sudo -S apt-get update && echo 'basebase' | sudo -S apt-get -y install default-mysql-client"

#Creating database
sshpass -p 'basebase' ssh -o StrictHostKeyChecking=no base@${bastion_ip} "mysql --host=${host_ip} --user=${default_username} --password=${default_password} -e 'create database ${database_name};'"

