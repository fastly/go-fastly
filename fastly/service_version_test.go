package fastly

import (
	"errors"
	"testing"
)

func TestClient_Versions(t *testing.T) {
	t.Parallel()

	var err error

	// Lock because we are creating a new version.
	testVersionLock.Lock()

	// Create
	var v *Version
	Record(t, "versions/create", func(c *Client) {
		v, err = c.CreateVersion(&CreateVersionInput{
			ServiceID: TestDeliveryServiceID,
			Comment:   ToPointer("test comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *v.Number == 0 {
		t.Errorf("bad number: %q", *v.Number)
	}
	if *v.Comment != "test comment" {
		t.Errorf("bad comment: %q", *v.Comment)
	}

	// Unlock and let other parallel tests go!
	testVersionLock.Unlock()

	// List
	var vs []*Version
	Record(t, "versions/list", func(c *Client) {
		vs, err = c.ListVersions(&ListVersionsInput{
			ServiceID: TestDeliveryServiceID,
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
	Record(t, "versions/get", func(c *Client) {
		nv, err = c.GetVersion(&GetVersionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *nv.Number == 0 {
		t.Errorf("bad number: %q", *nv.Number)
	}
	if *nv.Comment != *v.Comment {
		t.Errorf("bad comment: %q", *v.Comment)
	}

	// Update
	var uv *Version
	Record(t, "versions/update", func(c *Client) {
		uv, err = c.UpdateVersion(&UpdateVersionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
			Comment:        ToPointer("new comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uv.Comment != "new comment" {
		t.Errorf("bad comment: %q", *uv.Comment)
	}

	// Lock
	var vl *Version
	Record(t, "versions/lock", func(c *Client) {
		vl, err = c.LockVersion(&LockVersionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !*vl.Locked {
		t.Errorf("bad lock: %t", *vl.Locked)
	}

	// Clone
	var cv *Version
	Record(t, "versions/clone", func(c *Client) {
		cv, err = c.CloneVersion(&CloneVersionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *v.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *cv.Active {
		t.Errorf("bad clone: %t", *cv.Active)
	}
	if *cv.Comment != *uv.Comment {
		t.Errorf("bad comment: %q", *uv.Comment)
	}
}

func TestClient_ListVersions_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListVersions(&ListVersionsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateVersion(&CreateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetVersion(&GetVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetVersion(&GetVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateVersion(&UpdateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateVersion(&UpdateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ActivateVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.ActivateVersion(&ActivateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ActivateVersion(&ActivateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeactivateVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.DeactivateVersion(&DeactivateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.DeactivateVersion(&DeactivateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CloneVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.CloneVersion(&CloneVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CloneVersion(&CloneVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ValidateVersion_validation(t *testing.T) {
	var err error
	_, _, err = TestClient.ValidateVersion(&ValidateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, _, err = TestClient.ValidateVersion(&ValidateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_LockVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.LockVersion(&LockVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.LockVersion(&LockVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
