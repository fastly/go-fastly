package fastly

import "testing"

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
			Service:   testServiceID,
			Version:   tv.Number,
			Name:      "test-logentries",
			Port:      1234,
			UseTLS:    CBool(true),
			Token:     "abcd1234",
			Format:    "format",
			Placement: "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "logentries/delete", func(c *Client) {
			c.DeleteLogentries(&DeleteLogentriesInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-logentries",
			})

			c.DeleteLogentries(&DeleteLogentriesInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-logentries",
			})
		})
	}()

	if le.Name != "test-logentries" {
		t.Errorf("bad name: %q", le.Name)
	}
	if le.Port != 1234 {
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

	// List
	var les []*Logentries
	record(t, "logentries/list", func(c *Client) {
		les, err = c.ListLogentries(&ListLogentriesInput{
			Service: testServiceID,
			Version: tv.Number,
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
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-logentries",
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
			Service:       testServiceID,
			Version:       tv.Number,
			Name:          "test-logentries",
			NewName:       "new-test-logentries",
			FormatVersion: 2,
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

	// Delete
	record(t, "logentries/delete", func(c *Client) {
		err = c.DeleteLogentries(&DeleteLogentriesInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-logentries",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListLogentries_validation(t *testing.T) {
	var err error
	_, err = testClient.ListLogentries(&ListLogentriesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListLogentries(&ListLogentriesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateLogentries_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateLogentries(&CreateLogentriesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateLogentries(&CreateLogentriesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetLogentries_validation(t *testing.T) {
	var err error
	_, err = testClient.GetLogentries(&GetLogentriesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetLogentries(&GetLogentriesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetLogentries(&GetLogentriesInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateLogentries_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateLogentries(&UpdateLogentriesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateLogentries(&UpdateLogentriesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateLogentries(&UpdateLogentriesInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteLogentries_validation(t *testing.T) {
	var err error
	err = testClient.DeleteLogentries(&DeleteLogentriesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteLogentries(&DeleteLogentriesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteLogentries(&DeleteLogentriesInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
