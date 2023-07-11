module "firewall_rules" {
  source       = "../../modules/firewall-rules"
  project_id   = var.host_project_id
  network_name = local.network_name
  rules = [{
    name          = "allow-ssh-ingress"
    description   = "Creates firewall rule targeting tagged instances"
    direction     = "INGRESS"
    source_ranges = ["0.0.0.0/0"]
    target_tags   = var.gce_tags
    allow = [{
      protocol = "tcp"
      ports    = ["22"]
    }]
    log_config = {
      metadata = "INCLUDE_ALL_METADATA"
    }
    },
    {
      name                    = "cloudsql-ingress"
      description             = "Creates firewall rule targeting tagged instances for cloudsql"
      direction               = "INGRESS"
      source_service_accounts = [module.gce_sa.email]
      allow = [{
        protocol = "tcp"
        ports    = ["3306"]
      }]
      log_config = {
        metadata = "INCLUDE_ALL_METADATA"
      }
  }]
}
