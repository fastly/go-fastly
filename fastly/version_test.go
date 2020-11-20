package fastly

import (
	"sort"
	"testing"
)

func TestClient_Versions(t *testing.T) {
	t.Parallel()

	var err error

	// Lock because we are creating a new version.
	testVersionLock.Lock()

	// Create
	var v *Version
	record(t, "versions/create", func(c *Client) {
		v, err = c.CreateVersion(&CreateVersionInput{
			ServiceID: testServiceID,
			Comment:   "test comment",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if v.Number == 0 {
		t.Errorf("bad number: %q", v.Number)
	}
	if v.Comment != "test comment" {
		t.Errorf("bad comment: %q", v.Comment)
	}

	// Unlock and let other parallel tests go!
	testVersionLock.Unlock()

	// List
	var vs []*Version
	record(t, "versions/list", func(c *Client) {
		vs, err = c.ListVersions(&ListVersionsInput{
			ServiceID: testServiceID,
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
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if nv.Number == 0 {
		t.Errorf("bad number: %q", nv.Number)
	}
	if nv.Comment != v.Comment {
		t.Errorf("bad comment: %q", v.Comment)
	}

	// Update
	var uv *Version
	record(t, "versions/update", func(c *Client) {
		uv, err = c.UpdateVersion(&UpdateVersionInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
			Comment:        "new comment",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uv.Comment != "new comment" {
		t.Errorf("bad comment: %q", uv.Comment)
	}

	// Lock
	var vl *Version
	record(t, "versions/lock", func(c *Client) {
		vl, err = c.LockVersion(&LockVersionInput{
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
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
			ServiceID:      testServiceID,
			ServiceVersion: v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if cv.Active != false {
		t.Errorf("bad clone: %t", cv.Active)
	}
	if cv.Comment != uv.Comment {
		t.Errorf("bad comment: %q", uv.Comment)
	}
}

func TestClient_SortVersions(t *testing.T) {
	versionsData := []*Version{
		{Number: 1},
		{Number: 201},
		{Number: 10},
		{Number: 2},
		{Number: 197},
	}
	sort.Sort(versionsByNumber(versionsData))
	if versionsData[0].Number != 1 || versionsData[1].Number != 2 || versionsData[2].Number != 10 || versionsData[3].Number != 197 || versionsData[4].Number != 201 {
		t.Fatalf("The sort.Sort did not work properly. Got: %#v\n", versionsData)
	}
}

func TestClient_ListVersions_validation(t *testing.T) {
	var err error
	_, err = testClient.ListVersions(&ListVersionsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateVersion(&CreateVersionInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.GetVersion(&GetVersionInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVersion(&GetVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateVersion(&UpdateVersionInput{
		ServiceID: "",
	})
	if err.Error() != "missing required field 'ServiceID'" {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVersion(&UpdateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err.Error() != "missing required field 'ServiceVersion'" {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ActivateVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.ActivateVersion(&ActivateVersionInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ActivateVersion(&ActivateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeactivateVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.DeactivateVersion(&DeactivateVersionInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.DeactivateVersion(&DeactivateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CloneVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.CloneVersion(&CloneVersionInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CloneVersion(&CloneVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ValidateVersion_validation(t *testing.T) {
	var err error
	_, _, err = testClient.ValidateVersion(&ValidateVersionInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, _, err = testClient.ValidateVersion(&ValidateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_LockVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.LockVersion(&LockVersionInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.LockVersion(&LockVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
