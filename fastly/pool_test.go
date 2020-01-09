package fastly

import "testing"

func TestClient_Pool(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "pool/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var p *Pool
	record(t, "pool/create", func(c *Client) {
		p, err = c.CreatePool(&CreatePoolInput{
			Service:        testServiceID,
			Version:        tv.Number,
			Name:           "test-pool",
                        Quorum:         75,
                        MaxConnDefault: 200,
			ConnectTimeout: 1000,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "pool/cleanup", func(c *Client) {
			c.DeletePool(&DeletePoolInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-pool",
			})

		})
	}()

	if p.Name != "test-pool" {
		t.Errorf("bad name: %q", p.Name)
	}
        if p.Quorum != 75 {
                t.Errorf("bad quorum: %d", p.Quorum)
        }
	if p.MaxConnDefault != 200 {
		t.Errorf("bad max_conn_default: %d", p.MaxConnDefault)
	}
	if p.ConnectTimeout != 1000 {
		t.Errorf("bad connect_timeout: %d", p.ConnectTimeout)
	}

	// List
	var ps []*Pool
	record(t, "pool/list", func(c *Client) {
		ps, err = c.ListPools(&ListPoolsInput{
			Service: testServiceID,
			Version: tv.Number,
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
	record(t, "pool/get", func(c *Client) {
		np, err = c.GetPool(&GetPoolInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-pool",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != np.Name {
		t.Errorf("bad name: %q (%q)", p.Name, np.Name)
	}
        if p.Quorum != np.Quorum {
                t.Errorf("bad quorum: %d (%d)", p.Quorum, np.Quorum)
        }
        if p.MaxConnDefault != np.MaxConnDefault {
                t.Errorf("bad max_conn_default: %d (%d)", p.MaxConnDefault, np.MaxConnDefault)
        }
	if p.ConnectTimeout != np.ConnectTimeout {
		t.Errorf("bad connect_timeout: %q (%q)", p.ConnectTimeout, np.ConnectTimeout)
	}
	if p.OverrideHost != np.OverrideHost {
		t.Errorf("bad override_host: %q (%q)", p.OverrideHost, np.OverrideHost)
	}

	// Update
	var up *Pool
	record(t, "pool/update", func(c *Client) {
		up, err = c.UpdatePool(&UpdatePoolInput{
			Service:      testServiceID,
			Version:      tv.Number,
			Name:         "test-pool",
			NewName:      "new-test-pool",
			OverrideHost: "www.example.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if up.Name != "new-test-pool" {
		t.Errorf("bad name: %q", up.Name)
	}
	if up.OverrideHost != "www.example.com" {
		t.Errorf("bad override_host: %q", up.OverrideHost)
	}

	// Delete
	record(t, "pool/delete", func(c *Client) {
		err = c.DeletePool(&DeletePoolInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-pool",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPools_validation(t *testing.T) {
	var err error
	_, err = testClient.ListPools(&ListPoolsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListPools(&ListPoolsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreatePool_validation(t *testing.T) {
	var err error
	_, err = testClient.CreatePool(&CreatePoolInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreatePool(&CreatePoolInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

        _, err = testClient.CreatePool(&CreatePoolInput{
                Service: "foo",
                Version: 1,
                Name:    "",
        })
        if err != ErrMissingName {
                t.Errorf("bad error: %s", err)
        }
}

func TestClient_GetPool_validation(t *testing.T) {
	var err error
	_, err = testClient.GetPool(&GetPoolInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPool(&GetPoolInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPool(&GetPoolInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePool_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdatePool(&UpdatePoolInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePool(&UpdatePoolInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePool(&UpdatePoolInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeletePool_validation(t *testing.T) {
	var err error
	err = testClient.DeletePool(&DeletePoolInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePool(&DeletePoolInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePool(&DeletePoolInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
