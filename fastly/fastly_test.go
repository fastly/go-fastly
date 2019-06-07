package fastly

import (
	"sync"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
)

// testClient is the test client.
var testClient = DefaultClient()

// testStatsClient is the test client for realtime stats.
var testStatsClient = NewRealtimeStatsClient()

// testServiceID is the ID of the testing service.
var testServiceID = "7i6HN3TK9wS159v2gPAZ8A"

// testVersionLock is a lock around version creation because the Fastly API
// kinda dies on concurrent requests to create a version.
var testVersionLock sync.Mutex

// testVersion is a new, blank version suitable for testing.
func testVersion(t *testing.T, c *Client) *Version {
	testVersionLock.Lock()
	defer testVersionLock.Unlock()

	v, err := c.CreateVersion(&CreateVersionInput{
		Service: testServiceID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func record(t *testing.T, fixture string, f func(*Client)) {
	r, err := recorder.New("fixtures/" + fixture)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := r.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	client := DefaultClient()
	client.HTTPClient.Transport = r

	f(client)
}

func recordRealtimeStats(t *testing.T, fixture string, f func(*RTSClient)) {
	r, err := recorder.New("fixtures/" + fixture)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := r.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	client := NewRealtimeStatsClient()
	client.client.HTTPClient.Transport = r

	f(client)
}

// privatekey returns a ASN.1 DER encoded key suitable for testing.
func privateKey() string {
	return `-----BEGIN PRIVATE KEY-----
MIICXAIBAAKBgQCukSu6ece/5jlgnNiLK7mIxJuLLeZ1FI/rn0PK8XTO+vSZjV+o
vXpiSiXavXv2XhQymAf/tOpDV+uBxIzM3E1wyQ2qUxFjZ1zJHqBlmK+QJeg7pKmD
gVQrzdtfUYGTlpizMiucWyYJYHNdfIkoD9N+wSw/cM/kxHa8fYtGsN3M7wIDAQAB
AoGAOXiBtPqy0HKzISOCBw92HZjcvI13+cOzPhdI8l9b3WixbnwkqiD3UbSnkcQg
M5P1glKbD4w4M8OWPTrAQBGnTJa2iA7z9IqTGL7dwhVnQ04NYq14CpPs+XCKwOxK
O0gEjgbjlPNyE3OsMNBFMB5rnsVEI8uUukBmm/h6l8x7HUECQQDCeatYbfZZ8ra3
twTCmuc4qEXMzLhxBq7Ogyst90mC1fhfxFKiXwu0WVJqeGSFrK1upHPDfPJ3JN75
CDskG9YzAkEA5csoxhmPblUOG3e/Vt8dzjPlk/ZTgxEOKPb86BUprsR2J1aJVHVH
EoZjlAj5yo7iNSphp4cVXJd8I+ZsYSeaVQJBAJIF+5N9lcG6Tlop0SgyWbWgHDEH
8uHjS7SCpxRvnsHf2gxGhGmpBkfX3dtWJNx+aQcv8kBx/Dlb9RR2irm1MSMCQCmM
xICdWovuoTBiRJymlzMTuy032v3V9aN+lVg5i2HocBzIzugQlJtK5XJ89P2lPE20
rhemmzw0v+OV5H7ktEkCQDqLUZyqnGHX6qV+8eJMafLyy1AUfzSkbuZ/nX6hed8T
cfzsfBxi4bN4JOkJcA77FpXDecX/GDwzRN+yfwNs3+0=
-----END PRIVATE KEY-----
`
}
