package domain

import "errors"

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEditConflict       = errors.New("edit conflict")
)
