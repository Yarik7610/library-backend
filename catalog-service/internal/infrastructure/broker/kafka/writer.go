package kafka

import (
	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"
	"github.com/segmentio/kafka-go"
)

func NewWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(sharedKafka.KAFKA_NODE_1_ADDR, sharedKafka.KAFKA_NODE_2_ADDR, sharedKafka.KAFKA_NODE_3_ADDR),
		Topic: topic,
	}
}
