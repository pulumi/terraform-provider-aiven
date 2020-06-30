package aiven

import (
	"fmt"
	"github.com/aiven/aiven-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"os"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("aiven_kafka_topic", &resource.Sweeper{
		Name: "aiven_kafka_topic",
		F:    sweepKafkaTopics,
	})
}

func sweepKafkaTopics(region string) error {
	client, err := sharedClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*aiven.Client)

	projects, err := conn.Projects.List()
	if err != nil {
		return fmt.Errorf("error retrieving a list of projects : %s", err)
	}

	for _, project := range projects {
		if strings.Contains(project.Name, "test-acc-") {
			services, err := conn.Services.List(project.Name)
			if err != nil {
				return fmt.Errorf("error retrieving a list of services for a project `%s`: %s", project.Name, err)
			}

			for _, service := range services {
				if service.Type != "kafka" {
					continue
				}

				topics, err := conn.KafkaTopics.List(project.Name, service.Name)
				if err != nil {
					log.Printf("[ERROR] error retrieving a list of kafka topics for a service `%s`: %s", service.Name, err)
					continue
				}

				for _, topic := range topics {
					err = conn.KafkaTopics.Delete(project.Name, service.Name, topic.TopicName)
					if err != nil {
						return fmt.Errorf("error destroying kafka topic %s during sweep: %s", topic.TopicName, err)
					}
				}
			}
		}
	}

	return nil
}

func TestAccAivenKafkaTopic_basic(t *testing.T) {
	resourceName := "aiven_kafka_topic.foo"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAivenKafkaTopicResourceDestroy,
		Steps: []resource.TestStep{
			// basic Kafka Topic test
			{
				Config: testAccKafkaTopicResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAivenKafkaTopicAttributes("data.aiven_kafka_topic.topic"),
					resource.TestCheckResourceAttr(resourceName, "project", fmt.Sprintf("test-acc-pr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "topic_name", fmt.Sprintf("test-acc-topic-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "partitions", "3"),
					resource.TestCheckResourceAttr(resourceName, "replication", "2"),
					resource.TestCheckResourceAttr(resourceName, "termination_protection", "false"),
				),
			},
			// custom TF client timeouts test
			{
				Config: testAccKafkaTopicCustomTimeoutsResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAivenKafkaTopicAttributes("data.aiven_kafka_topic.topic"),
					resource.TestCheckResourceAttr(resourceName, "project", fmt.Sprintf("test-acc-pr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "topic_name", fmt.Sprintf("test-acc-topic-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "partitions", "3"),
					resource.TestCheckResourceAttr(resourceName, "replication", "2"),
					resource.TestCheckResourceAttr(resourceName, "termination_protection", "false"),
				),
			},
			// termination protection test
			{
				Config:                    testAccKafkaTopicTerminationProtectionResource(rName),
				PreventPostDestroyRefresh: true,
				ExpectNonEmptyPlan:        true,
				PlanOnly:                  true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project", fmt.Sprintf("test-acc-pr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "topic_name", fmt.Sprintf("test-acc-topic-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "partitions", "3"),
					resource.TestCheckResourceAttr(resourceName, "replication", "2"),
					resource.TestCheckResourceAttr(resourceName, "termination_protection", "true"),
				),
			},
		},
	})
}

func TestAccAivenKafkaTopic_100topics(t *testing.T) {
	if os.Getenv("AIVEN_ACC_LONG") == "" {
		t.Skip("Acceptance tests skipped unless env AIVEN_ACC_LONG set")
	}
	t.Parallel()

	resourceName := "aiven_kafka_topic.foo"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAivenKafkaTopicResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafka101TopicResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAivenKafkaTopicAttributes("data.aiven_kafka_topic.topic"),
					resource.TestCheckResourceAttr(resourceName, "project", fmt.Sprintf("test-acc-pr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "topic_name", fmt.Sprintf("test-acc-topic-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "partitions", "3"),
					resource.TestCheckResourceAttr(resourceName, "replication", "2"),
				),
			},
		},
	})
}

func testAccKafka101TopicResource(name string) string {
	s := testAccKafkaTopicResource(name)

	// add extra 100 Kafka topics to test caching layer and creation waiter functionality
	for i := 1; i < 100; i++ {
		s += fmt.Sprintf(`
			resource "aiven_kafka_topic" "foo%s" {
				project = aiven_project.foo.project
				service_name = aiven_service.bar.service_name
				topic_name = "test-acc-topic-%s"
				partitions = 3
				replication = 2
			}
		`,
			acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum),
			acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum))
	}

	return s
}

func testAccKafkaTopicResource(name string) string {
	return fmt.Sprintf(`
		resource "aiven_project" "foo" {
			project = "test-acc-pr-%s"
			card_id="%s"	
		}

		resource "aiven_service" "bar" {
			project = aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-4"
			service_name = "test-acc-sr-%s"
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
		
		resource "aiven_kafka_topic" "foo" {
			project = aiven_project.foo.project
			service_name = aiven_service.bar.service_name
			topic_name = "test-acc-topic-%s"
			partitions = 3
			replication = 2
		}

		data "aiven_kafka_topic" "topic" {
			project = aiven_kafka_topic.foo.project
			service_name = aiven_kafka_topic.foo.service_name
			topic_name = aiven_kafka_topic.foo.topic_name
		}
		`, name, os.Getenv("AIVEN_CARD_ID"), name, name)
}

func testAccKafkaTopicCustomTimeoutsResource(name string) string {
	return fmt.Sprintf(`
		resource "aiven_project" "foo" {
			project = "test-acc-pr-%s"
			card_id="%s"	
		}

		resource "aiven_service" "bar" {
			project = aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-4"
			service_name = "test-acc-sr-%s"
			service_type = "kafka"
			maintenance_window_dow = "monday"
			maintenance_window_time = "10:00:00"

			timeouts {
				create = "25m"
				update = "20m"
			}
			
			kafka_user_config {
				kafka_version = "2.4"
			}
		}
		
		resource "aiven_kafka_topic" "foo" {
			project = aiven_project.foo.project
			service_name = aiven_service.bar.service_name
			topic_name = "test-acc-topic-%s"
			partitions = 3
			replication = 2

			timeouts {
				create = "5m"
				read = "5m"
			}
		}

		data "aiven_kafka_topic" "topic" {
			project = aiven_kafka_topic.foo.project
			service_name = aiven_kafka_topic.foo.service_name
			topic_name = aiven_kafka_topic.foo.topic_name
		}
		`, name, os.Getenv("AIVEN_CARD_ID"), name, name)
}

func testAccKafkaTopicTerminationProtectionResource(name string) string {
	return fmt.Sprintf(`
		resource "aiven_project" "foo" {
			project = "test-acc-pr-%s"
			card_id="%s"	
		}

		resource "aiven_service" "bar" {
			project = aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-4"
			service_name = "test-acc-sr-%s"
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
		
		resource "aiven_kafka_topic" "foo" {
			project = aiven_project.foo.project
			service_name = aiven_service.bar.service_name
			topic_name = "test-acc-topic-%s"
			partitions = 3
			replication = 2
			termination_protection = true
		}

		data "aiven_kafka_topic" "topic" {
			project = aiven_kafka_topic.foo.project
			service_name = aiven_kafka_topic.foo.service_name
			topic_name = aiven_kafka_topic.foo.topic_name
		}
		`, name, os.Getenv("AIVEN_CARD_ID"), name, name)
}

func testAccCheckAivenKafkaTopicAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		log.Printf("[DEBUG] kafka topic attributes %v", a)

		if a["project"] == "" {
			return fmt.Errorf("expected to get a project name from Aiven")
		}

		if a["service_name"] == "" {
			return fmt.Errorf("expected to get a service_name from Aiven")
		}

		if a["topic_name"] == "" {
			return fmt.Errorf("expected to get a topic_name from Aiven")
		}

		if a["partitions"] == "" {
			return fmt.Errorf("expected to get partitions from Aiven")
		}

		if a["replication"] == "" {
			return fmt.Errorf("expected to get a replication from Aiven")
		}

		return nil
	}
}

func testAccCheckAivenKafkaTopicResourceDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*aiven.Client)

	// loop through the resources in state, verifying each kafka topic is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aiven_kafka_topic" {
			continue
		}

		project, serviceName, topicName := splitResourceID3(rs.Primary.ID)

		_, err := c.Services.Get(project, serviceName)
		if err != nil {
			if err.(aiven.Error).Status != 404 {
				return err
			}
		}

		t, err := c.KafkaTopics.Get(project, serviceName, topicName)
		if err != nil {
			if err.(aiven.Error).Status != 404 {
				return err
			}
		}

		if t != nil {
			return fmt.Errorf("kafka topic (%s) still exists, id %s", topicName, rs.Primary.ID)
		}
	}

	return nil
}
