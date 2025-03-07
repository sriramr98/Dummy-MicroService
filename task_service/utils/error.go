package utils

import "fmt"

type ApiErrorCode int

var (
	ErrInvalidRequestPayload = ApiErrorCode(1000)
	ErrInternalServer        = ApiErrorCode(1001)
)

type ApplicationError interface {
	ErrMap() map[string]interface{}
	HttpStatusCode() int
}

type ApiError struct {
	StatusCode int
	Code       ApiErrorCode `json:"code"`
	Message    string       `json:"message"`
}

type ValidationError struct {
	StatusCode int
	Code       ApiErrorCode      `json:"code"`
	Errors     map[string]string `json:"errors"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

func (e ApiError) ErrMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}
}

func (e ApiError) HttpStatusCode() int {
	return e.StatusCode
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%+v", e.Errors)
}

func (e ValidationError) ErrMap() map[string]interface{} {
	return map[string]interface{}{
		"code":   e.Code,
		"errors": e.Errors,
	}
}

func (e ValidationError) HttpStatusCode() int {
	return e.StatusCode
}
