package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Card struct {
	ID              int64           `json:"id,omitempty"`
	PatientID       int64           `json:"patient_id,omitempty"`
	ChronicDiseases ChronicDiseases `json:"chronic_diseases,omitempty"`
	DisabilityGroup *int            `json:"disability_group,omitempty"`
	BloodType       int             `json:"blood_type,omitempty"`
	RhFactor        bool            `json:"rh_factor,omitempty"`
	Consultations   Consultations   `json:"consultations,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ChronicDiseases []string

type Consultation struct {
	ID              int64     `json:"id,omitempty"`
	CardID          int64     `json:"card_id,omitempty"`
	DoctorID        string    `json:"doctor_id,omitempty"`
	FullName        string    `json:"full_name,omitempty"`
	Complaints      string    `json:"complaints,omitempty"`
	Descriptions    string    `json:"descriptions,omitempty"`
	Recommendations string    `json:"recommendations,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Consultations []Consultation

func (c ChronicDiseases) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ChronicDiseases) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

func (c Consultations) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Consultations) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}
