package fastly

import "testing"

func TestClient_Backends(t *testing.T) {
	tv := testVersion(t)

	// Create
	b, err := testClient.CreateBackend(&CreateBackendInput{
		Service:        testServiceID,
		Version:        tv.Number,
		Name:           "test-backend",
		Address:        "integ-test.hashicorp.com",
		Port:           1234,
		ConnectTimeout: 1500,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteBackend(&DeleteBackendInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-backend",
		})

		testClient.DeleteBackend(&DeleteBackendInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-backend",
		})
	}()

	if b.Name != "test-backend" {
		t.Errorf("bad name: %q", b.Name)
	}
	if b.Address != "integ-test.hashicorp.com" {
		t.Errorf("bad address: %q", b.Address)
	}
	if b.Port != 1234 {
		t.Errorf("bad port: %d", b.Port)
	}
	if b.ConnectTimeout != 1500 {
		t.Errorf("bad connect_timeout: %d", b.ConnectTimeout)
	}

	// List
	bs, err := testClient.ListBackends(&ListBackendsInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) < 1 {
		t.Errorf("bad backends: %v", bs)
	}

	// Get
	nb, err := testClient.GetBackend(&GetBackendInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-backend",
	})
	if err != nil {
		t.Fatal(err)
	}
	if b.Name != nb.Name {
		t.Errorf("bad name: %q (%q)", b.Name, nb.Name)
	}
	if b.Address != nb.Address {
		t.Errorf("bad address: %q (%q)", b.Address, nb.Address)
	}
	if b.Port != nb.Port {
		t.Errorf("bad port: %q (%q)", b.Port, nb.Port)
	}
	if b.ConnectTimeout != nb.ConnectTimeout {
		t.Errorf("bad connect_timeout: %q (%q)", b.ConnectTimeout, nb.ConnectTimeout)
	}

	// Update
	ub, err := testClient.UpdateBackend(&UpdateBackendInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-backend",
		NewName: "new-test-backend",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ub.Name != "new-test-backend" {
		t.Errorf("bad name: %q", ub.Name)
	}

	// Delete
	if err := testClient.DeleteBackend(&DeleteBackendInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "new-test-backend",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListBackends_validation(t *testing.T) {
	var err error
	_, err = testClient.ListBackends(&ListBackendsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListBackends(&ListBackendsInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateBackend_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateBackend(&CreateBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateBackend(&CreateBackendInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetBackend_validation(t *testing.T) {
	var err error
	_, err = testClient.GetBackend(&GetBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetBackend(&GetBackendInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetBackend(&GetBackendInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateBackend_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateBackend(&UpdateBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateBackend(&UpdateBackendInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateBackend(&UpdateBackendInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteBackend_validation(t *testing.T) {
	var err error
	err = testClient.DeleteBackend(&DeleteBackendInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteBackend(&DeleteBackendInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteBackend(&DeleteBackendInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
