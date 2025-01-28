package fastly

import (
	"errors"
	"net/http"
	"testing"
)

func TestClient_HealthChecks(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "health_checks/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var hc *HealthCheck
	Record(t, "health_checks/create", func(c *Client) {
		hc, err = c.CreateHealthCheck(&CreateHealthCheckInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-healthcheck"),
			Method:         ToPointer(http.MethodHead),
			Headers: ToPointer([]string{
				"Foo: Bar",
				"Baz: Qux",
			}),
			Host:             ToPointer("example.com"),
			Path:             ToPointer("/foo"),
			HTTPVersion:      ToPointer("1.1"),
			Timeout:          ToPointer(1500),
			CheckInterval:    ToPointer(2500),
			ExpectedResponse: ToPointer(http.StatusOK),
			Window:           ToPointer(5000),
			Threshold:        ToPointer(10),
			Initial:          ToPointer(10),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "health_checks/cleanup", func(c *Client) {
			_ = c.DeleteHealthCheck(&DeleteHealthCheckInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-healthcheck",
			})

			_ = c.DeleteHealthCheck(&DeleteHealthCheckInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-healthcheck",
			})
		})
	}()

	if *hc.Name != "test-healthcheck" {
		t.Errorf("bad name: %q", *hc.Name)
	}
	if *hc.Method != http.MethodHead {
		t.Errorf("bad address: %q", *hc.Method)
	}
	if *hc.Host != "example.com" {
		t.Errorf("bad host: %q", *hc.Host)
	}
	if *hc.Path != "/foo" {
		t.Errorf("bad path: %q", *hc.Path)
	}
	if *hc.HTTPVersion != "1.1" {
		t.Errorf("bad http_version: %q", *hc.HTTPVersion)
	}
	if *hc.Timeout != 1500 {
		t.Errorf("bad timeout: %q", *hc.Timeout)
	}
	if *hc.CheckInterval != 2500 {
		t.Errorf("bad check_interval: %q", *hc.CheckInterval)
	}
	if *hc.ExpectedResponse != http.StatusOK {
		t.Errorf("bad timeout: %q", *hc.ExpectedResponse)
	}
	if *hc.Window != 5000 {
		t.Errorf("bad window: %q", *hc.Window)
	}
	if *hc.Threshold != 10 {
		t.Errorf("bad threshold: %q", *hc.Threshold)
	}
	if *hc.Initial != 10 {
		t.Errorf("bad initial: %q", *hc.Initial)
	}
	if len(hc.Headers) != 2 {
		t.Errorf("bad headers: %q", hc.Headers)
	}

	// List
	var hcs []*HealthCheck
	Record(t, "health_checks/list", func(c *Client) {
		hcs, err = c.ListHealthChecks(&ListHealthChecksInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hcs) < 1 {
		t.Errorf("bad health checks: %v", hcs)
	}

	// Get
	var nhc *HealthCheck
	Record(t, "health_checks/get", func(c *Client) {
		nhc, err = c.GetHealthCheck(&GetHealthCheckInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-healthcheck",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *hc.Name != *nhc.Name {
		t.Errorf("bad name: %q (%q)", *hc.Name, *nhc.Name)
	}
	if *hc.Method != *nhc.Method {
		t.Errorf("bad address: %q", *hc.Method)
	}
	if *hc.Host != *nhc.Host {
		t.Errorf("bad host: %q", *hc.Host)
	}
	if *hc.Path != *nhc.Path {
		t.Errorf("bad path: %q", *hc.Path)
	}
	if *hc.HTTPVersion != *nhc.HTTPVersion {
		t.Errorf("bad http_version: %q", *hc.HTTPVersion)
	}
	if *hc.Timeout != *nhc.Timeout {
		t.Errorf("bad timeout: %q", *hc.Timeout)
	}
	if *hc.CheckInterval != *nhc.CheckInterval {
		t.Errorf("bad check_interval: %q", *hc.CheckInterval)
	}
	if *hc.ExpectedResponse != *nhc.ExpectedResponse {
		t.Errorf("bad timeout: %q", *hc.ExpectedResponse)
	}
	if *hc.Window != *nhc.Window {
		t.Errorf("bad window: %q", *hc.Window)
	}
	if *hc.Threshold != *nhc.Threshold {
		t.Errorf("bad threshold: %q", *hc.Threshold)
	}
	if *hc.Initial != *nhc.Initial {
		t.Errorf("bad initial: %q", *hc.Initial)
	}
	if len(nhc.Headers) != 2 {
		t.Errorf("bad headers: %q", nhc.Headers)
	}
	if hc.Headers[0] != nhc.Headers[0] || hc.Headers[1] != nhc.Headers[1] {
		t.Errorf("bad headers: %q", nhc.Headers)
	}

	// Update
	var uhc *HealthCheck
	Record(t, "health_checks/update", func(c *Client) {
		uhc, err = c.UpdateHealthCheck(&UpdateHealthCheckInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-healthcheck",
			NewName:        ToPointer("new-test-healthcheck"),
			Headers:        ToPointer([]string{"Beep: Boop"}),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uhc.Name != "new-test-healthcheck" {
		t.Errorf("bad update name: %q", *uhc.Name)
	}
	if len(uhc.Headers) != 1 {
		t.Errorf("bad headers: %q", uhc.Headers)
	}

	// Delete
	Record(t, "health_checks/delete", func(c *Client) {
		err = c.DeleteHealthCheck(&DeleteHealthCheckInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-healthcheck",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListHealthChecks_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListHealthChecks(&ListHealthChecksInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListHealthChecks(&ListHealthChecksInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHealthCheck_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateHealthCheck(&CreateHealthCheckInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateHealthCheck(&CreateHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHealthCheck_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetHealthCheck(&GetHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHealthCheck(&GetHealthCheckInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHealthCheck(&GetHealthCheckInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHealthCheck_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateHealthCheck(&UpdateHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHealthCheck(&UpdateHealthCheckInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHealthCheck(&UpdateHealthCheckInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHealthCheck_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteHealthCheck(&DeleteHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHealthCheck(&DeleteHealthCheckInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHealthCheck(&DeleteHealthCheckInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
