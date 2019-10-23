---
layout: "aiven"
page_title: "Aiven: aiven_service"
description: |-
  Manages an Aiven service resource.
---

# Resource: aiven_service

## Example Usage

```hcl
resource "aiven_service" "myservice" {
    project = "${aiven_project.myproject.project}"
    cloud_name = "google-europe-west1"
    plan = "business-8"
    service_name = "<SERVICE_NAME>"
    service_type = "pg"
    project_vpc_id = "${aiven_project_vpc.vpc_gcp_europe_west1.id}"
    termination_protection = true
    pg_user_config {
        ip_filter = ["0.0.0.0/0"]
        pg_version = "10"
    }
}
```

## Argument Reference

`project` identifies the project the service belongs to. To set up proper dependency
between the project and the service, refer to the project as shown in the above example.
Project cannot be changed later without destroying and re-creating the service.

`cloud_name` defines where the cloud provider and region where the service is hosted
in. This can be changed freely after service is created. Changing the value will trigger
a potentially lenghty migration process for the service. Format is cloud provider name
(`aws`, `azure`, `do` `google`, `upcloud`, etc.), dash, and the cloud provider
specific region name. These are documented on each Cloud provider's own support articles,
like [here for Google](https://cloud.google.com/compute/docs/regions-zones/) and
[here for AWS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html).

`plan` defines what kind of computing resources are allocated for the service. It can
be changed after creation, though there are some restrictions when going to a smaller
plan such as the new plan must have sufficient amount of disk space to store all current
data and switching to a plan with fewer nodes might not be supported. The basic plan
names are `hobbyist`, `startup-x`, `business-x` and `premium-x` where `x` is
(roughly) the amount of memory on each node (also other attributes like number of CPUs
and amount of disk space varies but naming is based on memory). The exact options can be
seen from the Aiven web console's Create Service dialog.

`service_name` specifies the actual name of the service. The name cannot be changed
later without destroying and re-creating the service so name should be picked based on
intended service usage rather than current attributes.

`service_type` is the actual service that is being provided. Currently available
options are `cassadra`, `elasticsearch`, `grafana`, `influxdb`, `kafka`,
`pg` (PostreSQL) and `redis`. This value cannot be changed after service creation.

`project_vpc_id` optionally specifies the VPC the service should run in. If the value
is not set the service is not run inside a VPC. When set, the value should be given as a
reference as shown above to set up dependencies correctly and the VPC must be in the same
cloud and region as the service itself. Project can be freely moved to and from VPC after
creation but doing so triggers migration to new servers so the operation can take
significant amount of time to complete if the service has a lot of data.

`termination_protection` prevents the service from being deleted. It is recommended to
set this to `true` for all production services to prevent unintentional service
deletions. This does not shield against deleting databases or topics but for services
with backups much of the content can at least be restored from backup in case accidental
deletion is done.

`x_user_config` defines service specific additional configuration options. These
options can be found from the [JSON schema description](aiven/templates/service_user_config_schema.json).

For services that support different versions the version information must be specified in
the user configuration. By the time of writing these services are Elasticsearch, Kafka
and PostgreSQL. These services should have configuration like

```
elasticsearch_user_config {
    elasticsearch_version = "6"
}
```

```
kafka_user_config {
    kafka_version = "2.0"
}
```

```
pg_user_config {
    pg_version = "10"
}
```

Some (very few) of the user configuration options have a dot (`.`) in their name.
That is not supported by Terraform so the provider converts any literal dots to the
text string `__dot__`. So if you want to set `foo.bar = "abc"` you need to instead
set `foo__dot__bar = "abc"`.

`service_(uri|host|port|username|password)` are computed properties that define the
URI for connecting to the service and the same info split into various parts. These
values cannot be set, only read.

`x` defines service specific additional computed values for connecting to the service
(where `x` is the type of the service). E.g. `elasticsearch.0.kibana_uri` specifies
the Kibana URI for Elasticsearch service (while `service_uri` at main level is the
connection URI for the actual Elasticsearch service itself). Note the need for using
`.0` when accessing the values due to Terraform's restrictions in defining nested
schematized values. These values cannot be set, only read.

`service_integrations` can be used to define service integrations that must exist
immediately upon service creation. By the time of writing the only such integration is
defining that MySQL service is a read-replica of another service. To define a read-
replica the following configuration needs to be added:

```
service_integrations {
    integration_type = "read_replica"
    source_service_name = "${aiven_service.mysourceservice.service_name}"
}
```

Making changes to the service integrations as well as removing the service integration
requires defining an explicit `aiven_service_integration` resource with the same
attributes (plus `project` and `destination_service_name` attributes); the backend
will handle creation of an existing read-replica integration as a no-op and will just
return the identifier of the existing integration.

## Import

Services can be imported using their ID in the format `<project_name>/<service_name>` , e.g.

```
$ terraform import aiven_service.my_service test-project/test-service
```
