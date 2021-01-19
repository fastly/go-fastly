package fastly

import (
	"testing"
)

func TestClient_S3s(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "s3s/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var s3CreateResp1, s3CreateResp2, s3CreateResp3 *S3
	record(t, "s3s/create", func(c *Client) {
		s3CreateResp1, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    testServiceID,
			ServiceVersion:               tv.Number,
			Name:                         "test-s3",
			BucketName:                   "bucket-name",
			Domain:                       "s3.us-east-1.amazonaws.com",
			AccessKey:                    "AKIAIOSFODNN7EXAMPLE",
			SecretKey:                    "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:                         "/path",
			Period:                       12,
			CompressionCodec:             "snappy",
			Format:                       "format",
			FormatVersion:                2,
			ResponseCondition:            "",
			TimestampFormat:              "%Y",
			MessageType:                  "classic",
			Redundancy:                   S3RedundancyReduced,
			Placement:                    "waf_debug",
			PublicKey:                    pgpPublicKey(),
			ServerSideEncryptionKMSKeyID: "1234",
			ServerSideEncryption:         S3ServerSideEncryptionKMS,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "s3s/create2", func(c *Client) {
		s3CreateResp2, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    testServiceID,
			ServiceVersion:               tv.Number,
			Name:                         "test-s3-2",
			BucketName:                   "bucket-name",
			Domain:                       "s3.us-east-1.amazonaws.com",
			AccessKey:                    "AKIAIOSFODNN7EXAMPLE",
			SecretKey:                    "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:                         "/path",
			Period:                       12,
			CompressionCodec:             "snappy",
			GzipLevel:                    8,
			Format:                       "format",
			FormatVersion:                2,
			ResponseCondition:            "",
			TimestampFormat:              "%Y",
			MessageType:                  "classic",
			Redundancy:                   S3RedundancyReduced,
			Placement:                    "waf_debug",
			PublicKey:                    pgpPublicKey(),
			ServerSideEncryptionKMSKeyID: "1234",
			ServerSideEncryption:         S3ServerSideEncryptionKMS,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "s3s/create3", func(c *Client) {
		s3CreateResp3, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    testServiceID,
			ServiceVersion:               tv.Number,
			Name:                         "test-s3-3",
			BucketName:                   "bucket-name",
			Domain:                       "s3.us-east-1.amazonaws.com",
			AccessKey:                    "AKIAIOSFODNN7EXAMPLE",
			SecretKey:                    "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:                         "/path",
			Period:                       12,
			CompressionCodec:             "snappy",
			Format:                       "format",
			FormatVersion:                2,
			ResponseCondition:            "",
			TimestampFormat:              "%Y",
			MessageType:                  "classic",
			Redundancy:                   S3RedundancyReduced,
			Placement:                    "waf_debug",
			PublicKey:                    pgpPublicKey(),
			ServerSideEncryptionKMSKeyID: "1234",
			ServerSideEncryption:         S3ServerSideEncryptionKMS,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "s3s/cleanup", func(c *Client) {
			c.DeleteS3(&DeleteS3Input{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-s3",
			})

			c.DeleteS3(&DeleteS3Input{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-s3-2",
			})

			c.DeleteS3(&DeleteS3Input{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-s3-3",
			})

			c.DeleteS3(&DeleteS3Input{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-s3",
			})
		})
	}()

	if s3CreateResp1.Name != "test-s3" {
		t.Errorf("bad name: %q", s3CreateResp1.Name)
	}
	if s3CreateResp1.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", s3CreateResp1.BucketName)
	}
	if s3CreateResp1.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("bad access_key: %q", s3CreateResp1.AccessKey)
	}
	if s3CreateResp1.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", s3CreateResp1.SecretKey)
	}
	if s3CreateResp1.Domain != "s3.us-east-1.amazonaws.com" {
		t.Errorf("bad domain: %q", s3CreateResp1.Domain)
	}
	if s3CreateResp1.Path != "/path" {
		t.Errorf("bad path: %q", s3CreateResp1.Path)
	}
	if s3CreateResp1.Period != 12 {
		t.Errorf("bad period: %q", s3CreateResp1.Period)
	}
	if s3CreateResp1.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", s3CreateResp1.CompressionCodec)
	}
	if s3CreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", s3CreateResp1.GzipLevel)
	}
	if s3CreateResp1.Format != "format" {
		t.Errorf("bad format: %q", s3CreateResp1.Format)
	}
	if s3CreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", s3CreateResp1.FormatVersion)
	}
	if s3CreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", s3CreateResp1.TimestampFormat)
	}
	if s3CreateResp1.Redundancy != S3RedundancyReduced {
		t.Errorf("bad redundancy: %q", s3CreateResp1.Redundancy)
	}
	if s3CreateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", s3CreateResp1.MessageType)
	}
	if s3CreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", s3CreateResp1.Placement)
	}
	if s3CreateResp1.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", s3CreateResp1.ResponseCondition)
	}
	if s3CreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", s3CreateResp1.PublicKey)
	}
	if s3CreateResp1.ServerSideEncryption != S3ServerSideEncryptionKMS {
		t.Errorf("bad server_side_encryption: %q", s3CreateResp1.ServerSideEncryption)
	}
	if s3CreateResp1.ServerSideEncryptionKMSKeyID != "1234" {
		t.Errorf("bad server_side_encryption_kms_key_id: %q", s3CreateResp1.ServerSideEncryptionKMSKeyID)
	}

	if s3CreateResp2.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", s3CreateResp1.CompressionCodec)
	}
	if s3CreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", s3CreateResp1.GzipLevel)
	}

	if s3CreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", s3CreateResp1.CompressionCodec)
	}
	if s3CreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", s3CreateResp1.GzipLevel)
	}

	// List
	var s3s []*S3
	record(t, "s3s/list", func(c *Client) {
		s3s, err = c.ListS3s(&ListS3sInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(s3s) < 1 {
		t.Errorf("bad s3s: %v", s3s)
	}

	// Get
	var s3GetResp *S3
	record(t, "s3s/get", func(c *Client) {
		s3GetResp, err = c.GetS3(&GetS3Input{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-s3",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s3CreateResp1.Name != s3GetResp.Name {
		t.Errorf("bad name: %q", s3CreateResp1.Name)
	}
	if s3CreateResp1.BucketName != s3GetResp.BucketName {
		t.Errorf("bad bucket_name: %q", s3CreateResp1.BucketName)
	}
	if s3CreateResp1.AccessKey != s3GetResp.AccessKey {
		t.Errorf("bad access_key: %q", s3CreateResp1.AccessKey)
	}
	if s3CreateResp1.SecretKey != s3GetResp.SecretKey {
		t.Errorf("bad secret_key: %q", s3CreateResp1.SecretKey)
	}
	if s3CreateResp1.Domain != s3GetResp.Domain {
		t.Errorf("bad domain: %q", s3CreateResp1.Domain)
	}
	if s3CreateResp1.Path != s3GetResp.Path {
		t.Errorf("bad path: %q", s3CreateResp1.Path)
	}
	if s3CreateResp1.Period != s3GetResp.Period {
		t.Errorf("bad period: %q", s3CreateResp1.Period)
	}
	if s3CreateResp1.CompressionCodec != s3GetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", s3CreateResp1.CompressionCodec)
	}
	if s3CreateResp1.GzipLevel != s3GetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", s3CreateResp1.GzipLevel)
	}
	if s3CreateResp1.Format != s3GetResp.Format {
		t.Errorf("bad format: %q", s3CreateResp1.Format)
	}
	if s3CreateResp1.FormatVersion != s3GetResp.FormatVersion {
		t.Errorf("bad format_version: %q", s3CreateResp1.FormatVersion)
	}
	if s3CreateResp1.TimestampFormat != s3GetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", s3CreateResp1.TimestampFormat)
	}
	if s3CreateResp1.Redundancy != s3GetResp.Redundancy {
		t.Errorf("bad redundancy: %q", s3CreateResp1.Redundancy)
	}
	if s3CreateResp1.Placement != s3GetResp.Placement {
		t.Errorf("bad placement: %q", s3CreateResp1.Placement)
	}
	if s3CreateResp1.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", s3CreateResp1.ResponseCondition)
	}
	if s3CreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", s3CreateResp1.PublicKey)
	}
	if s3CreateResp1.ServerSideEncryption != s3GetResp.ServerSideEncryption {
		t.Errorf("bad server_side_encryption: %q", s3CreateResp1.ServerSideEncryption)
	}
	if s3CreateResp1.ServerSideEncryptionKMSKeyID != s3GetResp.ServerSideEncryptionKMSKeyID {
		t.Errorf("bad server_side_encryption_kms_key_id: %q", s3CreateResp1.ServerSideEncryptionKMSKeyID)
	}

	// Update
	var s3UpdateResp1, s3UpdateResp2, s3UpdateResp3 *S3
	record(t, "s3s/update", func(c *Client) {
		s3UpdateResp1, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-s3",
			NewName:          String("new-test-s3"),
			PublicKey:        String(pgpPublicKeyUpdate()),
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that CompressionCodec can be set for a an endpoint where
	// GzipLevel was specified at creation time.
	record(t, "s3s/update2", func(c *Client) {
		s3UpdateResp2, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-s3-2",
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that GzipLevel can be set for an endpoint where CompressionCodec
	// was set at creation time.
	record(t, "s3s/update3", func(c *Client) {
		s3UpdateResp3, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-s3-3",
			GzipLevel:      Uint(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if s3UpdateResp1.Name != "new-test-s3" {
		t.Errorf("bad name: %q", s3UpdateResp1.Name)
	}
	if s3UpdateResp1.PublicKey != pgpPublicKeyUpdate() {
		t.Errorf("bad public_key: %q", s3UpdateResp1.PublicKey)
	}
	if s3UpdateResp1.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", s3UpdateResp1.CompressionCodec)
	}
	if s3UpdateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", s3UpdateResp1.GzipLevel)
	}

	if s3UpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", s3UpdateResp2.CompressionCodec)
	}
	if s3UpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", s3UpdateResp2.GzipLevel)
	}

	if s3UpdateResp3.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", s3UpdateResp3.CompressionCodec)
	}
	if s3UpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", s3UpdateResp3.GzipLevel)
	}

	// Delete
	record(t, "s3s/delete", func(c *Client) {
		err = c.DeleteS3(&DeleteS3Input{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-s3",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListS3s_validation(t *testing.T) {
	var err error
	_, err = testClient.ListS3s(&ListS3sInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListS3s(&ListS3sInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateS3_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateS3(&CreateS3Input{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateS3(&CreateS3Input{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateS3(&CreateS3Input{
		ServiceID:                    "foo",
		ServiceVersion:               1,
		Name:                         "test-service",
		ServerSideEncryption:         S3ServerSideEncryptionKMS,
		ServerSideEncryptionKMSKeyID: "",
	})
	if err != ErrMissingServerSideEncryptionKMSKeyID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetS3_validation(t *testing.T) {
	var err error
	_, err = testClient.GetS3(&GetS3Input{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetS3(&GetS3Input{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetS3(&GetS3Input{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateS3_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateS3(&UpdateS3Input{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateS3(&UpdateS3Input{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateS3(&UpdateS3Input{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateS3(&UpdateS3Input{
		ServiceID:                    "foo",
		ServiceVersion:               1,
		Name:                         "test-service",
		ServerSideEncryption:         S3ServerSideEncryptionKMS,
		ServerSideEncryptionKMSKeyID: String(""),
	})
	if err != ErrMissingServerSideEncryptionKMSKeyID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteS3_validation(t *testing.T) {
	var err error
	err = testClient.DeleteS3(&DeleteS3Input{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteS3(&DeleteS3Input{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteS3(&DeleteS3Input{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
