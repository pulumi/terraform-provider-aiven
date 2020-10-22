package aiven

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"
)

func TestAccAiven_m3db(t *testing.T) {
	resourceName := "aiven_m3db.bar"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAivenServiceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccM3DBResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAivenServiceCommonAttributes("data.aiven_m3db.service"),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "project", os.Getenv("AIVEN_PROJECT_NAME")),
					resource.TestCheckResourceAttr(resourceName, "service_type", "m3db"),
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

func testAccM3DBResource(name string) string {
	return fmt.Sprintf(`
		data "aiven_project" "foo" {
			project = "%s"
		}
		
		resource "aiven_m3db" "bar" {
			project = data.aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-8"
			service_name = "test-acc-sr-%s"
			maintenance_window_dow = "monday"
			maintenance_window_time = "10:00:00"
			
			m3db_user_config {
				m3_version = 0.15
							
				namespaces {
					name = "test-acc-%s"
					type = "unaggregated"
				}
			}
		}

		resource "aiven_pg" "pg1" {
			project = data.aiven_project.foo.project
			cloud_name = "google-europe-west1"
			service_name = "test-acc-sr-pg-%s"
			plan = "startup-4"

			pg_user_config {
				pg_version = 12.4
			}
		}
		
		resource "aiven_service_integration" "int-m3db-pg" {
			project = data.aiven_project.foo.project
			integration_type = "metrics"
			source_service_name = aiven_pg.pg1.service_name
			destination_service_name = aiven_m3db.bar.service_name
		}

		resource "aiven_grafana" "grafana1" {
			project = data.aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "startup-4"
			service_name = "test-acc-sr-g-%s"

			grafana_user_config {
				alerting_enabled = true
				
				public_access {
					grafana = true
				}
			}
		}
		
		resource "aiven_service_integration" "int-grafana-m3db" {
			project = data.aiven_project.foo.project
			integration_type = "dashboard"
			source_service_name = aiven_grafana.grafana1.service_name
			destination_service_name = aiven_m3db.bar.service_name
		}
		
		data "aiven_m3db" "service" {
			service_name = aiven_m3db.bar.service_name
			project = aiven_m3db.bar.project
		}
		`, os.Getenv("AIVEN_PROJECT_NAME"), name, name, name, name)
}
