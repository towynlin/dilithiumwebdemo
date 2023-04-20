provider "kubernetes" {
  host  = "https://${google_container_cluster.primary.endpoint}"
  token = var.gke_token

  client_certificate     = google_container_cluster.primary.master_auth.0.client_certificate
  client_key             = google_container_cluster.primary.master_auth.0.client_key
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}

resource "kubernetes_deployment_v1" "dilithiumwebdemo" {
  metadata {
    name = "dilithiumwebdemo"
    labels = {
      App = "dilithiumwebdemo"
    }
  }

  spec {
    replicas = 4
    selector {
      match_labels = {
        App = "dilithiumwebdemo"
      }
    }
    template {
      metadata {
        labels = {
          App = "dilithiumwebdemo"
        }
      }
      spec {
        container {
          image = "towynlin/dilithiumwebdemo:0.3"
          name  = "dilithiumwebdemo"
          env {
            name  = "DATABASE_URL"
            value = var.database_url
          }

          port {
            container_port = 1323
          }

          #   liveness_probe {
          #     http_get {
          #       path = "/health"
          #       port = 1323
          #     }
          #   }
          #   readiness_probe {
          #     http_get {
          #       path = "/health"
          #       port = 1323
          #     }
          #     failure_threshold = 40
          #   }

          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "50Mi"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "dilithiumwebdemo" {
  metadata {
    name = "dilithiumwebdemo"
    # annotations = {
    #   "cloud.google.com/neg" = jsonencode({ ingress = true })

    # }
  }
  spec {
    selector = {
      App = kubernetes_deployment_v1.dilithiumwebdemo.spec.0.template.0.metadata[0].labels.App
    }
    port {
      port        = 80
      target_port = 1323
    }

    type = "ClusterIP"
  }
  # lifecycle {
  #   ignore_changes = [
  #     metadata[0].annotations,
  #   ]
  # }
}

resource "kubernetes_ingress_v1" "dilithiumwebdemo" {
  wait_for_load_balancer = true
  metadata {
    name = "dilithium-lb"
    annotations = {
      "kubernetes.io/ingress.class" = "nginx"
    }
  }
  spec {
    default_backend {
      service {
        name = kubernetes_service_v1.dilithiumwebdemo.metadata[0].name
        port {
          number = kubernetes_service_v1.dilithiumwebdemo.spec[0].port[0].port
        }
      }
    }
    tls {
      secret_name = "pqcr-tls-secret3"
      hosts = [
        "postquantumcryptography.rocks"
      ]
    }
    rule {
      host = "postquantumcryptography.rocks"

      http {
        path {
          path      = "/"
          path_type = "Prefix"
          backend {
            service {
              name = "dilithiumwebdemo"
              port {
                number = 80
              }
            }
          }
        }
      }
    }
  }
}
