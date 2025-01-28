package fastly

import (
	"errors"
	"strings"
	"testing"
)

func TestClient_Syslogs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "syslogs/version", func(c *Client) {
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
	var s *Syslog
	Record(t, "syslogs/create", func(c *Client) {
		s, err = c.CreateSyslog(&CreateSyslogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-syslog"),
			Address:        ToPointer("example.com"),
			Hostname:       ToPointer("example.com"),
			Port:           ToPointer(1234),
			UseTLS:         ToPointer(Compatibool(true)),
			TLSCACert:      ToPointer(caCert),
			TLSHostname:    ToPointer("example.com"),
			TLSClientCert:  ToPointer(clientCert),
			TLSClientKey:   ToPointer(clientKey),
			Token:          ToPointer("abcd1234"),
			Format:         ToPointer("format"),
			FormatVersion:  ToPointer(2),
			MessageType:    ToPointer("classic"),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "syslogs/cleanup", func(c *Client) {
			_ = c.DeleteSyslog(&DeleteSyslogInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-syslog",
			})

			_ = c.DeleteSyslog(&DeleteSyslogInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-syslog",
			})
		})
	}()

	if *s.Name != "test-syslog" {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.Address != "example.com" {
		t.Errorf("bad address: %q", *s.Address)
	}
	if *s.Hostname != "example.com" {
		t.Errorf("bad hostname: %q", *s.Hostname)
	}
	if *s.Port != 1234 {
		t.Errorf("bad port: %q", *s.Port)
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
	if *s.Token != "abcd1234" {
		t.Errorf("bad token: %q", *s.Token)
	}
	if *s.Format != "format" {
		t.Errorf("bad format: %q", *s.Format)
	}
	if *s.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", *s.FormatVersion)
	}
	if *s.MessageType != "classic" {
		t.Errorf("bad message_type: %s", *s.MessageType)
	}
	if *s.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *s.Placement)
	}

	// List
	var ss []*Syslog
	Record(t, "syslogs/list", func(c *Client) {
		ss, err = c.ListSyslogs(&ListSyslogsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad syslogs: %v", ss)
	}

	// Get
	var ns *Syslog
	Record(t, "syslogs/get", func(c *Client) {
		ns, err = c.GetSyslog(&GetSyslogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-syslog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *s.Name != *ns.Name {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.Address != *ns.Address {
		t.Errorf("bad address: %q", *s.Address)
	}
	if *s.Hostname != *ns.Hostname {
		t.Errorf("bad hostname: %q", *s.Hostname)
	}
	if *s.Port != *ns.Port {
		t.Errorf("bad port: %q", *s.Port)
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
	if *s.Token != *ns.Token {
		t.Errorf("bad token: %q", *s.Token)
	}
	if *s.Format != *ns.Format {
		t.Errorf("bad format: %q", *s.Format)
	}
	if *s.FormatVersion != *ns.FormatVersion {
		t.Errorf("bad format_version: %q", *s.FormatVersion)
	}
	if *s.MessageType != *ns.MessageType {
		t.Errorf("bad message_type: %q", *s.MessageType)
	}
	if *s.Placement != *ns.Placement {
		t.Errorf("bad placement: %q", *s.Placement)
	}

	// Update
	var us *Syslog
	Record(t, "syslogs/update", func(c *Client) {
		us, err = c.UpdateSyslog(&UpdateSyslogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-syslog",
			NewName:        ToPointer("new-test-syslog"),
			FormatVersion:  ToPointer(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Name != "new-test-syslog" {
		t.Errorf("bad name: %q", *us.Name)
	}

	if *us.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", *us.FormatVersion)
	}

	// Delete
	Record(t, "syslogs/delete", func(c *Client) {
		err = c.DeleteSyslog(&DeleteSyslogInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-syslog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSyslogs_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListSyslogs(&ListSyslogsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListSyslogs(&ListSyslogsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSyslog_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateSyslog(&CreateSyslogInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateSyslog(&CreateSyslogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSyslog_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetSyslog(&GetSyslogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSyslog(&GetSyslogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSyslog(&GetSyslogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSyslog_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateSyslog(&UpdateSyslogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSyslog(&UpdateSyslogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSyslog(&UpdateSyslogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSyslog_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteSyslog(&DeleteSyslogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteSyslog(&DeleteSyslogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteSyslog(&DeleteSyslogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
