package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_Honeycombs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "honeycombs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var h *Honeycomb
	Record(t, "honeycombs/create", func(c *Client) {
		h, err = c.CreateHoneycomb(context.TODO(), &CreateHoneycombInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-honeycomb"),
			Format:         ToPointer("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:  ToPointer(2),
			Placement:      ToPointer("waf_debug"),
			Token:          ToPointer("super-secure-token"),
			Dataset:        ToPointer("testDataset"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "honeycombs/cleanup", func(c *Client) {
			_ = c.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-honeycomb",
			})

			_ = c.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-honeycomb",
			})
		})
	}()

	if *h.Name != "test-honeycomb" {
		t.Errorf("bad name: %q", *h.Name)
	}
	if *h.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", *h.Format)
	}
	if *h.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *h.FormatVersion)
	}
	if *h.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *h.Placement)
	}
	if *h.Token != "super-secure-token" {
		t.Errorf("bad token: %q", *h.Token)
	}
	if *h.Dataset != "testDataset" {
		t.Errorf("bad dataset: %q", *h.Dataset)
	}

	// List
	var hs []*Honeycomb
	Record(t, "honeycombs/list", func(c *Client) {
		hs, err = c.ListHoneycombs(context.TODO(), &ListHoneycombsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "honeycombs/get", func(c *Client) {
		nh, err = c.GetHoneycomb(context.TODO(), &GetHoneycombInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-honeycomb",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *h.Name != *nh.Name {
		t.Errorf("bad name: %q", *h.Name)
	}
	if *h.Format != *nh.Format {
		t.Errorf("bad format: %q", *h.Format)
	}
	if *h.FormatVersion != *nh.FormatVersion {
		t.Errorf("bad format_version: %q", *h.FormatVersion)
	}
	if *h.Placement != *nh.Placement {
		t.Errorf("bad placement: %q", *h.Placement)
	}
	if *h.Token != *nh.Token {
		t.Errorf("bad token: %q", *h.Token)
	}
	if *h.Dataset != *nh.Dataset {
		t.Errorf("bad dataset: %q", *h.Dataset)
	}

	// Update
	var us *Honeycomb
	Record(t, "honeycombs/update", func(c *Client) {
		us, err = c.UpdateHoneycomb(context.TODO(), &UpdateHoneycombInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-honeycomb",
			NewName:          ToPointer("new-test-honeycomb"),
			Token:            ToPointer("new-token"),
			Dataset:          ToPointer("newDataset"),
			ProcessingRegion: ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Name != "new-test-honeycomb" {
		t.Errorf("bad name: %q", *us.Name)
	}
	if *us.Token != "new-token" {
		t.Errorf("bad token: %q", *us.Token)
	}
	if *us.Dataset != "newDataset" {
		t.Errorf("bad dataset: %q", *us.Dataset)
	}
	if *us.ProcessingRegion != "eu" {
		t.Errorf("bad log_processing_region: %q", *us.ProcessingRegion)
	}

	// Delete
	Record(t, "honeycombs/delete", func(c *Client) {
		err = c.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-honeycomb",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Honeycombs_Compute(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "honeycombs/compute/version", func(c *Client) {
		tv = testVersionCompute(t, c)
	})

	// Create
	var h *Honeycomb
	Record(t, "honeycombs/compute/create", func(c *Client) {
		h, err = c.CreateHoneycomb(context.TODO(), &CreateHoneycombInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-honeycomb"),
			Format:         ToPointer("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:  ToPointer(2),
			Placement:      ToPointer("none"),
			Token:          ToPointer("super-secure-token"),
			Dataset:        ToPointer("testDataset"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "honeycombs/compute/cleanup", func(c *Client) {
			_ = c.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
				ServiceID:      TestComputeServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-honeycomb",
			})

			_ = c.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
				ServiceID:      TestComputeServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-honeycomb",
			})
		})
	}()

	if *h.Name != "test-honeycomb" {
		t.Errorf("bad name: %q", *h.Name)
	}
	if *h.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", *h.Format)
	}
	if *h.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *h.FormatVersion)
	}
	if *h.Placement != "none" {
		t.Errorf("bad placement: %q", *h.Placement)
	}
	if *h.Token != "super-secure-token" {
		t.Errorf("bad token: %q", *h.Token)
	}
	if *h.Dataset != "testDataset" {
		t.Errorf("bad dataset: %q", *h.Dataset)
	}

	// List
	var hs []*Honeycomb
	Record(t, "honeycombs/compute/list", func(c *Client) {
		hs, err = c.ListHoneycombs(context.TODO(), &ListHoneycombsInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "honeycombs/compute/get", func(c *Client) {
		nh, err = c.GetHoneycomb(context.TODO(), &GetHoneycombInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-honeycomb",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *h.Name != *nh.Name {
		t.Errorf("bad name: %q", *h.Name)
	}
	if *h.Format != *nh.Format {
		t.Errorf("bad format: %q", *h.Format)
	}
	if *h.FormatVersion != *nh.FormatVersion {
		t.Errorf("bad format_version: %q", *h.FormatVersion)
	}
	if *h.Placement != *nh.Placement {
		t.Errorf("bad placement: %q", *h.Placement)
	}
	if *h.Token != *nh.Token {
		t.Errorf("bad token: %q", *h.Token)
	}
	if *h.Dataset != *nh.Dataset {
		t.Errorf("bad dataset: %q", *h.Dataset)
	}

	// Update
	var us *Honeycomb
	Record(t, "honeycombs/compute/update", func(c *Client) {
		us, err = c.UpdateHoneycomb(context.TODO(), &UpdateHoneycombInput{
			ServiceID:        TestComputeServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-honeycomb",
			NewName:          ToPointer("new-test-honeycomb"),
			Token:            ToPointer("new-token"),
			Dataset:          ToPointer("newDataset"),
			ProcessingRegion: ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Name != "new-test-honeycomb" {
		t.Errorf("bad name: %q", *us.Name)
	}
	if *us.Token != "new-token" {
		t.Errorf("bad token: %q", *us.Token)
	}
	if *us.Dataset != "newDataset" {
		t.Errorf("bad dataset: %q", *us.Dataset)
	}
	if *us.ProcessingRegion != "eu" {
		t.Errorf("bad log_processing_region: %q", *us.ProcessingRegion)
	}

	// Delete
	Record(t, "honeycombs/compute/delete", func(c *Client) {
		err = c.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-honeycomb",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListHoneycombs_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListHoneycombs(context.TODO(), &ListHoneycombsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListHoneycombs(context.TODO(), &ListHoneycombsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHoneycomb_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateHoneycomb(context.TODO(), &CreateHoneycombInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateHoneycomb(context.TODO(), &CreateHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHoneycomb_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetHoneycomb(context.TODO(), &GetHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHoneycomb(context.TODO(), &GetHoneycombInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHoneycomb(context.TODO(), &GetHoneycombInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHoneycomb_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateHoneycomb(context.TODO(), &UpdateHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHoneycomb(context.TODO(), &UpdateHoneycombInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHoneycomb(context.TODO(), &UpdateHoneycombInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHoneycomb_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHoneycomb(context.TODO(), &DeleteHoneycombInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
