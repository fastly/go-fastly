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

// WAFRuleExclusionType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var WAFRuleExclusionType = reflect.TypeOf(new(WAFRuleExclusion))

const (
	// WAFRuleExclusionTypeRule is the type of WAF rule exclusions that excludes rules from the WAF based on certain conditions.
	WAFRuleExclusionTypeRule = "rule"
	// WAFRuleExclusionTypeWAF is the type of WAF rule exclusions that excludes WAF based on certain conditions.
	WAFRuleExclusionTypeWAF = "waf"
)

// WAFRuleExclusion is the information about a WAF rule exclusion object.
type WAFRuleExclusion struct {
	ID            string     `jsonapi:"primary,waf_exclusion"`
	Name          *string    `jsonapi:"attr,name"`
	ExclusionType *string    `jsonapi:"attr,exclusion_type"`
	Condition     *string    `jsonapi:"attr,condition"`
	Number        *int       `jsonapi:"attr,number"`
	Rules         []*WAFRule `jsonapi:"relation,waf_rules,omitempty"`
	CreatedAt     *time.Time `jsonapi:"attr,created_at,iso8601,omitempty"`
	UpdatedAt     *time.Time `jsonapi:"attr,updated_at,iso8601,omitempty"`
}

// WAFRuleExclusionResponse represents a list of rule exclusions - full response.
type WAFRuleExclusionResponse struct {
	Items []*WAFRuleExclusion
	Info  infoResponse
}

// ListWAFRuleExclusionsInput used as input for listing a WAF's rule exclusions.
type ListWAFRuleExclusionsInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// Limit results to exclusions with the specified exclusions type.
	FilterExclusionType *string
	// Limit results to exclusions with the specified exclusion name.
	FilterName *string
	// Limit results to exclusions that represent the specified ModSecurity modsec_rule_id.
	FilterModSedID *string
	// Limit the number of returned pages.
	PageSize *int
	// Request a specific page of exclusions.
	PageNumber *int
	// Include relationships. Optional, comma-separated values. Permitted values: waf_rule_revision and waf_firewall_version.
	Include *string
}

// ListAllWAFRuleExclusionsInput used as input for listing all WAF rule exclusions.
type ListAllWAFRuleExclusionsInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// Limit results to exclusions with the specified exclusions type.
	FilterExclusionType *string
	// Limit results to exclusions with the specified exclusion name.
	FilterName *string
	// Limit results to exclusions that represent the specified ModSecurity modsec_rule_id.
	FilterModSedID *string
	// Include relationships. Optional, comma-separated values. Permitted values: waf_rule_revision and waf_firewall_version.
	Include *string
}

// CreateWAFRuleExclusionInput used as input to create a WAF rule exclusion.
type CreateWAFRuleExclusionInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The Web Application Firewall's exclusion
	WAFRuleExclusion *WAFRuleExclusion
}

// UpdateWAFRuleExclusionInput is used for exclusions updates.
type UpdateWAFRuleExclusionInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The exclusion number.
	Number int
	// The WAF rule exclusion
	WAFRuleExclusion *WAFRuleExclusion
}

// DeleteWAFRuleExclusionInput used as input for removing WAF rule exclusions.
type DeleteWAFRuleExclusionInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The rule exclusion number.
	Number int
}

func (i *ListWAFRuleExclusionsInput) formatFilters() map[string]string {

	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[exclusion_type]":           i.FilterExclusionType,
		"filter[name]":                     i.FilterName,
		"filter[waf_rules.modsec_rule_id]": i.FilterModSedID,
		"page[size]":                       i.PageSize,
		"page[number]":                     i.PageNumber,
		"include":                          i.Include,
	}

	for key, value := range pairings {
		switch value := value.(type) {
		case *string:
			if value != nil {
				result[key] = *value
			}
		case *int:
			if value != nil {
				result[key] = strconv.Itoa(*value)
			}
		}
	}
	return result
}

// ListWAFRuleExclusions returns the list of exclusions for a given WAF ID.
func (c *Client) ListWAFRuleExclusions(i *ListWAFRuleExclusionsInput) (*WAFRuleExclusionResponse, error) {

	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/exclusions", i.WAFID, i.WAFVersionNumber)
	resp, err := c.Get(path, &RequestOptions{
		Params: i.formatFilters(),
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	info, err := getResponseInfo(tee)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(bytes.NewReader(buf.Bytes()), WAFRuleExclusionType)
	if err != nil {
		return nil, err
	}

	wafExclusions := make([]*WAFRuleExclusion, len(data))
	for i := range data {
		typed, ok := data[i].(*WAFRuleExclusion)
		if !ok {
			return nil, fmt.Errorf("got back a non-WAFRuleExclusion response")
		}
		wafExclusions[i] = typed
	}
	return &WAFRuleExclusionResponse{
		Items: wafExclusions,
		Info:  info,
	}, nil
}

// ListAllWAFRuleExclusions returns the complete list of WAF rule exclusions for a given WAF ID. It iterates through
// all existing pages to ensure all WAF rule exclusions are returned at once.
func (c *Client) ListAllWAFRuleExclusions(i *ListAllWAFRuleExclusionsInput) (*WAFRuleExclusionResponse, error) {

	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	currentPage := 1
	pageSize := WAFPaginationPageSize
	result := &WAFRuleExclusionResponse{Items: []*WAFRuleExclusion{}}
	for {
		r, err := c.ListWAFRuleExclusions(&ListWAFRuleExclusionsInput{
			WAFID:               i.WAFID,
			WAFVersionNumber:    i.WAFVersionNumber,
			PageNumber:          &currentPage,
			PageSize:            &pageSize,
			Include:             i.Include,
			FilterName:          i.FilterName,
			FilterModSedID:      i.FilterModSedID,
			FilterExclusionType: i.FilterExclusionType,
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

// CreateWAFRuleExclusion used to create a particular WAF rule exclusion.
func (c *Client) CreateWAFRuleExclusion(i *CreateWAFRuleExclusionInput) (*WAFRuleExclusion, error) {

	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	if i.WAFRuleExclusion == nil {
		return nil, ErrMissingWAFRuleExclusion
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/exclusions", i.WAFID, i.WAFVersionNumber)
	resp, err := c.PostJSONAPI(path, i.WAFRuleExclusion, nil)
	if err != nil {
		return nil, err
	}

	var wafExclusion WAFRuleExclusion
	if err := jsonapi.UnmarshalPayload(resp.Body, &wafExclusion); err != nil {
		return nil, err
	}
	return &wafExclusion, nil
}

// UpdateWAFRuleExclusion used to update a particular WAF rule exclusion.
func (c *Client) UpdateWAFRuleExclusion(i *UpdateWAFRuleExclusionInput) (*WAFRuleExclusion, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	if i.Number == 0 {
		return nil, ErrMissingWAFRuleExclusionNumber
	}

	if i.WAFRuleExclusion == nil {
		return nil, ErrMissingWAFRuleExclusion
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/exclusions/%d", i.WAFID, i.WAFVersionNumber, i.Number)
	resp, err := c.PatchJSONAPI(path, i.WAFRuleExclusion, nil)
	if err != nil {
		return nil, err
	}

	var exc *WAFRuleExclusion
	if err := decodeBodyMap(resp.Body, &exc); err != nil {
		return nil, err
	}
	return exc, nil
}

// DeleteWAFExclusions removes rules from a particular WAF.
func (c *Client) DeleteWAFRuleExclusion(i *DeleteWAFRuleExclusionInput) error {
	if i.WAFID == "" {
		return ErrMissingWAFID
	}
	if i.WAFVersionNumber == 0 {
		return ErrMissingWAFVersionNumber
	}
	if i.Number == 0 {
		return ErrMissingWAFRuleExclusionNumber
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/exclusions/%d", i.WAFID, i.WAFVersionNumber, i.Number)
	_, err := c.Delete(path, nil)
	return err
}
