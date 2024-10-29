package fastly

import (
	"testing"
)

func TestClient_GrafanaCloudLogs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "grafanacloudlogs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var d *GrafanaCloudLogs
	record(t, "grafanacloudlogs/create", func(c *Client) {
		d, err = c.CreateGrafanaCloudLogs(&CreateGrafanaCloudLogsInput{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-grafanacloudlogs"),
			URL:            ToPointer("https://test123.grafana.net"),
			User:           ToPointer("123456"),
			Token:          ToPointer("abcd1234"),
			Format:         ToPointer("format"),
			Placement:      ToPointer("waf_debug"),
			Index:          ToPointer("{\"env\": \"prod\"}"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "grafanacloudlogs/delete", func(c *Client) {
			_ = c.DeleteGrafanaCloudLogs(&DeleteGrafanaCloudLogsInput{
				ServiceID:      testServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-grafanacloudlogs",
			})

			_ = c.DeleteGrafanaCloudLogs(&DeleteGrafanaCloudLogsInput{
				ServiceID:      testServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-grafanacloudlogs",
			})
		})
	}()

	if *d.Name != "test-grafanacloudlogs" {
		t.Errorf("bad name: %q", *d.Name)
	}
	if *d.Token != "abcd1234" {
		t.Errorf("bad token: %q", *d.Token)
	}
	if *d.URL != "https://test123.grafana.net" {
		t.Errorf("bad URL: %q", *d.URL)
	}
	if *d.Index != "{\"env\": \"prod\"}" {
		t.Errorf("bad index: %q", *d.Index)
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
	var ld []*GrafanaCloudLogs
	record(t, "grafanacloudlogs/list", func(c *Client) {
		ld, err = c.ListGrafanaCloudLogss(&ListGrafanaCloudLogssInput{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ld) < 1 {
		t.Errorf("bad GrafanaCloudLogss: %v", ld)
	}

	// Get
	var nd *GrafanaCloudLogs
	record(t, "grafanacloudlogs/get", func(c *Client) {
		nd, err = c.GetGrafanaCloudLogs(&GetGrafanaCloudLogsInput{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-grafanacloudlogs",
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
	var ud *GrafanaCloudLogs
	record(t, "grafanacloudlogs/update", func(c *Client) {
		ud, err = c.UpdateGrafanaCloudLogs(&UpdateGrafanaCloudLogsInput{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-grafanacloudlogs",
			NewName:        ToPointer("new-test-grafanacloudlogs"),
			FormatVersion:  ToPointer(2),
			URL:            ToPointer("https://test456.grafana.net"),
			Token:          ToPointer("abcd6789"),
			Placement:      ToPointer("waf_debug"),
			Index:          ToPointer("{\"env\": \"staging\"}"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ud.Name != "new-test-grafanacloudlogs" {
		t.Errorf("bad name: %q", *ud.Name)
	}
	if *ud.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *ud.FormatVersion)
	}
	if *ud.URL != "https://test456.grafana.net" {
		t.Errorf("bad url: %q", *ud.URL)
	}
	if *ud.Index != "{\"env\": \"staging\"}" {
		t.Errorf("bad index: %q", *d.Index)
	}

	// Delete
	record(t, "grafanacloudlogs/delete", func(c *Client) {
		err = c.DeleteGrafanaCloudLogs(&DeleteGrafanaCloudLogsInput{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-grafanacloudlogs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListGrafanaCloudLogs_validation(t *testing.T) {
	var err error

	_, err = testClient.ListGrafanaCloudLogss(&ListGrafanaCloudLogssInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListGrafanaCloudLogss(&ListGrafanaCloudLogssInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateGrafanaCloudLogs_validation(t *testing.T) {
	var err error

	_, err = testClient.CreateGrafanaCloudLogs(&CreateGrafanaCloudLogsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateGrafanaCloudLogs(&CreateGrafanaCloudLogsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetGrafanaCloudLogs_validation(t *testing.T) {
	var err error

	_, err = testClient.GetGrafanaCloudLogs(&GetGrafanaCloudLogsInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetGrafanaCloudLogs(&GetGrafanaCloudLogsInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetGrafanaCloudLogs(&GetGrafanaCloudLogsInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateGrafanaCloudLogs_validation(t *testing.T) {
	var err error

	_, err = testClient.UpdateGrafanaCloudLogs(&UpdateGrafanaCloudLogsInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateGrafanaCloudLogs(&UpdateGrafanaCloudLogsInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateGrafanaCloudLogs(&UpdateGrafanaCloudLogsInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteGrafanaCloudLogs_validation(t *testing.T) {
	var err error

	err = testClient.DeleteGrafanaCloudLogs(&DeleteGrafanaCloudLogsInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteGrafanaCloudLogs(&DeleteGrafanaCloudLogsInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteGrafanaCloudLogs(&DeleteGrafanaCloudLogsInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
