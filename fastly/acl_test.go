package fastly

import (
	"errors"
	"testing"
)

func TestClient_ACLs(t *testing.T) {
	t.Parallel()

	fixtureBase := "acls/"

	testVersion := CreateTestVersion(t, fixtureBase+"version", TestDeliveryServiceID)

	// Create
	var err error
	var a *ACL
	Record(t, fixtureBase+"create", func(c *Client) {
		a, err = c.CreateACL(&CreateACLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Name:           ToPointer("test_acl"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create with expected error
	var errExpected error
	Record(t, fixtureBase+"create_expected_error", func(c *Client) {
		_, errExpected = c.CreateACL(&CreateACLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if errExpected == nil {
		t.Error("expected API error, got nil")
	}

	// Ensure deleted
	defer func() {
		Record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteACL(&DeleteACLInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				Name:           "test_acl",
			})

			_ = c.DeleteACL(&DeleteACLInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				Name:           "new_test_acl",
			})
		})
	}()

	if *a.Name != "test_acl" {
		t.Errorf("bad name: %q", *a.Name)
	}

	// List
	var as []*ACL
	Record(t, fixtureBase+"list", func(c *Client) {
		as, err = c.ListACLs(&ListACLsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(as) != 1 {
		t.Errorf("bad ACLs: %v", as)
	}

	// Get
	var na *ACL
	Record(t, fixtureBase+"get", func(c *Client) {
		na, err = c.GetACL(&GetACLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Name:           "test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *a.Name != *na.Name {
		t.Errorf("bad name: %q (%q)", *a.Name, *na.Name)
	}

	// Update
	var ua *ACL
	Record(t, fixtureBase+"update", func(c *Client) {
		ua, err = c.UpdateACL(&UpdateACLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Name:           "test_acl",
			NewName:        ToPointer("new_test_acl"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ua.Name != "new_test_acl" {
		t.Errorf("Bad name after update %s", *ua.Name)
	}

	if *a.ACLID != *ua.ACLID {
		t.Errorf("bad ACL id: %q (%q)", *a.ACLID, *ua.ACLID)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteACL(&DeleteACLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Name:           "new_test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListACLs_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListACLs(&ListACLsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListACLs(&ListACLsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateACL_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateACL(&CreateACLInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateACL(&CreateACLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetACL_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetACL(&GetACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetACL(&GetACLInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetACL(&GetACLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateACL_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateACL(&UpdateACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateACL(&UpdateACLInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateACL(&UpdateACLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteACL_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteACL(&DeleteACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteACL(&DeleteACLInput{
		Name:           "test",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteACL(&DeleteACLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
