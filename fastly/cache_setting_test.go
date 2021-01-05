package fastly

import (
	"testing"
)

func TestClient_CacheSettings(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "cache_settings/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var cacheSetting *CacheSetting
	record(t, "cache_settings/create", func(c *Client) {
		cacheSetting, err = c.CreateCacheSetting(&CreateCacheSettingInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-cache-setting",
			Action:         CacheSettingActionCache,
			TTL:            1234,
			StaleTTL:       1500,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "cache_settings/cleanup", func(c *Client) {
			c.DeleteCacheSetting(&DeleteCacheSettingInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-cache-setting",
			})

			c.DeleteCacheSetting(&DeleteCacheSettingInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-cache-setting",
			})
		})
	}()

	if cacheSetting.Name != "test-cache-setting" {
		t.Errorf("bad name: %q", cacheSetting.Name)
	}
	if cacheSetting.Action != CacheSettingActionCache {
		t.Errorf("bad action: %q", cacheSetting.Action)
	}
	if cacheSetting.TTL != 1234 {
		t.Errorf("bad ttl: %d", cacheSetting.TTL)
	}
	if cacheSetting.StaleTTL != 1500 {
		t.Errorf("bad stale_ttl: %d", cacheSetting.StaleTTL)
	}

	// List
	var cacheSettings []*CacheSetting
	record(t, "cache_settings/list", func(c *Client) {
		cacheSettings, err = c.ListCacheSettings(&ListCacheSettingsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
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
	record(t, "cache_settings/get", func(c *Client) {
		newCacheSetting, err = c.GetCacheSetting(&GetCacheSettingInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-cache-setting",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if cacheSetting.Name != newCacheSetting.Name {
		t.Errorf("bad name: %q (%q)", cacheSetting.Name, newCacheSetting.Name)
	}
	if cacheSetting.Action != CacheSettingActionCache {
		t.Errorf("bad action: %q", cacheSetting.Action)
	}
	if cacheSetting.TTL != 1234 {
		t.Errorf("bad ttl: %d", cacheSetting.TTL)
	}
	if cacheSetting.StaleTTL != 1500 {
		t.Errorf("bad stale_ttl: %d", cacheSetting.StaleTTL)
	}

	// Update
	var updatedCacheSetting *CacheSetting
	record(t, "cache_settings/update", func(c *Client) {
		updatedCacheSetting, err = c.UpdateCacheSetting(&UpdateCacheSettingInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-cache-setting",
			NewName:        String("new-test-cache-setting"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedCacheSetting.Name != "new-test-cache-setting" {
		t.Errorf("bad name: %q", updatedCacheSetting.Name)
	}

	// Delete
	record(t, "cache_settings/delete", func(c *Client) {
		err = c.DeleteCacheSetting(&DeleteCacheSettingInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-cache-setting",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListCacheSettings_validation(t *testing.T) {
	var err error
	_, err = testClient.ListCacheSettings(&ListCacheSettingsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListCacheSettings(&ListCacheSettingsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateCacheSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateCacheSetting(&CreateCacheSettingInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateCacheSetting(&CreateCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetCacheSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.GetCacheSetting(&GetCacheSettingInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCacheSetting(&GetCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCacheSetting(&GetCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCacheSetting_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCacheSetting(&UpdateCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCacheSetting_validation(t *testing.T) {
	var err error
	err = testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCacheSetting(&DeleteCacheSettingInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
