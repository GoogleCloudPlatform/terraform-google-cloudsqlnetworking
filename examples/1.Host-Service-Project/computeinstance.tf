

module "instance_template" {
  source               = "terraform-google-modules/vm/google//modules/instance_template"
  project_id           = var.service_project_id
  subnetwork           = local.subnetwork_id
  service_account      = local.vm_service_account
  subnetwork_project   = var.service_project_id
  tags                 = var.gce_tags
  source_image_project = "ubuntu-os-cloud"
  source_image_family  = "ubuntu-2204-lts"
  metadata             = { "enable-oslogin" : true }
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
  access_config = [{
    nat_ip       = "",
    network_tier = var.network_tier
  }]
  depends_on = [
    google_compute_shared_vpc_service_project.service1
  ]
}

module "compute_instance" {
  source              = "terraform-google-modules/vm/google//modules/compute_instance"
  region              = var.region
  zone                = var.zone
  num_instances       = var.target_size
  instance_template   = module.instance_template.self_link
  deletion_protection = false
  depends_on = [
    google_compute_shared_vpc_service_project.service1
  ]
}
