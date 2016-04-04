package fastly

import "testing"

func TestClient_GCSs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "gcses/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var gcs *GCS
	record(t, "gcses/create", func(c *Client) {
		gcs, err = c.CreateGCS(&CreateGCSInput{
			Service:         testServiceID,
			Version:         tv.Number,
			Name:            "test-gcs",
			Bucket:          "bucket",
			User:            "user",
			SecretKey:       "key",
			Path:            "/path",
			Period:          12,
			GzipLevel:       9,
			Format:          "format",
			TimestampFormat: "%Y",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "gcses/cleanup", func(c *Client) {
			c.DeleteGCS(&DeleteGCSInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-gcs",
			})

			c.DeleteGCS(&DeleteGCSInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-gcs",
			})
		})
	}()

	if gcs.Name != "test-gcs" {
		t.Errorf("bad name: %q", gcs.Name)
	}
	if gcs.Bucket != "bucket" {
		t.Errorf("bad bucket: %q", gcs.Bucket)
	}
	if gcs.User != "user" {
		t.Errorf("bad user: %q", gcs.User)
	}
	if gcs.SecretKey != "key" {
		t.Errorf("bad secret_key: %q", gcs.SecretKey)
	}
	if gcs.Path != "/path" {
		t.Errorf("bad path: %q", gcs.Path)
	}
	if gcs.Period != 12 {
		t.Errorf("bad period: %q", gcs.Period)
	}
	if gcs.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", gcs.GzipLevel)
	}
	if gcs.Format != "format" {
		t.Errorf("bad format: %q", gcs.Format)
	}
	if gcs.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", gcs.TimestampFormat)
	}

	// List
	var gcses []*GCS
	record(t, "gcses/list", func(c *Client) {
		gcses, err = c.ListGCSs(&ListGCSsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(gcses) < 1 {
		t.Errorf("bad gcses: %v", gcses)
	}

	// Get
	var ngcs *GCS
	record(t, "gcses/get", func(c *Client) {
		ngcs, err = c.GetGCS(&GetGCSInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-gcs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gcs.Name != "test-gcs" {
		t.Errorf("bad name: %q", gcs.Name)
	}
	if gcs.Bucket != "bucket" {
		t.Errorf("bad bucket: %q", gcs.Bucket)
	}
	if gcs.User != "user" {
		t.Errorf("bad user: %q", gcs.User)
	}
	if gcs.SecretKey != "key" {
		t.Errorf("bad secret_key: %q", gcs.SecretKey)
	}
	if gcs.Path != "/path" {
		t.Errorf("bad path: %q", gcs.Path)
	}
	if gcs.Period != 12 {
		t.Errorf("bad period: %q", gcs.Period)
	}
	if gcs.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", gcs.GzipLevel)
	}
	if gcs.Format != "format" {
		t.Errorf("bad format: %q", gcs.Format)
	}
	if gcs.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", gcs.TimestampFormat)
	}

	// Update
	var ugcs *GCS
	record(t, "gcses/update", func(c *Client) {
		ugcs, err = c.UpdateGCS(&UpdateGCSInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-gcs",
			NewName: "new-test-gcs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ugcs.Name != "new-test-gcs" {
		t.Errorf("bad name: %q", ugcs.Name)
	}

	// Delete
	record(t, "gcses/delete", func(c *Client) {
		err = c.DeleteGCS(&DeleteGCSInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-gcs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListGCSs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListGCSs(&ListGCSsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListGCSs(&ListGCSsInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateGCS(&CreateGCSInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateGCS(&CreateGCSInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.GetGCS(&GetGCSInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetGCS(&GetGCSInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetGCS(&GetGCSInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteGCS_validation(t *testing.T) {
	var err error
	err = testClient.DeleteGCS(&DeleteGCSInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteGCS(&DeleteGCSInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteGCS(&DeleteGCSInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
