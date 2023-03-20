package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	PatientID int64
	CreatedAt time.Time
	ExpiredAt time.Time
}
