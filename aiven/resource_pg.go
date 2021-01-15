package aiven

import (
	"context"
	"github.com/aiven/aiven-go-client"
	"github.com/aiven/terraform-provider-aiven/aiven/templates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

func aivenPGSchema() map[string]*schema.Schema {
	schemaPG := serviceCommonSchema()
	schemaPG[ServiceTypePG] = &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Computed:    true,
		Description: "PostgreSQL specific server provided values",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"replica_uri": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "PostgreSQL replica URI for services with a replica",
					Sensitive:   true,
				},
				"uri": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "PostgreSQL master connection URI",
					Optional:    true,
					Sensitive:   true,
				},
				"dbname": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Primary PostgreSQL database name",
				},
				"host": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "PostgreSQL master node host IP or name",
				},
				"password": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "PostgreSQL admin user password",
					Sensitive:   true,
				},
				"port": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "PostgreSQL port",
				},
				"sslmode": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "PostgreSQL sslmode setting (currently always \"require\")",
				},
				"user": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "PostgreSQL admin user name",
				},
			},
		},
	}
	schemaPG[ServiceTypePG+"_user_config"] = &schema.Schema{
		Type:             schema.TypeList,
		MaxItems:         1,
		Optional:         true,
		Description:      "PostgreSQL specific user configurable settings",
		DiffSuppressFunc: emptyObjectDiffSuppressFunc,
		Elem: &schema.Resource{
			Schema: GenerateTerraformUserConfigSchema(
				templates.GetUserConfigSchema("service")[ServiceTypePG].(map[string]interface{})),
		},
	}

	return schemaPG
}

func resourcePG() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreateWrapper(ServiceTypePG),
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServicePGUpdate,
		DeleteContext: resourceServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServiceState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(20 * time.Minute),
			Update:  schema.DefaultTimeout(20 * time.Minute),
			Default: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: aivenPGSchema(),
	}
}

func resourceServicePGUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*aiven.Client)

	projectName, serviceName := splitResourceID2(d.Id())
	userConfig := ConvertTerraformUserConfigToAPICompatibleFormat("service", "pg", false, d)

	if userConfig["pg_version"] != nil {
		service, err := client.Services.Get(projectName, serviceName)
		if err != nil {
			return diag.Errorf("cannot get a service: %s", err)
		}

		if userConfig["pg_version"].(string) != service.UserConfig["pg_version"].(string) {
			t, err := client.ServiceTask.Create(projectName, serviceName, aiven.ServiceTaskRequest{
				TargetVersion: userConfig["pg_version"].(string),
				TaskType:      "upgrade_check",
			})
			if err != nil {
				return diag.Errorf("cannot create PG upgrade check task: %s", err)
			}

			w := &ServiceTaskWaiter{
				Client:      m.(*aiven.Client),
				Project:     projectName,
				ServiceName: serviceName,
				TaskId:      t.Task.Id,
			}

			taskI, err := w.Conf(d.Timeout(schema.TimeoutDefault)).WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for Aiven service task to be DONE: %s", err)
			}

			task := taskI.(*aiven.ServiceTaskResponse)
			if !*task.Task.Success {
				return diag.Errorf(
					"PG service upgrade check error, version upgrade from %s to %s, result: %s",
					task.Task.SourcePgVersion, task.Task.TargetPgVersion, task.Task.Result)
			}

			log.Printf("[DEBUG] PG service upgrade check result: %s", task.Task.Result)
		}
	}

	return resourceServiceUpdate(ctx, d, m)
}

// ServiceTaskWaiter is used to refresh the Aiven Service Task endpoints when
// provisioning.
type ServiceTaskWaiter struct {
	Client      *aiven.Client
	Project     string
	ServiceName string
	TaskId      string
}

// RefreshFunc will call the Aiven client and refresh its state.
func (w *ServiceTaskWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		t, err := w.Client.ServiceTask.Get(
			w.Project,
			w.ServiceName,
			w.TaskId,
		)
		if err != nil {
			return nil, "", err
		}

		if t.Task.Success == nil {
			return nil, "IN_PROGRESS", nil
		}

		return t, "DONE", nil
	}
}

// Conf sets up the configuration to refresh.
func (w *ServiceTaskWaiter) Conf(timeout time.Duration) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:                   []string{"IN_PROGRESS"},
		Target:                    []string{"DONE"},
		Refresh:                   w.RefreshFunc(),
		Delay:                     10 * time.Second,
		Timeout:                   timeout,
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 3,
	}
}
