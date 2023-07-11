module "host-vpc" {
  count      = var.create_network == true ? 1 : 0
  source     = "../../modules/net-vpc"
  project_id = var.host_project_id
  name       = var.network_name
  vpc_create = var.create_network
  subnets = var.create_subnetwork == true ? [
    {
      name                  = var.subnetwork_name
      region                = var.region
      ip_cidr_range         = var.subnetwork_ip_cidr
      enable_private_access = true
      flow_logs_config = {
        flow_sampling        = 0.5
        aggregation_interval = "INTERVAL_10_MIN"
        metadata             = "INCLUDE_ALL_METADATA"
      }
    },
  ] : []
  shared_vpc_host = true
  shared_vpc_service_projects = [
    var.service_project_id
  ]
  psa_config = {
    ranges = {
      "${var.cloudsql_private_range_name}" = "${var.cloudsql_private_range_cidr}/${var.cloudsql_private_range_prefix_length}"
    }
    export_routes = true
    import_routes = true
  }
  depends_on = [
    module.host_project,
    module.project_services
  ]
}
