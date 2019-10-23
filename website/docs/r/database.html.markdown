---
layout: "aiven"
page_title: "Aiven: aiven_database"
description: |-
  Manages an Aiven database resource.
---

# Resource: aiven_database

## Example Usage

```hcl
resource "aiven_database" "mydatabase" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    database_name = "<DATABASE_NAME>"
}
```

## Argument Reference

`project` and `service_name` define the project and service the database belongs to.
They should be defined using reference as shown above to set up dependencies correctly.

`database_name` is the actual name of the database.

None of the database properties can currently be changed after creation. Doing so will
result in the old database getting dropped and a new database created.

## Import

Database can be imported using their ID in the format `<project_name>/<service_name>/<database_name>` , e.g.

```
$ terraform import aiven_database.my_database test-project/test-service/mydb
```
