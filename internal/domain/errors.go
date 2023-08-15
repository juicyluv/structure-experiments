package domain

import "errors"

var (
	ErrRequired     = errors.New("required")
	ErrInvalidValue = errors.New("invalid value")
)
