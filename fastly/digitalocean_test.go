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
	var digitalocean *DigitalOcean
	record(t, "digitaloceans/create", func(c *Client) {
		digitalocean, err = c.CreateDigitalOcean(&CreateDigitalOceanInput{
			ServiceID:       testServiceID,
			ServiceVersion:  tv.Number,
			Name:            "test-digitalocean",
			BucketName:      "bucket-name",
			Domain:          "fra1.digitaloceanspaces.com",
			AccessKey:       "AKIAIOSFODNN7EXAMPLE",
			SecretKey:       "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			Path:            "/path",
			Period:          12,
			GzipLevel:       9,
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
				Name:           "new-test-digitalocean",
			})
		})
	}()

	if digitalocean.Name != "test-digitalocean" {
		t.Errorf("bad name: %q", digitalocean.Name)
	}
	if digitalocean.BucketName != "bucket-name" {
		t.Errorf("bad bucket_name: %q", digitalocean.BucketName)
	}
	if digitalocean.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("bad access_key: %q", digitalocean.AccessKey)
	}
	if digitalocean.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("bad secret_key: %q", digitalocean.SecretKey)
	}
	if digitalocean.Domain != "fra1.digitaloceanspaces.com" {
		t.Errorf("bad domain: %q", digitalocean.Domain)
	}
	if digitalocean.Path != "/path" {
		t.Errorf("bad path: %q", digitalocean.Path)
	}
	if digitalocean.Period != 12 {
		t.Errorf("bad period: %q", digitalocean.Period)
	}
	if digitalocean.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", digitalocean.GzipLevel)
	}
	if digitalocean.Format != "format" {
		t.Errorf("bad format: %q", digitalocean.Format)
	}
	if digitalocean.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", digitalocean.FormatVersion)
	}
	if digitalocean.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", digitalocean.TimestampFormat)
	}
	if digitalocean.MessageType != "classic" {
		t.Errorf("bad message_type: %q", digitalocean.MessageType)
	}
	if digitalocean.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", digitalocean.Placement)
	}
	if digitalocean.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", digitalocean.PublicKey)
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
	var ndigitalocean *DigitalOcean
	record(t, "digitaloceans/get", func(c *Client) {
		ndigitalocean, err = c.GetDigitalOcean(&GetDigitalOceanInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-digitalocean",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if digitalocean.Name != ndigitalocean.Name {
		t.Errorf("bad name: %q", digitalocean.Name)
	}
	if digitalocean.BucketName != ndigitalocean.BucketName {
		t.Errorf("bad bucket_name: %q", digitalocean.BucketName)
	}
	if digitalocean.AccessKey != ndigitalocean.AccessKey {
		t.Errorf("bad access_key: %q", digitalocean.AccessKey)
	}
	if digitalocean.SecretKey != ndigitalocean.SecretKey {
		t.Errorf("bad secret_key: %q", digitalocean.SecretKey)
	}
	if digitalocean.Domain != ndigitalocean.Domain {
		t.Errorf("bad domain: %q", digitalocean.Domain)
	}
	if digitalocean.Path != ndigitalocean.Path {
		t.Errorf("bad path: %q", digitalocean.Path)
	}
	if digitalocean.Period != ndigitalocean.Period {
		t.Errorf("bad period: %q", digitalocean.Period)
	}
	if digitalocean.GzipLevel != ndigitalocean.GzipLevel {
		t.Errorf("bad gzip_level: %q", digitalocean.GzipLevel)
	}
	if digitalocean.Format != ndigitalocean.Format {
		t.Errorf("bad format: %q", digitalocean.Format)
	}
	if digitalocean.FormatVersion != ndigitalocean.FormatVersion {
		t.Errorf("bad format_version: %q", digitalocean.FormatVersion)
	}
	if digitalocean.TimestampFormat != ndigitalocean.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", digitalocean.TimestampFormat)
	}
	if digitalocean.Placement != ndigitalocean.Placement {
		t.Errorf("bad placement: %q", digitalocean.Placement)
	}
	if digitalocean.PublicKey != ndigitalocean.PublicKey {
		t.Errorf("bad public_key: %q", digitalocean.PublicKey)
	}

	// Update
	var udigitalocean *DigitalOcean
	record(t, "digitaloceans/update", func(c *Client) {
		udigitalocean, err = c.UpdateDigitalOcean(&UpdateDigitalOceanInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-digitalocean",
			NewName:        String("new-test-digitalocean"),
			Domain:         String("nyc3.digitaloceanspaces.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if udigitalocean.Name != "new-test-digitalocean" {
		t.Errorf("bad name: %q", udigitalocean.Name)
	}
	if udigitalocean.Domain != "nyc3.digitaloceanspaces.com" {
		t.Errorf("bad domain: %q", udigitalocean.Domain)
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
