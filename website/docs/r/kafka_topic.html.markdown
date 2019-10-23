---
layout: "aiven"
page_title: "Aiven: aiven_kafka_topic"
description: |-
  Manages an Aiven Kafka Topic.
---

# Resource: aiven_kafka_topic

## Example Usage

```hcl
resource "aiven_kafka_topic" "mytesttopic" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    topic_name = "<TOPIC_NAME>"
    partitions = 5
    replication = 3
    retention_bytes = -1
    retention_hours = 72
    minimum_in_sync_replicas = 2
    cleanup_policy = "delete"
}
```

## Argument Reference

`project` and `service_name` define the project and service the topic belongs to.
They should be defined using reference as shown above to set up dependencies correctly.
These properties cannot be changed once the service is created. Doing so will result in
the topic being deleted and new one created instead.

`topic_name` is the actual name of the topic account. This propery cannot be changed
once the service is created. Doing so will result in the topic being deleted and new one
created instead.

Other properties should be self-explanatory. They can be changed after the topic has been
created.

## Import

Kafka topics can be imported using their ID in the format `<project_name>/<service_name>/<topic_name>` , e.g.

```
$ terraform import aiven_kafka_topic.topic test-project/test-service/billing-topic
```
