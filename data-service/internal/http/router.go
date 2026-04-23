package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/http/handler"
)

func NewRouter(apiHandler *handler.DataHandler) nethttp.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/appointments", apiHandler.SearchAppointments)

		r.Route("/reports", func(r chi.Router) {
			r.Get("/daily", apiHandler.GetDailyReport)
			r.Get("/top-doctors", apiHandler.GetTopDoctors)
			r.Get("/top-patients", apiHandler.GetTopPatients)
		})
	})

	return r
}
