package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/domain"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/kafka/dto"
	kafkago "github.com/segmentio/kafka-go"
)

type AppointmentCommand interface {
	Create(ctx context.Context, appointment domain.Appointment) error
	Update(ctx context.Context, appointment domain.Appointment) error
}

type Consumer struct {
	reader  *kafkago.Reader
	service AppointmentCommand
}

func NewConsumer(broker, topic, groupID string, service AppointmentCommand) *Consumer {
	return &Consumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
		service: service,
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("read kafka message: %w", err)
		}

		if err := c.handleMessage(ctx, msg.Value); err != nil {
			log.Printf("kafka message skipped: %v", err)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, raw []byte) error {
	var event dto.AppointmentEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		return fmt.Errorf("unmarshal kafka event: %w", err)
	}

	appointment, err := toDomain(event.Appointment)
	if err != nil {
		return err
	}

	switch event.Action {
	case "create":
		return c.service.Create(ctx, appointment)
	case "update":
		return c.service.Update(ctx, appointment)
	default:
		return fmt.Errorf("unknown action: %s", event.Action)
	}
}

func toDomain(src dto.AppointmentDTO) (domain.Appointment, error) {
	id, err := uuid.Parse(src.ID)
	if err != nil {
		return domain.Appointment{}, fmt.Errorf("parse appointment id: %w", err)
	}

	patientID, err := uuid.Parse(src.PatientID)
	if err != nil {
		return domain.Appointment{}, fmt.Errorf("parse patient id: %w", err)
	}

	doctorID, err := uuid.Parse(src.DoctorID)
	if err != nil {
		return domain.Appointment{}, fmt.Errorf("parse doctor id: %w", err)
	}

	return domain.Appointment{
		ID:            id,
		PatientID:     patientID,
		DoctorID:      doctorID,
		AppointmentAt: src.AppointmentAt,
		Status:        domain.AppointmentStatus(src.Status),
	}, nil
}
