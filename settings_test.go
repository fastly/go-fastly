package fastly

import "testing"

func TestClient_Settings(t *testing.T) {
	t.Parallel()

	tv := testVersion(t)

	// Get
	nb, err := testClient.GetSettings(&GetSettingsInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if nb.DefaultTTL == 0 {
		t.Errorf("bad default_ttl: %d", nb.DefaultTTL)
	}

	// Update
	ub, err := testClient.UpdateSettings(&UpdateSettingsInput{
		Service:    testServiceID,
		Version:    tv.Number,
		DefaultTTL: 1800,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ub.DefaultTTL != 1800 {
		t.Errorf("bad default_ttl: %d", ub.DefaultTTL)
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
