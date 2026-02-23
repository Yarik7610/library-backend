package kafka

import (
	"context"

	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/config"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type OtelReader struct {
	reader      *kafka.Reader
	serviceName string
}

func NewOtelReader(config *config.Config, topic, groupID string) *OtelReader {
	return &OtelReader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{sharedKafka.KAFKA_NODE_1_ADDRESS, sharedKafka.KAFKA_NODE_2_ADDRESS, sharedKafka.KAFKA_NODE_3_ADDRESS},
			Topic:   topic,
			GroupID: groupID,
		}),
		serviceName: config.ServiceName,
	}
}

func (r *OtelReader) FetchMessage(ctx context.Context) (kafka.Message, context.Context, trace.Span, error) {
	message, err := r.reader.FetchMessage(ctx)
	if err != nil {
		return message, ctx, trace.SpanFromContext(ctx), err
	}

	// Use key-value buffer because HTTP headers aren't working in Kafka (raw bytes allowed only)
	carrier := propagation.MapCarrier{}
	for _, h := range message.Headers {
		carrier[h.Key] = string(h.Value)
	}

	// Regain parent context from another microservice that came here
	parentCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer(r.serviceName)
	spanCtx, span := tracer.Start(parentCtx, "kafka.consume",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.source", r.reader.Config().Topic),
		),
	)

	return message, spanCtx, span, err
}

func (r *OtelReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return r.reader.CommitMessages(ctx, msgs...)
}

func (r *OtelReader) Close() error {
	return r.reader.Close()
}
