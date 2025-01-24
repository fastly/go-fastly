package fastly

import (
	"errors"
	"testing"
)

func TestClient_RequestSettings(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "request_settings/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var rs *RequestSetting
	Record(t, "request_settings/create", func(c *Client) {
		rs, err = c.CreateRequestSetting(&CreateRequestSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-request-setting"),
			ForceMiss:      ToPointer(Compatibool(true)),
			ForceSSL:       ToPointer(Compatibool(true)),
			Action:         ToPointer(RequestSettingActionLookup),
			BypassBusyWait: ToPointer(Compatibool(true)),
			MaxStaleAge:    ToPointer(30),
			HashKeys:       ToPointer("a,b,c"),
			XForwardedFor:  ToPointer(RequestSettingXFFLeave),
			TimerSupport:   ToPointer(Compatibool(true)),
			GeoHeaders:     ToPointer(Compatibool(true)),
			DefaultHost:    ToPointer("example.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "request_settings/cleanup", func(c *Client) {
			_ = c.DeleteRequestSetting(&DeleteRequestSettingInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-request-setting",
			})

			_ = c.DeleteRequestSetting(&DeleteRequestSettingInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-request-setting",
			})
		})
	}()

	if *rs.Name != "test-request-setting" {
		t.Errorf("bad name: %q", *rs.Name)
	}
	if !*rs.ForceMiss {
		t.Errorf("bad force_miss: %t", *rs.ForceMiss)
	}
	if !*rs.ForceSSL {
		t.Errorf("bad force_ssl: %t", *rs.ForceSSL)
	}
	if *rs.Action != RequestSettingActionLookup {
		t.Errorf("bad action: %q", *rs.Action)
	}
	if !*rs.BypassBusyWait {
		t.Errorf("bad bypass_busy_wait: %t", *rs.BypassBusyWait)
	}
	if *rs.MaxStaleAge != 30 {
		t.Errorf("bad max_stale_age: %d", *rs.MaxStaleAge)
	}
	if *rs.HashKeys != "a,b,c" {
		t.Errorf("bad has_keys: %q", *rs.HashKeys)
	}
	if *rs.XForwardedFor != RequestSettingXFFLeave {
		t.Errorf("bad xff: %q", *rs.XForwardedFor)
	}
	if !*rs.TimerSupport {
		t.Errorf("bad timer_support: %t", *rs.TimerSupport)
	}
	if !*rs.GeoHeaders {
		t.Errorf("bad geo_headers: %t", *rs.GeoHeaders)
	}
	if *rs.DefaultHost != "example.com" {
		t.Errorf("bad default_host: %q", *rs.DefaultHost)
	}

	// List
	var rss []*RequestSetting
	Record(t, "request_settings/list", func(c *Client) {
		rss, err = c.ListRequestSettings(&ListRequestSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rss) < 1 {
		t.Errorf("bad request settings: %v", rss)
	}

	// Get
	var nrs *RequestSetting
	Record(t, "request_settings/get", func(c *Client) {
		nrs, err = c.GetRequestSetting(&GetRequestSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-request-setting",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *rs.Name != *nrs.Name {
		t.Errorf("bad name: %q (%q)", *rs.Name, *nrs.Name)
	}
	if *rs.ForceMiss != *nrs.ForceMiss {
		t.Errorf("bad force_miss: %t (%t)", *rs.ForceMiss, *nrs.ForceMiss)
	}
	if *rs.ForceSSL != *nrs.ForceSSL {
		t.Errorf("bad force_ssl: %t (%t)", *rs.ForceSSL, *nrs.ForceSSL)
	}
	if *rs.Action != *nrs.Action {
		t.Errorf("bad action: %q (%q)", *rs.Action, *nrs.Action)
	}
	if *rs.BypassBusyWait != *nrs.BypassBusyWait {
		t.Errorf("bad bypass_busy_wait: %t (%t)", *rs.BypassBusyWait, *nrs.BypassBusyWait)
	}
	if *rs.MaxStaleAge != *nrs.MaxStaleAge {
		t.Errorf("bad max_stale_age: %d (%d)", *rs.MaxStaleAge, *nrs.MaxStaleAge)
	}
	if *rs.HashKeys != *nrs.HashKeys {
		t.Errorf("bad has_keys: %q (%q)", *rs.HashKeys, *nrs.HashKeys)
	}
	if *rs.XForwardedFor != *nrs.XForwardedFor {
		t.Errorf("bad xff: %q (%q)", *rs.XForwardedFor, *nrs.XForwardedFor)
	}
	if *rs.TimerSupport != *nrs.TimerSupport {
		t.Errorf("bad timer_support: %t (%t)", *rs.TimerSupport, *nrs.TimerSupport)
	}
	if *rs.GeoHeaders != *nrs.GeoHeaders {
		t.Errorf("bad geo_headers: %t (%t)", *rs.GeoHeaders, *nrs.GeoHeaders)
	}
	if *rs.DefaultHost != *nrs.DefaultHost {
		t.Errorf("bad default_host: %q (%q)", *rs.DefaultHost, *nrs.DefaultHost)
	}

	// Update
	var urs *RequestSetting
	Record(t, "request_settings/update", func(c *Client) {
		urs, err = c.UpdateRequestSetting(&UpdateRequestSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-request-setting",
			NewName:        ToPointer("new-test-request-setting"),
			Action:         ToPointer(RequestSettingActionPass),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *urs.Name != "new-test-request-setting" {
		t.Errorf("bad name: %q", *urs.Name)
	}
	if *urs.Action != RequestSettingActionPass {
		t.Errorf("bad action: %q", *urs.Action)
	}

	// Update 2 (wrap empty string with RequestSettingAction)
	var urs2 *RequestSetting
	Record(t, "request_settings/update-2", func(c *Client) {
		urs2, err = c.UpdateRequestSetting(&UpdateRequestSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-request-setting",
			Action:         ToPointer(RequestSettingAction("")),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *urs2.Action != "" {
		t.Errorf("bad action: %q", *urs2.Action)
	}

	// Update 3 (use explicit RequestSettingActionUnset type)
	var urs3 *RequestSetting
	Record(t, "request_settings/update-3", func(c *Client) {
		urs3, err = c.UpdateRequestSetting(&UpdateRequestSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-request-setting",
			Action:         ToPointer(RequestSettingActionUnset),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *urs3.Action != RequestSettingActionUnset {
		t.Errorf("bad action: %q", *urs3.Action)
	}

	// Delete
	Record(t, "request_settings/delete", func(c *Client) {
		err = c.DeleteRequestSetting(&DeleteRequestSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-request-setting",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListRequestSettings_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListRequestSettings(&ListRequestSettingsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListRequestSettings(&ListRequestSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateRequestSetting_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateRequestSetting(&CreateRequestSettingInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateRequestSetting(&CreateRequestSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetRequestSetting_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetRequestSetting(&GetRequestSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetRequestSetting(&GetRequestSettingInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetRequestSetting(&GetRequestSettingInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateRequestSetting_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateRequestSetting(&UpdateRequestSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateRequestSetting(&UpdateRequestSettingInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateRequestSetting(&UpdateRequestSettingInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteRequestSetting_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteRequestSetting(&DeleteRequestSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteRequestSetting(&DeleteRequestSettingInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteRequestSetting(&DeleteRequestSettingInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
