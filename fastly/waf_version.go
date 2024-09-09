package fastly

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// WAFVersionType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var WAFVersionType = reflect.TypeOf(new(WAFVersion))

const (
	// WAFPaginationPageSize is used as the default pagination page size by the WAF related requests.
	WAFPaginationPageSize = 100

	// WAFBatchModifyMaximumOperations is used as the default batch maximum operations.
	WAFBatchModifyMaximumOperations = 500

	// WAFVersionDeploymentStatusPending is the string value representing pending state for last_deployment_status for a WAF version.
	WAFVersionDeploymentStatusPending = "pending"

	// WAFVersionDeploymentStatusInProgress is the string value representing in-progress state for last_deployment_status for a WAF version.
	WAFVersionDeploymentStatusInProgress = "in progress"

	// WAFVersionDeploymentStatusCompleted is the string value representing completed state for last_deployment_status for a WAF versions.
	WAFVersionDeploymentStatusCompleted = "completed"

	// WAFVersionDeploymentStatusFailed is the string value representing failed state for last_deployment_status for a WAF versions.
	WAFVersionDeploymentStatusFailed = "failed"
)

// WAFVersion is the information about a WAF version object.
type WAFVersion struct {
	Active                           bool       `jsonapi:"attr,active"`
	ActiveRulesFastlyBlockCount      int        `jsonapi:"attr,active_rules_fastly_block_count"`
	ActiveRulesFastlyLogCount        int        `jsonapi:"attr,active_rules_fastly_log_count"`
	ActiveRulesOWASPBlockCount       int        `jsonapi:"attr,active_rules_owasp_block_count"`
	ActiveRulesOWASPLogCount         int        `jsonapi:"attr,active_rules_owasp_log_count"`
	ActiveRulesOWASPScoreCount       int        `jsonapi:"attr,active_rules_owasp_score_count"`
	ActiveRulesTrustwaveBlockCount   int        `jsonapi:"attr,active_rules_trustwave_block_count"`
	ActiveRulesTrustwaveLogCount     int        `jsonapi:"attr,active_rules_trustwave_log_count"`
	AllowedHTTPVersions              string     `jsonapi:"attr,allowed_http_versions"`
	AllowedMethods                   string     `jsonapi:"attr,allowed_methods"`
	AllowedRequestContentType        string     `jsonapi:"attr,allowed_request_content_type"`
	AllowedRequestContentTypeCharset string     `jsonapi:"attr,allowed_request_content_type_charset"`
	ArgLength                        int        `jsonapi:"attr,arg_length"`
	ArgNameLength                    int        `jsonapi:"attr,arg_name_length"`
	CRSValidateUTF8Encoding          bool       `jsonapi:"attr,crs_validate_utf8_encoding"`
	CombinedFileSizes                int        `jsonapi:"attr,combined_file_sizes"`
	Comment                          string     `jsonapi:"attr,comment"`
	CreatedAt                        *time.Time `jsonapi:"attr,created_at,iso8601"`
	CriticalAnomalyScore             int        `jsonapi:"attr,critical_anomaly_score"`
	DeployedAt                       *time.Time `jsonapi:"attr,deployed_at,iso8601"`
	Error                            string     `jsonapi:"attr,error"`
	ErrorAnomalyScore                int        `jsonapi:"attr,error_anomaly_score"`
	HTTPViolationScoreThreshold      int        `jsonapi:"attr,http_violation_score_threshold"`
	HighRiskCountryCodes             string     `jsonapi:"attr,high_risk_country_codes"`
	ID                               string     `jsonapi:"primary,waf_firewall_version"`
	InboundAnomalyScoreThreshold     int        `jsonapi:"attr,inbound_anomaly_score_threshold"`
	LFIScoreThreshold                int        `jsonapi:"attr,lfi_score_threshold"`
	LastDeploymentStatus             string     `jsonapi:"attr,last_deployment_status"`
	Locked                           bool       `jsonapi:"attr,locked"`
	MaxFileSize                      int        `jsonapi:"attr,max_file_size"`
	MaxNumArgs                       int        `jsonapi:"attr,max_num_args"`
	NoticeAnomalyScore               int        `jsonapi:"attr,notice_anomaly_score"`
	Number                           int        `jsonapi:"attr,number"`
	PHPInjectionScoreThreshold       int        `jsonapi:"attr,php_injection_score_threshold"`
	ParanoiaLevel                    int        `jsonapi:"attr,paranoia_level"`
	RCEScoreThreshold                int        `jsonapi:"attr,rce_score_threshold"`
	RFIScoreThreshold                int        `jsonapi:"attr,rfi_score_threshold"`
	RestrictedExtensions             string     `jsonapi:"attr,restricted_extensions"`
	RestrictedHeaders                string     `jsonapi:"attr,restricted_headers"`
	SQLInjectionScoreThreshold       int        `jsonapi:"attr,sql_injection_score_threshold"`
	SessionFixationScoreThreshold    int        `jsonapi:"attr,session_fixation_score_threshold"`
	TotalArgLength                   int        `jsonapi:"attr,total_arg_length"`
	UpdatedAt                        *time.Time `jsonapi:"attr,updated_at,iso8601"`
	WarningAnomalyScore              int        `jsonapi:"attr,warning_anomaly_score"`
	XSSScoreThreshold                int        `jsonapi:"attr,xss_score_threshold"`
}

// WAFVersionResponse represents a list WAF versions full response.
type WAFVersionResponse struct {
	Info  infoResponse
	Items []*WAFVersion
}

// ListWAFVersionsInput used as input for listing WAF versions.
type ListWAFVersionsInput struct {
	// Include captures relationships. Optional, comma-separated values. Permitted values: waf_firewall_versions.
	Include string
	// PageNumber requests a specific page of WAFs.
	PageNumber int
	// PageSize limits the number records returned.
	PageSize int
	// WAFID is the Web Application Firewall's ID.
	WAFID string
}

func (i *ListWAFVersionsInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"page[size]":   i.PageSize,
		"page[number]": i.PageNumber,
		"include":      i.Include,
	}

	for key, value := range pairings {
		switch t := reflect.TypeOf(value).String(); t {
		case "string":
			if value != "" {
				v, _ := value.(string) // type assert to avoid runtime panic (v will have zero value for its type)
				result[key] = v
			}
		case "int":
			if value != 0 {
				v, _ := value.(int) // type assert to avoid runtime panic (v will have zero value for its type)
				result[key] = strconv.Itoa(v)
			}
		}
	}
	return result
}

// ListWAFVersions retrieves all resources.
func (c *Client) ListWAFVersions(i *ListWAFVersionsInput) (*WAFVersionResponse, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions")

	resp, err := c.Get(path, &RequestOptions{
		Params: i.formatFilters(),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	info, err := getResponseInfo(tee)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(bytes.NewReader(buf.Bytes()), WAFVersionType)
	if err != nil {
		return nil, err
	}

	wafVersions := make([]*WAFVersion, len(data))
	for i := range data {
		typed, ok := data[i].(*WAFVersion)
		if !ok {
			return nil, fmt.Errorf("got back a non-WAFVersion response")
		}
		wafVersions[i] = typed
	}
	return &WAFVersionResponse{
		Items: wafVersions,
		Info:  info,
	}, nil
}

// ListAllWAFVersionsInput used as input for listing all WAF versions.
type ListAllWAFVersionsInput struct {
	// Include captures relationships. Optional, comma-separated values. Permitted values: waf_firewall_versions.
	Include string
	// WAFID is the Web Application Firewall's ID.
	WAFID string
}

// ListAllWAFVersions retrieves all resources.
func (c *Client) ListAllWAFVersions(i *ListAllWAFVersionsInput) (*WAFVersionResponse, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	currentPage := 1
	result := &WAFVersionResponse{Items: []*WAFVersion{}}
	for {
		ptr, err := c.ListWAFVersions(&ListWAFVersionsInput{
			WAFID:      i.WAFID,
			Include:    i.Include,
			PageNumber: currentPage,
			PageSize:   WAFPaginationPageSize,
		})
		if err != nil {
			return nil, err
		}
		if ptr == nil {
			return nil, fmt.Errorf("error: unexpected nil pointer")
		}

		currentPage++
		result.Items = append(result.Items, ptr.Items...)

		if ptr.Info.Links.Next == "" || len(ptr.Items) == 0 {
			return result, nil
		}
	}
}

// GetWAFVersionInput used as input for GetWAFVersion function.
type GetWAFVersionInput struct {
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// GetWAFVersion retrieves the specified resource.
func (c *Client) GetWAFVersion(i *GetWAFVersionInput) (*WAFVersion, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions", strconv.Itoa(i.WAFVersionNumber))

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wafVer WAFVersion
	if err := jsonapi.UnmarshalPayload(resp.Body, &wafVer); err != nil {
		return nil, err
	}
	return &wafVer, nil
}

// UpdateWAFVersionInput is used as input to the UpdateWAFVersion function. This struct uses pointers due to the problem
// detailed on this issue https://github.com/opencredo/go-fastly/pull/17.
type UpdateWAFVersionInput struct {
	// AllowedHTTPVersions is allowed HTTP versions.
	AllowedHTTPVersions *string `jsonapi:"attr,allowed_http_versions,omitempty"`
	// AllowedMethods is a space-separated list of HTTP method names.
	AllowedMethods *string `jsonapi:"attr,allowed_methods,omitempty"`
	// AllowedRequestContentType is allowed request content types.
	AllowedRequestContentType *string `jsonapi:"attr,allowed_request_content_type,omitempty"`
	// AllowedRequestContentTypeCharset is allowed request content type charset.
	AllowedRequestContentTypeCharset *string `jsonapi:"attr,allowed_request_content_type_charset,omitempty"`
	// ArgLength is the maximum allowed length of an argument.
	ArgLength *int `jsonapi:"attr,arg_length,omitempty"`
	// ArgNameLength is the maximum allowed argument name length.
	ArgNameLength *int `jsonapi:"attr,arg_name_length,omitempty"`
	// CRSValidateUTF8Encoding is the CRS validate UTF8 encoding.
	CRSValidateUTF8Encoding *bool `jsonapi:"attr,crs_validate_utf8_encoding,omitempty"`
	// CombinedFileSizes is the maximum allowed size of all files (in bytes).
	CombinedFileSizes *int `jsonapi:"attr,combined_file_sizes,omitempty"`
	// Comment is a freeform descriptive note.
	Comment *string `jsonapi:"attr,comment,omitempty"`
	// CriticalAnomalyScore is the score value to add for critical anomalies.
	CriticalAnomalyScore *int `jsonapi:"attr,critical_anomaly_score,omitempty"`
	// ErrorAnomalyScore is the score value to add for error anomalies.
	ErrorAnomalyScore *int `jsonapi:"attr,error_anomaly_score,omitempty"`
	// HTTPViolationScoreThreshold is the HTTP violation threshold.
	HTTPViolationScoreThreshold *int `jsonapi:"attr,http_violation_score_threshold,omitempty"`
	// HighRiskCountryCodes is a space-separated list of country codes in ISO 3166-1 (two-letter) format.
	HighRiskCountryCodes *string `jsonapi:"attr,high_risk_country_codes,omitempty"`
	// InboundAnomalyScoreThreshold is the inbound anomaly threshold.
	InboundAnomalyScoreThreshold *int `jsonapi:"attr,inbound_anomaly_score_threshold,omitempty"`
	// LFIScoreThreshold is the local file inclusion attack threshold.
	LFIScoreThreshold *int `jsonapi:"attr,lfi_score_threshold,omitempty"`
	// MaxFileSize is the maximum allowed file size, in bytes.
	MaxFileSize *int `jsonapi:"attr,max_file_size,omitempty"`
	// MaxNumArgs is the maximum number of arguments allowed.
	MaxNumArgs *int `jsonapi:"attr,max_num_args,omitempty"`
	// NoticeAnomalyScore is the score value to add for notice anomalies.
	NoticeAnomalyScore *int `jsonapi:"attr,notice_anomaly_score,omitempty"`
	// PHPInjectionScoreThreshold is the PHP injection threshold.
	PHPInjectionScoreThreshold *int `jsonapi:"attr,php_injection_score_threshold,omitempty"`
	// ParanoiaLevel is the configured paranoia level.
	ParanoiaLevel *int `jsonapi:"attr,paranoia_level,omitempty"`
	// RCEScoreThreshold is the remote code execution threshold.
	RCEScoreThreshold *int `jsonapi:"attr,rce_score_threshold,omitempty"`
	// RFIScoreThreshold is the remote file inclusion attack threshold.
	RFIScoreThreshold *int `jsonapi:"attr,rfi_score_threshold,omitempty"`
	// RestrictedExtensions is a space-separated list of allowed file extensions.
	RestrictedExtensions *string `jsonapi:"attr,restricted_extensions,omitempty"`
	// RestrictedHeaders is a space-separated list of allowed header names.
	RestrictedHeaders *string `jsonapi:"attr,restricted_headers,omitempty"`
	// SQLInjectionScoreThreshold is the SQL injection attack threshold.
	SQLInjectionScoreThreshold *int `jsonapi:"attr,sql_injection_score_threshold,omitempty"`
	// SessionFixationScoreThreshold is the session fixation attack threshold.
	SessionFixationScoreThreshold *int `jsonapi:"attr,session_fixation_score_threshold,omitempty"`
	// TotalArgLength is the maximum size of argument names and values.
	TotalArgLength *int `jsonapi:"attr,total_arg_length,omitempty"`
	// WAFID is the Web Application Firewall's ID.
	WAFID *string
	// WAFVersionID is the Web Application Firewall's version ID.
	WAFVersionID *string `jsonapi:"primary,waf_firewall_version"`
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber *int
	// WarningAnomalyScore is the score value to add for warning anomalies.
	WarningAnomalyScore *int `jsonapi:"attr,warning_anomaly_score,omitempty"`
	// XSSScoreThreshold is the XSS attack threshold.
	XSSScoreThreshold *int `jsonapi:"attr,xss_score_threshold,omitempty"`
}

// HasChanges checks that UpdateWAFVersionInput has changed in terms of configuration, which means - if it has configuration fields populated.
// if UpdateWAFVersionInput is updated to have a slice this method will not longer work as it is.
// if a slice is introduced the "!=" must be replaced with !DeepEquals.
func (i UpdateWAFVersionInput) HasChanges() bool {
	return i != UpdateWAFVersionInput{
		WAFID:            i.WAFID,
		WAFVersionNumber: i.WAFVersionNumber,
		WAFVersionID:     i.WAFVersionID,
	}
}

// UpdateWAFVersion updates the specified resource.
func (c *Client) UpdateWAFVersion(i *UpdateWAFVersionInput) (*WAFVersion, error) {
	if i.WAFID == nil || *i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == nil || *i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	if i.WAFVersionID == nil || *i.WAFVersionID == "" {
		return nil, ErrMissingWAFVersionID
	}

	path := ToSafeURL("waf", "firewalls", *i.WAFID, "versions", strconv.Itoa(*i.WAFVersionNumber))

	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var waf WAFVersion
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// LockWAFVersionInput used as input for locking a WAF version.
type LockWAFVersionInput struct {
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// LockWAFVersion locks a specific WAF version.
func (c *Client) LockWAFVersion(i *LockWAFVersionInput) (*WAFVersion, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions", strconv.Itoa(i.WAFVersionNumber))

	resp, err := c.PatchJSONAPI(path, &struct {
		Locked bool `jsonapi:"attr,locked"`
	}{true}, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var waf WAFVersion
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// CloneWAFVersionInput used as input for cloning a WAF version.
type CloneWAFVersionInput struct {
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// CloneWAFVersion clones a specific WAF version.
func (c *Client) CloneWAFVersion(i *CloneWAFVersionInput) (*WAFVersion, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions", strconv.Itoa(i.WAFVersionNumber), "clone")

	resp, err := c.PutJSONAPI(path, &CloneWAFVersionInput{}, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var waf WAFVersion
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// DeployWAFVersionInput used as input for deploying a WAF version.
type DeployWAFVersionInput struct {
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// DeployWAFVersion deploys a specific WAF version.
func (c *Client) DeployWAFVersion(i *DeployWAFVersionInput) error {
	if i.WAFID == "" {
		return ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return ErrMissingWAFVersionNumber
	}

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions", strconv.Itoa(i.WAFVersionNumber), "activate")

	_, err := c.PutJSONAPI(path, &DeployWAFVersionInput{}, nil)
	return err
}

// CreateEmptyWAFVersionInput creates a new resource.
type CreateEmptyWAFVersionInput struct {
	// WAFID is the Web Application Firewall's ID.
	WAFID string
}

// CreateEmptyWAFVersion creates a new resource.
//
// There are no rules and all config options are set to their default values.
func (c *Client) CreateEmptyWAFVersion(i *CreateEmptyWAFVersionInput) (*WAFVersion, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions")

	resp, err := c.PostJSONAPI(path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var waf WAFVersion
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}
