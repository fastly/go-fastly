package fastly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
)

// FieldError represents a custom error type for API data fields.
type FieldError struct {
	kind    string
	message string
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
// field was either missing or some other type of error occurred.
func (e *FieldError) Error() string {
	if e.message != "" {
		return fmt.Sprintf("problem with field '%s': %s", e.kind, e.message)
	}

	return fmt.Sprintf("missing required field '%s'", e.kind)
}

// Message prints the error message.
func (e *FieldError) Message(msg string) *FieldError {
	e.message = msg
	return e
}

// NewFieldError returns an error that formats as the given text.
func NewFieldError(kind string) *FieldError {
	return &FieldError{
		kind: kind,
	}
}

const emptyTokenInvalid string = "the token value cannot be empty"

// ErrTokenEmpty is an error that is returned when an input struct
// specifies an "Token" key value which the user has set to an empty string.
var ErrTokenEmpty = NewFieldError("Token").Message(emptyTokenInvalid)

const batchModifyMaxExceeded string = "batch modify maximum operations exceeded"

// ErrMaxExceededEntries is an error that is returned when an input struct
// specifies an "Entries" key value exceeding the maximum allowed.
var ErrMaxExceededEntries = NewFieldError("Entries").Message(batchModifyMaxExceeded)

// ErrMaxExceededItems is an error that is returned when an input struct
// specifies an "Items" key value exceeding the maximum allowed.
var ErrMaxExceededItems = NewFieldError("Items").Message(batchModifyMaxExceeded)

// ErrMaxExceededRules is an error that is returned when an input struct
// specifies an "Rules" key value exceeding the maximum allowed.
var ErrMaxExceededRules = NewFieldError("Rules").Message(batchModifyMaxExceeded)

// ErrMissingACLID is an error that is returned when an input struct
// requires a "ACLID" key, but one was not set.
var ErrMissingACLID = NewFieldError("ACLID")

// ErrMissingAddress is an error that is returned when an input struct
// requires a "Address" key, but one was not set.
var ErrMissingAddress = NewFieldError("Address")

// ErrMissingBackend is an error that is returned when an input struct
// requires a "Backend" key, but one was not set.
var ErrMissingBackend = NewFieldError("Backend")

// ErrMissingCertBlob is an error that is returned when an input struct
// requires a "CertBlob" key, but one was not set.
var ErrMissingCertBlob = NewFieldError("CertBlob")

// ErrMissingContent is an error that is returned when an input struct
// requires a "Content" key, but one was not set.
var ErrMissingContent = NewFieldError("Content")

// ErrMissingType is an error that is returned when an input struct
// requires a "Type" key, but one was not set.
var ErrMissingType = NewFieldError("Type")

// ErrMissingCustomerID is an error that is returned when an input struct
// requires a "CustomerID" key, but one was not set.
var ErrMissingCustomerID = NewFieldError("CustomerID")

// ErrMissingDictionaryID is an error that is returned when an input struct
// requires a "DictionaryID" key, but one was not set.
var ErrMissingDictionaryID = NewFieldError("DictionaryID")

// ErrMissingDirector is an error that is returned when an input struct
// requires a "Director" key, but one was not set.
var ErrMissingDirector = NewFieldError("Director")

// ErrMissingEventID is an error that is returned when an input struct
// requires a "EventID" key, but one was not set.
var ErrMissingEventID = NewFieldError("EventID")

// ErrMissingFrom is an error that is returned when an input struct
// requires a "From" key, but one was not set.
var ErrMissingFrom = NewFieldError("From")

// ErrMissingTokenID is an error that is returned when an input struct requires a
// "TokenID" key, but one was not set.
var ErrMissingTokenID = errors.New("missing required field 'TokenID'")

// ErrMissingID is an error that is returned when an input struct
// requires a "ID" key, but one was not set.
var ErrMissingID = NewFieldError("ID")

// ErrMissingERLID is an error that is returned when an input struct
// requires an "ERLID" key, but one was not set.
var ErrMissingERLID = NewFieldError("ERLID")

// ErrMissingIP is an error that is returned when an input struct
// requires a "IP" key, but one was not set.
var ErrMissingIP = NewFieldError("IP")

// ErrMissingIntermediatesBlob is an error that is returned when an input struct
// requires a "IntermediatesBlob" key, but one was not set.
var ErrMissingIntermediatesBlob = NewFieldError("IntermediatesBlob")

// ErrMissingItemKey is an error that is returned when an input struct
// requires a "ItemKey" key, but one was not set.
var ErrMissingItemKey = NewFieldError("ItemKey")

// ErrMissingKey is an error that is returned when an input struct
// requires a "Key" key, but one was not set.
var ErrMissingKey = NewFieldError("Key")

// ErrMissingKeys is an error that is returned when an input struct
// requires a "Keys" key, but one was not set.
var ErrMissingKeys = NewFieldError("Keys")

// ErrMissingLogin is an error that is returned when an input struct
// requires a "Login" key, but one was not set.
var ErrMissingLogin = NewFieldError("Login")

// ErrMissingMonth is an error that is returned when an input struct
// requires a "Month" key, but one was not set.
var ErrMissingMonth = NewFieldError("Month")

// ErrMissingName is an error that is returned when an input struct
// requires a "Name" key, but one was not set.
var ErrMissingName = NewFieldError("Name")

// ErrMissingNameValue is an error that is returned when an input struct
// requires a "Name" key, but one was not set.
var ErrMissingNameValue = NewFieldError("Name").Message("service name can't be an empty value")

// ErrMissingNewName is an error that is returned when an input struct
// requires a "NewName" key, but one was not set.
var ErrMissingNewName = NewFieldError("NewName")

// ErrMissingNumber is an error that is returned when an input struct
// requires a "Number" key, but one was not set.
var ErrMissingNumber = NewFieldError("Number")

// ErrMissingPoolID is an error that is returned when an input struct
// requires a "PoolID" key, but one was not set.
var ErrMissingPoolID = NewFieldError("PoolID")

// ErrMissingSecret is an error that is returned when an input struct
// requires a "Secret" key, but one was not set.
var ErrMissingSecret = NewFieldError("Secret")

// ErrMissingServer is an error that is returned when an input struct
// requires a "Server" key, but one was not set.
var ErrMissingServer = NewFieldError("Server")

// ErrMissingServerSideEncryptionKMSKeyID is an error that is returned when an
// input struct requires a "ServerSideEncryptionKMSKeyID" key, but one was not set.
var ErrMissingServerSideEncryptionKMSKeyID = NewFieldError("ServerSideEncryptionKMSKeyID")

// ErrMissingServiceID is an error that is returned when an input struct
// requires a "ServiceID" key, but one was not set.
var ErrMissingServiceID = NewFieldError("ServiceID")

// ErrMissingServiceAuthorizationsService is an error that is returned when an input struct
// requires a "Service" key of type SAService, but one was not set or was misconfigured.
var ErrMissingServiceAuthorizationsService = NewFieldError("Service").Message("SAService requires an ID")

// ErrMissingServiceAuthorizationsUser is an error that is returned when an input struct
// requires a "User" key of type SAUser, but one was not set or was misconfigured.
var ErrMissingServiceAuthorizationsUser = NewFieldError("User").Message("SAUser requires an ID")

// ErrMissingStoreID is an error that is returned when an input struct
// requires a "StoreID" key, but one was not set.
var ErrMissingStoreID = NewFieldError("StoreID")

// ErrMissingUserID is an error that is returned when an input struct
// requires a "UserID" key, but one was not set
var ErrMissingUserID = NewFieldError("UserID")

// ErrMissingPermission is an error that is returned when an input struct
// requires a "Permission" key, but one was not set
var ErrMissingPermission = NewFieldError("Permission")

// ErrMissingServiceVersion is an error that is returned when an input struct
// requires a "ServiceVersion" key, but one was not set.
var ErrMissingServiceVersion = NewFieldError("ServiceVersion")

// ErrMissingTLSCertificate is an error that is returned when an input struct
// requires a "TLSCertificate" key, but one was not set.
var ErrMissingTLSCertificate = NewFieldError("TLSCertificate")

// ErrMissingTLSConfiguration is an error that is returned when an input
// struct requires a "TLSConfiguration" key, but one was not set.
var ErrMissingTLSConfiguration = NewFieldError("TLSConfiguration")

// ErrMissingTLSDomain is an error that is returned when an input struct
// requires a "TLSDomain" key, but one was not set.
var ErrMissingTLSDomain = NewFieldError("TLSDomain")

// ErrCommonNameNotInDomains is an error that is returned when an input struct
// requires that the domain in "CommonName" is also in "Domains"
var ErrCommonNameNotInDomains = NewFieldError("CommonName").Message("CommonName must be in Domains")

// ErrMissingTo is an error that is returned when an input struct
// requires a "To" key, but one was not set.
var ErrMissingTo = NewFieldError("To")

// ErrMissingKind is an error that is returned when an input struct requires a
// "Kind" key, but one was not set.
var ErrMissingKind = NewFieldError("Kind")

// ErrMissingURL is an error that is returned when an input struct
// requires a "URL" key, but one was not set.
var ErrMissingURL = NewFieldError("URL")

// ErrMissingValue is an error that is returned when an input struct
// requires a "Value" key, but one was not set.
var ErrMissingValue = NewFieldError("Value")

// ErrMissingWAFActiveRule is an error that is returned when an input struct
// requires a "Rules" key, but there needs to be at least one WAFActiveRule entry.
var ErrMissingWAFActiveRule = NewFieldError("Rules").Message("expect at least one WAFActiveRule")

// ErrMissingWAFID is an error that is returned when an input struct
// requires a "WAFID" key, but one was not set.
var ErrMissingWAFID = NewFieldError("WAFID")

// ErrMissingWAFRuleExclusion is an error that is returned when an input struct
// requires a "WAFRuleExclusion" key, but one was not set.
var ErrMissingWAFRuleExclusion = NewFieldError("WAFRuleExclusion")

// ErrMissingWAFRuleExclusionNumber is an error that is returned when an input
// struct requires a "WAFRuleExclusionNumber" key, but one was not set.
var ErrMissingWAFRuleExclusionNumber = NewFieldError("WAFRuleExclusionNumber")

// ErrMissingWAFVersionID is an error that is returned when an input struct
// requires a "WAFVersionID" key, but one was not set.
var ErrMissingWAFVersionID = NewFieldError("WAFVersionID")

// ErrMissingWAFVersionNumber is an error that is returned when an input
// struct requires a "WAFVersionNumber" key, but one was not set.
var ErrMissingWAFVersionNumber = NewFieldError("WAFVersionNumber")

// ErrMissingYear is an error that is returned when an input struct requires a
// "Year" key, but one was not set.
var ErrMissingYear = NewFieldError("Year")

// ErrMissingOptionalNameComment is an error that is returned when an input
// struct requires either a "Name" or "Comment" key, but one was not set.
var ErrMissingOptionalNameComment = NewFieldError("Name, Comment").Message("at least one of the available 'optional' fields is required")

// ErrMissingTokensValue is an error that is returned when an input struct
// requires a "Tokens" key, but there needs to be at least one token entry.
var ErrMissingTokensValue = NewFieldError("Tokens").Message("expect at least one token")

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

// ErrMissingToken is an error that is returned when an input struct
// requires a "Token" key, but one was not set.
var ErrMissingToken = NewFieldError("Token")

// ErrMissingProductID is an error that is returned when an input struct
// requires a "ProductID" key, but one was not set.
var ErrMissingProductID = NewFieldError("ProductID")

// Ensure HTTPError is, in fact, an error.
var _ error = (*HTTPError)(nil)

// HTTPError is a custom error type that wraps an HTTP status code with some
// helper functions.
type HTTPError struct {
	Errors []*ErrorObject `mapstructure:"errors"`
	// StatusCode is the HTTP status code (2xx-5xx).
	StatusCode int
}

// ErrorObject is a single error.
type ErrorObject struct {
	Code   string                  `mapstructure:"code"`
	Detail string                  `mapstructure:"detail"`
	ID     string                  `mapstructure:"id"`
	Meta   *map[string]interface{} `mapstructure:"meta"`
	Status string                  `mapstructure:"status"`
	Title  string                  `mapstructure:"title"`
}

// legacyError represents the older-style errors from Fastly. It is private
// because it is automatically converted to a jsonapi error.
type legacyError struct {
	Detail  string `mapstructure:"detail"`
	Message string `mapstructure:"msg"`
}

// NewHTTPError creates a new HTTP error from the given code.
func NewHTTPError(resp *http.Response) *HTTPError {
	var e HTTPError
	e.StatusCode = resp.StatusCode

	if resp.Body == nil {
		return &e
	}

	// Save a copy of the body as it's read/decoded.
	// If decoding fails, it can then be used (via addDecodeErr)
	// to create a generic error containing the body's read contents.
	var bodyCp bytes.Buffer
	body := io.TeeReader(resp.Body, &bodyCp)
	addDecodeErr := func() {
		// There are 2 errors at this point:
		//  1. The response error.
		//  2. The error decoding the response.
		// The response error is still most relevant to users (just unable to be decoded).
		// Provide the response's body verbatim as the error 'Detail' with the assumption
		// that it may contain useful information, e.g. 'Bad Gateway'.
		// The decode error could be conflated with the response error, so it is omitted.
		e.Errors = append(e.Errors, &ErrorObject{
			Title:  "Undefined error",
			Detail: bodyCp.String(),
		})
	}

	switch resp.Header.Get("Content-Type") {
	case jsonapi.MediaType:
		// If this is a jsonapi response, decode it accordingly.
		if err := decodeBodyMap(body, &e); err != nil {
			addDecodeErr()
		}

	case "application/problem+json":
		// Response is a "problem detail" as defined in RFC 7807.
		var problemDetail struct {
			Detail string `json:"detail,omitempty"` // A human-readable description of the specific error, aiming to help the user correct the problem
			Status int    `json:"status"`           // HTTP status code
			Title  string `json:"title,omitempty"`  // A short name for the error type, which remains constant from occurrence to occurrence
			URL    string `json:"type,omitempty"`   // URL to a human-readable document describing this specific error condition
		}
		if err := json.NewDecoder(body).Decode(&problemDetail); err != nil {
			addDecodeErr()
		} else {
			e.Errors = append(e.Errors, &ErrorObject{
				Title:  problemDetail.Title,
				Detail: problemDetail.Detail,
				Status: strconv.Itoa(problemDetail.Status),
			})
		}

	default:
		var lerr *legacyError
		if err := decodeBodyMap(body, &lerr); err != nil {
			addDecodeErr()
		} else if lerr != nil {
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
