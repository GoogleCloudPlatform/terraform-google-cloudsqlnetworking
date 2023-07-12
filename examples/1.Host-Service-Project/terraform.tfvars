host_project_id        = "pm-host-networking"
service_project_id     = "pm-service1-networking"
database_version       = "MYSQL_8_0"
cloudsql_instance_name = "sqldb1"
region                 = "us-central1"
zone                   = "us-central1-a"
create_network         = true
create_subnetwork      = true
network_name           = "cloudsql-easy"
subnetwork_name        = "cloudsql-easy-subnet"
subnetwork_ip_cidr     = "10.2.0.0/16"

