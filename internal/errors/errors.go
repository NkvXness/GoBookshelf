package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeBadRequest     ErrorType = "BAD_REQUEST"
	ErrorTypeInternalServer ErrorType = "INTERNAL_SERVER_ERROR"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewNotFoundError(message string) AppError {
	return AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
	}
}

func NewBadRequestError(message string) AppError {
	return AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
	}
}

func NewInternalServerError(message string, err error) AppError {
	return AppError{
		Type:    ErrorTypeInternalServer,
		Message: message,
		Err:     err,
	}
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	appErr, ok := err.(AppError)
	if !ok {
		appErr = NewInternalServerError("An unexpected error occurred", err)
	}

	statusCode := http.StatusInternalServerError
	switch appErr.Type {
	case ErrorTypeNotFound:
		statusCode = http.StatusNotFound
	case ErrorTypeBadRequest:
		statusCode = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"error": "%s", "message": "%s"}`, appErr.Type, appErr.Message)
}
