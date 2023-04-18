terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.62.0"
    }
  }
}

provider "google" {
  credentials = file(var.credentials_file)

  project = var.project
  region  = var.region
  zone    = var.zone
}

resource "google_compute_network" "default" {
  name = "dilithium-network"
}

resource "google_compute_instance" "default" {
  name         = "dilithium-instance"
  machine_type = "c3-highcpu-4" // free during preview

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-minimal-2204-jammy-v20230413"
    }
  }

  network_interface {
    network = google_compute_network.default.name
    access_config {
    }
  }
}

resource "google_sql_database" "default" {
  name     = "dilithium-db"
  instance = google_sql_database_instance.default.name
}

resource "google_sql_database_instance" "default" {
  name             = "dilithium-db-instance"
  region           = var.region
  database_version = "POSTGRES_14"
  settings {
    tier = "db-f1-micro"
  }

  deletion_protection = false
}
