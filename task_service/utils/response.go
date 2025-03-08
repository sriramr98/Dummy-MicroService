package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteSuccessMessage[T any](w http.ResponseWriter, statusCode int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(map[string]any{"data": data, "success": true})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteErrorResponse(w http.ResponseWriter, apiErr ApplicationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.HttpStatusCode())

	if err := json.NewEncoder(w).Encode(apiErr.ErrMap()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
