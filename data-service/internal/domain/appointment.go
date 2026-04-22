package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrAppointmentNotFound = errors.New("appointment not found")

type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
	ID            uuid.UUID
	PatientID     uuid.UUID
	DoctorID      uuid.UUID
	AppointmentAt time.Time
	Status        AppointmentStatus
	CreatedAt     time.Time
}

type AppointmentSearchFilter struct {
	PatientID *uuid.UUID
	DoctorID  *uuid.UUID
	Status    *AppointmentStatus
	Day       *time.Time
	Limit     int
	Offset    int
}
