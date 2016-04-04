package fastly

import "testing"

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
			Service:    testServiceID,
			Version:    tv.Number,
			DefaultTTL: 1800,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.DefaultTTL != 1800 {
		t.Errorf("bad default_ttl: %d", us.DefaultTTL)
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
		Version: "",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}
