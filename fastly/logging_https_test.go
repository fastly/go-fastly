package fastly

import (
	"context"
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
	var httpsCreateResp1, httpsCreateResp2, httpsCreateResp3 *HTTPS
	Record(t, "https/create", func(c *Client) {
		httpsCreateResp1, err = c.CreateHTTPS(context.TODO(), &CreateHTTPSInput{
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
			Period:            ToPointer(5),
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

	Record(t, "https/create2", func(c *Client) {
		httpsCreateResp2, err = c.CreateHTTPS(context.TODO(), &CreateHTTPSInput{
			ServiceID:         TestDeliveryServiceID,
			ServiceVersion:    *tv.Number,
			Name:              ToPointer("test-https-2"),
			Format:            ToPointer("format"),
			URL:               ToPointer("https://example.com/"),
			RequestMaxEntries: ToPointer(1),
			RequestMaxBytes:   ToPointer(1000),
			ContentType:       ToPointer(JSONMimeType),
			HeaderName:        ToPointer("X-Example-Header"),
			HeaderValue:       ToPointer("ExampleValue"),
			Method:            ToPointer(http.MethodPut),
			JSONFormat:        ToPointer("2"),
			Period:            ToPointer(5),
			Placement:         ToPointer("waf_debug"),
			TLSCACert:         ToPointer(caCert),
			TLSClientCert:     ToPointer(clientCert),
			TLSClientKey:      ToPointer(clientKey),
			TLSHostname:       ToPointer("example.com"),
			MessageType:       ToPointer("blank"),
			FormatVersion:     ToPointer(2),
			GzipLevel:         ToPointer(8),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "https/create3", func(c *Client) {
		httpsCreateResp3, err = c.CreateHTTPS(context.TODO(), &CreateHTTPSInput{
			ServiceID:         TestDeliveryServiceID,
			ServiceVersion:    *tv.Number,
			Name:              ToPointer("test-https-3"),
			Format:            ToPointer("format"),
			URL:               ToPointer("https://example.com/"),
			RequestMaxEntries: ToPointer(1),
			RequestMaxBytes:   ToPointer(1000),
			ContentType:       ToPointer(JSONMimeType),
			HeaderName:        ToPointer("X-Example-Header"),
			HeaderValue:       ToPointer("ExampleValue"),
			Method:            ToPointer(http.MethodPut),
			JSONFormat:        ToPointer("2"),
			Period:            ToPointer(5),
			Placement:         ToPointer("waf_debug"),
			TLSCACert:         ToPointer(caCert),
			TLSClientCert:     ToPointer(clientCert),
			TLSClientKey:      ToPointer(clientKey),
			TLSHostname:       ToPointer("example.com"),
			MessageType:       ToPointer("blank"),
			FormatVersion:     ToPointer(2),
			CompressionCodec:  ToPointer("snappy"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail because both CompressionCodec and
	// GzipLevel are present.
	Record(t, "https/create4", func(c *Client) {
		_, err = c.CreateHTTPS(context.TODO(), &CreateHTTPSInput{
			ServiceID:         TestDeliveryServiceID,
			ServiceVersion:    *tv.Number,
			Name:              ToPointer("test-https-4"),
			Format:            ToPointer("format"),
			URL:               ToPointer("https://example.com/"),
			RequestMaxEntries: ToPointer(1),
			RequestMaxBytes:   ToPointer(1000),
			ContentType:       ToPointer(JSONMimeType),
			HeaderName:        ToPointer("X-Example-Header"),
			HeaderValue:       ToPointer("ExampleValue"),
			Method:            ToPointer(http.MethodPut),
			JSONFormat:        ToPointer("2"),
			Period:            ToPointer(5),
			Placement:         ToPointer("waf_debug"),
			TLSCACert:         ToPointer(caCert),
			TLSClientCert:     ToPointer(clientCert),
			TLSClientKey:      ToPointer(clientKey),
			TLSHostname:       ToPointer("example.com"),
			MessageType:       ToPointer("blank"),
			FormatVersion:     ToPointer(2),
			GzipLevel:         ToPointer(8),
			CompressionCodec:  ToPointer("snappy"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// ensure deleted
	defer func() {
		Record(t, "https/cleanup", func(c *Client) {
			_ = c.DeleteHTTPS(context.TODO(), &DeleteHTTPSInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-https",
			})

			// ensure that renamed endpoint created in Update test is deleted
			_ = c.DeleteHTTPS(context.TODO(), &DeleteHTTPSInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-https",
			})

			// ensure that renamed endpoint created in Update test is deleted
			_ = c.DeleteHTTPS(context.TODO(), &DeleteHTTPSInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-https-2",
			})
			// ensure that renamed endpoint created in Update test is deleted
			_ = c.DeleteHTTPS(context.TODO(), &DeleteHTTPSInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-https-3",
			})
		})
	}()

	if *httpsCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", *httpsCreateResp1.Format)
	}
	if *httpsCreateResp1.URL != "https://example.com/" {
		t.Errorf("bad url: %q", *httpsCreateResp1.URL)
	}
	if *httpsCreateResp1.RequestMaxEntries != 1 {
		t.Errorf("bad request_max_entries: %q", *httpsCreateResp1.RequestMaxEntries)
	}
	if *httpsCreateResp1.RequestMaxBytes != 1000 {
		t.Errorf("bad request_max_bytes: %q", *httpsCreateResp1.RequestMaxBytes)
	}
	if *httpsCreateResp1.ContentType != JSONMimeType {
		t.Errorf("bad content_type: %q", *httpsCreateResp1.ContentType)
	}
	if *httpsCreateResp1.HeaderName != "X-Example-Header" {
		t.Errorf("bad *h.ader_name: %q", *httpsCreateResp1.HeaderName)
	}
	if *httpsCreateResp1.HeaderValue != "ExampleValue" {
		t.Errorf("bad *h.ader_value: %q", *httpsCreateResp1.HeaderValue)
	}
	if *httpsCreateResp1.Method != http.MethodPut {
		t.Errorf("bad met*h.d: %q", *httpsCreateResp1.Method)
	}
	if *httpsCreateResp1.JSONFormat != "2" {
		t.Errorf("bad json_format: %q", *httpsCreateResp1.JSONFormat)
	}
	if *httpsCreateResp1.Period != 5 {
		t.Errorf("bad period: %q", *httpsCreateResp1.Period)
	}
	if *httpsCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *httpsCreateResp1.Placement)
	}
	if *httpsCreateResp1.TLSCACert != caCert {
		t.Errorf("bad tls_ca_cert: %q", *httpsCreateResp1.TLSCACert)
	}
	if *httpsCreateResp1.TLSHostname != "example.com" {
		t.Errorf("bad tls_httpsCreateResp1.stname: %q", *httpsCreateResp1.TLSHostname)
	}
	if *httpsCreateResp1.TLSClientCert != clientCert {
		t.Errorf("bad tls_client_cert: %q", *httpsCreateResp1.TLSClientCert)
	}
	if *httpsCreateResp1.TLSClientKey != clientKey {
		t.Errorf("bad tls_client_key: %q", *httpsCreateResp1.TLSClientKey)
	}
	if *httpsCreateResp1.MessageType != "blank" {
		t.Errorf("bad message_type: %s", *httpsCreateResp1.MessageType)
	}
	if *httpsCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", *httpsCreateResp1.FormatVersion)
	}
	if httpsCreateResp2.CompressionCodec != nil {
		t.Errorf("bad compression_codec: %q", *httpsCreateResp2.CompressionCodec)
	}
	if *httpsCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", *httpsCreateResp2.GzipLevel)
	}
	if *httpsCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", *httpsCreateResp3.CompressionCodec)
	}
	if *httpsCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *httpsCreateResp3.GzipLevel)
	}

	// List
	var hs []*HTTPS
	Record(t, "https/list", func(c *Client) {
		hs, err = c.ListHTTPS(context.TODO(), &ListHTTPSInput{
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
		nh, err = c.GetHTTPS(context.TODO(), &GetHTTPSInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-https",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *httpsCreateResp1.Name != *nh.Name {
		t.Errorf("bad name: %q", *httpsCreateResp1.Name)
	}
	if *httpsCreateResp1.Format != *nh.Format {
		t.Errorf("bad format: %q", *httpsCreateResp1.Format)
	}
	if *httpsCreateResp1.URL != *nh.URL {
		t.Errorf("bad url: %q", *httpsCreateResp1.URL)
	}
	if *httpsCreateResp1.RequestMaxEntries != *nh.RequestMaxEntries {
		t.Errorf("bad request_max_entries: %q", *httpsCreateResp1.RequestMaxEntries)
	}
	if *httpsCreateResp1.RequestMaxBytes != *nh.RequestMaxBytes {
		t.Errorf("bad request_max_bytes: %q", *httpsCreateResp1.RequestMaxBytes)
	}
	if *httpsCreateResp1.ContentType != *nh.ContentType {
		t.Errorf("bad content_type: %q", *httpsCreateResp1.ContentType)
	}
	if *httpsCreateResp1.HeaderName != *nh.HeaderName {
		t.Errorf("bad *httpsCreateResp1.ader_name: %q", *httpsCreateResp1.HeaderName)
	}
	if *httpsCreateResp1.HeaderValue != *nh.HeaderValue {
		t.Errorf("bad *httpsCreateResp1.ader_value: %q", *httpsCreateResp1.HeaderValue)
	}
	if *httpsCreateResp1.Method != *nh.Method {
		t.Errorf("bad met*httpsCreateResp1.d: %q", *httpsCreateResp1.Method)
	}
	if *httpsCreateResp1.JSONFormat != *nh.JSONFormat {
		t.Errorf("bad json_format: %q", *httpsCreateResp1.JSONFormat)
	}
	if *httpsCreateResp1.Period != *nh.Period {
		t.Errorf("bad period: %q", *httpsCreateResp1.Period)
	}
	if *httpsCreateResp1.Placement != *nh.Placement {
		t.Errorf("bad placement: %q", *httpsCreateResp1.Placement)
	}
	if *httpsCreateResp1.TLSCACert != *nh.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", *httpsCreateResp1.TLSCACert)
	}
	if *httpsCreateResp1.TLSHostname != *nh.TLSHostname {
		t.Errorf("bad tls_*httpsCreateResp1.stname: %q", *httpsCreateResp1.TLSHostname)
	}
	if *httpsCreateResp1.TLSClientCert != *nh.TLSClientCert {
		t.Errorf("bad tls_client_cert: %q", *httpsCreateResp1.TLSClientCert)
	}
	if *httpsCreateResp1.TLSClientKey != *nh.TLSClientKey {
		t.Errorf("bad tls_client_key: %q", *httpsCreateResp1.TLSClientKey)
	}
	if *httpsCreateResp1.MessageType != *nh.MessageType {
		t.Errorf("bad message_type: %s", *httpsCreateResp1.MessageType)
	}
	if *httpsCreateResp1.FormatVersion != *nh.FormatVersion {
		t.Errorf("bad format_version: %d", *httpsCreateResp1.FormatVersion)
	}

	// Update
	var uh *HTTPS
	Record(t, "https/update", func(c *Client) {
		uh, err = c.UpdateHTTPS(context.TODO(), &UpdateHTTPSInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-https",
			NewName:          ToPointer("new-test-https"),
			Method:           ToPointer(http.MethodPost),
			ProcessingRegion: ToPointer("eu"),
			CompressionCodec: ToPointer("zstd"),
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
	if *uh.ProcessingRegion != "eu" {
		t.Errorf("bad log_processing_region: %q", *uh.ProcessingRegion)
	}
	if *uh.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", *uh.CompressionCodec)
	}

	// Update
	var uh2 *HTTPS
	Record(t, "https/update2", func(c *Client) {
		uh2, err = c.UpdateHTTPS(context.TODO(), &UpdateHTTPSInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-https",
			GzipLevel:      ToPointer(3),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uh2.GzipLevel != 3 {
		t.Errorf("bad gzip_level: %q", *uh2.GzipLevel)
	}

	// Delete
	Record(t, "https/delete", func(c *Client) {
		err = c.DeleteHTTPS(context.TODO(), &DeleteHTTPSInput{
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
	_, err = TestClient.ListHTTPS(context.TODO(), &ListHTTPSInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHTTPS_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateHTTPS(context.TODO(), &CreateHTTPSInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateHTTPS(context.TODO(), &CreateHTTPSInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHTTPS_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetHTTPS(context.TODO(), &GetHTTPSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHTTPS(context.TODO(), &GetHTTPSInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHTTPS(context.TODO(), &GetHTTPSInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHTTPS_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateHTTPS(context.TODO(), &UpdateHTTPSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHTTPS(context.TODO(), &UpdateHTTPSInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHTTPS(context.TODO(), &UpdateHTTPSInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHTTPS_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteHTTPS(context.TODO(), &DeleteHTTPSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHTTPS(context.TODO(), &DeleteHTTPSInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHTTPS(context.TODO(), &DeleteHTTPSInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
