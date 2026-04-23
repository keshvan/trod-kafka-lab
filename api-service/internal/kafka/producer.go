package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/keshvan/trod-kafka-lab/api-service/internal/dto"
	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(broker, topic string) *Producer {
	return &Producer{
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(broker),
			Topic:        topic,
			RequiredAcks: kafkago.RequireOne,
			BatchTimeout: 200 * time.Millisecond,
		},
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) PublishBatch(ctx context.Context, items []dto.AppointmentEvent) error {
	messages := make([]kafkago.Message, 0, len(items))

	for _, item := range items {
		payload, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("marshal kafka event: %w", err)
		}

		messages = append(messages, kafkago.Message{
			Key:   []byte(item.Appointment.ID),
			Value: payload,
		})
	}

	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("write kafka messages: %w", err)
	}

	return nil
}
