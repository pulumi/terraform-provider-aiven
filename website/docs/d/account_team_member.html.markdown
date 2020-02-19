---
layout: "aiven"
page_title: "Aiven: aiven_account_team_member"
description: |-
  Gets information on an Aiven account team member
---

# Data Source: aiven_account_team_member

## Example Usage

```hcl
data "aiven_account_team_member" "foo" {
  account_id = "${aiven_account.developers.account_id}"
  team_id = "${aiven_account.developers.account_id}"
  user_email = "user+1@example.com"
}
```

## Argument Reference

`account_id` is an unique account id.

`team_id` is an account team id.

`user_email` is a user email address that first will be invited, and after accepting an invitation, he or she becomes a member of a team.
