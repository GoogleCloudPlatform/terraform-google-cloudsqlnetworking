###############################################################################
#########################  GCE for the Host Project ###########################
###############################################################################

module "google_compute_instance" {
  source               = "../../modules/computeinstance"
  project_id           = var.service_project_id
  subnetwork_id        = local.subnetwork_id
  vm_service_account   = local.vm_service_account
  region               = var.region
  zone                 = var.zone
  subnetwork_project   = var.service_project_id
  gce_tags             = var.gce_tags
  source_image         = var.source_image
  source_image_project = var.source_image_project
  source_image_family  = var.source_image_family
  target_size          = var.target_size
  deletion_protection  = var.deletion_protection
  startup_script       = <<-EOT
                          #!/bin/sh
                          echo " ====== setting up the my sql cient ===== "
                          sudo apt-get -y update
                          sudo apt-get -y install mariadb-client-10.6

                          mysql --version

                          echo " ====== setting up the cloud sql proxy ===== "
                          curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.4.0/cloud-sql-proxy.linux.amd64
                          chmod +x cloud-sql-proxy

                          ./cloud-sql-proxy --version

                          echo " ====== Startup Script Execution Complete ===== "
                      EOT
  metadata = {
    "enable-oslogin" : true
  }
  access_config = [{
    nat_ip       = "",
    network_tier = var.network_tier
  }]
  depends_on = [
    module.host-vpc
  ]
}


###############################################################################
#########################  GCE for the User Project ###########################
###############################################################################

module "user_google_compute_instance" {
  source               = "../../modules/computeinstance"
  project_id           = var.user_project_id
  subnetwork_id        = local.uservpc_subnetwork_id
  vm_service_account   = local.user_vm_service_account
  region               = var.user_region
  zone                 = var.zone
  subnetwork_project   = var.user_project_id
  gce_tags             = var.gce_tags
  source_image         = var.source_image
  source_image_project = var.source_image_project
  source_image_family  = var.source_image_family
  target_size          = var.target_size
  deletion_protection  = false
  startup_script       = <<-EOT
                          #!/bin/sh
                          echo " ====== setting up the my sql cient ===== "
                          sudo apt-get -y update
                          sudo apt-get -y install mariadb-client-10.6

                          mysql --version

                          echo " ====== setting up the cloud sql proxy ===== "
                          curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.4.0/cloud-sql-proxy.linux.amd64
                          chmod +x cloud-sql-proxy

                          ./cloud-sql-proxy --version

                          echo " ====== Startup Script Execution Complete ===== "
                      EOT
  metadata = {
    "enable-oslogin" : true
  }
  access_config = [{
    nat_ip       = "",
    network_tier = var.network_tier
  }]
  depends_on = [
    module.user-vpc
  ]
}

