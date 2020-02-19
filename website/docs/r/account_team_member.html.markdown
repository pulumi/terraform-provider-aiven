---
layout: "aiven"
page_title: "Aiven: aiven_account_team_member"
description: |-
  Manages an Aiven account team member
---

# Resource: aiven_account_team_member

During the creation of `aiven_account_team_member` resource, an email invitation will be sent
to a user using `user_email` address. If the user accepts an invitation, he or she will become a member of the account team. 
The deletion of `aiven_account_team_member` will not only delete invitation if one was sent but not yet accepted by the 
user, and it will also eliminate an account team member if one has accepted an invitation previously.

## Example Usage

```hcl
resource "aiven_account_team_member" "foo" {
  account_id = "${aiven_account.developers.account_id}"
  team_id = "${aiven_account.developers.account_id}"
  user_email = "user+1@example.com"
}
```

## Argument Reference

`account_id` is an unique account id.

`team_id` is an account team id.

`user_email` is a user email address that first will be invited, and after accepting an invitation, he or she becomes a member of a team.

`accepted` is a boolean flag that determines whether an invitation was accepted or not by the user. `false` value means that the 
invitation was sent to the user but not yet accepted. `true` means that the user accepted the invitation and now a member of an account team
