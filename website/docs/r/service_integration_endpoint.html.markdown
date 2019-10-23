---
layout: "aiven"
page_title: "Aiven: aiven_service_integration_endpoint"
description: |-
  Manages an Aiven service integration endpoint.
---

# Resource: aiven_service_integration_endpoint

## Example Usage

```hcl
resource "aiven_service_integration_endpoint" "myendpoint" {
    project = "${aiven_project.myproject.project}"
    endpoint_name = "<ENDPOINT_NAME>"
    endpoint_type = "datadog"
    datadog_user_config {
        datadog_api_key = "<DATADOG_API_KEY>"
    }
}
```

## Argument Reference

`project` defines the project the endpoint is associated with.

`endpoint_name` is the name of the endpoint. This value has no effect beyond being used
to identify different integration endpoints.

`endpoint_type` is the type of the external service this endpoint is associated with.
By the time of writing the only available option is `datadog`.

`x_user_config` defines endpoint type specific configuration. `x` is the type of the
endpoint. The available configuration options are documented in
[this JSON file](aiven/templates/integration_endpoints_user_config_schema.json)

## Import

Service integration endpoints can be imported using their ID in the format `<project_name>/<endpoint_id>` , e.g.

```
$ terraform import aiven_service_integration_endpoint.myendpoint test-project/myServiceIntegrationEndpoint
```

**Please Note:** The endpoint identifier (UUID) is not directly visible in the Aiven web console.
