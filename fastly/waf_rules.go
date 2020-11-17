package fastly

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/jsonapi"
)

// WAFRuleType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var WAFRuleType = reflect.TypeOf(new(WAFRule))

// WAFRule is the information about a WAF rule object.
type WAFRule struct {
	ID        string             `jsonapi:"primary,waf_rule,omitempty"`
	ModSecID  int                `jsonapi:"attr,modsec_rule_id,omitempty"`
	Publisher string             `jsonapi:"attr,publisher,omitempty"`
	Type      string             `jsonapi:"attr,type,omitempty"`
	Revisions []*WAFRuleRevision `jsonapi:"relation,waf_rule_revisions,omitempty"`
}

// WAFRuleRevision is the information about a WAF rule revision object.
type WAFRuleRevision struct {
	ID            string `jsonapi:"primary,waf_rule_revision,omitempty"`
	Status        string `jsonapi:"attr,message,omitempty"`
	Severity      int    `jsonapi:"attr,severity,omitempty"`
	Revision      int    `jsonapi:"attr,revision,omitempty"`
	ParanoiaLevel int    `jsonapi:"attr,paranoia_level,omitempty"`
	ModSecID      int    `jsonapi:"attr,modsec_rule_id,omitempty"`
	State         string `jsonapi:"attr,state,omitempty"`
	Source        string `jsonapi:"attr,source,omitempty"`
	VCL           string `jsonapi:"attr,vcl,omitempty"`
}

// WAFRuleResponse represents a list WAF rules full response.
type WAFRuleResponse struct {
	Items []*WAFRule
	Info  infoResponse
}

// ListWAFRulesInput used as input for listing WAF rules.
type ListWAFRulesInput struct {
	// Limit the returned rules to a set linked to list of tags by name.
	FilterTagNames []string
	// Limit the returned rules to a set by publishers.
	FilterPublishers []string
	// Excludes individual rules by modsecurity rule IDs.
	ExcludeMocSecIDs []int
	// Limit the number of returned rules.
	PageSize int
	// Request a specific page of rules.
	PageNumber int
	// Include relationships. Optional, comma-separated values. Permitted values: waf_tags and waf_rule_revisions.
	Include string
}

func (i *ListWAFRulesInput) formatFilters() map[string]string {

	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[waf_tags][name][in]":  i.FilterTagNames,
		"filter[publisher][in]":       i.FilterPublishers,
		"filter[modsec_rule_id][not]": i.ExcludeMocSecIDs,
		"page[size]":                  i.PageSize,
		"page[number]":                i.PageNumber,
		"include":                     i.Include,
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
		case "[]string":
			if len(value.([]string)) > 0 {
				result[key] = strings.Join(value.([]string), ",")
			}
		case "[]int":
			if len(value.([]int)) > 0 {

				stringSlice := make([]string, len(value.([]int)))
				for i, id := range value.([]int) {
					stringSlice[i] = strconv.Itoa(id)
				}
				result[key] = strings.Join(stringSlice, ",")
			}
		}
	}
	return result
}

// ListWAFRules returns the list of VAF versions for a given WAF ID.
func (c *Client) ListWAFRules(i *ListWAFRulesInput) (*WAFRuleResponse, error) {

	resp, err := c.Get("/waf/rules", &RequestOptions{
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

	data, err := jsonapi.UnmarshalManyPayload(bytes.NewReader(buf.Bytes()), WAFRuleType)
	if err != nil {
		return nil, err
	}

	wafRules := make([]*WAFRule, len(data))
	for i := range data {
		typed, ok := data[i].(*WAFRule)
		if !ok {
			return nil, fmt.Errorf("got back a non-WAFRule response")
		}
		wafRules[i] = typed
	}
	return &WAFRuleResponse{
		Items: wafRules,
		Info:  info,
	}, nil
}

// ListAllWAFRulesInput used as input for listing all WAF rules.
type ListAllWAFRulesInput struct {
	// Limit the returned rules to a set linked to a tag by name.
	FilterTagNames []string
	// Limit the returned rules to a set by publishers.
	FilterPublishers []string
	// Excludes individual rules by modsecurity rule IDs.
	ExcludeMocSecIDs []int
	// Include relationships. Optional, comma-separated values. Permitted values: waf_tags and waf_rule_revisions.
	Include string
}

// ListAllWAFRules returns the complete list of WAF rules for the given filters. It iterates through
// all existing pages to ensure all WAF rules are returned at once.
func (c *Client) ListAllWAFRules(i *ListAllWAFRulesInput) (*WAFRuleResponse, error) {

	currentPage := 1
	result := &WAFRuleResponse{Items: []*WAFRule{}}
	for {
		r, err := c.ListWAFRules(&ListWAFRulesInput{
			FilterTagNames:   i.FilterTagNames,
			FilterPublishers: i.FilterPublishers,
			ExcludeMocSecIDs: i.ExcludeMocSecIDs,
			Include:          i.Include,
			PageNumber:       currentPage,
			PageSize:         WAFPaginationPageSize,
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
