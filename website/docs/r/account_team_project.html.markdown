---
layout: "aiven"
page_title: "Aiven: aiven_account_team_project"
description: |-
  Manages an Aiven account team project
---

# Resource: aiven_account_team_project

The account team project is intended to link and existing project to the existing account team. It is important to note 
that the project should have an `account_id` property set and equal to account team you are trying to link this project.

## Example Usage

```hcl
resource "aiven_project" "project1" {
  project = "project-1"
  account_id = "${aiven_account_team.developers.account_id}"
}

resource "aiven_account_team_project" "account_team_project1" {
    account_id = "${aiven_account_team.developers.account_id}"
    team_id = "${aiven_account_team.developers.team_id}"
    project_name = "${aiven_project.project1.project}"
    team_type = "admin"
}
```

## Argument Reference

`account_id` is an unique account id.

`team_id` is an account team id.

`project_name` is a project name of already existing project.

`team_type` is an account team project type, can one of the following values: `admin`, `developer`, `operator` and `read_only`.
