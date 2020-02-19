---
layout: "aiven"
page_title: "Aiven: aiven_account_team_project"
description: |-
  Gets information on an Aiven account team project
---

# Data Source: aiven_account_team_project

## Example Usage

```hcl
data "aiven_account_team_project" "account_team_project1" {
    account_id = "${aiven_account_team.developers.account_id}"
    team_id = "${aiven_account_team.developers.team_id}"
    project_name = "${aiven_project.project1.project}"
}
```

## Argument Reference

`account_id` is an unique account id.

`team_id` is an account team id.

`project_name` is a project name of already existing project.
