package fastly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// WAFConfigurationSet represents information about a configuration_set.
type WAFConfigurationSet struct {
	ID string `jsonapi:"primary,configuration_set"`
}

// WAF  is the information about a firewall object.
type WAF struct {
	ActiveRulesFastlyBlockCount    int        `jsonapi:"attr,active_rules_fastly_block_count"`
	ActiveRulesFastlyLogCount      int        `jsonapi:"attr,active_rules_fastly_log_count"`
	ActiveRulesOWASPBlockCount     int        `jsonapi:"attr,active_rules_owasp_block_count"`
	ActiveRulesOWASPLogCount       int        `jsonapi:"attr,active_rules_owasp_log_count"`
	ActiveRulesOWASPScoreCount     int        `jsonapi:"attr,active_rules_owasp_score_count"`
	ActiveRulesTrustwaveBlockCount int        `jsonapi:"attr,active_rules_trustwave_block_count"`
	ActiveRulesTrustwaveLogCount   int        `jsonapi:"attr,active_rules_trustwave_log_count"`
	CreatedAt                      *time.Time `jsonapi:"attr,created_at,iso8601"`
	Disabled                       bool       `jsonapi:"attr,disabled"`
	ID                             string     `jsonapi:"primary,waf_firewall"`
	PrefetchCondition              string     `jsonapi:"attr,prefetch_condition"`
	Response                       string     `jsonapi:"attr,response"`
	ServiceID                      string     `jsonapi:"attr,service_id"`
	ServiceVersion                 int        `jsonapi:"attr,service_version_number"`
	UpdatedAt                      *time.Time `jsonapi:"attr,updated_at,iso8601"`
}

// WAFResponse an object containing the list of WAF results.
type WAFResponse struct {
	Info  infoResponse
	Items []*WAF
}

// wafType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var wafType = reflect.TypeOf(new(WAF))

// ListWAFsInput is used as input to the ListWAFs function.
type ListWAFsInput struct {
	// FilterService specifies the service ID of the returned firewalls.
	FilterService string
	// FilterVersion specifies the version of the service for the firewalls.
	FilterVersion int
	// Include captures relationships. Optional, comma-separated values. Permitted values: waf_firewall_versions.
	Include string
	// PageNumber requests a specific page of firewalls.
	PageNumber int
	// PageSize limits the number of returned firewalls.
	PageSize int
}

func (i *ListWAFsInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"page[size]":                     i.PageSize,
		"page[number]":                   i.PageNumber,
		"filter[service_id]":             i.FilterService,
		"filter[service_version_number]": i.FilterVersion,
		"include":                        i.Include,
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

// ListWAFs retrieves all resources.
func (c *Client) ListWAFs(i *ListWAFsInput) (*WAFResponse, error) {
	resp, err := c.Get("/waf/firewalls", &RequestOptions{
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
	data, err := jsonapi.UnmarshalManyPayload(bytes.NewReader(buf.Bytes()), wafType)
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

	return &WAFResponse{
		Items: wafs,
		Info:  info,
	}, nil
}

// CreateWAFInput is used as input to the CreateWAF function.
type CreateWAFInput struct {
	// ID is an alphanumeric string identifying a WAF Firewall.
	ID string `jsonapi:"primary,waf_firewall"`
	// PrefetchCondition is the name of the corresponding condition object.
	PrefetchCondition string `jsonapi:"attr,prefetch_condition"`
	// Response is the name of the corresponding response object.
	Response string `jsonapi:"attr,response"`
	// ServiceID is the ID of the service (required).
	ServiceID string `jsonapi:"attr,service_id"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `jsonapi:"attr,service_version_number"`
}

// CreateWAF creates a new resource.
func (c *Client) CreateWAF(i *CreateWAFInput) (*WAF, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := "/waf/firewalls"
	resp, err := c.PostJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// GetWAFInput is used as input to the GetWAF function.
type GetWAFInput struct {
	// ID is the WAF's ID.
	ID string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetWAF retrieves the specified resource.
func (c *Client) GetWAF(i *GetWAFInput) (*WAF, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("waf", "firewalls", i.ID)

	resp, err := c.Get(path, &RequestOptions{
		Params: map[string]string{
			"filter[service_version_number]": strconv.Itoa(i.ServiceVersion),
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// UpdateWAFInput is used as input to the UpdateWAF function.
type UpdateWAFInput struct {
	// Disabled is the status of the firewall.
	Disabled *bool `jsonapi:"attr,disabled,omitempty"`
	// ID is an alphanumeric string identifying a WAF Firewall.
	ID string `jsonapi:"primary,waf_firewall"`
	// PrefetchCondition is the name of the corresponding condition object.
	PrefetchCondition *string `jsonapi:"attr,prefetch_condition,omitempty"`
	// Response is the name of the corresponding response object.
	Response *string `jsonapi:"attr,response,omitempty"`
	// ServiceID is the ID of the service.
	ServiceID *string `jsonapi:"attr,service_id,omitempty"`
	// ServiceVersion is the specific configuration version.
	ServiceVersion *int `jsonapi:"attr,service_version_number,omitempty"`
}

// UpdateWAF updates the specified resource.
func (c *Client) UpdateWAF(i *UpdateWAFInput) (*WAF, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	// 'Service' and 'Version' are mandatory if:
	// 		- 'Disabled' is not specified.
	// 		- 'PrefetchCondition' or 'Response' are NOT empty.
	if i.Disabled == nil || i.PrefetchCondition != nil || i.Response != nil {
		if i.ServiceID == nil {
			return nil, ErrMissingServiceID
		}

		if i.ServiceVersion == nil {
			return nil, ErrMissingServiceVersion
		}
	}

	path := ToSafeURL("waf", "firewalls", i.ID)

	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// DeleteWAFInput is used as input to the DeleteWAFInput function.
type DeleteWAFInput struct {
	// This is the WAF ID.
	ID string `jsonapi:"primary,waf_firewall"`
	// The service version.
	ServiceVersion int `jsonapi:"attr,service_version_number"`
}

// DeleteWAF deletes the specified resource.
func (c *Client) DeleteWAF(i *DeleteWAFInput) error {
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.ID == "" {
		return ErrMissingID
	}

	path := ToSafeURL("waf", "firewalls", i.ID)

	_, err := c.DeleteJSONAPI(path, i, nil)
	return err
}

// infoResponse is used to pull the links and meta from the result.
type infoResponse struct {
	Links paginationInfo `json:"links"`
	Meta  metaInfo       `json:"meta"`
}

// paginationInfo stores links to searches related to the current one, showing
// any information about additional results being stored on another page.
type paginationInfo struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
}

// metaInfo stores information about the result returned by the server.
type metaInfo struct {
	CurrentPage int `json:"current_page,omitempty"`
	PerPage     int `json:"per_page,omitempty"`
	RecordCount int `json:"record_count,omitempty"`
	TotalPages  int `json:"total_pages,omitempty"`
}

// getResponseInfo parses a response to get the pagination and metadata info.
func getResponseInfo(body io.Reader) (infoResponse, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return infoResponse{}, err
	}

	var info infoResponse
	if err := json.Unmarshal(bodyBytes, &info); err != nil {
		return infoResponse{}, err
	}
	return info, nil
}
