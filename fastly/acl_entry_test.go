package fastly

import "testing"

func TestClient_ACLEntries(t *testing.T) {

	fixtureBase := "acl_entries/"
	nameSuffix := "ACLEntries"

	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", testService.ID)

	testVersion := createTestVersion(t, fixtureBase+"version", testService.ID)

	testACL := createTestACL(t, fixtureBase+"acl", testService.ID, testVersion.Number, nameSuffix)
	defer deleteTestACL(t, testACL, fixtureBase+"delete_acl")

	// Create
	var err error
	var e *ACLEntry
	record(t, fixtureBase+"create", func(c *Client) {
		e, err = c.CreateACLEntry(&CreateACLEntryInput{
			Service: testService.ID,
			ACL:     testACL.ID,
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
	record(t, fixtureBase+"list", func(c *Client) {
		es, err = c.ListACLEntries(&ListACLEntriesInput{
			Service: testService.ID,
			ACL:     testACL.ID,
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
	record(t, fixtureBase+"get", func(c *Client) {
		ne, err = c.GetACLEntry(&GetACLEntryInput{
			Service: testService.ID,
			ACL:     testACL.ID,
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
	record(t, fixtureBase+"update", func(c *Client) {
		ue, err = c.UpdateACLEntry(&UpdateACLEntryInput{
			Service: testService.ID,
			ACL:     testACL.ID,
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
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteACLEntry(&DeleteACLEntryInput{
			Service: testService.ID,
			ACL:     testACL.ID,
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

func TestClient_BatchModifyACLEntries_validation(t *testing.T) {
	var err error
	err = testClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}
	err = testClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{
		Service: "foo",
		ACL:     "",
	})
	if err != ErrMissingACL {
		t.Errorf("bad error: %s", err)
	}

	oversizedACLEntries := make([]*BatchACLEntry, BatchModifyMaximumOperations+1)
	err = testClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{
		Service: "foo",
		ACL:     "bar",
		Entries: oversizedACLEntries,
	})
	if err != ErrBatchUpdateMaximumOperationsExceeded {
		t.Errorf("bad error: %s", err)
	}

}
