package errors

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	serverrors "dddstructure/service/errors"
)

var (
	// ErrNotFound is returned when a given endpoint could not be found.
	ErrNotFound = New(http.StatusNotFound, "", "The given endpoint could not be found")

	// ErrBadRequest is returned usually when sent parameters are invalid.
	ErrBadRequest = New(http.StatusBadRequest, "", "Bad request most likely due to invalid parameters")

	// ErrInternalServerError is returned when an internal server error occurs.
	ErrInternalServerError = New(http.StatusInternalServerError, "", "Internal server error")
)

var (
	// ErrOffsetInvalid is returned when the offset parameter is invalid.
	ErrOffsetInvalid = New(http.StatusBadRequest, "offset", "Offset parameter is invalid, must be an integer")

	// ErrLimitInvalid is returned when the limit parameter is invalid.
	ErrLimitInvalid = New(http.StatusBadRequest, "limit", "Limit parameter is invalid, must be an integer")

	// ErrLimitMax is returned when the limit parameter is greater than the
	// maximum allowable limit.
	ErrLimitMax = func(limit, max uint) *Error {
		return New(http.StatusBadRequest,
			"limit",
			fmt.Sprintf("Limit value of %d is greater than maximum allowable limit of %d", limit, max),
		)
	}
)

// Error defines the default API error type.
type Error struct {
	Status int    `json:"status,omitempty"`
	Param  string `json:"param,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// New returns a new Error.
func New(status int, param, detail string) *Error {
	return &Error{status, param, detail}
}

// Errors defines multiple API errors.
type Errors []*Error

// Add appends an Error.
func (es *Errors) Add(e *Error) {
	*es = append(*es, e)
}

// Length returns the number of errors.
func (es *Errors) Length() int {
	return len(*es)
}

// ErrorsWrap wraps the Errors type to produce JSON according to spec.
type ErrorsWrap struct {
	Errors *Errors `json:"errors"`
}

// Default renders an error using a default API error.
//
// The default error response will look like:
//
//	{
//	  errors: [
//	    {
//	      status: 404,
//	      detail: "The given endpoint could not be found"
//	    }
//	  ]
//	}
//
// If the Param field is not blank, this will also be included in the error
// response, used to signify specific parameter errors.
func Default(logger *slog.Logger, w http.ResponseWriter, e *Error) {
	// Wrap the error in top level errors array.
	wrap := ErrorsWrap{&Errors{e}}

	// Write the HTTP response code.
	w.WriteHeader(e.Status)

	// Create a new Encoder.
	enc := json.NewEncoder(w)

	// Handle HTML escape rule and indentation settings.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Set headers.
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := enc.Encode(wrap); err != nil {
		logger.Error("errors.Default() render.JSON() error",
			slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Multiple renders multiple API errors.
func Multiple(logger *slog.Logger, w http.ResponseWriter, status int, es *Errors) {
	// Wrap the error in top level errors array.
	wrap := ErrorsWrap{es}

	// Write the HTTP response code.
	w.WriteHeader(status)

	// Create a new Encoder.
	enc := json.NewEncoder(w)

	// Handle HTML escape rule and indentation settings.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Set headers.
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := enc.Encode(wrap); err != nil {
		logger.Error("errors.Default() render.JSON() error",
			slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Params renders API parameter errors.
//
// This is a helper function used by the API endpoint handlers to make it
// easier to render parameter errors returned from services.
func Params(logger *slog.Logger, w http.ResponseWriter, status int, pes *serverrors.ParamErrors) {
	// Create new Errors.
	errs := &Errors{}

	// Loop through each parameter error.
	for _, pe := range *pes {
		errs.Add(New(http.StatusBadRequest, pe.Name, pe.Error()))
	}

	Multiple(logger, w, http.StatusBadRequest, errs)
}
