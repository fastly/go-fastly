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

// WAFActiveRuleType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var WAFActiveRuleType = reflect.TypeOf(new(WAFActiveRule))

// WAFActiveRule is the information about a WAF active rule object.
type WAFActiveRule struct {
	ID             string     `jsonapi:"primary,waf_active_rule,omitempty"`
	Status         string     `jsonapi:"attr,status,omitempty"`
	ModSecID       int        `jsonapi:"attr,modsec_rule_id,omitempty"`
	Revision       int        `jsonapi:"attr,revision,omitempty"`
	Outdated       bool       `jsonapi:"attr,outdated,omitempty"`
	LatestRevision int        `jsonapi:"attr,latest_revision,omitempty"`
	CreatedAt      *time.Time `jsonapi:"attr,created_at,iso8601,omitempty"`
	UpdatedAt      *time.Time `jsonapi:"attr,updated_at,iso8601,omitempty"`
}

// WAFActiveRuleResponse represents a list of active rules - full response.
type WAFActiveRuleResponse struct {
	Items []*WAFActiveRule
	Info  infoResponse
}

// ListWAFActiveRulesInput used as input for listing a WAF's active rules.
type ListWAFActiveRulesInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// Limit results to active rules with the specified status.
	FilterStatus string
	// Limit results to active rules with the specified message.
	FilterMessage string
	// Limit results to active rules that represent the specified ModSecurity modsec_rule_id.
	FilterModSedID string
	// Limit the number of returned pages.
	PageSize int
	// Request a specific page of active rules.
	PageNumber int
	// Include relationships. Optional, comma-separated values. Permitted values: waf_rule_revision and waf_firewall_version.
	Include string
}

func (i *ListWAFActiveRulesInput) formatFilters() map[string]string {

	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[status]":                            i.FilterStatus,
		"filter[waf_rule_revision][message]":        i.FilterMessage,
		"filter[waf_rule_revision][modsec_rule_id]": i.FilterModSedID,
		"page[size]":                                i.PageSize,
		"page[number]":                              i.PageNumber,
		"include":                                   i.Include,
	}

	for key, value := range pairings {
		switch value := value.(type) {
		case string:
			if value != "" {
				result[key] = value
			}
		case int:
			if value != 0 {
				result[key] = strconv.Itoa(value)
			}
		}
	}
	return result
}

// ListWAFActiveRules returns the list of active rules for a given WAF ID.
func (c *Client) ListWAFActiveRules(i *ListWAFActiveRulesInput) (*WAFActiveRuleResponse, error) {

	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/active-rules", i.WAFID, i.WAFVersionNumber)
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

	data, err := jsonapi.UnmarshalManyPayload(bytes.NewReader(buf.Bytes()), WAFActiveRuleType)
	if err != nil {
		return nil, err
	}

	wafRules := make([]*WAFActiveRule, len(data))
	for i := range data {
		typed, ok := data[i].(*WAFActiveRule)
		if !ok {
			return nil, fmt.Errorf("got back a non-WAFActiveRule response")
		}
		wafRules[i] = typed
	}
	return &WAFActiveRuleResponse{
		Items: wafRules,
		Info:  info,
	}, nil
}

// ListAllWAFActiveRulesInput used as input for listing all WAF active rules.
type ListAllWAFActiveRulesInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// Limit results to active rules with the specified status.
	FilterStatus string
	// Limit results to active rules with the specified message.
	FilterMessage string
	// Limit results to active rules that represent the specified ModSecurity modsec_rule_id.
	FilterModSedID string
	// Include relationships. Optional, comma-separated values. Permitted values: waf_rule_revision and waf_firewall_version.
	Include string
}

// ListAllWAFActiveRules returns the complete list of WAF active rules for a given WAF ID. It iterates through
// all existing pages to ensure all WAF active rules are returned at once.
func (c *Client) ListAllWAFActiveRules(i *ListAllWAFActiveRulesInput) (*WAFActiveRuleResponse, error) {

	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	currentPage := 1
	result := &WAFActiveRuleResponse{Items: []*WAFActiveRule{}}
	for {
		r, err := c.ListWAFActiveRules(&ListWAFActiveRulesInput{
			WAFID:            i.WAFID,
			WAFVersionNumber: i.WAFVersionNumber,
			PageNumber:       currentPage,
			PageSize:         WAFPaginationPageSize,
			Include:          i.Include,
			FilterStatus:     i.FilterStatus,
			FilterModSedID:   i.FilterModSedID,
			FilterMessage:    i.FilterMessage,
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

// CreateWAFActiveRulesInput used as input for adding rules to a WAF.
type CreateWAFActiveRulesInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The list of WAF active rules (ModSecID, Status and Revision are required).
	Rules []*WAFActiveRule
}

// CreateWAFActiveRules adds rules to a particular WAF.
func (c *Client) CreateWAFActiveRules(i *CreateWAFActiveRulesInput) ([]*WAFActiveRule, error) {

	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	if len(i.Rules) == 0 {
		return nil, ErrMissingWAFActiveRule
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/active-rules", i.WAFID, i.WAFVersionNumber)
	resp, err := c.PostJSONAPIBulk(path, i.Rules, nil)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, WAFActiveRuleType)
	if err != nil {
		return nil, err
	}

	wafRules := make([]*WAFActiveRule, len(data))
	for i := range data {
		typed, ok := data[i].(*WAFActiveRule)
		if !ok {
			return nil, fmt.Errorf("got back a non-WAFActiveRule response")
		}
		wafRules[i] = typed
	}

	return wafRules, nil
}

// BatchModificationWAFActiveRulesInput is used for active rules batch modifications.
type BatchModificationWAFActiveRulesInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The list of WAF active rules (ModSecID, Status and Revision are required for upsert, ModSecID is required for delete).
	Rules []*WAFActiveRule
	// The batch operation to be performed (allowed operations are upsert and delete).
	OP BatchOperation
}

// BatchModificationWAFActiveRules is a generic function for creating or deleting WAF active rules in batches.
// Upsert and delete are the only operations allowed.
func (c *Client) BatchModificationWAFActiveRules(i *BatchModificationWAFActiveRulesInput) ([]*WAFActiveRule, error) {

	if len(i.Rules) > BatchModifyMaximumOperations {
		return nil, ErrMaxExceededRules
	}

	switch i.OP {
	case UpsertBatchOperation:
		return c.CreateWAFActiveRules(&CreateWAFActiveRulesInput{
			WAFID:            i.WAFID,
			WAFVersionNumber: i.WAFVersionNumber,
			Rules:            i.Rules,
		})
	case DeleteBatchOperation:
		return nil, c.DeleteWAFActiveRules(&DeleteWAFActiveRulesInput{
			WAFID:            i.WAFID,
			WAFVersionNumber: i.WAFVersionNumber,
			Rules:            i.Rules,
		})
	default:
		return nil, fmt.Errorf("operation %s not supported", i.OP)
	}
}

// DeleteWAFActiveRulesInput used as input for removing rules from a WAF.
type DeleteWAFActiveRulesInput struct {
	// The Web Application Firewall's ID.
	WAFID string
	// The Web Application Firewall's version number.
	WAFVersionNumber int
	// The list of WAF active rules (ModSecID is required).
	Rules []*WAFActiveRule
}

// DeleteWAFActiveRules removes rules from a particular WAF.
func (c *Client) DeleteWAFActiveRules(i *DeleteWAFActiveRulesInput) error {

	if i.WAFID == "" {
		return ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return ErrMissingWAFVersionNumber
	}

	if len(i.Rules) == 0 {
		return ErrMissingWAFActiveRule
	}

	path := fmt.Sprintf("/waf/firewalls/%s/versions/%d/active-rules", i.WAFID, i.WAFVersionNumber)
	_, err := c.DeleteJSONAPIBulk(path, i.Rules, nil)
	return err
}
