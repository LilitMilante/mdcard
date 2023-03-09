package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"medical-card/internal/entity"
	"medical-card/internal/service"
)

type Service interface {
	AddPatient(p entity.Patient) (entity.Patient, error)
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

	patient, err = h.srv.AddPatient(patient)
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
