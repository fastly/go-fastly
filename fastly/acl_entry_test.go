package fastly

import (
	"testing"
)

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
			ServiceID: testService.ID,
			ACLID:     testACL.ID,
			IP:        ToPointer("10.0.0.3"),
			Subnet:    ToPointer(8),
			Negated:   ToPointer(Compatibool(false)),
			Comment:   ToPointer("test entry"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if e.IP != "10.0.0.3" {
		t.Errorf("bad IP: %q", e.IP)
	}

	if *e.Subnet != 8 {
		t.Errorf("Bad subnet: %v", e.Subnet)
	}

	if e.Negated {
		t.Errorf("Bad negated flag: %t", e.Negated)
	}

	if e.Comment != "test entry" {
		t.Errorf("Bad comment: %q", e.Comment)
	}

	// List
	var es []*ACLEntry
	record(t, fixtureBase+"list", func(c *Client) {
		es, err = c.ListACLEntries(&ListACLEntriesInput{
			ACLID:     testACL.ID,
			Direction: "descend",
			ServiceID: testService.ID,
			Sort:      "created",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(es) != 1 {
		t.Errorf("Bad entries: %v", es)
	}

	// List with paginator
	var es2 []*ACLEntry
	var paginator *ListPaginator[ACLEntry]
	record(t, fixtureBase+"list2", func(c *Client) {
		paginator = c.GetACLEntries(&GetACLEntriesInput{
			ACLID:     testACL.ID,
			Direction: "ascend",
			PerPage:   50,
			ServiceID: testService.ID,
			Sort:      "ip",
		})

		for paginator.HasNext() {
			data, err := paginator.GetNext()
			if err != nil {
				t.Errorf("Bad paginator (remaining: %d): %s", paginator.Remaining(), err)
				return
			}
			es2 = append(es2, data...)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(es2) != 1 {
		t.Errorf("Bad entries: %v", es2)
	}
	if paginator.HasNext() {
		t.Errorf("Bad paginator (remaining: %v)", paginator.Remaining())
	}

	// Get
	var ne *ACLEntry
	record(t, fixtureBase+"get", func(c *Client) {
		ne, err = c.GetACLEntry(&GetACLEntryInput{
			ServiceID: testService.ID,
			ACLID:     testACL.ID,
			ID:        e.ID,
		})
	})

	if e.IP != ne.IP {
		t.Errorf("bad IP: %v", ne.IP)
	}

	if *e.Subnet != *ne.Subnet {
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
			ServiceID: testService.ID,
			ACLID:     testACL.ID,
			ID:        e.ID,
			IP:        ToPointer("10.0.0.4"),
			Negated:   ToPointer(Compatibool(true)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if ue.IP != "10.0.0.4" {
		t.Errorf("bad IP: %q", ue.IP)
	}
	if *e.Subnet != *ue.Subnet {
		t.Errorf("bad subnet: %v", ne.Subnet)
	}

	if !ue.Negated {
		t.Errorf("bad Negated flag: %v", ne.Negated)
	}

	if e.Comment != ue.Comment {
		t.Errorf("bad comment: %v", ne.Comment)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteACLEntry(&DeleteACLEntryInput{
			ServiceID: testService.ID,
			ACLID:     testACL.ID,
			ID:        e.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListACLEntries_validation(t *testing.T) {
	var err error

	_, err = testClient.ListACLEntries(&ListACLEntriesInput{})
	if err != ErrMissingACLID {
		t.Errorf("bad ACL ID: %s", err)
	}

	_, err = testClient.ListACLEntries(&ListACLEntriesInput{
		ACLID: "123",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad Service ID: %s", err)
	}
}

func TestClient_CreateACLEntry_validation(t *testing.T) {
	var err error

	_, err = testClient.CreateACLEntry(&CreateACLEntryInput{})
	if err != ErrMissingACLID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateACLEntry(&CreateACLEntryInput{
		ACLID: "123",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetACLEntry_validation(t *testing.T) {
	var err error

	_, err = testClient.GetACLEntry(&GetACLEntryInput{})
	if err != ErrMissingACLID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetACLEntry(&GetACLEntryInput{
		ACLID: "123",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetACLEntry(&GetACLEntryInput{
		ACLID: "123",
		ID:    "456",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateACLEntry_validation(t *testing.T) {
	var err error

	_, err = testClient.UpdateACLEntry(&UpdateACLEntryInput{})
	if err != ErrMissingACLID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateACLEntry(&UpdateACLEntryInput{
		ACLID: "123",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateACLEntry(&UpdateACLEntryInput{
		ACLID: "123",
		ID:    "456",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteACLEntry_validation(t *testing.T) {
	var err error

	err = testClient.DeleteACLEntry(&DeleteACLEntryInput{})
	if err != ErrMissingACLID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteACLEntry(&DeleteACLEntryInput{
		ACLID: "123",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteACLEntry(&DeleteACLEntryInput{
		ACLID: "123",
		ID:    "456",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_BatchModifyACLEntries_validation(t *testing.T) {
	var err error

	err = testClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{})
	if err != ErrMissingACLID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{
		ACLID: "123",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	oversizedACLEntries := make([]*BatchACLEntry, BatchModifyMaximumOperations+1)
	err = testClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{
		ACLID:     "123",
		ServiceID: "456",
		Entries:   oversizedACLEntries,
	})
	if err != ErrMaxExceededEntries {
		t.Errorf("bad error: %s", err)
	}
}
