package util

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the error response sent to clients.
type ErrorResponse struct {
	Ok      bool   `json:"ok"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HttpError writes a JSON error response to the client.
//
// The error message is formatted as:
// {
//   "ok": false,
//   "error": "http error code as string",
//   "message": "error message"
// }
func HttpError(w http.ResponseWriter, status int, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	res := ErrorResponse{
		Ok:      false,
		Error:   http.StatusText(status),
		Message: message,
	}
	json.NewEncoder(w).Encode(res)
}

// HttpJson writes a JSON response to the client.
func HttpJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}
