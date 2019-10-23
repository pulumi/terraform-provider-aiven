---
layout: "aiven"
page_title: "Aiven: aiven_database"
description: |-
  Gets information on an Aiven database resource.
---

# Data Source: aiven_database

## Example Usage

```hcl
data "aiven_database" "mydatabase" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    database_name = "<DATABASE_NAME>"
}
```

## Argument Reference

`project` and `service_name` define the project and service the database belongs to.

`database_name` is the actual name of the database.