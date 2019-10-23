---
layout: "aiven"
page_title: "Aiven: aiven_service_user"
description: |-
  Gets information on an Aiven service user resource.
---

# Data Source: aiven_service_user

## Example Usage

```hcl
data "aiven_service_user" "myserviceuser" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    username = "<USERNAME>"
}
```

## Argument Reference

`project` and `service_name` define the project and service the user belongs to.

`username` is the actual name of the user account.
