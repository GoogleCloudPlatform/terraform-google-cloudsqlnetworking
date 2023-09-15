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

module "mysql" {
  count               = var.create_mysql_db == true ? 1 : 0
  source              = "GoogleCloudPlatform/sql-db/google//modules/mysql"
  version             = "16.0.0"
  name                = "${var.name}-${random_string.mysql_random_suffix[0].result}"
  database_version    = var.database_version
  zone                = var.zone
  project_id          = var.project_id
  ip_configuration    = var.ip_configuration
  deletion_protection = var.deletion_protection
}

resource "random_string" "mysql_random_suffix" {
  count   = var.create_mysql_db == true ? 1 : 0
  length  = 3
  special = false
  upper   = false
}

module "mssql" {
  count               = var.create_mssql_db == true ? 1 : 0
  source              = "GoogleCloudPlatform/sql-db/google//modules/mssql"
  version             = "16.0.0"
  name                = "${var.name}-${random_string.mssql_random_suffix[0].result}"
  database_version    = var.database_version
  zone                = var.zone
  project_id          = var.project_id
  ip_configuration    = var.ip_configuration
  deletion_protection = var.deletion_protection
}
resource "random_string" "mssql_random_suffix" {
  count   = var.create_mssql_db == true ? 1 : 0
  length  = 3
  special = false
  upper   = false
}
module "postgresql" {
  count               = var.create_postgressql_db == true ? 1 : 0
  source              = "GoogleCloudPlatform/sql-db/google//modules/postgresql"
  version             = "16.0.0"
  name                = "${var.name}-${random_string.postgressql_random_suffix[0].result}"
  database_version    = var.database_version
  zone                = var.zone
  project_id          = var.project_id
  ip_configuration    = var.ip_configuration
  deletion_protection = var.deletion_protection
}
resource "random_string" "postgressql_random_suffix" {
  count   = var.create_postgressql_db == true ? 1 : 0
  length  = 3
  special = false
  upper   = false
}
