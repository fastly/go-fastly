package fastly

import "testing"

func TestClient_Versions(t *testing.T) {
	t.Parallel()

	var err error

	// Lock because we are creating a new version.
	testVersionLock.Lock()

	// Create
	var v *Version
	record(t, "versions/create", func(c *Client) {
		v, err = c.CreateVersion(&CreateVersionInput{
			Service: testServiceID,
		})
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
	var vs []*Version
	record(t, "versions/list", func(c *Client) {
		vs, err = c.ListVersions(&ListVersionsInput{
			Service: testServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) < 1 {
		t.Errorf("bad services: %v", vs)
	}

	// Get
	var nv *Version
	record(t, "versions/get", func(c *Client) {
		nv, err = c.GetVersion(&GetVersionInput{
			Service: testServiceID,
			Version: v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if nv.Number == "" {
		t.Errorf("bad number: %q", nv.Number)
	}

	// Update
	var uv *Version
	record(t, "versions/update", func(c *Client) {
		uv, err = c.UpdateVersion(&UpdateVersionInput{
			Service: testServiceID,
			Version: v.Number,
			Comment: "new comment",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uv.Comment != "new comment" {
		t.Errorf("bad name: %q", uv.Comment)
	}

	// Lock
	var vl *Version
	record(t, "versions/lock", func(c *Client) {
		vl, err = c.LockVersion(&LockVersionInput{
			Service: testServiceID,
			Version: v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if vl.Locked != true {
		t.Errorf("bad lock: %t", vl.Locked)
	}

	// Clone
	var cv *Version
	record(t, "versions/clone", func(c *Client) {
		cv, err = c.CloneVersion(&CloneVersionInput{
			Service: testServiceID,
			Version: v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if cv.Active != false {
		t.Errorf("bad clone: %t", cv.Active)
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
