package domain

import (
	"time"

	"github.com/google/uuid"
)

type AppointmentsPerDay struct {
	Day   time.Time
	Count int
}

type PatientAppointmentStat struct {
	PatientID uuid.UUID
	FullName  string
	Count     int
}

type DoctorAppointmentStat struct {
	DoctorID   uuid.UUID
	FullName   string
	Speciality string
	Count      int
}
