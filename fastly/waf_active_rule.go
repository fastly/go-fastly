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
	CreatedAt      *time.Time `jsonapi:"attr,created_at,iso8601,omitempty"`
	ID             string     `jsonapi:"primary,waf_active_rule,omitempty"`
	LatestRevision int        `jsonapi:"attr,latest_revision,omitempty"`
	ModSecID       int        `jsonapi:"attr,modsec_rule_id,omitempty"`
	Outdated       bool       `jsonapi:"attr,outdated,omitempty"`
	Revision       int        `jsonapi:"attr,revision,omitempty"`
	Status         string     `jsonapi:"attr,status,omitempty"`
	UpdatedAt      *time.Time `jsonapi:"attr,updated_at,iso8601,omitempty"`
}

// WAFActiveRuleResponse represents a list of active rules - full response.
type WAFActiveRuleResponse struct {
	Info  infoResponse
	Items []*WAFActiveRule
}

// ListWAFActiveRulesInput used as input for listing a WAF's active rules.
type ListWAFActiveRulesInput struct {
	// FilterMessage limits results to active rules with the specified message.
	FilterMessage string
	// FilterModSedID limits results to active rules that represent the specified ModSecurity modsec_rule_id.
	FilterModSedID string
	// FilterStatus limits results to active rules with the specified status.
	FilterStatus string
	// Include captures relationships. Optional, comma-separated values. Permitted values: waf_rule_revision and waf_firewall_version.
	Include string
	// PageNumber requests a specific page of active rules.
	PageNumber int
	// PageSize limits the number of returned pages.
	PageSize int
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

func (i *ListWAFActiveRulesInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"filter[status]":                            i.FilterStatus,
		"filter[waf_rule_revision][message]":        i.FilterMessage,
		"filter[waf_rule_revision][modsec_rule_id]": i.FilterModSedID,
		jsonapi.QueryParamPageSize:                  i.PageSize,
		jsonapi.QueryParamPageNumber:                i.PageNumber,
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

// ListWAFActiveRules retrieves all resources.
func (c *Client) ListWAFActiveRules(i *ListWAFActiveRulesInput) (*WAFActiveRuleResponse, error) {
	if i.WAFID == "" {
		return nil, ErrMissingWAFID
	}

	if i.WAFVersionNumber == 0 {
		return nil, ErrMissingWAFVersionNumber
	}

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions", strconv.Itoa(i.WAFVersionNumber), "active-rules")

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
	// FilterMessage limits results to active rules with the specified message.
	FilterMessage string
	// FilterModSedID limits results to active rules that represent the specified ModSecurity modsec_rule_id.
	FilterModSedID string
	// FilterStatus limits results to active rules with the specified status.
	FilterStatus string
	// Include captures relationships. Optional, comma-separated values. Permitted values: waf_rule_revision and waf_firewall_version.
	Include string
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// ListAllWAFActiveRules retrieves all resources.
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
		ptr, err := c.ListWAFActiveRules(&ListWAFActiveRulesInput{
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

// CreateWAFActiveRulesInput creates a new resource.
type CreateWAFActiveRulesInput struct {
	// Rules is the list of WAF active rules (ModSecID, Status and Revision are required).
	Rules []*WAFActiveRule
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// CreateWAFActiveRules creates a new resource.
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

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions", strconv.Itoa(i.WAFVersionNumber), "active-rules")

	resp, err := c.PostJSONAPIBulk(path, i.Rules, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	// OP is the batch operation to be performed (allowed operations are upsert and delete).
	OP BatchOperation
	// Rules is the list of WAF active rules (ModSecID, Status and Revision are required for upsert, ModSecID is required for delete).
	Rules []*WAFActiveRule
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// BatchModificationWAFActiveRules groups create/delete operations for the
// specified resource.
//
// This is a generic function for creating or deleting WAF active rules in
// batches. Upsert and delete are the only operations allowed.
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
	case CreateBatchOperation, UpdateBatchOperation:
		fallthrough
	default:
		return nil, fmt.Errorf("operation %s not supported", i.OP)
	}
}

// DeleteWAFActiveRulesInput used as input for removing rules from a WAF.
type DeleteWAFActiveRulesInput struct {
	// Rules is the list of WAF active rules (ModSecID is required).
	Rules []*WAFActiveRule
	// WAFID is the Web Application Firewall's ID.
	WAFID string
	// WAFVersionNumber is the Web Application Firewall's version number.
	WAFVersionNumber int
}

// DeleteWAFActiveRules deletes the specified resource.
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

	path := ToSafeURL("waf", "firewalls", i.WAFID, "versions", strconv.Itoa(i.WAFVersionNumber), "active-rules")

	ignored, err := c.DeleteJSONAPIBulk(path, i.Rules, nil)
	if err != nil {
		return err
	}
	defer ignored.Body.Close()
	return nil
}
