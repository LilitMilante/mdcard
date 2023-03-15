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

	CreateCard(ctx context.Context, c entity.Card) (entity.Card, error)
	CardByPatientPassportNumber(ctx context.Context, number string) (entity.Card, error)
	UpdateCard(ctx context.Context, id int64, c entity.Card) (entity.Card, error)
}

type PatientService struct {
	repo PatientRepository
}

func NewPatientService(repo PatientRepository) *PatientService {
	return &PatientService{repo: repo}
}

// Patient methods

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
	p.UpdatedAt = p.CreatedAt

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

// Card methods

func (s *PatientService) AddCard(ctx context.Context, c entity.Card) (entity.Card, error) {
	_, err := s.repo.PatientByID(ctx, c.PatientID)
	if err == nil {
		return c, fmt.Errorf("card with patient id %q: %w", c.PatientID, ErrAlreadyExists)
	}

	if !errors.Is(err, ErrNotFound) {
		return c, fmt.Errorf("card with patient %q: %w", c.PatientID, err)
	}

	c.CreatedAt = time.Now()
	c.UpdatedAt = c.CreatedAt

	card, err := s.repo.CreateCard(ctx, c)
	if err != nil {
		return card, fmt.Errorf("create card: %w", err)
	}

	return card, nil
}

func (s *PatientService) CardByPatientPassportNumber(ctx context.Context, number string) (entity.Card, error) {
	return s.repo.CardByPatientPassportNumber(ctx, number)
}

func (s *PatientService) UpdateCard(ctx context.Context, id int64, c entity.Card) (entity.Card, error) {
	_, err := s.repo.PatientByID(ctx, id)
	if err != nil {
		return c, fmt.Errorf("patient with id %d: %w", id, err)
	}

	c.UpdatedAt = time.Now()

	c, err = s.repo.UpdateCard(ctx, id, c)

	return c, err
}
