package fastly

import (
	"testing"
)

func TestClient_ListTLSDomains(t *testing.T) {
	t.Parallel()

	fixtureBase := "custom_tls_domain/"

	var err error

	// List
	var ldom []*TLSDomain
	record(t, fixtureBase+"list", func(c *Client) {
		ldom, err = c.ListTLSDomains(&ListTLSDomainsInput{
			PageSize: 10,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ldom) < 1 {
		t.Errorf("bad tls domains: %v", ldom)
	}

}
