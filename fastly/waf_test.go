package fastly

import (
	"testing"
	"fmt"
	"strconv"
)

func TestClient_WAFs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "wafs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Enable logging on the service - we cannot create wafs without logging
	// enabled
	var s *Syslog
	record(t, "wafs/logging/create", func(c *Client) {
		s, err = c.CreateSyslog(&CreateSyslogInput{
			Service:       testServiceID,
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
		record(t, "wafs/logging/cleanup", func(c *Client) {
			c.DeleteSyslog(&DeleteSyslogInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-syslog",
			})
		})
	}()

	// Create a condition - we cannot create a waf without attaching a condition
	var condition *Condition
	record(t, "wafs/condition/create", func(c *Client) {
		condition, err = c.CreateCondition(&CreateConditionInput{
			Service:   testServiceID,
			Version:   tv.Number,
			Name:      "test-waf-condition",
			Statement: "req.url~+\"index.html\"",
			Type:      "PREFETCH", // This must be a prefetch condition
			Priority:  1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		record(t, "wafs/condition/cleanup", func(c *Client) {
			c.DeleteCondition(&DeleteConditionInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-waf-condition",
			})
		})
	}()

	// Create a response object
	var ro *ResponseObject
	record(t, "wafs/response_object/create", func(c *Client) {
		ro, err = c.CreateResponseObject(&CreateResponseObjectInput{
			Service:     testServiceID,
			Version:     tv.Number,
			Name:        "test-response-object",
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
		record(t, "wafs/response_object/cleanup", func(c *Client) {
			c.DeleteResponseObject(&DeleteResponseObjectInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    ro.Name,
			})
		})
	}()

	// Create
	var waf *WAF
	record(t, "wafs/create", func(c *Client) {
		waf, err = c.CreateWAF(&CreateWAFInput{
			Service:           testServiceID,
			Version:           tv.Number,
			PrefetchCondition: condition.Name,
			Response:          ro.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// List
	var wafs []*WAF
	record(t, "wafs/list", func(c *Client) {
		wafs, err = c.ListWAFs(&ListWAFsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wafs) < 1 {
		t.Errorf("bad wafs: %v", wafs)
	}

	// Ensure deleted
	defer func() {
		record(t, "wafs/cleanup", func(c *Client) {
			c.DeleteWAF(&DeleteWAFInput{
				Service: testServiceID,
				Version: tv.Number,
				ID:      waf.ID,
			})
		})
	}()

	// Get
	var nwaf *WAF
	record(t, "wafs/get", func(c *Client) {
		nwaf, err = c.GetWAF(&GetWAFInput{
			Service: testServiceID,
			Version: tv.Number,
			ID:      waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if nwaf.ID != waf.ID {
		t.Errorf("expected %q to be %q", nwaf.ID, waf.ID)
	}

	// Update
	// Create a new response object to attach
	var nro *ResponseObject
	record(t, "wafs/response_object/create_another", func(c *Client) {
		nro, err = c.CreateResponseObject(&CreateResponseObjectInput{
			Service:     testServiceID,
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
		record(t, "wafs/response_object/cleanup_another", func(c *Client) {
			c.DeleteResponseObject(&DeleteResponseObjectInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    nro.Name,
			})
		})
	}()
	var uwaf *WAF
	record(t, "wafs/update", func(c *Client) {
		uwaf, err = c.UpdateWAF(&UpdateWAFInput{
			Service:  testServiceID,
			Version:  tv.Number,
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

	// Delete
	record(t, "wafs/delete", func(c *Client) {
		err = c.DeleteWAF(&DeleteWAFInput{
			Service: testServiceID,
			Version: tv.Number,
			ID:      waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	var rsro *RuleStatus
	record(t, "wafs/rule_status/get", func(c *Client) {
		rsro, err = c.GetWAFRuleStatus(&GetWAFRuleStatusInput{
			Service: testServiceID,
			RuleID: 933120,
			ID: waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedID := fmt.Sprintf("%v-%v", waf.ID, strconv.Itoa(933120))

	if rsro.ID !=  expectedID {
		t.Errorf("Get RuleStatus failed %s\n", rsro.ID)
	}

	var urs *RuleStatus
	record(t, "wafs/rule_status/update", func(c *Client) {
		urs, err = c.UpdateWAFRuleStatus(&UpdateWAFRuleStatusInput{
			Service: testServiceID,
			RuleID: 933120,
			ID: waf.ID,
		})
	})

}

func TestClient_ListWAFs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListWAFs(&ListWAFsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListWAFs(&ListWAFsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
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
		Version: 0,
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
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetWAF(&GetWAFInput{
		Service: "foo",
		Version: 1,
		ID:      "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}
}




//
// func TestClient_UpdateWAF_validation(t *testing.T) {
// 	var err error
// 	_, err = testClient.UpdateWAF(&UpdateWAFInput{
// 		Service: "",
// 	})
// 	if err != ErrMissingService {
// 		t.Errorf("bad error: %s", err)
// 	}
//
// 	_, err = testClient.UpdateWAF(&UpdateWAFInput{
// 		Service: "foo",
// 		Version: 0,
// 	})
// 	if err != ErrMissingVersion {
// 		t.Errorf("bad error: %s", err)
// 	}
//
// 	_, err = testClient.UpdateWAF(&UpdateWAFInput{
// 		Service: "foo",
// 		Version: 1,
// 		WAFID:   "",
// 	})
// 	if err != ErrMissingWAFID {
// 		t.Errorf("bad error: %s", err)
// 	}
// }
//
// func TestClient_DeleteWAF_validation(t *testing.T) {
// 	var err error
// 	err = testClient.DeleteWAF(&DeleteWAFInput{
// 		Service: "",
// 	})
// 	if err != ErrMissingService {
// 		t.Errorf("bad error: %s", err)
// 	}
//
// 	err = testClient.DeleteWAF(&DeleteWAFInput{
// 		Service: "foo",
// 		Version: 0,
// 	})
// 	if err != ErrMissingVersion {
// 		t.Errorf("bad error: %s", err)
// 	}
//
// 	err = testClient.DeleteWAF(&DeleteWAFInput{
// 		Service: "foo",
// 		Version: 1,
// 		WAFID:   "",
// 	})
// 	if err != ErrMissingWAFID {
// 		t.Errorf("bad error: %s", err)
// 	}
// }
