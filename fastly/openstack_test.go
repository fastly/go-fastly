package fastly

import (
	"testing"
)

func TestClient_Openstack(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "openstack/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var osCreateResp1, osCreateResp2, osCreateResp3 *Openstack
	record(t, "openstack/create", func(c *Client) {
		osCreateResp1, err = c.CreateOpenstack(&CreateOpenstackInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-openstack",
			User:             "user",
			AccessKey:        "secret-key",
			BucketName:       "bucket-name",
			URL:              "https://logs.example.com/v1.0",
			Path:             "/path",
			Period:           12,
			CompressionCodec: "snappy",
			Format:           "format",
			FormatVersion:    2,
			TimestampFormat:  "%Y",
			MessageType:      "classic",
			Placement:        "waf_debug",
			PublicKey:        pgpPublicKey(),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "openstack/create2", func(c *Client) {
		osCreateResp2, err = c.CreateOpenstack(&CreateOpenstackInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-openstack-2",
			User:             "user",
			AccessKey:        "secret-key",
			BucketName:       "bucket-name",
			URL:              "https://logs.example.com/v1.0",
			Path:             "/path",
			Period:           12,
			CompressionCodec: "snappy",
			GzipLevel:        8,
			Format:           "format",
			FormatVersion:    2,
			TimestampFormat:  "%Y",
			MessageType:      "classic",
			Placement:        "waf_debug",
			PublicKey:        pgpPublicKey(),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "openstack/create3", func(c *Client) {
		osCreateResp3, err = c.CreateOpenstack(&CreateOpenstackInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-openstack-3",
			User:             "user",
			AccessKey:        "secret-key",
			BucketName:       "bucket-name",
			URL:              "https://logs.example.com/v1.0",
			Path:             "/path",
			Period:           12,
			CompressionCodec: "snappy",
			Format:           "format",
			FormatVersion:    2,
			TimestampFormat:  "%Y",
			MessageType:      "classic",
			Placement:        "waf_debug",
			PublicKey:        pgpPublicKey(),
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
				Name:           "test-openstack-2",
			})

			c.DeleteOpenstack(&DeleteOpenstackInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-openstack-3",
			})

			c.DeleteOpenstack(&DeleteOpenstackInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-openstack",
			})
		})
	}()

	if osCreateResp1.Name != "test-openstack" {
		t.Errorf("bad name: %q", osCreateResp1.Name)
	}
	if osCreateResp1.User != "user" {
		t.Errorf("bad user: %q", osCreateResp1.User)
	}
	if osCreateResp1.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", osCreateResp1.BucketName)
	}
	if osCreateResp1.AccessKey != "secret-key" {
		t.Errorf("bad access_key: %q", osCreateResp1.AccessKey)
	}
	if osCreateResp1.Path != "/path" {
		t.Errorf("bad path: %q", osCreateResp1.Path)
	}
	if osCreateResp1.URL != "https://logs.example.com/v1.0" {
		t.Errorf("bad url: %q", osCreateResp1.URL)
	}
	if osCreateResp1.Period != 12 {
		t.Errorf("bad period: %q", osCreateResp1.Period)
	}
	if osCreateResp1.CompressionCodec != "snappy" {
		t.Errorf("bad comprssion_codec: %q", osCreateResp1.CompressionCodec)
	}
	if osCreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", osCreateResp1.GzipLevel)
	}
	if osCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", osCreateResp1.Format)
	}
	if osCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", osCreateResp1.FormatVersion)
	}
	if osCreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", osCreateResp1.TimestampFormat)
	}
	if osCreateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", osCreateResp1.MessageType)
	}
	if osCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", osCreateResp1.Placement)
	}
	if osCreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", osCreateResp1.PublicKey)
	}

	if osCreateResp2.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", osCreateResp2.CompressionCodec)
	}
	if osCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", osCreateResp2.GzipLevel)
	}

	if osCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", osCreateResp3.CompressionCodec)
	}
	if osCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", osCreateResp3.GzipLevel)
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
	var osGetResp *Openstack
	record(t, "openstack/get", func(c *Client) {
		osGetResp, err = c.GetOpenstack(&GetOpenstackInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-openstack",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if osCreateResp1.Name != osGetResp.Name {
		t.Errorf("bad name: %q", osCreateResp1.Name)
	}
	if osCreateResp1.User != osGetResp.User {
		t.Errorf("bad user: %q", osCreateResp1.User)
	}
	if osCreateResp1.BucketName != osGetResp.BucketName {
		t.Errorf("bad bucket_name: %q", osCreateResp1.BucketName)
	}
	if osCreateResp1.AccessKey != osGetResp.AccessKey {
		t.Errorf("bad access_key: %q", osCreateResp1.AccessKey)
	}
	if osCreateResp1.Path != osGetResp.Path {
		t.Errorf("bad path: %q", osCreateResp1.Path)
	}
	if osCreateResp1.URL != osGetResp.URL {
		t.Errorf("bad url: %q", osCreateResp1.URL)
	}
	if osCreateResp1.Period != osGetResp.Period {
		t.Errorf("bad period: %q", osCreateResp1.Period)
	}
	if osCreateResp1.CompressionCodec != osGetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", osCreateResp1.CompressionCodec)
	}
	if osCreateResp1.GzipLevel != osGetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", osCreateResp1.GzipLevel)
	}
	if osCreateResp1.Format != osGetResp.Format {
		t.Errorf("bad format: %q", osCreateResp1.Format)
	}
	if osCreateResp1.FormatVersion != osGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", osCreateResp1.FormatVersion)
	}
	if osCreateResp1.TimestampFormat != osGetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", osCreateResp1.TimestampFormat)
	}
	if osCreateResp1.MessageType != osGetResp.MessageType {
		t.Errorf("bad message_type: %q", osCreateResp1.MessageType)
	}
	if osCreateResp1.Placement != osGetResp.Placement {
		t.Errorf("bad placement: %q", osCreateResp1.Placement)
	}
	if osCreateResp1.PublicKey != osGetResp.PublicKey {
		t.Errorf("bad public_key: %q", osCreateResp1.PublicKey)
	}

	// Update
	var osUpdateResp1, osUpdateResp2, osUpdateResp3 *Openstack
	record(t, "openstack/update", func(c *Client) {
		osUpdateResp1, err = c.UpdateOpenstack(&UpdateOpenstackInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-openstack",
			User:             String("new-user"),
			NewName:          String("new-test-openstack"),
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "openstack/update2", func(c *Client) {
		osUpdateResp2, err = c.UpdateOpenstack(&UpdateOpenstackInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-openstack-2",
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "openstack/update3", func(c *Client) {
		osUpdateResp3, err = c.UpdateOpenstack(&UpdateOpenstackInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-openstack-3",
			GzipLevel:      Uint(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if osUpdateResp1.Name != "new-test-openstack" {
		t.Errorf("bad name: %q", osUpdateResp1.Name)
	}
	if osUpdateResp1.User != "new-user" {
		t.Errorf("bad user: %q", osUpdateResp1.User)
	}
	if osUpdateResp1.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", osUpdateResp1.CompressionCodec)
	}
	if osUpdateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", osUpdateResp1.GzipLevel)
	}

	if osUpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", osUpdateResp2.CompressionCodec)
	}
	if osUpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", osUpdateResp2.GzipLevel)
	}

	if osUpdateResp3.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", osUpdateResp3.CompressionCodec)
	}
	if osUpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", osUpdateResp3.GzipLevel)
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
