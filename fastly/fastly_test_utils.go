package fastly

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

// TestClient is the test client.
var TestClient = DefaultClient()

// TestStatsClient is the test client for realtime stats.
var TestStatsClient = NewRealtimeStatsClient()

// ID of the Delivery service for testing.
var TestDeliveryServiceID = deliveryServiceIDForTest()

// ID of the Compute service for testing.
var TestComputeServiceID = computeServiceIDForTest()

// ID of the NGWAF workspace for testing.
var TestNGWAFWorkspaceID = ngwafWorkspaceIDForTest()

// ID of the default Delivery service for testing.
var DefaultDeliveryTestServiceID = "kKJb5bOFI47uHeBVluGfX1"

// ID of the default Compute service for testing.
var defaultComputeTestServiceID = "XsjdElScZGjmfCcTwsYRC1"

// ID of the default NGWAF workspace for testing.
var defaultNGWAFWorkspaceID = "alk6DTsYKHKucJCOIavaJM"

const (
	// ServiceTypeVCL is the type for VCL services.
	ServiceTypeVCL = "vcl"
	// ServiceTypeWasm is the type for Wasm services.
	ServiceTypeWasm = "wasm"
)

// testVersionLock is a lock around version creation because the Fastly API
// kinda dies on concurrent requests to create a version.
var testVersionLock sync.Mutex

func deliveryServiceIDForTest() string {
	if tsid := os.Getenv("FASTLY_TEST_DELIVERY_SERVICE_ID"); tsid != "" {
		return tsid
	}

	return DefaultDeliveryTestServiceID
}

func computeServiceIDForTest() string {
	if tsid := os.Getenv("FASTLY_TEST_COMPUTE_SERVICE_ID"); tsid != "" {
		return tsid
	}

	return defaultComputeTestServiceID
}

func ngwafWorkspaceIDForTest() string {
	if tsid := os.Getenv("FASTLY_TEST_NGWAF_WORKSPACE_ID"); tsid != "" {
		return tsid
	}
	return defaultNGWAFWorkspaceID
}

func vcrDisabled() bool {
	vcrDisable := os.Getenv("VCR_DISABLE")

	return vcrDisable != ""
}

func Record(t *testing.T, fixture string, f func(*Client)) {
	client := DefaultClient()

	if vcrDisabled() {
		f(client)
	} else {
		r := getRecorder(t, fixture)
		defer stopRecorder(t, r)
		client.HTTPClient.Transport = r
		f(client)
	}
}

func RecordIgnoreBody(t *testing.T, fixture string, f func(*Client)) {
	client := DefaultClient()

	if vcrDisabled() {
		f(client)
	} else {
		r := getRecorder(t, fixture)
		defer stopRecorder(t, r)

		r.AddFilter(func(i *cassette.Interaction) error {
			i.Request.Body = ""
			return nil
		})

		client.HTTPClient.Transport = r
		f(client)
	}
}

func getRecorder(t *testing.T, fixture string) *recorder.Recorder {
	r, err := recorder.New("fixtures/" + fixture)
	if err != nil {
		t.Fatal(err)
	}

	// Add a filter which removes Fastly-Key header from all recorded requests.
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Fastly-Key")
		return nil
	})

	return r
}

func stopRecorder(t *testing.T, r *recorder.Recorder) {
	if err := r.Stop(); err != nil {
		t.Fatal(err)
	}
}

func RecordRealtimeStats(t *testing.T, fixture string, f func(*RTSClient)) {
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

func createTestService(t *testing.T, serviceFixture, serviceNameSuffix string) *Service {
	var err error
	var service *Service

	Record(t, serviceFixture, func(client *Client) {
		service, err = client.CreateService(&CreateServiceInput{
			Name:    ToPointer(fmt.Sprintf("test_service_%s", serviceNameSuffix)),
			Comment: ToPointer("go-fastly client test"),
			Type:    ToPointer(ServiceTypeVCL),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	return service
}

func createTestServiceWasm(t *testing.T, serviceFixture, serviceNameSuffix string) *Service {
	var err error
	var service *Service

	Record(t, serviceFixture, func(client *Client) {
		service, err = client.CreateService(&CreateServiceInput{
			Name:    ToPointer(fmt.Sprintf("test_service_wasm_%s", serviceNameSuffix)),
			Comment: ToPointer("go-fastly wasm client test"),
			Type:    ToPointer(ServiceTypeWasm),
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
		ServiceID: TestDeliveryServiceID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func CreateTestVersion(t *testing.T, versionFixture, serviceID string) *Version {
	var err error
	var version *Version

	Record(t, versionFixture, func(client *Client) {
		testVersionLock.Lock()
		defer testVersionLock.Unlock()

		version, err = client.CreateVersion(&CreateVersionInput{
			ServiceID: serviceID,
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	return version
}

func createTestDictionary(t *testing.T, dictionaryFixture, serviceID string, version int, dictionaryNameSuffix string) *Dictionary {
	var err error
	var dictionary *Dictionary

	Record(t, dictionaryFixture, func(client *Client) {
		dictionary, err = client.CreateDictionary(&CreateDictionaryInput{
			ServiceID:      serviceID,
			ServiceVersion: version,
			Name:           ToPointer(fmt.Sprintf("test_dictionary_%s", dictionaryNameSuffix)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return dictionary
}

func deleteTestDictionary(t *testing.T, dictionary *Dictionary, deleteFixture string) {
	var err error

	Record(t, deleteFixture, func(client *Client) {
		err = client.DeleteDictionary(&DeleteDictionaryInput{
			ServiceID:      *dictionary.ServiceID,
			ServiceVersion: *dictionary.ServiceVersion,
			Name:           *dictionary.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createTestACL(t *testing.T, createFixture, serviceID string, version int, aclNameSuffix string) *ACL {
	var err error
	var acl *ACL

	Record(t, createFixture, func(client *Client) {
		acl, err = client.CreateACL(&CreateACLInput{
			ServiceID:      serviceID,
			ServiceVersion: version,
			Name:           ToPointer(fmt.Sprintf("test_acl_%s", aclNameSuffix)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return acl
}

func deleteTestACL(t *testing.T, acl *ACL, deleteFixture string) {
	var err error

	Record(t, deleteFixture, func(client *Client) {
		err = client.DeleteACL(&DeleteACLInput{
			ServiceID:      *acl.ServiceID,
			ServiceVersion: *acl.ServiceVersion,
			Name:           *acl.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createTestPool(t *testing.T, createFixture, serviceID string, version int, poolNameSuffix string) *Pool {
	var err error
	var pool *Pool

	Record(t, createFixture, func(client *Client) {
		pool, err = client.CreatePool(&CreatePoolInput{
			ServiceID:      serviceID,
			ServiceVersion: version,
			Name:           ToPointer(fmt.Sprintf("test_pool_%s", poolNameSuffix)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func createTestLogging(t *testing.T, fixture, serviceID string, serviceNumber int) {
	var err error

	Record(t, fixture, func(c *Client) {
		_, err = c.CreateSyslog(&CreateSyslogInput{
			ServiceID:      serviceID,
			ServiceVersion: serviceNumber,
			Name:           ToPointer("test-syslog"),
			Address:        ToPointer("example.com"),
			Hostname:       ToPointer("example.com"),
			Port:           ToPointer(1234),
			Token:          ToPointer("abcd1234"),
			Format:         ToPointer("format"),
			FormatVersion:  ToPointer(2),
			MessageType:    ToPointer("classic"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestPool(t *testing.T, pool *Pool, deleteFixture string) {
	var err error

	Record(t, deleteFixture, func(client *Client) {
		err = client.DeletePool(&DeletePoolInput{
			ServiceID:      *pool.ServiceID,
			ServiceVersion: *pool.ServiceVersion,
			Name:           *pool.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestLogging(t *testing.T, fixture, serviceID string, serviceNumber int) {
	var err error

	Record(t, fixture, func(c *Client) {
		err = c.DeleteSyslog(&DeleteSyslogInput{
			ServiceID:      serviceID,
			ServiceVersion: serviceNumber,
			Name:           "test-syslog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createTestWAFCondition(t *testing.T, fixture, serviceID, name string, serviceNumber int) *Condition {
	var err error
	var condition *Condition

	Record(t, fixture, func(c *Client) {
		condition, err = c.CreateCondition(&CreateConditionInput{
			ServiceID:      serviceID,
			ServiceVersion: serviceNumber,
			Name:           ToPointer(name),
			Statement:      ToPointer("req.url~+\"index.html\""),
			Type:           ToPointer("PREFETCH"), // This must be a prefetch condition
			Priority:       ToPointer(1),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return condition
}

func deleteTestCondition(t *testing.T, fixture, serviceID, name string, serviceNumber int) {
	var err error

	Record(t, fixture, func(c *Client) {
		err = c.DeleteCondition(&DeleteConditionInput{
			ServiceID:      serviceID,
			ServiceVersion: serviceNumber,
			Name:           name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createTestWAFResponseObject(t *testing.T, fixture, serviceID, name string, serviceNumber int) *ResponseObject {
	var err error
	var ro *ResponseObject

	Record(t, fixture, func(c *Client) {
		ro, err = c.CreateResponseObject(&CreateResponseObjectInput{
			ServiceID:      serviceID,
			ServiceVersion: serviceNumber,
			Name:           ToPointer(name),
			Status:         ToPointer(403),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return ro
}

func deleteTestResponseObject(t *testing.T, fixture, serviceID, name string, serviceNumber int) {
	var err error

	Record(t, fixture, func(c *Client) {
		err = c.DeleteResponseObject(&DeleteResponseObjectInput{
			ServiceID:      serviceID,
			ServiceVersion: serviceNumber,
			Name:           name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createWAF(t *testing.T, fixture, serviceID, condition, response string, serviceNumber int) *WAF {
	var err error
	var waf *WAF

	Record(t, fixture, func(c *Client) {
		waf, err = c.CreateWAF(&CreateWAFInput{
			ServiceID:         serviceID,
			ServiceVersion:    serviceNumber,
			PrefetchCondition: condition,
			Response:          response,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return waf
}

func deleteWAF(t *testing.T, fixture, wafID string) {
	var err error

	Record(t, fixture, func(c *Client) {
		err = c.DeleteWAF(&DeleteWAFInput{
			ID:             wafID,
			ServiceVersion: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestService(t *testing.T, cleanupFixture, serviceID string) {
	var err error

	Record(t, cleanupFixture, func(client *Client) {
		err = client.DeleteService(&DeleteServiceInput{
			ServiceID: serviceID,
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

// pgpPublicKeyUpdate returns a PEM encoded PGP public key suitable for testing
// updates to the public key field.
func pgpPublicKeyUpdate() string {
	return `
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: OpenPGP.js v3.0.9
Comment: https://openpgpjs.org

xsFNBF7XxDYBEACwKFfEjdDsm0lmbdSEM6/RTiUNldEQIpbxIN0pK1z38Rhh
u9VSmO0Jf/x3xTyAvFyc+D5vtxUBp0DlwILVMeTCZ6KX0Yk05hQWTqSMWk6t
QaGXtncvmD9xaHodjuQJ9xJo+lKRhZlli9QnF/x3EcIPSb9fv8BpQRw7M/Kg
O32vkiEV6GvAcZetgLEOy6EkFVx/Y1/kq5crP+gwFxseVcNsH8bB0tMY2maz
TwIKzReyCSjSZn6t+TwdCFzCR6x3cLRegPRqSmpTurL+YtADrHse9hFlPmA0
nHJFTeF96pZQr9N7Hdm/fep1Lb3fEoVgmNyNtsEgZpu8XGB+2deo9A5C9ai+
mkSxkIqv0UDyTuJd4M7Uvt9N1LlWYxeDv2kFo1MDqZF/9LwHAhNk+tpp3zNp
68G5OX/TL5MQoXHlkXmoTBh38sby8IeTIPgS3UhIRUj6a1+6BhoSEAIlGOog
fPkG5ED/WK3aAUcqf8jIJWA01qW/GtwKkx4QTzDN4xYWXtgKrOyaPbpyCDFe
RpX9idrDeoXxA+q822WWJ9P+fIZyxoPi9dyRnysC+1fLTQtw97940KnuYC0N
yQRyn4s24gec9P1BVPS2/1HlyJZxTg3eJCtbPnyT1WbHW4P23hIQ3uChc2/x
h6ywc36ppl7zZPtcRhRnRRrBl/7UECAzC3oewwARAQABzRVmb28gPGJhckBl
eGFtcGxlLmNvbT7CwXUEEAEIACkFAl7XxDYGCwkHCAMCCRCyVpbnvRl/KAQV
CAoCAxYCAQIZAQIbAwIeAQAAWwoP/0abBY3HpYbUVFGHDqDHZ7sTHPD8m6eD
Za6lOI85ARie8Wb39a+KhoLxzo8jInizjCqXgXGLy8tL0sZA5s/K52k9dTPx
lwEjaldH8EeEgHBc9FgUK65MbJMEGXXZRKXescvb5vMu5sroHkC0kDY5sPEP
+rUk4PrgJ0VoOvBulXW/NCaOvaocUuYk5xENp9lUeHSSdDZKh2GJ5Ewbd5uG
BRUhby/QBAd7RAjgJDhC8l//RY9J1QjO693b/WBSJUrgvf8AgO24C7K5r1YE
P64hdSPWQO5ErDPPDqfdUVEWfyZh5flYdwJjk0hS5AGSSjdkKg/mxucb3DXJ
KXWY6In1Ws6Zu3TUqHW714Y0gKOmWYUaRjwvDCZ3+4xJ3FhCriL6IPAex9Dl
lIm7BgylkD2Q+KEm9pDUxTI35UBrmZ1Q89I/EjjUqd1qeoZ88p4YSC+9/RPA
9G03tLCjXvs9G224gxfMas2lhz6gsQIBZ77HDYwfbCAcWgLCD3S3KspfYqs+
M2JTnK5V8+FdRBYJvOj9q+qg9YUTtXYZ6+cEhG/I+6HZt2V6GIOa0vxzJTCk
ITiJqpitX6II0qkh0XpqWMFzoflTEVJxPJBYyHjPDK35pZkoEfvYBBAwoVF8
M6HiP9Eh3xBZbxRToJvAX3frwydyHJoxrZ0gx1/ht8LP3xsPlrP5zsFNBF7X
xDYBEAC4DHAZOhrp2wa8/v1p+uWdE9egUkscCkmFQUNzde5eIuFoGGv5ggi4
EPrj9JxO0S0J+UVCGL4rr0WxbalOSGRe1pgiBK2AIcYtMkESfRQQ7E0H1gCq
TvcimJ9U+NdkFrZfcIu0I/0Pf5Ml+afhl9FXecEwTPG8Ff06nu62PjhABwsx
8wnNhzjXM+xv74ElalccYV8TqWkdJCSc+tKcHXPT0VClWzdwwyfqgozP+5qQ
1Cj9OFGki3KMOn4kD9cX4zS6oC+hj+TIQzh3HIKSLsTDLLg+U0X8+qtYuX3I
2Fi/YLjpIlT/oYlHJlxozbAedLl9NPXmqmxnTEWhsqW1LmjLPwiCxP8plh2b
19QUEmQI48yyp6oWbiBnx0aeXLkSv/LwZ5UJJS+ia8aaN5WiipXZKXgvurS0
2+STGTeD1vEjoQCAKzk35gaz2cTfVZsFcG3ISgccQO3eGE/6+pl0ZpCw2OIs
WOQpnkuAaUkl7fMpLV/D1l660Hz8YmimKppxEow28vXinPt/nOUVjaCOeevJ
IGbjYUBFChu9Oykxg/ZaLRM7OTCe95+TCawcuuJnIXtVekisOE7lRUj/yhx4
cgfVuHtUxg7s+ssmkizV+T9bCpJj3UnUwaVpxPXjGfM1mn+X98xmFMh4zj9U
UESh9N7/EtnwAynAbsdfT/kT+McV3QARAQABwsFfBBgBCAATBQJe18Q2CRCy
VpbnvRl/KAIbDAAA5IwP+QEO6IOJLhr+XmDWLU+FSptWsuMNvWUiHq7wTfD2
70GmBFbdGX5QHtyg4CrbdnfAfjeifhqZYtLNDjLSEDeQDAPX5l+NnUEq6Luh
+ljFdaEFirVA2OCesK1zZ9utF4ZegI9Zukn92OSFYqmEKhhAIH3MaRhZZxQq
3cyXkw0fg2pPkn6PKcDSyLE7RwQMX6+P8X6u/t6kOcWyBJxK4OuJCB2sjFQo
Nphey2TwKv2swJC4mxbHL0S771HelGN+Ske64G5fAxrIEBc+drKtNHlnCZ30
0NRVPzkJGD6dM9WgMyfm/d80YAsnckHGy6WuZ0mQrKiUOoN4vib+tEGnmd+r
K/3afuGq/Shc0tiUFzEKGAUbzaSeWKO6ZgYYGnpxryr1/AEr5z3sPERk3Dm0
+jJ8wIrBBAxUDtuU2TxFRZvvFihrsvMFjBEzEjWNfJRTL7vCBRqR5Qem95tj
fofwzzgSSgf0mlHmyCvhY2Vcr0ozX9oucPJX90cBoWkRqLNCntOuSZSeFUys
IKGV3nYnpXB+VNYozayIoAhP/yVLO7DTibYbzpHgDe3fBG6NKBkQBenHCHpi
BA/4Ob0YNQ6SmKTFS4fIKbMMCnCVgLhMOSyeebVbAuvhP7V6lFdHC9agRXKh
4P1Y4fC/E3SHa+YSCvO4Usl7BoZnKjSULjaeLqeXnlIM
=Wto9
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

// caCert returns a CA certificate suitable for testing.
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
