package database

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal server error")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized - please authenticate")
	ErrForbidden    = errors.New("forbidden")
)
