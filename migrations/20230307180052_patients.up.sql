CREATE TYPE address AS (
    country TEXT,
    city TEXT,
    street TEXT,
    building TEXT,
    apartment TEXT
);

CREATE TABLE patients (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    data_of_born DATE NOT NULL,
    address ADDRESS NOT NULL,
    phone_number TEXT NOT NULL,
    passport_number TEXT NOT NULL,
    login TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

