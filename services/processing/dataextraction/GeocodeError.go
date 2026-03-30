package dataextraction

import "fmt"

type ErrorCode string

const (
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeInternalError  ErrorCode = "INTERNAL_ERROR"
	ErrCodeNetworkError   ErrorCode = "NETWORK_ERROR"
	ErrCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrCodeParsingError   ErrorCode = "PARSING_ERROR"
	ErrCodeFiltered       ErrorCode = "FILTERED"
)

type GeocodeError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *GeocodeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *GeocodeError) Unwrap() error {
	return e.Err
}

func NewGeocodeError(code ErrorCode, message string, err error) *GeocodeError {
	return &GeocodeError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
