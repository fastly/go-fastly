package fastly

import "testing"

func TestClient_Backends(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "backends/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var b *Backend
	record(t, "backends/create", func(c *Client) {
		b, err = c.CreateBackend(&CreateBackendInput{
			Service:        testServiceID,
			Version:        tv.Number,
			Name:           "test-backend",
			Address:        "integ-test.go-fastly.com",
			Port:           1234,
			ConnectTimeout: 1500,
			OverrideHost:   "origin.example.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "backends/cleanup", func(c *Client) {
			c.DeleteBackend(&DeleteBackendInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-backend",
			})

			c.DeleteBackend(&DeleteBackendInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-backend",
			})
		})
	}()

	if b.Name != "test-backend" {
		t.Errorf("bad name: %q", b.Name)
	}
	if b.Address != "integ-test.go-fastly.com" {
		t.Errorf("bad address: %q", b.Address)
	}
	if b.Port != 1234 {
		t.Errorf("bad port: %d", b.Port)
	}
	if b.ConnectTimeout != 1500 {
		t.Errorf("bad connect_timeout: %d", b.ConnectTimeout)
	}
	if b.OverrideHost != "origin.example.com" {
		t.Errorf("bad override_host: %q", b.OverrideHost)
	}

	// List
	var bs []*Backend
	record(t, "backends/list", func(c *Client) {
		bs, err = c.ListBackends(&ListBackendsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) < 1 {
		t.Errorf("bad backends: %v", bs)
	}

	// Get
	var nb *Backend
	record(t, "backends/get", func(c *Client) {
		nb, err = c.GetBackend(&GetBackendInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if b.Name != nb.Name {
		t.Errorf("bad name: %q (%q)", b.Name, nb.Name)
	}
	if b.Address != nb.Address {
		t.Errorf("bad address: %q (%q)", b.Address, nb.Address)
	}
	if b.Port != nb.Port {
		t.Errorf("bad port: %q (%q)", b.Port, nb.Port)
	}
	if b.ConnectTimeout != nb.ConnectTimeout {
		t.Errorf("bad connect_timeout: %q (%q)", b.ConnectTimeout, nb.ConnectTimeout)
	}
	if b.OverrideHost != nb.OverrideHost {
		t.Errorf("bad override_host: %q (%q)", b.OverrideHost, nb.OverrideHost)
	}

	// Update
	var ub *Backend
	record(t, "backends/update", func(c *Client) {
		ub, err = c.UpdateBackend(&UpdateBackendInput{
			Service:      testServiceID,
			Version:      tv.Number,
			Name:         "test-backend",
			NewName:      "new-test-backend",
			OverrideHost: "www.example.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ub.Name != "new-test-backend" {
		t.Errorf("bad name: %q", ub.Name)
	}
	if ub.OverrideHost != "www.example.com" {
		t.Errorf("bad override_host: %q", ub.OverrideHost)
	}

	// Delete
	record(t, "backends/delete", func(c *Client) {
		err = c.DeleteBackend(&DeleteBackendInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListBackends_validation(t *testing.T) {
	var err error
	_, err = testClient.ListBackends(&ListBackendsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListBackends(&ListBackendsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateBackend_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateBackend(&CreateBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateBackend(&CreateBackendInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetBackend_validation(t *testing.T) {
	var err error
	_, err = testClient.GetBackend(&GetBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetBackend(&GetBackendInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetBackend(&GetBackendInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateBackend_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateBackend(&UpdateBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateBackend(&UpdateBackendInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateBackend(&UpdateBackendInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteBackend_validation(t *testing.T) {
	var err error
	err = testClient.DeleteBackend(&DeleteBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteBackend(&DeleteBackendInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteBackend(&DeleteBackendInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
