---
layout: "aiven"
page_title: "Aiven: aiven_service_integration"
description: |-
  Manages an Aiven service integration.
---

# Resource: aiven_service_integration

## Example Usage

```hcl
resource "aiven_service_integration" "myintegration" {
    project = "${aiven_project.myproject.project}"
    destination_endpoint_id = "${aiven_service_integration_endpoint.myendpoint.id}"
    destination_service_name = ""
    integration_type = "datadog"
    source_endpoint_id = ""
    source_service_name = "${aiven_service.testkafka.service_name}"
}
```

## Argument Reference

`project` defines the project the integration belongs to.

`destination_endpoint_id` or `destination_service_name` identifies the target side of
the integration. Only either endpoint identifier or service name must be specified. In
either case the target needs to be defined using the reference syntax described above to
set up the dependency correctly.

`integration_type` identifies the type of integration that is set up. Possible values
include `dashboard`, `datadog`, `logs`, `metrics` and `mirrormaker`.

`source_endpoint_id` or `source_service_name` identifies the source side of the
integration. Only either endpoint identifier or service name must be specified. In either
case the source needs to be defined using the reference syntax described above to set up
the dependency correctly.

`x_user_config` defines integration specific configuration. `x` is the type of the
integration. The available configuration options are documented in
[this JSON file](aiven/templates/integrations_user_config_schema.json). Not all integration
types have any configurable settings.

## Import

Service integrations can be imported using their ID in the format `<project_name>/<integration_id>` , e.g.

```
$ terraform import aiven_service_integration.myinteration test-project/myServiceIntegration
```

**Please Note:** The integration identifier (UUID) is not directly visible in the Aiven web console.
