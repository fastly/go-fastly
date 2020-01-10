package fastly

import "testing"

func TestClient_Servers(t *testing.T) {
	var err error
	var tv *Version
	record(t, "servers/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	testPool := createTestPool(t, "servers/create_pool", testServiceID, tv.Number, "servers22")

	// Create
	var server *Server
	var altServer *Server
	record(t, "servers/create", func(c *Client) {
		server, err = c.CreateServer(&CreateServerInput{
			Service: testServiceID,
			Pool:    testPool.ID,
			Address: "127.0.0.1",
		})
		if err != nil {
			t.Fatal(err)
		}

		// additional pool server for DeleteServer usage
		altServer, err = c.CreateServer(&CreateServerInput{
			Service: testServiceID,
			Pool:    testPool.ID,
			Address: "altserver.example.com",
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	// Ensure deleted
	defer func() {
		// Delete the pool from this version.
		deleteTestPool(t, testPool, "servers/delete_pool")

		record(t, "servers/cleanup", func(c *Client) {
			// Expected to fail as this was explicitly deleted in the test.
			c.DeleteServer(&DeleteServerInput{
				Service: testServiceID,
				Pool:    testPool.ID,
				Server:  altServer.ID,
			})

			// Expected to fail as the API forbids deleting the last server in
			// the pool. The pool is deleted from this version but it still
			// exists as it may be associated with other versions.
			c.DeleteServer(&DeleteServerInput{
				Service: testServiceID,
				Pool:    testPool.ID,
				Server:  server.ID,
			})
		})
	}()

	if server.ServiceID != testServiceID {
		t.Errorf("bad server service: %q", server.ServiceID)
	}
	if server.PoolID != testPool.ID {
		t.Errorf("bad server pool: %q", server.PoolID)
	}
	if server.Address != "127.0.0.1" {
		t.Errorf("bad server address: %q", server.Address)
	}

	// List
	var ss []*Server
	record(t, "servers/list", func(c *Client) {
		ss, err = c.ListServers(&ListServersInput{
			Service: testServiceID,
			Pool:    testPool.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad servers: %v", ss)
	}

	// Get
	var ns *Server
	record(t, "servers/get", func(c *Client) {
		ns, err = c.GetServer(&GetServerInput{
			Service: testServiceID,
			Pool:    testPool.ID,
			Server:  server.ID,
		})
	})
	if server.ID != ns.ID {
		t.Errorf("bad ID: %q (%q)", server.ID, ns.ID)
	}
	if server.Address != ns.Address {
		t.Errorf("bad address: %q (%q)", server.Address, ns.Address)
	}

	// Update
	var us *Server
	record(t, "servers/update", func(c *Client) {
		us, err = c.UpdateServer(&UpdateServerInput{
			Service: testServiceID,
			Pool:    testPool.ID,
			Server:  server.ID,
			Address: "0.0.0.0",
			Weight:  50,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Address == server.Address {
		t.Errorf("bad address: %s", us.Address)
	}
	if us.Weight != 50 {
		t.Errorf("bad weight: %q", 50)
	}

	// Delete
	record(t, "servers/delete", func(c *Client) {
		err = c.DeleteServer(&DeleteServerInput{
			Service: testServiceID,
			Pool:    testPool.ID,
			Server:  altServer.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListServers_validation(t *testing.T) {
	var err error
	_, err = testClient.ListServers(&ListServersInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListServers(&ListServersInput{
		Service: "foo",
		Pool:    "",
	})
	if err != ErrMissingPool {
		t.Errorf("bad error: %q", err)
	}
}

func TestClient_CreateServer_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateServer(&CreateServerInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateServer(&CreateServerInput{
		Service: "foo",
		Pool:    "",
	})
	if err != ErrMissingPool {
		t.Errorf("bad error: %q", err)
	}
}

func TestClient_GetServer_validation(t *testing.T) {
	var err error
	_, err = testClient.GetServer(&GetServerInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetServer(&GetServerInput{
		Service: "foo",
		Pool:    "",
	})
	if err != ErrMissingPool {
		t.Errorf("bad error: %q", err)
	}

	_, err = testClient.GetServer(&GetServerInput{
		Service: "foo",
		Pool:    "bar",
		Server:  "",
	})
	if err != ErrMissingServer {
		t.Errorf("bad error: %q", err)
	}
}

func TestClient_UpdateServer_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateServer(&UpdateServerInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateServer(&UpdateServerInput{
		Service: "foo",
		Pool:    "",
	})
	if err != ErrMissingPool {
		t.Errorf("bad error: %q", err)
	}

	_, err = testClient.UpdateServer(&UpdateServerInput{
		Service: "foo",
		Pool:    "bar",
		Server:  "",
	})
	if err != ErrMissingServer {
		t.Errorf("bad error: %q", err)
	}
}

func TestClient_DeleteServer_validation(t *testing.T) {
	var err error
	err = testClient.DeleteServer(&DeleteServerInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteServer(&DeleteServerInput{
		Service: "foo",
		Pool:    "",
	})
	if err != ErrMissingPool {
		t.Errorf("bad error: %q", err)
	}

	err = testClient.DeleteServer(&DeleteServerInput{
		Service: "foo",
		Pool:    "bar",
		Server:  "",
	})
	if err != ErrMissingServer {
		t.Errorf("bad error: %q", err)
	}
}
