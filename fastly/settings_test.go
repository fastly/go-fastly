package fastly

import (
	"bytes"
	"testing"

	"github.com/ajg/form"
)

func TestClient_Settings(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "settings/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Get
	var ns *Settings
	record(t, "settings/get", func(c *Client) {
		ns, err = c.GetSettings(&GetSettingsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ns.DefaultTTL == 0 {
		t.Errorf("bad default_ttl: %d", ns.DefaultTTL)
	}

	// Update
	var us *Settings
	record(t, "settings/update", func(c *Client) {
		us, err = c.UpdateSettings(&UpdateSettingsInput{
			Service:         testServiceID,
			Version:         tv.Number,
			DefaultTTL:      1800,
			StaleIfError:    true,
			StaleIfErrorTTL: 57600,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.DefaultTTL != 1800 {
		t.Errorf("bad default_ttl: %d", us.DefaultTTL)
	}
	if us.StaleIfError != true {
		t.Errorf("bad stale_if_error: %t", us.StaleIfError)
	}
	if us.StaleIfErrorTTL != 57600 {
		t.Errorf("bad stale_if_error_ttl %d", us.StaleIfErrorTTL)
	}
}

// Tests if we can update a default_ttl to 0 as reported in issue #20
func TestClient_UpdateSettingsInput_default_ttl(t *testing.T) {
	t.Parallel()
	s := UpdateSettingsInput{Service: "foo", Version: 1, DefaultTTL: 0}
	buf := new(bytes.Buffer)
	form.NewEncoder(buf).KeepZeros(true).DelimitWith('|').Encode(s)
	if buf.String() != "Service=foo&Version=1&general.default_ttl=0" {
		t.Errorf("Update request should contain a default_ttl. Got: %s", buf.String())
	}
}

func TestClient_GetSettings_validation(t *testing.T) {
	var err error
	_, err = testClient.GetSettings(&GetSettingsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSettings(&GetSettingsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSettings_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateSettings(&UpdateSettingsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSettings(&UpdateSettingsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}
