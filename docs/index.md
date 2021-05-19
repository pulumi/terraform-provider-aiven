# Aiven Terraform Provider
The Terraform provider for [Aiven.io](https://aiven.io/), an open source data platform as a service. 

## 🚨 A word of caution 🚨 
Recreating stateful services with Terraform will possibly **delete** the service and all its data before creating it again. Whenever the Terraform plan indicates that a service will be **deleted** or **replaced**, a catastrophic action is possibly about to happen.

Some properties, like **project** and the **resource name**, cannot be changed and it will trigger a resource replacement.

To avoid any issues, **please set the `termination_protection` property to `true` on all production services**, it will prevent Terraform to remove the service until the flag is set back to `false` again. While it prevents a service to be deleted, any logical databases, topics or other configurations may be removed **even when this section is enabled**. Be very careful! 

## Authentication token
[Signup at Aiven](https://console.aiven.io/signup) and see the [official instructions](https://help.aiven.io/en/articles/2059201-authentication-tokens) to create an API Authentication Token.

## Example usage
_Only available for Terraform v0.13 and above. For older versions, see [here](guides/install-terraform-v012.md)._

```hcl
terraform {
  required_providers {
    aiven = {
      source = "aiven/aiven"
      version = "2.X.X"
    }
  }
}

provider "aiven" {
  api_token = var.aiven_api_token
}
```

Then, initialize your Terraform workspace by running `terraform init`.

The `api_token` is the only parameter for the provider configuration. Make sure the owner of the API Authentication Token has admin permissions in Aiven.

You can also set the environment variable `AIVEN_TOKEN` for the `api_token` property.

## More examples
Look at the [Sample Project Guide](guides/sample-project.md) and the [Examples Guide](guides/examples.md) for more examples on how to use the various Aiven resources.

## Resource options
The list of options in this document is not comprehensive, for most part they map directly to the [Aiven REST API](https://api.aiven.io/doc/) properties.

For various objects called `x_user_config`, the exact configuration options are available in [Service User Config](https://github.com/aiven/terraform-provider-aiven/tree/master/aiven/templates/service_user_config_schema.json), [Integration User Config](https://github.com/aiven/terraform-provider-aiven/tree/master/aiven/templates/integrations_user_config_schema.json) and in [Integration Endpoint User Config](https://github.com/aiven/terraform-provider-aiven/tree/master/aiven/templates/integration_endpoints_user_config_schema.json) schema files.