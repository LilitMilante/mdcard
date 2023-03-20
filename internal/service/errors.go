package service

import (
	"errors"
)

var ErrAlreadyExists = errors.New("already exists")
var ErrNotFound = errors.New("not found")
var ErrInternal = errors.New("internal error")
