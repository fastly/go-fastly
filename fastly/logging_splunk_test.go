package fastly

import (
	"errors"
	"strings"
	"testing"
)

func TestClient_Splunks(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "splunks/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	caCert := strings.TrimSpace(`
-----BEGIN CERTIFICATE-----
MIICUTCCAfugAwIBAgIBADANBgkqhkiG9w0BAQQFADBXMQswCQYDVQQGEwJDTjEL
MAkGA1UECBMCUE4xCzAJBgNVBAcTAkNOMQswCQYDVQQKEwJPTjELMAkGA1UECxMC
VU4xFDASBgNVBAMTC0hlcm9uZyBZYW5nMB4XDTA1MDcxNTIxMTk0N1oXDTA1MDgx
NDIxMTk0N1owVzELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAlBOMQswCQYDVQQHEwJD
TjELMAkGA1UEChMCT04xCzAJBgNVBAsTAlVOMRQwEgYDVQQDEwtIZXJvbmcgWWFu
ZzBcMA0GCSqGSIb3DQEBAQUAA0sAMEgCQQCp5hnG7ogBhtlynpOS21cBewKE/B7j
V14qeyslnr26xZUsSVko36ZnhiaO/zbMOoRcKK9vEcgMtcLFuQTWDl3RAgMBAAGj
gbEwga4wHQYDVR0OBBYEFFXI70krXeQDxZgbaCQoR4jUDncEMH8GA1UdIwR4MHaA
FFXI70krXeQDxZgbaCQoR4jUDncEoVukWTBXMQswCQYDVQQGEwJDTjELMAkGA1UE
CBMCUE4xCzAJBgNVBAcTAkNOMQswCQYDVQQKEwJPTjELMAkGA1UECxMCVU4xFDAS
BgNVBAMTC0hlcm9uZyBZYW5nggEAMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEE
BQADQQA/ugzBrjjK9jcWnDVfGHlk3icNRq0oV7Ri32z/+HQX67aRfgZu7KWdI+Ju
Wm7DCfrPNGVwFWUQOmsPue9rZBgO
-----END CERTIFICATE-----
`)
	clientCert := strings.TrimSpace(certificate())
	clientKey := strings.TrimSpace(privateKey())

	// Create
	var s *Splunk
	Record(t, "splunks/create", func(c *Client) {
		s, err = c.CreateSplunk(&CreateSplunkInput{
			ServiceID:         TestDeliveryServiceID,
			ServiceVersion:    *tv.Number,
			Name:              ToPointer("test-splunk"),
			URL:               ToPointer("https://mysplunkendpoint.example.com/services/collector/event"),
			RequestMaxEntries: ToPointer(1),
			RequestMaxBytes:   ToPointer(1000),
			Format:            ToPointer("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:     ToPointer(2),
			Placement:         ToPointer("waf_debug"),
			Token:             ToPointer("super-secure-token"),
			UseTLS:            ToPointer(Compatibool(true)),
			TLSCACert:         ToPointer(caCert),
			TLSHostname:       ToPointer("example.com"),
			TLSClientCert:     ToPointer(clientCert),
			TLSClientKey:      ToPointer(clientKey),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "splunks/cleanup", func(c *Client) {
			_ = c.DeleteSplunk(&DeleteSplunkInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-splunk",
			})

			_ = c.DeleteSplunk(&DeleteSplunkInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-splunk",
			})
		})
	}()

	if *s.Name != "test-splunk" {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.URL != "https://mysplunkendpoint.example.com/services/collector/event" {
		t.Errorf("bad url: %q", *s.URL)
	}
	if *s.RequestMaxEntries != 1 {
		t.Errorf("bad request_max_entries: %q", *s.RequestMaxEntries)
	}
	if *s.RequestMaxBytes != 1000 {
		t.Errorf("bad request_max_bytes: %q", *s.RequestMaxBytes)
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
	if *s.Token != "super-secure-token" {
		t.Errorf("bad token: %q", *s.Token)
	}
	if !*s.UseTLS {
		t.Errorf("bad use_tls: %t", *s.UseTLS)
	}
	if *s.TLSCACert != caCert {
		t.Errorf("bad tls_ca_cert: %q", *s.TLSCACert)
	}
	if *s.TLSHostname != "example.com" {
		t.Errorf("bad tls_hostname: %q", *s.TLSHostname)
	}
	if *s.TLSClientCert != clientCert {
		t.Errorf("bad tls_client_cert: %q", *s.TLSClientCert)
	}
	if *s.TLSClientKey != clientKey {
		t.Errorf("bad tls_client_key: %q", *s.TLSClientKey)
	}

	// List
	var ss []*Splunk
	Record(t, "splunks/list", func(c *Client) {
		ss, err = c.ListSplunks(&ListSplunksInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad splunks: %v", ss)
	}

	// Get
	var ns *Splunk
	Record(t, "splunks/get", func(c *Client) {
		ns, err = c.GetSplunk(&GetSplunkInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-splunk",
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
	if *s.RequestMaxEntries != *ns.RequestMaxEntries {
		t.Errorf("bad request_max_entries: %q", *s.RequestMaxEntries)
	}
	if *s.RequestMaxBytes != *ns.RequestMaxBytes {
		t.Errorf("bad request_max_bytes: %q", *s.RequestMaxBytes)
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
	if *s.Token != *ns.Token {
		t.Errorf("bad token: %q", *s.Token)
	}
	if *s.UseTLS != *ns.UseTLS {
		t.Errorf("bad use_tls: %t", *s.UseTLS)
	}
	if *s.TLSCACert != *ns.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", *s.TLSCACert)
	}
	if *s.TLSHostname != *ns.TLSHostname {
		t.Errorf("bad tls_hostname: %q", *s.TLSHostname)
	}
	if *s.TLSClientCert != *ns.TLSClientCert {
		t.Errorf("bad tls_client_cert: %q", *s.TLSClientCert)
	}
	if *s.TLSClientKey != *ns.TLSClientKey {
		t.Errorf("bad tls_client_key: %q", *s.TLSClientKey)
	}

	// Update
	var us *Splunk
	Record(t, "splunks/update", func(c *Client) {
		us, err = c.UpdateSplunk(&UpdateSplunkInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-splunk",
			NewName:        ToPointer("new-test-splunk"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Name != "new-test-splunk" {
		t.Errorf("bad name: %q", *us.Name)
	}

	// Delete
	Record(t, "splunks/delete", func(c *Client) {
		err = c.DeleteSplunk(&DeleteSplunkInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-splunk",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSplunks_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListSplunks(&ListSplunksInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListSplunks(&ListSplunksInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSplunk_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateSplunk(&CreateSplunkInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateSplunk(&CreateSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSplunk_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetSplunk(&GetSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSplunk(&GetSplunkInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSplunk(&GetSplunkInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSplunk_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateSplunk(&UpdateSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSplunk(&UpdateSplunkInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSplunk(&UpdateSplunkInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSplunk_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteSplunk(&DeleteSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteSplunk(&DeleteSplunkInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteSplunk(&DeleteSplunkInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
