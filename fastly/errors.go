package fastly

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
)

// ErrMissingService is an error that is returned when an input struct requires
// a "Service" key, but one was not set.
var ErrMissingService = errors.New("missing required field 'Service'")

// ErrMissingStatus is an error that is returned when an input struct requires
// a "Status" key, but one was not set.
var ErrMissingStatus = errors.New("missing required field 'Status'")

// ErrMissingTag is an error that is returned when an input struct requires
// a "Tag" key, but one was not set.
var ErrMissingTag = errors.New("missing required field 'Tag'")

// ErrMissingVersion is an error that is returned when an input struct requires
// a "Version" key, but one was not set.
var ErrMissingVersion = errors.New("missing required field 'Version'")

// ErrMissingContent is an error that is returned when an input struct requires a
// "Content" key, but one was not set.
var ErrMissingContent = errors.New("missing required field 'Content'")

// ErrMissingLogin is an error that is returned when an input struct requires a
// "Login" key, but one was not set.
var ErrMissingLogin = errors.New("missing required field 'Login'")

// ErrMissingName is an error that is returned when an input struct requires a
// "Name" key, but one was not set.
var ErrMissingName = errors.New("missing required field 'Name'")

// ErrMissingKey is an error that is returned when an input struct requires a
// "Key" key, but one was not set.
var ErrMissingKey = errors.New("missing required field 'Key'")

// ErrMissingURL is an error that is returned when an input struct requires a
// "URL" key, but one was not set.
var ErrMissingURL = errors.New("missing required field 'URL'")

// ErrMissingID is an error that is returned when an input struct requires an
// "ID" key, but one was not set.
var ErrMissingID = errors.New("missing required field 'ID'")

// ErrMissingDictionary is an error that is returned when an input struct
// requires a "Dictionary" key, but one was not set.
var ErrMissingDictionary = errors.New("missing required field 'Dictionary'")

// ErrMissingItemKey is an error that is returned when an input struct
// requires a "ItemKey" key, but one was not set.
var ErrMissingItemKey = errors.New("missing required field 'ItemKey'")

// ErrMissingFrom is an error that is returned when an input struct
// requires a "From" key, but one was not set.
var ErrMissingFrom = errors.New("missing required field 'From'")

// ErrMissingTo is an error that is returned when an input struct
// requires a "To" key, but one was not set.
var ErrMissingTo = errors.New("missing required field 'To'")

// ErrMissingDirector is an error that is returned when an input struct
// requires a "From" key, but one was not set.
var ErrMissingDirector = errors.New("missing required field 'Director'")

// ErrMissingBackend is an error that is returned when an input struct
// requires a "Backend" key, but one was not set.
var ErrMissingBackend = errors.New("missing required field 'Backend'")

// ErrMissingYear is an error that is returned when an input struct
// requires a "Year" key, but one was not set.
var ErrMissingYear = errors.New("missing required field 'Year'")

// ErrMissingMonth is an error that is returned when an input struct
// requires a "Month" key, but one was not set.
var ErrMissingMonth = errors.New("missing required field 'Month'")

// ErrMissingNewName is an erorr that is returned when an input struct
// requires a "NewName" key, but one was not set
var ErrMissingNewName = errors.New("missing required field 'NewName'")

// ErrMissingAcl is an error that is returned when an input struct
// required an "Acl" key, but one is not set
var ErrMissingACL = errors.New("missing required field 'ACL'")

// ErrMissingIP is an error that is returned when an input struct
// required an "IP" key, but one is not set
var ErrMissingIP = errors.New("missing required field 'IP'")

// ErrMissingCustomerID is an error that is returned was an input struct
// requires a "CustomerID" key, but one was not set
var ErrMissingCustomerID = errors.New("missing required field 'CustomerID'")

// ErrMissingEventID is an error that is returned was an input struct
// requires a "EventID" key, but one was not set
var ErrMissingEventID = errors.New("missing required field 'EventID'")

// ErrMissingWAFID is an error that is returned when an input struct
// requires a "WAFID" key, but one was not set.
var ErrMissingWAFID = errors.New("missing required field 'WAFID'")

// ErrMissingWAFVersionNumber is an error that is returned when an input struct
// requires a "WAFVersionNumber" key, but one was not set.
var ErrMissingWAFVersionNumber = errors.New("missing required field 'WAFVersionNumber'")

// ErrMissingWAFVersionID is an error that is returned when an input struct
// requires a "WAFVersionID" key, but one was not set.
var ErrMissingWAFVersionID = errors.New("missing required field 'WAFVersionID'")

// ErrMissingWAFActiveRuleList is an error that is returned when an input struct
// requires a list of WAF active rules, but it is empty.
var ErrMissingWAFActiveRuleList = errors.New("WAF active rules slice is empty")

// ErrMissingWAFExclusionNumber is an error that is returned when an input struct
// requires a "WAFExclusionNumber" key, but one was not set.
var ErrMissingWAFExclusionNumber = errors.New("missing required field 'WAFExclusionNumber'")

// ErrMissingWAFExclusion is an error that is returned when an input struct
// requires a "WAFExclusion" key, but one was not set.
var ErrMissingWAFExclusion = errors.New("missing required field 'WAFExclusion'")

// ErrMissingOWASPID is an error that is returned was an input struct
// requires a "OWASPID" key, but one was not set
var ErrMissingOWASPID = errors.New("missing required field 'OWASPID'")

// ErrMissingRuleID is an error that is returned was an input struct
// requires a "RuleID" key, but one was not set
var ErrMissingRuleID = errors.New("missing required field 'RuleID'")

// ErrMissingConfigSetID is an error that is returned was an input struct
// requires a "ConfigSetID" key, but one was not set
var ErrMissingConfigSetID = errors.New("missing required field 'ConfigSetID'")

// ErrMissingWAFList is an error that is returned was an input struct
// requires a list of WAF id's, but it is empty
var ErrMissingWAFList = errors.New("WAF slice is empty")

// ErrMissingPool is an error that is returned when an input struct requires
// a "Pool" key, but one was not set.
var ErrMissingPool = errors.New("missing required field 'Pool'")

// ErrMissingServer is an error that is returned when an input struct requires
// a "Server" key, but one was not set.
var ErrMissingServer = errors.New("missing required field 'Server'")

// ErrMissingAddress is an error that is returned when an input struct requires
// a "Address" key, but one was not set.
var ErrMissingAddress = errors.New("missing required field 'Address'")

// ErrBatchUpdateMaximumItemsExceeded is an error that indicates that too many batch operations are being executed.
// The Fastly API specifies an maximum limit.
var ErrBatchUpdateMaximumOperationsExceeded = errors.New("batch modify maximum operations exceeded")

// ErrMissingKMSKeyID is an error that is returned from an input struct that requires
// a "ServerSideEncryptionKMSKeyID" key, but one was not set.
var ErrMissingKMSKeyID = errors.New("missing required field 'ServerSideEncryptionKMSKeyID'")

// ErrMissingCertBlob is an error that is returned from an input struct that requires
// a "CertBlob" key, but one was not set.
var ErrMissingCertBlob = errors.New("missing required field 'CertBlob'")

// ErrMissingIntermediatesBlob is an error that is returned from an input struct that requires
// a "IntermediatesBlob" key, but one was not set.
var ErrMissingIntermediatesBlob = errors.New("missing required field 'IntermediatesBlob'")

// ErrStatusNotOk is an error that indicates that indicates that the response body returned
// by the Fastly API was not `{"status": "ok"}`
var ErrStatusNotOk = errors.New("unexpected 'status' field in API response body")

// ErrNotOK is a generic error indicating that something is not okay.
var ErrNotOK = errors.New("not ok")

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
