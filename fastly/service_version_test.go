package fastly

import (
	"context"
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
		v, err = c.CreateVersion(context.TODO(), &CreateVersionInput{
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
		vs, err = c.ListVersions(context.TODO(), &ListVersionsInput{
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
		nv, err = c.GetVersion(context.TODO(), &GetVersionInput{
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
		uv, err = c.UpdateVersion(context.TODO(), &UpdateVersionInput{
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
		vl, err = c.LockVersion(context.TODO(), &LockVersionInput{
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
		cv, err = c.CloneVersion(context.TODO(), &CloneVersionInput{
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

func TestClient_Versions_Compute(t *testing.T) {
	t.Parallel()

	var err error

	// Lock because we are creating a new version.
	testVersionLock.Lock()

	// Create
	var v *Version
	Record(t, "versions/compute/create", func(c *Client) {
		v, err = c.CreateVersion(context.TODO(), &CreateVersionInput{
			ServiceID: TestComputeServiceID,
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
	Record(t, "versions/compute/list", func(c *Client) {
		vs, err = c.ListVersions(context.TODO(), &ListVersionsInput{
			ServiceID: TestComputeServiceID,
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
	Record(t, "versions/compute/get", func(c *Client) {
		nv, err = c.GetVersion(context.TODO(), &GetVersionInput{
			ServiceID:      TestComputeServiceID,
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
	Record(t, "versions/compute/update", func(c *Client) {
		uv, err = c.UpdateVersion(context.TODO(), &UpdateVersionInput{
			ServiceID:      TestComputeServiceID,
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
	Record(t, "versions/compute/lock", func(c *Client) {
		vl, err = c.LockVersion(context.TODO(), &LockVersionInput{
			ServiceID:      TestComputeServiceID,
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
	Record(t, "versions/compute/clone", func(c *Client) {
		cv, err = c.CloneVersion(context.TODO(), &CloneVersionInput{
			ServiceID:      TestComputeServiceID,
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
	_, err = TestClient.ListVersions(context.TODO(), &ListVersionsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateVersion(context.TODO(), &CreateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetVersion(context.TODO(), &GetVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetVersion(context.TODO(), &GetVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateVersion(context.TODO(), &UpdateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateVersion(context.TODO(), &UpdateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ActivateVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.ActivateVersion(context.TODO(), &ActivateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ActivateVersion(context.TODO(), &ActivateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeactivateVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.DeactivateVersion(context.TODO(), &DeactivateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.DeactivateVersion(context.TODO(), &DeactivateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CloneVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.CloneVersion(context.TODO(), &CloneVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CloneVersion(context.TODO(), &CloneVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ValidateVersion_validation(t *testing.T) {
	var err error
	_, _, err = TestClient.ValidateVersion(context.TODO(), &ValidateVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, _, err = TestClient.ValidateVersion(context.TODO(), &ValidateVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_LockVersion_validation(t *testing.T) {
	var err error
	_, err = TestClient.LockVersion(context.TODO(), &LockVersionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.LockVersion(context.TODO(), &LockVersionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
