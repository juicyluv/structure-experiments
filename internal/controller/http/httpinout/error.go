package httpinout

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/juicyluv/structure-experiments/internal/pkg/logger"
	"github.com/juicyluv/structure-experiments/internal/service/apperror"
)

func InternalError(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, "Internal error", "", w, r, http.StatusInternalServerError)
}

func BadRequest(err error, msg string, field string, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, msg, field, w, r, http.StatusBadRequest)
}

func NotUnique(err error, msg string, field string, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, msg, field, w, r, http.StatusConflict)
}

func NotFound(err error, msg string, field string, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, msg, field, w, r, http.StatusNotFound)
}

func httpRespondWithError(err error, msg string, field string, w http.ResponseWriter, r *http.Request, status int) {
	logger.Get().Error(msg, "url", r.URL.Path, "status", status, "err", err, "field", field)

	resp := ErrorResponse{msg, field, status}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func JsonError(err error, w http.ResponseWriter, r *http.Request) {
	appError := apperror.NewInvalidRequestError(err, "Invalid JSON provided.", "")
	RespondWithError(appError, w, r)
}

type ErrorResponse struct {
	Error      string `json:"error"`
	Field      string `json:"field,omitempty"`
	httpStatus int
}

func RespondWithError(err error, w http.ResponseWriter, r *http.Request) {
	var appError apperror.AppError

	ok := errors.As(err, &appError)
	if !ok {
		InternalError(err, w, r)
		return
	}

	switch appError.ErrorType() {
	case apperror.ErrorTypeInvalidRequest:
		BadRequest(appError.Source(), appError.Error(), appError.Field(), w, r)
	case apperror.ErrorTypeNotUnique:
		NotUnique(appError.Source(), appError.Error(), appError.Field(), w, r)
	case apperror.ErrorTypeNotFound:
		NotFound(appError.Source(), appError.Error(), appError.Field(), w, r)
	default:
		InternalError(appError.Source(), w, r)
	}
}
