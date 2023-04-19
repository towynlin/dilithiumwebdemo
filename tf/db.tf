resource "google_sql_database" "default" {
  name     = "${var.project_name}-db"
  instance = google_sql_database_instance.default.name
}

resource "google_sql_database_instance" "default" {
  name             = "${var.project_name}-db-instance"
  region           = var.region
  database_version = "POSTGRES_14"
  settings {
    tier = "db-f1-micro"
  }
}
