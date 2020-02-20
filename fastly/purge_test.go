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

func TestClient_createPurgeRequest(t *testing.T) {
	t.Parallel()

	c := DefaultClient()
	purgeUrl := "http://fastly.com/robots.txt"
	req, err := c.createPurgeRequest(purgeUrl)
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != "POST" {
		t.Error("method should be POST for non-HTTPS purges")
	}

	if req.URL.Path != "/purge/"+purgeUrl {
		t.Errorf("url path for POST purges should be prefixed with purge/(url); got %s", req.URL.Path)
	}

	purgeUrl = "https://fastly.com/robots.txt"
	req, err = c.createPurgeRequest(purgeUrl)
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != "PURGE" {
		t.Error("method should be PURGE for HTTPS purges")
	}

	if req.URL.String() != purgeUrl {
		t.Errorf("url path for PURGE purges should be the URL; got %s", req.URL.String())
	}
}
