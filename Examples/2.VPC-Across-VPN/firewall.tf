###############################################################################
#########################  Firewall for the Host Project ######################
###############################################################################

resource "google_compute_firewall" "rules" {
  name        = "cloudsql-firewall-rule"
  project     = var.host_project_id
  network     = local.network_id
  description = "Creates firewall rule targeting tagged instances for cloudsql"
  allow {
    protocol = "tcp"
    ports    = ["3306"]
  }
  log_config {
    metadata = "INCLUDE_ALL_METADATA"
  }
  source_service_accounts = [
    module.gce_sa.email,
    module.user_gce_sa.email
  ]
}

resource "google_compute_firewall" "enable_ssh" {
  name        = "ssh-firewall-rule"
  project     = var.host_project_id
  network     = local.network_id
  description = "Creates firewall rule targeting tagged instances"
  allow {
    protocol = "tcp"
    ports    = ["80", "443", "22"]
  }
  allow {
    protocol = "icmp"
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags   = var.gce_tags
}

###############################################################################
#########################  Firewall for the User Project ######################
###############################################################################

resource "google_compute_firewall" "uservpc_enable_ssh" {
  name        = "ssh-firewall-rule"
  project     = var.user_project_id
  network     = local.uservpc_network_id
  description = "Creates firewall rule targeting tagged instances"
  allow {
    protocol = "tcp"
    ports    = ["80", "443", "22"]
  }
  allow {
    protocol = "icmp"
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags   = var.gce_tags
}
