package fastly

import "testing"

func TestClient_Purge(t *testing.T) {
	t.Parallel()

	var err error
	var purge *Purge
	record(t, "purges/purge_by_key", func(c *Client) {
		purge, err = c.PurgeKey(&PurgeKeyInput{
			Service: testServiceID,
			Key:     "foo",
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
