CREATE TABLE sessions (
    id UUID NOT NULL ,
    patient_id BIGINT NOT NULL REFERENCES patients(id),
    created_at TIMESTAMPTZ NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL
);
