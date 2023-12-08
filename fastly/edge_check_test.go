package fastly

import "testing"

func TestClient_EdgeCheck(t *testing.T) {
	t.Parallel()

	var err error
	var edges []*EdgeCheck
	record(t, "content/check", func(c *Client) {
		edges, err = c.EdgeCheck(&EdgeCheckInput{
			URL: "releases.hashicorp.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) < 1 {
		t.Errorf("bad edge check: %d", len(edges))
	}
}
