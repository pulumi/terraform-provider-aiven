---
layout: "aiven"
page_title: "Aiven: aiven_project"
description: |-
  Manages an Aiven project resource.
---

# Resource: aiven_project

## Example Usage

```hcl
resource "aiven_project" "myproject" {
    project = "<PROJECT_NAME>"
    card_id = "<FULL_CARD_ID/LAST4_DIGITS>"
}
```

## Argument Reference

The following arguments are supported:

`project` defines the name of the project. Name must be globally unique (between all
Aiven customers) and cannot be changed later without destroying and re-creating the
project, including all sub-resources.

`card_id` is either the full card UUID or the last 4 digits of the card. As the full
UUID is not shown in the UI it is typically easier to use the last 4 digits to identify
the card. This can be omitted if `copy_from_project` is used to copy billing info from
another project.

`copy_from_project` is the name of another project used to copy billing information and
some other project attributes like technical contacts from. This is mostly relevant when
an existing project has billing type set to invoice and that needs to be copied over to a
new project. (Setting billing is otherwise not allowed over the API.) This only has
effect when the project is created.

`ca_cert` is a computed property that can be used to read the CA certificate of the
project. This is required for configuring clients that connect to certain services like
Kafka. This value cannot be set, only read.

## Import

Projects can be imported using their name , e.g.

```
$ terraform import aiven_project.my_project test-project
```
