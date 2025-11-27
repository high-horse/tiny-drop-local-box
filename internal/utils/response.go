package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code int `json:"code"`
	Success bool `jaon:"success"`
	Data interface{} `json:"data,omitempty"`
}

type SuccessResponse struct {
	Error string `json:"error"`
	Code int `json:"code"`
	Success bool `jaon:"success"`
	Data interface{} `json:"data,omitempty"`
}

func SendError(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Error: message,
		Code: statusCode,
		Success: false,
		Data: data,
	}

	json.NewEncoder(w).Encode(resp)
}


func SendSuccess(w http.ResponseWriter, statusCode int, message string, data any) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	if message == "" {
		message = "Success"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := SuccessResponse{
		Error: message,
		Code: statusCode,
		Success: false,
		Data: data,
	}

	json.NewEncoder(w).Encode(resp)
}