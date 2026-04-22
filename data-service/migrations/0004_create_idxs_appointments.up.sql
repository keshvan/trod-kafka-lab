create index idx_appointments_appointment_at
    on appointments (appointment_at);

create index idx_appointments_patient_id
    on appointments (patient_id);

create index idx_appointments_doctor_id
    on appointments (doctor_id);

create index idx_appointments_status
    on appointments (status);

create unique index uq_appointments_doctor_slot
    on appointments (doctor_id, appointment_at)
    where status = 'scheduled';
