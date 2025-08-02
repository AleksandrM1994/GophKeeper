package kafka

import "github.com/segmentio/kafka-go"

const (
	GophKeeperPrivateDataSavedTopic = "gophkeeper.privateDataSaved"
)

func getTopics() []kafka.TopicConfig {
	return []kafka.TopicConfig{
		{
			Topic:             GophKeeperPrivateDataSavedTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
}
