CREATE TABLE patients (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    data_of_born DATE NOT NULL,
    address JSONB NOT NULL,
    phone_number TEXT NOT NULL,
    passport_number TEXT UNIQUE NOT NULL,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

