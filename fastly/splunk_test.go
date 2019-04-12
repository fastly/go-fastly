package fastly

import "testing"

func TestClient_Splunks(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "splunks/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var s *Splunk
	record(t, "splunks/create", func(c *Client) {
		s, err = c.CreateSplunk(&CreateSplunkInput{
			Service:       testServiceID,
			Version:       tv.Number,
			Name:          "test-splunk",
			URL:           "https://mysplunkendpoint.example.com/services/collector/event",
			Format:        "%h %l %u %t \"%r\" %>s %b",
			FormatVersion: 2,
			Placement:     "waf_debug",
			Token:         "super-secure-token",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "splunks/cleanup", func(c *Client) {
			c.DeleteSplunk(&DeleteSplunkInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-splunk",
			})

			c.DeleteSplunk(&DeleteSplunkInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-splunk",
			})
		})
	}()

	if s.Name != "test-splunk" {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.URL != "https://mysplunkendpoint.example.com/services/collector/event" {
		t.Errorf("bad url: %q", s.URL)
	}
	if s.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", s.FormatVersion)
	}
	if s.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", s.Placement)
	}
	if s.Token != "super-secure-token" {
		t.Errorf("bad token: %q", s.Token)
	}

	// List
	var ss []*Splunk
	record(t, "splunks/list", func(c *Client) {
		ss, err = c.ListSplunks(&ListSplunksInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad splunks: %v", ss)
	}

	// Get
	var ns *Splunk
	record(t, "splunks/get", func(c *Client) {
		ns, err = c.GetSplunk(&GetSplunkInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-splunk",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != ns.Name {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.URL != ns.URL {
		t.Errorf("bad url: %q", s.URL)
	}
	if s.Format != ns.Format {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != ns.FormatVersion {
		t.Errorf("bad format_version: %q", s.FormatVersion)
	}
	if s.Placement != ns.Placement {
		t.Errorf("bad placement: %q", s.Placement)
	}
	if s.Token != ns.Token {
		t.Errorf("bad token: %q", s.Token)
	}

	// Update
	var us *Splunk
	record(t, "splunks/update", func(c *Client) {
		us, err = c.UpdateSplunk(&UpdateSplunkInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-splunk",
			NewName: "new-test-splunk",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-splunk" {
		t.Errorf("bad name: %q", us.Name)
	}

	// Delete
	record(t, "splunks/delete", func(c *Client) {
		err = c.DeleteSplunk(&DeleteSplunkInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-splunk",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSplunks_validation(t *testing.T) {
	var err error
	_, err = testClient.ListSplunks(&ListSplunksInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListSplunks(&ListSplunksInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSplunk_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateSplunk(&CreateSplunkInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateSplunk(&CreateSplunkInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSplunk_validation(t *testing.T) {
	var err error
	_, err = testClient.GetSplunk(&GetSplunkInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSplunk(&GetSplunkInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSplunk(&GetSplunkInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSplunk_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateSplunk(&UpdateSplunkInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSplunk(&UpdateSplunkInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSplunk(&UpdateSplunkInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSplunk_validation(t *testing.T) {
	var err error
	err = testClient.DeleteSplunk(&DeleteSplunkInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSplunk(&DeleteSplunkInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSplunk(&DeleteSplunkInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
