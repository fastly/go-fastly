package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_Logshuttles(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "logshuttles/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var l *Logshuttle
	Record(t, "logshuttles/create", func(c *Client) {
		l, err = c.CreateLogshuttle(context.TODO(), &CreateLogshuttleInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-logshuttle"),
			Format:         ToPointer("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:  ToPointer(2),
			Placement:      ToPointer("waf_debug"),
			Token:          ToPointer("super-secure-token"),
			URL:            ToPointer("https://logs.example.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "logshuttles/cleanup", func(c *Client) {
			_ = c.DeleteLogshuttle(context.TODO(), &DeleteLogshuttleInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-logshuttle",
			})

			_ = c.DeleteLogshuttle(context.TODO(), &DeleteLogshuttleInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-logshuttle",
			})
		})
	}()

	if *l.Name != "test-logshuttle" {
		t.Errorf("bad name: %q", *l.Name)
	}
	if *l.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", *l.Format)
	}
	if *l.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *l.FormatVersion)
	}
	if *l.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *l.Placement)
	}
	if *l.Token != "super-secure-token" {
		t.Errorf("bad token: %q", *l.Token)
	}
	if *l.URL != "https://logs.example.com" {
		t.Errorf("bad url: %q", *l.URL)
	}

	// List
	var ls []*Logshuttle
	Record(t, "logshuttles/list", func(c *Client) {
		ls, err = c.ListLogshuttles(context.TODO(), &ListLogshuttlesInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "logshuttles/get", func(c *Client) {
		nl, err = c.GetLogshuttle(context.TODO(), &GetLogshuttleInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-logshuttle",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *l.Name != *nl.Name {
		t.Errorf("bad name: %q", *l.Name)
	}
	if *l.Format != *nl.Format {
		t.Errorf("bad format: %q", *l.Format)
	}
	if *l.FormatVersion != *nl.FormatVersion {
		t.Errorf("bad format_version: %q", *l.FormatVersion)
	}
	if *l.Placement != *nl.Placement {
		t.Errorf("bad placement: %q", *l.Placement)
	}
	if *l.Token != *nl.Token {
		t.Errorf("bad token: %q", *l.Token)
	}
	if *l.URL != *nl.URL {
		t.Errorf("bad url: %q", *l.URL)
	}

	// Update
	var ul *Logshuttle
	Record(t, "logshuttles/update", func(c *Client) {
		ul, err = c.UpdateLogshuttle(context.TODO(), &UpdateLogshuttleInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-logshuttle",
			NewName:          ToPointer("new-test-logshuttle"),
			Token:            ToPointer("new-token"),
			URL:              ToPointer("https://logs2.example.com"),
			ProcessingRegion: ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ul.Name != "new-test-logshuttle" {
		t.Errorf("bad name: %q", *ul.Name)
	}
	if *ul.Token != "new-token" {
		t.Errorf("bad token: %q", *ul.Token)
	}
	if *ul.URL != "https://logs2.example.com" {
		t.Errorf("bad url: %q", *ul.URL)
	}
	if *ul.ProcessingRegion != "eu" {
		t.Errorf("bad log_processing_region: %q", *ul.ProcessingRegion)
	}

	// Delete
	Record(t, "logshuttles/delete", func(c *Client) {
		err = c.DeleteLogshuttle(context.TODO(), &DeleteLogshuttleInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-logshuttle",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListLogshuttles_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListLogshuttles(context.TODO(), &ListLogshuttlesInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListLogshuttles(context.TODO(), &ListLogshuttlesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateLogshuttle_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateLogshuttle(context.TODO(), &CreateLogshuttleInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateLogshuttle(context.TODO(), &CreateLogshuttleInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetLogshuttle_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetLogshuttle(context.TODO(), &GetLogshuttleInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetLogshuttle(context.TODO(), &GetLogshuttleInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetLogshuttle(context.TODO(), &GetLogshuttleInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateLogshuttle_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateLogshuttle(context.TODO(), &UpdateLogshuttleInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateLogshuttle(context.TODO(), &UpdateLogshuttleInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateLogshuttle(context.TODO(), &UpdateLogshuttleInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteLogshuttle_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteLogshuttle(context.TODO(), &DeleteLogshuttleInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteLogshuttle(context.TODO(), &DeleteLogshuttleInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteLogshuttle(context.TODO(), &DeleteLogshuttleInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
