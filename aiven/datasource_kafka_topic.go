// Copyright (c) 2019 Aiven, Helsinki, Finland. https://aiven.io/
package aiven

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceKafkaTopicRead,
		Schema: resourceSchemaAsDatasourceSchema(aivenKafkaTopicSchema,
			"project", "service_name", "topic_name"),
	}
}

func datasourceKafkaTopicRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectName := d.Get("project").(string)
	serviceName := d.Get("service_name").(string)
	topicName := d.Get("topic_name").(string)

	d.SetId(buildResourceID(projectName, serviceName, topicName))

	return resourceKafkaTopicRead(ctx, d, m)
}
