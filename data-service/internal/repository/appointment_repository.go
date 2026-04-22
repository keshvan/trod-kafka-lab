package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/domain"
)

type AppointmentRepository struct {
	pool *pgxpool.Pool
}

func NewAppointmentRepository(pool *pgxpool.Pool) *AppointmentRepository {
	return &AppointmentRepository{pool: pool}
}

func (r *AppointmentRepository) Create(ctx context.Context, a domain.Appointment) error {
	const query = `
		INSERT INTO appointments (
			id,
			patient_id,
			doctor_id,
			appointment_at,
			status
		)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, a.ID, a.PatientID, a.DoctorID, a.AppointmentAt, a.Status)
	if err != nil {
		return fmt.Errorf("create appointment: %w", err)
	}

	return nil
}

func (r *AppointmentRepository) Update(ctx context.Context, a domain.Appointment) error {
	const query = `
		UPDATE appointments
		SET
			patient_id = $2,
			doctor_id = $3,
			appointment_at = $4,
			status = $5
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, a.ID, a.PatientID, a.DoctorID, a.AppointmentAt, a.Status)
	if err != nil {
		return fmt.Errorf("update appointment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("update appointment: %w", domain.ErrAppointmentNotFound)
	}

	return nil
}

func (r *AppointmentRepository) Search(ctx context.Context, filter domain.AppointmentSearchFilter) ([]domain.Appointment, error) {
	var (
		query strings.Builder
		args  []any
	)

	query.WriteString(`
		SELECT
			id,
			patient_id,
			doctor_id,
			appointment_at,
			status,
			created_at
		FROM appointments
		WHERE 1 = 1
	`)

	if filter.PatientID != nil {
		args = append(args, *filter.PatientID)
		query.WriteString(fmt.Sprintf(" AND patient_id = $%d", len(args)))
	}

	if filter.DoctorID != nil {
		args = append(args, *filter.DoctorID)
		query.WriteString(fmt.Sprintf(" AND doctor_id = $%d", len(args)))
	}

	if filter.Status != nil {
		args = append(args, *filter.Status)
		query.WriteString(fmt.Sprintf(" AND status = $%d", len(args)))
	}

	if filter.Day != nil {
		from := beginningOfDay(*filter.Day)
		to := from.AddDate(0, 0, 1)

		args = append(args, from)
		query.WriteString(fmt.Sprintf(" AND appointment_at >= $%d", len(args)))

		args = append(args, to)
		query.WriteString(fmt.Sprintf(" AND appointment_at < $%d", len(args)))
	}

	query.WriteString(" ORDER BY appointment_at DESC, created_at DESC")

	args = append(args, filter.Limit)
	query.WriteString(fmt.Sprintf(" LIMIT $%d", len(args)))

	args = append(args, filter.Offset)
	query.WriteString(fmt.Sprintf(" OFFSET $%d", len(args)))

	rows, err := r.pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("search appointments: %w", err)
	}
	defer rows.Close()

	result := make([]domain.Appointment, 0)
	for rows.Next() {
		var item domain.Appointment
		if err = rows.Scan(
			&item.ID,
			&item.PatientID,
			&item.DoctorID,
			&item.AppointmentAt,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan appointment: %w", err)
		}
		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate appointments: %w", err)
	}

	return result, nil
}

func (r *AppointmentRepository) CountByRange(ctx context.Context, from, to time.Time) ([]domain.AppointmentsPerDay, error) {
	const query = `
		SELECT DATE(appointment_at) AS day, COUNT(*) AS count
		FROM appointments
		WHERE appointment_at >= $1
		  AND appointment_at < $2
		GROUP BY DATE(appointment_at)
		ORDER BY day
	`

	rows, err := r.pool.Query(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("count appointments by day: %w", err)
	}
	defer rows.Close()

	result := make([]domain.AppointmentsPerDay, 0)
	for rows.Next() {
		var item domain.AppointmentsPerDay
		if err = rows.Scan(&item.Day, &item.Count); err != nil {
			return nil, fmt.Errorf("scan appointments by day: %w", err)
		}
		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate appointments by day: %w", err)
	}

	return result, nil
}

func (r *AppointmentRepository) TopPatients(ctx context.Context, limit int) ([]domain.PatientAppointmentStat, error) {
	const query = `
		SELECT
			p.id,
			p.first_name || ' ' || p.last_name AS full_name,
			COUNT(a.id) AS count
		FROM appointments a
		JOIN patients p ON p.id = a.patient_id
		GROUP BY p.id, p.first_name, p.last_name
		ORDER BY count DESC, full_name ASC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("top patients: %w", err)
	}
	defer rows.Close()

	result := make([]domain.PatientAppointmentStat, 0)
	for rows.Next() {
		var item domain.PatientAppointmentStat
		if err = rows.Scan(&item.PatientID, &item.FullName, &item.Count); err != nil {
			return nil, fmt.Errorf("scan top patients: %w", err)
		}
		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate top patients: %w", err)
	}

	return result, nil
}

func (r *AppointmentRepository) TopDoctors(ctx context.Context, limit int) ([]domain.DoctorAppointmentStat, error) {
	const query = `
		SELECT
			d.id,
			d.first_name || ' ' || d.last_name AS full_name,
			d.speciality,
			COUNT(a.id) AS count
		FROM appointments a
		JOIN doctors d ON d.id = a.doctor_id
		GROUP BY d.id, d.first_name, d.last_name, d.speciality
		ORDER BY count DESC, full_name ASC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("top doctors: %w", err)
	}
	defer rows.Close()

	result := make([]domain.DoctorAppointmentStat, 0)
	for rows.Next() {
		var item domain.DoctorAppointmentStat
		if err = rows.Scan(&item.DoctorID, &item.FullName, &item.Speciality, &item.Count); err != nil {
			return nil, fmt.Errorf("scan top doctors: %w", err)
		}
		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate top doctors: %w", err)
	}

	return result, nil
}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
