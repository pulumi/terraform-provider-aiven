---
layout: "provider"
page_title: "Provider: Aiven"
sidebar_current: "docs-aiven-index"
description: |-
  This   provider allows you to conveniently manage all resources for Aiven. The provider needs to be configured with an Aiven API token before it can be used.
---

# Aiven Provider

```
provider "aiven" {
    api_token = "<AIVEN_API_TOKEN>"
}
```

The Aiven provider currently only supports a single configuration option, `api_token`.
The Aiven web console can be used to create named, never expiring API tokens that should
be used for this kind of purposes. If Terraform is used for managing existing project(s),
the API token must belong to a user with admin privileges for those project(s). For new
projects the user will be automatically granted admin role. For projects with credit card
billing this account must be the one in possession with the credit card used to pay for
the services.