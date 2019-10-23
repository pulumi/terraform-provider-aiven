---
layout: "aiven"
page_title: "Aiven: aiven_vpc_peering_connection"
description: |-
  Gets information on an Aiven vpc peering connection.
---

# Data Source: aiven_vpc_peering_connection

## Example Usage

```hcl
data "aiven_vpc_peering_connection" "mypeeringconnection" {
    vpc_id = data.aiven_project_vpc.vpc_id
    peer_cloud_account = "<PEER_ACCOUNT_ID>"
    peer_vpc = "<PEER_VPC_ID/NAME>"
}
```

## Argument Reference

`vpc_id` is the Aiven VPC the peering connection is associated with.

`peer_cloud_account` defines the identifier of the cloud account the VPC has been
peered with.

`peer_vpc` defines the identifier or name of the remote VPC.