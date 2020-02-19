---
layout: "aiven"
page_title: "Aiven: aiven_account"
description: |-
  Gets information on an Aiven account
---

# Data Source: aiven_account

## Example Usage

```hcl
data "aiven_account" "account1" {
    name = "<ACCOUNT_NAME>"
}
```

## Argument Reference

`name` defines an account name.
