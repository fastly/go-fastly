package fastly

import (
	"testing"
)

func TestClient_Datadog(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "datadog/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var d *Datadog
	record(t, "datadog/create", func(c *Client) {
		d, err = c.CreateDatadog(&CreateDatadogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-datadog",
			Region:         "US",
			Token:          "abcd1234",
			Format:         "format",
			Placement:      "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "datadog/delete", func(c *Client) {
			c.DeleteDatadog(&DeleteDatadogInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-datadog",
			})

			c.DeleteDatadog(&DeleteDatadogInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-datadog",
			})
		})
	}()

	if d.Name != "test-datadog" {
		t.Errorf("bad name: %q", d.Name)
	}
	if d.Token != "abcd1234" {
		t.Errorf("bad token: %q", d.Token)
	}
	if d.Region != "US" {
		t.Errorf("bad token: %q", d.Region)
	}
	if d.Format != "format" {
		t.Errorf("bad format: %q", d.Format)
	}
	if d.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", d.FormatVersion)
	}
	if d.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", d.Placement)
	}

	// List
	var ld []*Datadog
	record(t, "datadog/list", func(c *Client) {
		ld, err = c.ListDatadog(&ListDatadogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
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
	record(t, "datadog/get", func(c *Client) {
		nd, err = c.GetDatadog(&GetDatadogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-datadog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != nd.Name {
		t.Errorf("bad name: %q", d.Name)
	}
	if d.Token != nd.Token {
		t.Errorf("bad token: %q", d.Token)
	}
	if d.Format != nd.Format {
		t.Errorf("bad format: %q", d.Format)
	}
	if d.FormatVersion != nd.FormatVersion {
		t.Errorf("bad format_version: %q", d.FormatVersion)
	}
	if d.Placement != nd.Placement {
		t.Errorf("bad placement: %q", d.Placement)
	}

	// Update
	var ud *Datadog
	record(t, "datadog/update", func(c *Client) {
		ud, err = c.UpdateDatadog(&UpdateDatadogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-datadog",
			NewName:        String("new-test-datadog"),
			Region:         String("EU"),
			FormatVersion:  Uint(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.Name != "new-test-datadog" {
		t.Errorf("bad name: %q", ud.Name)
	}
	if ud.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", ud.FormatVersion)
	}
	if ud.Region != "EU" {
		t.Errorf("bad region: %q", ud.Region)
	}

	// Delete
	record(t, "datadog/delete", func(c *Client) {
		err = c.DeleteDatadog(&DeleteDatadogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-datadog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDatadog_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDatadog(&ListDatadogInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDatadog(&ListDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDatadog_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDatadog(&CreateDatadogInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDatadog(&CreateDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDatadog_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDatadog(&GetDatadogInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDatadog(&GetDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDatadog(&GetDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDatadog_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDatadog(&UpdateDatadogInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDatadog(&UpdateDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDatadog(&UpdateDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDatadog_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDatadog(&DeleteDatadogInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDatadog(&DeleteDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDatadog(&DeleteDatadogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
