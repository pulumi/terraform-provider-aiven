---
layout: "aiven"
page_title: "Aiven: aiven_vpc_peering_connection"
description: |-
  Manages an Aiven vpc peering connection.
---

# Resource: aiven_vpc_peering_connection

## Example Usage

```hcl
resource "aiven_vpc_peering_connection" "mypeeringconnection" {
    vpc_id = "${aiven_project_vpc.myvpc.id}"
    peer_cloud_account = "<PEER_ACCOUNT_ID>"
    peer_vpc = "<PEER_VPC_ID/NAME>"
    peer_region = "<PEER_REGION>"
}
```

## Argument Reference

`vpc_id` is the Aiven VPC the peering connection is associated with.

`peer_cloud_account` defines the identifier of the cloud account the VPC is being
peered with.

`peer_vpc` defines the identifier or name of the remote VPC.

`peer_region` defines the region of the remote VPC if it is not in the same region as Aiven VPC.

Computed property `state` tells the current state of the VPC. This property cannot be
set, only read.

## Import

VPC Peering connections can be imported by using their ID in the format `<project_name>/<VPC_UUID>/<peer_cloud_account>/<peer_vpc>`, e.g.

```
$ terraform import aiven_vpc_peering_connection.myconnection test-project/myVpcId/testAccount/peerVpc
```

**Please Note:** The UUID is not directly visible in the Aiven web console.