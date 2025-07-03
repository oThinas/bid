package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oThinas/bid/internal/validator"
)

// EncodeJSON writes the given data as a JSON response with the specified status code.
// It sets the Content-Type header to application/json and encodes the data to the response writer.
// Returns an error if encoding fails.
func EncodeJSON[T any](w http.ResponseWriter, r *http.Request, statusCode int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}

	return nil
}

// DecodeJSON decodes the JSON request body into a struct that implements validator.Validator.
// It returns the decoded data, any validation errors, and an error if decoding or validation fails.
func DecodeJSON[T validator.Validator](r *http.Request) (T, validator.Evaluator, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, nil, fmt.Errorf("failed to decode JSON request body: %w", err)
	}

	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}
