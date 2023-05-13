package productapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ValidateJSON returns an error if the provided JSON is not valid
func ValidateJSON(v any) error {
	validate := validator.New()
	return validate.Struct(v)
}

// ReadJSON returns an error if decoding JSON fails
func ReadJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// WriteJSON returns an error if encoding JSON fails
func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

// WriteERR returns a map containing an error | {"error": "example error message"}
func WriteERR(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	err := map[string]string{"error": msg}
	_ = json.NewEncoder(w).Encode(err)
}
