package fastly

import (
	"errors"
	"testing"
)

func TestClient_ACLEntries(t *testing.T) {
	fixtureBase := "acl_entries/"
	nameSuffix := "ACLEntries"

	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"version", *testService.ServiceID)

	testACL := createTestACL(t, fixtureBase+"acl", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestACL(t, testACL, fixtureBase+"delete_acl")

	// Create
	var err error
	var e *ACLEntry
	Record(t, fixtureBase+"create", func(c *Client) {
		e, err = c.CreateACLEntry(&CreateACLEntryInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
			IP:        ToPointer("10.0.0.3"),
			Subnet:    ToPointer(8),
			Negated:   ToPointer(Compatibool(false)),
			Comment:   ToPointer("test entry"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *e.IP != "10.0.0.3" {
		t.Errorf("bad IP: %q", *e.IP)
	}

	if *e.Subnet != 8 {
		t.Errorf("Bad subnet: %v", *e.Subnet)
	}

	if *e.Negated {
		t.Errorf("Bad negated flag: %t", *e.Negated)
	}

	if *e.Comment != "test entry" {
		t.Errorf("Bad comment: %q", *e.Comment)
	}

	// List
	var es []*ACLEntry
	Record(t, fixtureBase+"list", func(c *Client) {
		es, err = c.ListACLEntries(&ListACLEntriesInput{
			ACLID:     *testACL.ACLID,
			Direction: ToPointer("descend"),
			ServiceID: *testService.ServiceID,
			Sort:      ToPointer("created"),
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
	Record(t, fixtureBase+"list2", func(c *Client) {
		paginator = c.GetACLEntries(&GetACLEntriesInput{
			ACLID:     *testACL.ACLID,
			Direction: ToPointer("ascend"),
			PerPage:   ToPointer(50),
			ServiceID: *testService.ServiceID,
			Sort:      ToPointer("ip"),
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
	Record(t, fixtureBase+"get", func(c *Client) {
		ne, err = c.GetACLEntry(&GetACLEntryInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
			EntryID:   *e.EntryID,
		})
	})

	if *e.IP != *ne.IP {
		t.Errorf("bad IP: %v", *ne.IP)
	}

	if *e.Subnet != *ne.Subnet {
		t.Errorf("bad subnet: %v", *ne.Subnet)
	}

	if *e.Negated != *ne.Negated {
		t.Errorf("bad Negated flag: %v", *ne.Negated)
	}

	if *e.Comment != *ne.Comment {
		t.Errorf("bad comment: %v", *ne.Comment)
	}

	// Update
	var ue *ACLEntry
	Record(t, fixtureBase+"update", func(c *Client) {
		ue, err = c.UpdateACLEntry(&UpdateACLEntryInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
			EntryID:   *e.EntryID,
			IP:        ToPointer("10.0.0.4"),
			Negated:   ToPointer(Compatibool(true)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *ue.IP != "10.0.0.4" {
		t.Errorf("bad IP: %q", *ue.IP)
	}
	if *e.Subnet != *ue.Subnet {
		t.Errorf("bad subnet: %v", *ue.Subnet)
	}

	if !*ue.Negated {
		t.Errorf("bad Negated flag: %v", *ue.Negated)
	}

	if *e.Comment != *ue.Comment {
		t.Errorf("bad comment: %v", *ue.Comment)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteACLEntry(&DeleteACLEntryInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
			EntryID:   *e.EntryID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListACLEntries_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListACLEntries(&ListACLEntriesInput{})
	if !errors.Is(err, ErrMissingACLID) {
		t.Errorf("bad ACL ID: %s", err)
	}

	_, err = TestClient.ListACLEntries(&ListACLEntriesInput{
		ACLID: "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad Service ID: %s", err)
	}
}

func TestClient_CreateACLEntry_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateACLEntry(&CreateACLEntryInput{})
	if !errors.Is(err, ErrMissingACLID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateACLEntry(&CreateACLEntryInput{
		ACLID: "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetACLEntry_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetACLEntry(&GetACLEntryInput{})
	if !errors.Is(err, ErrMissingACLID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetACLEntry(&GetACLEntryInput{
		ACLID: "123",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetACLEntry(&GetACLEntryInput{
		ACLID:   "123",
		EntryID: "456",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateACLEntry_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateACLEntry(&UpdateACLEntryInput{})
	if !errors.Is(err, ErrMissingACLID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateACLEntry(&UpdateACLEntryInput{
		ACLID: "123",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateACLEntry(&UpdateACLEntryInput{
		ACLID:   "123",
		EntryID: "456",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteACLEntry_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteACLEntry(&DeleteACLEntryInput{})
	if !errors.Is(err, ErrMissingACLID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteACLEntry(&DeleteACLEntryInput{
		ACLID: "123",
	})
	if !errors.Is(err, ErrMissingEntryID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteACLEntry(&DeleteACLEntryInput{
		ACLID:   "123",
		EntryID: "456",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_BatchModifyACLEntries_validation(t *testing.T) {
	var err error

	err = TestClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{})
	if !errors.Is(err, ErrMissingACLID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{
		ACLID: "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	oversizedACLEntries := make([]*BatchACLEntry, BatchModifyMaximumOperations+1)
	err = TestClient.BatchModifyACLEntries(&BatchModifyACLEntriesInput{
		ACLID:     "123",
		ServiceID: "456",
		Entries:   oversizedACLEntries,
	})
	if !errors.Is(err, ErrMaxExceededEntries) {
		t.Errorf("bad error: %s", err)
	}
}
