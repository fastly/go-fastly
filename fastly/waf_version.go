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
	// See documentation here https://developer.fastly.com/reference/api/waf/ngwaf/#api-section-ngwaf_firewall_versions
	ID                               string     `jsonapi:"primary,waf_firewall_version"`
	Number                           int        `jsonapi:"attr,number"`
	Active                           bool       `jsonapi:"attr,active"`
	Locked                           bool       `jsonapi:"attr,locked"`
	CRSValidateUTF8Encoding          bool       `jsonapi:"attr,crs_validate_utf8_encoding"`
	Comment                          string     `jsonapi:"attr,comment"`
	Error                            string     `jsonapi:"attr,error"`
	LastDeploymentStatus             string     `jsonapi:"attr,last_deployment_status"`
	AllowedHTTPVersions              string     `jsonapi:"attr,allowed_http_versions"`
	AllowedMethods                   string     `jsonapi:"attr,allowed_methods"`
	AllowedRequestContentType        string     `jsonapi:"attr,allowed_request_content_type"`
	AllowedRequestContentTypeCharset string     `jsonapi:"attr,allowed_request_content_type_charset"`
	HighRiskCountryCodes             string     `jsonapi:"attr,high_risk_country_codes"`
	RestrictedExtensions             string     `jsonapi:"attr,restricted_extensions"`
	RestrictedHeaders                string     `jsonapi:"attr,restricted_headers"`
	ArgLength                        int        `jsonapi:"attr,arg_length"`
	ArgNameLength                    int        `jsonapi:"attr,arg_name_length"`
	CombinedFileSizes                int        `jsonapi:"attr,combined_file_sizes"`
	CriticalAnomalyScore             int        `jsonapi:"attr,critical_anomaly_score"`
	ErrorAnomalyScore                int        `jsonapi:"attr,error_anomaly_score"`
	HTTPViolationScoreThreshold      int        `jsonapi:"attr,http_violation_score_threshold"`
	InboundAnomalyScoreThreshold     int        `jsonapi:"attr,inbound_anomaly_score_threshold"`
	LFIScoreThreshold                int        `jsonapi:"attr,lfi_score_threshold"`
	MaxFileSize                      int        `jsonapi:"attr,max_file_size"`
	MaxNumArgs                       int        `jsonapi:"attr,max_num_args"`
	NoticeAnomalyScore               int        `jsonapi:"attr,notice_anomaly_score"`
	ParanoiaLevel                    int        `jsonapi:"attr,paranoia_level"`
	PHPInjectionScoreThreshold       int        `jsonapi:"attr,php_injection_score_threshold"`
	RCEScoreThreshold                int        `jsonapi:"attr,rce_score_threshold"`
	RFIScoreThreshold                int        `jsonapi:"attr,rfi_score_threshold"`
	SessionFixationScoreThreshold    int        `jsonapi:"attr,session_fixation_score_threshold"`
	SQLInjectionScoreThreshold       int        `jsonapi:"attr,sql_injection_score_threshold"`
	TotalArgLength                   int        `jsonapi:"attr,total_arg_length"`
	WarningAnomalyScore              int        `jsonapi:"attr,warning_anomaly_score"`
	XSSScoreThreshold                int        `jsonapi:"attr,xss_score_threshold"`
	ActiveRulesTrustwaveLogCount     int        `jsonapi:"attr,active_rules_trustwave_log_count"`
	ActiveRulesTrustwaveBlockCount   int        `jsonapi:"attr,active_rules_trustwave_block_count"`
	ActiveRulesFastlyLogCount        int        `jsonapi:"attr,active_rules_fastly_log_count"`
	ActiveRulesFastlyBlockCount      int        `jsonapi:"attr,active_rules_fastly_block_count"`
	ActiveRulesOWASPLogCount         int        `jsonapi:"attr,active_rules_owasp_log_count"`
	ActiveRulesOWASPBlockCount       int        `jsonapi:"attr,active_rules_owasp_block_count"`
	ActiveRulesOWASPScoreCount       int        `jsonapi:"attr,active_rules_owasp_score_count"`
	DeployedAt                       *time.Time `jsonapi:"attr,deployed_at,iso8601"`
	CreatedAt                        *time.Time `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt                        *time.Time `jsonapi:"attr,updated_at,iso8601"`
}

// WAFVersionResponse represents a list WAF versions full response.
type WAFVersionResponse struct {
	Items []*WAFVersion
	Info  infoResponse
}

// ListWAFVersionsInput used as input for listing WAF versions.
type ListWAFVersionsInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// Limit the number records returned.
	PageSize int
	// Request a specific page of WAFs.
	PageNumber int
	// Include relationships. Optional, comma-separated values. Permitted values: waf_firewall_versions.
	Include string
}

func (i *ListWAFVersionsInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"page[size]":   i.PageSize,
		"page[number]": i.PageNumber,
		"include":      i.Include,
	}

	for key, value := range pairings {
		switch t := reflect.TypeOf(value).String(); t {
		case "string":
			if value != "" {
				result[key] = value.(string)
			}
		case "int":
			if value != 0 {
				result[key] = strconv.Itoa(value.(int))
			}
		}
	}
	return result
}

// ListWAFVersions returns the list of VAF versions for a given WAF ID.
func (c *Client) ListWAFVersions(i *ListWAFVersionsInput) (*WAFVersionResponse, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions", i.WAFID)
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
	// The Web Application Firewall's ID.
	WAFID string
	// Include relationships. Optional, comma-separated values. Permitted values: waf_firewall_versions.
	Include string
}

// ListAllWAFVersions returns the complete list of WAF versions for a given WAF ID. It iterates through
// all existing pages to ensure all WAF versions are returned at once.
func (c *Client) ListAllWAFVersions(i *ListAllWAFVersionsInput) (*WAFVersionResponse, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	currentPage := 1
	result := &WAFVersionResponse{Items: []*WAFVersion{}}
	for {
		r, err := c.ListWAFVersions(&ListWAFVersionsInput{
			WAFID:      i.WAFID,
			Include:    i.Include,
			PageNumber: currentPage,
			PageSize:   WAFPaginationPageSize,
		})
		if err != nil {
			return r, err
		}

		currentPage++
		result.Items = append(result.Items, r.Items...)

		if r.Info.Links.Next == "" || len(r.Items) == 0 {
			return result, nil
		}
	}
}

// GetWAFVersionInput used as input for GetWAFVersion function.
type GetWAFVersionInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
}

// GetWAFVersion gets details for given WAF version.
func (c *Client) GetWAFVersion(i *GetWAFVersionInput) (*WAFVersion, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d", i.WAFID, i.WAFVersionNumber)
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
	// The Web Application Firewall's ID.
	WAFID *string
	// The Web Application Firewall's version number.
	WAFVersionNumber *int
	// The Web Application Firewall's version ID.
	WAFVersionID                     *string `jsonapi:"primary,waf_firewall_version"`
	Comment                          *string `jsonapi:"attr,comment,omitempty"`
	CRSValidateUTF8Encoding          *bool   `jsonapi:"attr,crs_validate_utf8_encoding,omitempty"`
	AllowedHTTPVersions              *string `jsonapi:"attr,allowed_http_versions,omitempty"`
	AllowedMethods                   *string `jsonapi:"attr,allowed_methods,omitempty"`
	AllowedRequestContentType        *string `jsonapi:"attr,allowed_request_content_type,omitempty"`
	AllowedRequestContentTypeCharset *string `jsonapi:"attr,allowed_request_content_type_charset,omitempty"`
	HighRiskCountryCodes             *string `jsonapi:"attr,high_risk_country_codes,omitempty"`
	RestrictedExtensions             *string `jsonapi:"attr,restricted_extensions,omitempty"`
	RestrictedHeaders                *string `jsonapi:"attr,restricted_headers,omitempty"`
	ArgLength                        *int    `jsonapi:"attr,arg_length,omitempty"`
	ArgNameLength                    *int    `jsonapi:"attr,arg_name_length,omitempty"`
	CombinedFileSizes                *int    `jsonapi:"attr,combined_file_sizes,omitempty"`
	CriticalAnomalyScore             *int    `jsonapi:"attr,critical_anomaly_score,omitempty"`
	ErrorAnomalyScore                *int    `jsonapi:"attr,error_anomaly_score,omitempty"`
	HTTPViolationScoreThreshold      *int    `jsonapi:"attr,http_violation_score_threshold,omitempty"`
	InboundAnomalyScoreThreshold     *int    `jsonapi:"attr,inbound_anomaly_score_threshold,omitempty"`
	LFIScoreThreshold                *int    `jsonapi:"attr,lfi_score_threshold,omitempty"`
	MaxFileSize                      *int    `jsonapi:"attr,max_file_size,omitempty"`
	MaxNumArgs                       *int    `jsonapi:"attr,max_num_args,omitempty"`
	NoticeAnomalyScore               *int    `jsonapi:"attr,notice_anomaly_score,omitempty"`
	ParanoiaLevel                    *int    `jsonapi:"attr,paranoia_level,omitempty"`
	PHPInjectionScoreThreshold       *int    `jsonapi:"attr,php_injection_score_threshold,omitempty"`
	RCEScoreThreshold                *int    `jsonapi:"attr,rce_score_threshold,omitempty"`
	RFIScoreThreshold                *int    `jsonapi:"attr,rfi_score_threshold,omitempty"`
	SessionFixationScoreThreshold    *int    `jsonapi:"attr,session_fixation_score_threshold,omitempty"`
	SQLInjectionScoreThreshold       *int    `jsonapi:"attr,sql_injection_score_threshold,omitempty"`
	TotalArgLength                   *int    `jsonapi:"attr,total_arg_length,omitempty"`
	WarningAnomalyScore              *int    `jsonapi:"attr,warning_anomaly_score,omitempty"`
	XSSScoreThreshold                *int    `jsonapi:"attr,xss_score_threshold,omitempty"`
}

// HasChanges checks that UpdateWAFVersionInput has changed in terms of configuration, which means - if it has configuration fields populated.
// if UpdateWAFVersionInput is updated to have a slice this method will not longer work as it is.
// if a slice is introduced the "!=" must be replaced with !DeepEquals
func (i UpdateWAFVersionInput) HasChanges() bool {
	return i != UpdateWAFVersionInput{
		WAFID:            i.WAFID,
		WAFVersionNumber: i.WAFVersionNumber,
		WAFVersionID:     i.WAFVersionID,
	}
}

// UpdateWAFVersion updates a specific WAF version.
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

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d", *i.WAFID, *i.WAFVersionNumber)
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
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
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

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d", i.WAFID, i.WAFVersionNumber)
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
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
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

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/clone", i.WAFID, i.WAFVersionNumber)
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
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
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

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/activate", i.WAFID, i.WAFVersionNumber)
	_, err := c.PutJSONAPI(path, &DeployWAFVersionInput{}, nil)
	if err != nil {
		return err
	}
	return nil
}

// CreateEmptyWAFVersionInput used as input for creating an empty WAF version.
type CreateEmptyWAFVersionInput struct {
	// The Web Application Firewall's ID.
	WAFID string
}

// CreateEmptyWAFVersion creates an empty WAF version,
//  which means a version without rules and all config options set to their default values.
func (c *Client) CreateEmptyWAFVersion(i *CreateEmptyWAFVersionInput) (*WAFVersion, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions", i.WAFID)
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
