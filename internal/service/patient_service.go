package service

import (
	"errors"
	"fmt"
	"time"

	"medical-card/internal/entity"
)

type PatientRepository interface {
	PatientByPassportNumber(passNumber string) (entity.Patient, error)
	PatientByLogin(login string) (entity.Patient, error)
	CreatePatient(p entity.Patient) (entity.Patient, error)
	Patients() ([]entity.Patient, error)
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
		return p, fmt.Errorf("patient with passport %q: %w", p.PassportNumber, ErrAlreadyExists)
	}

	if !errors.Is(err, ErrNotFound) {
		return p, fmt.Errorf("patient with passport %q: %w", p.PassportNumber, err)
	}

	_, err = s.repo.PatientByLogin(p.Login)
	if err == nil {
		return p, fmt.Errorf("patient with login %q: %w", p.Login, ErrAlreadyExists)
	}

	if !errors.Is(err, ErrNotFound) {
		return p, fmt.Errorf("patient with login %q: %w", p.Login, err)
	}

	p.CreatedAt = time.Now()

	p, err = s.repo.CreatePatient(p)
	if err != nil {
		return p, fmt.Errorf("create patient: %w", err)
	}

	return p, nil
}

func (s *PatientService) Patients() ([]entity.Patient, error) {
	return s.repo.Patients()
}

func (s *PatientService) PatientByPassportNumber(passNumber string) (entity.Patient, error) {
	return s.repo.PatientByPassportNumber(passNumber)
}
