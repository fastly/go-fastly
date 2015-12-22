package fastly

import "testing"

func TestClient_EdgeCheck(t *testing.T) {
	t.Parallel()

	e, err := testClient.EdgeCheck(&EdgeCheckInput{
		URL: "releases.hashicorp.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(e) < 1 {
		t.Errorf("bad edge check: %d", len(e))
	}
}
