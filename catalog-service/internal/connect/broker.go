package connect

import (
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(sharedconstants.KAFKA_NODE_1_ADDR, sharedconstants.KAFKA_NODE_2_ADDR, sharedconstants.KAFKA_NODE_3_ADDR),
		Topic: topic,
	}
}
