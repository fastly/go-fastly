package fastly

import (
	"errors"
	"testing"
)

func TestClient_Datadog(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "datadog/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var d *Datadog
	Record(t, "datadog/create", func(c *Client) {
		d, err = c.CreateDatadog(&CreateDatadogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-datadog"),
			Region:         ToPointer("US"),
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
		Record(t, "datadog/delete", func(c *Client) {
			_ = c.DeleteDatadog(&DeleteDatadogInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-datadog",
			})

			_ = c.DeleteDatadog(&DeleteDatadogInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-datadog",
			})
		})
	}()

	if *d.Name != "test-datadog" {
		t.Errorf("bad name: %q", *d.Name)
	}
	if *d.Token != "abcd1234" {
		t.Errorf("bad token: %q", *d.Token)
	}
	if *d.Region != "US" {
		t.Errorf("bad token: %q", *d.Region)
	}
	if *d.Format != "format" {
		t.Errorf("bad format: %q", *d.Format)
	}
	if *d.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *d.FormatVersion)
	}
	if *d.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *d.Placement)
	}

	// List
	var ld []*Datadog
	Record(t, "datadog/list", func(c *Client) {
		ld, err = c.ListDatadog(&ListDatadogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ld) < 1 {
		t.Errorf("bad datadogs: %v", ld)
	}

	// Get
	var nd *Datadog
	Record(t, "datadog/get", func(c *Client) {
		nd, err = c.GetDatadog(&GetDatadogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-datadog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *d.Name != *nd.Name {
		t.Errorf("bad name: %q", *d.Name)
	}
	if *d.Token != *nd.Token {
		t.Errorf("bad token: %q", *d.Token)
	}
	if *d.Format != *nd.Format {
		t.Errorf("bad format: %q", *d.Format)
	}
	if *d.FormatVersion != *nd.FormatVersion {
		t.Errorf("bad format_version: %q", *d.FormatVersion)
	}
	if *d.Placement != *nd.Placement {
		t.Errorf("bad placement: %q", *d.Placement)
	}

	// Update
	var ud *Datadog
	Record(t, "datadog/update", func(c *Client) {
		ud, err = c.UpdateDatadog(&UpdateDatadogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-datadog",
			NewName:        ToPointer("new-test-datadog"),
			Region:         ToPointer("EU"),
			FormatVersion:  ToPointer(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ud.Name != "new-test-datadog" {
		t.Errorf("bad name: %q", *ud.Name)
	}
	if *ud.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *ud.FormatVersion)
	}
	if *ud.Region != "EU" {
		t.Errorf("bad region: %q", *ud.Region)
	}

	// Delete
	Record(t, "datadog/delete", func(c *Client) {
		err = c.DeleteDatadog(&DeleteDatadogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-datadog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDatadog_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListDatadog(&ListDatadogInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListDatadog(&ListDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDatadog_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateDatadog(&CreateDatadogInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateDatadog(&CreateDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDatadog_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetDatadog(&GetDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDatadog(&GetDatadogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDatadog(&GetDatadogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDatadog_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateDatadog(&UpdateDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDatadog(&UpdateDatadogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDatadog(&UpdateDatadogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDatadog_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteDatadog(&DeleteDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDatadog(&DeleteDatadogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDatadog(&DeleteDatadogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
