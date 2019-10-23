---
layout: "aiven"
page_title: "Aiven: aiven_service_user"
description: |-
  Manages an Aiven service user resource.
---

# Resource: aiven_service_user

## Example Usage

```hcl
resource "aiven_service_user" "myserviceuser" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    username = "<USERNAME>"
}
```

## Argument Reference

`project` and `service_name` define the project and service the user belongs to.
They should be defined using reference as shown above to set up dependencies correctly.

`username` is the actual name of the user account.

None of the service user properties can currently be changed after creation. Doing so
will result in the old database getting dropped and a new database created.

Service users have several computed properties that cannot be set, only read:

`password` is the password of the user (not applicable for all services).

`access_cert` is the access certificate of the user (not applicable for all services).

`access_key` is the access key of the user (not applicable for all services).

`type` tells whether the user is primary account or regular account.

## Import

Service users can be imported using their ID in the format `<project_name>/<service_name>/<username>` , e.g.

```
$ terraform import aiven_service_user.user test-project/test-service/testuser
```
