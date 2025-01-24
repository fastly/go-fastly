package fastly

import (
	"errors"
	"testing"
)

func TestClient_Scalyrs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "scalyrs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var s *Scalyr
	Record(t, "scalyrs/create", func(c *Client) {
		s, err = c.CreateScalyr(&CreateScalyrInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-scalyr"),
			Format:         ToPointer("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:  ToPointer(2),
			Placement:      ToPointer("waf_debug"),
			ProjectID:      ToPointer("logplex"),
			Region:         ToPointer("US"),
			Token:          ToPointer("super-secure-token"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "scalyrs/cleanup", func(c *Client) {
			_ = c.DeleteScalyr(&DeleteScalyrInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-scalyr",
			})

			_ = c.DeleteScalyr(&DeleteScalyrInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-scalyr",
			})
		})
	}()

	if *s.Name != "test-scalyr" {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", *s.Format)
	}
	if *s.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *s.FormatVersion)
	}
	if *s.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *s.Placement)
	}
	if *s.ProjectID != "logplex" {
		t.Errorf("bad project_id: %q", *s.Placement)
	}
	if *s.Region != "US" {
		t.Errorf("bad region: %q", *s.Region)
	}
	if *s.Token != "super-secure-token" {
		t.Errorf("bad token: %q", *s.Token)
	}

	// List
	var ss []*Scalyr
	Record(t, "scalyrs/list", func(c *Client) {
		ss, err = c.ListScalyrs(&ListScalyrsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "scalyrs/get", func(c *Client) {
		ns, err = c.GetScalyr(&GetScalyrInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-scalyr",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *s.Name != *ns.Name {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.Format != *ns.Format {
		t.Errorf("bad format: %q", *s.Format)
	}
	if *s.FormatVersion != *ns.FormatVersion {
		t.Errorf("bad format_version: %q", *s.FormatVersion)
	}
	if *s.Placement != *ns.Placement {
		t.Errorf("bad placement: %q", *s.Placement)
	}
	if *s.ProjectID != *ns.ProjectID {
		t.Errorf("bad project_id: %q", *s.ProjectID)
	}
	if *s.Region != "US" {
		t.Errorf("bad region: %q", *s.Region)
	}
	if *s.Token != *ns.Token {
		t.Errorf("bad token: %q", *s.Token)
	}

	// Update
	var us *Scalyr
	Record(t, "scalyrs/update", func(c *Client) {
		us, err = c.UpdateScalyr(&UpdateScalyrInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-scalyr",
			NewName:        ToPointer("new-test-scalyr"),
			ProjectID:      ToPointer("app-name"),
			Region:         ToPointer("EU"),
			Token:          ToPointer("new-token"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Name != "new-test-scalyr" {
		t.Errorf("bad name: %q", *us.Name)
	}
	if *us.ProjectID != "app-name" {
		t.Errorf("bad project_id: %q", *us.ProjectID)
	}
	if *us.Region != "EU" {
		t.Errorf("bad region: %q", *us.Region)
	}
	if *us.Token != "new-token" {
		t.Errorf("bad token: %q", *us.Token)
	}

	// Delete
	Record(t, "scalyrs/delete", func(c *Client) {
		err = c.DeleteScalyr(&DeleteScalyrInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-scalyr",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListScalyrs_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListScalyrs(&ListScalyrsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListScalyrs(&ListScalyrsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateScalyr_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateScalyr(&CreateScalyrInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateScalyr(&CreateScalyrInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetScalyr_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetScalyr(&GetScalyrInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetScalyr(&GetScalyrInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetScalyr(&GetScalyrInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateScalyr_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateScalyr(&UpdateScalyrInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateScalyr(&UpdateScalyrInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateScalyr(&UpdateScalyrInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteScalyr_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteScalyr(&DeleteScalyrInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteScalyr(&DeleteScalyrInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteScalyr(&DeleteScalyrInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
