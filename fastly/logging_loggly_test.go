package fastly

import (
	"errors"
	"testing"
)

func TestClient_Loggly(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "loggly/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var lg *Loggly
	Record(t, "loggly/create", func(c *Client) {
		lg, err = c.CreateLoggly(&CreateLogglyInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-loggly"),
			Token:          ToPointer("abcd1234"),
			Format:         ToPointer("format"),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "loggly/delete", func(c *Client) {
			_ = c.DeleteLoggly(&DeleteLogglyInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-loggly",
			})

			_ = c.DeleteLoggly(&DeleteLogglyInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-loggly",
			})
		})
	}()

	if *lg.Name != "test-loggly" {
		t.Errorf("bad name: %q", *lg.Name)
	}
	if *lg.Token != "abcd1234" {
		t.Errorf("bad token: %q", *lg.Token)
	}
	if *lg.Format != "format" {
		t.Errorf("bad format: %q", *lg.Format)
	}
	if *lg.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *lg.FormatVersion)
	}
	if *lg.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *lg.Placement)
	}

	// List
	var les []*Loggly
	Record(t, "loggly/list", func(c *Client) {
		les, err = c.ListLoggly(&ListLogglyInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(les) < 1 {
		t.Errorf("bad logglys: %v", les)
	}

	// Get
	var nlg *Loggly
	Record(t, "loggly/get", func(c *Client) {
		nlg, err = c.GetLoggly(&GetLogglyInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-loggly",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *lg.Name != *nlg.Name {
		t.Errorf("bad name: %q", *lg.Name)
	}
	if *lg.Token != *nlg.Token {
		t.Errorf("bad token: %q", *lg.Token)
	}
	if *lg.Format != *nlg.Format {
		t.Errorf("bad format: %q", *lg.Format)
	}
	if *lg.FormatVersion != *nlg.FormatVersion {
		t.Errorf("bad format_version: %q", *lg.FormatVersion)
	}
	if *lg.Placement != *nlg.Placement {
		t.Errorf("bad placement: %q", *lg.Placement)
	}

	// Update
	var ulg *Loggly
	Record(t, "loggly/update", func(c *Client) {
		ulg, err = c.UpdateLoggly(&UpdateLogglyInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-loggly",
			NewName:        ToPointer("new-test-loggly"),
			FormatVersion:  ToPointer(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ulg.Name != "new-test-loggly" {
		t.Errorf("bad name: %q", *ulg.Name)
	}
	if *ulg.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *ulg.FormatVersion)
	}

	// Delete
	Record(t, "loggly/delete", func(c *Client) {
		err = c.DeleteLoggly(&DeleteLogglyInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-loggly",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListLoggly_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListLoggly(&ListLogglyInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListLoggly(&ListLogglyInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateLoggly_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateLoggly(&CreateLogglyInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateLoggly(&CreateLogglyInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetLoggly_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetLoggly(&GetLogglyInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetLoggly(&GetLogglyInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetLoggly(&GetLogglyInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateLoggly_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateLoggly(&UpdateLogglyInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateLoggly(&UpdateLogglyInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateLoggly(&UpdateLogglyInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteLoggly_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteLoggly(&DeleteLogglyInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteLoggly(&DeleteLogglyInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteLoggly(&DeleteLogglyInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
