
resource "google_compute_instance" "compute_instance" {
  name                = var.compute_instance_name
  project             = var.project_id
  machine_type        = var.machine_type
  zone                = var.zone
  tags                = var.gce_tags
  deletion_protection = var.deletion_protection
  boot_disk {
    initialize_params {
      image = "${var.source_image_project}/${var.source_image_family}"
    }
  }

  network_interface {
    subnetwork = var.subnetwork_id
    access_config {
      network_tier = var.network_tier
    }
  }

  metadata                = var.metadata
  metadata_startup_script = var.startup_script
  service_account {
    email  = var.vm_service_account.email
    scopes = var.vm_service_account.scopes
  }
}
