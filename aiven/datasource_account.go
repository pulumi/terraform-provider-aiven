package aiven

import (
	"context"
	"github.com/aiven/aiven-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceAccountRead,
		Schema:      resourceSchemaAsDatasourceSchema(aivenAccountSchema, "name"),
	}
}

func datasourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*aiven.Client)

	name := d.Get("name").(string)

	r, err := client.Accounts.List()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, ac := range r.Accounts {
		if ac.Name == name {
			d.SetId(ac.Id)
			return resourceAccountRead(ctx, d, m)
		}
	}

	return diag.Errorf("account %s not found", name)
}
