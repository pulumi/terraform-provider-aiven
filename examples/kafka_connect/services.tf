# Kafka service
resource "aiven_kafka" "kafka-service" {
  project = aiven_project.kafka-con-project1.project
  cloud_name = "google-europe-west1"
  plan = "business-4"
  service_name = "kafka-service1"
  maintenance_window_dow = "monday"
  maintenance_window_time = "10:00:00"

  kafka_user_config {
    kafka_version = "2.4"
  }
}

# Kafka connect service
resource "aiven_kafka_connect" "kafka_connect" {
  project = aiven_project.kafka-con-project1.project
  cloud_name = "google-europe-west1"
  plan = "startup-4"
  service_name = "kafka-connect1"
  maintenance_window_dow = "monday"
  maintenance_window_time = "10:00:00"

  kafka_connect_user_config {
    kafka_connect {
      consumer_isolation_level = "read_committed"
    }

    public_access {
      kafka_connect = true
    }
  }
}

// Kafka connect service integration
resource "aiven_service_integration" "i1" {
  project = aiven_project.kafka-con-project1.project
  integration_type = "kafka_connect"
  source_service_name = aiven_kafka.kafka-service.service_name
  destination_service_name = aiven_kafka_connect.kafka_connect.service_name

  kafka_connect_user_config {
    kafka_connect {
      group_id = "connect"
      status_storage_topic = "__connect_status"
      offset_storage_topic = "__connect_offsets"
    }
  }
}