package fastly

import (
	"context"
	"testing"
)

func TestClient_EdgeCheck(t *testing.T) {
	t.Parallel()

	var err error
	var edges []*EdgeCheck
	Record(t, "content/check", func(c *Client) {
		edges, err = c.EdgeCheck(context.TODO(), &EdgeCheckInput{
			URL: "go-fastly-deliver-test.global.ssl.fastly.net",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) < 1 {
		t.Errorf("bad edge check: %d", len(edges))
	}
}
