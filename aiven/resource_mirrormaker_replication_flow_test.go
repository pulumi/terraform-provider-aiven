package aiven

import (
	"fmt"
	"github.com/aiven/aiven-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccAivenMirrorMakerReplicationFlow_basic(t *testing.T) {
	resourceName := "aiven_mirrormaker_replication_flow.foo"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAivenMirrorMakerReplicationFlowResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMirrorMakerReplicationFlowResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAivenMirrorMakerReplicationFlowAttributes("data.aiven_mirrormaker_replication_flow.flow"),
					resource.TestCheckResourceAttr(resourceName, "project", fmt.Sprintf("test-acc-pr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-mm-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "source_cluster", "source"),
					resource.TestCheckResourceAttr(resourceName, "target_cluster", "target"),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
				),
			},
		},
	})
}

func testAccCheckAivenMirrorMakerReplicationFlowResourceDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*aiven.Client)

	// loop through the resources in state, verifying each kafka mirror maker
	// replication flow is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aiven_mirrormaker_replication_flow" {
			continue
		}

		project, serviceName, sourceCluster, targetCluster := splitResourceID4(rs.Primary.ID)

		s, err := c.Services.Get(project, serviceName)
		if err != nil {
			if err.(aiven.Error).Status != 404 {
				return err
			}
			return nil
		}

		if s.Type == "kafka_mirrormaker" {
			f, err := c.KafkaMirrorMakerReplicationFlow.Get(project, serviceName, sourceCluster, targetCluster)
			if err != nil {
				if err.(aiven.Error).Status != 404 {
					return err
				}
			}

			if f != nil {
				return fmt.Errorf("kafka mirror maker replication flow still exists, id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccMirrorMakerReplicationFlowResource(name string) string {
	return fmt.Sprintf(`
		resource "aiven_project" "foo" {
			project = "test-acc-pr-%s"
			card_id="%s"	
		}
		
		resource "aiven_service" "source" {
			project = aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-4"
			service_name = "test-acc-sr-source-%s"
			service_type = "kafka"
			maintenance_window_dow = "monday"
			maintenance_window_time = "10:00:00"
			
			kafka_user_config {
				kafka_version = "2.4"
				kafka {
				  group_max_session_timeout_ms = 70000
				  log_retention_bytes = 1000000000
				}
			}
		}
		
		resource "aiven_kafka_topic" "source" {
			project = aiven_project.foo.project
			service_name = aiven_service.source.service_name
			topic_name = "test-acc-topic-a-%s"
			partitions = 3
			replication = 2
		}

		resource "aiven_service" "target" {
			project = aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-4"
			service_name = "test-acc-sr-target-%s"
			service_type = "kafka"
			maintenance_window_dow = "monday"
			maintenance_window_time = "10:00:00"
			
			kafka_user_config {
				kafka_version = "2.4"
				kafka {
				  group_max_session_timeout_ms = 70000
				  log_retention_bytes = 1000000000
				}
			}
		}
		
		resource "aiven_kafka_topic" "target" {
			project = aiven_project.foo.project
			service_name = aiven_service.target.service_name
			topic_name = "test-acc-topic-b-%s"
			partitions = 3
			replication = 2
		}

		resource "aiven_service" "mm" {
			project = aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "startup-4"
			service_name = "test-acc-sr-mm-%s"
			service_type = "kafka_mirrormaker"
			
			kafka_mirrormaker_user_config {
				ip_filter = ["0.0.0.0/0"]

				kafka_mirrormaker {
					refresh_groups_interval_seconds = 600
					refresh_topics_enabled = true
					refresh_topics_interval_seconds = 600
				}
			}
		}

		resource "aiven_service_integration" "bar" {
			project = aiven_project.foo.project
			integration_type = "kafka_mirrormaker"
			source_service_name = aiven_service.source.service_name
			destination_service_name = aiven_service.mm.service_name
	
			kafka_mirrormaker_user_config {
				cluster_alias = "source"
			}
		}

		resource "aiven_service_integration" "i2" {
			project = aiven_project.foo.project
			integration_type = "kafka_mirrormaker"
			source_service_name = aiven_service.target.service_name
			destination_service_name = aiven_service.mm.service_name
	
			kafka_mirrormaker_user_config {
				cluster_alias = "target"
			}
		}

		resource "aiven_mirrormaker_replication_flow" "foo" {
			project = aiven_project.foo.project
			service_name = aiven_service.mm.service_name
			source_cluster = "source"
			target_cluster = "target"
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

		data "aiven_mirrormaker_replication_flow" "flow" {
			project = aiven_project.foo.project
			service_name = aiven_service.mm.service_name
			source_cluster = aiven_mirrormaker_replication_flow.foo.source_cluster
			target_cluster = aiven_mirrormaker_replication_flow.foo.target_cluster
		}
		`, name, os.Getenv("AIVEN_CARD_ID"), name, name, name, name, name)
}

func testAccCheckAivenMirrorMakerReplicationFlowAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["project"] == "" {
			return fmt.Errorf("expected to get a project name from Aiven")
		}

		if a["service_name"] == "" {
			return fmt.Errorf("expected to get a service_name from Aiven")
		}

		if a["source_cluster"] != "source" {
			return fmt.Errorf("expected to get a source_cluster from Aiven")
		}

		if a["target_cluster"] != "target" {
			return fmt.Errorf("expected to get target_cluster from Aiven")
		}

		if a["enable"] != "true" {
			return fmt.Errorf("expected to get a correct enable from Aiven")
		}

		return nil
	}
}
