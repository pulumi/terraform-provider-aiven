---
title: Home
nav_exclude: true
---

# Terraform Aiven

[Terraform](https://www.terraform.io/) provider for [Aiven.io](https://aiven.io/). This
provider allows you to conveniently manage all resources for Aiven.

**A word of cautions**: While Terraform is an extremely powerful tool that can make
managing your infrastructure a breeze, great care must be taken when comparing the
changes that are about to applied to your infrastructure. When it comes to stateful
services, you cannot just re-create a resource and have it in the original state;
deleting a service deletes all data associated with it and there's often no way to
recover the data later. Whenever the Terraform plan indicates that a service, database,
topic or other such central construct is about to be deleted, something catastrophic is
quite possibly about to happen unless you're dealing with some throwaway test
environments or are deliberately retiring the service/database/topic.

There are many properties that cannot be changed after a resource is created and changing
those values later is handled by deleting the original resource and creating a new one.
These properties include such as the project a service is associated with, the name of a
service, etc. Unless the system contains no relevant data, such changes must not be
performed.

To allow mitigating this problem, the service resource supports
`termination_protection` property. It is recommended to set this property to `true`
for all production services to avoid them being accidentally deleted. With this setting
enabled service deletion, both intentional and unintentional, will fail until an explicit
update is done to change the setting to `false`. Note that while this does prevent the
service itself from being deleted, any databases, topics or such that have been configured
with Terraform can still be deleted and they will be deleted before the service itself is
attempted to be deleted so even with this setting enabled you need to be very careful
with the changes that are to be applied.

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) v0.12.X or greater
- [Go](https://golang.org/doc/install) 1.14.X or greater

## Installation

Work is underway to host the Aiven Terraform Provider under the official Terraform registry.
Until then, the easiest methods are:

**Terraform client version 0.12 and below**
```shell script
go get -u github.com/aiven/terraform-provider-aiven

# Global Install
mkdir -p ~/.terraform.d/plugins
cp $GOPATH/bin/terraform-provider-aiven ~/.terraform.d/plugins/ #Use a version here as a suffix if you specify a version in your plan

# Local Install
mkdir -p $PWD/terraform.d/plugins/linux_amd64 #darwin_amd64 if on OSX but linux will be needed for TF Cloud
cp $GOPATH/bin/terraform-provider-aiven $PWD/terraform.d/plugins/linux_amd64/terraform-provider-aiven_v1.2.4 # Change to release ver
```

**Terraform client version 0.13 and below**

Download the Aiven provider for Linux AMD64. 
```shell script
curl -Lo terraform-provider-aiven https://github.com/aiven/terraform-provider-aiven/releases/download/v2.X.X/terraform-provider-aiven-linux-amd64_v2.X.X

chmod +x terraform-provider-aiven
```
*Please specify version that you like to use and if you are not using Linux AMD64, 
download the correct binary for your system from the [release page](https://github.com/aiven/terraform-provider-aiven/releases).*

Third-party provider plugins — locally installed providers, not on the registry — need to be 
assigned an (arbitrary) source and placed in the appropriate subdirectory for Terraform to find and use them.
Create the appropriate subdirectory within the user plugins directory for the Aiven provider.
```shell script
mkdir -p ~/.terraform.d/plugins/aiven.io/provider/aiven/2.X/linux_amd64
```
*If you are not using Linux AMD64 replace `linux_amd64` with your `$OS_$ARCH`, and the same for the provider version.*

Finally, move the Aiven provider binary into the newly created directory.
```shell script
mv terraform-provider-aiven ~/.terraform.d/plugins/aiven.io/provider/aiven/2.X/linux_amd64
```

Now Aiven provider is in your user plugins directory, you can use the provider in your Terraform configuration.
```hcl-terraform
terraform {
  required_providers {
    aiven = {
      versions = [
        "2.X"
      ]
      source = "aiven.io/provider/aiven"
    }
  }
}

provider "aiven" {
  api_token = var.aiven_api_token
}
```
Then, initialize your Terraform workspace by running `terraform init`. If your Aiven provider  
is located in the correct directory, it should successfully initialize. Otherwise, move your 
Aiven provider to the correct directory: `~/.terraform.d/plugins/aiven.io/provider/aiven/$VERSION/$OS_$ARCH/`.

## Sample project

There is a [sample project](sample.tf) which sets up a project, defines Kafka,
PostgreSQL, InfluxDB and Grafana services, one PG database and user, one Kafka topic and
user, and metrics and dashboard integration for the Kafka and PG databases.

Make sure you have a look at the [variables](terraform.tfvars.sample) and copy it over to
`terraform.tfvars` with your own settings.

Other examples can be found in the [examples](examples) folder that provides examples to:
* [Getting Started](examples/getting-started.tf)
* [Account, projects, teams, and member management](examples/account)
* [Elasticsearch deployment and configuration](examples/elasticsearch)
* [Standalone Kafka connect deployment with custom config](examples/kafka_connect)
* [Deploying Kafka and Elasticsearch with a Kafka Connect Elasticsearch Sink connector](examples/kafka_connectors/es_sink)
* [Deploying Kafka and Elasticsearch with a Kafka Connect Mongo Sink connector](examples/kafka_connectors/mongo_sink)
* [Deploying Kafka with Schema Registry enabled and providing a schema](examples/kafka_schemas)
* [Deploying Cassandra and forking (cloning the service, config and data) into a new service with a higher plan](examples/cassandra_fork)
* [Deploying a Grafana service](examples/service)
* [Deploying a MirrorMaker service](examples/kafka_mirrormaker)
* [Deploying PostgreSQL services to multiple clouds and regions](examples/postgres)

## Importing existing infrastructure

All resources support importing so if you have already manually created an Aiven
environment using the web console or otherwise, it is possible to import all resources
and start managing them using Terraform. The documentation below mentions the ID format
for each resource but typically it is `<project_name>/<resource_name>` for resources
that are directly under project level and `<project_name>/<service_name>/<resource_name>`
for resources that belong to specific service. E.g. to import a database called `mydb`
belonging to service `myservice` in project `myproject` you'd do something like

```
terraform import aiven_database.mydb myproject/myservice/mydb
```

In some cases the internal identifiers are not shown in the Aiven web console. In such
cases the easiest way to obtain identifiers is typically to check network requests and
responses with your browser's debugging tools, as the raw responses do contain the IDs.

Note that as Terraform does not support creating configuration automatically you will
still need to manually create the Terraform configuration files. The import will just
match the existing resources with the ones defined in the configuration.

## Using datasources

Alternatively you can define already existing, or externally created and managed, resources
as datasources. Most of the resources are also available as datasources.

## Resource options

This section describes all the different resources that can be managed with this provider.
The list of options in this document is not comprehensive. For most part the options map
directly to [Aiven REST API](https://api.aiven.io/doc/) properties and that can be
consulted for details. For various objects called x_user_config the exact configuration
options are available in [Service User Config](aiven/templates/service_user_config_schema.json),
[Integration User Config](aiven/templates/integrations_user_config_schema.json) and in
[Integration Endpoint User Config](aiven/templates/integration_endpoints_user_config_schema.json).

### Provider

```
provider "aiven" {
    api_token = "<AIVEN_API_TOKEN>"
}
```

The Aiven provider currently only supports a single configuration option, `api_token`. This can also be specified with the AIVEN_TOKEN shell environment variable. 
The Aiven web console can be used to create named, never expiring API tokens that should
be used for this kind of purposes. If Terraform is used for managing existing project(s),
the API token must belong to a user with admin privileges for those project(s). For new
projects the user will be automatically granted admin role. For projects with credit card
billing this account must be the one in possession with the credit card used to pay for
the services.

### Resource Project

```
resource "aiven_project" "myproject" {
    project = "<PROJECT_NAME>"
    card_id = "<FULL_CARD_ID/LAST4_DIGITS>"
    account_id = aiven_account_team.<ACCOUNT_RESOURCE>.account_id
}
```

`project` defines the name of the project. Name must be globally unique (between all
Aiven customers) and cannot be changed later without destroying and re-creating the
project, including all sub-resources.

`card_id` is either the full card UUID or the last 4 digits of the card. As the full
UUID is not shown in the UI it is typically easier to use the last 4 digits to identify
the card. This can be omitted if `copy_from_project` is used to copy billing info from
another project.

`account_id` is an optional property to link a project to already an existing account by 
using account ID

`copy_from_project` is the name of another project used to copy billing information and
some other project attributes like technical contacts from. This is mostly relevant when
an existing project has billing type set to invoice and that needs to be copied over to a
new project. (Setting billing is otherwise not allowed over the API.) This only has
effect when the project is created.

`ca_cert` is a computed property that can be used to read the CA certificate of the
project. This is required for configuring clients that connect to certain services like
Kafka. This value cannot be set, only read.

Aiven ID format when importing existing resource: name of the project as is.

### Resource Account
```
resource "aiven_account" "account1" {
    name = "<ACCOUNT_NAME>"
}

```

`name` defines an account name.

`account_id` is an auto-generated unique account id.

### Resource Account Authentication
```
resource "aiven_account_authentication" "foo" {
    account_id = aiven_account.<ACCOUNT_RESOURCE>.account_id
    name = "auth-1"
    type = "saml"
    enabled = true
    saml_certificate = "---CERTIFICATE---"
    saml_entity_id = "https://example.com/00000"
    saml_idp_url = "https://example.com/sso/saml"
}
```

`account_id` is an unique account id.

`name` is an account authentication name.

`type` is an account authentication type, can be one of `internal` and `saml`.

`enabled` defines an authentication method enabled or not. 

`saml_certificate` is a SAML Certificate.

`saml_entity_id` is a SAML Entity Id.

`saml_idp_url` is a SAML Idp URL.

`saml_acs_url` is a SAML Assertion Consumer Service URL.

`saml_metadata_url` is a SAML Metadata URL

### Resource Account Team
```
resource "aiven_account_team" "account_team1" {
    account_id = aiven_account.<ACCOUNT_RESOURCE>.account_id
    name = "account_team1"
}
```

`name` defines an account team name.

`account_id` is an unique account id.

`team_id` is an auto-generated unique account team id.

### Resource Account Team Project
The account team project is intended to link and existing project to the existing account team. 
It is important to note that the project should have an `account_id` property set and equal to
account team you are trying to link this project. 

```
resource "aiven_project" "<PROJECT>" {
  project = "project-1"
  account_id = aiven_account_team.<ACCOUNT_RESOURCE>.account_id
}

resource "aiven_account_team_project" "account_team_project1" {
    account_id = aiven_account.<ACCOUNT_RESOURCE>.account_id
    team_id = aiven_account_team.<TEAM_RESOURCE>.team_id
    project_name = aiven_project.<PROJECT>.project
    team_type = "admin"
}
```

`account_id` is an unique account id.

`team_id` is an account team id.

`project_name` is a project name of already existing project.

`team_type` is an account team project type, can one of the following values: `admin`, 
`developer`, `operator` and `read_only`.

### Resource Account Team Member
During the creation of `aiven_account_team_member` resource, an email invitation will be sent  
to a user using `user_email` address. If the user accepts an invitation, he or she will become 
a member of the account team. The deletion of `aiven_account_team_member` will not only 
delete invitation if one was sent but not yet accepted by the user, and it will also 
eliminate an account team member if one has accepted an invitation previously.
 

```
resource "aiven_account_team_member" "foo" {
  account_id = aiven_account.<ACCOUNT_RESOURCE>.account_id
  team_id = aiven_account_team.<TEAM_RESOURCE>.team_id
  user_email = "user+1@example.com"
}

```

`account_id` is an unique account id.

`team_id` is an account team id.

`user_email` is a user email address that first will be invited, and after accepting an invitation,
he or she becomes a member of a team.

`accepted` is a boolean flag that determines whether an invitation was accepted or not by the user. 
`false` value means that the invitation was sent to the user but not yet accepted. 
`true` means that the user accepted the invitation and now a member of an account team 

### Resource Service

```
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
    
    timeouts {
        create = "20m"
        update = "15m"
    }
}
```

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
`pg` (PostgreSQL) and `redis`. This value cannot be changed after service creation.

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

`timeouts` a custom client timeouts.

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

Aiven ID format when importing existing resource: `<project_name>/<service_name>`.

### Separate Service Resources

Starting from the version 2 Aiven Provider supports separate Terraform resources for each 
service type available in Aiven Cloud. 

List of available resources:
- `aiven_pg` PostgreSQL service 
- `aiven_cassandra` Cassandra service
- `aiven_elasticsearch` Elasticsearch service
- `aiven_grafana` Grafana service
- `aiven_influxdb` Influxdb service
- `aiven_redis` Redis service
- `aiven_mysql` MySQL service
- `aiven_kafka` Kafka service
- `aiven_kafka_connect` Kafka Connect service
- `aiven_kafka_mirrormaker` Kafka Mirrormaker 2 service
 
Instructions on how to use them are similar to `aiven_service` with the exception that 
for example `aiven_kafka` contains only all the necessary configuration options related 
to this service type.

Each resource for certain service type has the following structure:
```
resource "aiven_<TYPE>" "my-service" {
    project = aiven_project.my-project.project
    cloud_name = "google-europe-west1"
    plan = "business-4"
    service_name = "my-service1"
    maintenance_window_dow = "monday"
    maintenance_window_time = "10:00:00"
    
    <TYPE>_user_config {
        ...
    
        <TYPE> {
            ...
        }
    }
}
```

### Resource Database

```
resource "aiven_database" "mydatabase" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    database_name = "<DATABASE_NAME>"
}
```

`project` and `service_name` define the project and service the database belongs to.
They should be defined using reference as shown above to set up dependencies correctly.

`database_name` is the actual name of the database.

None of the database properties can currently be changed after creation. Doing so will
result in the old database getting dropped and a new database created.

Aiven ID format when importing existing resource: `<project_name>/<service_name>/<database_name>`

### Resource Service User

```
resource "aiven_service_user" "myserviceuser" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    username = "<USERNAME>"
}
```

`project` and `service_name` define the project and service the user belongs to.
They should be defined using reference as shown above to set up dependencies correctly.

`username` is the actual name of the user account.

None of the service user properties can currently be changed after creation. Doing so
will result in the old database getting dropped and a new database created.

Service users have several computed properties that cannot be set, only read:

`password` is the password of the user (not applicable for all services).

`access_cert` is the access certificate of the user (not applicable for all services).

`access_key` is the access key of the user (not applicable for all services).

`type` tells whether the user is primary account or regular account.

Aiven ID format when importing existing resource: `<project_name>/<service_name>/<username>`

### Resource Connection Pool

```
resource "aiven_connection_pool" "mytestpool" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    database_name = "${aiven_database.mydatabase.database_name}"
    pool_mode = "transaction"
    pool_name = "mypool"
    pool_size = 10
    username = "${aiven_service_user.myserviceuser.username}"
}
```

`project` and `service_name` define the project and service the connection pool
belongs to. They should be defined using reference as shown above to set up dependencies
correctly. These properties cannot be changed once the service is created. Doing so will
result in the connection pool being deleted and new one created instead.

`database_name` is the name of the database the pool connects to. This should be
defined using reference as shown above to set up dependencies correctly.

`pool_mode` is the mode the pool operates in (session, transaction, statement).

`pool_name` is the name of the pool.

`pool_size` is the number of connections the pool may create towards the backend
server. This does not affect the number of incoming connections, which is always a much
larger number.

`username` is the name of the service user used to connect to the database. This should
be defined using reference as shown above to set up dependencies correctly.

`connection_uri` is a computed property that tells the URI for connecting to the pool.
This value cannot be set, only read.

Aiven ID format when importing existing resource: `<project_name>/<service_name>/<pool_name>`

### Resource Kafka Topic

```
resource "aiven_kafka_topic" "mytesttopic" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    topic_name = "<TOPIC_NAME>"
    partitions = 5
    replication = 3
    retention_bytes = -1
    retention_hours = 72
    minimum_in_sync_replicas = 2
    cleanup_policy = "delete"
    termination_protection = true

    timeouts {
        create = "1m"
        read = "5m"
    }
}
```

`project` and `service_name` define the project and service the topic belongs to.
They should be defined using reference as shown above to set up dependencies correctly.
These properties cannot be changed once the service is created. Doing so will result in
the topic being deleted and new one created instead.

`topic_name` is the actual name of the topic account. This propery cannot be changed
once the service is created. Doing so will result in the topic being deleted and new one
created instead.

`termination_protection` is a Terraform client-side deletion protection, which prevents a Kafka  
topic from being deleted. It is recommended to enable this for any production Kafka topic 
containing critical data.

Other properties should be self-explanatory. They can be changed after the topic has been
created.

Aiven ID format when importing existing resource: `<project_name>/<service_name>/<topic_name>`

`timeouts` a custom client timeouts.

### Resource Kafka ACL

```
resource "aiven_kafka_acl" "mytestacl" {
    project = "${aiven_project.myproject.project}"
    service_name = "${aiven_service.myservice.service_name}"
    topic = "<TOPIC_NAME_PATTERN>"
    permission = "admin"
    username = "<USERNAME_PATTERN>"
}
```

`project` and `service_name` define the project and service the ACL belongs to.
They should be defined using reference as shown above to set up dependencies correctly.
These properties cannot be changed once the service is created. Doing so will result in
the topic being deleted and new one created instead.

`topic` is a topic name pattern the ACL entry matches to.

`permission` is the level of permission the matching users are given to the matching
topics (admin, read, readwrite, write).

`username` is a username pattern the ACL entry matches to.

Aiven ID format when importing existing resource: `<project_name>/<service_name>/<acl_id>`.
The ACL ID is not directly visible in the Aiven web console.

### Resource Kafka Schemas
```hcl-terraform
# Kafka Schema global configuration
 resource "aiven_kafka_schema_configuration" "config" {
   project = aiven_project.kafka-schemas-project1.project
   service_name = aiven_service.kafka-service1.service_name
   compatibility_level = "BACKWARD"
 }
 
 # Kafka Schema
 resource "aiven_kafka_schema" "kafka-schema1" {
   project = aiven_project.kafka-schemas-project1.project
   service_name = aiven_service.kafka-service1.service_name
   subject_name = "kafka-schema1"
   compatibility_level = "FORWARD"
 
   schema = <<EOT
    	  {
           "doc": "example",
           "fields": [{
               "default": 5,
               "doc": "my test number",
               "name": "test",
               "namespace": "test",
               "type": "int"
           }],
           "name": "example",
           "namespace": "example",
           "type": "record"
       }
     EOT
 }
```

`project` and `service_name` define the project and service the Kafka Schemas belongs to. 
They should be defined using reference as shown above to set up dependencies correctly.

`compatibility_level` is the Global Kafka Schema configuration compatibility level when defined 
for `aiven_kafka_schema_configuration` resource. Also, Kafka Schema configuration 
compatibility level can be overridden for a specific subject when used in `aiven_kafka_schema` 
resource. If the compatibility level not specified for the individual subject by default, 
it takes a global value. Allowed values: `BACKWARD`, `BACKWARD_TRANSITIVE`, `FORWARD`, 
`FORWARD_TRANSITIVE`, `FULL`, `FULL_TRANSITIVE`, `NONE`.

`subject_name` is Kafka Schema subject name.

`schema` is Kafka Schema configuration should be a valid Avro Schema JSON format.

### Resource Kafka connectors
```
resource "aiven_kafka_connector" "kafka-es-con1" {
  project = aiven_project.kafka-con-project1.project
  service_name = aiven_service.kafka-service1.service_name
  connector_name = "kafka-es-con1"

  config = {
    "topics" = aiven_kafka_topic.kafka-topic1.topic_name
    "connector.class" : "io.aiven.connect.elasticsearch.ElasticsearchSinkConnector"
    "type.name" = "es-connector"
    "name" = "kafka-es-con1"
    "connection.url" = aiven_service.es-service1.service_uri
  }
}
```

`project` and `service_name` define the project and service the Kafka Connectors belongs to. 
They should be defined using reference as shown above to set up dependencies correctly.

`connector_name` is the Kafka connector name

`config` is the Kafka Connector configuration parameters, where `topics`, `connector.class` and `name` 
are required parameters but the rest of them are connector type specific. 

### Resource MirrorMaker 2 Replication Flow
```
resource "aiven_mirrormaker_replication_flow" "f1" {
  project = aiven_project.kafka-mm-project1.project
  service_name = aiven_service.mm.service_name
  source_cluster = aiven_service.source.service_name
  target_cluster = aiven_service.target.service_name
  enable = true

  topics = [
    ".*",
  ]

  topics_blacklist = [
    ".*[\\-\\.]internal",
    ".*\\.replica",
    "__.*"
  ]
}
```

`project` and `service_name` define the project and service the Kafka MirrorMaker Replication 
Flow belongs to. They should be defined using reference as shown above to set up dependencies correctly.

`source_cluster` is a source cluster alias.

`target_cluster` is a target cluster alias.

`enable` enable of disable replication flows for a mirror maker service 

`topics` is a list of topics and/or regular expressions to replicate.

`topics_blacklist` is a list of topics and/or regular expressions to not replicate.

### Resource Elasticsearch ACL

```
resource "aiven_elasticsearch_acl" "es-acls" {
     project = aiven_project.es-project.project
     service_name = aiven_service.es.service_name
     enabled = true
     extended_acl = false
     acl {
       username = aiven_service_user.es-user.username
       rule {
         index = "_*"
         permission = "admin"
       }
   
       rule {
         index = "*"
         permission = "admin"
       }
     }
   
     acl {
       username = "avnadmin"
       rule {
         index = "_*"
         permission = "read"
       }
   
       rule {
         index = "*"
         permission = "read"
       }
     }
   }
```

`project` and `service_name` define the project and service the ACL belongs to. 
They should be defined using reference as shown above to set up dependencies correctly.

All other properties except `project` and `service_name` can be changed after creation of the 
resource and will not trigger recreation of Elasticsearch entire ACL's. 

`enabled` enables of disables Elasticsearch ACL's.

`extended_acl` Index rules can be applied in a limited fashion to the _mget, _msearch and _bulk APIs 
(and only those) by enabling the ExtendedAcl option for the service. When it is enabled, users can use 
 these APIs as long as all operations only target indexes they have been granted access to.
 
`acl.username` is the name of the existing service user, and service user must be preliminary added 
to the Elasticsearch service; this can be done using `aiven_service_user` resource. Aiven has a 
default user `avnadmin` which is automatically created as a part of the creation process of Elasticsearch service. 

Elasticsearch ACL support multiple rules for a single user.

`acl.rule.index` is the Elasticsearch index pattern.

`acl.rule.permission` is the Elasticsearch permission, list of supported permissions: "deny", "admin", "read", "readwrite", "write"
### Resource Service Integration Endpoint

Service Integration Endpoint is an external endpoint that can be used as a source or
destination for a Service Integration. These are only defined for non-Aiven resources.

```
resource "aiven_service_integration_endpoint" "myendpoint" {
    project = "${aiven_project.myproject.project}"
    endpoint_name = "<ENDPOINT_NAME>"
    endpoint_type = "datadog"
    datadog_user_config {
        datadog_api_key = "<DATADOG_API_KEY>"
    }
}
```

`project` defines the project the endpoint is associated with.

`endpoint_name` is the name of the endpoint. This value has no effect beyond being used
to identify different integration endpoints.

`endpoint_type` is the type of the external service this endpoint is associated with.
By the time of writing the only available option is `datadog`.

`x_user_config` defines endpoint type specific configuration. `x` is the type of the
endpoint. The available configuration options are documented in
[this JSON file](aiven/templates/integration_endpoints_user_config_schema.json)

Aiven ID format when importing existing resource: `<project_name>/<endpoint_id>`. The
endpoint identifier (UUID) is not directly visible in the Aiven web console.

### Resource Service Integration

Service Integration defines an integration between two Aiven services or between Aiven
service and an external integration endpoint. Integration could be for example sending
metrics from Kafka service to an InfluxDB service, getting metrics from an InfluxDB
service to a Grafana service to show dashboards, sending logs from any service to
Elasticsearch, etc.

```
resource "aiven_service_integration" "myintegration" {
    project = "${aiven_project.myproject.project}"
    destination_endpoint_id = "${aiven_service_integration_endpoint.myendpoint.id}"
    destination_service_name = ""
    integration_type = "datadog"
    source_endpoint_id = ""
    source_service_name = "${aiven_service.testkafka.service_name}"
}
```

`project` defines the project the integration belongs to.

`destination_endpoint_id` or `destination_service_name` identifies the target side of
the integration. Only either endpoint identifier or service name must be specified. In
either case the target needs to be defined using the reference syntax described above to
set up the dependency correctly.

`integration_type` identifies the type of integration that is set up. Possible values
include `dashboard`, `datadog`, `logs`, `metrics` and `mirrormaker`.

`source_endpoint_id` or `source_service_name` identifies the source side of the
integration. Only either endpoint identifier or service name must be specified. In either
case the source needs to be defined using the reference syntax described above to set up
the dependency correctly.

`x_user_config` defines integration specific configuration. `x` is the type of the
integration. The available configuration options are documented in
[this JSON file](aiven/templates/integrations_user_config_schema.json). Not all integration
types have any configurable settings.

Aiven ID format when importing existing resource: `<project_name>/<integration_id>`.
The integration identifier (UUID) is not directly visible in the Aiven web console.

### Resource Project User

```
resource "aiven_project_user" "mytestuser" {
    project = "${aiven_project.myproject.project}"
    email = "john.doe@example.com"
    member_type = "admin"
}
```

`project` defines the project the user is a member of.

`email` identifies the email address of the user.

`member_type` defines the access level the user has to the project.

Computed property `accepted` tells whether the user has accepted the request to join
the project; adding user to a project sends an invitation to the target user and the
actual membership is only created once the user accepts the invitation. This property
cannot be set, only read.

Aiven ID format when importing existing resource: `<project_name>/<email>`

### Resource Project VPC

```
resource "aiven_project_vpc" "myvpc" {
    project = "${aiven_project.myproject.project}"
    cloud_name = "google-europe-west1"
    network_cidr = "192.168.0.1/24"

    timeouts {
        create = "5m"
    }
}
```

`project` defines the project the VPC belongs to.

`cloud_name` defines where the cloud provider and region where the service is hosted
in. See the Service resource for additional information.

`network_cidr` defines the network CIDR of the VPC.

`timeouts` a custom client timeouts.

Computed property `state` tells the current state of the VPC. This property cannot be
set, only read.

Aiven ID format when importing existing resource: `<project_name>/<VPC_UUID>`. The UUID
is not directly visible in the Aiven web console.

## Resource VPC Peering Connection

```
resource "aiven_vpc_peering_connection" "mypeeringconnection" {
    vpc_id = "${aiven_project_vpc.myvpc.id}"
    peer_cloud_account = "<PEER_ACCOUNT_ID>"
    peer_vpc = "<PEER_VPC_ID/NAME>"
    peer_region = "<PEER_REGION>"

    timeouts {
        create = "10m"
    }
}
```

`vpc_id` is the Aiven VPC the peering connection is associated with.

`peer_cloud_account` defines the identifier of the cloud account the VPC is being
peered with.

`peer_vpc` defines the identifier or name of the remote VPC.

`peer_region` defines the region of the remote VPC if it is not in the same region as Aiven VPC.

`timeouts` a custom client timeouts.

`state` is the state of the peering connection. This property is computed by Aiven 
therefore cannot be set, only read. Where state can be one of: `APPROVED`, 
`PENDING_PEER`, `ACTIVE`, `DELETED`, `DELETED_BY_PEER`, `REJECTED_BY_PEER` and 
`INVALID_SPECIFICATION`. 

Happy path of the `state` includes the following transitions of a VPC peering connection: 
`APPROVED` -> `PENDING_PEER` -> `ACTIVE`.

More details regarding each state:

- `APPROVED` is the initial state after the user does a successful creation of a 
peering connection resource via Terraform. 

- `PENDING_PEER` the connection enters the `PENDING_PEER` state from `APPROVED` once the 
Aiven platform has created a connection to the specified peer successfully in the cloud, 
but the connection is not active until the user completes the setup in their cloud account. 
The steps needed in the user cloud account depend on the used cloud provider.

- `ACTIVE` stands for a VPC peering connection whose setup has been completed

- `DELETED` means a user deleted the peering connection through the Aiven Terraform provider, 
or Aiven Web Console or directly via Aiven API.

- `DELETED_BY_PEER` appears when a user deleted the VPC peering connection through their cloud 
account. That is, the user deleted the peering cloud resource in their account. There are no 
transitions from this state

- `REJECTED_BY_PEER` an AWS specific state, when VPC peering connection was in the `PENDING_PEER` state, 
and the user rejected the AWS peering connection request.

- `INVALID_SPECIFICATION` is a VPC peering connection that was in the `APPROVED` state but could not be  
successfully created because of something in the user's control, for example, the peer cloud account of VPC 
doesn't exist, overlapping IP ranges, or the Aiven cloud account doesn't have permissions to peer 
there. `state_info` field contains more details about the particular issue.

`state_info` state-specific help or error information.

Aiven ID format when importing existing resource: `<project_name>/<VPC_UUID>/<peer_cloud_account>/<peer_vpc>`.
Aiven ID format when importing existing cross-region resource: `<project_name>/<VPC_UUID>/<peer_cloud_account>/<peer_vpc>/peer_region`.
The UUID is not directly visible in the Aiven web console.

## Resource Transit Gateway VPC Attachment

```
resource "aiven_transit_gateway_vpc_attachment" "attachment" {
    vpc_id = aiven_project_vpc.bar.id
    peer_cloud_account = "<PEER_ACCOUNT_ID>"
    peer_vpc = "google-project1"
    peer_region = "google-europe-west1"
    user_peer_network_cidrs = [ "10.0.0.0/24" ]
}
```

`vpc_id` is the Aiven VPC the peering connection is associated with.

`peer_cloud_account` AWS account ID of the peered VPC.

`peer_vpc` Transit gateway ID

`peer_region` AWS region of the peered VPC (if not in the same region as Aiven VPC).

`user_peer_network_cidrs` List of private IPv4 ranges to route through the peering connection.

`timeouts` a custom client timeouts.

`state` is the state of the peering connection.

`state_info` state-specific help or error information.

`peering_connection_id` Cloud provider identifier for the peering connection if available.

## Datasource options

This section describes all the different datasources supported by the provider. The options
listed describe the options required to identify the datasource, see the equivalent resource
for options that can be referenced via the datasource.

(All examples are written using datasource references.)

### Datasource Project

```
data "aiven_project" "myproject" {
    project = "<PROJECT_NAME>"
}
```

`project` defines the name of the project.

### Datasource Service

```
data "aiven_service" "myservice" {
    project = data.aiven_project.myproject.project
    service_name = "<SERVICE_NAME>"
}
```

`project` identifies the project the service belongs to.

`service_name` specifies the actual name of the service.

### Datasource Database

```
data "aiven_database" "mydatabase" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    database_name = "<DATABASE_NAME>"
    termination_protection = true
}
```

`project` and `service_name` define the project and service the database belongs to.

`database_name` is the actual name of the database.

`termination_protection` is a Terraform client-side deletion protections, which prevents 
the database from being deleted by Terraform. It is recommended to enable this for any 
production databases containing critical data.

### Datasource Service User

```
data "aiven_service_user" "myserviceuser" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    username = "<USERNAME>"
}
```

`project` and `service_name` define the project and service the user belongs to.

`username` is the actual name of the user account.

### Datasource Connection Pool

```
data "aiven_connection_pool" "mytestpool" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    pool_name = "<POOLNAME>"
}
```

`project` and `service_name` define the project and service the connection pool
belongs to.

`pool_name` is the name of the pool.

### Datasource Kafka Topic

```
data "aiven_kafka_topic" "mytesttopic" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    topic_name = "<TOPIC_NAME>"
}
```

`project` and `service_name` define the project and service the topic belongs to.

`topic_name` is the actual name of the topic account.

### Datasource Kafka ACL

```
data "aiven_kafka_acl" "mytestacl" {
    project = data.aiven_service.myservice.project
    service_name = data.aiven_service.myservice.service_name
    topic = "<TOPIC_NAME_PATTERN>"
    username = "<USERNAME_PATTERN>"
}
```

`project` and `service_name` define the project and service the ACL belongs to.

`topic` is a topic name pattern the ACL entry matches to.

`username` is a username pattern the ACL entry matches to.

### Datasource Service Integration Endpoint

```
data "aiven_service_integration_endpoint" "myendpoint" {
    project = data.aiven_project.myproject.project
    endpoint_name = "<ENDPOINT_NAME>"
}
```

`project` defines the project the endpoint is associated with.

`endpoint_name` is the name of the endpoint.

### Datasource Project User

```
data "aiven_project_user" "mytestuser" {
    project = data.aiven_project.myproject.project
    email = "john.doe@example.com"
}
```

`project` defines the project the user is a member of.

`email` identifies the email address of the user.

### Datasource Project VPC

```
data "aiven_project_vpc" "myvpc" {
    project = data.aiven_project.myproject.project
    cloud_name = "google-europe-west1"
}
```

`project` defines the project the VPC belongs to.

`cloud_name` defines the cloud provider and region where the vpc resides.

### Datasource VPC Peering Connection

```
data "aiven_vpc_peering_connection" "mypeeringconnection" {
    vpc_id = data.aiven_project_vpc.vpc_id
    peer_cloud_account = "<PEER_ACCOUNT_ID>"
    peer_vpc = "<PEER_VPC_ID/NAME>"
}
```

`vpc_id` is the Aiven VPC the peering connection is associated with.

`peer_cloud_account` defines the identifier of the cloud account the VPC has been
peered with.

`peer_vpc` defines the identifier or name of the remote VPC.

### Datasource Service Component
Service components can be defined to get the connection info for specific service. 
Services may support multiple different access routes (VPC peering and public access), 
have additional components or support various authentication methods. Each of these 
may be represented by different DNS name or TCP port and the specific component to 
match can be selected by specifying appropriate filters as shown below.

```hcl-terraform
data "aiven_service_component" "sc1" {
    project = aiven_kafka.project1.project
    service_name = aiven_kafka.service1.service_name
    component = "kafka"
    route = "dynamic"
    kafka_authentication_method = "certificate"
}
```

`project` and `service_name` define the project and service the service component
belongs to.

`component` is a service component name. Component may match the name of the service 
(`cassandra`, `elasticsearch`, `grafana`, `influxdb`, `kafka`, `kafka_connect`, `mysql`, 
`pg` and `redis`), in which case the connection info of the service itself is returned. 
Some service types support additional service specific components like `kibana` for 
Elasticsearch, `kafka_connect`, `kafka_rest` and `schema_registry` for Kafka, and 
`pgbouncer` for PostgreSQL. Most service types also support `prometheus`.

`route` is network access route. The route may be one of `dynamic`, `public`, and `private`. 
Usually, you'll want to use `dynamic`, which for services that are not in a private network 
identifies the regular public DNS name of the service and for services in a private network  
the private DNS name. If the service is in a private network but has also public access  
enabled the `public` route type can be used to get the public DNS name of the service. The  
`private` option should typically not be used.

`host` is DNS name for connecting to the service component.

`port` is port number for connecting to the service component.

`kafka_authentication_method` is Kafka authentication method. This is a value specific 
to the 'kafka' service components. And has the following available options: `certificate` 
and `sasl`. If not set by the user only entries with empty `kafka_authentication_method` 
will be selected.

`ssl` whether the endpoint is encrypted or accepts plaintext. By default endpoints are
always encrypted and this property is only included for service components they may
disable encryption. If not set by the user only entries with empty `ssl` or `ssl` set 
to true will be selected.

`usage` is DNS usage name, and can be one of `primary`, `replica` or `syncing`. `replica` 
is used by services that have separate master and standby roles for which it identifies 
the `replica` DNS name. `syncing` is used by limited set of services to expose nodes 
before they have finished restoring state but may already be partially available, for 
example a PostgreSQL node that is streaming WAL segments from backup or current master 
but hasn't yet fully caught up.

## Credits

The original version of the Aiven Terraform provider was written and maintained by
Jelmer Snoeck (https://github.com/jelmersnoeck).
