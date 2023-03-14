package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"medical-card/internal/entity"
)

type PatientRepository interface {
	PatientByPassportNumber(ctx context.Context, passNumber string) (entity.Patient, error)
	PatientByLogin(ctx context.Context, login string) (entity.Patient, error)
	PatientByID(ctx context.Context, id int64) (entity.Patient, error)
	CreatePatient(ctx context.Context, p entity.Patient) (entity.Patient, error)
	Patients(ctx context.Context) ([]entity.Patient, error)
	UpdatePatient(ctx context.Context, id int64, p entity.Patient) (entity.Patient, error)
	DeletePatient(ctx context.Context, id int64) error
}

type PatientService struct {
	repo PatientRepository
}

func NewPatientService(repo PatientRepository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) AddPatient(ctx context.Context, p entity.Patient) (entity.Patient, error) {
	_, err := s.repo.PatientByPassportNumber(ctx, p.PassportNumber)
	if err == nil {
		return p, fmt.Errorf("patient with passport %q: %w", p.PassportNumber, ErrAlreadyExists)
	}

	if !errors.Is(err, ErrNotFound) {
		return p, fmt.Errorf("patient with passport %q: %w", p.PassportNumber, err)
	}

	_, err = s.repo.PatientByLogin(ctx, p.Login)
	if err == nil {
		return p, fmt.Errorf("patient with login %q: %w", p.Login, ErrAlreadyExists)
	}

	if !errors.Is(err, ErrNotFound) {
		return p, fmt.Errorf("patient with login %q: %w", p.Login, err)
	}

	p.CreatedAt = time.Now()

	p, err = s.repo.CreatePatient(ctx, p)
	if err != nil {
		return p, fmt.Errorf("create patient: %w", err)
	}

	return p, nil
}

func (s *PatientService) Patients(ctx context.Context) ([]entity.Patient, error) {
	return s.repo.Patients(ctx)
}

func (s *PatientService) PatientByPassportNumber(ctx context.Context, passNumber string) (entity.Patient, error) {
	return s.repo.PatientByPassportNumber(ctx, passNumber)
}

func (s *PatientService) UpdatePatient(ctx context.Context, id int64, p entity.Patient) (entity.Patient, error) {
	_, err := s.repo.PatientByID(ctx, id)
	if err != nil {
		return p, fmt.Errorf("patient with id %d: %w", id, err)
	}

	p.UpdatedAt = time.Now()

	p, err = s.repo.UpdatePatient(ctx, id, p)
	if err != nil {
		return entity.Patient{}, err
	}

	return p, nil
}

func (s *PatientService) DeletePatient(ctx context.Context, id int64) error {
	_, err := s.repo.PatientByID(ctx, id)
	if err != nil {
		return fmt.Errorf("patient with id %d: %w", id, err)
	}

	return s.repo.DeletePatient(ctx, id)
}
