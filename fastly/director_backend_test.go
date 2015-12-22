package fastly

import "testing"

func TestClient_DirectorBackends(t *testing.T) {
	tv := testVersion(t)

	// Create
	b, err := testClient.CreateDirectorBackend(&CreateDirectorBackendInput{
		Service:  testServiceID,
		Version:  tv.Number,
		Director: "director",
		Backend:  "backend",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteDirectorBackend(&DeleteDirectorBackendInput{
			Service:  testServiceID,
			Version:  tv.Number,
			Director: "director",
			Backend:  "backend",
		})
	}()

	if b.Director != "director" {
		t.Errorf("bad director: %q", b.Director)
	}
	if b.Backend != "backend" {
		t.Errorf("bad backend: %q", b.Backend)
	}

	// Get
	nb, err := testClient.GetDirectorBackend(&GetDirectorBackendInput{
		Service:  testServiceID,
		Version:  tv.Number,
		Director: "director",
		Backend:  "backend",
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
	if err := testClient.DeleteDirectorBackend(&DeleteDirectorBackendInput{
		Service:  testServiceID,
		Version:  tv.Number,
		Director: "director",
		Backend:  "backend",
	}); err != nil {
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
		Version: "",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirectorBackend(&GetDirectorBackendInput{
		Service:  "foo",
		Version:  "1",
		Director: "",
	})
	if err != ErrMissingDirector {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirectorBackend(&GetDirectorBackendInput{
		Service:  "foo",
		Version:  "1",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirectorBackend(&DeleteDirectorBackendInput{
		Service:  "foo",
		Version:  "1",
		Director: "",
	})
	if err != ErrMissingDirector {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirectorBackend(&DeleteDirectorBackendInput{
		Service:  "foo",
		Version:  "1",
		Director: "director",
		Backend:  "",
	})
	if err != ErrMissingBackend {
		t.Errorf("bad error: %s", err)
	}
}
