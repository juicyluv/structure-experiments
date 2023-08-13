package apperror

import "fmt"

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
	return fmt.Sprintf("%s: %v", a.message, a.source)
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
