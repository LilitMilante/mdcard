package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Patient struct {
	ID             int64     `json:"id,omitempty"`
	FullName       string    `json:"full_name,omitempty"`
	DateOfBorn     time.Time `json:"date_of_born"`
	Address        Address   `json:"address"`
	PhoneNumber    string    `json:"phone_number,omitempty"`
	PassportNumber string    `json:"passport_number,omitempty"`
	Login          string    `json:"login,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type Address struct {
	Country   string `json:"country,omitempty"`
	City      string `json:"city,omitempty"`
	Street    string `json:"street,omitempty"`
	Building  string `json:"building,omitempty"`
	Apartment string `json:"apartment,omitempty"`
}

func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Address) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
