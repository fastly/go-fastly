package fastly

import (
	"testing"
)

func TestClient_Pools(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "pools/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var p *Pool
	record(t, "pools/create", func(c *Client) {
		p, err = c.CreatePool(&CreatePoolInput{
			ServiceID:       testServiceID,
			ServiceVersion:  tv.Number,
			Name:            "test_pool",
			Comment:         "test pool",
			Quorum:          50,
			UseTLS:          true,
			TLSCertHostname: "example.com",
			Type:            PoolTypeRandom,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "pools/cleanup", func(c *Client) {
			c.DeletePool(&DeletePoolInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test_pool",
			})

			c.DeletePool(&DeletePoolInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new_test_pool",
			})
		})
	}()

	if p.Name != "test_pool" {
		t.Errorf("bad name: %q", p.Name)
	}

	if p.Quorum != 50 {
		t.Errorf("bad quorum: %q", p.Quorum)
	}

	if p.UseTLS != true {
		t.Errorf("bad use_tls: %t", p.UseTLS)
	}

	if p.TLSCertHostname != "example.com" {
		t.Errorf("bad tls_cert_hostname: %q", p.TLSCertHostname)
	}

	if p.Type != PoolTypeRandom {
		t.Errorf("bad type: %q", p.Type)
	}

	// List
	var ps []*Pool
	record(t, "pools/list", func(c *Client) {
		ps, err = c.ListPools(&ListPoolsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) < 1 {
		t.Errorf("bad pools: %v", ps)
	}

	// Get
	var np *Pool
	record(t, "pools/get", func(c *Client) {
		np, err = c.GetPool(&GetPoolInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test_pool",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != np.Name {
		t.Errorf("bad name: %q (%q)", p.Name, np.Name)
	}
	if p.Quorum != np.Quorum {
		t.Errorf("bad quorum: %q (%q)", p.Quorum, np.Quorum)
	}
	if p.Type != np.Type {
		t.Errorf("bad type: %q (%q)", p.Type, np.Type)
	}

	// Update
	var up *Pool
	record(t, "pools/update", func(c *Client) {
		up, err = c.UpdatePool(&UpdatePoolInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test_pool",
			NewName:        String("new_test_pool"),
			Quorum:         Uint(0),
			Type:           PPoolType(PoolTypeHash),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if up.Name != "new_test_pool" {
		t.Errorf("bad name: %q", up.Name)
	}
	if up.Quorum != 0 {
		t.Errorf("bad quorum: %q", up.Quorum)
	}

	// Delete
	record(t, "pools/delete", func(c *Client) {
		err = c.DeletePool(&DeletePoolInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new_test_pool",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPools_validation(t *testing.T) {
	var err error
	_, err = testClient.ListPools(&ListPoolsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListPools(&ListPoolsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreatePool_validation(t *testing.T) {
	var err error
	_, err = testClient.CreatePool(&CreatePoolInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreatePool(&CreatePoolInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreatePool(&CreatePoolInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetPool_validation(t *testing.T) {
	var err error
	_, err = testClient.GetPool(&GetPoolInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPool(&GetPoolInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPool(&GetPoolInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePool_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdatePool(&UpdatePoolInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePool(&UpdatePoolInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePool(&UpdatePoolInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeletePool_validation(t *testing.T) {
	var err error
	err = testClient.DeletePool(&DeletePoolInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePool(&DeletePoolInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePool(&DeletePoolInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
