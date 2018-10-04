package fastly

import "testing"

func TestClient_Directors(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "directors/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var b *Director
	record(t, "directors/create", func(c *Client) {
		b, err = c.CreateDirector(&CreateDirectorInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-director",
			Quorum:  50,
			Type:    DirectorTypeRandom,
			Retries: 5,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "directors/cleanup", func(c *Client) {
			c.DeleteDirector(&DeleteDirectorInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-director",
			})

			c.DeleteDirector(&DeleteDirectorInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-director",
			})
		})
	}()

	if b.Name != "test-director" {
		t.Errorf("bad name: %q", b.Name)
	}
	if b.Quorum != 50 {
		t.Errorf("bad quorum: %q", b.Quorum)
	}
	if b.Type != DirectorTypeRandom {
		t.Errorf("bad type: %d", b.Type)
	}
	if b.Retries != 5 {
		t.Errorf("bad retries: %d", b.Retries)
	}

	// List
	var bs []*Director
	record(t, "directors/list", func(c *Client) {
		bs, err = c.ListDirectors(&ListDirectorsInput{
			Service: testServiceID,
			Version: tv.Number,
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
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-director",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if b.Name != nb.Name {
		t.Errorf("bad name: %q (%q)", b.Name, nb.Name)
	}
	if b.Quorum != nb.Quorum {
		t.Errorf("bad quorum: %q (%q)", b.Quorum, nb.Quorum)
	}
	if b.Type != nb.Type {
		t.Errorf("bad type: %q (%q)", b.Type, nb.Type)
	}
	if b.Retries != nb.Retries {
		t.Errorf("bad retries: %q (%q)", b.Retries, nb.Retries)
	}

	// Update
	var ub *Director
	record(t, "directors/update", func(c *Client) {
		ub, err = c.UpdateDirector(&UpdateDirectorInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-director",
			NewName: "new-test-director",
			Quorum:  100,
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
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-director",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDirectors_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDirectors(&ListDirectorsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDirectors(&ListDirectorsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDirector_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDirector(&CreateDirectorInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDirector(&CreateDirectorInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDirector_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDirector(&GetDirectorInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirector(&GetDirectorInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDirector(&GetDirectorInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDirector_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDirector(&UpdateDirectorInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDirector(&UpdateDirectorInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDirector(&UpdateDirectorInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDirector_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDirector(&DeleteDirectorInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirector(&DeleteDirectorInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDirector(&DeleteDirectorInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
