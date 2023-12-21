host_project_id = ""
// Use 'null' for the service project if service project topology is not being used.
service_project_id     = null
database_version       = ""
cloudsql_instance_name = "cn-sqlinstance10"
region                 = "us-west2"
zone                   = "us-west2-a"
create_network         = true
create_subnetwork      = true
network_name           = "cloudsql-easy"
subnetwork_name        = "cloudsql-easy-subnet"
subnetwork_ip_cidr     = "10.2.0.0/16"

# Variables for Interconnect
interconnect_project_id  = ""
first_interconnect_name  = ""
second_interconnect_name = ""
ic_router_bgp_asn        = 65009

//first vlan attachment configuration values
first_va_asn       = ""
first_va_bandwidth = ""
first_va_bgp_range = ""
first_vlan_tag     = 609


//second vlan attachment configuration values
second_va_asn       = ""
second_va_bandwidth = ""
second_va_bgp_range = ""
second_vlan_tag     = 609




