package fastly

import (
	"testing"
)

func TestClient_Cloudfiles(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "cloudfiles/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var cloudfilesCreateResp1, cloudfilesCreateResp2, cloudfilesCreateResp3 *Cloudfiles
	record(t, "cloudfiles/create", func(c *Client) {
		cloudfilesCreateResp1, err = c.CreateCloudfiles(&CreateCloudfilesInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-cloudfiles",
			User:             "user",
			AccessKey:        "secret-key",
			BucketName:       "bucket-name",
			Path:             "/path",
			Region:           "DFW",
			Period:           12,
			Format:           "format",
			FormatVersion:    1,
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

	record(t, "cloudfiles/create2", func(c *Client) {
		cloudfilesCreateResp2, err = c.CreateCloudfiles(&CreateCloudfilesInput{
			ServiceID:       testServiceID,
			ServiceVersion:  tv.Number,
			Name:            "test-cloudfiles-2",
			User:            "user",
			AccessKey:       "secret-key",
			BucketName:      "bucket-name",
			Path:            "/path",
			Region:          "DFW",
			Period:          12,
			GzipLevel:       8,
			Format:          "format",
			FormatVersion:   1,
			TimestampFormat: "%Y",
			MessageType:     "classic",
			Placement:       "waf_debug",
			PublicKey:       pgpPublicKey(),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "cloudfiles/create3", func(c *Client) {
		cloudfilesCreateResp3, err = c.CreateCloudfiles(&CreateCloudfilesInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-cloudfiles-3",
			User:             "user",
			AccessKey:        "secret-key",
			BucketName:       "bucket-name",
			Path:             "/path",
			Region:           "DFW",
			Period:           12,
			Format:           "format",
			FormatVersion:    1,
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
	record(t, "cloudfiles/create4", func(c *Client) {
		_, err = c.CreateCloudfiles(&CreateCloudfilesInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-cloudfiles-4",
			User:             "user",
			AccessKey:        "secret-key",
			BucketName:       "bucket-name",
			Path:             "/path",
			Region:           "DFW",
			Period:           12,
			GzipLevel:        8,
			Format:           "format",
			FormatVersion:    1,
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
		record(t, "cloudfiles/cleanup", func(c *Client) {
			c.DeleteCloudfiles(&DeleteCloudfilesInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-cloudfiles",
			})

			c.DeleteCloudfiles(&DeleteCloudfilesInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-cloudfiles-2",
			})

			c.DeleteCloudfiles(&DeleteCloudfilesInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-cloudfiles-3",
			})

			c.DeleteCloudfiles(&DeleteCloudfilesInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-cloudfiles",
			})
		})
	}()

	if cloudfilesCreateResp1.Name != "test-cloudfiles" {
		t.Errorf("bad name: %q", cloudfilesCreateResp1.Name)
	}
	if cloudfilesCreateResp1.User != "user" {
		t.Errorf("bad user: %q", cloudfilesCreateResp1.User)
	}
	if cloudfilesCreateResp1.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", cloudfilesCreateResp1.BucketName)
	}
	if cloudfilesCreateResp1.AccessKey != "secret-key" {
		t.Errorf("bad access_key: %q", cloudfilesCreateResp1.AccessKey)
	}
	if cloudfilesCreateResp1.Path != "/path" {
		t.Errorf("bad path: %q", cloudfilesCreateResp1.Path)
	}
	if cloudfilesCreateResp1.Region != "DFW" {
		t.Errorf("bad region: %q", cloudfilesCreateResp1.Region)
	}
	if cloudfilesCreateResp1.Period != 12 {
		t.Errorf("bad period: %q", cloudfilesCreateResp1.Period)
	}
	if cloudfilesCreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", cloudfilesCreateResp1.GzipLevel)
	}
	if cloudfilesCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", cloudfilesCreateResp1.Format)
	}
	if cloudfilesCreateResp1.FormatVersion != 1 {
		t.Errorf("bad format_version: %q", cloudfilesCreateResp1.FormatVersion)
	}
	if cloudfilesCreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", cloudfilesCreateResp1.TimestampFormat)
	}
	if cloudfilesCreateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", cloudfilesCreateResp1.MessageType)
	}
	if cloudfilesCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", cloudfilesCreateResp1.Placement)
	}
	if cloudfilesCreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", cloudfilesCreateResp1.PublicKey)
	}

	if cloudfilesCreateResp2.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", cloudfilesCreateResp2.CompressionCodec)
	}
	if cloudfilesCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", cloudfilesCreateResp2.GzipLevel)
	}

	if cloudfilesCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", cloudfilesCreateResp3.CompressionCodec)
	}
	if cloudfilesCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", cloudfilesCreateResp3.GzipLevel)
	}

	// List
	var lc []*Cloudfiles
	record(t, "cloudfiles/list", func(c *Client) {
		lc, err = c.ListCloudfiles(&ListCloudfilesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lc) < 1 {
		t.Errorf("bad cloudfiles: %v", lc)
	}

	// Get
	var cloudfilesGetResp *Cloudfiles
	record(t, "cloudfiles/get", func(c *Client) {
		cloudfilesGetResp, err = c.GetCloudfiles(&GetCloudfilesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-cloudfiles",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if cloudfilesCreateResp1.Name != cloudfilesGetResp.Name {
		t.Errorf("bad name: %q", cloudfilesCreateResp1.Name)
	}
	if cloudfilesCreateResp1.User != cloudfilesGetResp.User {
		t.Errorf("bad user: %q", cloudfilesCreateResp1.User)
	}
	if cloudfilesCreateResp1.BucketName != cloudfilesGetResp.BucketName {
		t.Errorf("bad bucket_name: %q", cloudfilesCreateResp1.BucketName)
	}
	if cloudfilesCreateResp1.AccessKey != cloudfilesGetResp.AccessKey {
		t.Errorf("bad access_key: %q", cloudfilesCreateResp1.AccessKey)
	}
	if cloudfilesCreateResp1.Path != cloudfilesGetResp.Path {
		t.Errorf("bad path: %q", cloudfilesCreateResp1.Path)
	}
	if cloudfilesCreateResp1.Region != cloudfilesGetResp.Region {
		t.Errorf("bad region: %q", cloudfilesCreateResp1.Region)
	}
	if cloudfilesCreateResp1.Period != cloudfilesGetResp.Period {
		t.Errorf("bad period: %q", cloudfilesCreateResp1.Period)
	}
	if cloudfilesCreateResp1.GzipLevel != cloudfilesGetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", cloudfilesCreateResp1.GzipLevel)
	}
	if cloudfilesCreateResp1.Format != cloudfilesGetResp.Format {
		t.Errorf("bad format: %q", cloudfilesCreateResp1.Format)
	}
	if cloudfilesCreateResp1.FormatVersion != cloudfilesGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", cloudfilesCreateResp1.FormatVersion)
	}
	if cloudfilesCreateResp1.TimestampFormat != cloudfilesGetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", cloudfilesCreateResp1.TimestampFormat)
	}
	if cloudfilesCreateResp1.MessageType != cloudfilesGetResp.MessageType {
		t.Errorf("bad message_type: %q", cloudfilesCreateResp1.MessageType)
	}
	if cloudfilesCreateResp1.Placement != cloudfilesGetResp.Placement {
		t.Errorf("bad placement: %q", cloudfilesCreateResp1.Placement)
	}
	if cloudfilesCreateResp1.PublicKey != cloudfilesGetResp.PublicKey {
		t.Errorf("bad public_key: %q", cloudfilesCreateResp1.PublicKey)
	}
	if cloudfilesCreateResp1.CompressionCodec != cloudfilesGetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", cloudfilesCreateResp1.CompressionCodec)
	}

	// Update
	var cloudfilesUpdateResp1, cloudfilesUpdateResp2, cloudfilesUpdateResp3 *Cloudfiles
	record(t, "cloudfiles/update", func(c *Client) {
		cloudfilesUpdateResp1, err = c.UpdateCloudfiles(&UpdateCloudfilesInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-cloudfiles",
			NewName:          String("new-test-cloudfiles"),
			User:             String("new-user"),
			Period:           Uint(0),
			FormatVersion:    Uint(2),
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "cloudfiles/update2", func(c *Client) {
		cloudfilesUpdateResp2, err = c.UpdateCloudfiles(&UpdateCloudfilesInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-cloudfiles-2",
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "cloudfiles/update3", func(c *Client) {
		cloudfilesUpdateResp3, err = c.UpdateCloudfiles(&UpdateCloudfilesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-cloudfiles-3",
			GzipLevel:      Uint(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if cloudfilesUpdateResp1.Name != "new-test-cloudfiles" {
		t.Errorf("bad name: %q", cloudfilesUpdateResp1.Name)
	}
	if cloudfilesUpdateResp1.User != "new-user" {
		t.Errorf("bad user: %q", cloudfilesUpdateResp1.User)
	}
	if cloudfilesUpdateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", cloudfilesUpdateResp1.GzipLevel)
	}
	if cloudfilesUpdateResp1.Period != 0 {
		t.Errorf("bad period: %q", cloudfilesUpdateResp1.Period)
	}
	if cloudfilesUpdateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", cloudfilesUpdateResp1.FormatVersion)
	}
	if cloudfilesUpdateResp1.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", cloudfilesUpdateResp1.CompressionCodec)
	}
	if cloudfilesUpdateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", cloudfilesUpdateResp1.GzipLevel)
	}

	if cloudfilesUpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", cloudfilesUpdateResp2.CompressionCodec)
	}
	if cloudfilesUpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", cloudfilesUpdateResp2.GzipLevel)
	}

	if cloudfilesUpdateResp3.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", cloudfilesUpdateResp3.CompressionCodec)
	}
	if cloudfilesUpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", cloudfilesUpdateResp3.GzipLevel)
	}

	// Delete
	record(t, "cloudfiles/delete", func(c *Client) {
		err = c.DeleteCloudfiles(&DeleteCloudfilesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-cloudfiles",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListCloudfiles_validation(t *testing.T) {
	var err error
	_, err = testClient.ListCloudfiles(&ListCloudfilesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListCloudfiles(&ListCloudfilesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateCloudfiles_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateCloudfiles(&CreateCloudfilesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateCloudfiles(&CreateCloudfilesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetCloudfiles_validation(t *testing.T) {
	var err error
	_, err = testClient.GetCloudfiles(&GetCloudfilesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCloudfiles(&GetCloudfilesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCloudfiles(&GetCloudfilesInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCloudfiles_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateCloudfiles(&UpdateCloudfilesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCloudfiles(&UpdateCloudfilesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCloudfiles(&UpdateCloudfilesInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCloudfiles_validation(t *testing.T) {
	var err error
	err = testClient.DeleteCloudfiles(&DeleteCloudfilesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCloudfiles(&DeleteCloudfilesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCloudfiles(&DeleteCloudfilesInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
