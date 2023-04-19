terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.62.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.19.0"
    }
  }

  required_version = ">= 1.4.5"
}
