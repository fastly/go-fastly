package fastly

import "testing"

func TestClient_IPs(t *testing.T) {
	t.Parallel()

	var err error
	var ips IPAddrs
	Record(t, "ips/list", func(c *Client) {
		ips, err = c.IPs()
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ips) == 0 {
		t.Fatal("missing ips")
	}
}

func TestClient_IPsV6(t *testing.T) {
	t.Parallel()

	var err error
	var ips IPAddrs
	Record(t, "ips/list", func(c *Client) {
		ips, err = c.IPsV6()
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ips) == 0 {
		t.Fatal("missing ips")
	}
}

func TestClient_AllIPs(t *testing.T) {
	t.Parallel()

	var err error
	var v4 IPAddrs
	var v6 IPAddrs
	Record(t, "ips/list", func(c *Client) {
		v4, v6, err = c.AllIPs()
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(v4) == 0 {
		t.Fatal("missing v4 ips")
	}
	if len(v6) == 0 {
		t.Fatal("missing v6 ips")
	}
}
