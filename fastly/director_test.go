package fastly

import (
	"testing"
)

func TestClient_Directors(t *testing.T) {
	t.Parallel()

	var tv *Version
	record(t, "directors/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var (
		b                  *Backend
		d                  *Director
		errBackend         error
		errDirector        error
		errDirectorBackend error
	)
	record(t, "directors/create", func(c *Client) {
		b, errBackend = c.CreateBackend(&CreateBackendInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-backend",
			Address:        "integ-test.go-fastly.com",
			Port:           Uint(1234),
			ConnectTimeout: Uint(1500),
			OverrideHost:   "origin.example.com",
			SSLCiphers:     "DHE-RSA-AES256-SHA:DHE-RSA-CAMELLIA256-SHA:AES256-GCM-SHA384",
		})
		d, errDirector = c.CreateDirector(&CreateDirectorInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-director",
			Quorum:         Uint(50),
			Type:           DirectorTypeRandom,
			Retries:        Uint(5),
		})
		_, errDirectorBackend = c.CreateDirectorBackend(&CreateDirectorBackendInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Director:       "test-director",
			Backend:        b.Name,
		})
	})
	if errBackend != nil {
		t.Fatal(errBackend)
	}
	if errDirector != nil {
		t.Fatal(errDirector)
	}
	if errDirectorBackend != nil {
		t.Fatal(errDirectorBackend)
	}

	// Ensure deleted
	defer func() {
		record(t, "directors/cleanup", func(c *Client) {
			c.DeleteDirectorBackend(&DeleteDirectorBackendInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Director:       d.Name,
				Backend:        b.Name,
			})

			c.DeleteBackend(&DeleteBackendInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           b.Name,
			})

			c.DeleteDirector(&DeleteDirectorInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-director",
			})

			c.DeleteDirector(&DeleteDirectorInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-director",
			})
		})
	}()

	if d.Name != "test-director" {
		t.Errorf("bad name: %q", d.Name)
	}
	if d.Quorum != 50 {
		t.Errorf("bad quorum: %q", d.Quorum)
	}
	if d.Type != DirectorTypeRandom {
		t.Errorf("bad type: %d", d.Type)
	}
	if d.Retries != 5 {
		t.Errorf("bad retries: %d", d.Retries)
	}

	// List
	var (
		bs  []*Director
		err error
	)
	record(t, "directors/list", func(c *Client) {
		bs, err = c.ListDirectors(&ListDirectorsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) < 1 {
		t.Errorf("bad directors: %v", bs)
	}

	// Get
	var nb *Director
	record(t, "directors/get", func(c *Client) {
		nb, err = c.GetDirector(&GetDirectorInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-director",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != nb.Name {
		t.Errorf("bad name: %q (%q)", d.Name, nb.Name)
	}
	if d.Quorum != nb.Quorum {
		t.Errorf("bad quorum: %q (%q)", d.Quorum, nb.Quorum)
	}
	if d.Type != nb.Type {
		t.Errorf("bad type: %q (%q)", d.Type, nb.Type)
	}
	if d.Retries != nb.Retries {
		t.Errorf("bad retries: %q (%q)", d.Retries, nb.Retries)
	}
	if len(nb.Backends) == 0 || nb.Backends[0] != b.Name {
		t.Error("bad backend: expected a backend")
	}

	// Update
	var ub *Director
	record(t, "directors/update", func(c *Client) {
		ub, err = c.UpdateDirector(&UpdateDirectorInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-director",
			NewName:        String("new-test-director"),
			Quorum:         Uint(100),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ub.Quorum != 100 {
		t.Errorf("bad quorum: %q", ub.Quorum)
	}

	// Delete
	record(t, "directors/delete", func(c *Client) {
		err = c.DeleteDirector(&DeleteDirectorInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-director",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDirectors_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDirectors(&ListDirectorsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDirectors(&ListDirectorsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDirector_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDirector(&CreateDirectorInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDirector(&CreateDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDirector_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDirector(&GetDirectorInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirector(&GetDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirector(&GetDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDirector_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDirector(&UpdateDirectorInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDirector(&UpdateDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDirector(&UpdateDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDirector_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDirector(&DeleteDirectorInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirector(&DeleteDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirector(&DeleteDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
