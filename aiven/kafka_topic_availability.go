package aiven

import (
	"fmt"
	"github.com/aiven/terraform-provider-aiven/pkg/cache"
	"log"
	"time"

	"github.com/aiven/aiven-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// KafkaTopicAvailabilityWaiter is used to refresh the Aiven Kafka Topic endpoints when
// provisioning.
type KafkaTopicAvailabilityWaiter struct {
	Client      *aiven.Client
	Project     string
	ServiceName string
	TopicName   string
}

// RefreshFunc will call the Aiven client and refresh it's state.
func (w *KafkaTopicAvailabilityWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		if w.Project == "" {
			return nil, "WRONG_INPUT", fmt.Errorf("project name of the kafka topic resource cannot be empty `%s`", w.Project)
		}

		if w.ServiceName == "" {
			return nil, "WRONG_INPUT", fmt.Errorf("service name of the kafka topic resource cannot be empty `%s`", w.ServiceName)
		}

		if w.TopicName == "" {
			return nil, "WRONG_INPUT", fmt.Errorf("topic name of the kafka topic resource cannot be empty `%s`", w.TopicName)
		}

		topicCache := cache.GetTopicCache()
		topic, ok := topicCache.LoadByTopicName(w.Project, w.ServiceName, w.TopicName)

		if !ok {
			list, err := w.Client.KafkaTopics.List(w.Project, w.ServiceName)

			for _, item := range list {
				log.Printf("[TRACE] got a topic `%s` from aiven API with the status `%s`", item.TopicName, item.State)
			}

			if err != nil {
				aivenError, ok := err.(aiven.Error)
				// Topic creation is asynchronous so it is possible for the creation call to
				// have completed successfully yet fetcing topic info fails with 404.
				if ok && aivenError.Status == 404 {
					return nil, "CONFIGURING", nil
				}
				// Getting topic info can sometimes temporarily fail with 501 and 502. Don't
				// treat that as fatal error but keep on retrying instead.
				if (ok && aivenError.Status == 501) || (ok && aivenError.Status == 502) {
					return nil, "CONFIGURING", nil
				}
				return nil, "", err
			}

			topicCache.StoreByProjectAndServiceName(w.Project, w.ServiceName, list)

			topic, ok = topicCache.LoadByTopicName(w.Project, w.ServiceName, w.TopicName)
			if !ok {
				return nil, "CONFIGURING", fmt.Errorf("topic `%s` for project `%s` and service `%s` not found",
					w.TopicName,
					w.Project,
					w.ServiceName)
			}
		}

		log.Printf("[DEBUG] Got `%s` state while waiting for topic `%s` to be up.", topic.State, w.TopicName)

		return topic, topic.State, nil
	}
}

// Conf sets up the configuration to refresh.
func (w *KafkaTopicAvailabilityWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{"CONFIGURING"},
		Target:     []string{"ACTIVE"},
		Refresh:    w.RefreshFunc(),
		Delay:      10 * time.Second,
		Timeout:    4 * time.Minute,
		MinTimeout: 1 * time.Second,
	}
}
