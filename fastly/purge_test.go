package fastly

import (
	"testing"
)

// TestClient_Purge validates no runtime panics are raised by the Purge method.
//
// Specifically, we're ensuring that the setting of the `Soft` key to `true`
// (which will require assigning a header to the RequestOptions.Headers field)
// doesn't cause `panic: assignment to entry in nil map`.
func TestClient_Purge(t *testing.T) {
	t.Parallel()

	client := DefaultClient()
	url := "http://gofastly.fastly-testing.com/foo/bar"

	_, err := client.Purge(&PurgeInput{
		URL:  url,
		Soft: true,
	})
	if err == nil {
		t.Fatalf("we've accidentally purged a real URL: %s", url)
	}
}

func TestClient_PurgeKey(t *testing.T) {
	t.Parallel()

	var err error
	var purge *Purge
	record(t, "purges/purge_by_key", func(c *Client) {
		purge, err = c.PurgeKey(&PurgeKeyInput{
			ServiceID: testServiceID,
			Key:       "foo",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if purge.Status != "ok" {
		t.Error("bad status")
	}
	if len(purge.ID) == 0 {
		t.Error("bad id")
	}
}
