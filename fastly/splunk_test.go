package fastly

import (
	"strings"
	"testing"
)

func TestClient_Splunks(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "splunks/version", func(c *Client) {
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

	// Create
	var s *Splunk
	record(t, "splunks/create", func(c *Client) {
		s, err = c.CreateSplunk(&CreateSplunkInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           String("test-splunk"),
			URL:            String("https://mysplunkendpoint.example.com/services/collector/event"),
			Format:         String("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:  Uint(2),
			Placement:      String("waf_debug"),
			Token:          String("super-secure-token"),
			TLSCACert:      String(caCert),
			TLSHostname:    String("example.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "splunks/cleanup", func(c *Client) {
			c.DeleteSplunk(&DeleteSplunkInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-splunk",
			})

			c.DeleteSplunk(&DeleteSplunkInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-splunk",
			})
		})
	}()

	if s.Name != "test-splunk" {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.URL != "https://mysplunkendpoint.example.com/services/collector/event" {
		t.Errorf("bad url: %q", s.URL)
	}
	if s.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", s.FormatVersion)
	}
	if s.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", s.Placement)
	}
	if s.Token != "super-secure-token" {
		t.Errorf("bad token: %q", s.Token)
	}
	if s.TLSCACert != caCert {
		t.Errorf("bad tls_ca_cert: %q", s.TLSCACert)
	}
	if s.TLSHostname != "example.com" {
		t.Errorf("bad tls_hostname: %q", s.TLSHostname)
	}

	// List
	var ss []*Splunk
	record(t, "splunks/list", func(c *Client) {
		ss, err = c.ListSplunks(&ListSplunksInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
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
	record(t, "splunks/get", func(c *Client) {
		ns, err = c.GetSplunk(&GetSplunkInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-splunk",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != ns.Name {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.URL != ns.URL {
		t.Errorf("bad url: %q", s.URL)
	}
	if s.Format != ns.Format {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != ns.FormatVersion {
		t.Errorf("bad format_version: %q", s.FormatVersion)
	}
	if s.Placement != ns.Placement {
		t.Errorf("bad placement: %q", s.Placement)
	}
	if s.Token != ns.Token {
		t.Errorf("bad token: %q", s.Token)
	}
	if s.TLSCACert != ns.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", s.TLSCACert)
	}
	if s.TLSHostname != ns.TLSHostname {
		t.Errorf("bad tls_hostname: %q", s.TLSHostname)
	}

	// Update
	var us *Splunk
	record(t, "splunks/update", func(c *Client) {
		us, err = c.UpdateSplunk(&UpdateSplunkInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-splunk",
			NewName:        String("new-test-splunk"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-splunk" {
		t.Errorf("bad name: %q", us.Name)
	}

	// Delete
	record(t, "splunks/delete", func(c *Client) {
		err = c.DeleteSplunk(&DeleteSplunkInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-splunk",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSplunks_validation(t *testing.T) {
	var err error
	_, err = testClient.ListSplunks(&ListSplunksInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListSplunks(&ListSplunksInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSplunk_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateSplunk(&CreateSplunkInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateSplunk(&CreateSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSplunk_validation(t *testing.T) {
	var err error
	_, err = testClient.GetSplunk(&GetSplunkInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSplunk(&GetSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSplunk(&GetSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSplunk_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateSplunk(&UpdateSplunkInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSplunk(&UpdateSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSplunk(&UpdateSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSplunk_validation(t *testing.T) {
	var err error
	err = testClient.DeleteSplunk(&DeleteSplunkInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSplunk(&DeleteSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSplunk(&DeleteSplunkInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
