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

// WAFExclusionType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var WAFExclusionType = reflect.TypeOf(new(WAFExclusion))

const (
	// WAFExclusionTypeRule is the type of WAF exclusions that excludes rules from the WAF based on certain conditions
	WAFExclusionTypeRule = "rule"
	// WAFExclusionTypeWAF is the type of WAF exclusions that excludes WAF based on certain conditions
	WAFExclusionTypeWAF = "waf"
)

// WAFExclusion is the information about a WAF exclusion object.
type WAFExclusion struct {
	ID            string     `jsonapi:"primary,waf_exclusion"`
	Name          *string    `jsonapi:"attr,name"`
	ExclusionType *string    `jsonapi:"attr,exclusion_type"`
	Condition     *string    `jsonapi:"attr,condition"`
	Number        *int       `jsonapi:"attr,number"`
	Rules         []*WAFRule `jsonapi:"relation,waf_rules,omitempty"`
	CreatedAt     *time.Time `jsonapi:"attr,created_at,iso8601,omitempty"`
	UpdatedAt     *time.Time `jsonapi:"attr,updated_at,iso8601,omitempty"`
}

// WAFExclusionResponse represents a list of exclusions - full response.
type WAFExclusionResponse struct {
	Items []*WAFExclusion
	Info  infoResponse
}

// ListWAFExclusionsInput used as input for listing a WAF's exclusions.
type ListWAFExclusionsInput struct {
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

// ListAllWAFExclusionsInput used as input for listing all WAF exclusions.
type ListAllWAFExclusionsInput struct {
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

// CreateWAFExclusionInput used as input to create a WAF exclusion.
type CreateWAFExclusionInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The Web Application Firewall's exclusion
	WAFExclusion *WAFExclusion
}

// UpdateWAFExclusionInput is used for exclusions updates.
type UpdateWAFExclusionInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The exclusion number.
	Number int
	// The WAF exclusion
	WAFExclusion *WAFExclusion
}

// DeleteWAFExclusionInput used as input for removing WAF exclusions.
type DeleteWAFExclusionInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The exclusion number.
	Number int
}

func (i *ListWAFExclusionsInput) formatFilters() map[string]string {

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

// ListWAFExclusions returns the list of exclusions for a given WAF ID.
func (c *Client) ListWAFExclusions(i *ListWAFExclusionsInput) (*WAFExclusionResponse, error) {

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

	data, err := jsonapi.UnmarshalManyPayload(bytes.NewReader(buf.Bytes()), WAFExclusionType)
	if err != nil {
		return nil, err
	}

	wafExclusions := make([]*WAFExclusion, len(data))
	for i := range data {
		typed, ok := data[i].(*WAFExclusion)
		if !ok {
			return nil, fmt.Errorf("got back a non-WAFExclusion response")
		}
		wafExclusions[i] = typed
	}
	return &WAFExclusionResponse{
		Items: wafExclusions,
		Info:  info,
	}, nil
}

// ListAllWAFExclusions returns the complete list of WAF exclusions for a given WAF ID. It iterates through
// all existing pages to ensure all WAF exclusions are returned at once.
func (c *Client) ListAllWAFExclusions(i *ListAllWAFExclusionsInput) (*WAFExclusionResponse, error) {

	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	currentPage := 1
	pageSize := WAFPaginationPageSize
	result := &WAFExclusionResponse{Items: []*WAFExclusion{}}
	for {
		r, err := c.ListWAFExclusions(&ListWAFExclusionsInput{
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

// CreateWAFExclusion used to create a particular WAF exclusion.
func (c *Client) CreateWAFExclusion(i *CreateWAFExclusionInput) (*WAFExclusion, error) {

	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	if i.WAFExclusion == nil {
		return nil, ErrMissingWAFExclusion
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/exclusions", i.WAFID, i.WAFVersionNumber)
	resp, err := c.PostJSONAPI(path, i.WAFExclusion, nil)
	if err != nil {
		return nil, err
	}

	var wafExclusion WAFExclusion
	if err := jsonapi.UnmarshalPayload(resp.Body, &wafExclusion); err != nil {
		return nil, err
	}
	return &wafExclusion, nil
}

// UpdateWAFExclusion used to update a particular WAF exclusion.
func (c *Client) UpdateWAFExclusion(i *UpdateWAFExclusionInput) (*WAFExclusion, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	if i.Number == 0 {
		return nil, ErrMissingWAFExclusionNumber
	}

	if i.WAFExclusion == nil {
		return nil, ErrMissingWAFExclusion
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/exclusions/%d", i.WAFID, i.WAFVersionNumber, i.Number)
	resp, err := c.PatchJSONAPI(path, i.WAFExclusion, nil)
	if err != nil {
		return nil, err
	}

	var exc *WAFExclusion
	if err := decodeBodyMap(resp.Body, &exc); err != nil {
		return nil, err
	}
	return exc, nil
}

// DeleteWAFExclusions removes rules from a particular WAF.
func (c *Client) DeleteWAFExclusion(i *DeleteWAFExclusionInput) error {
	if i.WAFID == "" {
		return ErrMissingWAFID
	}
	if i.WAFVersionNumber == 0 {
		return ErrMissingWAFVersionNumber
	}
	if i.Number == 0 {
		return ErrMissingWAFExclusionNumber
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/exclusions/%d", i.WAFID, i.WAFVersionNumber, i.Number)
	_, err := c.Delete(path, nil)
	return err
}
