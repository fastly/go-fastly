package fastly

import "testing"

func TestClient_Purge(t *testing.T) {
	t.Parallel()

	resp, err := testClient.Purge(&PurgeInput{
		URL: "https://releases.hashicorp.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "ok" {
		t.Error("bad status")
	}
	if len(resp.ID) == 0 {
		t.Error("bad id")
	}
}
