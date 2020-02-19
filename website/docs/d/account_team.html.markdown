---
layout: "aiven"
page_title: "Aiven: aiven_account_team"
description: |-
  Gets information on an Aiven account team
---

# Data Source: aiven_account_team

## Example Usage

```hcl
data "aiven_account_team" "account_team1" {
    account_id = "${aiven_account.team.account_id}"
    name = "account_team1"
}
```

## Argument Reference

`name` defines an account team name.

`account_id` is an unique account id.
