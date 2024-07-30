package models

type APIError struct {
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(status int, msg string) *APIError {
	return &APIError{
		StatusCode: status,
		Message:    msg,
	}
}
