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
			Comment:         ToPointer("test pool"),
			Name:            ToPointer("test_pool"),
			Quorum:          ToPointer(50),
			ServiceID:       testServiceID,
			ServiceVersion:  *tv.Number,
			TLSCertHostname: ToPointer("example.com"),
			Type:            ToPointer(PoolTypeRandom),
			UseTLS:          ToPointer(Compatibool(true)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "pools/cleanup", func(c *Client) {
			_ = c.DeletePool(&DeletePoolInput{
				ServiceID:      testServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test_pool",
			})

			_ = c.DeletePool(&DeletePoolInput{
				ServiceID:      testServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new_test_pool",
			})
		})
	}()

	if *p.Name != "test_pool" {
		t.Errorf("bad name: %q", *p.Name)
	}
	if *p.Quorum != 50 {
		t.Errorf("bad quorum: %q", *p.Quorum)
	}
	if !*p.UseTLS {
		t.Errorf("bad use_tls: %t", *p.UseTLS)
	}
	if *p.TLSCertHostname != "example.com" {
		t.Errorf("bad tls_cert_hostname: %q", *p.TLSCertHostname)
	}
	if *p.Type != PoolTypeRandom {
		t.Errorf("bad type: %q", *p.Type)
	}

	// List
	var ps []*Pool
	record(t, "pools/list", func(c *Client) {
		ps, err = c.ListPools(&ListPoolsInput{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
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
			ServiceVersion: *tv.Number,
			Name:           "test_pool",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *p.Name != *np.Name {
		t.Errorf("bad name: %q (%q)", *p.Name, *np.Name)
	}
	if *p.Quorum != *np.Quorum {
		t.Errorf("bad quorum: %q (%q)", *p.Quorum, *np.Quorum)
	}
	if *p.Type != *np.Type {
		t.Errorf("bad type: %q (%q)", *p.Type, *np.Type)
	}

	// Update
	var up *Pool
	record(t, "pools/update", func(c *Client) {
		up, err = c.UpdatePool(&UpdatePoolInput{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test_pool",
			NewName:        ToPointer("new_test_pool"),
			Quorum:         ToPointer(0),
			Type:           ToPointer(PoolTypeHash),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *up.Name != "new_test_pool" {
		t.Errorf("bad name: %q", *up.Name)
	}
	if *up.Quorum != 0 {
		t.Errorf("bad quorum: %q", *up.Quorum)
	}

	// Delete
	record(t, "pools/delete", func(c *Client) {
		err = c.DeletePool(&DeletePoolInput{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
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
		Name:           ToPointer("test"),
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreatePool(&CreatePoolInput{
		Name:      ToPointer("test"),
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetPool_validation(t *testing.T) {
	var err error

	_, err = testClient.GetPool(&GetPoolInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPool(&GetPoolInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPool(&GetPoolInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePool_validation(t *testing.T) {
	var err error

	_, err = testClient.UpdatePool(&UpdatePoolInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePool(&UpdatePoolInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePool(&UpdatePoolInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeletePool_validation(t *testing.T) {
	var err error

	err = testClient.DeletePool(&DeletePoolInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePool(&DeletePoolInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePool(&DeletePoolInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
