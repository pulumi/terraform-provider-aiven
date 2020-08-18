package aiven

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccAivenKafkaSchemaConfiguration_basic(t *testing.T) {
	resourceName := "aiven_kafka_schema_configuration.foo"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckAivenKafkaSchemaConfigurationResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaSchemaConfigurationResource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project", os.Getenv("AIVEN_PROJECT_NAME")),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "BACKWARD"),
				),
			},
		},
	})
}

// Kafka Schemas configuration cannot be deleted
func testAccCheckAivenKafkaSchemaConfigurationResourceDestroy(_ *terraform.State) error {
	return nil
}

func testAccKafkaSchemaConfigurationResource(name string) string {
	return fmt.Sprintf(`
		data "aiven_project" "foo" {
			project = "%s"
		}

		resource "aiven_service" "bar" {
			project = data.aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-4"
			service_name = "test-acc-sr-%s"
			service_type = "kafka"
			maintenance_window_dow = "monday"
			maintenance_window_time = "10:00:00"
			
			kafka_user_config {
				schema_registry = true
				kafka_version = "2.4"
				kafka {
				  group_max_session_timeout_ms = 70000
				  log_retention_bytes = 1000000000
				}
			}
		}
		
		resource "aiven_kafka_schema_configuration" "foo" {
			project = data.aiven_project.foo.project
			service_name = aiven_service.bar.service_name
			compatibility_level = "BACKWARD"
		}
		`, os.Getenv("AIVEN_PROJECT_NAME"), name)
}
