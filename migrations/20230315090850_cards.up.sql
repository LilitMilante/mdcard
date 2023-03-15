CREATE TABLE cards (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id),
    chronic_diseases JSONB NOT NULL,
    disability_group INT,
    blood_type INT NOT NULL,
    rh_factor BOOLEAN NOT NULL,
    consultations JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
