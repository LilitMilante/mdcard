package api

import (
	"context"
	"encoding/json"
	"errors"
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
	UpdatePatient(ctx context.Context, id int64, p entity.Patient) (entity.Patient, error)
	DeletePatient(ctx context.Context, id int64) error
}

type PatientHandler struct {
	srv Service
}

func NewPatientHandler(srv Service) *PatientHandler {
	return &PatientHandler{srv: srv}
}

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

	patient, err = h.srv.UpdatePatient(r.Context(), int64(id), patient)
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
