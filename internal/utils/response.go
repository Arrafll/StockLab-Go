package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"`            // success / error
	Message string      `json:"message,omitempty"` // pesan deskriptif
	Data    interface{} `json:"data,omitempty"`    // payload
}

func RespondJSON(w http.ResponseWriter, statusCode int, status string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

func RespondError(w http.ResponseWriter, statusCode int, message string) {
	RespondJSON(w, statusCode, "error", message, nil)
}

func RespondSuccess(w http.ResponseWriter, data interface{}, message string) {
	RespondJSON(w, http.StatusOK, "success", message, data)
}
