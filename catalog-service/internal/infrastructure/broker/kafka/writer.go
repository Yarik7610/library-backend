package kafka

import (
	"context"

	sharedKafka "github.com/Yarik7610/library-backend-common/broker/kafka"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/tracing"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type OtelWriter struct {
	writer      *kafka.Writer
	serviceName string
}

func NewOtelWriter(config *config.Config, topic string) *OtelWriter {
	return &OtelWriter{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(sharedKafka.KAFKA_NODE_1_ADDRESS, sharedKafka.KAFKA_NODE_2_ADDRESS, sharedKafka.KAFKA_NODE_3_ADDRESS),
			Topic: topic,
		},
		serviceName: config.ServiceName,
	}
}

func (w *OtelWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	tracer := otel.Tracer(w.serviceName)

	ctx, span := tracer.Start(ctx, "kafka.produce",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", w.writer.Topic),
		),
	)
	defer span.End()

	// Use key-value buffer because HTTP headers aren't working in Kafka (raw bytes allowed only)
	carrier := propagation.MapCarrier{}
	// Enrich carrier with current ctx
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	for i := range msgs {
		for k, v := range carrier {
			msgs[i].Headers = append(msgs[i].Headers, kafka.Header{Key: k, Value: []byte(v)})
		}
	}

	if err := w.writer.WriteMessages(ctx, msgs...); err != nil {
		tracing.Error(span, err)
		return err
	}
	return nil
}
