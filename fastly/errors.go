package fastly

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
)

// FieldError represents a custom error type for API data fields.
type FieldError struct {
	kind   string
	custom string
}

// Error fulfills the error interface.
//
// NOTE: some fields are optional but still need to present an error depending
// on the API they are associated with. For example, when updating a service
// the 'name' and 'comment' fields are both optional, but at least one of them
// needs to be provided for the API call to have any purpose (otherwise the API
// backend will just reject the call, thus being a waste of network resources).
//
// Because of this we allow modifying the error message to reflect whether the
// missing field was either 'required' or just missing a value.
func (e *FieldError) Error() string {
	if e.custom != "" {
		return fmt.Sprintf("problem with field '%s': %s", e.kind, e.custom)
	}

	return fmt.Sprintf("missing required field '%s'", e.kind)
}

func (e *FieldError) Custom(msg string) *FieldError {
	e.custom = msg
	return e
}

// NewFieldError returns an error that formats as the given text.
func NewFieldError(kind string) *FieldError {
	return &FieldError{
		kind: kind,
	}
}

// ErrStatusNotOk is an error that indicates the response body returned by the
// Fastly API was not `{"status": "ok"}`
var ErrStatusNotOk = errors.New("unexpected 'status' field in API response body")

// ErrNotOK is a generic error indicating that something is not okay.
var ErrNotOK = errors.New("not ok")

// ErrNotImplemented is a generic error indicating that something is not yet implemented.
var ErrNotImplemented = errors.New("not implemented")

// ErrManagedLoggingEnabled is an error that indicates that managed logging was
// already enabled for a service.
var ErrManagedLoggingEnabled = errors.New("managed logging already enabled")

// Ensure HTTPError is, in fact, an error.
var _ error = (*HTTPError)(nil)

// HTTPError is a custom error type that wraps an HTTP status code with some
// helper functions.
type HTTPError struct {
	// StatusCode is the HTTP status code (2xx-5xx).
	StatusCode int

	Errors []*ErrorObject `mapstructure:"errors"`
}

// ErrorObject is a single error.
type ErrorObject struct {
	ID     string `mapstructure:"id"`
	Title  string `mapstructure:"title"`
	Detail string `mapstructure:"detail"`
	Status string `mapstructure:"status"`
	Code   string `mapstructure:"code"`

	Meta *map[string]interface{} `mapstructure:"meta"`
}

// legacyError represents the older-style errors from Fastly. It is private
// because it is automatically converted to a jsonapi error.
type legacyError struct {
	Message string `mapstructure:"msg"`
	Detail  string `mapstructure:"detail"`
}

// NewHTTPError creates a new HTTP error from the given code.
func NewHTTPError(resp *http.Response) *HTTPError {
	var e HTTPError
	e.StatusCode = resp.StatusCode

	if resp.Body == nil {
		return &e
	}

	// If this is a jsonapi response, decode it accordingly
	if resp.Header.Get("Content-Type") == jsonapi.MediaType {
		if err := decodeBodyMap(resp.Body, &e); err != nil {
			panic(err)
		}
	} else {
		var lerr *legacyError
		decodeBodyMap(resp.Body, &lerr)
		if lerr != nil {
			e.Errors = append(e.Errors, &ErrorObject{
				Title:  lerr.Message,
				Detail: lerr.Detail,
			})
		}
	}

	return &e
}

// Error implements the error interface and returns the string representing the
// error text that includes the status code and the corresponding status text.
func (e *HTTPError) Error() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "%d - %s:", e.StatusCode, http.StatusText(e.StatusCode))

	for _, e := range e.Errors {
		fmt.Fprintf(&b, "\n")

		if e.ID != "" {
			fmt.Fprintf(&b, "\n    ID:     %s", e.ID)
		}

		if e.Title != "" {
			fmt.Fprintf(&b, "\n    Title:  %s", e.Title)
		}

		if e.Detail != "" {
			fmt.Fprintf(&b, "\n    Detail: %s", e.Detail)
		}

		if e.Code != "" {
			fmt.Fprintf(&b, "\n    Code:   %s", e.Code)
		}

		if e.Meta != nil {
			fmt.Fprintf(&b, "\n    Meta:   %v", *e.Meta)
		}
	}

	return b.String()
}

// String implements the stringer interface and returns the string representing
// the string text that includes the status code and corresponding status text.
func (e *HTTPError) String() string {
	return e.Error()
}

// IsNotFound returns true if the HTTP error code is a 404, false otherwise.
func (e *HTTPError) IsNotFound() bool {
	return e.StatusCode == 404
}
