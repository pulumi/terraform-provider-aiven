---
layout: "aiven"
page_title: "Aiven: aiven_project_vpc"
description: |-
  Manages an Aiven project vpc.
---

# Resource: aiven_project_vpc

## Example Usage

```hcl
resource "aiven_project_vpc" "myvpc" {
    project = "${aiven_project.myproject.project}"
    cloud_name = "google-europe-west1"
    network_cidr = "192.168.0.1/24"
}
```

## Argument Reference

`project` defines the project the VPC belongs to.

`cloud_name` defines where the cloud provider and region where the service is hosted
in. See the Service resource for additional information.

`network_cidr` defines the network CIDR of the VPC.

Computed property `state` tells the current state of the VPC. This property cannot be
set, only read.

## Import

Project VPCs can be imported using their ID in the format `<project_name>/<VPC_UUID>`, e.g.

```
$ terraform import aiven_project_vpc.myvpc test-project/myVpcId
```

**Please Note:** The UUID is not directly visible in the Aiven web console.