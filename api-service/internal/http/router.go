package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler interface {
	AddBatch(w http.ResponseWriter, r *http.Request)
	SearchAppointments(w http.ResponseWriter, r *http.Request)
	GetDailyReport(w http.ResponseWriter, r *http.Request)
	GetTopDoctors(w http.ResponseWriter, r *http.Request)
	GetTopPatients(w http.ResponseWriter, r *http.Request)
}

func NewRouter(apiHandler Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/appointments/batch", apiHandler.AddBatch)
		r.Get("/appointments", apiHandler.SearchAppointments)

		r.Route("/reports", func(r chi.Router) {
			r.Get("/daily", apiHandler.GetDailyReport)
			r.Get("/top-doctors", apiHandler.GetTopDoctors)
			r.Get("/top-patients", apiHandler.GetTopPatients)
		})
	})

	return r
}
