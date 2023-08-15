package apperror

import (
	"errors"
	"fmt"
)

type ErrorType string

const (
	ErrorTypeNotFound       ErrorType = "NotFound"
	ErrorTypeNotUnique      ErrorType = "NotUnique"
	ErrorTypeInvalidRequest ErrorType = "InvalidRequest"
)

type AppError struct {
	message string
	field   string
	source  error
	etype   ErrorType
}

func (a AppError) Error() string {
	return fmt.Errorf("%s: %w", a.message, a.source).Error()
}

func (a AppError) Is(err error) bool {
	var ae AppError
	if !errors.As(err, &ae) {
		return errors.Is(err, a.source)
	}

	return a.etype == ae.etype && a.message == ae.message && a.field == ae.field && errors.Is(err, a.source)
}

func (a AppError) ErrorType() ErrorType {
	return a.etype
}

func (a AppError) Field() string {
	return a.field
}

func (a AppError) Source() error {
	return a.source
}

func NewInvalidRequestError(err error, msg string, field string) AppError {
	return AppError{
		message: msg,
		source:  err,
		field:   field,
		etype:   ErrorTypeInvalidRequest,
	}
}

func NewNotFoundError(err error, msg string, field string) AppError {
	return AppError{
		message: msg,
		source:  err,
		field:   field,
		etype:   ErrorTypeNotFound,
	}
}

func NewNotUniqueError(err error, msg string, field string) AppError {
	return AppError{
		message: msg,
		source:  err,
		field:   field,
		etype:   ErrorTypeNotUnique,
	}
}
