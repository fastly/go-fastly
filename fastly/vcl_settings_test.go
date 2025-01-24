package fastly

import (
	"errors"
	"testing"

	"github.com/google/go-querystring/query"
)

func TestClient_Settings(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "settings/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Get
	var ns *Settings
	Record(t, "settings/get", func(c *Client) {
		ns, err = c.GetSettings(&GetSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ns.DefaultTTL == 0 {
		t.Errorf("bad default_ttl: %d", ns.DefaultTTL)
	}

	// Update
	var us *Settings
	Record(t, "settings/update", func(c *Client) {
		us, err = c.UpdateSettings(&UpdateSettingsInput{
			ServiceID:       TestDeliveryServiceID,
			ServiceVersion:  *tv.Number,
			DefaultTTL:      ToPointer(uint(1800)),
			StaleIfError:    ToPointer(true),
			StaleIfErrorTTL: ToPointer(uint(57600)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.DefaultTTL != 1800 {
		t.Errorf("bad default_ttl: %d", *us.DefaultTTL)
	}
	if !*us.StaleIfError {
		t.Errorf("bad stale_if_error: %t", *us.StaleIfError)
	}
	if *us.StaleIfErrorTTL != 57600 {
		t.Errorf("bad stale_if_error_ttl %d", *us.StaleIfErrorTTL)
	}
}

// Tests if we can update a default_ttl to 0 as reported in issue #20.
func TestClient_UpdateSettingsInput_default_ttl(t *testing.T) {
	t.Parallel()
	s := UpdateSettingsInput{
		DefaultTTL:     ToPointer(uint(0)),
		ServiceID:      "foo",
		ServiceVersion: 1,
	}

	v, err := query.Values(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	body := v.Encode()

	if body != "ServiceID=foo&ServiceVersion=1&general.default_ttl=0" {
		t.Errorf("Update request should contain a default_ttl. Got: %s", body)
	}
}

func TestClient_GetSettings_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetSettings(&GetSettingsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSettings(&GetSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSettings_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateSettings(&UpdateSettingsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSettings(&UpdateSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
