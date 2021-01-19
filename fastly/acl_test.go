package fastly

import (
	"testing"
)

func TestClient_ACLs(t *testing.T) {
	t.Parallel()

	fixtureBase := "acls/"

	testVersion := createTestVersion(t, fixtureBase+"version", testServiceID)

	// Create
	var err error
	var a *ACL
	record(t, fixtureBase+"create", func(c *Client) {
		a, err = c.CreateACL(&CreateACLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           "test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteACL(&DeleteACLInput{
				ServiceID:      testServiceID,
				ServiceVersion: testVersion.Number,
				Name:           "test_acl",
			})

			c.DeleteACL(&DeleteACLInput{
				ServiceID:      testServiceID,
				ServiceVersion: testVersion.Number,
				Name:           "new_test_acl",
			})
		})
	}()

	if a.Name != "test_acl" {
		t.Errorf("bad name: %q", a.Name)
	}

	// List
	var as []*ACL
	record(t, fixtureBase+"list", func(c *Client) {
		as, err = c.ListACLs(&ListACLsInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
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
	record(t, fixtureBase+"get", func(c *Client) {
		na, err = c.GetACL(&GetACLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           "test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != na.Name {
		t.Errorf("bad name: %q (%q)", a.Name, na.Name)
	}

	// Update
	var ua *ACL
	record(t, fixtureBase+"update", func(c *Client) {
		ua, err = c.UpdateACL(&UpdateACLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           "test_acl",
			NewName:        "new_test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ua.Name != "new_test_acl" {
		t.Errorf("Bad name after update %s", ua.Name)
	}

	if a.ID != ua.ID {
		t.Errorf("bad ACL id: %q (%q)", a.ID, ua.ID)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteACL(&DeleteACLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           "new_test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestClient_ListACLs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListACLs(&ListACLsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListACLs(&ListACLsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateACL_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateACL(&CreateACLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateACL(&CreateACLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetACL_validation(t *testing.T) {
	var err error
	_, err = testClient.GetACL(&GetACLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetACL(&GetACLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetACL(&GetACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateACL_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateACL(&UpdateACLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateACL(&UpdateACLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateACL(&UpdateACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
	_, err = testClient.UpdateACL(&UpdateACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "acl",
		NewName:        "",
	})
	if err != ErrMissingNewName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteACL_validation(t *testing.T) {
	var err error
	err = testClient.DeleteACL(&DeleteACLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteACL(&DeleteACLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteACL(&DeleteACLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
