package fastly

import (
	"testing"
)

func TestClient_Logentries(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "logentries/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var le *Logentries
	record(t, "logentries/create", func(c *Client) {
		le, err = c.CreateLogentries(&CreateLogentriesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-logentries",
			Port:           0,
			UseTLS:         true,
			Token:          "abcd1234",
			Format:         "format",
			Placement:      "waf_debug",
			Region:         "us",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "logentries/delete", func(c *Client) {
			c.DeleteLogentries(&DeleteLogentriesInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-logentries",
			})

			c.DeleteLogentries(&DeleteLogentriesInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-logentries",
			})
		})
	}()

	if le.Name != "test-logentries" {
		t.Errorf("bad name: %q", le.Name)
	}
	if le.Port != 0 {
		t.Errorf("bad port: %q", le.Port)
	}
	if le.UseTLS != true {
		t.Errorf("bad use_tls: %t", le.UseTLS)
	}
	if le.Token != "abcd1234" {
		t.Errorf("bad token: %q", le.Token)
	}
	if le.Format != "format" {
		t.Errorf("bad format: %q", le.Format)
	}
	if le.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", le.FormatVersion)
	}
	if le.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", le.Placement)
	}
	if le.Region != "us" {
		t.Errorf("bad region: %q", le.Region)
	}

	// List
	var les []*Logentries
	record(t, "logentries/list", func(c *Client) {
		les, err = c.ListLogentries(&ListLogentriesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(les) < 1 {
		t.Errorf("bad logentriess: %v", les)
	}

	// Get
	var nle *Logentries
	record(t, "logentries/get", func(c *Client) {
		nle, err = c.GetLogentries(&GetLogentriesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-logentries",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if le.Name != nle.Name {
		t.Errorf("bad name: %q", le.Name)
	}
	if le.Port != nle.Port {
		t.Errorf("bad port: %q", le.Port)
	}
	if le.UseTLS != nle.UseTLS {
		t.Errorf("bad use_tls: %t", le.UseTLS)
	}
	if le.Token != nle.Token {
		t.Errorf("bad token: %q", le.Token)
	}
	if le.Format != nle.Format {
		t.Errorf("bad format: %q", le.Format)
	}
	if le.FormatVersion != nle.FormatVersion {
		t.Errorf("bad format_version: %q", le.FormatVersion)
	}
	if le.Placement != nle.Placement {
		t.Errorf("bad placement: %q", le.Placement)
	}

	// Update
	var ule *Logentries
	record(t, "logentries/update", func(c *Client) {
		ule, err = c.UpdateLogentries(&UpdateLogentriesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-logentries",
			NewName:        String("new-test-logentries"),
			FormatVersion:  Uint(2),
			Region:         String("ap"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ule.Name != "new-test-logentries" {
		t.Errorf("bad name: %q", ule.Name)
	}
	if ule.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", ule.FormatVersion)
	}
	if ule.Region != "ap" {
		t.Errorf("bad region: %q", ule.Region)
	}

	// Delete
	record(t, "logentries/delete", func(c *Client) {
		err = c.DeleteLogentries(&DeleteLogentriesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-logentries",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListLogentries_validation(t *testing.T) {
	var err error
	_, err = testClient.ListLogentries(&ListLogentriesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListLogentries(&ListLogentriesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateLogentries_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateLogentries(&CreateLogentriesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateLogentries(&CreateLogentriesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetLogentries_validation(t *testing.T) {
	var err error
	_, err = testClient.GetLogentries(&GetLogentriesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetLogentries(&GetLogentriesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetLogentries(&GetLogentriesInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateLogentries_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateLogentries(&UpdateLogentriesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateLogentries(&UpdateLogentriesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateLogentries(&UpdateLogentriesInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteLogentries_validation(t *testing.T) {
	var err error
	err = testClient.DeleteLogentries(&DeleteLogentriesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteLogentries(&DeleteLogentriesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteLogentries(&DeleteLogentriesInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
