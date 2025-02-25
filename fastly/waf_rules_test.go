package fastly

import (
	"reflect"
	"testing"

	"github.com/google/jsonapi"
)

func TestClient_WAF_Rules(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_rules/"

	var err error
	var rulesResp *WAFRuleResponse
	publisher := "owasp"
	Record(t, fixtureBase+"/list_owasp", func(c *Client) {
		rulesResp, err = c.ListWAFRules(&ListWAFRulesInput{
			FilterPublishers: []string{publisher},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesResp.Items) < 1 {
		t.Errorf("expected many rules: got %d", len(rulesResp.Items))
	}

	for _, r := range rulesResp.Items {
		if r.Publisher != publisher {
			t.Errorf("expected rule publisher %s: got %s", publisher, r.Publisher)
		}
	}

	publisher = "fastly"
	var fastlyRulesNumber int
	Record(t, fixtureBase+"/list_all_fastly", func(c *Client) {
		rulesResp, err = c.ListAllWAFRules(&ListAllWAFRulesInput{
			FilterPublishers: []string{publisher},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesResp.Items) < 1 {
		t.Errorf("expected many rules: got %d", len(rulesResp.Items))
	}

	for _, r := range rulesResp.Items {
		if r.Publisher != publisher {
			t.Errorf("expected rule publisher %s: got %s", publisher, r.Publisher)
		}
	}
	fastlyRulesNumber = len(rulesResp.Items)

	Record(t, fixtureBase+"/list_all_fastly_exclusion", func(c *Client) {
		rulesResp, err = c.ListAllWAFRules(&ListAllWAFRulesInput{
			FilterPublishers: []string{publisher},
			ExcludeMocSecIDs: []int{4170020},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesResp.Items) < 1 {
		t.Errorf("expected many rules: got %d", len(rulesResp.Items))
	}

	for _, r := range rulesResp.Items {
		if r.Publisher != publisher {
			t.Errorf("expected rule publisher %s: got %s", publisher, r.Publisher)
		}
	}

	if fastlyRulesNumber-1 != len(rulesResp.Items) {
		t.Errorf("expected %d rules: got %d", fastlyRulesNumber-1, len(rulesResp.Items))
	}

	Record(t, fixtureBase+"/list_all_fastly_filter_by_rule_ids", func(c *Client) {
		rulesResp, err = c.ListAllWAFRules(&ListAllWAFRulesInput{
			FilterModSecIDs: []int{1010060, 1010070},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesResp.Items) != 2 {
		t.Errorf("expected 2 rules: got %d", len(rulesResp.Items))
	}
}

func TestClient_listWAFRules_formatFilters(t *testing.T) {
	cases := []struct {
		remote *ListWAFRulesInput
		local  map[string]string
	}{
		{
			remote: &ListWAFRulesInput{
				FilterTagNames:   []string{"tag1", "tag2"},
				FilterPublishers: []string{"owasp", "trustwave"},
				FilterModSecIDs:  []int{1010060, 1010070},
				ExcludeModSecIDs: []int{123456, 1234567},
				PageSize:         2,
				PageNumber:       2,
				Include:          "included",
			},
			local: map[string]string{
				"filter[waf_tags][name][in]":  "tag1,tag2",
				"filter[publisher][in]":       "owasp,trustwave",
				"filter[modsec_rule_id][in]":  "1010060,1010070",
				"filter[modsec_rule_id][not]": "123456,1234567",
				jsonapi.QueryParamPageSize:    "2",
				jsonapi.QueryParamPageNumber:  "2",
				"include":                     "included",
			},
		},
	}
	for _, c := range cases {
		out := c.remote.formatFilters()
		if !reflect.DeepEqual(out, c.local) {
			t.Fatalf("Error matching:\nexpected: %#v\n     got: %#v", c.local, out)
		}
	}
}
