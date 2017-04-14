package fastly

import (
	"sync"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
)

// testClient is the test client.
var testClient = DefaultClient()

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
