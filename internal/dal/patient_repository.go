package dal

import (
	"database/sql"
	"errors"

	"medical-card/internal/entity"
	"medical-card/internal/service"
)

var _ service.PatientRepository = (*PatientRepository)(nil)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{
		db: db,
	}
}

func (r *PatientRepository) PatientByPassportNumber(n string) (entity.Patient, error) {
	var p entity.Patient

	q := `
SELECT id, full_name, data_of_born, address, phone_number, passport_number, login, created_at
FROM patients 
WHERE passport_number = $1
`

	err := r.db.QueryRow(q, n).
		Scan(
			&p.ID,
			&p.FullName,
			&p.DateOfBorn,
			&p.Address,
			&p.PhoneNumber,
			&p.PassportNumber,
			&p.Login,
			&p.CreatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, service.ErrNotFound
		}

		return p, err
	}

	return p, nil
}

func (r *PatientRepository) CreatePatient(p entity.Patient) (entity.Patient, error) {
	q := `
INSERT INTO patients (full_name, data_of_born, address, phone_number, passport_number, login, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
`
	err := r.db.QueryRow(
		q,
		p.FullName,
		p.DateOfBorn,
		p.Address,
		p.PhoneNumber,
		p.PassportNumber,
		p.Login,
		p.CreatedAt).
		Scan(&p.ID)
	if err != nil {
		return p, err
	}

	return p, nil
}
