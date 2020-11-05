package fastly

import (
	"strings"
	"testing"
)

func TestClient_Elasticsearch(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "elasticsearch/version", func(c *Client) {
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
	var es *Elasticsearch
	record(t, "elasticsearch/create", func(c *Client) {
		es, err = c.CreateElasticsearch(&CreateElasticsearchInput{
			ServiceID:         testServiceID,
			ServiceVersion:    tv.Number,
			Name:              String("test-elasticsearch"),
			Format:            String("format"),
			Index:             String("#{%F}"),
			URL:               String("https://example.com/"),
			Pipeline:          String("my_pipeline_id"),
			User:              String("user"),
			Password:          String("password"),
			RequestMaxEntries: Uint(1),
			RequestMaxBytes:   Uint(1000),
			Placement:         String("waf_debug"),
			TLSCACert:         String(caCert),
			TLSClientCert:     String(clientCert),
			TLSClientKey:      String(clientKey),
			TLSHostname:       String("example.com"),
			FormatVersion:     Uint(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// ensure deleted
	defer func() {
		record(t, "elasticsearch/cleanup", func(c *Client) {
			c.DeleteElasticsearch(&DeleteElasticsearchInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-elasticsearch",
			})

			// ensure that renamed endpoint created in Update test is deleted
			c.DeleteElasticsearch(&DeleteElasticsearchInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-elasticsearch",
			})
		})
	}()

	if es.Name != "test-elasticsearch" {
		t.Errorf("bad name: %q", es.Name)
	}
	if es.Format != "format" {
		t.Errorf("bad format: %q", es.Format)
	}
	if es.Index != "#{%F}" {
		t.Errorf("bad index: %q", es.Index)
	}
	if es.URL != "https://example.com/" {
		t.Errorf("bad url: %q", es.URL)
	}
	if es.Pipeline != "my_pipeline_id" {
		t.Errorf("bad pipeline: %q", es.Pipeline)
	}
	if es.User != "user" {
		t.Errorf("bad user: %q", es.User)
	}
	if es.Password != "password" {
		t.Errorf("bad password: %q", es.Password)
	}
	if es.RequestMaxEntries != 1 {
		t.Errorf("bad request_max_entries: %q", es.RequestMaxEntries)
	}
	if es.RequestMaxBytes != 1000 {
		t.Errorf("bad request_max_bytes: %q", es.RequestMaxBytes)
	}
	if es.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", es.Placement)
	}
	if es.TLSCACert != caCert {
		t.Errorf("bad tls_ca_cert: %q", es.TLSCACert)
	}
	if es.TLSHostname != "example.com" {
		t.Errorf("bad tls_hostname: %q", es.TLSHostname)
	}
	if es.TLSClientCert != clientCert {
		t.Errorf("bad tls_client_cert: %q", es.TLSClientCert)
	}
	if es.TLSClientKey != clientKey {
		t.Errorf("bad tls_client_key: %q", es.TLSClientKey)
	}
	if es.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", es.FormatVersion)
	}

	// List
	var ess []*Elasticsearch
	record(t, "elasticsearch/list", func(c *Client) {
		ess, err = c.ListElasticsearch(&ListElasticsearchInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ess) < 1 {
		t.Errorf("bad elasticsearch: %v", ess)
	}

	// Get
	var nes *Elasticsearch
	record(t, "elasticsearch/get", func(c *Client) {
		nes, err = c.GetElasticsearch(&GetElasticsearchInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-elasticsearch",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if es.Name != nes.Name {
		t.Errorf("bad name: %q", es.Name)
	}
	if es.Format != nes.Format {
		t.Errorf("bad format: %q", es.Format)
	}
	if es.Index != nes.Index {
		t.Errorf("bad index: %q", es.Index)
	}
	if es.URL != nes.URL {
		t.Errorf("bad url: %q", es.URL)
	}
	if es.Pipeline != nes.Pipeline {
		t.Errorf("bad pipeline: %q", es.Pipeline)
	}
	if es.User != nes.User {
		t.Errorf("bad user: %q", es.User)
	}
	if es.Password != nes.Password {
		t.Errorf("bad password: %q", es.Password)
	}
	if es.RequestMaxEntries != nes.RequestMaxEntries {
		t.Errorf("bad request_max_entries: %q", es.RequestMaxEntries)
	}
	if es.RequestMaxBytes != nes.RequestMaxBytes {
		t.Errorf("bad request_max_bytes: %q", es.RequestMaxBytes)
	}
	if es.Placement != nes.Placement {
		t.Errorf("bad placement: %q", es.Placement)
	}
	if es.TLSCACert != nes.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", es.TLSCACert)
	}
	if es.TLSHostname != nes.TLSHostname {
		t.Errorf("bad tls_hostname: %q", es.TLSHostname)
	}
	if es.TLSClientCert != nes.TLSClientCert {
		t.Errorf("bad tls_client_cert: %q", es.TLSClientCert)
	}
	if es.TLSClientKey != nes.TLSClientKey {
		t.Errorf("bad tls_client_key: %q", es.TLSClientKey)
	}
	if es.FormatVersion != nes.FormatVersion {
		t.Errorf("bad format_version: %d", es.FormatVersion)
	}

	// Update
	var ues *Elasticsearch
	record(t, "elasticsearch/update", func(c *Client) {
		ues, err = c.UpdateElasticsearch(&UpdateElasticsearchInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-elasticsearch",
			NewName:        String("new-test-elasticsearch"),
			Pipeline:       String("my_new_pipeline_id"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ues.Name != "new-test-elasticsearch" {
		t.Errorf("bad name: %q", ues.Name)
	}
	if ues.Pipeline != "my_new_pipeline_id" {
		t.Errorf("bad pipeline: %q", ues.Pipeline)
	}

	// Delete
	record(t, "elasticsearch/delete", func(c *Client) {
		err = c.DeleteElasticsearch(&DeleteElasticsearchInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-elasticsearch",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListElasticsearch_validation(t *testing.T) {
	var err error
	_, err = testClient.ListElasticsearch(&ListElasticsearchInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListElasticsearch(&ListElasticsearchInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateElasticsearch_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateElasticsearch(&CreateElasticsearchInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateElasticsearch(&CreateElasticsearchInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetElasticsearch_validation(t *testing.T) {
	var err error
	_, err = testClient.GetElasticsearch(&GetElasticsearchInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetElasticsearch(&GetElasticsearchInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetElasticsearch(&GetElasticsearchInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateElasticsearch_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateElasticsearch(&UpdateElasticsearchInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateElasticsearch(&UpdateElasticsearchInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateElasticsearch(&UpdateElasticsearchInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteElasticsearch_validation(t *testing.T) {
	var err error
	err = testClient.DeleteElasticsearch(&DeleteElasticsearchInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteElasticsearch(&DeleteElasticsearchInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteElasticsearch(&DeleteElasticsearchInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
