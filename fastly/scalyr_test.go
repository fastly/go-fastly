package fastly

import (
	"testing"
)

func TestClient_Scalyrs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "scalyrs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var s *Scalyr
	record(t, "scalyrs/create", func(c *Client) {
		s, err = c.CreateScalyr(&CreateScalyrInput{
			Service:       testServiceID,
			Version:       tv.Number,
			Name:          String("test-scalyr"),
			Format:        String("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion: Uint(2),
			Placement:     String("waf_debug"),
			Token:         String("super-secure-token"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "scalyrs/cleanup", func(c *Client) {
			c.DeleteScalyr(&DeleteScalyrInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-scalyr",
			})

			c.DeleteScalyr(&DeleteScalyrInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-scalyr",
			})
		})
	}()

	if s.Name != "test-scalyr" {
		t.Errorf("bad name: %q", s.Name)
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
	var ss []*Scalyr
	record(t, "scalyrs/list", func(c *Client) {
		ss, err = c.ListScalyrs(&ListScalyrsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad scalyrs: %v", ss)
	}

	// Get
	var ns *Scalyr
	record(t, "scalyrs/get", func(c *Client) {
		ns, err = c.GetScalyr(&GetScalyrInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-scalyr",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != ns.Name {
		t.Errorf("bad name: %q", s.Name)
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
	var us *Scalyr
	record(t, "scalyrs/update", func(c *Client) {
		us, err = c.UpdateScalyr(&UpdateScalyrInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-scalyr",
			NewName: String("new-test-scalyr"),
			Token:   String("new-token"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-scalyr" {
		t.Errorf("bad name: %q", us.Name)
	}
	if us.Token != "new-token" {
		t.Errorf("bad token: %q", us.Token)
	}

	// Delete
	record(t, "scalyrs/delete", func(c *Client) {
		err = c.DeleteScalyr(&DeleteScalyrInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-scalyr",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListScalyrs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListScalyrs(&ListScalyrsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListScalyrs(&ListScalyrsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateScalyr_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateScalyr(&CreateScalyrInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateScalyr(&CreateScalyrInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetScalyr_validation(t *testing.T) {
	var err error
	_, err = testClient.GetScalyr(&GetScalyrInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetScalyr(&GetScalyrInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetScalyr(&GetScalyrInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateScalyr_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateScalyr(&UpdateScalyrInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateScalyr(&UpdateScalyrInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateScalyr(&UpdateScalyrInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteScalyr_validation(t *testing.T) {
	var err error
	err = testClient.DeleteScalyr(&DeleteScalyrInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteScalyr(&DeleteScalyrInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteScalyr(&DeleteScalyrInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
