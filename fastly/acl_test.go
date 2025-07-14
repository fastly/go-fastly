package fastly

import (
	"context"
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
		a, err = c.CreateACL(context.TODO(), &CreateACLInput{
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
		_, errExpected = c.CreateACL(context.TODO(), &CreateACLInput{
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
			_ = c.DeleteACL(context.TODO(), &DeleteACLInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				Name:           "test_acl",
			})

			_ = c.DeleteACL(context.TODO(), &DeleteACLInput{
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
		as, err = c.ListACLs(context.TODO(), &ListACLsInput{
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
		na, err = c.GetACL(context.TODO(), &GetACLInput{
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
		ua, err = c.UpdateACL(context.TODO(), &UpdateACLInput{
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
		err = c.DeleteACL(context.TODO(), &DeleteACLInput{
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
	_, err = TestClient.ListACLs(context.TODO(), &ListACLsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListACLs(context.TODO(), &ListACLsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateACL_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateACL(context.TODO(), &CreateACLInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateACL(context.TODO(), &CreateACLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetACL_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetACL(context.TODO(), &GetACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetACL(context.TODO(), &GetACLInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetACL(context.TODO(), &GetACLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateACL_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateACL(context.TODO(), &UpdateACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateACL(context.TODO(), &UpdateACLInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateACL(context.TODO(), &UpdateACLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteACL_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteACL(context.TODO(), &DeleteACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteACL(context.TODO(), &DeleteACLInput{
		Name:           "test",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteACL(context.TODO(), &DeleteACLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
