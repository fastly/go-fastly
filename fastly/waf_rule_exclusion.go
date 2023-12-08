package fastly

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
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
	Condition     *string    `jsonapi:"attr,condition"`
	CreatedAt     *time.Time `jsonapi:"attr,created_at,iso8601,omitempty"`
	ExclusionType *string    `jsonapi:"attr,exclusion_type"`
	ID            string     `jsonapi:"primary,waf_exclusion"`
	Name          *string    `jsonapi:"attr,name"`
	Number        *int       `jsonapi:"attr,number"`
	Rules         []*WAFRule `jsonapi:"relation,waf_rules,omitempty"`
	UpdatedAt     *time.Time `jsonapi:"attr,updated_at,iso8601,omitempty"`
}

// WAFRuleExclusionResponse represents a list of rule exclusions - full response.
type WAFRuleExclusionResponse struct {
	Info  infoResponse
	Items []*WAFRuleExclusion
}

// ListWAFRuleExclusionsInput used as input for listing a WAF's rule exclusions.
type ListWAFRuleExclusionsInput struct {
	// FilterExclusionType limits results to exclusions with the specified exclusions type.
	FilterExclusionType *string
	// FilterModSedID limits results to exclusions that represent the specified ModSecurity modsec_rule_id.
	FilterModSedID *string
	// FilterName limits results to exclusions with the specified exclusion name.
	FilterName *string
	// Include captures relationships. Optional. Permitted values: waf_rules.
	Include []string
	// PageNumber requests a specific page of exclusions.
	PageNumber *int
	// PageSize limits the number of returned pages.
	PageSize *int
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// ListAllWAFRuleExclusionsInput used as input for listing all WAF rule exclusions.
type ListAllWAFRuleExclusionsInput struct {
	// FilterExclusionType limits results to exclusions with the specified exclusions type.
	FilterExclusionType *string
	// FilterModSedID limits results to exclusions that represent the specified ModSecurity modsec_rule_id.
	FilterModSedID *string
	// FilterName limits results to exclusions with the specified exclusion name.
	FilterName *string
	// Include captures relationships. Optional. Permitted values: waf_rules.
	Include []string
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// CreateWAFRuleExclusionInput creates a new resource.
type CreateWAFRuleExclusionInput struct {
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFRuleExclusion is the Web Application Firewall's exclusion
	WAFRuleExclusion *WAFRuleExclusion
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// UpdateWAFRuleExclusionInput is used for exclusions updates.
type UpdateWAFRuleExclusionInput struct {
	// Number is the rule exclusion number.
	Number int
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFRuleExclusion is the Web Application Firewall's exclusion
	WAFRuleExclusion *WAFRuleExclusion
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// DeleteWAFRuleExclusionInput used as input for removing WAF rule exclusions.
type DeleteWAFRuleExclusionInput struct {
	// Number is the rule exclusion number.
	Number int
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

func (i *ListWAFRuleExclusionsInput) formatFilters() map[string]string {
	include := strings.Join(i.Include, ",")

	result := map[string]string{}
	pairings := map[string]any{
		"filter[exclusion_type]":           i.FilterExclusionType,
		"filter[name]":                     i.FilterName,
		"filter[waf_rules.modsec_rule_id]": i.FilterModSedID,
		"page[size]":                       i.PageSize,
		"page[number]":                     i.PageNumber,
		"include":                          include,
	}

	for key, value := range pairings {
		switch value := value.(type) {
		case string:
			if value != "" {
				result[key] = value
			}
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

// ListWAFRuleExclusions retrieves all resources.
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
	defer resp.Body.Close()

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

// ListAllWAFRuleExclusions retrieves all resources.
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
		ptr, err := c.ListWAFRuleExclusions(&ListWAFRuleExclusionsInput{
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

// CreateWAFRuleExclusion creates a new resource.
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
	defer resp.Body.Close()

	var wafExclusion WAFRuleExclusion
	if err := jsonapi.UnmarshalPayload(resp.Body, &wafExclusion); err != nil {
		return nil, err
	}
	return &wafExclusion, nil
}

// UpdateWAFRuleExclusion updates the specified resource.
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
	defer resp.Body.Close()

	var exc *WAFRuleExclusion
	if err := decodeBodyMap(resp.Body, &exc); err != nil {
		return nil, err
	}
	return exc, nil
}

// DeleteWAFRuleExclusion deletes the specified resource.
func (c *Client) DeleteWAFRuleExclusion(i *DeleteWAFRuleExclusionInput) error {
	if i.WAFID == "" {
		return ErrMissingWAFID
	}
	if i.WAFVersionNumber == 0 {
		return ErrMissingWAFVersionNumber
	}
	if i.Number == 0 {
		return ErrMissingNumber
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/exclusions/%d", i.WAFID, i.WAFVersionNumber, i.Number)
	_, err := c.Delete(path, nil)
	return err
}
