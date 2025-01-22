package fastly

import (
	"errors"
	"testing"
)

func TestClient_Backends(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "backends/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var b *Backend
	Record(t, "backends/create", func(c *Client) {
		b, err = c.CreateBackend(&CreateBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-backend"),
			Address:        ToPointer("integ-test.go-fastly.com"),
			ConnectTimeout: ToPointer(1500),
			OverrideHost:   ToPointer("origin.example.com"),
			SSLCheckCert:   ToPointer(Compatibool(false)),
			SSLCiphers:     ToPointer("DHE-RSA-AES256-SHA:DHE-RSA-CAMELLIA256-SHA:AES256-GCM-SHA384"),
			SSLSNIHostname: ToPointer("ssl-hostname.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "backends/cleanup", func(c *Client) {
			_ = c.DeleteBackend(&DeleteBackendInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-backend",
			})

			_ = c.DeleteBackend(&DeleteBackendInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-backend",
			})
		})
	}()

	if *b.Name != "test-backend" {
		t.Errorf("bad name: %q", *b.Name)
	}
	if *b.Address != "integ-test.go-fastly.com" {
		t.Errorf("bad address: %q", *b.Address)
	}
	if *b.Port != 80 {
		t.Errorf("bad port: %d", *b.Port)
	}
	if *b.ConnectTimeout != 1500 {
		t.Errorf("bad connect_timeout: %d", *b.ConnectTimeout)
	}
	if *b.OverrideHost != "origin.example.com" {
		t.Errorf("bad override_host: %q", *b.OverrideHost)
	}
	if b.ShareKey != nil {
		t.Errorf("bad share_key: %s", *b.ShareKey)
	}
	if *b.SSLCheckCert {
		t.Errorf("bad ssl_check_cert: %t", *b.SSLCheckCert) // API defaults to true and we want to allow setting false
	}
	if *b.SSLSNIHostname != "ssl-hostname.com" {
		t.Errorf("bad ssl_sni_hostname: %q", *b.SSLSNIHostname)
	}

	// List
	var bs []*Backend
	Record(t, "backends/list", func(c *Client) {
		bs, err = c.ListBackends(&ListBackendsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "backends/get", func(c *Client) {
		nb, err = c.GetBackend(&GetBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *b.Name != *nb.Name {
		t.Errorf("bad name: %q (%q)", *b.Name, *nb.Name)
	}
	if *b.Address != *nb.Address {
		t.Errorf("bad address: %q (%q)", *b.Address, *nb.Address)
	}
	if *b.Port != *nb.Port {
		t.Errorf("bad port: %q (%q)", *b.Port, *nb.Port)
	}
	if *b.ConnectTimeout != *nb.ConnectTimeout {
		t.Errorf("bad connect_timeout: %q (%q)", *b.ConnectTimeout, *nb.ConnectTimeout)
	}
	if *b.OverrideHost != *nb.OverrideHost {
		t.Errorf("bad override_host: %q (%q)", *b.OverrideHost, *nb.OverrideHost)
	}

	// Update
	var ub *Backend
	Record(t, "backends/update", func(c *Client) {
		ub, err = c.UpdateBackend(&UpdateBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-backend",
			NewName:        ToPointer("new-test-backend"),
			OverrideHost:   ToPointer("www.example.com"),
			Port:           ToPointer(1234),
			ShareKey:       ToPointer("shared-key"),
			SSLCiphers:     ToPointer("RC4:!COMPLEMENTOFDEFAULT"),
			SSLCheckCert:   ToPointer(Compatibool(false)),
			SSLSNIHostname: ToPointer("ssl-hostname-updated.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ub.Name != "new-test-backend" {
		t.Errorf("bad name: %q", *ub.Name)
	}
	if *ub.OverrideHost != "www.example.com" {
		t.Errorf("bad override_host: %q", *ub.OverrideHost)
	}
	if *ub.Port != 1234 {
		t.Errorf("bad port: %d", *ub.Port)
	}
	if *ub.ShareKey == "" || *ub.ShareKey != "shared-key" {
		t.Errorf("bad share_key: %s", *ub.ShareKey)
	}
	if *ub.SSLCheckCert {
		t.Errorf("bad ssl_check_cert: %t", *ub.SSLCheckCert)
	}
	if *ub.SSLSNIHostname != "ssl-hostname-updated.com" {
		t.Errorf("bad ssl_sni_hostname: %q", *ub.SSLSNIHostname)
	}

	// NOTE: The following test validates empty values are NOT sent.
	Record(t, "backends/update_ignore_empty_values", func(c *Client) {
		ub, err = c.UpdateBackend(&UpdateBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ub.OverrideHost != "www.example.com" {
		t.Errorf("bad override_host: %q", *ub.OverrideHost)
	}
	if *ub.Port != 1234 {
		t.Errorf("bad port: %d", *ub.Port)
	}

	// NOTE: The following test validates empty values ARE sent.
	//
	// e.g. Although OverrideHost and Port are set to the type's zero values, and
	// the UpdateBackendInput struct fields set omitempty, they're pointer types
	// and so the JSON unmarshal recognises that empty values are allowed.
	Record(t, "backends/update_allow_empty_values", func(c *Client) {
		ub, err = c.UpdateBackend(&UpdateBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-backend",
			OverrideHost:   ToPointer(""),
			Port:           ToPointer(0),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ub.OverrideHost != nil {
		t.Errorf("bad override_host: %q", *ub.OverrideHost)
	}
	if *ub.Port != 0 {
		t.Errorf("bad port: %d", *ub.Port)
	}

	// Delete
	Record(t, "backends/delete", func(c *Client) {
		err = c.DeleteBackend(&DeleteBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListBackends_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListBackends(&ListBackendsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListBackends(&ListBackendsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateBackend_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateBackend(&CreateBackendInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateBackend(&CreateBackendInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetBackend_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetBackend(&GetBackendInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetBackend(&GetBackendInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetBackend(&GetBackendInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateBackend_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateBackend(&UpdateBackendInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateBackend(&UpdateBackendInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateBackend(&UpdateBackendInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteBackend_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteBackend(&DeleteBackendInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteBackend(&DeleteBackendInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteBackend(&DeleteBackendInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
