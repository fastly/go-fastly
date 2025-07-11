package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_DirectorBackends(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "director_backends/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var b *DirectorBackend
	Record(t, "director_backends/create", func(c *Client) {
		b, err = c.CreateDirectorBackend(context.TODO(), &CreateDirectorBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Director:       "director",
			Backend:        "backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "director_backends/cleanup", func(c *Client) {
			_ = c.DeleteDirectorBackend(context.TODO(), &DeleteDirectorBackendInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Director:       "director",
				Backend:        "backend",
			})
		})
	}()

	if *b.Director != "director" {
		t.Errorf("bad director: %q", *b.Director)
	}
	if *b.Backend != "backend" {
		t.Errorf("bad backend: %q", *b.Backend)
	}

	// Get
	var nb *DirectorBackend
	Record(t, "director_backends/get", func(c *Client) {
		nb, err = c.GetDirectorBackend(context.TODO(), &GetDirectorBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Director:       "director",
			Backend:        "backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *b.Director != *nb.Director {
		t.Errorf("bad director: %q", *b.Director)
	}
	if *b.Backend != *nb.Backend {
		t.Errorf("bad backend: %q", *b.Backend)
	}

	// Delete
	Record(t, "director_backends/delete", func(c *Client) {
		err = c.DeleteDirectorBackend(context.TODO(), &DeleteDirectorBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Director:       "director",
			Backend:        "backend",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateDirectorBackend_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateDirectorBackend(context.TODO(), &CreateDirectorBackendInput{})
	if !errors.Is(err, ErrMissingBackend) {
		t.Errorf("bad error: %s", err)
	}
	_, err = TestClient.CreateDirectorBackend(context.TODO(), &CreateDirectorBackendInput{
		Backend: "foo",
	})
	if !errors.Is(err, ErrMissingDirector) {
		t.Errorf("bad error: %s", err)
	}
	_, err = TestClient.CreateDirectorBackend(context.TODO(), &CreateDirectorBackendInput{
		Backend:  "foo",
		Director: "bar",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
	_, err = TestClient.CreateDirectorBackend(context.TODO(), &CreateDirectorBackendInput{
		Backend:   "foo",
		Director:  "bar",
		ServiceID: "baz",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDirectorBackend_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetDirectorBackend(context.TODO(), &GetDirectorBackendInput{})
	if !errors.Is(err, ErrMissingBackend) {
		t.Errorf("bad error: %s", err)
	}
	_, err = TestClient.GetDirectorBackend(context.TODO(), &GetDirectorBackendInput{
		Backend: "foo",
	})
	if !errors.Is(err, ErrMissingDirector) {
		t.Errorf("bad error: %s", err)
	}
	_, err = TestClient.GetDirectorBackend(context.TODO(), &GetDirectorBackendInput{
		Backend:  "foo",
		Director: "bar",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
	_, err = TestClient.GetDirectorBackend(context.TODO(), &GetDirectorBackendInput{
		Backend:   "foo",
		Director:  "bar",
		ServiceID: "baz",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDirectorBackend_validation(t *testing.T) {
	var err error
	err = TestClient.DeleteDirectorBackend(context.TODO(), &DeleteDirectorBackendInput{})
	if !errors.Is(err, ErrMissingBackend) {
		t.Errorf("bad error: %s", err)
	}
	err = TestClient.DeleteDirectorBackend(context.TODO(), &DeleteDirectorBackendInput{
		Backend: "foo",
	})
	if !errors.Is(err, ErrMissingDirector) {
		t.Errorf("bad error: %s", err)
	}
	err = TestClient.DeleteDirectorBackend(context.TODO(), &DeleteDirectorBackendInput{
		Backend:  "foo",
		Director: "bar",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
	err = TestClient.DeleteDirectorBackend(context.TODO(), &DeleteDirectorBackendInput{
		Backend:   "foo",
		Director:  "bar",
		ServiceID: "baz",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
