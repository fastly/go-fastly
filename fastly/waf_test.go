package fastly

import (
	"reflect"
	"testing"
)

func TestClient_WAFs(t *testing.T) {
	t.Parallel()

	fixtureBase := "wafs/"

	testService := createTestService(t, fixtureBase+"service/create", "service2")
	defer deleteTestService(t, fixtureBase+"/service/delete", testService.ID)

	tv := createTestVersion(t, fixtureBase+"/service/version", testService.ID)

	prefetch := "WAF_Prefetch"
	condition := createTestWAFCondition(t, fixtureBase+"/condition/create", testService.ID, prefetch, tv.Number)
	defer deleteTestCondition(t, fixtureBase+"/condition/delete", testService.ID, prefetch, tv.Number)

	responseName := "WAF_Response"
	ro := createTestWAFResponseObject(t, fixtureBase+"/response_object/create", testService.ID, responseName, tv.Number)
	defer deleteTestResponseObject(t, fixtureBase+"/response_object/delete", testService.ID, responseName, tv.Number)

	responseName2 := "WAF_Response2"
	nro := createTestWAFResponseObject(t, fixtureBase+"/response_object/create_another", testService.ID, responseName2, tv.Number)
	defer deleteTestResponseObject(t, fixtureBase+"/response_object/cleanup_another", testService.ID, responseName2, tv.Number)

	var err error
	var waf *WAF
	record(t, fixtureBase+"/create", func(c *Client) {
		waf, err = c.CreateWAF(&CreateWAFInput{
			ServiceID:         testService.ID,
			ServiceVersion:    tv.Number,
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
				ServiceVersion: tv.Number,
				ID:             waf.ID,
			})
		})
	}()

	record(t, fixtureBase+"/deploy", func(c *Client) {
		err = c.DeployWAFVersion(&DeployWAFVersionInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get
	var nwaf *WAF
	record(t, fixtureBase+"/get", func(c *Client) {
		nwaf, err = c.GetWAF(&GetWAFInput{
			ServiceID:      testService.ID,
			ServiceVersion: tv.Number,
			ID:             waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if nwaf.ID != waf.ID {
		t.Errorf("expected %q to be %q", nwaf.ID, waf.ID)
	}
	if nwaf.Disabled {
		t.Errorf("expected disabled false, got : %v", nwaf.Disabled)
	}

	var uwaf *WAF
	record(t, fixtureBase+"/update", func(c *Client) {
		uwaf, err = c.UpdateWAF(&UpdateWAFInput{
			ServiceID:      &testService.ID,
			ServiceVersion: &tv.Number,
			ID:             waf.ID,
			Response:       &nro.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uwaf.Response != responseName2 {
		t.Errorf("bad name: %q", uwaf.Response)
	}

	var dwaf *WAF
	record(t, fixtureBase+"/disable", func(c *Client) {
		dwaf, err = c.UpdateWAF(&UpdateWAFInput{
			ID:       waf.ID,
			Disabled: Bool(true),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !dwaf.Disabled {
		t.Errorf("expected disabled true, got : %v", dwaf.Disabled)
	}

	var ewaf *WAF
	record(t, fixtureBase+"/enable", func(c *Client) {
		ewaf, err = c.UpdateWAF(&UpdateWAFInput{
			ID:       waf.ID,
			Disabled: Bool(false),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ewaf.Disabled {
		t.Errorf("expected disabled false, got : %v", ewaf.Disabled)
	}

	// Delete
	record(t, fixtureBase+"/delete", func(c *Client) {
		err = c.DeleteWAF(&DeleteWAFInput{
			ServiceVersion: tv.Number,
			ID:             waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateWAF_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateWAF(&CreateWAFInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateWAF(&CreateWAFInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetWAF_validation(t *testing.T) {
	var err error
	_, err = testClient.GetWAF(&GetWAFInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetWAF(&GetWAFInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetWAF(&GetWAFInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAF_validation(t *testing.T) {
	var err error

	_, err = testClient.UpdateWAF(&UpdateWAFInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateWAF(&UpdateWAFInput{
		ID: "123999",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	serviceID := "foo"

	_, err = testClient.UpdateWAF(&UpdateWAFInput{
		ID:        "123",
		ServiceID: &serviceID,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteWAF_validation(t *testing.T) {
	var err error
	err = testClient.DeleteWAF(&DeleteWAFInput{
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteWAF(&DeleteWAFInput{
		ServiceVersion: 1,
		ID:             "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAF_Enable_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateWAF(&UpdateWAFInput{
		ID:       "",
		Disabled: Bool(false),
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAF_Disable_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateWAF(&UpdateWAFInput{
		ID:       "",
		Disabled: Bool(true),
	})
	if err != ErrMissingID {
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
