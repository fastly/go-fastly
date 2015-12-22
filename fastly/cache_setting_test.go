package fastly

import "testing"

func TestClient_CacheSettings(t *testing.T) {
	tv := testVersion(t)

	// Create
	c, err := testClient.CreateCacheSetting(&CreateCacheSettingInput{
		Service:  testServiceID,
		Version:  tv.Number,
		Name:     "test-cache-setting",
		Action:   CacheSettingActionCache,
		TTL:      1234,
		StaleTTL: 1500,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-cache-setting",
		})

		testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-cache-setting",
		})
	}()

	if c.Name != "test-cache-setting" {
		t.Errorf("bad name: %q", c.Name)
	}
	if c.Action != CacheSettingActionCache {
		t.Errorf("bad action: %q", c.Action)
	}
	if c.TTL != 1234 {
		t.Errorf("bad ttl: %d", c.TTL)
	}
	if c.StaleTTL != 1500 {
		t.Errorf("bad stale_ttl: %d", c.StaleTTL)
	}

	// List
	cs, err := testClient.ListCacheSettings(&ListCacheSettingsInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cs) < 1 {
		t.Errorf("bad cache settings: %v", cs)
	}

	// Get
	nc, err := testClient.GetCacheSetting(&GetCacheSettingInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-cache-setting",
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != nc.Name {
		t.Errorf("bad name: %q (%q)", c.Name, nc.Name)
	}
	if c.Action != CacheSettingActionCache {
		t.Errorf("bad action: %q", c.Action)
	}
	if c.TTL != 1234 {
		t.Errorf("bad ttl: %d", c.TTL)
	}
	if c.StaleTTL != 1500 {
		t.Errorf("bad stale_ttl: %d", c.StaleTTL)
	}

	// Update
	uc, err := testClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-cache-setting",
		NewName: "new-test-cache-setting",
	})
	if err != nil {
		t.Fatal(err)
	}
	if uc.Name != "new-test-cache-setting" {
		t.Errorf("bad name: %q", uc.Name)
	}

	// Delete
	if err := testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "new-test-cache-setting",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListCacheSettings_validation(t *testing.T) {
	var err error
	_, err = testClient.ListCacheSettings(&ListCacheSettingsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListCacheSettings(&ListCacheSettingsInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateCacheSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateCacheSetting(&CreateCacheSettingInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateCacheSetting(&CreateCacheSettingInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetCacheSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.GetCacheSetting(&GetCacheSettingInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCacheSetting(&GetCacheSettingInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCacheSetting(&GetCacheSettingInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCacheSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCacheSetting_validation(t *testing.T) {
	var err error
	err = testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
