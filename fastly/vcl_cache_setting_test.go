package fastly

import (
	"errors"
	"testing"
)

func TestClient_CacheSettings(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "cache_settings/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var cacheSetting *CacheSetting
	Record(t, "cache_settings/create", func(c *Client) {
		cacheSetting, err = c.CreateCacheSetting(&CreateCacheSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-cache-setting"),
			Action:         ToPointer(CacheSettingActionCache),
			TTL:            ToPointer(1234),
			StaleTTL:       ToPointer(1500),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "cache_settings/cleanup", func(c *Client) {
			_ = c.DeleteCacheSetting(&DeleteCacheSettingInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-cache-setting",
			})

			_ = c.DeleteCacheSetting(&DeleteCacheSettingInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-cache-setting",
			})
		})
	}()

	if *cacheSetting.Name != "test-cache-setting" {
		t.Errorf("bad name: %q", *cacheSetting.Name)
	}
	if *cacheSetting.Action != CacheSettingActionCache {
		t.Errorf("bad action: %q", *cacheSetting.Action)
	}
	if *cacheSetting.TTL != 1234 {
		t.Errorf("bad ttl: %d", *cacheSetting.TTL)
	}
	if *cacheSetting.StaleTTL != 1500 {
		t.Errorf("bad stale_ttl: %d", *cacheSetting.StaleTTL)
	}

	// List
	var cacheSettings []*CacheSetting
	Record(t, "cache_settings/list", func(c *Client) {
		cacheSettings, err = c.ListCacheSettings(&ListCacheSettingsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cacheSettings) < 1 {
		t.Errorf("bad cache settings: %v", cacheSettings)
	}

	// Get
	var newCacheSetting *CacheSetting
	Record(t, "cache_settings/get", func(c *Client) {
		newCacheSetting, err = c.GetCacheSetting(&GetCacheSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-cache-setting",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *cacheSetting.Name != *newCacheSetting.Name {
		t.Errorf("bad name: %q (%q)", *cacheSetting.Name, *newCacheSetting.Name)
	}
	if *cacheSetting.Action != CacheSettingActionCache {
		t.Errorf("bad action: %q", *cacheSetting.Action)
	}
	if *cacheSetting.TTL != 1234 {
		t.Errorf("bad ttl: %d", *cacheSetting.TTL)
	}
	if *cacheSetting.StaleTTL != 1500 {
		t.Errorf("bad stale_ttl: %d", *cacheSetting.StaleTTL)
	}

	// Update
	var updatedCacheSetting *CacheSetting
	Record(t, "cache_settings/update", func(c *Client) {
		updatedCacheSetting, err = c.UpdateCacheSetting(&UpdateCacheSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-cache-setting",
			NewName:        ToPointer("new-test-cache-setting"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *updatedCacheSetting.Name != "new-test-cache-setting" {
		t.Errorf("bad name: %q", *updatedCacheSetting.Name)
	}

	// Delete
	Record(t, "cache_settings/delete", func(c *Client) {
		err = c.DeleteCacheSetting(&DeleteCacheSettingInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-cache-setting",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListCacheSettings_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListCacheSettings(&ListCacheSettingsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListCacheSettings(&ListCacheSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateCacheSetting_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateCacheSetting(&CreateCacheSettingInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateCacheSetting(&CreateCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetCacheSetting_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetCacheSetting(&GetCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetCacheSetting(&GetCacheSettingInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetCacheSetting(&GetCacheSettingInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCacheSetting_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCacheSetting_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
