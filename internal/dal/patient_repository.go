package dal

import (
	"database/sql"

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

func (r *PatientRepository) PatientByPassportNumber(pn string) (entity.Patient, error) {
	return entity.Patient{}, nil
}

func (r *PatientRepository) CreatePatient(p entity.Patient) (entity.Patient, error) {
	return entity.Patient{}, nil
}
