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

module "sql_db" {
  source              = "../../modules/cloudsql"
  name                = var.cloudsql_instance_name
  database_version    = var.database_version
  zone                = var.zone
  region              = var.region
  project_id          = var.producer_project_id
  ip_configuration    = local.ip_configuration
  deletion_protection = var.deletion_protection
}
