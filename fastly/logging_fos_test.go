package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_FOSes(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "foses/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var fosCreateResp *FOS
	Record(t, "foses/create", func(c *Client) {
		fosCreateResp, err = c.CreateFOS(context.TODO(), &CreateFOSInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             new("test-fos"),
			BucketName:       new("bucket-name"),
			Region:           new("us-east"),
			AccessKey:        new("AKIAIOSFODNN7EXAMPLE"),
			SecretKey:        new("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			Path:             new("/path"),
			Period:           new(12),
			CompressionCodec: new("snappy"),
			Format:           new("format"),
			FormatVersion:    new(2),
			TimestampFormat:  new("%Y"),
			MessageType:      new("classic"),
			Placement:        new("none"),
			PublicKey:        new(pgpPublicKey()),
			FileMaxBytes:     new(MiB),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "foses/cleanup", func(c *Client) {
			_ = c.DeleteFOS(context.TODO(), &DeleteFOSInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-fos",
			})

			_ = c.DeleteFOS(context.TODO(), &DeleteFOSInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-fos",
			})
		})
	}()

	if *fosCreateResp.Name != "test-fos" {
		t.Errorf("bad name: %q", *fosCreateResp.Name)
	}
	if *fosCreateResp.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", *fosCreateResp.BucketName)
	}
	if *fosCreateResp.Region != "us-east" {
		t.Errorf("bad region: %q", *fosCreateResp.Region)
	}
	if fosCreateResp.AccessKey == nil || *fosCreateResp.AccessKey == "" {
		t.Errorf("missing access_key in response")
	}
	if *fosCreateResp.Path != "/path" {
		t.Errorf("bad path: %q", *fosCreateResp.Path)
	}
	if *fosCreateResp.Period != 12 {
		t.Errorf("bad period: %d", *fosCreateResp.Period)
	}
	if *fosCreateResp.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", *fosCreateResp.CompressionCodec)
	}
	if *fosCreateResp.Format != "format" {
		t.Errorf("bad format: %q", *fosCreateResp.Format)
	}
	if *fosCreateResp.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", *fosCreateResp.FormatVersion)
	}
	if *fosCreateResp.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", *fosCreateResp.TimestampFormat)
	}
	if *fosCreateResp.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", *fosCreateResp.PublicKey)
	}
	if *fosCreateResp.FileMaxBytes != MiB {
		t.Errorf("bad file_max_bytes: %d", *fosCreateResp.FileMaxBytes)
	}

	// List
	var foses []*FOS
	Record(t, "foses/list", func(c *Client) {
		foses, err = c.ListFOSes(context.TODO(), &ListFOSesInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(foses) < 1 {
		t.Errorf("bad foses: %v", foses)
	}

	// Get
	var fosGetResp *FOS
	Record(t, "foses/get", func(c *Client) {
		fosGetResp, err = c.GetFOS(context.TODO(), &GetFOSInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-fos",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *fosCreateResp.Name != *fosGetResp.Name {
		t.Errorf("bad name: %q", *fosGetResp.Name)
	}
	if *fosCreateResp.BucketName != *fosGetResp.BucketName {
		t.Errorf("bad bucket_name: %q", *fosGetResp.BucketName)
	}
	if *fosCreateResp.Region != *fosGetResp.Region {
		t.Errorf("bad region: %q", *fosGetResp.Region)
	}

	// Update
	var fosUpdateResp *FOS
	Record(t, "foses/update", func(c *Client) {
		fosUpdateResp, err = c.UpdateFOS(context.TODO(), &UpdateFOSInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-fos",
			NewName:          new("new-test-fos"),
			PublicKey:        new(pgpPublicKeyUpdate()),
			CompressionCodec: new("zstd"),
			FileMaxBytes:     new(5 * MiB),
			Placement:        NullValue[string](),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *fosUpdateResp.Name != "new-test-fos" {
		t.Errorf("bad name: %q", *fosUpdateResp.Name)
	}
	if *fosUpdateResp.PublicKey != pgpPublicKeyUpdate() {
		t.Errorf("bad public_key: %q", *fosUpdateResp.PublicKey)
	}
	if *fosUpdateResp.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", *fosUpdateResp.CompressionCodec)
	}
	if *fosUpdateResp.FileMaxBytes != 5*MiB {
		t.Errorf("bad file_max_bytes: %d", *fosUpdateResp.FileMaxBytes)
	}
	if fosUpdateResp.Placement != nil {
		t.Errorf("bad placement: %q", *fosUpdateResp.Placement)
	}

	// Delete
	Record(t, "foses/delete", func(c *Client) {
		err = c.DeleteFOS(context.TODO(), &DeleteFOSInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-fos",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListFOSes_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListFOSes(context.TODO(), &ListFOSesInput{ServiceID: ""})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListFOSes(context.TODO(), &ListFOSesInput{ServiceID: "foo", ServiceVersion: 0})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateFOS_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateFOS(context.TODO(), &CreateFOSInput{Name: new("test-service"), ServiceVersion: 1})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateFOS(context.TODO(), &CreateFOSInput{Name: new("test-service"), ServiceID: "foo"})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetFOS_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetFOS(context.TODO(), &GetFOSInput{Name: "", ServiceID: "foo", ServiceVersion: 1})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetFOS(context.TODO(), &GetFOSInput{Name: "test-fos", ServiceID: "", ServiceVersion: 1})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetFOS(context.TODO(), &GetFOSInput{Name: "test-fos", ServiceID: "foo", ServiceVersion: 0})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateFOS_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateFOS(context.TODO(), &UpdateFOSInput{Name: "", ServiceID: "foo", ServiceVersion: 1})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateFOS(context.TODO(), &UpdateFOSInput{Name: "test-fos", ServiceID: "", ServiceVersion: 1})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateFOS(context.TODO(), &UpdateFOSInput{Name: "test-fos", ServiceID: "foo", ServiceVersion: 0})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteFOS_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteFOS(context.TODO(), &DeleteFOSInput{Name: "", ServiceID: "foo", ServiceVersion: 1})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteFOS(context.TODO(), &DeleteFOSInput{Name: "test-fos", ServiceID: "", ServiceVersion: 1})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteFOS(context.TODO(), &DeleteFOSInput{Name: "test-fos", ServiceID: "foo", ServiceVersion: 0})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
