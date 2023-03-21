package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"medical-card/internal/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type PatientRepository interface {
	PatientByPassportNumber(ctx context.Context, passNumber string) (entity.Patient, error)
	PatientByLogin(ctx context.Context, login string) (entity.Patient, error)
	PatientByID(ctx context.Context, id int64) (entity.Patient, error)
	CreatePatient(ctx context.Context, p entity.Patient) (entity.Patient, error)
	Patients(ctx context.Context) ([]entity.Patient, error)
	UpdatePatient(ctx context.Context, id int64, p entity.Patient) error
	DeletePatient(ctx context.Context, id int64) error

	CreateCard(ctx context.Context, c entity.Card) (entity.Card, error)
	CardByID(ctx context.Context, id int64) (entity.Card, error)
	UpdateCard(ctx context.Context, id int64, c entity.Card) error

	Login(ctx context.Context, sess entity.Session) error
	SessionByID(ctx context.Context, id string) (entity.Session, error)
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

	p.Password, err = s.hashPassword(p.Password)
	if err != nil {
		return p, err
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

func (s *PatientService) PatientByLogin(ctx context.Context, login string) (entity.Patient, error) {
	return s.repo.PatientByLogin(ctx, login)
}

func (s *PatientService) UpdatePatient(ctx context.Context, id int64, p entity.Patient) error {
	_, err := s.repo.PatientByID(ctx, id)
	if err != nil {
		return fmt.Errorf("patient with id %d: %w", id, err)
	}

	p.UpdatedAt = time.Now()

	err = s.repo.UpdatePatient(ctx, id, p)
	if err != nil {
		return err
	}

	return nil
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
		return c, fmt.Errorf("patient with id %d: %w", c.PatientID, ErrAlreadyExists)
	}

	if !errors.Is(err, ErrNotFound) {
		return c, fmt.Errorf("patient with id  %d: %w", c.PatientID, err)
	}

	c.CreatedAt = time.Now()
	c.UpdatedAt = c.CreatedAt

	card, err := s.repo.CreateCard(ctx, c)
	if err != nil {
		return card, fmt.Errorf("create card: %w", err)
	}

	return card, nil
}

func (s *PatientService) UpdateCard(ctx context.Context, id int64, c entity.Card) error {
	_, err := s.repo.CardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("card with id %d: %w", id, err)
	}

	c.UpdatedAt = time.Now()

	return s.repo.UpdateCard(ctx, id, c)
}

func (s *PatientService) hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

// Session

func (s *PatientService) Login(ctx context.Context, patientID int64) (entity.Session, error) {
	sess := entity.Session{
		ID:        uuid.New(),
		PatientID: patientID,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 1),
	}

	err := s.repo.Login(ctx, sess)
	if err != nil {
		return entity.Session{}, err
	}

	return sess, nil
}

func (s *PatientService) PatientBySessionID(ctx context.Context, ssid string) (entity.Patient, error) {
	session, err := s.repo.SessionByID(ctx, ssid)
	if err != nil {
		return entity.Patient{}, err
	}

	if time.Now().After(session.ExpiredAt) {
		return entity.Patient{}, fmt.Errorf("%w: session expired", ErrUnauthorized)
	}

	patient, err := s.repo.PatientByID(ctx, session.PatientID)
	if err != nil {
		return patient, err
	}

	return patient, nil
}
