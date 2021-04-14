package fastly

import (
	"testing"
)

func TestClient_DigitalOceans(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "digitaloceans/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var digitaloceanCreateResp1, digitaloceanCreateResp2, digitaloceanCreateResp3 *DigitalOcean
	record(t, "digitaloceans/create", func(c *Client) {
		digitaloceanCreateResp1, err = c.CreateDigitalOcean(&CreateDigitalOceanInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-digitalocean",
			BucketName:       "bucket-name",
			Domain:           "fra1.digitaloceanspaces.com",
			AccessKey:        "AKIAIOSFODNN7EXAMPLE",
			SecretKey:        "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:             "/path",
			Period:           12,
			Format:           "format",
			FormatVersion:    2,
			TimestampFormat:  "%Y",
			MessageType:      "classic",
			Placement:        "waf_debug",
			PublicKey:        pgpPublicKey(),
			CompressionCodec: "snappy",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "digitaloceans/create2", func(c *Client) {
		digitaloceanCreateResp2, err = c.CreateDigitalOcean(&CreateDigitalOceanInput{
			ServiceID:       testServiceID,
			ServiceVersion:  tv.Number,
			Name:            "test-digitalocean-2",
			BucketName:      "bucket-name",
			Domain:          "fra1.digitaloceanspaces.com",
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretKey:       "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:            "/path",
			Period:          12,
			GzipLevel:       8,
			Format:          "format",
			FormatVersion:   2,
			TimestampFormat: "%Y",
			MessageType:     "classic",
			Placement:       "waf_debug",
			PublicKey:       pgpPublicKey(),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "digitaloceans/create3", func(c *Client) {
		digitaloceanCreateResp3, err = c.CreateDigitalOcean(&CreateDigitalOceanInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-digitalocean-3",
			BucketName:       "bucket-name",
			Domain:           "fra1.digitaloceanspaces.com",
			AccessKey:        "AKIAIOSFODNN7EXAMPLE",
			SecretKey:        "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:             "/path",
			Period:           12,
			Format:           "format",
			FormatVersion:    2,
			TimestampFormat:  "%Y",
			MessageType:      "classic",
			Placement:        "waf_debug",
			PublicKey:        pgpPublicKey(),
			CompressionCodec: "snappy",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail because both CompressionCodec and
	// GzipLevel are present.
	record(t, "digitaloceans/create4", func(c *Client) {
		_, err = c.CreateDigitalOcean(&CreateDigitalOceanInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-digitalocean-4",
			BucketName:       "bucket-name",
			Domain:           "fra1.digitaloceanspaces.com",
			AccessKey:        "AKIAIOSFODNN7EXAMPLE",
			SecretKey:        "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:             "/path",
			Period:           12,
			GzipLevel:        8,
			Format:           "format",
			FormatVersion:    2,
			TimestampFormat:  "%Y",
			MessageType:      "classic",
			Placement:        "waf_debug",
			PublicKey:        pgpPublicKey(),
			CompressionCodec: "snappy",
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "digitaloceans/cleanup", func(c *Client) {
			c.DeleteDigitalOcean(&DeleteDigitalOceanInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-digitalocean",
			})

			c.DeleteDigitalOcean(&DeleteDigitalOceanInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-digitalocean-2",
			})

			c.DeleteDigitalOcean(&DeleteDigitalOceanInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-digitalocean-3",
			})

			c.DeleteDigitalOcean(&DeleteDigitalOceanInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-digitalocean",
			})
		})
	}()

	if digitaloceanCreateResp1.Name != "test-digitalocean" {
		t.Errorf("bad name: %q", digitaloceanCreateResp1.Name)
	}
	if digitaloceanCreateResp1.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", digitaloceanCreateResp1.BucketName)
	}
	if digitaloceanCreateResp1.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("bad access_key: %q", digitaloceanCreateResp1.AccessKey)
	}
	if digitaloceanCreateResp1.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", digitaloceanCreateResp1.SecretKey)
	}
	if digitaloceanCreateResp1.Domain != "fra1.digitaloceanspaces.com" {
		t.Errorf("bad domain: %q", digitaloceanCreateResp1.Domain)
	}
	if digitaloceanCreateResp1.Path != "/path" {
		t.Errorf("bad path: %q", digitaloceanCreateResp1.Path)
	}
	if digitaloceanCreateResp1.Period != 12 {
		t.Errorf("bad period: %q", digitaloceanCreateResp1.Period)
	}
	if digitaloceanCreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", digitaloceanCreateResp1.GzipLevel)
	}
	if digitaloceanCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", digitaloceanCreateResp1.Format)
	}
	if digitaloceanCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", digitaloceanCreateResp1.FormatVersion)
	}
	if digitaloceanCreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", digitaloceanCreateResp1.TimestampFormat)
	}
	if digitaloceanCreateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", digitaloceanCreateResp1.MessageType)
	}
	if digitaloceanCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", digitaloceanCreateResp1.Placement)
	}
	if digitaloceanCreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", digitaloceanCreateResp1.PublicKey)
	}

	if digitaloceanCreateResp2.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", digitaloceanCreateResp2.CompressionCodec)
	}
	if digitaloceanCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", digitaloceanCreateResp2.GzipLevel)
	}

	if digitaloceanCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", digitaloceanCreateResp3.CompressionCodec)
	}
	if digitaloceanCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", digitaloceanCreateResp3.GzipLevel)
	}

	// List
	var digitaloceans []*DigitalOcean
	record(t, "digitaloceans/list", func(c *Client) {
		digitaloceans, err = c.ListDigitalOceans(&ListDigitalOceansInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(digitaloceans) < 1 {
		t.Errorf("bad digitaloceans: %v", digitaloceans)
	}

	// Get
	var digitaloceanGetResp *DigitalOcean
	record(t, "digitaloceans/get", func(c *Client) {
		digitaloceanGetResp, err = c.GetDigitalOcean(&GetDigitalOceanInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-digitalocean",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if digitaloceanCreateResp1.Name != digitaloceanGetResp.Name {
		t.Errorf("bad name: %q", digitaloceanCreateResp1.Name)
	}
	if digitaloceanCreateResp1.BucketName != digitaloceanGetResp.BucketName {
		t.Errorf("bad bucket_name: %q", digitaloceanCreateResp1.BucketName)
	}
	if digitaloceanCreateResp1.AccessKey != digitaloceanGetResp.AccessKey {
		t.Errorf("bad access_key: %q", digitaloceanCreateResp1.AccessKey)
	}
	if digitaloceanCreateResp1.SecretKey != digitaloceanGetResp.SecretKey {
		t.Errorf("bad secret_key: %q", digitaloceanCreateResp1.SecretKey)
	}
	if digitaloceanCreateResp1.Domain != digitaloceanGetResp.Domain {
		t.Errorf("bad domain: %q", digitaloceanCreateResp1.Domain)
	}
	if digitaloceanCreateResp1.Path != digitaloceanGetResp.Path {
		t.Errorf("bad path: %q", digitaloceanCreateResp1.Path)
	}
	if digitaloceanCreateResp1.Period != digitaloceanGetResp.Period {
		t.Errorf("bad period: %q", digitaloceanCreateResp1.Period)
	}
	if digitaloceanCreateResp1.GzipLevel != digitaloceanGetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", digitaloceanCreateResp1.GzipLevel)
	}
	if digitaloceanCreateResp1.Format != digitaloceanGetResp.Format {
		t.Errorf("bad format: %q", digitaloceanCreateResp1.Format)
	}
	if digitaloceanCreateResp1.FormatVersion != digitaloceanGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", digitaloceanCreateResp1.FormatVersion)
	}
	if digitaloceanCreateResp1.TimestampFormat != digitaloceanGetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", digitaloceanCreateResp1.TimestampFormat)
	}
	if digitaloceanCreateResp1.Placement != digitaloceanGetResp.Placement {
		t.Errorf("bad placement: %q", digitaloceanCreateResp1.Placement)
	}
	if digitaloceanCreateResp1.PublicKey != digitaloceanGetResp.PublicKey {
		t.Errorf("bad public_key: %q", digitaloceanCreateResp1.PublicKey)
	}
	if digitaloceanCreateResp1.CompressionCodec != digitaloceanGetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", digitaloceanCreateResp1.CompressionCodec)
	}

	// Update
	var digitaloceanUpdateResp1, digitaloceanUpdateResp2, digitaloceanUpdateResp3 *DigitalOcean
	record(t, "digitaloceans/update", func(c *Client) {
		digitaloceanUpdateResp1, err = c.UpdateDigitalOcean(&UpdateDigitalOceanInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-digitalocean",
			NewName:          String("new-test-digitalocean"),
			Domain:           String("nyc3.digitaloceanspaces.com"),
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "digitaloceans/update2", func(c *Client) {
		digitaloceanUpdateResp2, err = c.UpdateDigitalOcean(&UpdateDigitalOceanInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-digitalocean-2",
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "digitaloceans/update3", func(c *Client) {
		digitaloceanUpdateResp3, err = c.UpdateDigitalOcean(&UpdateDigitalOceanInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-digitalocean-3",
			GzipLevel:      Uint(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if digitaloceanUpdateResp1.Name != "new-test-digitalocean" {
		t.Errorf("bad name: %q", digitaloceanUpdateResp1.Name)
	}
	if digitaloceanUpdateResp1.Domain != "nyc3.digitaloceanspaces.com" {
		t.Errorf("bad domain: %q", digitaloceanUpdateResp1.Domain)
	}
	if digitaloceanUpdateResp1.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", digitaloceanUpdateResp1.CompressionCodec)
	}
	if digitaloceanUpdateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", digitaloceanUpdateResp1.GzipLevel)
	}

	if digitaloceanUpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", digitaloceanUpdateResp2.CompressionCodec)
	}
	if digitaloceanUpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", digitaloceanUpdateResp2.GzipLevel)
	}

	if digitaloceanUpdateResp3.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", digitaloceanUpdateResp3.CompressionCodec)
	}
	if digitaloceanUpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", digitaloceanUpdateResp3.GzipLevel)
	}

	// Delete
	record(t, "digitaloceans/delete", func(c *Client) {
		err = c.DeleteDigitalOcean(&DeleteDigitalOceanInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-digitalocean",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDigitalOceans_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDigitalOceans(&ListDigitalOceansInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDigitalOceans(&ListDigitalOceansInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDigitalOcean_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDigitalOcean(&CreateDigitalOceanInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDigitalOcean(&CreateDigitalOceanInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDigitalOcean_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDigitalOcean(&GetDigitalOceanInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDigitalOcean(&GetDigitalOceanInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDigitalOcean(&GetDigitalOceanInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDigitalOcean_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDigitalOcean(&UpdateDigitalOceanInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDigitalOcean(&UpdateDigitalOceanInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDigitalOcean(&UpdateDigitalOceanInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDigitalOcean_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDigitalOcean(&DeleteDigitalOceanInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDigitalOcean(&DeleteDigitalOceanInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDigitalOcean(&DeleteDigitalOceanInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
