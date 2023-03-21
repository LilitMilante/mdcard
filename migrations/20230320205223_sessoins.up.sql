CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id),
    created_at TIMESTAMPTZ NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL
);
