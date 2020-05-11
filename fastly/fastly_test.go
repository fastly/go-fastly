package fastly

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

// testClient is the test client.
var testClient = DefaultClient()

// testStatsClient is the test client for realtime stats.
var testStatsClient = NewRealtimeStatsClient()

// testServiceID is the ID of the testing service.
var testServiceID = serviceIDForTest()

// Default ID of the testing service.
var defaultTestServiceID = "7i6HN3TK9wS159v2gPAZ8A"

// testVersionLock is a lock around version creation because the Fastly API
// kinda dies on concurrent requests to create a version.
var testVersionLock sync.Mutex

func serviceIDForTest() string {
	tsid := os.Getenv("FASTLY_TEST_SERVICE_ID")

	if tsid != "" {
		return tsid
	}

	return defaultTestServiceID
}

func vcrDisabled() bool {
	vcrDisable := os.Getenv("VCR_DISABLE")

	return vcrDisable != ""
}

func record(t *testing.T, fixture string, f func(*Client)) {
	client := DefaultClient()

	if vcrDisabled() {
		f(client)
	} else {
		r, err := recorder.New("fixtures/" + fixture)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := r.Stop(); err != nil {
				t.Fatal(err)
			}
		}()

		// Add a filter which removes Fastly-Key header from all recorded requests.
		r.AddFilter(func(i *cassette.Interaction) error {
			delete(i.Request.Headers, "Fastly-Key")
			return nil
		})

		client.HTTPClient.Transport = r

		f(client)
	}
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

func createTestService(t *testing.T, serviceFixture string, serviceNameSuffix string) *Service {

	var err error
	var service *Service

	record(t, serviceFixture, func(client *Client) {
		service, err = client.CreateService(&CreateServiceInput{
			Name:    fmt.Sprintf("test_service_%s", serviceNameSuffix),
			Comment: "go-fastly client test",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	return service
}

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

func createTestVersion(t *testing.T, versionFixture string, serviceId string) *Version {

	var err error
	var version *Version

	record(t, versionFixture, func(client *Client) {
		testVersionLock.Lock()
		defer testVersionLock.Unlock()

		version, err = client.CreateVersion(&CreateVersionInput{
			Service: serviceId,
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	return version
}

func createTestDictionary(t *testing.T, dictionaryFixture string, serviceId string, version int, dictionaryNameSuffix string) *Dictionary {

	var err error
	var dictionary *Dictionary

	record(t, dictionaryFixture, func(client *Client) {
		dictionary, err = client.CreateDictionary(&CreateDictionaryInput{
			Service: serviceId,
			Version: version,
			Name:    fmt.Sprintf("test_dictionary_%s", dictionaryNameSuffix),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return dictionary
}

func deleteTestDictionary(t *testing.T, dictionary *Dictionary, deleteFixture string) {

	var err error

	record(t, deleteFixture, func(client *Client) {
		err = client.DeleteDictionary(&DeleteDictionaryInput{
			Service: dictionary.ServiceID,
			Version: dictionary.Version,
			Name:    dictionary.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createTestACL(t *testing.T, createFixture string, serviceId string, version int, aclNameSuffix string) *ACL {

	var err error
	var acl *ACL

	record(t, createFixture, func(client *Client) {
		acl, err = client.CreateACL(&CreateACLInput{
			Service: serviceId,
			Version: version,
			Name:    fmt.Sprintf("test_acl_%s", aclNameSuffix),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return acl
}

func deleteTestACL(t *testing.T, acl *ACL, deleteFixture string) {

	var err error

	record(t, deleteFixture, func(client *Client) {
		err = client.DeleteACL(&DeleteACLInput{
			Service: acl.ServiceID,
			Version: acl.Version,
			Name:    acl.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createTestPool(t *testing.T, createFixture string, serviceId string, version int, poolNameSuffix string) *Pool {

	var err error
	var pool *Pool

	record(t, createFixture, func(client *Client) {
		pool, err = client.CreatePool(&CreatePoolInput{
			Service: serviceId,
			Version: version,
			Name:    fmt.Sprintf("test_pool_%s", poolNameSuffix),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func deleteTestPool(t *testing.T, pool *Pool, deleteFixture string) {

	var err error

	record(t, deleteFixture, func(client *Client) {
		err = client.DeletePool(&DeletePoolInput{
			Service: pool.ServiceID,
			Version: pool.Version,
			Name:    pool.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestService(t *testing.T, cleanupFixture string, serviceId string) {

	var err error

	record(t, cleanupFixture, func(client *Client) {
		err = client.DeleteService(&DeleteServiceInput{
			ID: serviceId,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

// pgpPublicKey returns a PEM encoded PGP public key suitable for testing.
func pgpPublicKey() string {
	return `-----BEGIN PGP PUBLIC KEY BLOCK-----

mQENBFyUD8sBCACyFnB39AuuTygseek+eA4fo0cgwva6/FSjnWq7riouQee8GgQ/
ibXTRyv4iVlwI12GswvMTIy7zNvs1R54i0qvsLr+IZ4GVGJqs6ZJnvQcqe3xPoR4
8AnBfw90o32r/LuHf6QCJXi+AEu35koNlNAvLJ2B+KACaNB7N0EeWmqpV/1V2k9p
lDYk+th7LcCuaFNGqKS/PrMnnMqR6VDLCjHhNx4KR79b0Twm/2qp6an3hyNRu8Gn
dwxpf1/BUu3JWf+LqkN4Y3mbOmSUL3MaJNvyQguUzTfS0P0uGuBDHrJCVkMZCzDB
89ag55jCPHyGeHBTd02gHMWzsg3WMBWvCsrzABEBAAG0JXRlcnJhZm9ybSAodGVz
dCkgPHRlc3RAdGVycmFmb3JtLmNvbT6JAU4EEwEIADgWIQSHYyc6Kj9l6HzQsau6
vFFc9jxV/wUCXJQPywIbAwULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRC6vFFc
9jxV/815CAClb32OxV7wG01yF97TzlyTl8TnvjMtoG29Mw4nSyg+mjM3b8N7iXm9
OLX59fbDAWtBSldSZE22RXd3CvlFOG/EnKBXSjBtEqfyxYSnyOPkMPBYWGL/ApkX
SvPYJ4LKdvipYToKFh3y9kk2gk1DcDBDyaaHvR+3rv1u3aoy7/s2EltAfDS3ZQIq
7/cWTLJml/lleeB/Y6rPj8xqeCYhE5ahw9gsV/Mdqatl24V9Tks30iijx0Hhw+Gx
kATUikMGr2GDVqoIRga5kXI7CzYff4rkc0Twn47fMHHHe/KY9M2yVnMHUXmAZwbG
M1cMI/NH1DjevCKdGBLcRJlhuLPKF/anuQENBFyUD8sBCADIpd7r7GuPd6n/Ikxe
u6h7umV6IIPoAm88xCYpTbSZiaK30Svh6Ywra9jfE2KlU9o6Y/art8ip0VJ3m07L
4RSfSpnzqgSwdjSq5hNour2Fo/BzYhK7yaz2AzVSbe33R0+RYhb4b/6N+bKbjwGF
ftCsqVFMH+PyvYkLbvxyQrHlA9woAZaNThI1ztO5rGSnGUR8xt84eup28WIFKg0K
UEGUcTzz+8QGAwAra+0ewPXo/AkO+8BvZjDidP417u6gpBHOJ9qYIcO9FxHeqFyu
YrjlrxowEgXn5wO8xuNz6Vu1vhHGDHGDsRbZF8pv1d5O+0F1G7ttZ2GRRgVBZPwi
kiyRABEBAAGJATYEGAEIACAWIQSHYyc6Kj9l6HzQsau6vFFc9jxV/wUCXJQPywIb
DAAKCRC6vFFc9jxV/9YOCACe8qmOSnKQpQfW+PqYOqo3dt7JyweTs3FkD6NT8Zml
dYy/vkstbTjPpX6aTvUZjkb46BVi7AOneVHpD5GBqvRsZ9iVgDYHaehmLCdKiG5L
3Tp90NN+QY5WDbsGmsyk6+6ZMYejb4qYfweQeduOj27aavCJdLkCYMoRKfcFYI8c
FaNmEfKKy/r1PO20NXEG6t9t05K/frHy6ZG8bCNYdpagfFVot47r9JaQqWlTNtIR
5+zkkSq/eG9BEtRij3a6cTdQbktdBzx2KBeI0PYc1vlZR0LpuFKZqY9vlE6vTGLR
wMfrTEOvx0NxUM3rpaCgEmuWbB1G1Hu371oyr4srrr+N
=28dr
-----END PGP PUBLIC KEY BLOCK-----
`
}

// privatekey returns a ASN.1 DER encoded key suitable for testing.
func privateKey() string {
	return `-----BEGIN PRIVATE KEY-----
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
ImXgKS5UmJB53s9yiomen/n7cUXvrAk=
-----END PRIVATE KEY-----
`
}

// certificate returns a ASN.1 DER encoded certificate for the private key suitable for testing.
func certificate() string {
	return `-----BEGIN CERTIFICATE-----
MIIE6DCCAtACCQCzHO2a8qU6KzANBgkqhkiG9w0BAQsFADA2MRIwEAYDVQQDDAls
b2NhbGhvc3QxIDAeBgNVBAoMF0NsaWVudCBDZXJ0aWZpY2F0ZSBEZW1vMB4XDTE5
MTIwNTE3MjY1N1oXDTIwMTIwNDE3MjY1N1owNjESMBAGA1UEAwwJbG9jYWxob3N0
MSAwHgYDVQQKDBdDbGllbnQgQ2VydGlmaWNhdGUgRGVtbzCCAiIwDQYJKoZIhvcN
AQEBBQADggIPADCCAgoCggIBAJ3iM9y8yWa9P+MRVj9Jj+rpNGDDrS1z2kl7SgDB
ZisFYCHccvLfnz10OCWkI4fUWeAbsg2QJZ8i1EdFDtWnqhPXyqUsHRZRAI92ujvf
4g9r7ID0MENfVKOIY7ye3aphHUEW5lAkkJae12QcxyjTbrbfE0kQXpBHYaqFPRZs
CTXZOOZfSuEVxcrRJMjlt1J8thnGGmaa3yH6o1yt6hQdTp9JNzzeUwWe78PaAms8
RaLcaDXC7fsgByI8j9coKOluQdinYxkBLAHMpo7x6QuOYA3oguZXAUpMXJMJCQD2
Nd0WL33Dy4XHcPrhw+lfCVW7E1sbWO3Ka1ZVeu0FEF0erOKAmfQk2eoMJzR9qKwm
kbM2rYq5qrN19a5Rdtxxov4zOuOyvI6b1Uz5PajN1dyXzuKImOXFmiEL7ykC8kxD
u8HA90pC2VK3V0mx6tsR/H6zMBIg2je51nJ/11VCmIgS5/+jso1h+oUtqHAsWi1/
5u8lrQzMC3CR3VKLGCWhfF7NpQ82DYLnBh60t/E4mY0WX7GDAY8QTMKd4dnmKU6d
nKYDzXZR1he1c08+6NX+pdzJxih8Q/EG0PkNNla0MabMDsi7eFMUCjSPOUG99vGW
PNvMqP/EvCmCW7VKmph4NSNHjwkxTTOQD/ZGX9IpQWZkCxyCIMxAi1hCgu0zKR++
U+7BAgMBAAEwDQYJKoZIhvcNAQELBQADggIBAC8av9I4ezwlmM7ysaJvC1IfCzNN
VawIK1U7bfj9Oyjl49Bn/yTwbbiQ8j5VjOza4umIwnYp1HP6/mlBO+ey8WFYPmDM
JAspk6sYEQW7MrbZ9QOmq24YAkwMzgK1hDASCKq4GJCzGDym3Zx6fvPnMCPdei2c
jgtjzzBmyewE0zcegOHDrFXTaUIfoSbduTbV9zClJ/pJDkTklRX1cYBtIox77gpZ
1cnIC803gi1rVCHRNdq85ltOTjoF1/wVamLy5c6CYlp5IPyVOm0nqbqra47QIwss
QSGxn5l52BC1jP1l3eK1mEr64+dbMhqX3ZQwhfuiQ9VmdovNN1NarPWfmQia6Spq
TfxN+3VhloKFUh+fgwNzWYLKCMnwBuPVhVGcpclvrj50MsyeiT2IfE8pqWw26g6g
0xu85AbqYKePaZ7wPoDddbwCIGr6BBT87Nsu+AqtnkH3uw34FDDcjWR1CmNuI1mP
ac6d1jdfbkL5ZUJTpTJi0BxWbTGupv8VzufteFRNa7U2h1O6+kyPmEpA3heEZcEO
Hq5zIfW6QTUmBXDfMFzQ9h0764oBVwm29bjZ59bU3RhcAZtL8fi5BapNtoKAy55d
P/0WahbwNjP68QYVLBeK9Sfo0XxLU0hJP4RJUZSXy9kUuZ8xhAM/6PdE04cDq71v
Zfq6/HA3phy85qyj
-----END CERTIFICATE-----
`
}

// caCert returns a CA certificate suitable for testing
func caCert() string {
	return `-----BEGIN CERTIFICATE-----
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
`
}
