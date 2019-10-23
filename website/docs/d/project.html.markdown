---
layout: "aiven"
page_title: "Aiven: aiven_project"
description: |-
  Gets information on an Aiven project resource.
---

# Data Source: aiven_project

## Example Usage

```hcl
data "aiven_project" "myproject" {
    project = "<PROJECT_NAME>"
}
```

## Argument Reference

`project` defines the name of the project.