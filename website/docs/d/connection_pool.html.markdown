---
layout: "aiven"
page_title: "Aiven: aiven_connection_pool"
description: |-
  Gets information on an Aiven connection pool resource.
---

# Data Source: aiven_connection_pool

## Example Usage

```hcl
data "aiven_connection_pool" "mytestpool" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    pool_name = "<POOLNAME>"
}
```

## Argument Reference

`project` and `service_name` define the project and service the connection pool
belongs to.

`pool_name` is the name of the pool.
