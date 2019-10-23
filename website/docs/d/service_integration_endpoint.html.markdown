---
layout: "aiven"
page_title: "Aiven: aiven_service_integration_endpoint"
description: |-
  Gets integration on an Aiven service integration endpoint.
---

# Data Source: aiven_service_integration_endpoint

## Example Usage

```hcl
data "aiven_service_integration_endpoint" "myendpoint" {
    project = data.aiven_project.myproject.project
    endpoint_name = "<ENDPOINT_NAME>"
}
```

## Argument Reference

`project` defines the project the endpoint is associated with.

`endpoint_name` is the name of the endpoint.