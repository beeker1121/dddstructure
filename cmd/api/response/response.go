package response

import (
	"encoding/json"
	"net/http"
)

// JSON responds with the given data encoded to JSON.
func JSON(w http.ResponseWriter, indent bool, v any) error {
	// Create a new Encoder.
	enc := json.NewEncoder(w)

	// Handle HTML escape rule and indentation settings.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Set headers.
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return enc.Encode(v)
}
