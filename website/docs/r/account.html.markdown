---
layout: "aiven"
page_title: "Aiven: aiven_account"
description: |-
  Manages an Aiven account
---

# Resource: aiven_account

## Example Usage

```hcl
resource "aiven_account" "account1" {
    name = "<ACCOUNT_NAME>"
}
```

## Argument Reference

`name` defines an account name.

`account_id` is an auto-generated unique account id.
