package fastly

import "testing"

func TestClient_DirectorBackends(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "director_backends/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var b *DirectorBackend
	record(t, "director_backends/create", func(c *Client) {
		b, err = c.CreateDirectorBackend(&CreateDirectorBackendInput{
			Service:  testServiceID,
			Version:  tv.Number,
			Director: "director",
			Backend:  "backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "director_backends/cleanup", func(c *Client) {
			c.DeleteDirectorBackend(&DeleteDirectorBackendInput{
				Service:  testServiceID,
				Version:  tv.Number,
				Director: "director",
				Backend:  "backend",
			})
		})
	}()

	if b.Director != "director" {
		t.Errorf("bad director: %q", b.Director)
	}
	if b.Backend != "backend" {
		t.Errorf("bad backend: %q", b.Backend)
	}

	// Get
	var nb *DirectorBackend
	record(t, "director_backends/get", func(c *Client) {
		nb, err = c.GetDirectorBackend(&GetDirectorBackendInput{
			Service:  testServiceID,
			Version:  tv.Number,
			Director: "director",
			Backend:  "backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if b.Director != nb.Director {
		t.Errorf("bad director: %q", b.Director)
	}
	if b.Backend != nb.Backend {
		t.Errorf("bad backend: %q", b.Backend)
	}

	// Delete
	record(t, "director_backends/delete", func(c *Client) {
		err = c.DeleteDirectorBackend(&DeleteDirectorBackendInput{
			Service:  testServiceID,
			Version:  tv.Number,
			Director: "director",
			Backend:  "backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateDirectorBackend_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDirectorBackend(&CreateDirectorBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDirectorBackend(&CreateDirectorBackendInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDirectorBackend_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDirectorBackend(&GetDirectorBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirectorBackend(&GetDirectorBackendInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirectorBackend(&GetDirectorBackendInput{
		Service:  "foo",
		Version:  1,
		Director: "",
	})
	if err != ErrMissingDirector {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirectorBackend(&GetDirectorBackendInput{
		Service:  "foo",
		Version:  1,
		Director: "director",
		Backend:  "",
	})
	if err != ErrMissingBackend {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDirectorBackend_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDirectorBackend(&DeleteDirectorBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirectorBackend(&DeleteDirectorBackendInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirectorBackend(&DeleteDirectorBackendInput{
		Service:  "foo",
		Version:  1,
		Director: "",
	})
	if err != ErrMissingDirector {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirectorBackend(&DeleteDirectorBackendInput{
		Service:  "foo",
		Version:  1,
		Director: "director",
		Backend:  "",
	})
	if err != ErrMissingBackend {
		t.Errorf("bad error: %s", err)
	}
}
