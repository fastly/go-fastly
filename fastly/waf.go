package fastly

import (
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
)

// WAFConfigurationSet represents information about a configuration_set.
type WAFConfigurationSet struct {
	ID string `jsonapi:"primary,configuration_set"`
}

// WAF is the information about a firewall object.
type WAF struct {
	ID                string `jsonapi:"primary,waf"`
	Version           int    `jsonapi:"attr,version"`
	PrefetchCondition string `jsonapi:"attr,prefetch_condition"`
	Response          string `jsonapi:"attr,response"`
	LastPush          string `jsonapi:"attr,last_push"`

	ConfigurationSet *WAFConfigurationSet `jsonapi:"relation,configuration_set"`
}

// wafType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var wafType = reflect.TypeOf(new(WAF))

// ListWAFsInput is used as input to the ListWAFs function.
type ListWAFsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListWAFs returns the list of wafs for the configuration version.
func (c *Client) ListWAFs(i *ListWAFsInput) ([]*WAF, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, wafType)
	if err != nil {
		return nil, err
	}

	wafs := make([]*WAF, len(data))
	for i := range data {
		typed, ok := data[i].(*WAF)
		if !ok {
			return nil, fmt.Errorf("got back a non-WAF response")
		}
		wafs[i] = typed
	}
	return wafs, nil
}

// CreateWAFInput is used as input to the CreateWAF function.
type CreateWAFInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	ID                string `jsonapi:"primary,waf"`
	PrefetchCondition string `jsonapi:"attr,prefetch_condition,omitempty"`
	Response          string `jsonapi:"attr,response,omitempty"`
}

// CreateWAF creates a new Fastly WAF.
func (c *Client) CreateWAF(i *CreateWAFInput) (*WAF, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs", i.Service, i.Version)
	resp, err := c.PostJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// GetWAFInput is used as input to the GetWAF function.
type GetWAFInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// ID is the id of the WAF to get.
	ID string
}

func (c *Client) GetWAF(i *GetWAFInput) (*WAF, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs/%s", i.Service, i.Version, i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// UpdateWAFInput is used as input to the UpdateWAF function.
type UpdateWAFInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	ID                string `jsonapi:"primary,waf"`
	PrefetchCondition string `jsonapi:"attr,prefetch_condition,omitempty"`
	Response          string `jsonapi:"attr,response,omitempty"`
}

// UpdateWAF updates a specific WAF.
func (c *Client) UpdateWAF(i *UpdateWAFInput) (*WAF, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs/%s", i.Service, i.Version, i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// DeleteWAFInput is used as input to the DeleteWAFInput function.
type DeleteWAFInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// ID is the id of the WAF to delete.
	ID string
}

func (c *Client) DeleteWAF(i *DeleteWAFInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.ID == "" {
		return ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs/%s", i.Service, i.Version, i.ID)
	_, err := c.Delete(path, nil)
	return err
}

// OWASP is the information about an OWASP object.
type OWASP struct {
	ID                            string `jsonapi:"primary,owasp"`
	AllowedHTTPVersions           string `jsonapi:"attr,allowed_http_versions"`
	AllowedMethods                string `jsonapi:"attr,allowed_methods"`
	AllowedRequestContentType     string `jsonapi:"attr,allowed_request_content_type"`
	ArgLength                     int    `jsonapi:"attr,arg_length"`
	ArgNameLength                 int    `jsonapi:"attr,arg_name_length"`
	CombinedFileSizes             int    `jsonapi:"attr,combined_file_sizes"`
	CreatedAt                     string `jsonapi:"attr,created_at"`
	CriticalAnomalyScore          int    `jsonapi:"attr,critical_anomaly_score"`
	CRSValidateUTF8Encoding       bool   `jsonapi:"attr,crs_validate_utf8_encoding"`
	ErrorAnomalyScore             int    `jsonapi:"attr,error_anomaly_score"`
	HighRiskCountryCodes          string `jsonapi:"attr,high_risk_country_codes"`
	HTTPViolationScoreThreshold   int    `jsonapi:"attr,http_violation_score_threshold"`
	InboundAnomalyScoreThreshold  int    `jsonapi:"attr,inbound_anomaly_score_threshold"`
	LFIScoreThreshold             int    `jsonapi:"attr,lfi_score_threshold"`
	MaxFileSize                   int    `jsonapi:"attr,max_file_size"`
	MaxNumArgs                    int    `jsonapi:"attr,max_num_args"`
	NoticeAnomalyScore            int    `jsonapi:"attr,notice_anomaly_score"`
	ParanoiaLevel                 int    `jsonapi:"attr,paranoia_level"`
	PHPInjectionScoreThreshold    int    `jsonapi:"attr,php_injection_score_threshold"`
	RCEScoreThreshold             int    `jsonapi:"attr,rce_score_threshold"`
	RestrictedExtensions          string `jsonapi:"attr,restricted_extensions"`
	RestrictedHeaders             string `jsonapi:"attr,restricted_headers"`
	RFIScoreThreshold             int    `jsonapi:"attr,rfi_score_threshold"`
	SessionFixationScoreThreshold int    `jsonapi:"attr,session_fixation_score_threshold"`
	SQLInjectionScoreThreshold    int    `jsonapi:"attr,sql_injection_score_threshold"`
	TotalArgLength                int    `jsonapi:"attr,total_arg_length"`
	UpdatedAt                     string `jsonapi:"attr,updated_at"`
	WarningAnomalyScore           int    `jsonapi:"attr,warning_anomaly_score"`
	XDDScoreThreshold             int    `jsonapi:"attr,xss_score_threshold"`
}

// GetOWASPInput is used as input to the GetOWASP function.
type GetOWASPInput struct {
	// Service is the ID of the service. WafID is the ID of the firewall.
	// Both fields are required.
	Service string
	ID      string
}

// GetOWASP gets OWASP settings for a service firewall object.
func (c *Client) GetOWASP(i *GetOWASPInput) (*OWASP, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/wafs/%s/owasp", i.Service, i.ID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var owasp OWASP
	if err := jsonapi.UnmarshalPayload(resp.Body, &owasp); err != nil {
		return nil, err
	}
	return &owasp, nil
}

// CreateOWASPInput is used as input to the CreateOWASP function.
type CreateOWASPInput struct {
	// Service is the ID of the service. ID is the ID of the firewall.
	// Both fields are required.
	Service string
	ID      string `jsonapi:"primary,owasp"`

	Type string `jsonapi:"attr,type`
}

// CreateOWASP creates an OWASP settings object for a service firewall object.
func (c *Client) CreateOWASP(i *CreateOWASPInput) (*OWASP, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/wafs/%s/owasp", i.Service, i.ID)
	resp, err := c.PostJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var owasp OWASP
	if err := jsonapi.UnmarshalPayload(resp.Body, &owasp); err != nil {
		return nil, err
	}
	return &owasp, nil
}

// CreateOWASPInput is used as input to the CreateOWASP function.
type UpdateOWASPInput struct {
	// Service is the ID of the service. WafID is the ID of the firewall.
	// Both fields are required.
	Service string
	ID      string
	OWASPID string `jsonapi:"primary,owasp,omitempty"`

	Type                          string `jsonapi:"attr,type`
	AllowedHTTPVersions           string `jsonapi:"attr,allowed_http_versions,omitempty"`
	AllowedMethods                string `jsonapi:"attr,allowed_methods,omitempty"`
	AllowedRequestContentType     string `jsonapi:"attr,allowed_request_content_type,omitempty"`
	ArgLength                     int    `jsonapi:"attr,arg_length,omitempty"`
	ArgNameLength                 int    `jsonapi:"attr,arg_name_length,omitempty"`
	CombinedFileSizes             int    `jsonapi:"attr,combined_file_sizes,omitempty"`
	CreatedAt                     string `jsonapi:"attr,created_at,omitempty"`
	CriticalAnomalyScore          int    `jsonapi:"attr,critical_anomaly_score,omitempty"`
	CRSValidateUTF8Encoding       bool   `jsonapi:"attr,crs_validate_utf8_encoding,omitempty"`
	ErrorAnomalyScore             int    `jsonapi:"attr,error_anomaly_score,omitempty"`
	HighRiskCountryCodes          string `jsonapi:"attr,high_risk_country_codes,omitempty"`
	HTTPViolationScoreThreshold   int    `jsonapi:"attr,http_violation_score_threshold,omitempty"`
	InboundAnomalyScoreThreshold  int    `jsonapi:"attr,inbound_anomaly_score_threshold,omitempty"`
	LFIScoreThreshold             int    `jsonapi:"attr,lfi_score_threshold,omitempty"`
	MaxFileSize                   int    `jsonapi:"attr,max_file_size,omitempty"`
	MaxNumArgs                    int    `jsonapi:"attr,max_num_args,omitempty"`
	NoticeAnomalyScore            int    `jsonapi:"attr,notice_anomaly_score,omitempty"`
	ParanoiaLevel                 int    `jsonapi:"attr,paranoia_level,omitempty"`
	PHPInjectionScoreThreshold    int    `jsonapi:"attr,php_injection_score_threshold,omitempty"`
	RCEScoreThreshold             int    `jsonapi:"attr,rce_score_threshold,omitempty"`
	RestrictedExtensions          string `jsonapi:"attr,restricted_extensions,omitempty"`
	RestrictedHeaders             string `jsonapi:"attr,restricted_headers,omitempty"`
	RFIScoreThreshold             int    `jsonapi:"attr,rfi_score_threshold,omitempty"`
	SessionFixationScoreThreshold int    `jsonapi:"attr,session_fixation_score_threshold,omitempty"`
	SQLInjectionScoreThreshold    int    `jsonapi:"attr,sql_injection_score_threshold,omitempty"`
	TotalArgLength                int    `jsonapi:"attr,total_arg_length,omitempty"`
	UpdatedAt                     string `jsonapi:"attr,updated_at,omitempty"`
	WarningAnomalyScore           int    `jsonapi:"attr,warning_anomaly_score,omitempty"`
	XDDScoreThreshold             int    `jsonapi:"attr,xss_score_threshold,omitempty"`
}

// CreateOWASP creates an OWASP settings object for a service firewall object.
func (c *Client) UpdateOWASP(i *UpdateOWASPInput) (*OWASP, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	if i.OWASPID == "" {
		return nil, ErrMissingOWASPID
	}

	path := fmt.Sprintf("/service/%s/wafs/%s/owasp", i.Service, i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var owasp OWASP
	if err := jsonapi.UnmarshalPayload(resp.Body, &owasp); err != nil {
		return nil, err
	}
	return &owasp, nil
}

// Rules is the information about an WAF rules.
type Rule struct {
	ID       string `jsonapi:"primary,rule"`
	RuleID   string `jsonapi:"attr,rule_id,omitempty"`
	Severity int    `jsonapi:"attr,severity,omitempty"`
	Message  string `jsonapi:"attr,message,omitempty"`
}

// rulesType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var rulesType = reflect.TypeOf(new(Rule))

// GetRules returns the list of wafs for the configuration version.
func (c *Client) GetRules() ([]*Rule, error) {
	path := fmt.Sprintf("/wafs/rules")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, rulesType)
	if err != nil {
		return nil, err
	}

	rules := make([]*Rule, len(data))
	for i := range data {
		typed, ok := data[i].(*Rule)
		if !ok {
			return nil, fmt.Errorf("got back a non-Rules response")
		}
		rules[i] = typed
	}

	return rules, nil
}

// GetRuleVCLInput is used as input to the GetRuleVCL function.
type GetRuleInput struct {
	// RuleID is the ID of the rule and is required.
	RuleID string
}

// GetRule gets a Rule using the Rule ID.
func (c *Client) GetRule(i *GetRuleInput) (*Rule, error) {
	if i.RuleID == "" {
		return nil, ErrMissingRuleID
	}

	path := fmt.Sprintf("/wafs/rules/%s", i.RuleID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var rule Rule
	if err := jsonapi.UnmarshalPayload(resp.Body, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// RuleVCL is the information about a Rule's VCL.
type RuleVCL struct {
	ID  string `jsonapi:"primary,rule_vcl"`
	VCL string `jsonapi:"attr,vcl,omitempty"`
}

// GetRuleVCL gets the VCL for a Rule.
func (c *Client) GetRuleVCL(i *GetRuleInput) (*RuleVCL, error) {
	if i.RuleID == "" {
		return nil, ErrMissingRuleID
	}

	path := fmt.Sprintf("/wafs/rules/%s/vcl", i.RuleID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var vcl RuleVCL
	if err := jsonapi.UnmarshalPayload(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return &vcl, nil
}

// GetWAFRuleVCLInput is used as input to the GetWAFRuleVCL function.
type GetWAFRuleVCLInput struct {
	// ID is the ID of the firewall. RuleID is the ID of the rule.
	// Both are required.
	ID     string
	RuleID string
}

// GetWAFRuleVCL gets the VCL for a role associated with a firewall WAF.
func (c *Client) GetWAFRuleVCL(i *GetWAFRuleVCLInput) (*RuleVCL, error) {
	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	if i.RuleID == "" {
		return nil, ErrMissingRuleID
	}

	path := fmt.Sprintf("/wafs/%s/rules/%s/vcl", i.ID, i.RuleID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var vcl RuleVCL
	if err := jsonapi.UnmarshalPayload(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return &vcl, nil
}

// Ruleset is the information about a firewall object's ruleset.
type Ruleset struct {
	ID       string `jsonapi:"primary,ruleset"`
	VCL      string `jsonapi:"attr,vcl,omitempty"`
	LastPush string `jsonapi:"attr,last_push,omitempty"`
}

// GetWAFRuleRuleSetsInput is used as input to the GetWAFRuleRuleSets function.
type GetWAFRuleRuleSetsInput struct {
	// Service is the ID of the service. ID is the ID of the firewall.
	// Both fields are required.
	Service string
	ID      string
}

// GetWAFRuleRuleSets gets the VCL for rulesets associated with a firewall WAF.
func (c *Client) GetWAFRuleRuleSets(i *GetWAFRuleRuleSetsInput) (*Ruleset, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/wafs/%s/ruleset", i.Service, i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ruleset Ruleset
	if err := jsonapi.UnmarshalPayload(resp.Body, &ruleset); err != nil {
		return nil, err
	}
	return &ruleset, nil
}

// UpdateWAFRuleRuleSetsInput is used as input to the UpdateWafRuleSets function.
type UpdateWAFRuleRuleSetsInput struct {
	// Service is the ID of the service. ID is the ID of the firewall.
	// Both fields are required.
	Service string
	ID      string `jsonapi:"primary,ruleset"`
}

// UpdateWafRuleSets updates the rulesets for a role associated with a firewall WAF.
func (c *Client) UpdateWafRuleSets(i *UpdateWAFRuleRuleSetsInput) (*Ruleset, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/wafs/%s/ruleset", i.Service, i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var ruleset Ruleset
	if err := jsonapi.UnmarshalPayload(resp.Body, &ruleset); err != nil {
		return nil, err
	}
	return &ruleset, nil
}
