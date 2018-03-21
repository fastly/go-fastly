package fastly

import "testing"

func TestClient_ACLs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "acls/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Clean up
	defer func() {
		record(t, "acls/cleanup", func(c *Client) {
			c.DeleteACL(&DeleteACLInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test_acl",
			})

			c.DeleteACL(&DeleteACLInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new_test_acl",
			})
		})
	}()

	// Create
	var a *ACL
	record(t, "acls/create", func(c *Client) {
		a, err = c.CreateACL(&CreateACLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if a.Name != "test_acl" {
		t.Errorf("bad name: %q", a.Name)
	}

	// List
	var as []*ACL
	record(t, "acls/list", func(c *Client) {
		as, err = c.ListACLs(&ListACLsInput{
			Service: testServiceID,
			Version: tv.Number,
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
	record(t, "acls/get", func(c *Client) {
		na, err = c.GetACL(&GetACLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test_acl",
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
	record(t, "acls/update", func(c *Client) {
		ua, err = c.UpdateACL(&UpdateACLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test_acl",
			NewName: "new_test_acl",
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
	record(t, "acls/delete", func(c *Client) {
		err = c.DeleteACL(&DeleteACLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new_test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestClient_ListACLs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListACLs(&ListACLsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListACLs(&ListACLsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateACL_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateACL(&CreateACLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateACL(&CreateACLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetACL_validation(t *testing.T) {
	var err error
	_, err = testClient.GetACL(&GetACLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetACL(&GetACLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetACL(&GetACLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateACL_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateACL(&UpdateACLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateACL(&UpdateACLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateACL(&UpdateACLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
	_, err = testClient.UpdateACL(&UpdateACLInput{
		Service: "foo",
		Version: 1,
		Name:    "acl",
		NewName: "",
	})
	if err != ErrMissingNewName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteACL_validation(t *testing.T) {
	var err error
	err = testClient.DeleteACL(&DeleteACLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteACL(&DeleteACLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteACL(&DeleteACLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
