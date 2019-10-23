---
layout: "aiven"
page_title: "Aiven: aiven_kafka_acl"
description: |-
  Manages an Aiven Kafka ACL.
---

# Resource: aiven_kafka_acl

## Example Usage

```hcl
resource "aiven_kafka_acl" "mytestacl" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    topic = "<TOPIC_NAME_PATTERN>"
    permission = "admin"
    username = "<USERNAME_PATTERN>"
}
```

## Argument Reference

`project` and `service_name` define the project and service the ACL belongs to.
They should be defined using reference as shown above to set up dependencies correctly.
These properties cannot be changed once the service is created. Doing so will result in
the topic being deleted and new one created instead.

`topic` is a topic name pattern the ACL entry matches to.

`permission` is the level of permission the matching users are given to the matching
topics (admin, read, readwrite, write).

`username` is a username pattern the ACL entry matches to.

## Import

Kafka ACLs can be imported using their ID in the format `<project_name>/<service_name>/<acl_id>` , e.g.

```
$ terraform import aiven_kafka_acl.acl test-project/test-service/aclId
```

**Please Note:** The ACL ID is not directly visible in the Aiven web console.
