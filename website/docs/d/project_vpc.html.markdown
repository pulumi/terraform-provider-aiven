---
layout: "aiven"
page_title: "Aiven: aiven_project_vpc"
description: |-
  Gets information on an Aiven project vpc.
---

# Data Source: aiven_project_vpc

## Example Usage

```hcl
data "aiven_project_vpc" "myvpc" {
    project = data.aiven_project.myproject.project
    cloud_name = "google-europe-west1"
}
```

## Argument Reference

`project` defines the project the VPC belongs to.

`cloud_name` defines the cloud provider and region where the vpc resides.e in the Aiven web console.