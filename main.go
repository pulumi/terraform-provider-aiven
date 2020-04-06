package main

import (
	"github.com/aiven/terraform-provider-aiven/aiven"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aiven.Provider})
}
