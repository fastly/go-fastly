package fastly

import "testing"

func TestClient_Acls(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "acls/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Clean up
	defer func() {
		record(t, "acls/cleanup", func(c *Client) {
			c.DeleteAcl(&DeleteAclInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test_acl",
			})

			c.DeleteAcl(&DeleteAclInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new_test_acl",
			})
		})
	}()

	// Create
	var a *Acl
	record(t, "acls/create", func(c *Client) {
		a, err = c.CreateAcl(&CreateAclInput{
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
	var as []*Acl
	record(t, "acls/list", func(c *Client) {
		as, err = c.ListAcls(&ListAclsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(as) != 1 {
		t.Errorf("bad acls: %v", as)
	}

	// Get
	var na *Acl
	record(t, "acls/get", func(c *Client) {
		na, err = c.GetAcl(&GetAclInput{
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
	var ua *Acl
	record(t, "acls/update", func(c *Client) {
		ua, err = c.UpdateAcl(&UpdateAclInput{
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
		err = c.DeleteAcl(&DeleteAclInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new_test_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestClient_ListAcls_validation(t *testing.T) {
	var err error
	_, err = testClient.ListAcls(&ListAclsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListAcls(&ListAclsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateAcl_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateAcl(&CreateAclInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateAcl(&CreateAclInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetAcl_validation(t *testing.T) {
	var err error
	_, err = testClient.GetAcl(&GetAclInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetAcl(&GetAclInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetAcl(&GetAclInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateAcl_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateAcl(&UpdateAclInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateAcl(&UpdateAclInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateAcl(&UpdateAclInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
	_, err = testClient.UpdateAcl(&UpdateAclInput{
		Service: "foo",
		Version: 1,
		Name:    "acl",
		NewName: "",
	})
	if err != ErrMissingNewName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteAcl_validation(t *testing.T) {
	var err error
	err = testClient.DeleteAcl(&DeleteAclInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteAcl(&DeleteAclInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteAcl(&DeleteAclInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
