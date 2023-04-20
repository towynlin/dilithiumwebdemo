resource "google_sql_database" "default" {
  name     = "${var.project_name}-db"
  instance = google_sql_database_instance.default.name
}

resource "google_sql_database_instance" "default" {
  provider         = google-beta
  name             = "${var.project_name}-db-instance"
  region           = var.region
  database_version = "POSTGRES_14"
  # depends_on       = [google_service_networking_connection.private_vpc_connection]
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled                                  = true
      private_network                               = google_compute_network.private_network.id
      enable_private_path_for_google_cloud_services = true
    }
  }
}

resource "google_compute_network" "private_network" {
  provider = google-beta

  name = "private-network"
}

resource "google_compute_global_address" "private_ip_address" {
  provider = google-beta

  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.private_network.id
}

# resource "google_service_networking_connection" "private_vpc_connection" {
#   provider = google-beta

#   network                 = google_compute_network.private_network.id
#   service                 = "servicenetworking.googleapis.com"
#   reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
# }

resource "random_id" "db_name_suffix" {
  byte_length = 4
}
