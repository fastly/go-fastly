package fastly

import (
	"strings"
	"testing"
)

func TestClient_Syslogs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "syslogs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	ca_cert := strings.TrimSpace(`
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
	client_cert := strings.TrimSpace(`
-----BEGIN CERTIFICATE-----
MIICwjCCAaoCCQDpxERwAV8huTANBgkqhkiG9w0BAQsFADAjMQswCQYDVQQGEwJH
QjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wHhcNMTkxMjA1MTYxNjMzWhcNMjAwMTA0
MTYxNjMzWjAjMQswCQYDVQQGEwJHQjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDiQFmq5jJbGSXdRx7wF2Efj8Oc
ila+i0Q6FCoGFXGLV7nSJeY05FaBa+zrA/qQxASAg5pa+v/noOrC6Dq2Qmkdi3AM
NT4UtaWEIQ7IVpj+DjYsVp9R3l3sHZQGtflZOu6lwVmOJ+11FWszgiLTD1vrqzwl
2vu8mSpUPFGpECB0kyRB7/pBlf5cjHdlCiD4xh++2+aJmwrgXq/MTr65XwX6eeul
3+/KALhiDqajwMIYKi6fuAZ308yggpxVJCx2zNMqDpYQVS2sWaATwWngTTCq5t0O
oVliAbU5KHvdUZ/rwmd5kzRuznK79tSH5cMBQ11bRfRrshnaHTJFRbFDUKBdAgMB
AAEwDQYJKoZIhvcNAQELBQADggEBAMOJlf3QocN58Ni9VsCYAb17OXdsg5oDOPZU
b6TvYyEhF69rntH5Z05BvadWmFWixsVConPWSmIt9tLZQB2kCK/KbhC0697GKqy9
0FIbGukPdpSwu6nsdC7Pr+kDn8N5Hgxp6on61wC4hINmNYuu0HLMUpZNg9lIfsfZ
Hu6azlnfvbiIpb2ysM+W8TgW7nJtHv91OsjAqUlGX23L/nFxgobwmvNGXFxN+GBJ
jZt5HMD3u8LvwKjvUW5cGvf0eYPESi4CQOXB1KSPZ3n0toJVUKcLBTymd3VHrMNs
mcgsAByKMLeR0YsFh4pzOg9ji5KCijjfTf0VjBkoP9gLBRkgns4=
-----END CERTIFICATE-----
	
`)
	client_key := strings.TrimSpace(`
-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIFHzBJBgkqhkiG9w0BBQ0wPDAbBgkqhkiG9w0BBQwwDgQIaNbyz5arAFECAggA
MB0GCWCGSAFlAwQBKgQQwrD9J0Ja0Y92VCOWL/cdtgSCBNCMospVmhwttJ1Ad2Kx
lHqhfVVlN6dZw/QLShabtWkfzvVmSNdpaMkObcCpK58czmG+FOs9pPoxJWSyRHPF
GU5EnDd4TAIoIq7WkMiy0KIoM+k9KGL7pZ3YWh9j9B884H5d7iobKjxZjU1zHiKc
Spw4EpL1A5iTF4yP3RftoXbMXv4K0URGqZyvj/i+Dt/fYwmj28CU85VD83YifjFq
K2necIchFLzajuHmTpZfOZ6uDthIv8guE4oBEI4gRrpy0cvTEBm3oltV/ZfdsFGs
nX/WuVORlAUVIPV5OF3r/APTBdUuBC3ZsoH/EE6ySOrLl2hLvNU2iNBKLlSbsq0J
sattoPpQLenS6tUUvDAqtViaUMew1MagB4tdqM+rKw0lw7MMSCrxuhDhMX1ZTgqK
xIR9fCPTUssN01Bj8YC3lzOwmiY0xhG7kjnAMcfyky9SJfBh4/t7OasZASfEkmat
w/prYUEMPsFge8BPuhACpcsb6SwtVLYs17YVaU+OVwHIq1CEXPnw9O+HRHqNAP0C
zS9WcKQttVGWuxCgDV8CxJZSuuNeHFTTzVP/YUpNM20FBnFzG7j+eFGv0un7oOtm
Lo8tvWYITEwXqyJs5i2kYV8Gu1abGCOkwhFqQTNyas1gUUdk3zTymcklIpTntiRZ
vAYXgJoClk9fhvr4le7iIg0fFmocCfmzLO5uqQKny0ce9KDeWP5PTys+QIcQTNMw
CffmlH6HFmOBcehPvv5Ud6MrHKfIBCtHleUGOAAYaoGkpSEZodZCo2Kwzn8CKDJg
qKqjyZbsaC8Hv4AjzT+nJ5WmhTQhMaF/cngeQKY+3hQZJm11aNsXiNl+iFX71X/T
aPHom1mgXqXBWZktjLXYJwM5f8J2uA3EuxvkSXOp1giEt9XuioW6/Ponf9GVXdTI
iLHeFS98K8jSRWt7HAgFAuTtthFjmbIl3Pv6zGwKbC0b6OiGVKzodeRczXnT3UK/
g2KdeUy8xOF7CnCEchXvn4t6SzLr7MbmePZ5Bcjvnlr19WLvejmKDyojsGrsPoVG
gv4zA3tDsbMpBc1gyHQ2Q82+9RxhtXifSPdCq9YaFio9LXbUX8ORoDJ7OjaYLlWI
v2dStjaI4l2mvamS6+VXatyTpCelMOvlpL7UKOScQ8LIuFKw5TWYm5lUWVwm1tka
hldB9bXPV07/y8zYdKxNmNnFg0f3jF3bd/kmucU5TabxLOCbCH8zPlx+EJ/y6Y28
SnzvJQag/p32W9Es7ZO8tQ//m3CpHx8OONDIu5zuyHs0MKhnHYPe3viJWgaGDn94
wk0J/UjQso2X3yBOwK4tFP8nUU/5HX3khTL9vpBNu/7f7fxePHPDxUUW4baOn6Mx
dDFS9Bi0YMEW0mf57RPkGID4Aam8OEjE5aGJ/2qZtbsWeKKQowoER/vSni576KD7
Kxkz6Jv73/tSGmvfJCLUzVTLzDh2gW7jx9710RoqTXlzuqhiv1mtdq87m6GNpxAR
uYy4NNy0STH8Bmm0t4meRKcGx899nj5DKdZ0zcBST6tmBgFDbC+yXk5XnudbuMyn
exp+TUodr1uA292KuOVW0aekVYjTNLxlkawHyvxrhO4rhvpjfx6Rw/LZfQmV43WO
04gxLgnMwQDFSeh5LOhxpbcizg==
-----END ENCRYPTED PRIVATE KEY-----
`)

	// Create
	var s *Syslog
	record(t, "syslogs/create", func(c *Client) {
		s, err = c.CreateSyslog(&CreateSyslogInput{
			Service:       testServiceID,
			Version:       tv.Number,
			Name:          "test-syslog",
			Address:       "example.com",
			Hostname:      "example.com",
			Port:          1234,
			UseTLS:        CBool(true),
			TLSCACert:     ca_cert,
			TLSHostname:   "example.com",
			TLSClientCert: client_cert,
			TLSClientKey:  client_key,
			Token:         "abcd1234",
			Format:        "format",
			FormatVersion: 2,
			MessageType:   "classic",
			Placement:     "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "syslogs/cleanup", func(c *Client) {
			c.DeleteSyslog(&DeleteSyslogInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-syslog",
			})

			c.DeleteSyslog(&DeleteSyslogInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-syslog",
			})
		})
	}()

	if s.Name != "test-syslog" {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.Address != "example.com" {
		t.Errorf("bad address: %q", s.Address)
	}
	if s.Hostname != "example.com" {
		t.Errorf("bad hostname: %q", s.Hostname)
	}
	if s.Port != 1234 {
		t.Errorf("bad port: %q", s.Port)
	}
	if s.UseTLS != true {
		t.Errorf("bad use_tls: %t", s.UseTLS)
	}
	if s.TLSCACert != ca_cert {
		t.Errorf("bad tls_ca_cert: %q", s.TLSCACert)
	}
	if s.TLSHostname != "example.com" {
		t.Errorf("bad tls_hostname: %q", s.TLSHostname)
	}
	if s.TLSClientCert != client_cert {
		t.Errorf("bad tls_client_cert: %q", s.TLSClientCert)
	}
	if s.TLSClientKey != client_key {
		t.Errorf("bad tls_client_key: %q", s.TLSClientKey)
	}
	if s.Token != "abcd1234" {
		t.Errorf("bad token: %q", s.Token)
	}
	if s.Format != "format" {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", s.FormatVersion)
	}
	if s.MessageType != "classic" {
		t.Errorf("bad message_type: %s", s.MessageType)
	}
	if s.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", s.Placement)
	}

	// List
	var ss []*Syslog
	record(t, "syslogs/list", func(c *Client) {
		ss, err = c.ListSyslogs(&ListSyslogsInput{
			Service: testServiceID,
			Version: tv.Number,
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
	record(t, "syslogs/get", func(c *Client) {
		ns, err = c.GetSyslog(&GetSyslogInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-syslog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != ns.Name {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.Address != ns.Address {
		t.Errorf("bad address: %q", s.Address)
	}
	if s.Hostname != ns.Hostname {
		t.Errorf("bad hostname: %q", s.Hostname)
	}
	if s.Port != ns.Port {
		t.Errorf("bad port: %q", s.Port)
	}
	if s.UseTLS != ns.UseTLS {
		t.Errorf("bad use_tls: %t", s.UseTLS)
	}
	if s.TLSCACert != ns.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", s.TLSCACert)
	}
	if s.TLSHostname != ns.TLSHostname {
		t.Errorf("bad tls_hostname: %q", s.TLSHostname)
	}
	if s.TLSClientCert != ns.TLSClientCert {
		t.Errorf("bad tls_client_cert: %q", s.TLSClientCert)
	}
	if s.TLSClientKey != ns.TLSClientKey {
		t.Errorf("bad tls_client_key: %q", s.TLSClientKey)
	}
	if s.Token != ns.Token {
		t.Errorf("bad token: %q", s.Token)
	}
	if s.Format != ns.Format {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != ns.FormatVersion {
		t.Errorf("bad format_version: %q", s.FormatVersion)
	}
	if s.MessageType != ns.MessageType {
		t.Errorf("bad message_type: %q", s.MessageType)
	}
	if s.Placement != ns.Placement {
		t.Errorf("bad placement: %q", s.Placement)
	}

	// Update
	var us *Syslog
	record(t, "syslogs/update", func(c *Client) {
		us, err = c.UpdateSyslog(&UpdateSyslogInput{
			Service:       testServiceID,
			Version:       tv.Number,
			Name:          "test-syslog",
			NewName:       "new-test-syslog",
			FormatVersion: 2,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-syslog" {
		t.Errorf("bad name: %q", us.Name)
	}

	if us.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", us.FormatVersion)
	}

	// Delete
	record(t, "syslogs/delete", func(c *Client) {
		err = c.DeleteSyslog(&DeleteSyslogInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-syslog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSyslogs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListSyslogs(&ListSyslogsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListSyslogs(&ListSyslogsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSyslog_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateSyslog(&CreateSyslogInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateSyslog(&CreateSyslogInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSyslog_validation(t *testing.T) {
	var err error
	_, err = testClient.GetSyslog(&GetSyslogInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSyslog(&GetSyslogInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSyslog(&GetSyslogInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSyslog_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateSyslog(&UpdateSyslogInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSyslog(&UpdateSyslogInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSyslog(&UpdateSyslogInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSyslog_validation(t *testing.T) {
	var err error
	err = testClient.DeleteSyslog(&DeleteSyslogInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSyslog(&DeleteSyslogInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSyslog(&DeleteSyslogInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
