---
layout: "aiven"
page_title: "Aiven: aiven_kafka_acl"
description: |-
  Gets information on an Aiven Kafka ACL.
---

# Data Source: aiven_kafka_acl

## Example Usage

```hcl
data "aiven_kafka_acl" "mytestacl" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    topic = "<TOPIC_NAME_PATTERN>"
    username = "<USERNAME_PATTERN>"
}
```

## Argument Reference

`project` and `service_name` define the project and service the ACL belongs to.

`topic` is a topic name pattern the ACL entry matches to.

`username` is a username pattern the ACL entry matches to.
