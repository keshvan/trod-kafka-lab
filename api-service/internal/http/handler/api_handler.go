package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/keshvan/trod-kafka-lab/api-service/internal/dto"
)

type Service interface {
	AddBatch(ctx context.Context, items []dto.AppointmentEvent) error
	SearchAppointments(ctx context.Context, query url.Values) ([]byte, error)
	GetDailyReport(ctx context.Context, query url.Values) ([]byte, error)
	GetTopDoctors(ctx context.Context, query url.Values) ([]byte, error)
	GetTopPatients(ctx context.Context, query url.Values) ([]byte, error)
}

type ApiHandler struct {
	service Service
}

func NewApiHandler(service Service) *ApiHandler {
	return &ApiHandler{
		service: service,
	}
}

func (h *ApiHandler) AddBatch(w http.ResponseWriter, r *http.Request) {
	var req dto.AddBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		http.Error(w, "items must not be empty", http.StatusBadRequest)
		return
	}

	if err := h.service.AddBatch(r.Context(), req.Items); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"status": "accepted",
		"count":  len(req.Items),
	})
}

func (h *ApiHandler) SearchAppointments(w http.ResponseWriter, r *http.Request) {
	body, err := h.service.SearchAppointments(r.Context(), r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeRawJSON(w, http.StatusOK, body)
}

func (h *ApiHandler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	body, err := h.service.GetDailyReport(r.Context(), r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeRawJSON(w, http.StatusOK, body)
}

func (h *ApiHandler) GetTopDoctors(w http.ResponseWriter, r *http.Request) {
	body, err := h.service.GetTopDoctors(r.Context(), r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeRawJSON(w, http.StatusOK, body)
}

func (h *ApiHandler) GetTopPatients(w http.ResponseWriter, r *http.Request) {
	body, err := h.service.GetTopPatients(r.Context(), r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeRawJSON(w, http.StatusOK, body)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeRawJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
