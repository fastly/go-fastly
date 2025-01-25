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
	Revisions []*WAFRuleRevision `jsonapi:"relation,waf_rule_revisions,omitempty"`
	Type      string             `jsonapi:"attr,type,omitempty"`
}

// WAFRuleRevision is the information about a WAF rule revision object.
type WAFRuleRevision struct {
	ID            string `jsonapi:"primary,waf_rule_revision,omitempty"`
	ModSecID      int    `jsonapi:"attr,modsec_rule_id,omitempty"`
	ParanoiaLevel int    `jsonapi:"attr,paranoia_level,omitempty"`
	Revision      int    `jsonapi:"attr,revision,omitempty"`
	Severity      int    `jsonapi:"attr,severity,omitempty"`
	Source        string `jsonapi:"attr,source,omitempty"`
	State         string `jsonapi:"attr,state,omitempty"`
	Status        string `jsonapi:"attr,message,omitempty"`
	VCL           string `jsonapi:"attr,vcl,omitempty"`
}

// WAFRuleResponse represents a list WAF rules full response.
type WAFRuleResponse struct {
	Info  infoResponse
	Items []*WAFRule
}

// ListWAFRulesInput used as input for listing WAF rules.
type ListWAFRulesInput struct {
	// ExcludeModSecIDs excludes individual rules by modsecurity rule IDs.
	ExcludeModSecIDs []int
	// FilterModSecIDs limits the returned rules to a set by modsecurity rule IDs.
	FilterModSecIDs []int
	// FilterPublishers limits the returned rules to a set by publishers.
	FilterPublishers []string
	// FilterTagNames limits the returned rules to a set linked to list of tags by name.
	FilterTagNames []string
	// Include captures relationships. Optional, comma-separated values. Permitted values: waf_tags and waf_rule_revisions.
	Include string
	// PageNumber requests a specific page of rules.
	PageNumber int
	// PageSize limits the number of returned rules.
	PageSize int
}

func (i *ListWAFRulesInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"filter[waf_tags][name][in]":  i.FilterTagNames,
		"filter[publisher][in]":       i.FilterPublishers,
		"filter[modsec_rule_id][in]":  i.FilterModSecIDs,
		"filter[modsec_rule_id][not]": i.ExcludeModSecIDs,
		"page[size]":                  i.PageSize,
		"page[number]":                i.PageNumber,
		"include":                     i.Include,
	}

	for key, value := range pairings {
		switch v := value.(type) {
		case string:
			if v != "" {
				result[key] = v
			}
		case int:
			if v != 0 {
				result[key] = strconv.Itoa(v)
			}
		case []string:
			if len(v) > 0 {
				result[key] = strings.Join(v, ",")
			}
		case []int:
			if len(v) > 0 {
				stringSlice := make([]string, len(v))
				for i, id := range v {
					stringSlice[i] = strconv.Itoa(id)
				}
				result[key] = strings.Join(stringSlice, ",")
			}
		}
	}
	return result
}

// ListWAFRules retrieves all resources.
func (c *Client) ListWAFRules(i *ListWAFRulesInput) (*WAFRuleResponse, error) {
	resp, err := c.Get("/waf/rules", &RequestOptions{
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
	// ExcludeMocSecIDs excludes individual rules by modsecurity rule IDs.
	ExcludeMocSecIDs []int
	// FilterModSecIDs limits the returned rules to a set by modsecurity rule IDs.
	FilterModSecIDs []int
	// FilterPublishers limits the returned rules to a set by publishers.
	FilterPublishers []string
	// FilterTagNames limits the returned rules to a set linked to a tag by name.
	FilterTagNames []string
	// Include captures relationships. Optional, comma-separated values. Permitted values: waf_tags and waf_rule_revisions.
	Include string
}

// ListAllWAFRules retrieves all resources.
func (c *Client) ListAllWAFRules(i *ListAllWAFRulesInput) (*WAFRuleResponse, error) {
	currentPage := 1
	result := &WAFRuleResponse{Items: []*WAFRule{}}
	for {
		ptr, err := c.ListWAFRules(&ListWAFRulesInput{
			FilterTagNames:   i.FilterTagNames,
			FilterPublishers: i.FilterPublishers,
			FilterModSecIDs:  i.FilterModSecIDs,
			ExcludeModSecIDs: i.ExcludeMocSecIDs,
			Include:          i.Include,
			PageNumber:       currentPage,
			PageSize:         WAFPaginationPageSize,
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
