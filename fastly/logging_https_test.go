package fastly

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestClient_HTTPS(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "https/version", func(c *Client) {
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
	var h *HTTPS
	Record(t, "https/create", func(c *Client) {
		h, err = c.CreateHTTPS(&CreateHTTPSInput{
			ServiceID:         TestDeliveryServiceID,
			ServiceVersion:    *tv.Number,
			Name:              ToPointer("test-https"),
			Format:            ToPointer("format"),
			URL:               ToPointer("https://example.com/"),
			RequestMaxEntries: ToPointer(1),
			RequestMaxBytes:   ToPointer(1000),
			ContentType:       ToPointer(JSONMimeType),
			HeaderName:        ToPointer("X-Example-Header"),
			HeaderValue:       ToPointer("ExampleValue"),
			Method:            ToPointer(http.MethodPut),
			JSONFormat:        ToPointer("2"),
			Placement:         ToPointer("waf_debug"),
			TLSCACert:         ToPointer(caCert),
			TLSClientCert:     ToPointer(clientCert),
			TLSClientKey:      ToPointer(clientKey),
			TLSHostname:       ToPointer("example.com"),
			MessageType:       ToPointer("blank"),
			FormatVersion:     ToPointer(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// ensure deleted
	defer func() {
		Record(t, "https/cleanup", func(c *Client) {
			_ = c.DeleteHTTPS(&DeleteHTTPSInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-https",
			})

			// ensure that renamed endpoint created in Update test is deleted
			_ = c.DeleteHTTPS(&DeleteHTTPSInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-https",
			})
		})
	}()

	if *h.Name != "test-https" {
		t.Errorf("bad name: %q", *h.Name)
	}
	if *h.Format != "format" {
		t.Errorf("bad format: %q", *h.Format)
	}
	if *h.URL != "https://example.com/" {
		t.Errorf("bad url: %q", *h.URL)
	}
	if *h.RequestMaxEntries != 1 {
		t.Errorf("bad request_max_entries: %q", *h.RequestMaxEntries)
	}
	if *h.RequestMaxBytes != 1000 {
		t.Errorf("bad request_max_bytes: %q", *h.RequestMaxBytes)
	}
	if *h.ContentType != JSONMimeType {
		t.Errorf("bad content_type: %q", *h.ContentType)
	}
	if *h.HeaderName != "X-Example-Header" {
		t.Errorf("bad *h.ader_name: %q", *h.HeaderName)
	}
	if *h.HeaderValue != "ExampleValue" {
		t.Errorf("bad *h.ader_value: %q", *h.HeaderValue)
	}
	if *h.Method != http.MethodPut {
		t.Errorf("bad met*h.d: %q", *h.Method)
	}
	if *h.JSONFormat != "2" {
		t.Errorf("bad json_format: %q", *h.JSONFormat)
	}
	if *h.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *h.Placement)
	}
	if *h.TLSCACert != caCert {
		t.Errorf("bad tls_ca_cert: %q", *h.TLSCACert)
	}
	if *h.TLSHostname != "example.com" {
		t.Errorf("bad tls_*h.stname: %q", *h.TLSHostname)
	}
	if *h.TLSClientCert != clientCert {
		t.Errorf("bad tls_client_cert: %q", *h.TLSClientCert)
	}
	if *h.TLSClientKey != clientKey {
		t.Errorf("bad tls_client_key: %q", *h.TLSClientKey)
	}
	if *h.MessageType != "blank" {
		t.Errorf("bad message_type: %s", *h.MessageType)
	}
	if *h.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", *h.FormatVersion)
	}

	// List
	var hs []*HTTPS
	Record(t, "https/list", func(c *Client) {
		hs, err = c.ListHTTPS(&ListHTTPSInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hs) < 1 {
		t.Errorf("bad https: %v", hs)
	}

	// Get
	var nh *HTTPS
	Record(t, "https/get", func(c *Client) {
		nh, err = c.GetHTTPS(&GetHTTPSInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-https",
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
	if *h.URL != *nh.URL {
		t.Errorf("bad url: %q", *h.URL)
	}
	if *h.RequestMaxEntries != *nh.RequestMaxEntries {
		t.Errorf("bad request_max_entries: %q", *h.RequestMaxEntries)
	}
	if *h.RequestMaxBytes != *nh.RequestMaxBytes {
		t.Errorf("bad request_max_bytes: %q", *h.RequestMaxBytes)
	}
	if *h.ContentType != *nh.ContentType {
		t.Errorf("bad content_type: %q", *h.ContentType)
	}
	if *h.HeaderName != *nh.HeaderName {
		t.Errorf("bad *h.ader_name: %q", *h.HeaderName)
	}
	if *h.HeaderValue != *nh.HeaderValue {
		t.Errorf("bad *h.ader_value: %q", *h.HeaderValue)
	}
	if *h.Method != *nh.Method {
		t.Errorf("bad met*h.d: %q", *h.Method)
	}
	if *h.JSONFormat != *nh.JSONFormat {
		t.Errorf("bad json_format: %q", *h.JSONFormat)
	}
	if *h.Placement != *nh.Placement {
		t.Errorf("bad placement: %q", *h.Placement)
	}
	if *h.TLSCACert != *nh.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", *h.TLSCACert)
	}
	if *h.TLSHostname != *nh.TLSHostname {
		t.Errorf("bad tls_*h.stname: %q", *h.TLSHostname)
	}
	if *h.TLSClientCert != *nh.TLSClientCert {
		t.Errorf("bad tls_client_cert: %q", *h.TLSClientCert)
	}
	if *h.TLSClientKey != *nh.TLSClientKey {
		t.Errorf("bad tls_client_key: %q", *h.TLSClientKey)
	}
	if *h.MessageType != *nh.MessageType {
		t.Errorf("bad message_type: %s", *h.MessageType)
	}
	if *h.FormatVersion != *nh.FormatVersion {
		t.Errorf("bad format_version: %d", *h.FormatVersion)
	}

	// Update
	var uh *HTTPS
	Record(t, "https/update", func(c *Client) {
		uh, err = c.UpdateHTTPS(&UpdateHTTPSInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-https",
			NewName:        ToPointer("new-test-https"),
			Method:         ToPointer(http.MethodPost),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uh.Name != "new-test-https" {
		t.Errorf("bad name: %q", *uh.Name)
	}
	if *uh.Method != http.MethodPost {
		t.Errorf("bad method: %q", *uh.Method)
	}

	// Delete
	Record(t, "https/delete", func(c *Client) {
		err = c.DeleteHTTPS(&DeleteHTTPSInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-https",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListHTTPS_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListHTTPS(&ListHTTPSInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHTTPS_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateHTTPS(&CreateHTTPSInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateHTTPS(&CreateHTTPSInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHTTPS_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetHTTPS(&GetHTTPSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHTTPS(&GetHTTPSInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHTTPS(&GetHTTPSInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHTTPS_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateHTTPS(&UpdateHTTPSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHTTPS(&UpdateHTTPSInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHTTPS(&UpdateHTTPSInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHTTPS_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteHTTPS(&DeleteHTTPSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHTTPS(&DeleteHTTPSInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHTTPS(&DeleteHTTPSInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
