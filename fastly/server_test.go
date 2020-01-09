package fastly

import "testing"

func TestClient_server(t *testing.T) {
	t.Parallel()


        var serverID = "srvsrvsrvsrvsrvsrv"
	var err error


        fixtureBase := "server/"
        nameSuffix := "test_pool"

        testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
        defer deleteTestService(t, fixtureBase+"delete_service", testService.ID)

        testVersion := createTestVersion(t, fixtureBase+"version", testService.ID)

        testPool := createTestPool(t, fixtureBase+"create_pool", testService.ID, testVersion.Number, nameSuffix)
        defer deleteTestPool(t, testPool, fixtureBase+"delete_pool")


	// Create
	var s *Server
	record(t, fixtureBase+"create_server", func(c *Client) {
		s, err = c.CreateServer(&CreateServerInput{
			Service:        testService.ID,
			Pool:           testPool.Id,
                        Weight:         "75",
        		MaxConn:        "777",
        		Port:           "1234",
        		Address:        "integ-test.go-fastly.com",
        		Disabled:       false,
        		OverrideHost:   "origin.example.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"delete_server", func(c *Client) {
			c.DeleteServer(&DeleteServerInput{
				Service: testServiceID,
				Pool:    testPoolID,
                                Id:      serverID,
			})

		})
	}()

        if s.Weight != "75" {
                t.Errorf("bad weight: %q", s.Weight)
        }
        if s.MaxConn != "777" {
                t.Errorf("bad maxconn: %q", s.MaxConn)
        }
	if s.Port != "1234" {
		t.Errorf("bad port: %q", s.Port)
	}
        if s.Address != "integ-test.go-fastly.com" {
                t.Errorf("bad address: %q", s.Address)
        }
        if s.Disabled != false {
                t.Errorf("bad disabled: %t", s.Disabled)
        }
	if s.OverrideHost != "origin.example.com" {
		t.Errorf("bad override_host: %q", s.OverrideHost)
	}

	// List
	var servers []*Server
	record(t, fixtureBase+"list_servers", func(c *Client) {
		servers, err = c.ListServers(&ListServersInput{
                                 Service:        testService.ID,
                                 Pool:           testPool.Id,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) < 1 {
		t.Errorf("bad servers: %v", servers)
	}

	// Get
	var ns *Server
	record(t, fixtureBase+"get_server", func(c *Client) {
		ns, err = c.GetServer(&GetServerInput{
			Service: testServiceID,
			Pool:    testPoolID,
			Id:      serverID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Id != ns.Id {
		t.Errorf("bad id: %q (%q)", s.Id, ns.Id)
	}
        if s.Weight != ns.Weight {
                t.Errorf("bad weigth: %q (%q)", s.Weight, ns.Weight)
        }
        if s.MaxConn != ns.MaxConn {
                t.Errorf("bad maxconn: %q (%q)", s.MaxConn, ns.MaxConn)
        }
        if s.Port != ns.Port {
                t.Errorf("bad port: %q (%q)", s.Port, ns.Port)
        }
	if s.Address != ns.Address {
		t.Errorf("bad address: %q (%q)", s.Address, ns.Address)
	}
        if s.Disabled != ns.Disabled {
                t.Errorf("bad disabled: %t (%t)", s.Disabled, ns.Disabled)
        }
	if s.OverrideHost != ns.OverrideHost {
		t.Errorf("bad override_host: %q (%q)", s.OverrideHost, ns.OverrideHost)
	}

	// Update
	var us *Server
	record(t, fixtureBase+"update_server", func(c *Client) {
		us, err = c.UpdateServer(&UpdateServerInput{
			Service:      testServiceID,
			Pool:         testPoolID,
			Id:           serverID,
			MaxConn:      "888",
			OverrideHost: "www.example.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.MaxConn != "888" {
		t.Errorf("bad maxconn: %q", us.MaxConn)
	}
	if us.OverrideHost != "www.example.com" {
		t.Errorf("bad override_host: %q", us.OverrideHost)
	}

	// Delete
	record(t, fixtureBase+"delete_server", func(c *Client) {
		err = c.DeleteServer(&DeleteServerInput{
			Service: testServiceID,
			Pool:    testPoolID,
			Id:      serverID,
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
		t.Errorf("bad error: %s", err)
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
		t.Errorf("bad error: %s", err)
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
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetServer(&GetServerInput{
		Service: "foo",
		Pool:    "bar",
		Id:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
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
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateServer(&UpdateServerInput{
		Service: "foo",
		Pool:    "bar",
		Id:        "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
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
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteServer(&DeleteServerInput{
		Service: "foo",
		Pool:    "bar",
		Id:        "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
