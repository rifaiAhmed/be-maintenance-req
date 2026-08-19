package helpers

import "errors"

var (
	ErrUnauthorized = errors.New("authentication required")
	ErrForbidden    = errors.New("you do not have permission to perform this action")
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource already exists")
	ErrInvalid      = errors.New("invalid request")
)
