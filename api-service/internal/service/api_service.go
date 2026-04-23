package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/keshvan/trod-kafka-lab/api-service/internal/dto"
)

const dataServiceAPIPrefix = "/api/v1"

type Producer interface {
	PublishBatch(ctx context.Context, items []dto.AppointmentEvent) error
}

type DataClient interface {
	Get(ctx context.Context, path string, query url.Values) ([]byte, error)
}

type APIService struct {
	producer   Producer
	dataClient DataClient
}

func NewAPIService(producer Producer, dataClient DataClient) *APIService {
	return &APIService{
		producer:   producer,
		dataClient: dataClient,
	}
}

func (s *APIService) AddBatch(ctx context.Context, items []dto.AppointmentEvent) error {
	prepared := make([]dto.AppointmentEvent, len(items))
	copy(prepared, items)

	for i := range prepared {
		if prepared[i].Appointment.ID != "" {
			continue
		}

		if prepared[i].Action != "create" {
			return fmt.Errorf("appointment id is required for action %q", prepared[i].Action)
		}

		prepared[i].Appointment.ID = uuid.NewString()
	}

	return s.producer.PublishBatch(ctx, prepared)
}

func (s *APIService) SearchAppointments(ctx context.Context, query url.Values) ([]byte, error) {
	return s.dataClient.Get(ctx, dataServiceAPIPrefix+"/appointments", query)
}

func (s *APIService) GetDailyReport(ctx context.Context, query url.Values) ([]byte, error) {
	return s.dataClient.Get(ctx, dataServiceAPIPrefix+"/reports/daily", query)
}

func (s *APIService) GetTopDoctors(ctx context.Context, query url.Values) ([]byte, error) {
	return s.dataClient.Get(ctx, dataServiceAPIPrefix+"/reports/top-doctors", query)
}

func (s *APIService) GetTopPatients(ctx context.Context, query url.Values) ([]byte, error) {
	return s.dataClient.Get(ctx, dataServiceAPIPrefix+"/reports/top-patients", query)
}
