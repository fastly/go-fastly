package fastly

import (
	"testing"
)

func TestClient_Herokus(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "herokus/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var h *Heroku
	record(t, "herokus/create", func(c *Client) {
		h, err = c.CreateHeroku(&CreateHerokuInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           String("test-heroku"),
			Format:         String("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:  Uint(2),
			Placement:      String("waf_debug"),
			Token:          String("super-secure-token"),
			URL:            String("https://1.us.logplex.io/logs"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "herokus/cleanup", func(c *Client) {
			c.DeleteHeroku(&DeleteHerokuInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-heroku",
			})

			c.DeleteHeroku(&DeleteHerokuInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-heroku",
			})
		})
	}()

	if h.Name != "test-heroku" {
		t.Errorf("bad name: %q", h.Name)
	}
	if h.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", h.Format)
	}
	if h.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", h.FormatVersion)
	}
	if h.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", h.Placement)
	}
	if h.Token != "super-secure-token" {
		t.Errorf("bad token: %q", h.Token)
	}
	if h.URL != "https://1.us.logplex.io/logs" {
		t.Errorf("bad url: %q", h.URL)
	}

	// List
	var hs []*Heroku
	record(t, "herokus/list", func(c *Client) {
		hs, err = c.ListHerokus(&ListHerokusInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hs) < 1 {
		t.Errorf("bad herokus: %v", hs)
	}

	// Get
	var nh *Heroku
	record(t, "herokus/get", func(c *Client) {
		nh, err = c.GetHeroku(&GetHerokuInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-heroku",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if h.Name != nh.Name {
		t.Errorf("bad name: %q", h.Name)
	}
	if h.Format != nh.Format {
		t.Errorf("bad format: %q", h.Format)
	}
	if h.FormatVersion != nh.FormatVersion {
		t.Errorf("bad format_version: %q", h.FormatVersion)
	}
	if h.Placement != nh.Placement {
		t.Errorf("bad placement: %q", h.Placement)
	}
	if h.Token != nh.Token {
		t.Errorf("bad token: %q", h.Token)
	}
	if h.URL != nh.URL {
		t.Errorf("bad url: %q", h.URL)
	}

	// Update
	var uh *Heroku
	record(t, "herokus/update", func(c *Client) {
		uh, err = c.UpdateHeroku(&UpdateHerokuInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-heroku",
			NewName:        String("new-test-heroku"),
			Token:          String("new-token"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uh.Name != "new-test-heroku" {
		t.Errorf("bad name: %q", uh.Name)
	}
	if uh.Token != "new-token" {
		t.Errorf("bad token: %q", uh.Token)
	}

	// Delete
	record(t, "herokus/delete", func(c *Client) {
		err = c.DeleteHeroku(&DeleteHerokuInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-heroku",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListHerokus_validation(t *testing.T) {
	var err error
	_, err = testClient.ListHerokus(&ListHerokusInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListHerokus(&ListHerokusInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHeroku_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateHeroku(&CreateHerokuInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateHeroku(&CreateHerokuInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHeroku_validation(t *testing.T) {
	var err error
	_, err = testClient.GetHeroku(&GetHerokuInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHeroku(&GetHerokuInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHeroku(&GetHerokuInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHeroku_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateHeroku(&UpdateHerokuInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHeroku(&UpdateHerokuInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHeroku(&UpdateHerokuInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHeroku_validation(t *testing.T) {
	var err error
	err = testClient.DeleteHeroku(&DeleteHerokuInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHeroku(&DeleteHerokuInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHeroku(&DeleteHerokuInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
