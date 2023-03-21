package service

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrInternal      = errors.New("internal error")
	ErrUnauthorized  = errors.New("unauthorized")
)
