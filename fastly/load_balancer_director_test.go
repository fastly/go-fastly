package fastly

import (
	"errors"
	"testing"
)

func TestClient_Directors(t *testing.T) {
	t.Parallel()

	var tv *Version
	Record(t, "directors/version", func(c *Client) {
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
	Record(t, "directors/create", func(c *Client) {
		b, errBackend = c.CreateBackend(&CreateBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-backend"),
			Address:        ToPointer("integ-test.go-fastly.com"),
			Port:           ToPointer(1234),
			ConnectTimeout: ToPointer(1500),
			OverrideHost:   ToPointer("origin.example.com"),
			SSLCiphers:     ToPointer("DHE-RSA-AES256-SHA:DHE-RSA-CAMELLIA256-SHA:AES256-GCM-SHA384"),
		})
		d, errDirector = c.CreateDirector(&CreateDirectorInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-director"),
			Quorum:         ToPointer(50),
			Type:           ToPointer(DirectorTypeRandom),
			Retries:        ToPointer(5),
		})
		_, errDirectorBackend = c.CreateDirectorBackend(&CreateDirectorBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Director:       "test-director",
			Backend:        *b.Name,
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
		Record(t, "directors/cleanup", func(c *Client) {
			_ = c.DeleteDirectorBackend(&DeleteDirectorBackendInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Director:       *d.Name,
				Backend:        *b.Name,
			})

			_ = c.DeleteBackend(&DeleteBackendInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           *b.Name,
			})

			_ = c.DeleteDirector(&DeleteDirectorInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-director",
			})

			_ = c.DeleteDirector(&DeleteDirectorInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-director",
			})
		})
	}()

	if *d.Name != "test-director" {
		t.Errorf("bad name: %q", *d.Name)
	}
	if *d.Quorum != 50 {
		t.Errorf("bad quorum: %q", *d.Quorum)
	}
	if *d.Type != DirectorTypeRandom {
		t.Errorf("bad type: %d", *d.Type)
	}
	if *d.Retries != 5 {
		t.Errorf("bad retries: %d", *d.Retries)
	}

	// List
	var (
		bs  []*Director
		err error
	)
	Record(t, "directors/list", func(c *Client) {
		bs, err = c.ListDirectors(&ListDirectorsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "directors/get", func(c *Client) {
		nb, err = c.GetDirector(&GetDirectorInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-director",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *d.Name != *nb.Name {
		t.Errorf("bad name: %q (%q)", *d.Name, *nb.Name)
	}
	if *d.Quorum != *nb.Quorum {
		t.Errorf("bad quorum: %q (%q)", *d.Quorum, *nb.Quorum)
	}
	if *d.Type != *nb.Type {
		t.Errorf("bad type: %q (%q)", *d.Type, *nb.Type)
	}
	if *d.Retries != *nb.Retries {
		t.Errorf("bad retries: %q (%q)", *d.Retries, *nb.Retries)
	}
	if len(nb.Backends) == 0 || nb.Backends[0] != *b.Name {
		t.Error("bad backend: expected a backend")
	}

	// Update
	var ub *Director
	Record(t, "directors/update", func(c *Client) {
		ub, err = c.UpdateDirector(&UpdateDirectorInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-director",
			NewName:        ToPointer("new-test-director"),
			Quorum:         ToPointer(100),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ub.Quorum != 100 {
		t.Errorf("bad quorum: %q", *ub.Quorum)
	}

	// Delete
	Record(t, "directors/delete", func(c *Client) {
		err = c.DeleteDirector(&DeleteDirectorInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-director",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDirectors_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListDirectors(&ListDirectorsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListDirectors(&ListDirectorsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDirector_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateDirector(&CreateDirectorInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateDirector(&CreateDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDirector_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetDirector(&GetDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDirector(&GetDirectorInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDirector(&GetDirectorInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDirector_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateDirector(&UpdateDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDirector(&UpdateDirectorInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDirector(&UpdateDirectorInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDirector_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteDirector(&DeleteDirectorInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDirector(&DeleteDirectorInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDirector(&DeleteDirectorInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
