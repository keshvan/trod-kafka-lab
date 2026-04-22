package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/domain"
)

type AppointmentService interface {
	Search(ctx context.Context, filter domain.AppointmentSearchFilter) ([]domain.Appointment, error)
	GetAppointmentsReportForDay(ctx context.Context, day time.Time) ([]domain.AppointmentsPerDay, error)
	GetTodayAppointmentsReport(ctx context.Context) ([]domain.AppointmentsPerDay, error)
	GetTopDoctors(ctx context.Context, limit int) ([]domain.DoctorAppointmentStat, error)
	GetTop10Doctors(ctx context.Context) ([]domain.DoctorAppointmentStat, error)
	GetTopPatients(ctx context.Context, limit int) ([]domain.PatientAppointmentStat, error)
	GetTop10Patients(ctx context.Context) ([]domain.PatientAppointmentStat, error)
}

type DataHandler struct {
	service AppointmentService
}

func NewDataHandler(service AppointmentService) *DataHandler {
	return &DataHandler{service: service}
}

func (h *DataHandler) SearchAppointments(w http.ResponseWriter, r *http.Request) {
	filter, err := parseSearchFilter(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items, err := h.service.Search(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *DataHandler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	dayRaw := r.URL.Query().Get("day")

	var (
		items []domain.AppointmentsPerDay
		err   error
	)

	if dayRaw == "" {
		items, err = h.service.GetTodayAppointmentsReport(r.Context())
	} else {
		day, parseErr := time.Parse("2006-01-02", dayRaw)
		if parseErr != nil {
			http.Error(w, "invalid day, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		items, err = h.service.GetAppointmentsReportForDay(r.Context(), day)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *DataHandler) GetTopDoctors(w http.ResponseWriter, r *http.Request) {
	limit, ok := parseOptionalLimit(r)
	if !ok {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	var (
		items []domain.DoctorAppointmentStat
		err   error
	)

	if limit == 0 {
		items, err = h.service.GetTop10Doctors(r.Context())
	} else {
		items, err = h.service.GetTopDoctors(r.Context(), limit)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *DataHandler) GetTopPatients(w http.ResponseWriter, r *http.Request) {
	limit, ok := parseOptionalLimit(r)
	if !ok {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	var (
		items []domain.PatientAppointmentStat
		err   error
	)

	if limit == 0 {
		items, err = h.service.GetTop10Patients(r.Context())
	} else {
		items, err = h.service.GetTopPatients(r.Context(), limit)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func parseOptionalLimit(r *http.Request) (int, bool) {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return 0, true
	}

	limit, err := strconv.Atoi(raw)
	if err != nil || limit <= 0 {
		return 0, false
	}

	return limit, true
}

func parseSearchFilter(r *http.Request) (domain.AppointmentSearchFilter, error) {
	query := r.URL.Query()
	filter := domain.AppointmentSearchFilter{
		Limit:  50,
		Offset: 0,
	}

	if raw := query.Get("patient_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return filter, err
		}
		filter.PatientID = &id
	}

	if raw := query.Get("doctor_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return filter, err
		}
		filter.DoctorID = &id
	}

	if raw := query.Get("status"); raw != "" {
		status := domain.AppointmentStatus(raw)
		filter.Status = &status
	}

	if raw := query.Get("day"); raw != "" {
		day, err := time.Parse("2006-01-02", raw)
		if err != nil {
			return filter, err
		}
		filter.Day = &day
	}

	if raw := query.Get("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil {
			return filter, err
		}
		filter.Limit = limit
	}

	if raw := query.Get("offset"); raw != "" {
		offset, err := strconv.Atoi(raw)
		if err != nil {
			return filter, err
		}
		filter.Offset = offset
	}

	return filter, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
