package models

import (
	"fmt"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// IMPLEMENTS "error" INTERFACE
func (a APIError) Error() string {
	return fmt.Sprint(a.Message)
}

func NewAPIError(statusCode int, msg string) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    msg,
	}
}
