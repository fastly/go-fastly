package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_Papertrails(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "papertrails/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var p *Papertrail
	Record(t, "papertrails/create", func(c *Client) {
		p, err = c.CreatePapertrail(context.TODO(), &CreatePapertrailInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-papertrail"),
			Address:        ToPointer("integ-test.go-fastly.com"),
			Port:           ToPointer(1234),
			FormatVersion:  ToPointer(2),
			Format:         ToPointer("format"),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "papertrails/cleanup", func(c *Client) {
			_ = c.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-papertrail",
			})

			_ = c.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-papertrail",
			})
		})
	}()

	if *p.Name != "test-papertrail" {
		t.Errorf("bad name: %q", *p.Name)
	}
	if *p.Address != "integ-test.go-fastly.com" {
		t.Errorf("bad address: %q", *p.Address)
	}
	if *p.Port != 1234 {
		t.Errorf("bad port: %q", *p.Port)
	}
	if *p.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *p.FormatVersion)
	}
	if *p.Format != "format" {
		t.Errorf("bad format: %q", *p.Format)
	}
	if *p.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *p.Placement)
	}

	// List
	var ps []*Papertrail
	Record(t, "papertrails/list", func(c *Client) {
		ps, err = c.ListPapertrails(context.TODO(), &ListPapertrailsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) < 1 {
		t.Errorf("bad papertrails: %v", ps)
	}

	// Get
	var np *Papertrail
	Record(t, "papertrails/get", func(c *Client) {
		np, err = c.GetPapertrail(context.TODO(), &GetPapertrailInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-papertrail",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *p.Name != *np.Name {
		t.Errorf("bad name: %q", *p.Name)
	}
	if *p.Address != *np.Address {
		t.Errorf("bad address: %q", *p.Address)
	}
	if *p.Port != *np.Port {
		t.Errorf("bad port: %q", *p.Port)
	}
	if *p.FormatVersion != *np.FormatVersion {
		t.Errorf("bad format_version: %q", *p.FormatVersion)
	}
	if *p.Format != *np.Format {
		t.Errorf("bad format: %q", *p.Format)
	}
	if *p.Placement != *np.Placement {
		t.Errorf("bad placement: %q", *p.Placement)
	}

	// Update
	var up *Papertrail
	Record(t, "papertrails/update", func(c *Client) {
		up, err = c.UpdatePapertrail(context.TODO(), &UpdatePapertrailInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-papertrail",
			NewName:          ToPointer("new-test-papertrail"),
			ProcessingRegion: ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *up.Name != "new-test-papertrail" {
		t.Errorf("bad name: %q", *up.Name)
	}
	if *up.ProcessingRegion != "eu" {
		t.Errorf("bad log_processing_region: %q", *up.ProcessingRegion)
	}

	// Delete
	Record(t, "papertrails/delete", func(c *Client) {
		err = c.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-papertrail",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Papertrails_Compute(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "papertrails/compute/version", func(c *Client) {
		tv = testVersionCompute(t, c)
	})

	// Create
	var p *Papertrail
	Record(t, "papertrails/compute/create", func(c *Client) {
		p, err = c.CreatePapertrail(context.TODO(), &CreatePapertrailInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-papertrail"),
			Address:        ToPointer("integ-test.go-fastly.com"),
			Port:           ToPointer(1234),
			FormatVersion:  ToPointer(2),
			Format:         ToPointer("format"),
			Placement:      ToPointer("none"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "papertrails/compute/cleanup", func(c *Client) {
			_ = c.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
				ServiceID:      TestComputeServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-papertrail",
			})

			_ = c.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
				ServiceID:      TestComputeServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-papertrail",
			})
		})
	}()

	if *p.Name != "test-papertrail" {
		t.Errorf("bad name: %q", *p.Name)
	}
	if *p.Address != "integ-test.go-fastly.com" {
		t.Errorf("bad address: %q", *p.Address)
	}
	if *p.Port != 1234 {
		t.Errorf("bad port: %q", *p.Port)
	}
	if *p.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *p.FormatVersion)
	}
	if *p.Format != "format" {
		t.Errorf("bad format: %q", *p.Format)
	}
	if *p.Placement != "none" {
		t.Errorf("bad placement: %q", *p.Placement)
	}

	// List
	var ps []*Papertrail
	Record(t, "papertrails/compute/list", func(c *Client) {
		ps, err = c.ListPapertrails(context.TODO(), &ListPapertrailsInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) < 1 {
		t.Errorf("bad papertrails: %v", ps)
	}

	// Get
	var np *Papertrail
	Record(t, "papertrails/compute/get", func(c *Client) {
		np, err = c.GetPapertrail(context.TODO(), &GetPapertrailInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-papertrail",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *p.Name != *np.Name {
		t.Errorf("bad name: %q", *p.Name)
	}
	if *p.Address != *np.Address {
		t.Errorf("bad address: %q", *p.Address)
	}
	if *p.Port != *np.Port {
		t.Errorf("bad port: %q", *p.Port)
	}
	if *p.FormatVersion != *np.FormatVersion {
		t.Errorf("bad format_version: %q", *p.FormatVersion)
	}
	if *p.Format != *np.Format {
		t.Errorf("bad format: %q", *p.Format)
	}
	if *p.Placement != *np.Placement {
		t.Errorf("bad placement: %q", *p.Placement)
	}

	// Update
	var up *Papertrail
	Record(t, "papertrails/compute/update", func(c *Client) {
		up, err = c.UpdatePapertrail(context.TODO(), &UpdatePapertrailInput{
			ServiceID:        TestComputeServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-papertrail",
			NewName:          ToPointer("new-test-papertrail"),
			ProcessingRegion: ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *up.Name != "new-test-papertrail" {
		t.Errorf("bad name: %q", *up.Name)
	}
	if *up.ProcessingRegion != "eu" {
		t.Errorf("bad log_processing_region: %q", *up.ProcessingRegion)
	}

	// Delete
	Record(t, "papertrails/compute/delete", func(c *Client) {
		err = c.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-papertrail",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPapertrails_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListPapertrails(context.TODO(), &ListPapertrailsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListPapertrails(context.TODO(), &ListPapertrailsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreatePapertrail_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreatePapertrail(context.TODO(), &CreatePapertrailInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreatePapertrail(context.TODO(), &CreatePapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetPapertrail_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetPapertrail(context.TODO(), &GetPapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetPapertrail(context.TODO(), &GetPapertrailInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetPapertrail(context.TODO(), &GetPapertrailInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePapertrail_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdatePapertrail(context.TODO(), &UpdatePapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdatePapertrail(context.TODO(), &UpdatePapertrailInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdatePapertrail(context.TODO(), &UpdatePapertrailInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeletePapertrail_validation(t *testing.T) {
	var err error

	err = TestClient.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeletePapertrail(context.TODO(), &DeletePapertrailInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
