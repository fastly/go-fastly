package fastly

import (
	"sort"
	"testing"
)

func TestClient_BatchModifyAclEntries_Create(t *testing.T) {

	fixtureBase := "acl_entries_batch/create/"
	nameSuffix := "BatchModifyAclEntries_Create"

	// Given: a test service with an ACL and a batch of create operations,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", testService.ID)

	testVersion := createTestVersion(t, fixtureBase+"create_version", testService.ID)

	testACL := createTestACL(t, fixtureBase+"create_acl", testService.ID, testVersion.Number, nameSuffix)
	defer deleteTestACL(t, testACL, fixtureBase+"delete_acl")

	batchCreateOperations := &BatchModifyACLEntriesInput{
		Service: testService.ID,
		ACL:     testACL.ID,
		Entries: []*BatchACLEntry{
			{
				Operation: CreateBatchOperation,
				IP:        "127.0.0.1",
				Subnet:    "24",
				Negated:   false,
				Comment:   "ACL Entry 1",
			},
			{
				Operation: CreateBatchOperation,
				IP:        "192.168.0.1",
				Subnet:    "24",
				Negated:   false,
				Comment:   "ACL Entry 2",
			},
		},
	}

	// When: I execute the batch create operations against the Fastly API,
	var err error
	record(t, fixtureBase+"create_acl_entries", func(c *Client) {

		err = c.BatchModifyACLEntries(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the created ACL entries.
	var actualACLEntries []*ACLEntry
	record(t, fixtureBase+"list_after_create", func(c *Client) {
		actualACLEntries, err = c.ListACLEntries(&ListACLEntriesInput{
			Service: testService.ID,
			ACL:     testACL.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(actualACLEntries, func(i, j int) bool {
		return actualACLEntries[i].IP < actualACLEntries[j].IP
	})

	actualNumberOfACLEntries := len(actualACLEntries)
	expectedNumberOfACLEntries := len(batchCreateOperations.Entries)
	if actualNumberOfACLEntries != expectedNumberOfACLEntries {
		t.Errorf("Incorrect number of ACL entries returned, expected: %d, got %d", expectedNumberOfACLEntries, actualNumberOfACLEntries)
	}

	for i, entry := range actualACLEntries {

		actualIp := entry.IP
		expectedIp := batchCreateOperations.Entries[i].IP

		if actualIp != expectedIp {
			t.Errorf("IP did not match, expected %s, got %s", expectedIp, actualIp)
		}

		actualSubnet := entry.Subnet
		expectedSubnet := batchCreateOperations.Entries[i].Subnet

		if actualSubnet != expectedSubnet {
			t.Errorf("Subnet did not match, expected %s, got %s", expectedSubnet, actualSubnet)
		}

		actualNegated := entry.Negated
		expectedNegated := batchCreateOperations.Entries[i].Negated

		if actualNegated != expectedNegated {
			t.Errorf("Negated did not match, expected %t, got %t", expectedNegated, actualNegated)
		}

		actualComment := entry.Comment
		expectedComment := batchCreateOperations.Entries[i].Comment

		if actualComment != expectedComment {
			t.Errorf("Comment did not match, expected %s, got %s", expectedComment, actualComment)
		}
	}

}

func TestClient_BatchModifyAclEntries_Delete(t *testing.T) {

	fixtureBase := "acl_entries_batch/delete/"
	nameSuffix := "BatchModifyAclEntries_Delete"

	// Given: a test service with an ACL and a batch of create operations,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", testService.ID)

	testVersion := createTestVersion(t, fixtureBase+"create_version", testService.ID)

	testACL := createTestACL(t, fixtureBase+"create_acl", testService.ID, testVersion.Number, nameSuffix)
	defer deleteTestACL(t, testACL, fixtureBase+"delete_acl")

	batchCreateOperations := &BatchModifyACLEntriesInput{
		Service: testService.ID,
		ACL:     testACL.ID,
		Entries: []*BatchACLEntry{
			{
				Operation: CreateBatchOperation,
				IP:        "127.0.0.1",
				Subnet:    "24",
				Negated:   false,
				Comment:   "ACL Entry 1",
			},
			{
				Operation: CreateBatchOperation,
				IP:        "192.168.0.1",
				Subnet:    "24",
				Negated:   false,
				Comment:   "ACL Entry 2",
			},
		},
	}

	var err error
	record(t, fixtureBase+"create_acl_entries", func(c *Client) {

		err = c.BatchModifyACLEntries(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	var createdACLEntries []*ACLEntry
	record(t, fixtureBase+"list_before_delete", func(client *Client) {
		createdACLEntries, err = client.ListACLEntries(&ListACLEntriesInput{
			Service: testService.ID,
			ACL:     testACL.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(createdACLEntries, func(i, j int) bool {
		return createdACLEntries[i].IP < createdACLEntries[j].IP
	})

	// When: I execute the batch delete operations against the Fastly API,
	batchDeleteOperations := &BatchModifyACLEntriesInput{
		Service: testService.ID,
		ACL:     testACL.ID,
		Entries: []*BatchACLEntry{
			{
				Operation: DeleteBatchOperation,
				ID: createdACLEntries[0].ID,
			},
		},
	}

	record(t, fixtureBase+"delete_acl_entries", func(c *Client) {

		err = c.BatchModifyACLEntries(batchDeleteOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list a single ACL entry.
	var actualACLEntries []*ACLEntry
	record(t, fixtureBase+"list_after_delete", func(client *Client) {
		actualACLEntries, err = client.ListACLEntries(&ListACLEntriesInput{
			Service:    testService.ID,
			ACL: testACL.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(actualACLEntries, func(i, j int) bool {
		return actualACLEntries[i].IP < actualACLEntries[j].IP
	})

	actualNumberOfACLEntries := len(actualACLEntries)
	expectedNumberOfACLEntries := len(batchDeleteOperations.Entries)
	if actualNumberOfACLEntries != expectedNumberOfACLEntries {
		t.Errorf("Incorrect number of ACL entries returned, expected: %d, got %d", expectedNumberOfACLEntries, actualNumberOfACLEntries)
	}
}

func TestClient_BatchModifyAclEntries_Update(t *testing.T) {

	fixtureBase := "acl_entries_batch/update/"
	nameSuffix := "BatchModifyAclEntries_Update"

	// Given: a test service with an ACL and ACL entries,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", testService.ID)

	testVersion := createTestVersion(t, fixtureBase+"create_version", testService.ID)

	testACL := createTestACL(t, fixtureBase+"create_acl", testService.ID, testVersion.Number, nameSuffix)
	defer deleteTestACL(t, testACL, fixtureBase+"delete_acl")

	batchCreateOperations := &BatchModifyACLEntriesInput{
		Service: testService.ID,
		ACL:     testACL.ID,
		Entries: []*BatchACLEntry{
			{
				Operation: CreateBatchOperation,
				IP:        "127.0.0.1",
				Subnet:    "24",
				Negated:   false,
				Comment:   "ACL Entry 1",
			},
			{
				Operation: CreateBatchOperation,
				IP:        "192.168.0.1",
				Subnet:    "24",
				Negated:   false,
				Comment:   "ACL Entry 2",
			},
		},
	}

	var err error
	record(t, fixtureBase+"create_acl_entries", func(c *Client) {

		err = c.BatchModifyACLEntries(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	var createdACLEntries []*ACLEntry
	record(t, fixtureBase+"list_before_update", func(client *Client) {
		createdACLEntries, err = client.ListACLEntries(&ListACLEntriesInput{
			Service: testService.ID,
			ACL:     testACL.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(createdACLEntries, func(i, j int) bool {
		return createdACLEntries[i].IP < createdACLEntries[j].IP
	})

	// When: I execute the batch update operations against the Fastly API,
	batchUpdateOperations := &BatchModifyACLEntriesInput{
		Service: testService.ID,
		ACL:     testACL.ID,
		Entries: []*BatchACLEntry{
			{
				Operation: UpdateBatchOperation,
				ID: createdACLEntries[0].ID,
				IP: "127.0.0.2",
				Subnet: "16",
				Negated: true,
				Comment: "Updated ACL Entry 1",
			},
		},
	}

	record(t, fixtureBase+"update_acl_entries", func(c *Client) {

		err = c.BatchModifyACLEntries(batchUpdateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the ACL entries with modifications applied to a single item.
	var actualACLEntries []*ACLEntry
	record(t, fixtureBase+"list_after_update", func(client *Client) {
		actualACLEntries, err = client.ListACLEntries(&ListACLEntriesInput{
			Service: testService.ID,
			ACL:     testACL.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(actualACLEntries, func(i, j int) bool {
		return actualACLEntries[i].IP < actualACLEntries[j].IP
	})

	actualNumberOfAclEntries := len(actualACLEntries)
	expectedNumberOfAclEntries := len(batchCreateOperations.Entries)
	if actualNumberOfAclEntries != expectedNumberOfAclEntries {
		t.Errorf("Incorrect number of ACL entries returned, expected: %d, got %d", expectedNumberOfAclEntries, actualNumberOfAclEntries)
	}

	actualID := actualACLEntries[0].ID
	expectedID := batchUpdateOperations.Entries[0].ID

	if actualID != expectedID {
		t.Errorf("First ID did not match, expected %s, got %s", expectedID, actualID)
	}

	actualIP := actualACLEntries[0].IP
	expectedIP := batchUpdateOperations.Entries[0].IP

	if actualIP != expectedIP {
		t.Errorf("First IP did not match, expected %s, got %s", expectedIP, actualIP)
	}

	actualSubnet := actualACLEntries[0].Subnet
	expectedSubnet := batchUpdateOperations.Entries[0].Subnet

	if actualSubnet != expectedSubnet {
		t.Errorf("First Subnet did not match, expected %s, got %s", expectedSubnet, actualSubnet)
	}

	actualNegated := actualACLEntries[0].Negated
	expectedNegated := batchUpdateOperations.Entries[0].Negated

	if actualNegated != expectedNegated {
		t.Errorf("First Subnet did not match, expected %t, got %t", expectedNegated, actualNegated)
	}

	actualComment := actualACLEntries[0].Comment
	expectedComment := batchUpdateOperations.Entries[0].Comment

	if actualComment != expectedComment {
		t.Errorf("First Comment did not match, expected %s, got %s", expectedComment, actualComment)
	}

}