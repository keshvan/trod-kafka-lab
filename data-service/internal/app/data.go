package app

import (
	"errors"
	"fmt"
	nethttp "net/http"

	"github.com/keshvan/trod-kafka-lab/data-service/internal/config"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/db"
	datahttp "github.com/keshvan/trod-kafka-lab/data-service/internal/http"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/http/handler"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/kafka"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/repository"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/service"
)

type Data struct {
	Config *config.Config

	Database *db.Database

	AppointmentRepository *repository.AppointmentRepository
	AppointmentService    *service.AppointmentService

	DataHandler *handler.DataHandler
	Router      nethttp.Handler
	HTTPServer  *datahttp.Server

	Consumer *kafka.Consumer
}

func New(cfg *config.Config) (*Data, error) {
	if cfg.KafkaHost == "" || cfg.KafkaPort == "" {
		return nil, fmt.Errorf("kafka host and port are required")
	}

	database, err := db.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	appointmentRepository := repository.NewAppointmentRepository(database.Pool)
	appointmentService := service.NewAppointmentService(appointmentRepository)

	dataHandler := handler.NewDataHandler(appointmentService)
	router := datahttp.NewRouter(dataHandler)
	httpServer := datahttp.NewServer(cfg, router)

	kafkaAddress := fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)
	consumer := kafka.NewConsumer(kafkaAddress, cfg.KafkaTopic, cfg.KafkaGroup, appointmentService)

	return &Data{
		Config: cfg,

		Database: database,

		AppointmentRepository: appointmentRepository,
		AppointmentService:    appointmentService,

		DataHandler: dataHandler,
		Router:      router,
		HTTPServer:  httpServer,

		Consumer: consumer,
	}, nil
}

func (d *Data) Shutdown() error {
	var errs []error

	if d.Consumer != nil {
		if err := d.Consumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close kafka consumer: %w", err))
		}
	}

	if d.Database != nil {
		d.Database.Close()
	}

	return errors.Join(errs...)
}
