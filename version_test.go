package fastly

import "testing"

func TestClient_Versions(t *testing.T) {
	t.Parallel()

	// Lock because we are creating a new version.
	testVersionLock.Lock()

	// Create
	v, err := testClient.CreateVersion(&CreateVersionInput{
		Service: testServiceID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if v.Number == "" {
		t.Errorf("bad number: %q", v.Number)
	}

	// Unlock and let other parallel tests go!
	testVersionLock.Unlock()

	// List
	vs, err := testClient.ListVersions(&ListVersionsInput{
		Service: testServiceID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) < 1 {
		t.Errorf("bad services: %v", vs)
	}

	// Get
	nv, err := testClient.GetVersion(&GetVersionInput{
		Service: testServiceID,
		Version: v.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if nv.Number == "" {
		t.Errorf("bad number: %q", nv.Number)
	}

	// Update
	uv, err := testClient.UpdateVersion(&UpdateVersionInput{
		Service: testServiceID,
		Version: v.Number,
		Comment: "new comment",
	})
	if err != nil {
		t.Fatal(err)
	}
	if uv.Comment != "new comment" {
		t.Errorf("bad name: %q", uv.Comment)
	}

	// Lock
	v, err = testClient.LockVersion(&LockVersionInput{
		Service: testServiceID,
		Version: v.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if v.Locked != true {
		t.Errorf("bad lock: %d", v.Locked)
	}

	// Clone
	nv, err = testClient.CloneVersion(&CloneVersionInput{
		Service: testServiceID,
		Version: v.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if v.Active != false {
		t.Errorf("bad clone: %d", v.Active)
	}
}

func TestClient_ListVersions_validation(t *testing.T) {
	var err error
	_, err = testClient.ListVersions(&ListVersionsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateVersion(&CreateVersionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.GetVersion(&GetVersionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVersion(&GetVersionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateVersion(&UpdateVersionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVersion(&UpdateVersionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ActivateVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.ActivateVersion(&ActivateVersionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ActivateVersion(&ActivateVersionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeactivateVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.DeactivateVersion(&DeactivateVersionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.DeactivateVersion(&DeactivateVersionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CloneVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.CloneVersion(&CloneVersionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CloneVersion(&CloneVersionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ValidateVersion_validation(t *testing.T) {
	var err error
	_, _, err = testClient.ValidateVersion(&ValidateVersionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, _, err = testClient.ValidateVersion(&ValidateVersionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_LockVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.LockVersion(&LockVersionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.LockVersion(&LockVersionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}
