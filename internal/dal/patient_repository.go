package dal

import (
	"database/sql"
	"errors"
	"fmt"

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

func (r *PatientRepository) PatientByPassportNumber(passNumber string) (entity.Patient, error) {
	return r.findPatientByColumn("passport_number", passNumber)
}

func (r *PatientRepository) PatientByLogin(login string) (entity.Patient, error) {
	return r.findPatientByColumn("login", login)
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

func (r *PatientRepository) Patients() ([]entity.Patient, error) {
	q := `
SELECT id, full_name, data_of_born, address, phone_number, passport_number, login, created_at
FROM patients 
`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []entity.Patient

	for rows.Next() {
		var p entity.Patient

		err = rows.Scan(
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
			return nil, err
		}

		patients = append(patients, p)
	}

	return patients, nil
}

func (r *PatientRepository) findPatientByColumn(col string, value any) (entity.Patient, error) {
	var p entity.Patient

	q := "SELECT id, full_name, data_of_born, address, phone_number, passport_number, login, created_at FROM patients"
	q = fmt.Sprintf("%s WHERE %s = $1", q, col)

	err := r.db.QueryRow(q, value).
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
