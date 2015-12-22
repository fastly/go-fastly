package fastly

import "testing"

func TestClient_IPs(t *testing.T) {
	t.Parallel()

	ips, err := testClient.IPs()
	if err != nil {
		t.Fatal(err)
	}
	if len(ips) == 0 {
		t.Fatal("missing ips")
	}
}
