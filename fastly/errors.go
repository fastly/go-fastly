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

const alertTypeDoesNotMatch string = "alert type does not match"

// ErrInvalidType is an error that is returned when an alert is being updated,
// but the alert is not of the correct type.
var ErrInvalidType = NewFieldError("Type").Message(alertTypeDoesNotMatch)

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

// ErrMissingAlertID is an error that is returned when an input struct
// requires a "AlertID" key, but one was not set.
var ErrMissingAlertID = NewFieldError("AlertID")

// ErrMissingBackend is an error that is returned when an input struct
// requires a "Backend" key, but one was not set.
var ErrMissingBackend = NewFieldError("Backend")

// ErrMissingCertBlob is an error that is returned when an input struct
// requires a "CertBlob" key, but one was not set.
var ErrMissingCertBlob = NewFieldError("CertBlob")

// ErrMissingCertBundle is an error that is returned when an input struct
// requires a "CertBundle" key, but one was not set.
var ErrMissingCertBundle = NewFieldError("CertBundle")

// ErrMissingConditions is an error that is returned when an input struct
// requires a "Conditions" key, but one was not set.
var ErrMissingConditions = NewFieldError("Conditions")

// ErrMissingComputeACLID is an error that is returned when an input struct
// requires a "ComputeACLID" key, but one was not set.
var ErrMissingComputeACLID = NewFieldError("ComputeACLID")

// ErrMissingComputeACLIP is an error that is returned when an input struct
// requires a "ComputeACLIP" key, but one was not set.
var ErrMissingComputeACLIP = NewFieldError("ComputeACLIP")

// ErrMissingConfig is an error that is returned when an input struct
// requires a "Config" key, but one was not set.
var ErrMissingConfig = NewFieldError("Config")

// ErrMissingWorkspaceID is an error that is returned when an input struct
// requires a "WorkspaceID" key, but one was not set.
var ErrMissingWorkspaceID = NewFieldError("WorkspaceID")

// ErrMissingListID is an error that is returned when an input struct
// requires a "ListID" key, but one was not set.
var ErrMissingListID = NewFieldError("ListID")

// ErrMissingEntries is an error that is returned when an input struct
// requires a "Entries" key, but one was not set.
var ErrMissingEntries = NewFieldError("Entries")

// ErrMissingIsExpired is an error that is returned when an input struct
// requires an "IsExpired" key, but one was not set.
var ErrMissingIsExpired = NewFieldError("IsExpired")

// ErrMissingRedactionID is an error that is returned when an input struct
// requires a "RedactionID" key, but one was not set.
var ErrMissingRedactionID = NewFieldError("RedactionID")

// ErrMissingThresholdID is an error that is returned when an input struct
// requires a "ThresholdID" key, but one was not set.
var ErrMissingThresholdID = NewFieldError("ThresholdID")

// ErrMissingAction is an error that is returned when an input struct
// requires a "Action" key, but one was not set.
var ErrMissingAction = NewFieldError("Action")

// ErrMissingLimit is an error that is returned when an input struct
// requires a "Limit" key, but one was not set.
var ErrMissingLimit = NewFieldError("Limit")

// ErrMissingInterval is an error that is returned when an input struct
// requires a "Interval" key, but one was not set.
var ErrMissingInterval = NewFieldError("Interval")

// ErrMissingSignal is an error that is returned when an input struct
// requires a "Signal" key, but one was not set.
var ErrMissingSignal = NewFieldError("Signal")

// ErrMissingRuleID is an error that is returned when an input struct
// requires a "RuleID" key, but one was not set.
var ErrMissingRuleID = NewFieldError("RuleID")

// ErrMissingSignalID is an error that is returned when an input struct
// requires a "SignalID" key, but one was not set.
var ErrMissingSignalID = NewFieldError("SignalID")

// ErrMissingRequestID is an error that is returned when an input struct
// requires a "RequestID" key, but one was not set.
var ErrMissingRequestID = NewFieldError("RequestID")

// ErrMissingField is an error that is returned when an input struct
// requires a "Field" key, but one was not set.
var ErrMissingField = NewFieldError("Field")

// ErrMissingContent is an error that is returned when an input struct
// requires a "Content" key, but one was not set.
var ErrMissingContent = NewFieldError("Content")

// ErrMissingType is an error that is returned when an input struct
// requires a "Type" key, but one was not set.
var ErrMissingType = NewFieldError("Type")

// ErrMissingCustomerID is an error that is returned when an input struct
// requires a "CustomerID" key, but one was not set.
var ErrMissingCustomerID = NewFieldError("CustomerID")

// ErrMissingAccessKeyID is an error that is returned when an input struct
// requires a "AccessKeyID" key, but one was not set.
var ErrMissingAccessKeyID = NewFieldError("AccessKeyID")

// ErrMissingDescription is an error that is returned when an input struct
// requires a "Description" key, but one was not set.
var ErrMissingDescription = NewFieldError("Description")

// ErrMissingScope is an error that is returned when an input struct
// requires a "Scope" key, but one was not set.
var ErrMissingScope = NewFieldError("Scope")

// ErrMissingDictionaryID is an error that is returned when an input struct
// requires a "DictionaryID" key, but one was not set.
var ErrMissingDictionaryID = NewFieldError("DictionaryID")

// ErrMissingDirector is an error that is returned when an input struct
// requires a "Director" key, but one was not set.
var ErrMissingDirector = NewFieldError("Director")

// ErrMissingEventID is an error that is returned when an input struct
// requires a "EventID" key, but one was not set.
var ErrMissingEventID = NewFieldError("EventID")

// ErrMissingEvents is an error that is returned when an input struct
// requires a "Events" key, but one was not set.
var ErrMissingEvents = NewFieldError("Events")

// ErrMissingFrom is an error that is returned when an input struct
// requires a "From" key, but one was not set.
var ErrMissingFrom = NewFieldError("From")

// ErrMissingTokenID is an error that is returned when an input struct requires a
// "TokenID" key, but one was not set.
var ErrMissingTokenID = errors.New("missing required field 'TokenID'")

// ErrMissingID is an error that is returned when an input struct
// requires a "ID" key, but one was not set.
var ErrMissingID = NewFieldError("ID")

// ErrMissingDomainID is an error that is returned when an input struct
// requires a "DomainID" key, but one was not set.
var ErrMissingDomainID = NewFieldError("DomainID")

// ErrMissingDomainQuery is an error that is returned when an input struct
// requires a "Query" key, but one was not set.
var ErrMissingDomainQuery = NewFieldError("Query")

// ErrMissingDomain is an error that is returned when an input struct
// requires a "Domain" key, but one was not set.
var ErrMissingDomain = NewFieldError("Domain")

// ErrMissingEntryID is an error that is returned when an input struct
// requires a "EntryID" key, but one was not set.
var ErrMissingEntryID = NewFieldError("EntryID")

// ErrMissingSnippetID is an error that is returned when an input struct
// requires a "SnippetID" key, but one was not set.
var ErrMissingSnippetID = NewFieldError("SnippetID")

// ErrMissingResourceID is an error that is returned when an input struct
// requires a "ResourceID" key, but one was not set.
var ErrMissingResourceID = NewFieldError("ResourceID")

// ErrMissingERLID is an error that is returned when an input struct
// requires an "ERLID" key, but one was not set.
var ErrMissingERLID = NewFieldError("ERLID")

// ErrMissingHost is an error that is returned when an input struct
// requires an "Host" key, but one was not set.
var ErrMissingHost = NewFieldError("Host")

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

// ErrMissingMetrics is an error that is returned when an input struct
// requires a "Metrics" key, but one was not set.
var ErrMissingMetrics = NewFieldError("Metrics")

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

// ErrMissingProject is an error that is returned when an input struct
// requires a "Project" key, but one was not set.
var ErrMissingProject = NewFieldError("Project")

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

// ErrMissingSite is an error that is returned when an input struct
// requires a "Site" key, but one was not set.
var ErrMissingSite = NewFieldError("Site")

// ErrMissingStart is an error that is returned when an input struct
// requires a "Start" key, but one was not set.
var ErrMissingStart = NewFieldError("Start")

// ErrMissingStoreID is an error that is returned when an input struct
// requires a "StoreID" key, but one was not set.
var ErrMissingStoreID = NewFieldError("StoreID")

// ErrMissingUserID is an error that is returned when an input struct
// requires a "UserID" key, but one was not set.
var ErrMissingUserID = NewFieldError("UserID")

// ErrMissingPermission is an error that is returned when an input struct
// requires a "Permission" key, but one was not set.
var ErrMissingPermission = NewFieldError("Permission")

// ErrInvalidPermission is an error that is returned when an input struct
// has a "Permission" key, but the one provided is invalid.
var ErrInvalidPermission = NewFieldError("Permission").Message("invalid")

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
// requires that the domain in "CommonName" is also in "Domains".
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

// ErrMissingUsername is an error that is returned when an input struct
// requires a "Username" key, but one was not set.
var ErrMissingUsername = NewFieldError("Username")

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

// ErrMissingWebhook is an error that is returned when an input
// struct requires a "Webhook" key, but one was not set.
var ErrMissingWebhook = NewFieldError("Webhook")

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
// Fastly API was not `{"status": "ok"}`.
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

// ErrInvalidMethod is an error that is returned when the input struct
// has an invalid Method value.
var ErrInvalidMethod = NewFieldError("Method").Message("invalid")

// ErrMissingCertificateMTLS is an error that is returned when an input struct
// requires either a "Certificate" or "MutualAuthentication" key, but neither
// was set.
var ErrMissingCertificateMTLS = NewFieldError("Certificate, MutualAuthentication").Message("at least one of the available optional fields is required")

// ErrMissingIntegrationID is an error that is returned when an input struct
// requires a "IntegrationID" key, but one was not set.
var ErrMissingIntegrationID = NewFieldError("IntegrationID")

// ErrMissingImageOptimizerSettings is an error that is returned when an input struct
// requires one of the optional Image Optimizer default settings, but none are set.
var ErrMissingImageOptimizerDefaultSetting = NewFieldError("ResizeFilter, Webp, WebpQuality, JpegType, JpegQuality, Upscale, AllowVideo").Message("at least one of the available optional fields is required")

// ErrMissingVirtualPatchID is an error that is returned when an input struct
// requires a "VirtualPatchID" key, but one was not set.
var ErrMissingVirtualPatchID = NewFieldError("Virtual Patch")

// ErrMissingMode is an error that is returned when an input struct
// requires a "Mode" key, but one was not set.
var ErrMissingMode = NewFieldError("Mode")

// Ensure HTTPError is, in fact, an error.
var _ error = (*HTTPError)(nil)

// HTTPError is a custom error type that wraps an HTTP status code with some
// helper functions.
type HTTPError struct {
	Errors []*ErrorObject `mapstructure:"errors"`
	// StatusCode is the HTTP status code (2xx-5xx).
	StatusCode int
	// RateLimitRemaining is the number of API requests remaining in the current
	// rate limit window. A `nil` value indicates the API returned no value for
	// the associated Fastly-RateLimit-Remaining response header.
	RateLimitRemaining *int
	// RateLimitReset is the time at which the current rate limit window resets,
	// as a Unix timestamp. A `nil` value indicates the API returned no value for
	// the associated Fastly-RateLimit-Reset response header.
	RateLimitReset *int
}

// ErrorObject is a single error.
type ErrorObject struct {
	Code   string          `mapstructure:"code" json:"code,omitempty"`
	Detail string          `mapstructure:"detail" json:"detail,omitempty"`
	ID     string          `mapstructure:"id" json:"id,omitempty"`
	Meta   *map[string]any `mapstructure:"meta" json:"meta,omitempty"`
	Status string          `mapstructure:"status" json:"status,omitempty"`
	Title  string          `mapstructure:"title" json:"title,omitempty"`
}

// legacyError represents the older-style errors from Fastly.
type legacyError struct {
	Errors  []map[string]any `mapstructure:"errors"`
	Detail  string           `mapstructure:"detail"`
	Message string           `mapstructure:"msg"`
	Title   string           `mapstructure:"title"`
}

// NewHTTPError creates a new HTTP error from the given code.
func NewHTTPError(resp *http.Response) *HTTPError {
	var e HTTPError
	e.StatusCode = resp.StatusCode

	if v, err := strconv.Atoi(resp.Header.Get("Fastly-RateLimit-Remaining")); err == nil {
		e.RateLimitRemaining = &v
	}
	if v, err := strconv.Atoi(resp.Header.Get("Fastly-RateLimit-Reset")); err == nil {
		e.RateLimitReset = &v
	}

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
		if err := DecodeBodyMap(body, &e); err != nil {
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
		if err := DecodeBodyMap(body, &lerr); err != nil {
			addDecodeErr()
		} else if lerr != nil {
			// This is for better handling the KV Store Bulk Insert endpoint.
			// https://developer.fastly.com/reference/api/services/resources/kv-store-item/#batch-create-keys
			if len(lerr.Errors) != 0 {
				for _, le := range lerr.Errors {
					var (
						code, detail string
						index        float64
					)
					// NOTE: We use `ok` second argument but _ it so as to avoid a panic.
					// Panics can happen if the service returns a 503 Service Unavailable.
					if c, ok := le["code"]; ok {
						code, _ = c.(string)
					}
					if r, ok := le["reason"]; ok {
						detail, _ = r.(string)
					}
					if d, ok := le["detail"]; ok {
						detail, _ = d.(string)
					}
					var title string
					if i, ok := le["index"]; ok {
						index, _ = i.(float64)
						title = fmt.Sprintf("error at index: %v", index)
					}
					if t, ok := le["title"]; ok {
						title, _ = t.(string)
					}
					e.Errors = append(e.Errors, &ErrorObject{
						Code:   code,
						Detail: detail,
						Title:  title,
					})
				}
			} else {
				msg := lerr.Message
				if msg == "" && lerr.Title != "" {
					msg = lerr.Title
				}
				e.Errors = append(e.Errors, &ErrorObject{
					Title:  msg,
					Detail: lerr.Detail,
				})
			}
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

	if e.RateLimitRemaining != nil {
		fmt.Fprintf(&b, "\n    RateLimitRemaining: %v", *e.RateLimitRemaining)
	}
	if e.RateLimitReset != nil {
		fmt.Fprintf(&b, "\n    RateLimitReset:     %v", *e.RateLimitReset)
	}

	return b.String()
}

// String implements the stringer interface and returns the string representing
// the error text that includes the status code and corresponding status text.
func (e *HTTPError) String() string {
	return e.Error()
}

// IsBadRequest returns true if the HTTP status code is 400, false otherwise.
func (e *HTTPError) IsBadRequest() bool {
	return e.StatusCode == http.StatusBadRequest
}

// IsNotFound returns true if the HTTP status code is 404, false otherwise.
func (e *HTTPError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsPreconditionFailed returns true if the HTTP status code is 412, false otherwise.
func (e *HTTPError) IsPreconditionFailed() bool {
	return e.StatusCode == http.StatusPreconditionFailed
}
