package aiven

import (
	"fmt"
	"github.com/aiven/aiven-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var aivenKafkaConnectorSchema = map[string]*schema.Schema{
	"project": {
		Type:        schema.TypeString,
		Description: "Project to link the kafka connector to",
		Required:    true,
		ForceNew:    true,
	},
	"service_name": {
		Type:        schema.TypeString,
		Description: "Service to link the kafka connector to",
		Required:    true,
		ForceNew:    true,
	},
	"connector_name": {
		Type:        schema.TypeString,
		Description: "Kafka connector name",
		Required:    true,
		ForceNew:    true,
	},
	"config": {
		Type:        schema.TypeMap,
		Description: "Kafka Connector configuration parameters",
		Required:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"plugin_author": {
		Type:        schema.TypeString,
		Description: "Kafka connector author",
		Computed:    true,
	},
	"plugin_class": {
		Type:        schema.TypeString,
		Description: "Kafka connector Java class",
		Computed:    true,
	},
	"plugin_doc_url": {
		Type:        schema.TypeString,
		Description: "Kafka connector documentation URL",
		Computed:    true,
	},
	"plugin_title": {
		Type:        schema.TypeString,
		Description: "Kafka connector title",
		Computed:    true,
	},
	"plugin_type": {
		Type:        schema.TypeString,
		Description: "Kafka connector type",
		Computed:    true,
	},
	"plugin_version": {
		Type:        schema.TypeString,
		Description: "Kafka connector version",
		Computed:    true,
	},
	"task": {
		Type:        schema.TypeSet,
		Description: "List of tasks of a connector",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"connector": {
					Type:        schema.TypeString,
					Description: "Related connector name",
					Computed:    true,
				},
				"task": {
					Type:        schema.TypeInt,
					Description: "Task id / number",
					Computed:    true,
				},
			},
		},
	},
}

func resourceKafkaConnector() *schema.Resource {
	return &schema.Resource{
		Create: resourceKafkaConnectorCreate,
		Read:   resourceKafkaConnectorRead,
		Update: resourceKafkaTConnectorUpdate,
		Delete: resourceKafkaConnectorDelete,
		Exists: resourceKafkaConnectorExists,
		Importer: &schema.ResourceImporter{
			State: resourceKafkaConnectorState,
		},

		Schema: aivenKafkaConnectorSchema,
	}
}

func flattenKafkaConnectorTasks(r *aiven.KafkaConnector) []map[string]interface{} {
	var tasks []map[string]interface{}

	for _, taskS := range r.Tasks {
		task := map[string]interface{}{
			"connector": taskS.Connector,
			"task":      taskS.Task,
		}

		tasks = append(tasks, task)
	}

	return tasks
}

func resourceKafkaConnectorRead(d *schema.ResourceData, m interface{}) error {
	project, serviceName, connectorName := splitResourceID3(d.Id())

	res, err := m.(*aiven.Client).KafkaConnectors.List(project, serviceName)
	if err != nil {
		return fmt.Errorf("cannot read Kafka Connector List resource %s: %s", d.Id(), err)
	}

	var found bool
	for _, r := range res.Connectors {
		if r.Name == connectorName {
			found = true
			if err := d.Set("project", project); err != nil {
				return fmt.Errorf("error setting Kafka Connector `project` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("service_name", serviceName); err != nil {
				return fmt.Errorf("error setting Kafka Connector `service_name` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("connector_name", connectorName); err != nil {
				return fmt.Errorf("error setting Kafka Connector `connector_name` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("config", r.Config); err != nil {
				return fmt.Errorf("error setting Kafka Connector `config` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("plugin_author", r.Plugin.Author); err != nil {
				return fmt.Errorf("error setting Kafka Connector `plugin_author` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("plugin_class", r.Plugin.Class); err != nil {
				return fmt.Errorf("error setting Kafka Connector `plugin_class` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("plugin_doc_url", r.Plugin.DocumentationURL); err != nil {
				return fmt.Errorf("error setting Kafka Connector `plugin_doc_url` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("plugin_title", r.Plugin.Title); err != nil {
				return fmt.Errorf("error setting Kafka Connector `plugin_title` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("plugin_type", r.Plugin.Type); err != nil {
				return fmt.Errorf("error setting Kafka Connector `plugin_type` for resource %s: %s", d.Id(), err)
			}
			if err := d.Set("plugin_version", r.Plugin.Version); err != nil {
				return fmt.Errorf("error setting Kafka Connector `plugin_version` for resource %s: %s", d.Id(), err)
			}

			tasks := flattenKafkaConnectorTasks(&r)
			if err := d.Set("task", tasks); err != nil {
				return fmt.Errorf("error setting Kafka Connector `task` array for resource %s: %s", d.Id(), err)
			}
		}
	}

	if !found {
		return fmt.Errorf("cannot read Kafka Connector resource with Id: %s not found in a Kafka Connectors list", d.Id())
	}

	return nil
}

func resourceKafkaConnectorCreate(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	serviceName := d.Get("service_name").(string)
	connectorName := d.Get("connector_name").(string)

	config := make(aiven.KafkaConnectorConfig)
	for k, cS := range d.Get("config").(map[string]interface{}) {
		config[k] = cS.(string)
	}

	err := m.(*aiven.Client).KafkaConnectors.Create(project, serviceName, config)
	if err != nil {
		return err
	}

	d.SetId(buildResourceID(project, serviceName, connectorName))

	return resourceKafkaConnectorRead(d, m)
}

func resourceKafkaConnectorDelete(d *schema.ResourceData, m interface{}) error {
	return m.(*aiven.Client).KafkaConnectors.Delete(
		splitResourceID3(d.Id()))
}

func resourceKafkaTConnectorUpdate(d *schema.ResourceData, m interface{}) error {
	project, serviceName, connectorName := splitResourceID3(d.Id())

	config := make(aiven.KafkaConnectorConfig)
	for k, cS := range d.Get("config").(map[string]interface{}) {
		config[k] = cS.(string)
	}

	_, err := m.(*aiven.Client).KafkaConnectors.Update(project, serviceName, connectorName, config)
	if err != nil {
		return err
	}

	return resourceKafkaConnectorRead(d, m)
}

func resourceKafkaConnectorExists(d *schema.ResourceData, m interface{}) (bool, error) {
	project, serviceName, connectorName := splitResourceID3(d.Id())

	r, err := m.(*aiven.Client).KafkaConnectors.List(project, serviceName)
	if err != nil {
		return false, err
	}

	if ok, err := resourceExists(err); err != nil {
		return ok, err
	}

	for _, c := range r.Connectors {
		if c.Name == connectorName {
			return true, nil
		}
	}

	return false, nil
}

func resourceKafkaConnectorState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	err := resourceKafkaConnectorRead(d, m)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
