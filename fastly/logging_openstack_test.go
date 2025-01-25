package fastly

import (
	"errors"
	"testing"
)

func TestClient_Openstack(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "openstack/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var osCreateResp1, osCreateResp2, osCreateResp3 *Openstack
	Record(t, "openstack/create", func(c *Client) {
		osCreateResp1, err = c.CreateOpenstack(&CreateOpenstackInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             ToPointer("test-openstack"),
			User:             ToPointer("user"),
			AccessKey:        ToPointer("secret-key"),
			BucketName:       ToPointer("bucket-name"),
			URL:              ToPointer("https://logs.example.com/v1.0"),
			Path:             ToPointer("/path"),
			Period:           ToPointer(12),
			CompressionCodec: ToPointer("snappy"),
			Format:           ToPointer("format"),
			FormatVersion:    ToPointer(2),
			TimestampFormat:  ToPointer("%Y"),
			MessageType:      ToPointer("classic"),
			Placement:        ToPointer("waf_debug"),
			PublicKey:        ToPointer(pgpPublicKey()),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "openstack/create2", func(c *Client) {
		osCreateResp2, err = c.CreateOpenstack(&CreateOpenstackInput{
			ServiceID:       TestDeliveryServiceID,
			ServiceVersion:  *tv.Number,
			Name:            ToPointer("test-openstack-2"),
			User:            ToPointer("user"),
			AccessKey:       ToPointer("secret-key"),
			BucketName:      ToPointer("bucket-name"),
			URL:             ToPointer("https://logs.example.com/v1.0"),
			Path:            ToPointer("/path"),
			Period:          ToPointer(12),
			GzipLevel:       ToPointer(8),
			Format:          ToPointer("format"),
			FormatVersion:   ToPointer(2),
			TimestampFormat: ToPointer("%Y"),
			MessageType:     ToPointer("classic"),
			Placement:       ToPointer("waf_debug"),
			PublicKey:       ToPointer(pgpPublicKey()),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "openstack/create3", func(c *Client) {
		osCreateResp3, err = c.CreateOpenstack(&CreateOpenstackInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             ToPointer("test-openstack-3"),
			User:             ToPointer("user"),
			AccessKey:        ToPointer("secret-key"),
			BucketName:       ToPointer("bucket-name"),
			URL:              ToPointer("https://logs.example.com/v1.0"),
			Path:             ToPointer("/path"),
			Period:           ToPointer(12),
			CompressionCodec: ToPointer("snappy"),
			Format:           ToPointer("format"),
			FormatVersion:    ToPointer(2),
			TimestampFormat:  ToPointer("%Y"),
			MessageType:      ToPointer("classic"),
			Placement:        ToPointer("waf_debug"),
			PublicKey:        ToPointer(pgpPublicKey()),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail because both CompressionCodec and
	// GzipLevel are present.
	Record(t, "openstack/create4", func(c *Client) {
		_, err = c.CreateOpenstack(&CreateOpenstackInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             ToPointer("test-openstack-4"),
			User:             ToPointer("user"),
			AccessKey:        ToPointer("secret-key"),
			BucketName:       ToPointer("bucket-name"),
			URL:              ToPointer("https://logs.example.com/v1.0"),
			Path:             ToPointer("/path"),
			Period:           ToPointer(12),
			CompressionCodec: ToPointer("snappy"),
			GzipLevel:        ToPointer(8),
			Format:           ToPointer("format"),
			FormatVersion:    ToPointer(2),
			TimestampFormat:  ToPointer("%Y"),
			MessageType:      ToPointer("classic"),
			Placement:        ToPointer("waf_debug"),
			PublicKey:        ToPointer(pgpPublicKey()),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "openstack/cleanup", func(c *Client) {
			_ = c.DeleteOpenstack(&DeleteOpenstackInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-openstack",
			})

			_ = c.DeleteOpenstack(&DeleteOpenstackInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-openstack-2",
			})

			_ = c.DeleteOpenstack(&DeleteOpenstackInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-openstack-3",
			})

			_ = c.DeleteOpenstack(&DeleteOpenstackInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-openstack",
			})
		})
	}()

	if *osCreateResp1.Name != "test-openstack" {
		t.Errorf("bad name: %q", *osCreateResp1.Name)
	}
	if *osCreateResp1.User != "user" {
		t.Errorf("bad user: %q", *osCreateResp1.User)
	}
	if *osCreateResp1.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", *osCreateResp1.BucketName)
	}
	if *osCreateResp1.AccessKey != "secret-key" {
		t.Errorf("bad access_key: %q", *osCreateResp1.AccessKey)
	}
	if *osCreateResp1.Path != "/path" {
		t.Errorf("bad path: %q", *osCreateResp1.Path)
	}
	if *osCreateResp1.URL != "https://logs.example.com/v1.0" {
		t.Errorf("bad url: %q", *osCreateResp1.URL)
	}
	if *osCreateResp1.Period != 12 {
		t.Errorf("bad period: %q", *osCreateResp1.Period)
	}
	if *osCreateResp1.CompressionCodec != "snappy" {
		t.Errorf("bad comprssion_codec: %q", *osCreateResp1.CompressionCodec)
	}
	if *osCreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *osCreateResp1.GzipLevel)
	}
	if *osCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", *osCreateResp1.Format)
	}
	if *osCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *osCreateResp1.FormatVersion)
	}
	if *osCreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", *osCreateResp1.TimestampFormat)
	}
	if *osCreateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", *osCreateResp1.MessageType)
	}
	if *osCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *osCreateResp1.Placement)
	}
	if *osCreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", *osCreateResp1.PublicKey)
	}
	if osCreateResp2.CompressionCodec != nil {
		t.Errorf("bad compression_codec: %q", *osCreateResp2.CompressionCodec)
	}
	if *osCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", *osCreateResp2.GzipLevel)
	}
	if *osCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", *osCreateResp3.CompressionCodec)
	}
	if *osCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *osCreateResp3.GzipLevel)
	}

	// List
	var lc []*Openstack
	Record(t, "openstack/list", func(c *Client) {
		lc, err = c.ListOpenstack(&ListOpenstackInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "openstack/get", func(c *Client) {
		osGetResp, err = c.GetOpenstack(&GetOpenstackInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-openstack",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *osCreateResp1.Name != *osGetResp.Name {
		t.Errorf("bad name: %q", *osCreateResp1.Name)
	}
	if *osCreateResp1.User != *osGetResp.User {
		t.Errorf("bad user: %q", *osCreateResp1.User)
	}
	if *osCreateResp1.BucketName != *osGetResp.BucketName {
		t.Errorf("bad bucket_name: %q", *osCreateResp1.BucketName)
	}
	if *osCreateResp1.AccessKey != *osGetResp.AccessKey {
		t.Errorf("bad access_key: %q", *osCreateResp1.AccessKey)
	}
	if *osCreateResp1.Path != *osGetResp.Path {
		t.Errorf("bad path: %q", *osCreateResp1.Path)
	}
	if *osCreateResp1.URL != *osGetResp.URL {
		t.Errorf("bad url: %q", *osCreateResp1.URL)
	}
	if *osCreateResp1.Period != *osGetResp.Period {
		t.Errorf("bad period: %q", *osCreateResp1.Period)
	}
	if *osCreateResp1.CompressionCodec != *osGetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", *osCreateResp1.CompressionCodec)
	}
	if *osCreateResp1.GzipLevel != *osGetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", *osCreateResp1.GzipLevel)
	}
	if *osCreateResp1.Format != *osGetResp.Format {
		t.Errorf("bad format: %q", *osCreateResp1.Format)
	}
	if *osCreateResp1.FormatVersion != *osGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", *osCreateResp1.FormatVersion)
	}
	if *osCreateResp1.TimestampFormat != *osGetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", *osCreateResp1.TimestampFormat)
	}
	if *osCreateResp1.MessageType != *osGetResp.MessageType {
		t.Errorf("bad message_type: %q", *osCreateResp1.MessageType)
	}
	if *osCreateResp1.Placement != *osGetResp.Placement {
		t.Errorf("bad placement: %q", *osCreateResp1.Placement)
	}
	if *osCreateResp1.PublicKey != *osGetResp.PublicKey {
		t.Errorf("bad public_key: %q", *osCreateResp1.PublicKey)
	}

	// Update
	var osUpdateResp1, osUpdateResp2, osUpdateResp3 *Openstack
	Record(t, "openstack/update", func(c *Client) {
		osUpdateResp1, err = c.UpdateOpenstack(&UpdateOpenstackInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-openstack",
			User:             ToPointer("new-user"),
			NewName:          ToPointer("new-test-openstack"),
			CompressionCodec: ToPointer("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "openstack/update2", func(c *Client) {
		osUpdateResp2, err = c.UpdateOpenstack(&UpdateOpenstackInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-openstack-2",
			CompressionCodec: ToPointer("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "openstack/update3", func(c *Client) {
		osUpdateResp3, err = c.UpdateOpenstack(&UpdateOpenstackInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-openstack-3",
			GzipLevel:      ToPointer(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *osUpdateResp1.Name != "new-test-openstack" {
		t.Errorf("bad name: %q", *osUpdateResp1.Name)
	}
	if *osUpdateResp1.User != "new-user" {
		t.Errorf("bad user: %q", *osUpdateResp1.User)
	}
	if *osUpdateResp1.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", *osUpdateResp1.CompressionCodec)
	}
	if osUpdateResp1.GzipLevel != nil {
		t.Errorf("bad gzip_level: %q", *osUpdateResp1.GzipLevel)
	}
	if *osUpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", *osUpdateResp2.CompressionCodec)
	}
	if *osUpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *osUpdateResp2.GzipLevel)
	}
	if osUpdateResp3.CompressionCodec != nil {
		t.Errorf("bad compression_codec: %q", *osUpdateResp3.CompressionCodec)
	}
	if *osUpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", *osUpdateResp3.GzipLevel)
	}

	// Delete
	Record(t, "openstack/delete", func(c *Client) {
		err = c.DeleteOpenstack(&DeleteOpenstackInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-openstack",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListOpenstack_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListOpenstack(&ListOpenstackInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListOpenstack(&ListOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateOpenstack_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateOpenstack(&CreateOpenstackInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateOpenstack(&CreateOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetOpenstack_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetOpenstack(&GetOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetOpenstack(&GetOpenstackInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetOpenstack(&GetOpenstackInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateOpenstack_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateOpenstack(&UpdateOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateOpenstack(&UpdateOpenstackInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateOpenstack(&UpdateOpenstackInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteOpenstack_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteOpenstack(&DeleteOpenstackInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteOpenstack(&DeleteOpenstackInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteOpenstack(&DeleteOpenstackInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
