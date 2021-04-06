package fastly

import (
	"testing"
)

func TestClient_Kinesis(t *testing.T) {
	t.Parallel()

	var err error
	var v *Version
	record(t, "kinesis/version", func(c *Client) {
		v = testVersion(t, c)
	})

	// Create
	var kinesisCreateResp1, kinesisCreateResp2 *Kinesis
	record(t, "kinesis/create", func(c *Client) {
		kinesisCreateResp1, err = c.CreateKinesis(&CreateKinesisInput{
			ServiceID:         testServiceID,
			ServiceVersion:    v.Number,
			Name:              "test-kinesis",
			StreamName:        "stream-name",
			Region:            "us-east-1",
			AccessKey:         "AKIAIOSFODNN7EXAMPLE",
			SecretKey:         "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Format:            "format",
			FormatVersion:     2,
			ResponseCondition: "",
			Placement:         "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "kinesis/create2", func(c *Client) {
		kinesisCreateResp2, err = c.CreateKinesis(&CreateKinesisInput{
			ServiceID:         testServiceID,
			ServiceVersion:    v.Number,
			Name:              "test-kinesis-2",
			StreamName:        "stream-name",
			Region:            "us-east-1",
			IAMRole:           "arn:aws:iam::123456789012:role/S3Access",
			Format:            "format",
			FormatVersion:     2,
			ResponseCondition: "",
			Placement:         "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail
	record(t, "kinesis/create3", func(c *Client) {
		_, err = c.CreateKinesis(&CreateKinesisInput{
			ServiceID:         testServiceID,
			ServiceVersion:    v.Number,
			Name:              "test-kinesis-3",
			StreamName:        "stream-name",
			Region:            "us-east-1",
			AccessKey:         "AKIAIOSFODNN7EXAMPLE",
			SecretKey:         "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			IAMRole:           "arn:aws:iam::123456789012:role/S3Access",
			Format:            "format",
			FormatVersion:     2,
			ResponseCondition: "",
			Placement:         "waf_debug",
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// This case is expected to fail
	record(t, "kinesis/create4", func(c *Client) {
		_, err = c.CreateKinesis(&CreateKinesisInput{
			ServiceID:         testServiceID,
			ServiceVersion:    v.Number,
			Name:              "test-kinesis-3",
			StreamName:        "stream-name",
			Region:            "us-east-1",
			IAMRole:           "badarn",
			Format:            "format",
			FormatVersion:     2,
			ResponseCondition: "",
			Placement:         "waf_debug",
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "kinesis/cleanup", func(c *Client) {
			c.DeleteKinesis(&DeleteKinesisInput{
				ServiceID:      testServiceID,
				ServiceVersion: v.Number,
				Name:           "test-kinesis",
			})

			c.DeleteKinesis(&DeleteKinesisInput{
				ServiceID:      testServiceID,
				ServiceVersion: v.Number,
				Name:           "test-kinesis-2",
			})

			c.DeleteKinesis(&DeleteKinesisInput{
				ServiceID:      testServiceID,
				ServiceVersion: v.Number,
				Name:           "new-test-kinesis",
			})
		})
	}()

	if kinesisCreateResp1.Name != "test-kinesis" {
		t.Errorf("bad name: %q", kinesisCreateResp1.Name)
	}
	if kinesisCreateResp1.StreamName != "stream-name" {
		t.Errorf("bad bucket_name: %q", kinesisCreateResp1.StreamName)
	}
	if kinesisCreateResp1.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("bad access_key: %q", kinesisCreateResp1.AccessKey)
	}
	if kinesisCreateResp1.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", kinesisCreateResp1.SecretKey)
	}
	if kinesisCreateResp1.IAMRole != "" {
		t.Errorf("bad iam_role: %q", kinesisCreateResp1.IAMRole)
	}
	if kinesisCreateResp1.Region != "us-east-1" {
		t.Errorf("bad domain: %q", kinesisCreateResp1.Region)
	}
	if kinesisCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", kinesisCreateResp1.Format)
	}
	if kinesisCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", kinesisCreateResp1.FormatVersion)
	}
	if kinesisCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", kinesisCreateResp1.Placement)
	}
	if kinesisCreateResp1.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", kinesisCreateResp1.ResponseCondition)
	}

	if kinesisCreateResp2.AccessKey != "" {
		t.Errorf("bad access_key: %q", kinesisCreateResp2.AccessKey)
	}
	if kinesisCreateResp2.SecretKey != "" {
		t.Errorf("bad secret_key: %q", kinesisCreateResp2.SecretKey)
	}
	if kinesisCreateResp2.IAMRole != "arn:aws:iam::123456789012:role/S3Access" {
		t.Errorf("bad iam_role: %q", kinesisCreateResp2.IAMRole)
	}

	// List
	var kineses []*Kinesis
	record(t, "kinesis/list", func(c *Client) {
		kineses, err = c.ListKinesis(&ListKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(kineses) < 1 {
		t.Errorf("bad kineses: %v", kineses)
	}

	// Get
	var kinesisGetResp, kinesisGetResp2 *Kinesis
	record(t, "kinesis/get", func(c *Client) {
		kinesisGetResp, err = c.GetKinesis(&GetKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "test-kinesis",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "kinesis/get2", func(c *Client) {
		kinesisGetResp2, err = c.GetKinesis(&GetKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "test-kinesis-2",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if kinesisCreateResp1.Name != kinesisGetResp.Name {
		t.Errorf("bad name: %q", kinesisCreateResp1.Name)
	}
	if kinesisCreateResp1.StreamName != kinesisGetResp.StreamName {
		t.Errorf("bad bucket_name: %q", kinesisCreateResp1.StreamName)
	}
	if kinesisCreateResp1.AccessKey != kinesisGetResp.AccessKey {
		t.Errorf("bad access_key: %q", kinesisCreateResp1.AccessKey)
	}
	if kinesisCreateResp1.SecretKey != kinesisGetResp.SecretKey {
		t.Errorf("bad secret_key: %q", kinesisCreateResp1.SecretKey)
	}
	if kinesisCreateResp1.IAMRole != kinesisGetResp.IAMRole {
		t.Errorf("bad iam_role: %q", kinesisCreateResp1.IAMRole)
	}
	if kinesisCreateResp1.Region != kinesisGetResp.Region {
		t.Errorf("bad domain: %q", kinesisCreateResp1.Region)
	}
	if kinesisCreateResp1.Format != kinesisGetResp.Format {
		t.Errorf("bad format: %q", kinesisCreateResp1.Format)
	}
	if kinesisCreateResp1.FormatVersion != kinesisGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", kinesisCreateResp1.FormatVersion)
	}
	if kinesisCreateResp1.Placement != kinesisGetResp.Placement {
		t.Errorf("bad placement: %q", kinesisCreateResp1.Placement)
	}
	if kinesisCreateResp1.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", kinesisCreateResp1.ResponseCondition)
	}

	if kinesisCreateResp2.AccessKey != kinesisGetResp2.AccessKey {
		t.Errorf("bad access_key: %q", kinesisGetResp2.AccessKey)
	}
	if kinesisCreateResp2.SecretKey != kinesisGetResp2.SecretKey {
		t.Errorf("bad secret_key: %q", kinesisGetResp2.SecretKey)
	}
	if kinesisCreateResp2.IAMRole != kinesisGetResp2.IAMRole {
		t.Errorf("bad iam_role: %q", kinesisGetResp2.IAMRole)
	}

	// Update
	var kinesisUpdateResp1, kinesisUpdateResp2, kinesisUpdateResp3 *Kinesis
	record(t, "kinesis/update", func(c *Client) {
		kinesisUpdateResp1, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "test-kinesis",
			NewName:        String("new-test-kinesis"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that a configuration using an access key/secret key can be
	// updated to use IAM role.
	record(t, "kinesis/update2", func(c *Client) {
		kinesisUpdateResp2, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "new-test-kinesis",
			AccessKey:      String(""),
			SecretKey:      String(""),
			IAMRole:        String("arn:aws:iam::123456789012:role/S3Access"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that a configuration using an IAM role can be updated to use
	// access key/secret key.
	record(t, "kinesis/update3", func(c *Client) {
		kinesisUpdateResp3, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "test-kinesis-2",
			AccessKey:      String("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:      String("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			IAMRole:        String(""),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that an invalid IAM role ARN is rejected. This case is expected
	// to fail.
	record(t, "kinesis/update4", func(c *Client) {
		_, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "test-kinesis",
			IAMRole:        String("badarn"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	if kinesisUpdateResp1.Name != "new-test-kinesis" {
		t.Errorf("bad name: %q", kinesisUpdateResp1.Name)
	}

	if kinesisUpdateResp2.AccessKey != "" {
		t.Errorf("bad access_key: %q", kinesisUpdateResp2.AccessKey)
	}
	if kinesisUpdateResp2.SecretKey != "" {
		t.Errorf("bad secret_key: %q", kinesisUpdateResp2.SecretKey)
	}
	if kinesisUpdateResp2.IAMRole != "arn:aws:iam::123456789012:role/S3Access" {
		t.Errorf("bad iam_role: %q", kinesisUpdateResp2.IAMRole)
	}

	if kinesisUpdateResp3.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("bad access_key: %q", kinesisUpdateResp3.AccessKey)
	}
	if kinesisUpdateResp3.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", kinesisUpdateResp3.SecretKey)
	}
	if kinesisUpdateResp3.IAMRole != "" {
		t.Errorf("bad iam_role: %q", kinesisUpdateResp3.IAMRole)
	}

	// Delete
	record(t, "kinesis/delete", func(c *Client) {
		err = c.DeleteKinesis(&DeleteKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "new-test-kinesis",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListKinesis_validation(t *testing.T) {
	var err error
	_, err = testClient.ListKinesis(&ListKinesisInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListKinesis(&ListKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateKinesis_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateKinesis(&CreateKinesisInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateKinesis(&CreateKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetKinesis_validation(t *testing.T) {
	var err error
	_, err = testClient.GetKinesis(&GetKinesisInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetKinesis(&GetKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetKinesis(&GetKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateKinesis_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateKinesis(&UpdateKinesisInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateKinesis(&UpdateKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateKinesis(&UpdateKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteKinesis_validation(t *testing.T) {
	var err error
	err = testClient.DeleteKinesis(&DeleteKinesisInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteKinesis(&DeleteKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteKinesis(&DeleteKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
