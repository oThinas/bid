package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
		// Provide more specific error messages for common JSON parsing issues
		var problems validator.Evaluator

		if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
			problems = validator.Evaluator{
				jsonErr.Field: fmt.Sprintf("expected %s, got %s", jsonErr.Type, jsonErr.Value),
			}
		} else if strings.Contains(err.Error(), "cannot unmarshal") {
			// Try to extract field name from error message
			problems = validator.Evaluator{
				"json": "invalid JSON format - check field types (numbers should not be quoted)",
			}
		} else {
			problems = validator.Evaluator{
				"json": "invalid JSON format",
			}
		}

		return data, problems, fmt.Errorf("failed to decode JSON request body: %w", err)
	}

	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}
