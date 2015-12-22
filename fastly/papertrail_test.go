package fastly

import "testing"

func TestClient_Papertrails(t *testing.T) {
	t.Parallel()

	tv := testVersion(t)

	// Create
	p, err := testClient.CreatePapertrail(&CreatePapertrailInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-papertrail",
		Address: "integ-test.go-fastly.com",
		Port:    1234,
		Format:  "format",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeletePapertrail(&DeletePapertrailInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-papertrail",
		})

		testClient.DeletePapertrail(&DeletePapertrailInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-papertrail",
		})
	}()

	if p.Name != "test-papertrail" {
		t.Errorf("bad name: %q", p.Name)
	}
	if p.Address != "integ-test.go-fastly.com" {
		t.Errorf("bad address: %q", p.Address)
	}
	if p.Port != 1234 {
		t.Errorf("bad port: %q", p.Port)
	}
	if p.Format != "format" {
		t.Errorf("bad format: %q", p.Format)
	}

	// List
	ps, err := testClient.ListPapertrails(&ListPapertrailsInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) < 1 {
		t.Errorf("bad papertrails: %v", ps)
	}

	// Get
	np, err := testClient.GetPapertrail(&GetPapertrailInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-papertrail",
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != np.Name {
		t.Errorf("bad name: %q", p.Name)
	}
	if p.Address != np.Address {
		t.Errorf("bad address: %q", p.Address)
	}
	if p.Port != np.Port {
		t.Errorf("bad port: %q", p.Port)
	}
	if p.Format != np.Format {
		t.Errorf("bad format: %q", p.Format)
	}

	// Update
	up, err := testClient.UpdatePapertrail(&UpdatePapertrailInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-papertrail",
		NewName: "new-test-papertrail",
	})
	if err != nil {
		t.Fatal(err)
	}
	if up.Name != "new-test-papertrail" {
		t.Errorf("bad name: %q", up.Name)
	}

	// Delete
	if err := testClient.DeletePapertrail(&DeletePapertrailInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "new-test-papertrail",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPapertrails_validation(t *testing.T) {
	var err error
	_, err = testClient.ListPapertrails(&ListPapertrailsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListPapertrails(&ListPapertrailsInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreatePapertrail_validation(t *testing.T) {
	var err error
	_, err = testClient.CreatePapertrail(&CreatePapertrailInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreatePapertrail(&CreatePapertrailInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetPapertrail_validation(t *testing.T) {
	var err error
	_, err = testClient.GetPapertrail(&GetPapertrailInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPapertrail(&GetPapertrailInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPapertrail(&GetPapertrailInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePapertrail_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdatePapertrail(&UpdatePapertrailInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePapertrail(&UpdatePapertrailInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePapertrail(&UpdatePapertrailInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeletePapertrail_validation(t *testing.T) {
	var err error
	err = testClient.DeletePapertrail(&DeletePapertrailInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePapertrail(&DeletePapertrailInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeletePapertrail(&DeletePapertrailInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
