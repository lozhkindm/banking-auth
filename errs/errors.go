package errs

import "net/http"

type AppError struct {
	Code    int    `json:",omitempty"`
	Message string `json:"message"`
}

func (e AppError) AsMessage() *AppError {
	return &AppError{
		Message: e.Message,
	}
}

func NewDatabaseError() *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "unexpected database error",
	}
}

func NewAuthorizationError(msg string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}
