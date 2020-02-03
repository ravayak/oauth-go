package errors

import (
	"errors"
	"net/http"
)

type RestError struct {
	Message string `json:"message"`
	Status  int    `json:"code"`
	Error   string `json:"error"`
}

func NewBadRequestError(message string, err error) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   err.Error(),
	}
}

func NewNotFoundError(message string, err error) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   err.Error(),
	}
}

func NewInternalServerError(message string, err error) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   err.Error(),
	}
}

func NewError(msg string) error {
	return errors.New(msg)
}
