package entity

import (
	"time"
)

type Patient struct {
	ID             int64
	FullName       string
	DateOfBorn     time.Time
	Address        Address
	PhoneNumber    string
	PassportNumber string
	Login          string
	CreatedAt      time.Time
}

type Address struct {
	Country   string
	City      string
	Street    string
	Building  string
	Apartment string
}
