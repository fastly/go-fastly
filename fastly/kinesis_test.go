package fastly

import "testing"

func TestClient_Kineses(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "kineses/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var kinesis *Kinesis
	record(t, "kineses/create", func(c *Client) {
		kinesis, err = c.CreateKinesis(&CreateKinesisInput{
			Service:           testServiceID,
			Version:           tv.Number,
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
		record(t, "kineses/cleanup", func(c *Client) {
			c.DeleteKinesis(&DeleteKinesisInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-kinesis",
			})

			c.DeleteKinesis(&DeleteKinesisInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-kinesis",
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
	record(t, "kineses/list", func(c *Client) {
		kineses, err = c.ListKineses(&ListKinesesInput{
			Service: testServiceID,
			Version: tv.Number,
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
	record(t, "kineses/get", func(c *Client) {
		nkinesis, err = c.GetKinesis(&GetKinesisInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-kinesis",
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
	record(t, "kineses/update", func(c *Client) {
		ukinesis, err = c.UpdateKinesis(&UpdateKinesisInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-kinesis",
			NewName: "new-test-kinesis",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ukinesis.Name != "new-test-kinesis" {
		t.Errorf("bad name: %q", ukinesis.Name)
	}

	// Delete
	record(t, "kineses/delete", func(c *Client) {
		err = c.DeleteKinesis(&DeleteKinesisInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-kinesis",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListKineses_validation(t *testing.T) {
	var err error
	_, err = testClient.ListKineses(&ListKinesesInput{
		Service: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListKineses(&ListKinesesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateKinesis_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateKinesis(&CreateKinesisInput{
		Service: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateKinesis(&CreateKinesisInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetKinesis_validation(t *testing.T) {
	var err error
	_, err = testClient.GetKinesis(&GetKinesisInput{
		Service: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetKinesis(&GetKinesisInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetKinesis(&GetKinesisInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateKinesis_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateKinesis(&UpdateKinesisInput{
		Service: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateKinesis(&UpdateKinesisInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateKinesis(&UpdateKinesisInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteKinesis_validation(t *testing.T) {
	var err error
	err = testClient.DeleteKinesis(&DeleteKinesisInput{
		Service: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteKinesis(&DeleteKinesisInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteKinesis(&DeleteKinesisInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
