package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data any) error {
	return JSONWithHeaders(w, status, data, nil)
}

func JSONWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// ErrorResponse represents the structure of an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// BadRequest sends a 400 Bad Request response with the given error
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	JSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
}

// ServerError sends a 500 Internal Server Error response
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	// In a production environment, you might want to log the error here
	JSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
}

// InvalidCredentials sends a 401 Unauthorized response for invalid credentials
func InvalidCredentials(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
}