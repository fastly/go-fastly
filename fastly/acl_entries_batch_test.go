package fastly

import (
	"sort"
	"testing"
)

func TestClient_BatchModifyACLEntries_Create(t *testing.T) {
	fixtureBase := "acl_entries_batch/create/"
	nameSuffix := "BatchModifyAclEntries_Create"

	// Given: a test service with an ACL and a batch of create operations,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"create_version", *testService.ServiceID)

	testACL := createTestACL(t, fixtureBase+"create_acl", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestACL(t, testACL, fixtureBase+"delete_acl")

	batchCreateOperations := &BatchModifyACLEntriesInput{
		ServiceID: *testService.ServiceID,
		ACLID:     *testACL.ACLID,
		Entries: []*BatchACLEntry{
			{
				Operation: ToPointer(CreateBatchOperation),
				IP:        ToPointer("127.0.0.1"),
				Subnet:    ToPointer(24),
				Negated:   ToPointer(Compatibool(false)),
				Comment:   ToPointer("ACL Entry 1"),
			},
			{
				Operation: ToPointer(CreateBatchOperation),
				IP:        ToPointer("192.168.0.1"),
				Subnet:    ToPointer(24),
				Negated:   ToPointer(Compatibool(true)),
				Comment:   ToPointer("ACL Entry 2"),
			},
		},
	}

	// When: I execute the batch create operations against the Fastly API,
	var err error
	Record(t, fixtureBase+"create_acl_entries", func(c *Client) {
		err = c.BatchModifyACLEntries(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the created ACL entries.
	var actualACLEntries []*ACLEntry
	Record(t, fixtureBase+"list_after_create", func(c *Client) {
		actualACLEntries, err = c.ListACLEntries(&ListACLEntriesInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(actualACLEntries, func(i, j int) bool {
		return *actualACLEntries[i].IP < *actualACLEntries[j].IP
	})

	actualNumberOfACLEntries := len(actualACLEntries)
	expectedNumberOfACLEntries := len(batchCreateOperations.Entries)
	if actualNumberOfACLEntries != expectedNumberOfACLEntries {
		t.Errorf("Incorrect number of ACL entries returned, expected: %d, got %d", expectedNumberOfACLEntries, actualNumberOfACLEntries)
	}

	for i, entry := range actualACLEntries {
		actualIP := entry.IP
		expectedIP := batchCreateOperations.Entries[i].IP

		if *actualIP != *expectedIP {
			t.Errorf("IP did not match, expected %v, got %v", *expectedIP, *actualIP)
		}

		actualSubnet := entry.Subnet
		expectedSubnet := batchCreateOperations.Entries[i].Subnet

		if *actualSubnet != *expectedSubnet {
			t.Errorf("Subnet did not match, expected %v, got %v", expectedSubnet, actualSubnet)
		}

		actualNegated := entry.Negated
		expectedNegated := bool(*batchCreateOperations.Entries[i].Negated)

		if *actualNegated != expectedNegated {
			t.Errorf("Negated did not match, expected %v, got %v", expectedNegated, *actualNegated)
		}

		actualComment := entry.Comment
		expectedComment := batchCreateOperations.Entries[i].Comment

		if *actualComment != *expectedComment {
			t.Errorf("Comment did not match, expected %v, got %v", *expectedComment, *actualComment)
		}
	}
}

func TestClient_BatchModifyACLEntries_Delete(t *testing.T) {
	fixtureBase := "acl_entries_batch/delete/"
	nameSuffix := "BatchModifyAclEntries_Delete"

	// Given: a test service with an ACL and a batch of create operations,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"create_version", *testService.ServiceID)

	testACL := createTestACL(t, fixtureBase+"create_acl", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestACL(t, testACL, fixtureBase+"delete_acl")

	batchCreateOperations := &BatchModifyACLEntriesInput{
		ServiceID: *testService.ServiceID,
		ACLID:     *testACL.ACLID,
		Entries: []*BatchACLEntry{
			{
				Operation: ToPointer(CreateBatchOperation),
				IP:        ToPointer("127.0.0.1"),
				Subnet:    ToPointer(24),
				Negated:   ToPointer(Compatibool(false)),
				Comment:   ToPointer("ACL Entry 1"),
			},
			{
				Operation: ToPointer(CreateBatchOperation),
				IP:        ToPointer("192.168.0.1"),
				Subnet:    ToPointer(24),
				Negated:   ToPointer(Compatibool(true)),
				Comment:   ToPointer("ACL Entry 2"),
			},
		},
	}

	var err error
	Record(t, fixtureBase+"create_acl_entries", func(c *Client) {
		err = c.BatchModifyACLEntries(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	var createdACLEntries []*ACLEntry
	Record(t, fixtureBase+"list_before_delete", func(client *Client) {
		createdACLEntries, err = client.ListACLEntries(&ListACLEntriesInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(createdACLEntries, func(i, j int) bool {
		return *createdACLEntries[i].IP < *createdACLEntries[j].IP
	})

	// When: I execute the batch delete operations against the Fastly API,
	batchDeleteOperations := &BatchModifyACLEntriesInput{
		ServiceID: *testService.ServiceID,
		ACLID:     *testACL.ACLID,
		Entries: []*BatchACLEntry{
			{
				Operation: ToPointer(DeleteBatchOperation),
				EntryID:   createdACLEntries[0].EntryID,
			},
		},
	}

	Record(t, fixtureBase+"delete_acl_entries", func(c *Client) {
		err = c.BatchModifyACLEntries(batchDeleteOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list a single ACL entry.
	var actualACLEntries []*ACLEntry
	Record(t, fixtureBase+"list_after_delete", func(client *Client) {
		actualACLEntries, err = client.ListACLEntries(&ListACLEntriesInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(actualACLEntries, func(i, j int) bool {
		return *actualACLEntries[i].IP < *actualACLEntries[j].IP
	})

	actualNumberOfACLEntries := len(actualACLEntries)
	expectedNumberOfACLEntries := len(batchDeleteOperations.Entries)
	if actualNumberOfACLEntries != expectedNumberOfACLEntries {
		t.Errorf("Incorrect number of ACL entries returned, expected: %d, got %d", expectedNumberOfACLEntries, actualNumberOfACLEntries)
	}
}

func TestClient_BatchModifyACLEntries_Update(t *testing.T) {
	fixtureBase := "acl_entries_batch/update/"
	nameSuffix := "BatchModifyAclEntries_Update"

	// Given: a test service with an ACL and ACL entries,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"create_version", *testService.ServiceID)

	testACL := createTestACL(t, fixtureBase+"create_acl", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestACL(t, testACL, fixtureBase+"delete_acl")

	batchCreateOperations := &BatchModifyACLEntriesInput{
		ServiceID: *testService.ServiceID,
		ACLID:     *testACL.ACLID,
		Entries: []*BatchACLEntry{
			{
				Operation: ToPointer(CreateBatchOperation),
				IP:        ToPointer("127.0.0.1"),
				Subnet:    ToPointer(24),
				Negated:   ToPointer(Compatibool(false)),
				Comment:   ToPointer("ACL Entry 1"),
			},
			{
				Operation: ToPointer(CreateBatchOperation),
				IP:        ToPointer("192.168.0.1"),
				Subnet:    ToPointer(24),
				Negated:   ToPointer(Compatibool(true)),
				Comment:   ToPointer("ACL Entry 2"),
			},
		},
	}

	var err error
	Record(t, fixtureBase+"create_acl_entries", func(c *Client) {
		err = c.BatchModifyACLEntries(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	var createdACLEntries []*ACLEntry
	Record(t, fixtureBase+"list_before_update", func(client *Client) {
		createdACLEntries, err = client.ListACLEntries(&ListACLEntriesInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(createdACLEntries, func(i, j int) bool {
		return *createdACLEntries[i].IP < *createdACLEntries[j].IP
	})

	// When: I execute the batch update operations against the Fastly API,
	batchUpdateOperations := &BatchModifyACLEntriesInput{
		ServiceID: *testService.ServiceID,
		ACLID:     *testACL.ACLID,
		Entries: []*BatchACLEntry{
			{
				Operation: ToPointer(UpdateBatchOperation),
				EntryID:   createdACLEntries[0].EntryID,
				IP:        ToPointer("127.0.0.2"),
				Subnet:    ToPointer(16),
				Negated:   ToPointer(Compatibool(true)),
				Comment:   ToPointer("Updated ACL Entry 1"),
			},
		},
	}

	Record(t, fixtureBase+"update_acl_entries", func(c *Client) {
		err = c.BatchModifyACLEntries(batchUpdateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the ACL entries with modifications applied to a single item.
	var actualACLEntries []*ACLEntry
	Record(t, fixtureBase+"list_after_update", func(client *Client) {
		actualACLEntries, err = client.ListACLEntries(&ListACLEntriesInput{
			ServiceID: *testService.ServiceID,
			ACLID:     *testACL.ACLID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(actualACLEntries, func(i, j int) bool {
		return *actualACLEntries[i].IP < *actualACLEntries[j].IP
	})

	actualNumberOfACLEntries := len(actualACLEntries)
	expectedNumberOfACLEntries := len(batchCreateOperations.Entries)
	if actualNumberOfACLEntries != expectedNumberOfACLEntries {
		t.Errorf("Incorrect number of ACL entries returned, expected: %d, got %d", expectedNumberOfACLEntries, actualNumberOfACLEntries)
	}

	actualID := actualACLEntries[0].EntryID
	expectedID := batchUpdateOperations.Entries[0].EntryID

	if *actualID != *expectedID {
		t.Errorf("First ID did not match, expected %v, got %v", *expectedID, *actualID)
	}

	actualIP := actualACLEntries[0].IP
	expectedIP := batchUpdateOperations.Entries[0].IP

	if *actualIP != *expectedIP {
		t.Errorf("First IP did not match, expected %v, got %v", *expectedIP, *actualIP)
	}

	actualSubnet := actualACLEntries[0].Subnet
	expectedSubnet := batchUpdateOperations.Entries[0].Subnet

	if *actualSubnet != *expectedSubnet {
		t.Errorf("First Subnet did not match, expected %v, got %v", expectedSubnet, actualSubnet)
	}

	actualNegated := actualACLEntries[0].Negated
	expectedNegated := bool(*batchUpdateOperations.Entries[0].Negated)

	if *actualNegated != expectedNegated {
		t.Errorf("First Subnet did not match, expected %v, got %v", expectedNegated, *actualNegated)
	}

	actualComment := actualACLEntries[0].Comment
	expectedComment := batchUpdateOperations.Entries[0].Comment

	if *actualComment != *expectedComment {
		t.Errorf("First Comment did not match, expected %v, got %v", expectedComment, *actualComment)
	}
}
