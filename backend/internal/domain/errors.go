package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input data")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrConflict     = errors.New("resource already exists")
	ErrInternal     = errors.New("internal server error")
)
