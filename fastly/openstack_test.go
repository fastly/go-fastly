package fastly

import "testing"

func TestClient_Openstack(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "openstack/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var openstack *Openstack
	record(t, "openstack/create", func(c *Client) {
		openstack, err = c.CreateOpenstack(&CreateOpenstackInput{
			ServiceID:       testServiceID,
			ServiceVersion:  tv.Number,
			Name:            String("test-openstack"),
			User:            String("user"),
			AccessKey:       String("secret-key"),
			BucketName:      String("bucket-name"),
			URL:             String("https://logs.example.com/v1.0"),
			Path:            String("/path"),
			Period:          Uint(12),
			GzipLevel:       Uint(9),
			Format:          String("format"),
			FormatVersion:   Uint(2),
			TimestampFormat: String("%Y"),
			MessageType:     String("classic"),
			Placement:       String("waf_debug"),
			PublicKey:       String(pgpPublicKey()),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "openstack/cleanup", func(c *Client) {
			c.DeleteOpenstack(&DeleteOpenstackInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-openstack",
			})

			c.DeleteOpenstack(&DeleteOpenstackInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-openstack",
			})
		})
	}()

	if openstack.Name != "test-openstack" {
		t.Errorf("bad name: %q", openstack.Name)
	}
	if openstack.User != "user" {
		t.Errorf("bad user: %q", openstack.User)
	}
	if openstack.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", openstack.BucketName)
	}
	if openstack.AccessKey != "secret-key" {
		t.Errorf("bad access_key: %q", openstack.AccessKey)
	}
	if openstack.Path != "/path" {
		t.Errorf("bad path: %q", openstack.Path)
	}
	if openstack.URL != "https://logs.example.com/v1.0" {
		t.Errorf("bad url: %q", openstack.URL)
	}
	if openstack.Period != 12 {
		t.Errorf("bad period: %q", openstack.Period)
	}
	if openstack.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", openstack.GzipLevel)
	}
	if openstack.Format != "format" {
		t.Errorf("bad format: %q", openstack.Format)
	}
	if openstack.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", openstack.FormatVersion)
	}
	if openstack.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", openstack.TimestampFormat)
	}
	if openstack.MessageType != "classic" {
		t.Errorf("bad message_type: %q", openstack.MessageType)
	}
	if openstack.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", openstack.Placement)
	}
	if openstack.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", openstack.PublicKey)
	}

	// List
	var lc []*Openstack
	record(t, "openstack/list", func(c *Client) {
		lc, err = c.ListOpenstack(&ListOpenstackInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lc) < 1 {
		t.Errorf("bad openstack: %v", lc)
	}

	// Get
	var nopenstack *Openstack
	record(t, "openstack/get", func(c *Client) {
		nopenstack, err = c.GetOpenstack(&GetOpenstackInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-openstack",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if openstack.Name != nopenstack.Name {
		t.Errorf("bad name: %q", openstack.Name)
	}
	if openstack.User != nopenstack.User {
		t.Errorf("bad user: %q", openstack.User)
	}
	if openstack.BucketName != nopenstack.BucketName {
		t.Errorf("bad bucket_name: %q", openstack.BucketName)
	}
	if openstack.AccessKey != nopenstack.AccessKey {
		t.Errorf("bad access_key: %q", openstack.AccessKey)
	}
	if openstack.Path != nopenstack.Path {
		t.Errorf("bad path: %q", openstack.Path)
	}
	if openstack.URL != nopenstack.URL {
		t.Errorf("bad url: %q", openstack.URL)
	}
	if openstack.Period != nopenstack.Period {
		t.Errorf("bad period: %q", openstack.Period)
	}
	if openstack.GzipLevel != nopenstack.GzipLevel {
		t.Errorf("bad gzip_level: %q", openstack.GzipLevel)
	}
	if openstack.Format != nopenstack.Format {
		t.Errorf("bad format: %q", openstack.Format)
	}
	if openstack.FormatVersion != nopenstack.FormatVersion {
		t.Errorf("bad format_version: %q", openstack.FormatVersion)
	}
	if openstack.TimestampFormat != nopenstack.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", openstack.TimestampFormat)
	}
	if openstack.MessageType != nopenstack.MessageType {
		t.Errorf("bad message_type: %q", openstack.MessageType)
	}
	if openstack.Placement != nopenstack.Placement {
		t.Errorf("bad placement: %q", openstack.Placement)
	}
	if openstack.PublicKey != nopenstack.PublicKey {
		t.Errorf("bad public_key: %q", openstack.PublicKey)
	}

	// Update
	var uopenstack *Openstack
	record(t, "openstack/update", func(c *Client) {
		uopenstack, err = c.UpdateOpenstack(&UpdateOpenstackInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-openstack",
			User:           String("new-user"),
			NewName:        String("new-test-openstack"),
			GzipLevel:      Uint(0),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uopenstack.Name != "new-test-openstack" {
		t.Errorf("bad name: %q", uopenstack.Name)
	}
	if uopenstack.User != "new-user" {
		t.Errorf("bad user: %q", uopenstack.User)
	}
	if uopenstack.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", uopenstack.GzipLevel)
	}

	// Delete
	record(t, "openstack/delete", func(c *Client) {
		err = c.DeleteOpenstack(&DeleteOpenstackInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-openstack",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListOpenstack_validation(t *testing.T) {
	var err error
	_, err = testClient.ListOpenstack(&ListOpenstackInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListOpenstack(&ListOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateOpenstack_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateOpenstack(&CreateOpenstackInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateOpenstack(&CreateOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetOpenstack_validation(t *testing.T) {
	var err error
	_, err = testClient.GetOpenstack(&GetOpenstackInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetOpenstack(&GetOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetOpenstack(&GetOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateOpenstack_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateOpenstack(&UpdateOpenstackInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateOpenstack(&UpdateOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateOpenstack(&UpdateOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteOpenstack_validation(t *testing.T) {
	var err error
	err = testClient.DeleteOpenstack(&DeleteOpenstackInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteOpenstack(&DeleteOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteOpenstack(&DeleteOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
