package fastly

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strconv"
	"testing"
)

func TestClient_WAFs(t *testing.T) {
	t.Parallel()

	fixtureBase := "wafs/"

	testService := createTestService(t, fixtureBase+"create_service", "service")
	defer deleteTestService(t, fixtureBase+"delete_service", testService.ID)

	tv := createTestVersion(t, fixtureBase+"/version", testService.ID)

	var err error
	// Enable logging on the service - we cannot create wafs without logging
	// enabled
	record(t, fixtureBase+"/logging/create", func(c *Client) {
		_, err = c.CreateSyslog(&CreateSyslogInput{
			Service:       testService.ID,
			Version:       tv.Number,
			Name:          "test-syslog",
			Address:       "example.com",
			Hostname:      "example.com",
			Port:          1234,
			Token:         "abcd1234",
			Format:        "format",
			FormatVersion: 2,
			MessageType:   "classic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		record(t, fixtureBase+"/logging/cleanup", func(c *Client) {
			c.DeleteSyslog(&DeleteSyslogInput{
				Service: testService.ID,
				Version: tv.Number,
				Name:    "test-syslog",
			})
		})
	}()

	// Create a condition - we cannot create a waf without attaching a condition
	var condition *Condition
	record(t, fixtureBase+"/condition/create", func(c *Client) {
		condition, err = c.CreateCondition(&CreateConditionInput{
			Service:   testService.ID,
			Version:   tv.Number,
			Name:      "WAF_Prefetch",
			Statement: "req.url~+\"index.html\"",
			Type:      "PREFETCH", // This must be a prefetch condition
			Priority:  1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		record(t, fixtureBase+"/condition/cleanup", func(c *Client) {
			c.DeleteCondition(&DeleteConditionInput{
				Service: testService.ID,
				Version: tv.Number,
				Name:    condition.Name,
			})
		})
	}()

	// Create a response object
	var ro *ResponseObject
	record(t, fixtureBase+"/response_object/create", func(c *Client) {
		ro, err = c.CreateResponseObject(&CreateResponseObjectInput{
			Service:     testService.ID,
			Version:     tv.Number,
			Name:        "WAf_Response",
			Status:      200,
			Response:    "Ok",
			Content:     "abcd",
			ContentType: "text/plain",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		record(t, fixtureBase+"/response_object/cleanup", func(c *Client) {
			c.DeleteResponseObject(&DeleteResponseObjectInput{
				Service: testService.ID,
				Version: tv.Number,
				Name:    ro.Name,
			})
		})
	}()

	// Create
	var waf *WAF
	record(t, fixtureBase+"/create", func(c *Client) {
		waf, err = c.CreateWAF(&CreateWAFInput{
			Service:           testService.ID,
			Version:           strconv.Itoa(tv.Number),
			PrefetchCondition: condition.Name,
			Response:          ro.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// List
	var wafsResp *WAFResponse
	record(t, fixtureBase+"/list", func(c *Client) {
		wafsResp, err = c.ListWAFs(&ListWAFsInput{
			FilterService: testService.ID,
			FilterVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wafsResp.Items) < 0 {
		t.Errorf("bad wafs: %v", wafsResp.Items)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"/cleanup", func(c *Client) {
			c.DeleteWAF(&DeleteWAFInput{
				Version: strconv.Itoa(tv.Number),
				ID:      waf.ID,
			})
		})
	}()

	// Get
	var nwaf *WAF
	record(t, fixtureBase+"/get", func(c *Client) {
		nwaf, err = c.GetWAF(&GetWAFInput{
			Service: testService.ID,
			Version: strconv.Itoa(tv.Number),
			ID:      waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if nwaf.ID != waf.ID {
		t.Errorf("expected %q to be %q", nwaf.ID, waf.ID)
	}

	// Get the WAF ruleset
	var ruleset *Ruleset
	record(t, "wafs/ruleset/get", func(c *Client) {
		ruleset, err = c.GetWAFRuleRuleSets(&GetWAFRuleRuleSetsInput{
			Service: testServiceID,
			ID:      waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ruleset.ID != waf.ID {
		t.Errorf("expected %q to be %q", ruleset.ID, waf.ID)
	}
	if ruleset.VCL == "" {
		t.Error("bad vcl")
	}
	if ruleset.LastPush != waf.LastPush {
		t.Errorf("expected %q to be %q", ruleset.LastPush, waf.LastPush)
	}
	if ruleset.Link != "" {
		t.Error("bad link")
	}

	// Update
	// Create a new response object to attach
	var nro *ResponseObject
	record(t, fixtureBase+"/response_object/create_another", func(c *Client) {
		nro, err = c.CreateResponseObject(&CreateResponseObjectInput{
			Service:     testService.ID,
			Version:     tv.Number,
			Name:        "test-response-object-2",
			Status:      200,
			Response:    "Ok",
			Content:     "efgh",
			ContentType: "text/plain",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		record(t, fixtureBase+"/response_object/cleanup_another", func(c *Client) {
			c.DeleteResponseObject(&DeleteResponseObjectInput{
				Service: testService.ID,
				Version: tv.Number,
				Name:    nro.Name,
			})
		})
	}()
	var uwaf *WAF
	record(t, fixtureBase+"/update", func(c *Client) {
		uwaf, err = c.UpdateWAF(&UpdateWAFInput{
			Service:  testService.ID,
			Version:  strconv.Itoa(tv.Number),
			ID:       waf.ID,
			Response: nro.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uwaf.Response != "test-response-object-2" {
		t.Errorf("bad name: %q", uwaf.Response)
	}

	// Update the WAF ruleset
	var uRuleset *Ruleset
	record(t, "wafs/ruleset/patch", func(c *Client) {
		uRuleset, err = c.UpdateWAFRuleSets(&UpdateWAFRuleRuleSetsInput{
			Service: testServiceID,
			ID:      uwaf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uRuleset.ID != uwaf.ID {
		t.Errorf("expected %q to be %q", uRuleset.ID, uwaf.ID)
	}
	if uRuleset.VCL != "" {
		t.Error("bad vcl")
	}
	if uRuleset.LastPush != nil {
		t.Error("bad last_push")
	}
	if uRuleset.Link == "" {
		t.Error("bad link")
	}

	// Delete
	record(t, fixtureBase+"/delete", func(c *Client) {
		err = c.DeleteWAF(&DeleteWAFInput{
			Version: strconv.Itoa(tv.Number),
			ID:      waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestClient_CreateWAF_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateWAF(&CreateWAFInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateWAF(&CreateWAFInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetWAF_validation(t *testing.T) {
	var err error
	_, err = testClient.GetWAF(&GetWAFInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetWAF(&GetWAFInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetWAF(&GetWAFInput{
		Service: "foo",
		Version: "1",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAF_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateWAF(&UpdateWAFInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateWAF(&UpdateWAFInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateWAF(&UpdateWAFInput{
		Service: "foo",
		Version: "1",
		ID:      "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteWAF_validation(t *testing.T) {
	var err error
	err = testClient.DeleteWAF(&DeleteWAFInput{
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteWAF(&DeleteWAFInput{
		Version: "1",
		ID:      "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_listWAFs_formatFilters(t *testing.T) {
	cases := []struct {
		remote *ListWAFsInput
		local  map[string]string
	}{
		{
			remote: &ListWAFsInput{
				FilterService: "service1",
				FilterVersion: 1,
			},
			local: map[string]string{
				"filter[service_id]":             "service1",
				"filter[service_version_number]": "1",
			},
		},
		{
			remote: &ListWAFsInput{
				FilterService: "service1",
				FilterVersion: 1,
				PageSize:      2,
				PageNumber:    2,
				Include:       "included",
			},
			local: map[string]string{
				"filter[service_id]":             "service1",
				"filter[service_version_number]": "1",
				"page[size]":                     "2",
				"page[number]":                   "2",
				"include":                        "included",
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

func TestUpdateWAFRuleStatusesInput_validate(t *testing.T) {
	tests := []struct {
		description string
		input       UpdateWAFRuleStatusInput
		expected    error
	}{
		{
			description: "Accepts valid input",
			input: UpdateWAFRuleStatusInput{
				ID:      "as098k-8104",
				RuleID:  8104,
				Service: "108asj1",
				WAF:     "as098k",
				Status:  "block",
			},
			expected: nil,
		},
		{
			description: "Rejects input with missing int field",
			input: UpdateWAFRuleStatusInput{
				ID:      "as098k-8104",
				Service: "108asj1",
				WAF:     "as098k",
				Status:  "block",
			},
			expected: ErrMissingRuleID,
		},
		{
			description: "Rejects input with missing string field",
			input: UpdateWAFRuleStatusInput{
				ID:     "as098k-8104",
				RuleID: 8104,
				WAF:    "as098k",
				Status: "block",
			},
			expected: ErrMissingService,
		},
	}
	for _, testcase := range tests {
		err := testcase.input.validate()
		if err != testcase.expected {
			t.Errorf("In test %s: Expected %v,got %v", testcase.description, testcase.expected, err)
		}
	}
}

func TestUpdateWAFRuleTagStatusInput_validate(t *testing.T) {
	tests := []struct {
		description string
		input       UpdateWAFRuleTagStatusInput
		expected    error
	}{
		{
			description: "Accepts valid input",
			input: UpdateWAFRuleTagStatusInput{
				Tag:     "lala tag la",
				Service: "108asj1",
				WAF:     "as098k",
				Status:  "block",
			},
			expected: nil,
		},
		{
			description: "Rejects input with missing string field",
			input: UpdateWAFRuleTagStatusInput{
				Service: "108asj1",
				WAF:     "as098k",
				Status:  "block",
			},
			expected: ErrMissingTag,
		},
	}
	for _, testcase := range tests {
		err := testcase.input.validate()
		if err != testcase.expected {
			t.Errorf("In test %s: Expected %v,got %v", testcase.description, testcase.expected, err)
		}
	}
}

func TestGetWAFRuleStatusesInput_formatFilters(t *testing.T) {
	tests := []struct {
		description string
		filters     GetWAFRuleStatusesFilters
		expected    map[string]string
	}{
		{
			description: "converts both strings and ints to strings",
			filters: GetWAFRuleStatusesFilters{
				Status:   "log",
				Accuracy: 10,
				Version:  "180ad",
			},
			expected: map[string]string{
				"filter[status]":         "log",
				"filter[rule][accuracy]": "10",
				"filter[rule][version]":  "180ad",
			},
		},
		{
			description: "converts arrays to strings",
			filters: GetWAFRuleStatusesFilters{
				Status:  "log",
				Version: "181ad",
				Tags:    []int{18, 1, 1093, 86308},
			},
			expected: map[string]string{
				"filter[status]":        "log",
				"filter[rule][version]": "181ad",
				"include":               "18,1,1093,86308",
			},
		},
	}
	for _, testcase := range tests {
		input := GetWAFRuleStatusesInput{
			Filters: testcase.filters,
		}
		answer := input.formatFilters()
		if len(answer) != len(testcase.expected) {
			t.Errorf("In test %s: Expected map with %d entries,got one with %d", testcase.description, len(testcase.expected), len(answer))
		}
		for key, value := range testcase.expected {
			if answer[key] != value {
				t.Errorf("In test %s: Expected %s key to have value %s, got %s", testcase.description, key, value, answer[key])
			}
		}
	}
}

func TestGetPages(t *testing.T) {
	tests := []struct {
		description   string
		input         string
		expectedPages paginationInfo
		expectedErr   error
	}{
		{
			description: "returns the next page",
			input:       `{"links": {"next": "https://google.com/2"}, "data": []}`,
			expectedPages: paginationInfo{
				Next: "https://google.com/2",
			},
		},
		{
			description: "returns multiple pages",
			input:       `{"links": {"next": "https://google.com/2", "first": "https://google.com/1"}, "data": []}`,
			expectedPages: paginationInfo{
				First: "https://google.com/1",
				Next:  "https://google.com/2",
			},
		},
		{
			description:   "returns no pages",
			input:         `{"data": []}`,
			expectedPages: paginationInfo{},
		},
	}
	for _, testcase := range tests {
		pages, reader, err := getPages(bytes.NewReader([]byte(testcase.input)))
		if pages != testcase.expectedPages {
			t.Errorf("Test %s: Expected pages %+v, got %+v", testcase.description, testcase.expectedPages, pages)
		}

		// we expect to be able to get the original input out again
		resultBytes, _ := ioutil.ReadAll(reader)
		if string(resultBytes) != testcase.input {
			t.Errorf("Test %s: Expected body %s, got %s", testcase.description, testcase.input, string(resultBytes))
		}
		if err != testcase.expectedErr {
			t.Errorf("Test %s: Expected error %v, got %v", testcase.description, testcase.expectedErr, err)
		}
	}
}

func TestUpdateWAFConfigSetInput_validate(t *testing.T) {
	tests := []struct {
		description string
		input       UpdateWAFConfigSetInput
		expected    error
	}{
		{
			description: "Accepts valid input",
			input: UpdateWAFConfigSetInput{
				WAFList:     []ConfigSetWAFs{{ID: "derpID"}},
				ConfigSetID: "derpConfigSet",
			},
			expected: nil,
		},
	}
	for _, testcase := range tests {
		err := testcase.input.validate()
		if err != testcase.expected {
			t.Errorf("In test %s: Expected %v,got %v", testcase.description, testcase.expected, err)
		}
	}
}
