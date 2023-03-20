package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Patient struct {
	ID                int64     `json:"id"`
	FullName          string    `json:"full_name"`
	DateOfBorn        time.Time `json:"date_of_born"`
	Address           Address   `json:"address"`
	PhoneNumber       string    `json:"phone_number"`
	PassportNumber    string    `json:"passport_number"`
	Login             string    `json:"login"`
	Password          string    `json:"password,omitempty"`
	EncryptedPassword string    `json:"-"`
	Card              *Card     `json:"card"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
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

func (p *Patient) Sanitize() {
	p.Password = ""
}

func (p *Patient) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.EncryptedPassword), []byte(password)) == nil
}
