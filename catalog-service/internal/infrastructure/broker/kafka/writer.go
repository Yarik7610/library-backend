package kafka

import (
	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"
	"github.com/segmentio/kafka-go"
)

func NewWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(sharedKafka.KAFKA_NODE_1_ADDRESS, sharedKafka.KAFKA_NODE_2_ADDRESS, sharedKafka.KAFKA_NODE_3_ADDRESS),
		Topic: topic,
	}
}
