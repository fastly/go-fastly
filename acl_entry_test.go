package fastly

import "testing"

func TestClient_AclEntries(t *testing.T) {

	var err error
	var tv *Version
	record(t, "acls/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create ACL before adding entries to it
	var a *Acl
	record(t, "acl_entries/create_acl", func(c *Client) {
		a, err = c.CreateAcl(&CreateAclInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "entry_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Clean up helper Acl
	defer func() {
		record(t, "acl_entries/cleanup", func(c *Client) {
			err = c.DeleteAcl(&DeleteAclInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "entry_acl",
			})
		})
	}()

	// Create
	var e *AclEntry
	record(t, "acl_entries/create", func(c *Client) {
		e, err = c.CreateAclEntry(&CreateAclEntryInput{
			Service: testServiceID,
			Acl:     a.ID,
			IP:      "10.0.0.3",
			Subnet:  "8",
			Negated: false,
			Comment: "test entry",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if e.IP != "10.0.0.3" {
		t.Errorf("bad IP: %q", e.IP)
	}

	if e.Subnet != "8" {
		t.Errorf("Bad subnet: %q", e.Subnet)
	}

	if e.Negated != false {
		t.Errorf("Bad negated flag: %t", e.Negated)
	}

	if e.Comment != "test entry" {
		t.Errorf("Bad comment: %q", e.Comment)
	}

	// List
	var es []*AclEntry
	record(t, "acl_entries/list", func(c *Client) {
		es, err = c.ListAclEntries(&ListAclEntriesInput{
			Service: testServiceID,
			Acl:     a.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(es) != 1 {
		t.Errorf("Bad entries: %v", es)
	}

	// Get
	var ne *AclEntry
	record(t, "acl_entries/get", func(c *Client) {
		ne, err = c.GetAclEntry(&GetAclEntryInput{
			Service: testServiceID,
			Acl:     a.ID,
			ID:      e.ID,
		})
	})

	if e.IP != ne.IP {
		t.Errorf("bad IP: %v", ne.IP)
	}

	if e.Subnet != ne.Subnet {
		t.Errorf("bad subnet: %v", ne.Subnet)
	}

	if e.Negated != ne.Negated {
		t.Errorf("bad Negated flag: %v", ne.Negated)
	}

	if e.Comment != ne.Comment {
		t.Errorf("bad comment: %v", ne.Comment)
	}

	// Update
	var ue *AclEntry
	record(t, "acl_entries/update", func(c *Client) {
		ue, err = c.UpdateAclEntry(&UpdateAclEntryInput{
			Service: testServiceID,
			Acl:     a.ID,
			ID:      e.ID,
			IP:      "10.0.0.4",
			Negated: true,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if ue.IP != "10.0.0.4" {
		t.Errorf("bad IP: %q", ue.IP)
	}
	if e.Subnet != ue.Subnet {
		t.Errorf("bad subnet: %v", ne.Subnet)
	}

	if ue.Negated != true {
		t.Errorf("bad Negated flag: %v", ne.Negated)
	}

	if e.Comment != ue.Comment {
		t.Errorf("bad comment: %v", ne.Comment)
	}

	// Delete
	record(t, "acl_entries/delete", func(c *Client) {
		err = c.DeleteAclEntry(&DeleteAclEntryInput{
			Service: testServiceID,
			Acl:     a.ID,
			ID:      e.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestClient_ListAclEntries_validation(t *testing.T) {
	var err error
	_, err = testClient.ListAclEntries(&ListAclEntriesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListAclEntries(&ListAclEntriesInput{
		Service: "foo",
		Acl:     "",
	})
	if err != ErrMissingAcl {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateAclEntry_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateAclEntry(&CreateAclEntryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateAclEntry(&CreateAclEntryInput{
		Service: "foo",
		Acl:     "",
	})
	if err != ErrMissingAcl {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetAclEntry_validation(t *testing.T) {
	var err error
	_, err = testClient.GetAclEntry(&GetAclEntryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetAclEntry(&GetAclEntryInput{
		Service: "foo",
		Acl:     "",
	})
	if err != ErrMissingAcl {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetAclEntry(&GetAclEntryInput{
		Service: "foo",
		Acl:     "acl",
		ID:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateAclEntry_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateAclEntry(&UpdateAclEntryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateAclEntry(&UpdateAclEntryInput{
		Service: "foo",
		Acl:     "",
	})
	if err != ErrMissingAcl {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateAclEntry(&UpdateAclEntryInput{
		Service: "foo",
		Acl:     "acl",
		ID:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteAclEntry_validation(t *testing.T) {
	var err error
	err = testClient.DeleteAclEntry(&DeleteAclEntryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteAclEntry(&DeleteAclEntryInput{
		Service: "foo",
		Acl:     "",
	})
	if err != ErrMissingAcl {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteAclEntry(&DeleteAclEntryInput{
		Service: "foo",
		Acl:     "acl",
		ID:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
