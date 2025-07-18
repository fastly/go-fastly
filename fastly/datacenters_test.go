package fastly

import (
	"context"
	"testing"
)

func TestDatacenters(t *testing.T) {
	t.Parallel()

	var err error
	var datacenters []Datacenter
	Record(t, "datacenters/list", func(c *Client) {
		datacenters, err = c.AllDatacenters(context.TODO())
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(datacenters) == 0 {
		t.Fatal("missing datacenters")
	}
}
