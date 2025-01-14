package fastly

import (
	"encoding/json"
	"net"
	"strings"
	"testing"
)

func TestClient_ComputeACL(t *testing.T) {
	t.Parallel()

	const computeACLTestName = "test-compute-acl"

	var computeACLsListResp *ListComputeACLsResponse
	var err error

	// List all compute ACLs.
	Record(t, "compute_acls/list_acls", func(c *Client) {
		computeACLsListResp, err = c.ListComputeACLs()
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the test compute acl we're going to create isn't among them.
	for _, acl := range computeACLsListResp.Data {
		if acl.Name == computeACLTestName {
			t.Errorf("found test compute ACL %q, aborting", computeACLTestName)
		}
	}

	// Create a compute ACL for testing.
	var acl *ComputeACL
	Record(t, "compute_acls/create_acl", func(c *Client) {
		acl, err = c.CreateComputeACL(&CreateComputeACLInput{
			Name: computeACLTestName,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if acl.Name != computeACLTestName {
		t.Errorf("unexpected acl name: got %q, expected %q", acl.Name, computeACLTestName)
	}

	// Ensure we delete the test compute ACL at the end.
	defer func() {
		Record(t, "compute_acls/cleanup", func(c *Client) {
			err = c.DeleteComputeACL(&DeleteComputeACLInput{
				ComputeACLID: acl.ComputeACLID,
			})
		})
		if err != nil {
			t.Errorf("error during compute ACL cleanup: %v", err)
		}
	}()

	// Describe the test compute ACL.
	var describeComputeACLResponse *ComputeACL
	Record(t, "compute_acls/describe_acl", func(c *Client) {
		describeComputeACLResponse, err = c.DescribeComputeACL(&DescribeComputeACLInput{
			ComputeACLID: acl.ComputeACLID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if describeComputeACLResponse.Name != acl.Name {
		t.Errorf("unexpected acl name: got %q, expected %q", describeComputeACLResponse.Name, acl.Name)
	}

	if describeComputeACLResponse.ComputeACLID != acl.ComputeACLID {
		t.Errorf("unexpected acl ID: got %q, expected %q", describeComputeACLResponse.ComputeACLID, acl.ComputeACLID)
	}

	// Add a bunch of entries to the test compute ACL.
	batchCreateEntries := struct {
		Entries []struct {
			Operation string `json:"op"`
			Prefix    string `json:"prefix"`
			Action    string `json:"action"`
		} `json:"entries"`
	}{
		Entries: []struct {
			Operation string `json:"op"`
			Prefix    string `json:"prefix"`
			Action    string `json:"action"`
		}{
			{
				Operation: "create",
				Prefix:    "1.2.3.0/24",
				Action:    "BLOCK",
			},
			{
				Operation: "update",
				Prefix:    "1.2.3.4/32",
				Action:    "ALLOW",
			},
			{
				Operation: "create",
				Prefix:    "23.23.23.23/32",
				Action:    "ALLOW",
			},
			{
				Operation: "update",
				Prefix:    "192.168.0.0/16",
				Action:    "BLOCK",
			},
		},
	}

	body, err := json.Marshal(batchCreateEntries)
	if err != nil {
		t.Fatal(err)
	}

	// Add the entries to the test compute ACL.
	Record(t, "compute_acls/batch_create_entries", func(c *Client) {
		err = c.BatchModifyComputeACLEntries(&BatchModifyComputeACLEntriesInput{
			ComputeACLID: acl.ComputeACLID,
			Body:         strings.NewReader(string(body)),
		})
	})
	if err != nil {
		t.Errorf("error creating multiple compute ACL entries: %v", err)
	}

	// List all entries of the test compute ACL and compare it to the input.
	var actualComputeACLEntries *ListComputeACLEntriesResponse
	Record(t, "compute_acls/list_entries_after_batch_create", func(c *Client) {
		actualComputeACLEntries, err = c.ListComputeACLEntries(&ListComputeACLEntriesInput{
			ComputeACLID: acl.ComputeACLID,
		})
	})
	if err != nil {
		t.Errorf("error fetching list of compute ACL entries: %v", err)
	}

	actualNumberOfComputeACLEntries := len(actualComputeACLEntries.Entries)
	expectedNumberOfComputeACLEntries := len(batchCreateEntries.Entries)
	if actualNumberOfComputeACLEntries != expectedNumberOfComputeACLEntries {
		t.Errorf("incorrect number of compute ACL entries returned, expected: %d, got %d", expectedNumberOfComputeACLEntries, actualNumberOfComputeACLEntries)
	}

	for i, entry := range actualComputeACLEntries.Entries {
		actualPrefix := entry.Prefix
		expectedPrefix := batchCreateEntries.Entries[i].Prefix

		if actualPrefix != expectedPrefix {
			t.Errorf("prefix does not match, expected %v, got %v", expectedPrefix, actualPrefix)
		}

		actualAction := entry.Action
		expectedAction := batchCreateEntries.Entries[i].Action

		if actualAction != expectedAction {
			t.Errorf("action does not match, expected %v, got %v", expectedAction, actualAction)
		}
	}

	// Lookup each entry of the test compute ACL with pagination and compare it to the input.
	const Limit = 1
	Record(t, "compute_acls/lookup_entries_after_batch_create", func(c *Client) {
		p := c.NewListComputeACLEntriesPaginator(&ListComputeACLEntriesInput{
			ComputeACLID: acl.ComputeACLID,
			Limit:        Limit,
		})

		var page int
		for p.Next() {
			numberOfEntries := len(p.Entries())
			if numberOfEntries != Limit {
				t.Errorf("wrong number of entries for page %d: got %d, requested %d", page, numberOfEntries, Limit)
			}

			ip, ipNet, err := net.ParseCIDR(p.Entries()[0].Prefix)
			if err != nil {
				t.Errorf("error parsing IP: %v", err)
			}

			entry, err := c.ComputeACLLookup(&ComputeACLLookupInput{
				ComputeACLID: acl.ComputeACLID,
				ComputeACLIP: ip.Mask(ipNet.Mask).String(),
			})
			if err != nil {
				t.Errorf("error during IP lookup: %v", err)
			}

			actualPrefix := entry.Prefix
			expectedPrefix := batchCreateEntries.Entries[page].Prefix

			if actualPrefix != expectedPrefix {
				t.Errorf("prefix does not match, expected %v, got %v", expectedPrefix, actualPrefix)
			}

			actualAction := entry.Action
			expectedAction := batchCreateEntries.Entries[page].Action

			if actualAction != expectedAction {
				t.Errorf("action does not match, expected %v, got %v", expectedAction, actualAction)
			}

			page++
		}
		if err := p.Err(); err != nil {
			t.Errorf("error during IP lookup pagination: %v", err)
		}
	})
}

func TestClient_CreateComputeACL_validation(t *testing.T) {
	_, err := TestClient.CreateComputeACL(&CreateComputeACLInput{
		Name: "",
	})
	if err != ErrMissingName {
		t.Errorf("expected ErrMissingName: got %q", err)
	}
}

func TestClient_DescribeComputeACL_validation(t *testing.T) {
	_, err := TestClient.DescribeComputeACL(&DescribeComputeACLInput{
		ComputeACLID: "",
	})
	if err != ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %q", err)
	}
}

func TestClient_DeleteComputeACL_validation(t *testing.T) {
	err := TestClient.DeleteComputeACL(&DeleteComputeACLInput{
		ComputeACLID: "",
	})
	if err != ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %q", err)
	}
}

func TestClient_ComputeACLLookup_validation(t *testing.T) {
	var err error
	_, err = TestClient.ComputeACLLookup(&ComputeACLLookupInput{
		ComputeACLID: "",
	})
	if err != ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %q", err)
	}

	_, err = TestClient.ComputeACLLookup(&ComputeACLLookupInput{
		ComputeACLID: "foo",
		ComputeACLIP: "",
	})
	if err != ErrMissingComputeACLIP {
		t.Errorf("expected ErrMissingComputeACLIP: got %q", err)
	}
}

func TestClient_ListComputeACLEntries(t *testing.T) {
	_, err := TestClient.ListComputeACLEntries(&ListComputeACLEntriesInput{
		ComputeACLID: "",
	})
	if err != ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %q", err)
	}
}

func TestClient_BatchModifyComputeACLEntries(t *testing.T) {
	err := TestClient.BatchModifyComputeACLEntries(&BatchModifyComputeACLEntriesInput{
		ComputeACLID: "",
	})
	if err != ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %q", err)
	}
}
