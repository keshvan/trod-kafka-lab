package service

import (
	"context"
	"net/url"

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
	return s.producer.PublishBatch(ctx, items)
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
