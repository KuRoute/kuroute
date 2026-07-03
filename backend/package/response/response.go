package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func OK[T any](w http.ResponseWriter, status int, data T) {
	writeJSON(w, status, SuccessResponse[T]{
		Success: true,
		Data:    data,
	})
}

func Fail(w http.ResponseWriter, status int, code, message string, details ...any) {
	var d any
	if len(details) > 0 {
		d = details[0]
	}

	writeJSON(w, status, ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: d,
		},
	})
}