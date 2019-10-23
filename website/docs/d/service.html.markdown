---
layout: "aiven"
page_title: "Aiven: aiven_service"
description: |-
  Gets information on an Aiven service resource.
---

# Data Source: aiven_service

## Example Usage

```hcl
data "aiven_service" "myservice" {
    project = data.aiven_project.myproject.project
    service_name = "<SERVICE_NAME>"
}
```

## Argument Reference

`project` identifies the project the service belongs to.

`service_name` specifies the actual name of the service.
