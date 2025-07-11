package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_Sumologics(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "sumologics/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var s *Sumologic
	Record(t, "sumologics/create", func(c *Client) {
		s, err = c.CreateSumologic(context.TODO(), &CreateSumologicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-sumologic"),
			URL:            ToPointer("https://foo.sumologic.com"),
			Format:         ToPointer("format"),
			FormatVersion:  ToPointer(1),
			MessageType:    ToPointer("classic"),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "sumologics/cleanup", func(c *Client) {
			_ = c.DeleteSumologic(context.TODO(), &DeleteSumologicInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-sumologic",
			})

			_ = c.DeleteSumologic(context.TODO(), &DeleteSumologicInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-sumologic",
			})
		})
	}()

	if *s.Name != "test-sumologic" {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.URL != "https://foo.sumologic.com" {
		t.Errorf("bad url: %q", *s.URL)
	}
	if *s.Format != "format" {
		t.Errorf("bad format: %q", *s.Format)
	}
	if *s.FormatVersion != 1 {
		t.Errorf("bad format version: %q", *s.FormatVersion)
	}
	if *s.MessageType != "classic" {
		t.Errorf("bad message type: %q", *s.MessageType)
	}
	if *s.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *s.Placement)
	}

	// List
	var ss []*Sumologic
	Record(t, "sumologics/list", func(c *Client) {
		ss, err = c.ListSumologics(context.TODO(), &ListSumologicsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "sumologics/get", func(c *Client) {
		ns, err = c.GetSumologic(context.TODO(), &GetSumologicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-sumologic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *s.Name != *ns.Name {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.URL != *ns.URL {
		t.Errorf("bad url: %q", *s.URL)
	}
	if *s.Format != *ns.Format {
		t.Errorf("bad format: %q", *s.Format)
	}
	if *s.FormatVersion != *ns.FormatVersion {
		t.Errorf("bad format version: %q", *s.FormatVersion)
	}
	if *s.MessageType != *ns.MessageType {
		t.Errorf("bad message type: %q", *s.MessageType)
	}
	if *s.Placement != *ns.Placement {
		t.Errorf("bad placement: %q", *s.Placement)
	}

	// Update
	var us *Sumologic
	Record(t, "sumologics/update", func(c *Client) {
		us, err = c.UpdateSumologic(context.TODO(), &UpdateSumologicInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-sumologic",
			NewName:          ToPointer("new-test-sumologic"),
			ProcessingRegion: ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Name != "new-test-sumologic" {
		t.Errorf("bad name: %q", *us.Name)
	}
	if *us.ProcessingRegion != "eu" {
		t.Errorf("bad log_processing_region: %q", *us.ProcessingRegion)
	}

	// Delete
	Record(t, "sumologics/delete", func(c *Client) {
		err = c.DeleteSumologic(context.TODO(), &DeleteSumologicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-sumologic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSumologics_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListSumologics(context.TODO(), &ListSumologicsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListSumologics(context.TODO(), &ListSumologicsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSumologic_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateSumologic(context.TODO(), &CreateSumologicInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateSumologic(context.TODO(), &CreateSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSumologic_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetSumologic(context.TODO(), &GetSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSumologic(context.TODO(), &GetSumologicInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSumologic(context.TODO(), &GetSumologicInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSumologic_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateSumologic(context.TODO(), &UpdateSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSumologic(context.TODO(), &UpdateSumologicInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSumologic(context.TODO(), &UpdateSumologicInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSumologic_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteSumologic(context.TODO(), &DeleteSumologicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteSumologic(context.TODO(), &DeleteSumologicInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteSumologic(context.TODO(), &DeleteSumologicInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
