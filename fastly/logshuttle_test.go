package fastly

import (
	"testing"
)

func TestClient_Logshuttles(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "logshuttles/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var l *Logshuttle
	record(t, "logshuttles/create", func(c *Client) {
		l, err = c.CreateLogshuttle(&CreateLogshuttleInput{
			Service:       testServiceID,
			Version:       tv.Number,
			Name:          String("test-logshuttle"),
			Format:        String("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion: Uint(2),
			Placement:     String("waf_debug"),
			Token:         String("super-secure-token"),
			URL:           String("https://logs.example.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "logshuttles/cleanup", func(c *Client) {
			c.DeleteLogshuttle(&DeleteLogshuttleInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-logshuttle",
			})

			c.DeleteLogshuttle(&DeleteLogshuttleInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-logshuttle",
			})
		})
	}()

	if l.Name != "test-logshuttle" {
		t.Errorf("bad name: %q", l.Name)
	}
	if l.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", l.Format)
	}
	if l.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", l.FormatVersion)
	}
	if l.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", l.Placement)
	}
	if l.Token != "super-secure-token" {
		t.Errorf("bad token: %q", l.Token)
	}
	if l.URL != "https://logs.example.com" {
		t.Errorf("bad url: %q", l.URL)
	}

	// List
	var ls []*Logshuttle
	record(t, "logshuttles/list", func(c *Client) {
		ls, err = c.ListLogshuttles(&ListLogshuttlesInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ls) < 1 {
		t.Errorf("bad logshuttles: %v", ls)
	}

	// Get
	var nl *Logshuttle
	record(t, "logshuttles/get", func(c *Client) {
		nl, err = c.GetLogshuttle(&GetLogshuttleInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-logshuttle",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if l.Name != nl.Name {
		t.Errorf("bad name: %q", l.Name)
	}
	if l.Format != nl.Format {
		t.Errorf("bad format: %q", l.Format)
	}
	if l.FormatVersion != nl.FormatVersion {
		t.Errorf("bad format_version: %q", l.FormatVersion)
	}
	if l.Placement != nl.Placement {
		t.Errorf("bad placement: %q", l.Placement)
	}
	if l.Token != nl.Token {
		t.Errorf("bad token: %q", l.Token)
	}
	if l.URL != nl.URL {
		t.Errorf("bad url: %q", l.URL)
	}

	// Update
	var ul *Logshuttle
	record(t, "logshuttles/update", func(c *Client) {
		ul, err = c.UpdateLogshuttle(&UpdateLogshuttleInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-logshuttle",
			NewName: String("new-test-logshuttle"),
			Token:   String("new-token"),
			URL:     String("https://logs2.example.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ul.Name != "new-test-logshuttle" {
		t.Errorf("bad name: %q", ul.Name)
	}
	if ul.Token != "new-token" {
		t.Errorf("bad token: %q", ul.Token)
	}
	if ul.URL != "https://logs2.example.com" {
		t.Errorf("bad url: %q", ul.URL)
	}

	// Delete
	record(t, "logshuttles/delete", func(c *Client) {
		err = c.DeleteLogshuttle(&DeleteLogshuttleInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-logshuttle",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListLogshuttles_validation(t *testing.T) {
	var err error
	_, err = testClient.ListLogshuttles(&ListLogshuttlesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListLogshuttles(&ListLogshuttlesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateLogshuttle_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateLogshuttle(&CreateLogshuttleInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateLogshuttle(&CreateLogshuttleInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetLogshuttle_validation(t *testing.T) {
	var err error
	_, err = testClient.GetLogshuttle(&GetLogshuttleInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetLogshuttle(&GetLogshuttleInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetLogshuttle(&GetLogshuttleInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateLogshuttle_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateLogshuttle(&UpdateLogshuttleInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateLogshuttle(&UpdateLogshuttleInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateLogshuttle(&UpdateLogshuttleInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteLogshuttle_validation(t *testing.T) {
	var err error
	err = testClient.DeleteLogshuttle(&DeleteLogshuttleInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteLogshuttle(&DeleteLogshuttleInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteLogshuttle(&DeleteLogshuttleInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
