package dal

import (
	"context"
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

// Patient methods

func (r *PatientRepository) PatientByPassportNumber(ctx context.Context, passNumber string) (entity.Patient, error) {
	return r.findPatientByColumn(ctx, "passport_number", passNumber)
}

func (r *PatientRepository) PatientByLogin(ctx context.Context, login string) (entity.Patient, error) {
	return r.findPatientByColumn(ctx, "login", login)
}

func (r *PatientRepository) PatientByID(ctx context.Context, id int64) (entity.Patient, error) {
	return r.findPatientByColumn(ctx, "id", id)
}

func (r *PatientRepository) CreatePatient(ctx context.Context, p entity.Patient) (entity.Patient, error) {
	q := `
INSERT INTO patients (full_name, data_of_born, address, phone_number, passport_number, login, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
`
	err := r.db.QueryRowContext(
		ctx,
		q,
		p.FullName,
		p.DateOfBorn,
		p.Address,
		p.PhoneNumber,
		p.PassportNumber,
		p.Login,
		p.Password,
		p.CreatedAt,
		p.UpdatedAt).
		Scan(&p.ID)

	return p, err
}

func (r *PatientRepository) Patients(ctx context.Context) ([]entity.Patient, error) {
	q := `
SELECT id, full_name, data_of_born, address, phone_number, passport_number, login, created_at
FROM patients 
`
	rows, err := r.db.QueryContext(ctx, q)
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

func (r *PatientRepository) UpdatePatient(ctx context.Context, id int64, p entity.Patient) error {
	q := `
UPDATE patients
SET full_name = $1, data_of_born = $2, address = $3, phone_number = $4, passport_number = $5, login = $6, updated_at = $7
WHERE id = $8
`

	_, err := r.db.ExecContext(
		ctx,
		q,
		p.FullName,
		p.DateOfBorn,
		p.Address,
		p.PhoneNumber,
		p.PassportNumber,
		p.Login,
		p.UpdatedAt,
		id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PatientRepository) DeletePatient(ctx context.Context, id int64) error {
	q := "DELETE FROM patients WHERE id = $1"

	_, err := r.db.ExecContext(ctx, q, id)

	err = r.deleteCard(ctx, id)

	return err
}

func (r *PatientRepository) findPatientByColumn(ctx context.Context, col string, value any) (entity.Patient, error) {
	var p entity.Patient

	q := "SELECT id, full_name, data_of_born, address, phone_number, passport_number, login, created_at, updated_at FROM patients"
	q = fmt.Sprintf("%s WHERE %s = $1", q, col)

	err := r.db.QueryRowContext(ctx, q, value).
		Scan(
			&p.ID,
			&p.FullName,
			&p.DateOfBorn,
			&p.Address,
			&p.PhoneNumber,
			&p.PassportNumber,
			&p.Login,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, fmt.Errorf("get patient by %s %v: %w", col, value, service.ErrNotFound)
		}

		return p, fmt.Errorf("get patient by %s %v: %w", col, value, err)
	}

	c, err := r.patientCard(ctx, p.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return p, nil
		}

		return p, fmt.Errorf("get patient %d card: %w", p.ID, err)
	}

	p.Card = &c

	return p, nil
}

// Card methods

func (r *PatientRepository) CreateCard(ctx context.Context, c entity.Card) (entity.Card, error) {
	q := `
INSERT INTO cards (patient_id, chronic_diseases, disability_group, blood_type, rh_factor, consultations, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
`
	err := r.db.QueryRowContext(
		ctx,
		q,
		c.PatientID,
		c.ChronicDiseases,
		c.DisabilityGroup,
		c.BloodType,
		c.RhFactor,
		c.Consultations,
		c.CreatedAt,
		c.UpdatedAt).Scan(&c.ID)

	return c, err
}

func (r *PatientRepository) CardByID(ctx context.Context, id int64) (entity.Card, error) {
	var c entity.Card

	q := `
SELECT id, patient_id, chronic_diseases, disability_group, blood_type, rh_factor, consultations, created_at, updated_at
FROM cards
WHERE id = $1
`

	err := r.db.QueryRowContext(ctx, q, id).
		Scan(
			&c.ID,
			&c.PatientID,
			&c.ChronicDiseases,
			&c.DisabilityGroup,
			&c.BloodType,
			&c.RhFactor,
			&c.Consultations,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c, service.ErrNotFound
		}

		return c, err
	}

	return c, nil
}

func (r *PatientRepository) UpdateCard(ctx context.Context, id int64, c entity.Card) error {
	q := `
UPDATE cards
SET patient_id = $1,  chronic_diseases = $2, disability_group = $3, blood_type = $4, rh_factor = $5, consultations = $6
WHERE id = $7
`

	_, err := r.db.ExecContext(
		ctx,
		q,
		&c.PatientID,
		&c.ChronicDiseases,
		&c.DisabilityGroup,
		&c.BloodType,
		&c.RhFactor,
		&c.Consultations,
		id,
	)

	return err
}

func (r *PatientRepository) patientCard(ctx context.Context, patientID int64) (entity.Card, error) {
	var c entity.Card

	q := `
SELECT id, patient_id, chronic_diseases, disability_group, blood_type, rh_factor, consultations, created_at, updated_at
FROM cards
WHERE patient_id = $1
`

	err := r.db.QueryRowContext(ctx, q, patientID).
		Scan(
			&c.ID,
			&c.PatientID,
			&c.ChronicDiseases,
			&c.DisabilityGroup,
			&c.BloodType,
			&c.RhFactor,
			&c.Consultations,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c, service.ErrNotFound
		}

		return c, err
	}

	return c, nil
}

func (r *PatientRepository) deleteCard(ctx context.Context, patientID int64) error {
	q := "DELETE FROM cards WHERE patient_id = $1"

	_, err := r.db.ExecContext(ctx, q, patientID)

	return err
}
