package service

import (
	"errors"
	"fmt"
	"time"

	"medical-card/internal/entity"
)

type PatientRepository interface {
	PatientByPassportNumber(pn string) (entity.Patient, error)
	CreatePatient(p entity.Patient) (entity.Patient, error)
}

type PatientService struct {
	repo PatientRepository
}

func NewPatientService(repo PatientRepository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) AddPatient(p entity.Patient) (entity.Patient, error) {
	_, err := s.repo.PatientByPassportNumber(p.PassportNumber)
	if err == nil {
		return p, fmt.Errorf("get patient with passport %q: %w", p.PassportNumber, ErrAlreadyExists)
	}

	if !errors.Is(err, ErrNotFound) {
		return p, fmt.Errorf("get patient with passport %q: %w", p.PassportNumber, err)
	}

	p.CreatedAt = time.Now()

	p, err = s.repo.CreatePatient(p)
	if err != nil {
		return p, fmt.Errorf("create patient: %w", err)
	}

	return p, nil
}