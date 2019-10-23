---
layout: "aiven"
page_title: "Aiven: aiven_project_user"
description: |-
  Gets information on an Aiven project user.
---

# Data Source: aiven_project_user

## Example Usage

```hcl
data "aiven_project_user" "mytestuser" {
    project = data.aiven_project.myproject.project
    email = "john.doe@example.com"
}
```

## Argument Reference

`project` defines the project the user is a member of.

`email` identifies the email address of the user.