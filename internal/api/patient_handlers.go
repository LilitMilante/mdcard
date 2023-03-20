package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"medical-card/internal/entity"
	"medical-card/internal/service"

	"github.com/gorilla/mux"
)

type Service interface {
	AddPatient(ctx context.Context, p entity.Patient) (entity.Patient, error)
	Patients(ctx context.Context) ([]entity.Patient, error)
	PatientByPassportNumber(ctx context.Context, passNumber string) (entity.Patient, error)
	PatientByLogin(ctx context.Context, login string) (entity.Patient, error)
	UpdatePatient(ctx context.Context, id int64, p entity.Patient) error
	DeletePatient(ctx context.Context, id int64) error

	AddCard(ctx context.Context, c entity.Card) (entity.Card, error)
	UpdateCard(ctx context.Context, id int64, c entity.Card) error

	Login(ctx context.Context, patientID int64) (entity.Session, error)
}

type PatientHandler struct {
	srv Service
}

func NewPatientHandler(srv Service) *PatientHandler {
	return &PatientHandler{srv: srv}
}

// Patient methods

func (h *PatientHandler) AddPatient(w http.ResponseWriter, r *http.Request) {
	var patient entity.Patient

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	patient, err = h.srv.AddPatient(r.Context(), patient)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			SendErr(w, http.StatusConflict, err)
			return
		}

		SendErr(w, http.StatusInternalServerError, err)
		return
	}

	patient.Sanitize()

	SendJSON(w, patient)
}

func (h *PatientHandler) Patients(w http.ResponseWriter, r *http.Request) {
	patients, err := h.srv.Patients(r.Context())
	if err != nil {
		SendErr(w, http.StatusInternalServerError, err)
		return
	}

	SendJSON(w, patients)
}

func (h *PatientHandler) PatientByPassportNumber(w http.ResponseWriter, r *http.Request) {
	passNumber := mux.Vars(r)["passport_number"]

	patient, err := h.srv.PatientByPassportNumber(r.Context(), passNumber)
	if err != nil {
		SendErr(w, http.StatusInsufficientStorage, err)
		return
	}

	SendJSON(w, patient)
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	var patient entity.Patient

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	err = h.srv.UpdatePatient(r.Context(), int64(id), patient)
	if err != nil {
		SendErr(w, http.StatusInsufficientStorage, err)
		return
	}

	SendJSON(w, patient)
}

func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	err = h.srv.DeletePatient(r.Context(), int64(id))
	if err != nil {
		SendErr(w, http.StatusInsufficientStorage, err)
		return
	}

	SendJSON(w, id)
}

// Card Methods

func (h *PatientHandler) AddCard(w http.ResponseWriter, r *http.Request) {
	var card entity.Card

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	card, err = h.srv.AddCard(r.Context(), card)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			SendErr(w, http.StatusConflict, err)
			return
		}

		SendErr(w, http.StatusInternalServerError, err)
		return
	}

	SendJSON(w, card)
}

func (h *PatientHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	var card entity.Card

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	err = h.srv.UpdateCard(r.Context(), int64(id), card)
	if err != nil {
		SendErr(w, http.StatusInsufficientStorage, err)
		return
	}

	SendJSON(w, card)
}

// Sessions

func (h *PatientHandler) Login(w http.ResponseWriter, r *http.Request) {
	var patient entity.Patient

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	patient, err = h.srv.PatientByLogin(r.Context(), patient.Login)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) || !patient.ComparePassword(patient.Password) {
			SendErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		SendErr(w, http.StatusInternalServerError, err)
		return
	}

	sess, err := h.srv.Login(r.Context(), patient.ID)
	if err != nil {
		SendErr(w, http.StatusInternalServerError, err)
		return
	}

	cookie := &http.Cookie{
		Name:    "session",
		Value:   sess.ID.String(),
		Expires: sess.ExpiredAt,
	}

	http.SetCookie(w, cookie)

	SendJSON(w, "added session!")
}
