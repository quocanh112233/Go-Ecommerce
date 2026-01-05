package errors

import "errors"

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrRecordNotFound      = errors.New("record not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInternalServerError = errors.New("internal server error")
)
