package fastly

import "testing"

func TestClient_RequestSettings(t *testing.T) {
	tv := testVersion(t)

	// Create
	rs, err := testClient.CreateRequestSetting(&CreateRequestSettingInput{
		Service:        testServiceID,
		Version:        tv.Number,
		Name:           "test-request-setting",
		ForceMiss:      true,
		ForceSSL:       true,
		Action:         RequestSettingActionLookup,
		BypassBusyWait: true,
		MaxStaleAge:    30,
		HashKeys:       "a,b,c",
		XForwardedFor:  RequestSettingXFFLeave,
		TimerSupport:   true,
		GeoHeaders:     true,
		DefaultHost:    "example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteRequestSetting(&DeleteRequestSettingInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-request-setting",
		})

		testClient.DeleteRequestSetting(&DeleteRequestSettingInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-request-setting",
		})
	}()

	if rs.Name != "test-request-setting" {
		t.Errorf("bad name: %q", rs.Name)
	}
	if rs.ForceMiss != true {
		t.Errorf("bad force_miss: %b", rs.ForceMiss)
	}
	if rs.ForceSSL != true {
		t.Errorf("bad force_ssl: %b", rs.ForceSSL)
	}
	if rs.Action != RequestSettingActionLookup {
		t.Errorf("bad action: %q", rs.Action)
	}
	if rs.BypassBusyWait != true {
		t.Errorf("bad bypass_busy_wait: %b", rs.BypassBusyWait)
	}
	if rs.MaxStaleAge != 30 {
		t.Errorf("bad max_stale_age: %d", rs.MaxStaleAge)
	}
	if rs.HashKeys != "a,b,c" {
		t.Errorf("bad has_keys: %q", rs.HashKeys)
	}
	if rs.XForwardedFor != RequestSettingXFFLeave {
		t.Errorf("bad xff: %q", rs.XForwardedFor)
	}
	if rs.TimerSupport != true {
		t.Errorf("bad timer_support: %b", rs.TimerSupport)
	}
	if rs.GeoHeaders != true {
		t.Errorf("bad geo_headers: %b", rs.GeoHeaders)
	}
	if rs.DefaultHost != "example.com" {
		t.Errorf("bad default_host: %q", rs.DefaultHost)
	}

	// List
	rss, err := testClient.ListRequestSettings(&ListRequestSettingsInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rss) < 1 {
		t.Errorf("bad request settings: %v", rss)
	}

	// Get
	nrs, err := testClient.GetRequestSetting(&GetRequestSettingInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-request-setting",
	})
	if err != nil {
		t.Fatal(err)
	}
	if rs.Name != nrs.Name {
		t.Errorf("bad name: %q (%q)", rs.Name, nrs.Name)
	}
	if rs.ForceMiss != nrs.ForceMiss {
		t.Errorf("bad force_miss: %b (%b)", rs.ForceMiss, nrs.ForceMiss)
	}
	if rs.ForceSSL != rs.ForceSSL {
		t.Errorf("bad force_ssl: %b (%b)", rs.ForceSSL, nrs.ForceSSL)
	}
	if rs.Action != nrs.Action {
		t.Errorf("bad action: %q (%q)", rs.Action, nrs.Action)
	}
	if rs.BypassBusyWait != nrs.BypassBusyWait {
		t.Errorf("bad bypass_busy_wait: %b (%b)", rs.BypassBusyWait, nrs.BypassBusyWait)
	}
	if rs.MaxStaleAge != nrs.MaxStaleAge {
		t.Errorf("bad max_stale_age: %d (%d)", rs.MaxStaleAge, nrs.MaxStaleAge)
	}
	if rs.HashKeys != nrs.HashKeys {
		t.Errorf("bad has_keys: %q (%q)", rs.HashKeys, nrs.HashKeys)
	}
	if rs.XForwardedFor != nrs.XForwardedFor {
		t.Errorf("bad xff: %q (%q)", rs.XForwardedFor, nrs.XForwardedFor)
	}
	if rs.TimerSupport != nrs.TimerSupport {
		t.Errorf("bad timer_support: %b (%b)", rs.TimerSupport, nrs.TimerSupport)
	}
	if rs.GeoHeaders != nrs.GeoHeaders {
		t.Errorf("bad geo_headers: %b (%b)", rs.GeoHeaders, nrs.GeoHeaders)
	}
	if rs.DefaultHost != nrs.DefaultHost {
		t.Errorf("bad default_host: %q (%q)", rs.DefaultHost, nrs.DefaultHost)
	}

	// Update
	urs, err := testClient.UpdateRequestSetting(&UpdateRequestSettingInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-request-setting",
		NewName: "new-test-request-setting",
	})
	if err != nil {
		t.Fatal(err)
	}
	if urs.Name != "new-test-request-setting" {
		t.Errorf("bad name: %q", urs.Name)
	}

	// Delete
	if err := testClient.DeleteRequestSetting(&DeleteRequestSettingInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "new-test-request-setting",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListRequestSettings_validation(t *testing.T) {
	var err error
	_, err = testClient.ListRequestSettings(&ListRequestSettingsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListRequestSettings(&ListRequestSettingsInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateRequestSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateRequestSetting(&CreateRequestSettingInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateRequestSetting(&CreateRequestSettingInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetRequestSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.GetRequestSetting(&GetRequestSettingInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetRequestSetting(&GetRequestSettingInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetRequestSetting(&GetRequestSettingInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateRequestSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateRequestSetting(&UpdateRequestSettingInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateRequestSetting(&UpdateRequestSettingInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateRequestSetting(&UpdateRequestSettingInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteRequestSetting_validation(t *testing.T) {
	var err error
	err = testClient.DeleteRequestSetting(&DeleteRequestSettingInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteRequestSetting(&DeleteRequestSettingInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteRequestSetting(&DeleteRequestSettingInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
