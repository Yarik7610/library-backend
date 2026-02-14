package kafka

import (
	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"
	"github.com/segmentio/kafka-go"
)

func NewReader(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{sharedKafka.KAFKA_NODE_1_ADDRESS, sharedKafka.KAFKA_NODE_2_ADDRESS, sharedKafka.KAFKA_NODE_3_ADDRESS},
		Topic:   topic,
		GroupID: groupID,
	})
}
