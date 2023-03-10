package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"medical-card/internal/api"
	"medical-card/internal/app"
	"medical-card/internal/dal"
	"medical-card/internal/entity"
	service2 "medical-card/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddPatient(t *testing.T) {
	c := app.Config{
		Database: app.DBConfig{
			Host:     "localhost",
			Port:     8181,
			Name:     "mdcard",
			User:     "dev",
			Password: "dev",
		},
	}
	db, err := app.NewPostgresClient(c.Database)
	require.NoError(t, err)
	repo := dal.NewPatientRepository(db)
	service := service2.NewPatientService(repo)
	handler := api.NewPatientHandler(service)

	payload := entity.Patient{
		FullName:   "test test",
		DateOfBorn: time.Now(),
		Address: entity.Address{
			Country:   "Belarus",
			City:      "Vitebsk",
			Street:    "Smolenskaya",
			Building:  "11",
			Apartment: "158",
		},
		PhoneNumber:    "1234567890",
		PassportNumber: "BM1234567",
		Login:          "@Nasta",
	}

	jsonPayload, err := json.Marshal(payload)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/patients", bytes.NewReader(jsonPayload))

	handler.AddPatient(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp entity.Patient

	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.NotZero(t, resp.ID)
	assert.False(t, resp.CreatedAt.IsZero())
}
