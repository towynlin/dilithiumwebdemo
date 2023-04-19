provider "kubernetes" {
  host  = "https://${google_container_cluster.primary.endpoint}"
  token = var.gke_token

  client_certificate     = google_container_cluster.primary.master_auth.0.client_certificate
  client_key             = google_container_cluster.primary.master_auth.0.client_key
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}

resource "kubernetes_deployment" "dilithiumwebdemo" {
  metadata {
    name = "dilithiumwebdemo"
    labels = {
      App = "dilithiumwebdemo"
    }
  }

  spec {
    replicas = 2
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
          image = "towynlin/dilithiumwebdemo:0.2"
          name  = "dilithiumwebdemo"
          env {
            name  = "DATABASE_URL"
            value = var.database_url
          }

          port {
            container_port = 1323
          }

          liveness_probe {
            http_get {
              path = "/health"
              port = 1323
            }
          }
          readiness_probe {
            http_get {
              path = "/health"
              port = 1323
            }
          }

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

resource "kubernetes_service" "dilithiumwebdemo" {
  metadata {
    name        = "dilithiumwebdemo"
    annotations = { "cloud.google.com/neg" = "{\"ingress\": true}" }
  }
  spec {
    selector = {
      App = kubernetes_deployment.dilithiumwebdemo.spec.0.template.0.metadata[0].labels.App
    }
    port {
      port        = 8080
      target_port = 1323
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_ingress_v1" "dilithiumwebdemo" {
  metadata {
    name = "dilithiumwebdemo"
  }
  spec {
    default_backend {
      service {
        name = kubernetes_service.dilithiumwebdemo.metadata[0].name
        port {
          number = kubernetes_service.dilithiumwebdemo.spec[0].port[0].port
        }
      }
    }
    tls {
    }
  }
}
