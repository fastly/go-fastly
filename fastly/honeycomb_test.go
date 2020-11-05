package fastly

import (
	"testing"
)

func TestClient_Honeycombs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "honeycombs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var h *Honeycomb
	record(t, "honeycombs/create", func(c *Client) {
		h, err = c.CreateHoneycomb(&CreateHoneycombInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           String("test-honeycomb"),
			Format:         String("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:  Uint(2),
			Placement:      String("waf_debug"),
			Token:          String("super-secure-token"),
			Dataset:        String("testDataset"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "honeycombs/cleanup", func(c *Client) {
			c.DeleteHoneycomb(&DeleteHoneycombInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-honeycomb",
			})

			c.DeleteHoneycomb(&DeleteHoneycombInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-honeycomb",
			})
		})
	}()

	if h.Name != "test-honeycomb" {
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
	if h.Dataset != "testDataset" {
		t.Errorf("bad dataset: %q", h.Dataset)
	}

	// List
	var hs []*Honeycomb
	record(t, "honeycombs/list", func(c *Client) {
		hs, err = c.ListHoneycombs(&ListHoneycombsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hs) < 1 {
		t.Errorf("bad honeycombs: %v", hs)
	}

	// Get
	var nh *Honeycomb
	record(t, "honeycombs/get", func(c *Client) {
		nh, err = c.GetHoneycomb(&GetHoneycombInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-honeycomb",
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
	if h.Dataset != nh.Dataset {
		t.Errorf("bad dataset: %q", h.Dataset)
	}

	// Update
	var us *Honeycomb
	record(t, "honeycombs/update", func(c *Client) {
		us, err = c.UpdateHoneycomb(&UpdateHoneycombInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-honeycomb",
			NewName:        String("new-test-honeycomb"),
			Token:          String("new-token"),
			Dataset:        String("newDataset"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-honeycomb" {
		t.Errorf("bad name: %q", us.Name)
	}
	if us.Token != "new-token" {
		t.Errorf("bad token: %q", us.Token)
	}
	if us.Dataset != "newDataset" {
		t.Errorf("bad dataset: %q", us.Dataset)
	}

	// Delete
	record(t, "honeycombs/delete", func(c *Client) {
		err = c.DeleteHoneycomb(&DeleteHoneycombInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-honeycomb",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListHoneycombs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListHoneycombs(&ListHoneycombsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListHoneycombs(&ListHoneycombsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHoneycomb_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateHoneycomb(&CreateHoneycombInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateHoneycomb(&CreateHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHoneycomb_validation(t *testing.T) {
	var err error
	_, err = testClient.GetHoneycomb(&GetHoneycombInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHoneycomb(&GetHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHoneycomb(&GetHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHoneycomb_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateHoneycomb(&UpdateHoneycombInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHoneycomb(&UpdateHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHoneycomb(&UpdateHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHoneycomb_validation(t *testing.T) {
	var err error
	err = testClient.DeleteHoneycomb(&DeleteHoneycombInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHoneycomb(&DeleteHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHoneycomb(&DeleteHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
