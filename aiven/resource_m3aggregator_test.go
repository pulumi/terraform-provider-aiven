package aiven

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccAiven_m3aggregator(t *testing.T) {
	resourceName := "aiven_m3aggregator.bar"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAivenServiceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccM3AggregatorResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAivenServiceCommonAttributes("data.aiven_m3aggregator.service"),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-m3a-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "project", os.Getenv("AIVEN_PROJECT_NAME")),
					resource.TestCheckResourceAttr(resourceName, "service_type", "m3aggregator"),
					resource.TestCheckResourceAttr(resourceName, "cloud_name", "google-europe-west1"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_window_dow", "monday"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_window_time", "10:00:00"),
					resource.TestCheckResourceAttr(resourceName, "state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "termination_protection", "false"),
				),
			},
		},
	})
}

func testAccM3AggregatorResource(name string) string {
	return fmt.Sprintf(`
		data "aiven_project" "foo" {
			project = "%s"
		}

		resource "aiven_m3db" "foo" {
			project = data.aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-8"
			service_name = "test-acc-m3d-%s"

			m3db_user_config {
				m3_version = 0.15
							
				namespaces {
					name = "test-acc-%s"
					type = "unaggregated"
				}
			}
		}
		
		resource "aiven_m3aggregator" "bar" {
			project = data.aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-8"
			service_name = "test-acc-m3a-%s"
			maintenance_window_dow = "monday"
			maintenance_window_time = "10:00:00"

			m3aggregator_user_config {
				m3_version = 0.15
			}
		}

		resource "aiven_service_integration" "int-m3db-aggr" {
			project = data.aiven_project.foo.project
			integration_type = "m3aggregator"
			source_service_name = aiven_m3db.foo.service_name
			destination_service_name = aiven_m3aggregator.bar.service_name
		}
		
		data "aiven_m3aggregator" "service" {
			service_name = aiven_m3aggregator.bar.service_name
			project = aiven_m3aggregator.bar.project

			depends_on = [aiven_m3aggregator.bar]
		}
		`, os.Getenv("AIVEN_PROJECT_NAME"), name, name, name)
}
