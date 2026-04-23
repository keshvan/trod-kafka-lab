package dto

import "time"

type AppointmentEvent struct {
	Action      string         `json:"action"`
	Appointment AppointmentDTO `json:"appointment"`
}

type AppointmentDTO struct {
	ID            string    `json:"id"`
	PatientID     string    `json:"patient_id"`
	DoctorID      string    `json:"doctor_id"`
	AppointmentAt time.Time `json:"appointment_at"`
	Status        string    `json:"status"`
}

type AddBatchRequest struct {
	Items []AppointmentEvent `json:"items"`
}
