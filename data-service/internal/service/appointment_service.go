package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/domain"
)

const top10Limit = 10

var (
	ErrInvalidAppointmentID = errors.New("invalid appointment id")
	ErrInvalidPatientID     = errors.New("invalid patient id")
	ErrInvalidDoctorID      = errors.New("invalid doctor id")
	ErrInvalidStatus        = errors.New("invalid appointment status")
	ErrInvalidAppointmentAt = errors.New("invalid appointment time")
	ErrAppointmentInPast    = errors.New("appointment time cannot be in the past")
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrInvalidLimit         = errors.New("invalid limit")
	ErrInvalidOffset        = errors.New("invalid offset")
)

type AppointmentRepository interface {
	Create(ctx context.Context, a domain.Appointment) error
	Update(ctx context.Context, a domain.Appointment) error
	Search(ctx context.Context, filter domain.AppointmentSearchFilter) ([]domain.Appointment, error)
	CountByRange(ctx context.Context, from, to time.Time) ([]domain.AppointmentsPerDay, error)
	TopPatients(ctx context.Context, limit int) ([]domain.PatientAppointmentStat, error)
	TopDoctors(ctx context.Context, limit int) ([]domain.DoctorAppointmentStat, error)
}

type AppointmentService struct {
	repo AppointmentRepository
	now  func() time.Time
}

func NewAppointmentService(repo AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *AppointmentService) Create(ctx context.Context, appointment domain.Appointment) error {
	if err := validateAppointment(appointment); err != nil {
		return err
	}

	if appointment.AppointmentAt.Before(s.now()) {
		return ErrAppointmentInPast
	}

	return s.repo.Create(ctx, appointment)
}

func (s *AppointmentService) Update(ctx context.Context, appointment domain.Appointment) error {
	if err := validateAppointment(appointment); err != nil {
		return err
	}

	return s.repo.Update(ctx, appointment)
}

func (s *AppointmentService) Search(ctx context.Context, filter domain.AppointmentSearchFilter) ([]domain.Appointment, error) {
	if err := validateSearchFilter(filter); err != nil {
		return nil, err
	}

	items, err := s.repo.Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("search appointments from repository: %w", err)
	}

	return items, nil
}

func (s *AppointmentService) GetAppointmentsReportForDay(ctx context.Context, day time.Time) ([]domain.AppointmentsPerDay, error) {
	if day.IsZero() {
		return nil, ErrInvalidDateRange
	}

	from := beginningOfDay(day)
	to := from.AddDate(0, 0, 1)

	items, err := s.repo.CountByRange(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("get daily appointments report from repository: %w", err)
	}

	return items, nil
}

func (s *AppointmentService) GetTodayAppointmentsReport(ctx context.Context) ([]domain.AppointmentsPerDay, error) {
	return s.GetAppointmentsReportForDay(ctx, s.now())
}

func (s *AppointmentService) GetTopPatients(ctx context.Context, limit int) ([]domain.PatientAppointmentStat, error) {
	if err := validateLimit(limit); err != nil {
		return nil, err
	}

	items, err := s.repo.TopPatients(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("get top patients from repository: %w", err)
	}

	return items, nil
}

func (s *AppointmentService) GetTop10Patients(ctx context.Context) ([]domain.PatientAppointmentStat, error) {
	return s.GetTopPatients(ctx, top10Limit)
}

func (s *AppointmentService) GetTopDoctors(ctx context.Context, limit int) ([]domain.DoctorAppointmentStat, error) {
	if err := validateLimit(limit); err != nil {
		return nil, err
	}

	items, err := s.repo.TopDoctors(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("get top doctors from repository: %w", err)
	}

	return items, nil
}

func (s *AppointmentService) GetTop10Doctors(ctx context.Context) ([]domain.DoctorAppointmentStat, error) {
	return s.GetTopDoctors(ctx, top10Limit)
}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func validateAppointment(appointment domain.Appointment) error {
	if appointment.ID == uuid.Nil {
		return ErrInvalidAppointmentID
	}

	if appointment.PatientID == uuid.Nil {
		return ErrInvalidPatientID
	}

	if appointment.DoctorID == uuid.Nil {
		return ErrInvalidDoctorID
	}

	if appointment.AppointmentAt.IsZero() {
		return ErrInvalidAppointmentAt
	}

	switch appointment.Status {
	case domain.AppointmentStatusScheduled,
		domain.AppointmentStatusCancelled,
		domain.AppointmentStatusCompleted:
		return nil
	default:
		return ErrInvalidStatus
	}
}

func validateLimit(limit int) error {
	if limit <= 0 {
		return ErrInvalidLimit
	}

	return nil
}

func validateSearchFilter(filter domain.AppointmentSearchFilter) error {
	if err := validateLimit(filter.Limit); err != nil {
		return err
	}

	if filter.Offset < 0 {
		return ErrInvalidOffset
	}

	if filter.Status != nil {
		switch *filter.Status {
		case domain.AppointmentStatusScheduled,
			domain.AppointmentStatusCancelled,
			domain.AppointmentStatusCompleted:
		default:
			return ErrInvalidStatus
		}
	}

	if filter.Day != nil && filter.Day.IsZero() {
		return ErrInvalidDateRange
	}

	return nil
}
