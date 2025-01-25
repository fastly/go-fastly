package fastly

import (
	"errors"
	"testing"
)

func TestClient_Kinesis(t *testing.T) {
	t.Parallel()

	var err error
	var v *Version
	Record(t, "kinesis/version", func(c *Client) {
		v = testVersion(t, c)
	})

	// Create
	//
	// NOTE: You can't send the API and empty ResponseCondition.
	var kinesisCreateResp1, kinesisCreateResp2 *Kinesis
	Record(t, "kinesis/create", func(c *Client) {
		kinesisCreateResp1, err = c.CreateKinesis(&CreateKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           ToPointer("test-kinesis"),
			StreamName:     ToPointer("stream-name"),
			Region:         ToPointer("us-east-1"),
			AccessKey:      ToPointer("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:      ToPointer("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			Format:         ToPointer("format"),
			FormatVersion:  ToPointer(2),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "kinesis/create2", func(c *Client) {
		kinesisCreateResp2, err = c.CreateKinesis(&CreateKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           ToPointer("test-kinesis-2"),
			StreamName:     ToPointer("stream-name"),
			Region:         ToPointer("us-east-1"),
			IAMRole:        ToPointer("arn:aws:iam::123456789012:role/S3Access"),
			Format:         ToPointer("format"),
			FormatVersion:  ToPointer(2),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail
	Record(t, "kinesis/create3", func(c *Client) {
		_, err = c.CreateKinesis(&CreateKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           ToPointer("test-kinesis-3"),
			StreamName:     ToPointer("stream-name"),
			Region:         ToPointer("us-east-1"),
			AccessKey:      ToPointer("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:      ToPointer("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			IAMRole:        ToPointer("arn:aws:iam::123456789012:role/S3Access"),
			Format:         ToPointer("format"),
			FormatVersion:  ToPointer(2),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// This case is expected to fail
	Record(t, "kinesis/create4", func(c *Client) {
		_, err = c.CreateKinesis(&CreateKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           ToPointer("test-kinesis-3"),
			StreamName:     ToPointer("stream-name"),
			Region:         ToPointer("us-east-1"),
			IAMRole:        ToPointer("badarn"),
			Format:         ToPointer("format"),
			FormatVersion:  ToPointer(2),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "kinesis/cleanup", func(c *Client) {
			_ = c.DeleteKinesis(&DeleteKinesisInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *v.Number,
				Name:           "test-kinesis",
			})

			_ = c.DeleteKinesis(&DeleteKinesisInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *v.Number,
				Name:           "test-kinesis-2",
			})

			_ = c.DeleteKinesis(&DeleteKinesisInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *v.Number,
				Name:           "new-test-kinesis",
			})
		})
	}()

	if *kinesisCreateResp1.Name != "test-kinesis" {
		t.Errorf("bad name: %q", *kinesisCreateResp1.Name)
	}
	if *kinesisCreateResp1.StreamName != "stream-name" {
		t.Errorf("bad bucket_name: %q", *kinesisCreateResp1.StreamName)
	}
	if *kinesisCreateResp1.AccessKey != "AKIAIOSFODNN7EXAMPLE" { // #nosec G101
		t.Errorf("bad access_key: %q", *kinesisCreateResp1.AccessKey)
	}
	if *kinesisCreateResp1.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", *kinesisCreateResp1.SecretKey)
	}
	if kinesisCreateResp1.IAMRole != nil {
		t.Errorf("bad iam_role: %q", *kinesisCreateResp1.IAMRole)
	}
	if *kinesisCreateResp1.Region != "us-east-1" {
		t.Errorf("bad domain: %q", *kinesisCreateResp1.Region)
	}
	if *kinesisCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", *kinesisCreateResp1.Format)
	}
	if *kinesisCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *kinesisCreateResp1.FormatVersion)
	}
	if *kinesisCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *kinesisCreateResp1.Placement)
	}
	if *kinesisCreateResp1.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", *kinesisCreateResp1.ResponseCondition)
	}
	if kinesisCreateResp2.AccessKey != nil {
		t.Errorf("bad access_key: %q", *kinesisCreateResp2.AccessKey)
	}
	if kinesisCreateResp2.SecretKey != nil {
		t.Errorf("bad secret_key: %q", *kinesisCreateResp2.SecretKey)
	}
	if *kinesisCreateResp2.IAMRole != "arn:aws:iam::123456789012:role/S3Access" {
		t.Errorf("bad iam_role: %q", *kinesisCreateResp2.IAMRole)
	}

	// List
	var kineses []*Kinesis
	Record(t, "kinesis/list", func(c *Client) {
		kineses, err = c.ListKinesis(&ListKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
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
	Record(t, "kinesis/get", func(c *Client) {
		kinesisGetResp, err = c.GetKinesis(&GetKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           "test-kinesis",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "kinesis/get2", func(c *Client) {
		kinesisGetResp2, err = c.GetKinesis(&GetKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           "test-kinesis-2",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *kinesisCreateResp1.Name != *kinesisGetResp.Name {
		t.Errorf("bad name: %q", *kinesisCreateResp1.Name)
	}
	if *kinesisCreateResp1.StreamName != *kinesisGetResp.StreamName {
		t.Errorf("bad bucket_name: %q", *kinesisCreateResp1.StreamName)
	}
	if *kinesisCreateResp1.AccessKey != *kinesisGetResp.AccessKey {
		t.Errorf("bad access_key: %q", *kinesisCreateResp1.AccessKey)
	}
	if *kinesisCreateResp1.SecretKey != *kinesisGetResp.SecretKey {
		t.Errorf("bad secret_key: %q", *kinesisCreateResp1.SecretKey)
	}
	if kinesisCreateResp1.IAMRole != nil && kinesisGetResp.IAMRole != nil {
		if *kinesisCreateResp1.IAMRole != *kinesisGetResp.IAMRole {
			t.Errorf("bad iam_role: %q", *kinesisCreateResp1.IAMRole)
		}
	}
	if *kinesisCreateResp1.Region != *kinesisGetResp.Region {
		t.Errorf("bad domain: %q", *kinesisCreateResp1.Region)
	}
	if *kinesisCreateResp1.Format != *kinesisGetResp.Format {
		t.Errorf("bad format: %q", *kinesisCreateResp1.Format)
	}
	if *kinesisCreateResp1.FormatVersion != *kinesisGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", *kinesisCreateResp1.FormatVersion)
	}
	if *kinesisCreateResp1.Placement != *kinesisGetResp.Placement {
		t.Errorf("bad placement: %q", *kinesisCreateResp1.Placement)
	}
	if *kinesisCreateResp1.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", *kinesisCreateResp1.ResponseCondition)
	}
	if kinesisCreateResp2.AccessKey != nil && kinesisGetResp2.AccessKey != nil {
		if *kinesisCreateResp2.AccessKey != *kinesisGetResp2.AccessKey {
			t.Errorf("bad access_key: %q", *kinesisGetResp2.AccessKey)
		}
	}
	if kinesisCreateResp2.SecretKey != nil && kinesisGetResp2.SecretKey != nil {
		if *kinesisCreateResp2.SecretKey != *kinesisGetResp2.SecretKey {
			t.Errorf("bad secret_key: %q", *kinesisGetResp2.SecretKey)
		}
	}
	if *kinesisCreateResp2.IAMRole != *kinesisGetResp2.IAMRole {
		t.Errorf("bad iam_role: %q", *kinesisGetResp2.IAMRole)
	}

	// Update
	var kinesisUpdateResp1, kinesisUpdateResp2, kinesisUpdateResp3 *Kinesis
	Record(t, "kinesis/update", func(c *Client) {
		kinesisUpdateResp1, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           "test-kinesis",
			NewName:        ToPointer("new-test-kinesis"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that a configuration using an access key/secret key can be
	// updated to use IAM role.
	Record(t, "kinesis/update2", func(c *Client) {
		kinesisUpdateResp2, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           "new-test-kinesis",
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
	Record(t, "kinesis/update3", func(c *Client) {
		kinesisUpdateResp3, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           "test-kinesis-2",
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
	Record(t, "kinesis/update4", func(c *Client) {
		_, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           "test-kinesis",
			IAMRole:        ToPointer("badarn"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	if *kinesisUpdateResp1.Name != "new-test-kinesis" {
		t.Errorf("bad name: %q", *kinesisUpdateResp1.Name)
	}
	if *kinesisUpdateResp2.AccessKey != "" {
		t.Errorf("bad access_key: %q", *kinesisUpdateResp2.AccessKey)
	}
	if *kinesisUpdateResp2.SecretKey != "" {
		t.Errorf("bad secret_key: %q", *kinesisUpdateResp2.SecretKey)
	}
	if *kinesisUpdateResp2.IAMRole != "arn:aws:iam::123456789012:role/S3Access" {
		t.Errorf("bad iam_role: %q", *kinesisUpdateResp2.IAMRole)
	}
	if *kinesisUpdateResp3.AccessKey != "AKIAIOSFODNN7EXAMPLE" { // #nosec G101
		t.Errorf("bad access_key: %q", *kinesisUpdateResp3.AccessKey)
	}
	if *kinesisUpdateResp3.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", *kinesisUpdateResp3.SecretKey)
	}
	if *kinesisUpdateResp3.IAMRole != "" {
		t.Errorf("bad iam_role: %q", *kinesisUpdateResp3.IAMRole)
	}

	// Delete
	Record(t, "kinesis/delete", func(c *Client) {
		err = c.DeleteKinesis(&DeleteKinesisInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Name:           "new-test-kinesis",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListKinesis_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListKinesis(&ListKinesisInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListKinesis(&ListKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateKinesis_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateKinesis(&CreateKinesisInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateKinesis(&CreateKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetKinesis_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetKinesis(&GetKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetKinesis(&GetKinesisInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetKinesis(&GetKinesisInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateKinesis_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateKinesis(&UpdateKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateKinesis(&UpdateKinesisInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateKinesis(&UpdateKinesisInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteKinesis_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteKinesis(&DeleteKinesisInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteKinesis(&DeleteKinesisInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteKinesis(&DeleteKinesisInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
