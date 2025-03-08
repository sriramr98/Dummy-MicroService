package utils

import (
	"errors"
	"log"
	"net/http"
)

type ApiHandler func(http.ResponseWriter, *http.Request) error

func ErrorHandler(handler ApiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			var apiErr ApplicationError
			if !errors.As(err, &apiErr) {
				log.Printf("Error: %v", err)
				apiErr = ApiError{
					StatusCode: http.StatusInternalServerError,
					Code:       ErrInternalServer,
					Message:    "Internal server error",
				}
			}

			WriteErrorResponse(w, apiErr)
		}
	}
}
