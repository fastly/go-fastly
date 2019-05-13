package fastly

import "testing"

func TestClient_S3s(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "s3s/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var s3 *S3
	record(t, "s3s/create", func(c *Client) {
		s3, err = c.CreateS3(&CreateS3Input{
			Service:         testServiceID,
			Version:         tv.Number,
			Name:            "test-s3",
			BucketName:      "bucket-name",
			Domain:          "s3.us-east-1.amazonaws.com",
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretKey:       "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:            "/path",
			Period:          12,
			GzipLevel:       9,
			Format:          "format",
			FormatVersion:   2,
			TimestampFormat: "%Y",
			MessageType:     "classic",
			Redundancy:      S3RedundancyReduced,
			Placement:       "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "s3s/cleanup", func(c *Client) {
			c.DeleteS3(&DeleteS3Input{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-s3",
			})

			c.DeleteS3(&DeleteS3Input{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-s3",
			})
		})
	}()

	if s3.Name != "test-s3" {
		t.Errorf("bad name: %q", s3.Name)
	}
	if s3.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", s3.BucketName)
	}
	if s3.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("bad access_key: %q", s3.AccessKey)
	}
	if s3.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", s3.SecretKey)
	}
	if s3.Domain != "s3.us-east-1.amazonaws.com" {
		t.Errorf("bad domain: %q", s3.Domain)
	}
	if s3.Path != "/path" {
		t.Errorf("bad path: %q", s3.Path)
	}
	if s3.Period != 12 {
		t.Errorf("bad period: %q", s3.Period)
	}
	if s3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", s3.GzipLevel)
	}
	if s3.Format != "format" {
		t.Errorf("bad format: %q", s3.Format)
	}
	if s3.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", s3.FormatVersion)
	}
	if s3.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", s3.TimestampFormat)
	}
	if s3.Redundancy != S3RedundancyReduced {
		t.Errorf("bad redundancy: %q", s3.Redundancy)
	}
	if s3.MessageType != "classic" {
		t.Errorf("bad message_type: %q", s3.MessageType)
	}
	if s3.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", s3.Placement)
	}

	// List
	var s3s []*S3
	record(t, "s3s/list", func(c *Client) {
		s3s, err = c.ListS3s(&ListS3sInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(s3s) < 1 {
		t.Errorf("bad s3s: %v", s3s)
	}

	// Get
	var ns3 *S3
	record(t, "s3s/get", func(c *Client) {
		ns3, err = c.GetS3(&GetS3Input{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-s3",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s3.Name != ns3.Name {
		t.Errorf("bad name: %q", s3.Name)
	}
	if s3.BucketName != ns3.BucketName {
		t.Errorf("bad bucket_name: %q", s3.BucketName)
	}
	if s3.AccessKey != ns3.AccessKey {
		t.Errorf("bad access_key: %q", s3.AccessKey)
	}
	if s3.SecretKey != ns3.SecretKey {
		t.Errorf("bad secret_key: %q", s3.SecretKey)
	}
	if s3.Domain != ns3.Domain {
		t.Errorf("bad domain: %q", s3.Domain)
	}
	if s3.Path != ns3.Path {
		t.Errorf("bad path: %q", s3.Path)
	}
	if s3.Period != ns3.Period {
		t.Errorf("bad period: %q", s3.Period)
	}
	if s3.GzipLevel != ns3.GzipLevel {
		t.Errorf("bad gzip_level: %q", s3.GzipLevel)
	}
	if s3.Format != ns3.Format {
		t.Errorf("bad format: %q", s3.Format)
	}
	if s3.FormatVersion != ns3.FormatVersion {
		t.Errorf("bad format_version: %q", s3.FormatVersion)
	}
	if s3.TimestampFormat != ns3.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", s3.TimestampFormat)
	}
	if s3.Redundancy != ns3.Redundancy {
		t.Errorf("bad redundancy: %q", s3.Redundancy)
	}
	if s3.Placement != ns3.Placement {
		t.Errorf("bad placement: %q", s3.Placement)
	}

	// Update
	var us3 *S3
	record(t, "s3s/update", func(c *Client) {
		us3, err = c.UpdateS3(&UpdateS3Input{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-s3",
			NewName: "new-test-s3",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us3.Name != "new-test-s3" {
		t.Errorf("bad name: %q", us3.Name)
	}

	// Delete
	record(t, "s3s/delete", func(c *Client) {
		err = c.DeleteS3(&DeleteS3Input{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-s3",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListS3s_validation(t *testing.T) {
	var err error
	_, err = testClient.ListS3s(&ListS3sInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListS3s(&ListS3sInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateS3_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateS3(&CreateS3Input{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateS3(&CreateS3Input{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetS3_validation(t *testing.T) {
	var err error
	_, err = testClient.GetS3(&GetS3Input{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetS3(&GetS3Input{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetS3(&GetS3Input{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateS3_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateS3(&UpdateS3Input{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateS3(&UpdateS3Input{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateS3(&UpdateS3Input{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteS3_validation(t *testing.T) {
	var err error
	err = testClient.DeleteS3(&DeleteS3Input{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteS3(&DeleteS3Input{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteS3(&DeleteS3Input{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
