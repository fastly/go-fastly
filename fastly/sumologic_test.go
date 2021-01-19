package fastly

import (
	"testing"
)

func TestClient_Sumologics(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "sumologics/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var s *Sumologic
	record(t, "sumologics/create", func(c *Client) {
		s, err = c.CreateSumologic(&CreateSumologicInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-sumologic",
			URL:            "https://foo.sumologic.com",
			Format:         "format",
			FormatVersion:  1,
			MessageType:    "classic",
			Placement:      "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "sumologics/cleanup", func(c *Client) {
			c.DeleteSumologic(&DeleteSumologicInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-sumologic",
			})

			c.DeleteSumologic(&DeleteSumologicInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-sumologic",
			})
		})
	}()

	if s.Name != "test-sumologic" {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.URL != "https://foo.sumologic.com" {
		t.Errorf("bad url: %q", s.URL)
	}
	if s.Format != "format" {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != 1 {
		t.Errorf("bad format version: %q", s.FormatVersion)
	}
	if s.MessageType != "classic" {
		t.Errorf("bad message type: %q", s.MessageType)
	}
	if s.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", s.Placement)
	}

	// List
	var ss []*Sumologic
	record(t, "sumologics/list", func(c *Client) {
		ss, err = c.ListSumologics(&ListSumologicsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad sumologics: %v", ss)
	}

	// Get
	var ns *Sumologic
	record(t, "sumologics/get", func(c *Client) {
		ns, err = c.GetSumologic(&GetSumologicInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-sumologic",
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
		t.Errorf("bad format version: %q", s.FormatVersion)
	}
	if s.MessageType != ns.MessageType {
		t.Errorf("bad message type: %q", s.MessageType)
	}
	if s.Placement != ns.Placement {
		t.Errorf("bad placement: %q", s.Placement)
	}

	// Update
	var us *Sumologic
	record(t, "sumologics/update", func(c *Client) {
		us, err = c.UpdateSumologic(&UpdateSumologicInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-sumologic",
			NewName:        String("new-test-sumologic"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-sumologic" {
		t.Errorf("bad name: %q", us.Name)
	}

	// Delete
	record(t, "sumologics/delete", func(c *Client) {
		err = c.DeleteSumologic(&DeleteSumologicInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-sumologic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSumologics_validation(t *testing.T) {
	var err error
	_, err = testClient.ListSumologics(&ListSumologicsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListSumologics(&ListSumologicsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSumologic_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateSumologic(&CreateSumologicInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateSumologic(&CreateSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSumologic_validation(t *testing.T) {
	var err error
	_, err = testClient.GetSumologic(&GetSumologicInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSumologic(&GetSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSumologic(&GetSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSumologic_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateSumologic(&UpdateSumologicInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSumologic(&UpdateSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSumologic(&UpdateSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSumologic_validation(t *testing.T) {
	var err error
	err = testClient.DeleteSumologic(&DeleteSumologicInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSumologic(&DeleteSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSumologic(&DeleteSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
