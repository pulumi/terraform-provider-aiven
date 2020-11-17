terraform {
  required_providers {
    aiven = {
      source  = "aiven/aiven"
      version = "2.X.X"
    }
  }
}


provider "aiven" {
  api_token = var.timescale_api_token
}

data "aiven_project" "test" {
  project = var.aiven_project_id
}

resource "aiven_pg" "timescaledb" {
  project      = data.aiven_project.test.id
  cloud_name   = "timescale-google-europe-north1"
  plan         = "timescale-basic-100-compute-optimized"
  service_name = "timescaledb"

  pg_user_config {
    pg_version = 12          # Choose an available postgres version
    variant    = "timescale" # Do not forget to specify this variant
  }
}

resource "aiven_grafana" "grafana" {
  project      = data.aiven_project.test.id
  cloud_name   = "timescale-google-europe-north1"
  plan         = "dashboard-1" # Grafana plans are also available in Timescale Cloud
  service_name = "grafana"
}
