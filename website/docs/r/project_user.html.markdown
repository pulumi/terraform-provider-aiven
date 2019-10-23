---
layout: "aiven"
page_title: "Aiven: aiven_project_user"
description: |-
  Manages an Aiven project user.
---

# Resource: aiven_project_user

## Example Usage

```hcl
resource "aiven_project_user" "mytestuser" {
    project = "${aiven_project.myproject.project}"
    email = "john.doe@example.com"
    member_type = "admin"
}
```

## Argument Reference

`project` defines the project the user is a member of.

`email` identifies the email address of the user.

`member_type` defines the access level the user has to the project.

Computed property `accepted` tells whether the user has accepted the request to join
the project; adding user to a project sends an invitation to the target user and the
actual membership is only created once the user accepts the invitation. This property
cannot be set, only read.

## Import

Project users can be imported using their ID in the format `<project_name>/<email>`, e.g.

```
$ terraform import aiven_project_user.testuser test-project/john.doe@example.com
```
