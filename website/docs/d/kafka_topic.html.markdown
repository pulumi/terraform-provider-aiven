---
layout: "aiven"
page_title: "Aiven: aiven_kafka_topic"
description: |-
  Gets information on an Aiven Kafka Topic.
---

# Data Source: aiven_kafka_topic

## Example Usage

```hcl
data "aiven_kafka_topic" "mytesttopic" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    topic_name = "<TOPIC_NAME>"
}
```

## Argument Reference

`project` and `service_name` define the project and service the topic belongs to.

`topic_name` is the actual name of the topic account.
