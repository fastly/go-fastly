package fastly

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
)

func TestClient_WAFs(t *testing.T) {
	t.Parallel()

	fixtureBase := "wafs/"

	testService := createTestService(t, fixtureBase+"service/create", "service2")
	defer deleteTestService(t, fixtureBase+"/service/delete", *testService.ServiceID)

	tv := CreateTestVersion(t, fixtureBase+"/service/version", *testService.ServiceID)

	prefetch := "WAF_Prefetch"
	condition := createTestWAFCondition(t, fixtureBase+"/condition/create", *testService.ServiceID, prefetch, *tv.Number)
	defer deleteTestCondition(t, fixtureBase+"/condition/delete", *testService.ServiceID, prefetch, *tv.Number)

	responseName := "WAF_Response"
	ro := createTestWAFResponseObject(t, fixtureBase+"/response_object/create", *testService.ServiceID, responseName, *tv.Number)
	defer deleteTestResponseObject(t, fixtureBase+"/response_object/delete", *testService.ServiceID, responseName, *tv.Number)

	responseName2 := "WAF_Response2"
	nro := createTestWAFResponseObject(t, fixtureBase+"/response_object/create_another", *testService.ServiceID, responseName2, *tv.Number)
	defer deleteTestResponseObject(t, fixtureBase+"/response_object/cleanup_another", *testService.ServiceID, responseName2, *tv.Number)

	var err error
	var waf *WAF
	Record(t, fixtureBase+"/create", func(c *Client) {
		waf, err = c.CreateWAF(&CreateWAFInput{
			ServiceID:         *testService.ServiceID,
			ServiceVersion:    *tv.Number,
			PrefetchCondition: *condition.Name,
			Response:          *ro.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// List
	var wafsResp *WAFResponse
	Record(t, fixtureBase+"/list", func(c *Client) {
		wafsResp, err = c.ListWAFs(&ListWAFsInput{
			FilterService: *testService.ServiceID,
			FilterVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wafsResp.Items) == 0 {
		t.Errorf("bad wafs: %v", wafsResp.Items)
	}

	// Ensure deleted
	defer func() {
		Record(t, fixtureBase+"/cleanup", func(c *Client) {
			_ = c.DeleteWAF(&DeleteWAFInput{
				ServiceVersion: *tv.Number,
				ID:             waf.ID,
			})
		})
	}()

	Record(t, fixtureBase+"/deploy", func(c *Client) {
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
	Record(t, fixtureBase+"/get", func(c *Client) {
		nwaf, err = c.GetWAF(&GetWAFInput{
			ServiceID:      *testService.ServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, fixtureBase+"/update", func(c *Client) {
		uwaf, err = c.UpdateWAF(&UpdateWAFInput{
			ServiceID:      testService.ServiceID,
			ServiceVersion: tv.Number,
			ID:             waf.ID,
			Response:       nro.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uwaf.Response != responseName2 {
		t.Errorf("bad name: %q", uwaf.Response)
	}

	var dwaf *WAF
	Record(t, fixtureBase+"/disable", func(c *Client) {
		dwaf, err = c.UpdateWAF(&UpdateWAFInput{
			ID:       waf.ID,
			Disabled: ToPointer(true),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !dwaf.Disabled {
		t.Errorf("expected disabled true, got : %v", dwaf.Disabled)
	}

	var ewaf *WAF
	Record(t, fixtureBase+"/enable", func(c *Client) {
		ewaf, err = c.UpdateWAF(&UpdateWAFInput{
			ID:       waf.ID,
			Disabled: ToPointer(false),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ewaf.Disabled {
		t.Errorf("expected disabled false, got : %v", ewaf.Disabled)
	}

	// Delete
	Record(t, fixtureBase+"/delete", func(c *Client) {
		err = c.DeleteWAF(&DeleteWAFInput{
			ServiceVersion: *tv.Number,
			ID:             waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateWAF_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateWAF(&CreateWAFInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateWAF(&CreateWAFInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetWAF_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetWAF(&GetWAFInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetWAF(&GetWAFInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetWAF(&GetWAFInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAF_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateWAF(&UpdateWAFInput{
		ID: "",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateWAF(&UpdateWAFInput{
		ID: "123999",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	serviceID := "foo"

	_, err = TestClient.UpdateWAF(&UpdateWAFInput{
		ID:        "123",
		ServiceID: &serviceID,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteWAF_validation(t *testing.T) {
	var err error
	err = TestClient.DeleteWAF(&DeleteWAFInput{
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteWAF(&DeleteWAFInput{
		ServiceVersion: 1,
		ID:             "",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAF_Enable_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateWAF(&UpdateWAFInput{
		ID:       "",
		Disabled: ToPointer(false),
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAF_Disable_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateWAF(&UpdateWAFInput{
		ID:       "",
		Disabled: ToPointer(true),
	})
	if !errors.Is(err, ErrMissingID) {
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
				jsonapi.QueryParamPageSize:       "2",
				jsonapi.QueryParamPageNumber:     "2",
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
