package fastly

import "testing"

func TestClient_IPs(t *testing.T) {
	t.Parallel()

	var err error
	var ips IPAddrs
	record(t, "ips/list", func(c *Client) {
		ips, err = c.IPs()
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ips) == 0 {
		t.Fatal("missing ips")
	}
}
