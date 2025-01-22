package fastly

import (
	"errors"
	"testing"
)

func TestClient_S3s(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "s3s/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	//
	// NOTE: You can't send the API and empty ResponseCondition.
	var s3CreateResp1, s3CreateResp2, s3CreateResp3, s3CreateResp4 *S3
	Record(t, "s3s/create", func(c *Client) {
		s3CreateResp1, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    TestDeliveryServiceID,
			ServiceVersion:               *tv.Number,
			Name:                         ToPointer("test-s3"),
			BucketName:                   ToPointer("bucket-name"),
			Domain:                       ToPointer("s3.us-east-1.amazonaws.com"),
			AccessKey:                    ToPointer("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:                    ToPointer("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			Path:                         ToPointer("/path"),
			Period:                       ToPointer(12),
			CompressionCodec:             ToPointer("snappy"),
			Format:                       ToPointer("format"),
			FormatVersion:                ToPointer(2),
			TimestampFormat:              ToPointer("%Y"),
			MessageType:                  ToPointer("classic"),
			Redundancy:                   ToPointer(S3RedundancyReduced),
			Placement:                    ToPointer("waf_debug"),
			PublicKey:                    ToPointer(pgpPublicKey()),
			ServerSideEncryptionKMSKeyID: ToPointer("1234"),
			ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
			ACL:                          ToPointer(S3AccessControlListPrivate),
			FileMaxBytes:                 ToPointer(MiB),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "s3s/create2", func(c *Client) {
		s3CreateResp2, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    TestDeliveryServiceID,
			ServiceVersion:               *tv.Number,
			Name:                         ToPointer("test-s3-2"),
			BucketName:                   ToPointer("bucket-name"),
			Domain:                       ToPointer("s3.us-east-1.amazonaws.com"),
			AccessKey:                    ToPointer("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:                    ToPointer("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			Path:                         ToPointer("/path"),
			Period:                       ToPointer(12),
			GzipLevel:                    ToPointer(8),
			Format:                       ToPointer("format"),
			FormatVersion:                ToPointer(2),
			TimestampFormat:              ToPointer("%Y"),
			MessageType:                  ToPointer("classic"),
			Redundancy:                   ToPointer(S3RedundancyOneZoneIA),
			Placement:                    ToPointer("waf_debug"),
			PublicKey:                    ToPointer(pgpPublicKey()),
			ServerSideEncryptionKMSKeyID: ToPointer("1234"),
			ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
			ACL:                          ToPointer(S3AccessControlListAuthenticatedRead),
			FileMaxBytes:                 ToPointer(10 * MiB),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	Record(t, "s3s/create3", func(c *Client) {
		s3CreateResp3, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    TestDeliveryServiceID,
			ServiceVersion:               *tv.Number,
			Name:                         ToPointer("test-s3-3"),
			BucketName:                   ToPointer("bucket-name"),
			Domain:                       ToPointer("s3.us-east-1.amazonaws.com"),
			AccessKey:                    ToPointer("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:                    ToPointer("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			Path:                         ToPointer("/path"),
			Period:                       ToPointer(12),
			Format:                       ToPointer("format"),
			CompressionCodec:             ToPointer("snappy"),
			FormatVersion:                ToPointer(2),
			TimestampFormat:              ToPointer("%Y"),
			MessageType:                  ToPointer("classic"),
			Redundancy:                   ToPointer(S3RedundancyStandardIA),
			Placement:                    ToPointer("waf_debug"),
			PublicKey:                    ToPointer(pgpPublicKey()),
			ServerSideEncryptionKMSKeyID: ToPointer("1234"),
			ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
			ACL:                          ToPointer(S3AccessControlListBucketOwnerFullControl),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "s3s/create4", func(c *Client) {
		s3CreateResp4, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    TestDeliveryServiceID,
			ServiceVersion:               *tv.Number,
			Name:                         ToPointer("test-s3-4"),
			BucketName:                   ToPointer("bucket-name"),
			Domain:                       ToPointer("s3.us-east-1.amazonaws.com"),
			IAMRole:                      ToPointer("arn:aws:iam::123456789012:role/S3Access"),
			Path:                         ToPointer("/path"),
			Period:                       ToPointer(12),
			Format:                       ToPointer("format"),
			CompressionCodec:             ToPointer("snappy"),
			FormatVersion:                ToPointer(2),
			TimestampFormat:              ToPointer("%Y"),
			MessageType:                  ToPointer("classic"),
			Placement:                    ToPointer("waf_debug"),
			Redundancy:                   ToPointer(S3RedundancyStandard),
			PublicKey:                    ToPointer(pgpPublicKey()),
			ServerSideEncryptionKMSKeyID: ToPointer("1234"),
			ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
			FileMaxBytes:                 ToPointer(10 * MiB),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail
	Record(t, "s3s/create5", func(c *Client) {
		_, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    TestDeliveryServiceID,
			ServiceVersion:               *tv.Number,
			Name:                         ToPointer("test-s3"),
			BucketName:                   ToPointer("bucket-name"),
			Domain:                       ToPointer("s3.us-east-1.amazonaws.com"),
			AccessKey:                    ToPointer("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:                    ToPointer("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			IAMRole:                      ToPointer("arn:aws:iam::123456789012:role/S3Access"),
			Path:                         ToPointer("/path"),
			Period:                       ToPointer(12),
			CompressionCodec:             ToPointer("snappy"),
			Format:                       ToPointer("format"),
			FormatVersion:                ToPointer(2),
			TimestampFormat:              ToPointer("%Y"),
			MessageType:                  ToPointer("classic"),
			Redundancy:                   ToPointer(S3RedundancyReduced),
			Placement:                    ToPointer("waf_debug"),
			PublicKey:                    ToPointer(pgpPublicKey()),
			ServerSideEncryptionKMSKeyID: ToPointer("1234"),
			ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// This case is expected to fail
	Record(t, "s3s/create6", func(c *Client) {
		_, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    TestDeliveryServiceID,
			ServiceVersion:               *tv.Number,
			Name:                         ToPointer("test-s3"),
			BucketName:                   ToPointer("bucket-name"),
			Domain:                       ToPointer("s3.us-east-1.amazonaws.com"),
			IAMRole:                      ToPointer("badarn"),
			Path:                         ToPointer("/path"),
			Period:                       ToPointer(12),
			CompressionCodec:             ToPointer("snappy"),
			Format:                       ToPointer("format"),
			FormatVersion:                ToPointer(2),
			TimestampFormat:              ToPointer("%Y"),
			MessageType:                  ToPointer("classic"),
			Redundancy:                   ToPointer(S3RedundancyReduced),
			Placement:                    ToPointer("waf_debug"),
			PublicKey:                    ToPointer(pgpPublicKey()),
			ServerSideEncryptionKMSKeyID: ToPointer("1234"),
			ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// This case is expected to fail because both CompressionCodec and
	// GzipLevel are present.
	Record(t, "s3s/create7", func(c *Client) {
		_, err = c.CreateS3(&CreateS3Input{
			ServiceID:                    TestDeliveryServiceID,
			ServiceVersion:               *tv.Number,
			Name:                         ToPointer("test-s3-2"),
			BucketName:                   ToPointer("bucket-name"),
			Domain:                       ToPointer("s3.us-east-1.amazonaws.com"),
			AccessKey:                    ToPointer("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:                    ToPointer("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			Path:                         ToPointer("/path"),
			Period:                       ToPointer(12),
			CompressionCodec:             ToPointer("snappy"),
			GzipLevel:                    ToPointer(8),
			Format:                       ToPointer("format"),
			FormatVersion:                ToPointer(2),
			TimestampFormat:              ToPointer("%Y"),
			MessageType:                  ToPointer("classic"),
			Redundancy:                   ToPointer(S3RedundancyReduced),
			Placement:                    ToPointer("waf_debug"),
			PublicKey:                    ToPointer(pgpPublicKey()),
			ServerSideEncryptionKMSKeyID: ToPointer("1234"),
			ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "s3s/cleanup", func(c *Client) {
			_ = c.DeleteS3(&DeleteS3Input{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-s3",
			})

			_ = c.DeleteS3(&DeleteS3Input{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-s3-3",
			})

			_ = c.DeleteS3(&DeleteS3Input{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-s3-4",
			})

			_ = c.DeleteS3(&DeleteS3Input{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-s3",
			})
		})
	}()

	if *s3CreateResp1.Name != "test-s3" {
		t.Errorf("bad name: %q", *s3CreateResp1.Name)
	}
	if *s3CreateResp1.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", *s3CreateResp1.BucketName)
	}
	if *s3CreateResp1.AccessKey != "AKIAIOSFODNN7EXAMPLE" { // #nosec G101
		t.Errorf("bad access_key: %q", *s3CreateResp1.AccessKey)
	}
	if *s3CreateResp1.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", *s3CreateResp1.SecretKey)
	}
	if s3CreateResp1.IAMRole != nil {
		t.Errorf("bad iam_role: %q", *s3CreateResp1.IAMRole)
	}
	if *s3CreateResp1.Domain != "s3.us-east-1.amazonaws.com" {
		t.Errorf("bad domain: %q", *s3CreateResp1.Domain)
	}
	if *s3CreateResp1.Path != "/path" {
		t.Errorf("bad path: %q", *s3CreateResp1.Path)
	}
	if *s3CreateResp1.Period != 12 {
		t.Errorf("bad period: %q", *s3CreateResp1.Period)
	}
	if *s3CreateResp1.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", *s3CreateResp1.CompressionCodec)
	}
	if *s3CreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *s3CreateResp1.GzipLevel)
	}
	if *s3CreateResp1.Format != "format" {
		t.Errorf("bad format: %q", *s3CreateResp1.Format)
	}
	if *s3CreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *s3CreateResp1.FormatVersion)
	}
	if *s3CreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", *s3CreateResp1.TimestampFormat)
	}
	if *s3CreateResp1.Redundancy != S3RedundancyReduced {
		t.Errorf("bad redundancy: %q", *s3CreateResp1.Redundancy)
	}
	if *s3CreateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", *s3CreateResp1.MessageType)
	}
	if *s3CreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *s3CreateResp1.Placement)
	}
	if *s3CreateResp1.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", *s3CreateResp1.ResponseCondition)
	}
	if *s3CreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", *s3CreateResp1.PublicKey)
	}
	if *s3CreateResp1.ServerSideEncryption != S3ServerSideEncryptionKMS {
		t.Errorf("bad server_side_encryption: %q", *s3CreateResp1.ServerSideEncryption)
	}
	if *s3CreateResp1.ServerSideEncryptionKMSKeyID != "1234" {
		t.Errorf("bad server_side_encryption_kms_key_id: %q", *s3CreateResp1.ServerSideEncryptionKMSKeyID)
	}
	if *s3CreateResp1.ACL != S3AccessControlListPrivate {
		t.Errorf("bad acl: %s", *s3CreateResp1.ACL)
	}
	if *s3CreateResp1.FileMaxBytes != MiB {
		t.Errorf("bad file_max_bytes: %q", *s3CreateResp1.FileMaxBytes)
	}
	if *s3CreateResp2.FileMaxBytes != 10*MiB {
		t.Errorf("bad file_max_bytes: %q", *s3CreateResp2.FileMaxBytes)
	}
	if *s3CreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", *s3CreateResp1.CompressionCodec)
	}
	if *s3CreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *s3CreateResp1.GzipLevel)
	}
	if *s3CreateResp3.Redundancy != S3RedundancyStandardIA {
		t.Errorf("bad acl: %s", *s3CreateResp3.Redundancy)
	}
	if *s3CreateResp3.ACL != S3AccessControlListBucketOwnerFullControl {
		t.Errorf("bad acl: %s", *s3CreateResp3.ACL)
	}
	if *s3CreateResp3.FileMaxBytes != 0 {
		t.Errorf("bad file_max_bytes: %q", *s3CreateResp3.FileMaxBytes)
	}
	if s3CreateResp4.AccessKey != nil {
		t.Errorf("bad access_key: %q", *s3CreateResp4.AccessKey)
	}
	if s3CreateResp4.SecretKey != nil {
		t.Errorf("bad secret_key: %q", *s3CreateResp4.SecretKey)
	}
	if *s3CreateResp4.IAMRole != "arn:aws:iam::123456789012:role/S3Access" {
		t.Errorf("bad iam_role: %q", *s3CreateResp4.IAMRole)
	}
	if *s3CreateResp4.Redundancy != S3RedundancyStandard {
		t.Errorf("bad acl: %s", *s3CreateResp4.Redundancy)
	}
	if s3CreateResp4.ACL != nil {
		t.Errorf("bad acl: %s", *s3CreateResp4.ACL)
	}
	if *s3CreateResp4.FileMaxBytes != 10*MiB {
		t.Errorf("bad file_max_bytes: %q", *s3CreateResp4.FileMaxBytes)
	}

	// List
	var s3s []*S3
	Record(t, "s3s/list", func(c *Client) {
		s3s, err = c.ListS3s(&ListS3sInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(s3s) < 1 {
		t.Errorf("bad s3s: %v", s3s)
	}

	// Get
	var s3GetResp, s3GetResp2 *S3
	Record(t, "s3s/get", func(c *Client) {
		s3GetResp, err = c.GetS3(&GetS3Input{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-s3",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Request the configuration using the IAM role
	Record(t, "s3s/get2", func(c *Client) {
		s3GetResp2, err = c.GetS3(&GetS3Input{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-s3-4",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *s3CreateResp1.Name != *s3GetResp.Name {
		t.Errorf("bad name: %q", *s3CreateResp1.Name)
	}
	if *s3CreateResp1.BucketName != *s3GetResp.BucketName {
		t.Errorf("bad bucket_name: %q", *s3CreateResp1.BucketName)
	}
	if *s3CreateResp1.AccessKey != *s3GetResp.AccessKey {
		t.Errorf("bad access_key: %q", *s3CreateResp1.AccessKey)
	}
	if *s3CreateResp1.SecretKey != *s3GetResp.SecretKey {
		t.Errorf("bad secret_key: %q", *s3CreateResp1.SecretKey)
	}
	if s3CreateResp1.IAMRole != s3GetResp.IAMRole {
		t.Errorf("bad iam_role: %q", *s3CreateResp1.IAMRole)
	}
	if *s3CreateResp1.Domain != *s3GetResp.Domain {
		t.Errorf("bad domain: %q", *s3CreateResp1.Domain)
	}
	if *s3CreateResp1.Path != *s3GetResp.Path {
		t.Errorf("bad path: %q", *s3CreateResp1.Path)
	}
	if *s3CreateResp1.Period != *s3GetResp.Period {
		t.Errorf("bad period: %q", *s3CreateResp1.Period)
	}
	if *s3CreateResp1.CompressionCodec != *s3GetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", *s3CreateResp1.CompressionCodec)
	}
	if *s3CreateResp1.GzipLevel != *s3GetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", *s3CreateResp1.GzipLevel)
	}
	if *s3CreateResp1.Format != *s3GetResp.Format {
		t.Errorf("bad format: %q", *s3CreateResp1.Format)
	}
	if *s3CreateResp1.FormatVersion != *s3GetResp.FormatVersion {
		t.Errorf("bad format_version: %q", *s3CreateResp1.FormatVersion)
	}
	if *s3CreateResp1.TimestampFormat != *s3GetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", *s3CreateResp1.TimestampFormat)
	}
	if *s3CreateResp1.Redundancy != *s3GetResp.Redundancy {
		t.Errorf("bad redundancy: %q", *s3CreateResp1.Redundancy)
	}
	if *s3CreateResp1.Placement != *s3GetResp.Placement {
		t.Errorf("bad placement: %q", *s3CreateResp1.Placement)
	}
	if *s3CreateResp1.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", *s3CreateResp1.ResponseCondition)
	}
	if *s3CreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", *s3CreateResp1.PublicKey)
	}
	if *s3CreateResp1.ServerSideEncryption != *s3GetResp.ServerSideEncryption {
		t.Errorf("bad server_side_encryption: %q", *s3CreateResp1.ServerSideEncryption)
	}
	if *s3CreateResp1.ServerSideEncryptionKMSKeyID != *s3GetResp.ServerSideEncryptionKMSKeyID {
		t.Errorf("bad server_side_encryption_kms_key_id: %q", *s3CreateResp1.ServerSideEncryptionKMSKeyID)
	}
	if *s3CreateResp1.ACL != *s3GetResp.ACL {
		t.Errorf("bad acl: %s", *s3CreateResp1.ACL)
	}
	if s3CreateResp4.AccessKey != s3GetResp2.AccessKey {
		t.Errorf("bad access_key: %q", *s3CreateResp4.AccessKey)
	}
	if s3CreateResp4.SecretKey != s3GetResp2.SecretKey {
		t.Errorf("bad secret_key: %q", *s3CreateResp4.SecretKey)
	}
	if *s3CreateResp4.IAMRole != *s3GetResp2.IAMRole {
		t.Errorf("bad iam_role: %q", *s3CreateResp4.IAMRole)
	}
	if *s3CreateResp4.Redundancy != *s3GetResp2.Redundancy {
		t.Errorf("bad redundancy: %q", *s3CreateResp4.Redundancy)
	}
	if s3CreateResp4.ACL != s3GetResp2.ACL {
		t.Errorf("bad acl: %s", *s3CreateResp4.ACL)
	}

	// Update
	var s3UpdateResp1, s3UpdateResp2, s3UpdateResp3, s3UpdateResp4, s3UpdateResp5 *S3
	Record(t, "s3s/update", func(c *Client) {
		s3UpdateResp1, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-s3",
			NewName:          ToPointer("new-test-s3"),
			PublicKey:        ToPointer(pgpPublicKeyUpdate()),
			CompressionCodec: ToPointer("zstd"),
			FileMaxBytes:     ToPointer(5 * MiB),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that CompressionCodec can be set for a an endpoint where
	// GzipLevel was specified at creation time.
	Record(t, "s3s/update2", func(c *Client) {
		s3UpdateResp2, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-s3-2",
			CompressionCodec: ToPointer("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that GzipLevel can be set for an endpoint where CompressionCodec
	// was set at creation time.
	Record(t, "s3s/update3", func(c *Client) {
		s3UpdateResp3, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-s3-3",
			GzipLevel:      ToPointer(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that a configuration using an access key/secret key can be
	// updated to use IAM role.
	Record(t, "s3s/update4", func(c *Client) {
		s3UpdateResp4, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-s3",
			AccessKey:      ToPointer(""),
			SecretKey:      ToPointer(""),
			IAMRole:        ToPointer("arn:aws:iam::123456789012:role/S3Access"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that a configuration using an IAM role can be updated to use
	// access key/secret key.
	Record(t, "s3s/update5", func(c *Client) {
		s3UpdateResp5, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-s3-4",
			AccessKey:      ToPointer("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:      ToPointer("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			IAMRole:        ToPointer(""),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that an invalid IAM role ARN is rejected. This case is expected
	// to fail.
	Record(t, "s3s/update6", func(c *Client) {
		_, err = c.UpdateS3(&UpdateS3Input{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-s3",
			IAMRole:        ToPointer("badarn"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	if *s3UpdateResp1.Name != "new-test-s3" {
		t.Errorf("bad name: %q", *s3UpdateResp1.Name)
	}
	if *s3UpdateResp1.PublicKey != pgpPublicKeyUpdate() {
		t.Errorf("bad public_key: %q", *s3UpdateResp1.PublicKey)
	}
	if *s3UpdateResp1.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", *s3UpdateResp1.CompressionCodec)
	}
	if s3UpdateResp1.GzipLevel != nil {
		t.Errorf("bad gzip_level: %q", *s3UpdateResp1.GzipLevel)
	}
	if *s3UpdateResp1.FileMaxBytes != 5*MiB {
		t.Errorf("bad file_max_bytes: %q", *s3UpdateResp1.FileMaxBytes)
	}
	if *s3UpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", *s3UpdateResp2.CompressionCodec)
	}
	if *s3UpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *s3UpdateResp2.GzipLevel)
	}
	if s3UpdateResp3.CompressionCodec != nil {
		t.Errorf("bad compression_codec: %q", *s3UpdateResp3.CompressionCodec)
	}
	if *s3UpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", *s3UpdateResp3.GzipLevel)
	}
	if *s3UpdateResp4.AccessKey != "" {
		t.Errorf("bad access_key: %q", *s3UpdateResp4.AccessKey)
	}
	if *s3UpdateResp4.SecretKey != "" {
		t.Errorf("bad secret_key: %q", *s3UpdateResp4.SecretKey)
	}
	if *s3UpdateResp4.IAMRole != "arn:aws:iam::123456789012:role/S3Access" {
		t.Errorf("bad iam_role: %q", *s3UpdateResp4.IAMRole)
	}
	if *s3UpdateResp5.AccessKey != "AKIAIOSFODNN7EXAMPLE" { // #nosec G101
		t.Errorf("bad access_key: %q", *s3UpdateResp5.AccessKey)
	}
	if *s3UpdateResp5.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", *s3UpdateResp5.SecretKey)
	}
	if *s3UpdateResp5.IAMRole != "" {
		t.Errorf("bad iam_role: %q", *s3UpdateResp5.IAMRole)
	}

	// Delete
	Record(t, "s3s/delete", func(c *Client) {
		err = c.DeleteS3(&DeleteS3Input{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-s3",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListS3s_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListS3s(&ListS3sInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListS3s(&ListS3sInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateS3_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateS3(&CreateS3Input{
		Name:                         ToPointer("test-service"),
		ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
		ServerSideEncryptionKMSKeyID: ToPointer("1234"),
		ServiceVersion:               1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateS3(&CreateS3Input{
		Name:                         ToPointer("test-service"),
		ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
		ServerSideEncryptionKMSKeyID: ToPointer("1234"),
		ServiceID:                    "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateS3(&CreateS3Input{
		Name:                         ToPointer("test-service"),
		ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
		ServerSideEncryptionKMSKeyID: ToPointer(""),
		ServiceID:                    "foo",
		ServiceVersion:               1,
	})
	if !errors.Is(err, ErrMissingServerSideEncryptionKMSKeyID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetS3_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetS3(&GetS3Input{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetS3(&GetS3Input{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetS3(&GetS3Input{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateS3_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateS3(&UpdateS3Input{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateS3(&UpdateS3Input{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateS3(&UpdateS3Input{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateS3(&UpdateS3Input{
		Name:                         "test-service",
		ServerSideEncryption:         ToPointer(S3ServerSideEncryptionKMS),
		ServerSideEncryptionKMSKeyID: ToPointer(""),
		ServiceID:                    "foo",
		ServiceVersion:               1,
	})
	if !errors.Is(err, ErrMissingServerSideEncryptionKMSKeyID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteS3_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteS3(&DeleteS3Input{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteS3(&DeleteS3Input{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteS3(&DeleteS3Input{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
