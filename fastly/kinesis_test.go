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
	var kinesis *Kinesis
	record(t, "kinesis/create", func(c *Client) {
		kinesis, err = c.CreateKinesis(&CreateKinesisInput{
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
				Name:           "new-test-kinesis",
			})
		})
	}()

	if kinesis.Name != "test-kinesis" {
		t.Errorf("bad name: %q", kinesis.Name)
	}
	if kinesis.StreamName != "stream-name" {
		t.Errorf("bad bucket_name: %q", kinesis.StreamName)
	}
	if kinesis.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("bad access_key: %q", kinesis.AccessKey)
	}
	if kinesis.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", kinesis.SecretKey)
	}
	if kinesis.Region != "us-east-1" {
		t.Errorf("bad domain: %q", kinesis.Region)
	}
	if kinesis.Format != "format" {
		t.Errorf("bad format: %q", kinesis.Format)
	}
	if kinesis.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", kinesis.FormatVersion)
	}
	if kinesis.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", kinesis.Placement)
	}
	if kinesis.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", kinesis.ResponseCondition)
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
	var nkinesis *Kinesis
	record(t, "kinesis/get", func(c *Client) {
		nkinesis, err = c.GetKinesis(&GetKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "test-kinesis",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if kinesis.Name != nkinesis.Name {
		t.Errorf("bad name: %q", kinesis.Name)
	}
	if kinesis.StreamName != nkinesis.StreamName {
		t.Errorf("bad bucket_name: %q", kinesis.StreamName)
	}
	if kinesis.AccessKey != nkinesis.AccessKey {
		t.Errorf("bad access_key: %q", kinesis.AccessKey)
	}
	if kinesis.SecretKey != nkinesis.SecretKey {
		t.Errorf("bad secret_key: %q", kinesis.SecretKey)
	}
	if kinesis.Region != nkinesis.Region {
		t.Errorf("bad domain: %q", kinesis.Region)
	}
	if kinesis.Format != nkinesis.Format {
		t.Errorf("bad format: %q", kinesis.Format)
	}
	if kinesis.FormatVersion != nkinesis.FormatVersion {
		t.Errorf("bad format_version: %q", kinesis.FormatVersion)
	}
	if kinesis.Placement != nkinesis.Placement {
		t.Errorf("bad placement: %q", kinesis.Placement)
	}
	if kinesis.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", kinesis.ResponseCondition)
	}

	// Update
	var ukinesis *Kinesis
	record(t, "kinesis/update", func(c *Client) {
		ukinesis, err = c.UpdateKinesis(&UpdateKinesisInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Name:           "test-kinesis",
			NewName:        String("new-test-kinesis"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ukinesis.Name != "new-test-kinesis" {
		t.Errorf("bad name: %q", ukinesis.Name)
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
