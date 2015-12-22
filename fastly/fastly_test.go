package fastly

import "testing"

// testClient is the test client.
var testClient = DefaultClient()

// testServiceID is the ID of the testing service.
var testServiceID = "3ywicaizKdg7GxRdRJL58Q"

// testVersion is a new, blank version suitable for testing.
func testVersion(t *testing.T) *Version {
	v, err := testClient.CreateVersion(&CreateVersionInput{
		Service: testServiceID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return v
}
