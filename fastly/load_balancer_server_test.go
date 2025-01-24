package fastly

import (
	"errors"
	"testing"
)

func TestClient_Servers(t *testing.T) {
	var err error
	var tv *Version
	Record(t, "servers/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	testPool := createTestPool(t, "servers/create_pool", TestDeliveryServiceID, *tv.Number, "servers22")

	// Create
	var server *Server
	var altServer *Server
	Record(t, "servers/create", func(c *Client) {
		server, err = c.CreateServer(&CreateServerInput{
			ServiceID: TestDeliveryServiceID,
			PoolID:    *testPool.PoolID,
			Address:   ToPointer("127.0.0.1"),
		})
		if err != nil {
			t.Fatal(err)
		}

		// additional pool server for DeleteServer usage
		altServer, err = c.CreateServer(&CreateServerInput{
			ServiceID: TestDeliveryServiceID,
			PoolID:    *testPool.PoolID,
			Address:   ToPointer("altserver.example.com"),
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	// Ensure deleted
	defer func() {
		// Delete the pool from this version.
		deleteTestPool(t, testPool, "servers/delete_pool")

		Record(t, "servers/cleanup", func(c *Client) {
			// Expected to fail as this was explicitly deleted in the test.
			_ = c.DeleteServer(&DeleteServerInput{
				ServiceID: TestDeliveryServiceID,
				PoolID:    *testPool.PoolID,
				Server:    *altServer.ServerID,
			})

			// Expected to fail as the API forbids deleting the last server in
			// the pool. The pool is deleted from this version but it still
			// exists as it may be associated with other versions.
			_ = c.DeleteServer(&DeleteServerInput{
				ServiceID: TestDeliveryServiceID,
				PoolID:    *testPool.PoolID,
				Server:    *server.ServerID,
			})
		})
	}()

	if *server.ServiceID != TestDeliveryServiceID {
		t.Errorf("bad server service: %q", *server.ServiceID)
	}
	if *server.PoolID != *testPool.PoolID {
		t.Errorf("bad server pool: %q", *server.PoolID)
	}
	if *server.Address != "127.0.0.1" {
		t.Errorf("bad server address: %q", *server.Address)
	}

	// List
	var ss []*Server
	Record(t, "servers/list", func(c *Client) {
		ss, err = c.ListServers(&ListServersInput{
			ServiceID: TestDeliveryServiceID,
			PoolID:    *testPool.PoolID,
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
	Record(t, "servers/get", func(c *Client) {
		ns, err = c.GetServer(&GetServerInput{
			ServiceID: TestDeliveryServiceID,
			PoolID:    *testPool.PoolID,
			Server:    *server.ServerID,
		})
	})
	if *server.ServerID != *ns.ServerID {
		t.Errorf("bad ID: %q (%q)", *server.ServerID, *ns.ServerID)
	}
	if *server.Address != *ns.Address {
		t.Errorf("bad address: %q (%q)", *server.Address, *ns.Address)
	}

	// Update
	var us *Server
	Record(t, "servers/update", func(c *Client) {
		us, err = c.UpdateServer(&UpdateServerInput{
			ServiceID: TestDeliveryServiceID,
			PoolID:    *testPool.PoolID,
			Server:    *server.ServerID,
			Address:   ToPointer("0.0.0.0"),
			Weight:    ToPointer(50),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Address == *server.Address {
		t.Errorf("bad address: %s", *us.Address)
	}
	if *us.Weight != 50 {
		t.Errorf("bad weight: %q", 50)
	}

	// Delete
	Record(t, "servers/delete", func(c *Client) {
		err = c.DeleteServer(&DeleteServerInput{
			ServiceID: TestDeliveryServiceID,
			PoolID:    *testPool.PoolID,
			Server:    *altServer.ServerID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListServers_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListServers(&ListServersInput{
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingPoolID) {
		t.Errorf("bad error: %q", err)
	}

	_, err = TestClient.ListServers(&ListServersInput{
		PoolID: "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateServer_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateServer(&CreateServerInput{
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingPoolID) {
		t.Errorf("bad error: %q", err)
	}

	_, err = TestClient.CreateServer(&CreateServerInput{
		PoolID: "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetServer_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetServer(&GetServerInput{
		Server:    "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingPoolID) {
		t.Errorf("bad error: %q", err)
	}

	_, err = TestClient.GetServer(&GetServerInput{
		PoolID:    "bar",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServer) {
		t.Errorf("bad error: %q", err)
	}

	_, err = TestClient.GetServer(&GetServerInput{
		PoolID: "bar",
		Server: "test",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateServer_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateServer(&UpdateServerInput{
		Server:    "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingPoolID) {
		t.Errorf("bad error: %q", err)
	}

	_, err = TestClient.UpdateServer(&UpdateServerInput{
		PoolID:    "bar",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServer) {
		t.Errorf("bad error: %q", err)
	}

	_, err = TestClient.UpdateServer(&UpdateServerInput{
		PoolID: "bar",
		Server: "test",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteServer_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteServer(&DeleteServerInput{
		Server:    "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingPoolID) {
		t.Errorf("bad error: %q", err)
	}

	err = TestClient.DeleteServer(&DeleteServerInput{
		PoolID:    "bar",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServer) {
		t.Errorf("bad error: %q", err)
	}

	err = TestClient.DeleteServer(&DeleteServerInput{
		PoolID: "bar",
		Server: "test",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}
