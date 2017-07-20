package fastly

import "testing"

func TestClient_ACLEntries(t *testing.T) {

	var err error
	var tv *Version
	record(t, "acls/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create ACL before adding entries to it
	var a *ACL
	record(t, "acl_entries/create_acl", func(c *Client) {
		a, err = c.CreateACL(&CreateACLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "entry_acl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Clean up helper ACL
	defer func() {
		record(t, "acl_entries/cleanup", func(c *Client) {
			err = c.DeleteACL(&DeleteACLInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "entry_acl",
			})
		})
	}()

	// Create
	var e *ACLEntry
	record(t, "acl_entries/create", func(c *Client) {
		e, err = c.CreateACLEntry(&CreateACLEntryInput{
			Service: testServiceID,
			ACL:     a.ID,
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
	var es []*ACLEntry
	record(t, "acl_entries/list", func(c *Client) {
		es, err = c.ListACLEntries(&ListACLEntriesInput{
			Service: testServiceID,
			ACL:     a.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(es) != 1 {
		t.Errorf("Bad entries: %v", es)
	}

	// Get
	var ne *ACLEntry
	record(t, "acl_entries/get", func(c *Client) {
		ne, err = c.GetACLEntry(&GetACLEntryInput{
			Service: testServiceID,
			ACL:     a.ID,
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
	var ue *ACLEntry
	record(t, "acl_entries/update", func(c *Client) {
		ue, err = c.UpdateACLEntry(&UpdateACLEntryInput{
			Service: testServiceID,
			ACL:     a.ID,
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
		err = c.DeleteACLEntry(&DeleteACLEntryInput{
			Service: testServiceID,
			ACL:     a.ID,
			ID:      e.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestClient_ListACLEntries_validation(t *testing.T) {
	var err error
	_, err = testClient.ListACLEntries(&ListACLEntriesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListACLEntries(&ListACLEntriesInput{
		Service: "foo",
		ACL:     "",
	})
	if err != ErrMissingACL {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateACLEntry_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateACLEntry(&CreateACLEntryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateACLEntry(&CreateACLEntryInput{
		Service: "foo",
		ACL:     "",
	})
	if err != ErrMissingACL {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetACLEntry_validation(t *testing.T) {
	var err error
	_, err = testClient.GetACLEntry(&GetACLEntryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetACLEntry(&GetACLEntryInput{
		Service: "foo",
		ACL:     "",
	})
	if err != ErrMissingACL {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetACLEntry(&GetACLEntryInput{
		Service: "foo",
		ACL:     "acl",
		ID:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateACLEntry_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateACLEntry(&UpdateACLEntryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateACLEntry(&UpdateACLEntryInput{
		Service: "foo",
		ACL:     "",
	})
	if err != ErrMissingACL {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateACLEntry(&UpdateACLEntryInput{
		Service: "foo",
		ACL:     "acl",
		ID:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteACLEntry_validation(t *testing.T) {
	var err error
	err = testClient.DeleteACLEntry(&DeleteACLEntryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteACLEntry(&DeleteACLEntryInput{
		Service: "foo",
		ACL:     "",
	})
	if err != ErrMissingACL {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteACLEntry(&DeleteACLEntryInput{
		Service: "foo",
		ACL:     "acl",
		ID:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
