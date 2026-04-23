package app

import (
	"errors"
	"fmt"
	nethttp "net/http"

	"github.com/keshvan/trod-kafka-lab/api-service/internal/client/dataservice"
	"github.com/keshvan/trod-kafka-lab/api-service/internal/config"
	apihttp "github.com/keshvan/trod-kafka-lab/api-service/internal/http"
	"github.com/keshvan/trod-kafka-lab/api-service/internal/http/handler"
	"github.com/keshvan/trod-kafka-lab/api-service/internal/kafka"
	"github.com/keshvan/trod-kafka-lab/api-service/internal/service"
)

type API struct {
	Config *config.Config

	DataClient *dataservice.DataClient
	Producer   *kafka.Producer
	APIService *service.APIService

	APIHandler *handler.ApiHandler
	Router     nethttp.Handler
	HTTPServer *apihttp.Server
}

func New(cfg *config.Config) (*API, error) {
	if cfg.KafkaHost == "" || cfg.KafkaPort == "" {
		return nil, fmt.Errorf("kafka host and port are required")
	}

	dataServiceURL := fmt.Sprintf("http://%s:%s", cfg.DataClientHost, cfg.DataClientPort)
	dataClient := dataservice.NewDataClient(dataServiceURL)

	kafkaAddress := fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)
	producer := kafka.NewProducer(kafkaAddress, cfg.KafkaTopic)

	apiService := service.NewAPIService(producer, dataClient)
	apiHandler := handler.NewApiHandler(apiService)
	router := apihttp.NewRouter(apiHandler)
	httpServer := apihttp.NewServer(cfg, router)

	return &API{
		Config: cfg,

		DataClient: dataClient,
		Producer:   producer,
		APIService: apiService,

		APIHandler: apiHandler,
		Router:     router,
		HTTPServer: httpServer,
	}, nil
}

func (a *API) Shutdown() error {
	var errs []error

	if a.Producer != nil {
		if err := a.Producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close kafka producer: %w", err))
		}
	}

	return errors.Join(errs...)
}
