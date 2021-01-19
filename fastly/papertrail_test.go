package fastly

import (
	"testing"
)

func TestClient_Papertrails(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "papertrails/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var p *Papertrail
	record(t, "papertrails/create", func(c *Client) {
		p, err = c.CreatePapertrail(&CreatePapertrailInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-papertrail",
			Address:        "integ-test.go-fastly.com",
			Port:           1234,
			FormatVersion:  2,
			Format:         "format",
			Placement:      "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "papertrails/cleanup", func(c *Client) {
			c.DeletePapertrail(&DeletePapertrailInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-papertrail",
			})

			c.DeletePapertrail(&DeletePapertrailInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-papertrail",
			})
		})
	}()

	if p.Name != "test-papertrail" {
		t.Errorf("bad name: %q", p.Name)
	}
	if p.Address != "integ-test.go-fastly.com" {
		t.Errorf("bad address: %q", p.Address)
	}
	if p.Port != 1234 {
		t.Errorf("bad port: %q", p.Port)
	}
	if p.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", p.FormatVersion)
	}
	if p.Format != "format" {
		t.Errorf("bad format: %q", p.Format)
	}
	if p.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", p.Placement)
	}

	// List
	var ps []*Papertrail
	record(t, "papertrails/list", func(c *Client) {
		ps, err = c.ListPapertrails(&ListPapertrailsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) < 1 {
		t.Errorf("bad papertrails: %v", ps)
	}

	// Get
	var np *Papertrail
	record(t, "papertrails/get", func(c *Client) {
		np, err = c.GetPapertrail(&GetPapertrailInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-papertrail",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != np.Name {
		t.Errorf("bad name: %q", p.Name)
	}
	if p.Address != np.Address {
		t.Errorf("bad address: %q", p.Address)
	}
	if p.Port != np.Port {
		t.Errorf("bad port: %q", p.Port)
	}
	if p.FormatVersion != np.FormatVersion {
		t.Errorf("bad format_version: %q", p.FormatVersion)
	}
	if p.Format != np.Format {
		t.Errorf("bad format: %q", p.Format)
	}
	if p.Placement != np.Placement {
		t.Errorf("bad placement: %q", p.Placement)
	}

	// Update
	var up *Papertrail
	record(t, "papertrails/update", func(c *Client) {
		up, err = c.UpdatePapertrail(&UpdatePapertrailInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-papertrail",
			NewName:        String("new-test-papertrail"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if up.Name != "new-test-papertrail" {
		t.Errorf("bad name: %q", up.Name)
	}

	// Delete
	record(t, "papertrails/delete", func(c *Client) {
		err = c.DeletePapertrail(&DeletePapertrailInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-papertrail",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPapertrails_validation(t *testing.T) {
	var err error
	_, err = testClient.ListPapertrails(&ListPapertrailsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListPapertrails(&ListPapertrailsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreatePapertrail_validation(t *testing.T) {
	var err error
	_, err = testClient.CreatePapertrail(&CreatePapertrailInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreatePapertrail(&CreatePapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetPapertrail_validation(t *testing.T) {
	var err error
	_, err = testClient.GetPapertrail(&GetPapertrailInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPapertrail(&GetPapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPapertrail(&GetPapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePapertrail_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdatePapertrail(&UpdatePapertrailInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePapertrail(&UpdatePapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePapertrail(&UpdatePapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeletePapertrail_validation(t *testing.T) {
	var err error
	err = testClient.DeletePapertrail(&DeletePapertrailInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePapertrail(&DeletePapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePapertrail(&DeletePapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
