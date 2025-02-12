package computeacls

import (
	"net"
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
)

func TestClient_ComputeACL(t *testing.T) {
	t.Parallel()

	const aclName = "test-compute-acl"

	var acls *ComputeACLs
	var err error

	// List all compute ACLs.
	fastly.Record(t, "list_acls", func(c *fastly.Client) {
		acls, err = ListACLs(c)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the test compute acl we're going to create isn't among them.
	for _, acl := range acls.Data {
		if acl.Name == aclName {
			t.Errorf("found test compute ACL %q, aborting", aclName)
		}
	}

	// Create a compute ACL for testing.
	var acl *ComputeACL
	fastly.Record(t, "create_acl", func(c *fastly.Client) {
		acl, err = Create(c, &CreateInput{
			Name: fastly.ToPointer(aclName),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if acl.Name != aclName {
		t.Errorf("unexpected acl name: got %q, expected %q", acl.Name, aclName)
	}

	// Ensure we delete the test compute ACL at the end.
	defer func() {
		fastly.Record(t, "delete_acl", func(c *fastly.Client) {
			err = Delete(c, &DeleteInput{
				ComputeACLID: fastly.ToPointer(acl.ComputeACLID),
			})
		})
		if err != nil {
			t.Errorf("error during compute ACL cleanup: %v", err)
		}
	}()

	// Describe the test compute ACL.
	var da *ComputeACL
	fastly.Record(t, "describe_acl", func(c *fastly.Client) {
		da, err = Describe(c, &DescribeInput{
			ComputeACLID: fastly.ToPointer(acl.ComputeACLID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if da.Name != acl.Name {
		t.Errorf("unexpected acl name: got %q, expected %q", da.Name, acl.Name)
	}
	if da.ComputeACLID != acl.ComputeACLID {
		t.Errorf("unexpected acl ID: got %q, expected %q", da.ComputeACLID, acl.ComputeACLID)
	}

	entries := []*BatchComputeACLEntry{
		{
			Operation: fastly.ToPointer("create"),
			Prefix:    fastly.ToPointer("1.2.3.0/24"),
			Action:    fastly.ToPointer("BLOCK"),
		},
		{
			Operation: fastly.ToPointer("update"),
			Prefix:    fastly.ToPointer("1.2.3.4/32"),
			Action:    fastly.ToPointer("ALLOW"),
		},
		{
			Operation: fastly.ToPointer("create"),
			Prefix:    fastly.ToPointer("23.23.23.23/32"),
			Action:    fastly.ToPointer("ALLOW"),
		},
		{
			Operation: fastly.ToPointer("update"),
			Prefix:    fastly.ToPointer("192.168.0.0/16"),
			Action:    fastly.ToPointer("BLOCK"),
		},
	}

	// Add the entries to the test compute ACL.
	fastly.Record(t, "update_acl", func(c *fastly.Client) {
		err = Update(c, &UpdateInput{
			ComputeACLID: fastly.ToPointer(acl.ComputeACLID),
			Entries:      entries,
		})
	})
	if err != nil {
		t.Errorf("error updating compute acl: %v", err)
	}

	// List all entries of the test compute ACL and compare it to the input.
	var actualACLEntries *ComputeACLEntries
	fastly.Record(t, "list_entries", func(c *fastly.Client) {
		actualACLEntries, err = ListEntries(c, &ListEntriesInput{
			ComputeACLID: fastly.ToPointer(acl.ComputeACLID),
		})
	})
	if err != nil {
		t.Errorf("error fetching list of compute ACL entries: %v", err)
	}

	actualNumberOfACLEntries := len(actualACLEntries.Entries)
	expectedNumberOfComputeACLEntries := len(entries)
	if actualNumberOfACLEntries != expectedNumberOfComputeACLEntries {
		t.Errorf("incorrect number of compute ACL entries returned, expected: %d, got %d", expectedNumberOfComputeACLEntries, actualNumberOfACLEntries)
	}

	for i, entry := range actualACLEntries.Entries {
		actualPrefix := entry.Prefix
		expectedPrefix := entries[i].Prefix

		if actualPrefix != *expectedPrefix {
			t.Errorf("prefix does not match, expected %v, got %v", expectedPrefix, actualPrefix)
		}

		actualAction := entry.Action
		expectedAction := entries[i].Action

		if actualAction != *expectedAction {
			t.Errorf("action does not match, expected %v, got %v", expectedAction, actualAction)
		}
	}

	// Lookup each entry of the test compute ACL with pagination and compare it to the input.
	const Limit = 1
	page := 0
	cursor := ""
	input := &ListEntriesInput{}
	fastly.Record(t, "lookup_entries", func(c *fastly.Client) {
		for {
			input.ComputeACLID = fastly.ToPointer(acl.ComputeACLID)
			input.Limit = fastly.ToPointer(int(Limit))

			actualACLEntries, err = ListEntries(c, input)
			if err != nil {
				t.Errorf("error fetching list of compute ACL entries: %v", err)
			}

			numberOfEntries := len(actualACLEntries.Entries)
			if numberOfEntries != Limit {
				t.Errorf("wrong number of entries for page %d: got %d, requested %d", page, numberOfEntries, Limit)
			}

			ip, ipNet, err := net.ParseCIDR(actualACLEntries.Entries[0].Prefix)
			if err != nil {
				t.Fatal(err)
			}

			entry, err := Lookup(c, &LookupInput{
				ComputeACLID: fastly.ToPointer(acl.ComputeACLID),
				ComputeACLIP: fastly.ToPointer(ip.Mask(ipNet.Mask).String()),
			})
			if err != nil {
				t.Errorf("error during IP lookup: %v", err)
			}

			actualPrefix := entry.Prefix
			expectedPrefix := entries[page].Prefix
			if actualPrefix != *expectedPrefix {
				t.Errorf("prefix does not match, expected %v, got %v", expectedPrefix, actualPrefix)
			}

			actualAction := entry.Action
			expectedAction := entries[page].Action
			if actualAction != *expectedAction {
				t.Errorf("action does not match, expected %v, got %v", expectedAction, actualAction)
			}

			if cursor = actualACLEntries.Meta.NextCursor; cursor == "" {
				break
			}

			input.Cursor = fastly.ToPointer(cursor)

			page++
		}
	})

	// Lookup a non-existing IP in the test compute ACL
	fastly.Record(t, "lookup_non_existing_ip", func(c *fastly.Client) {
		entry, err := Lookup(c, &LookupInput{
			ComputeACLID: fastly.ToPointer(acl.ComputeACLID),
			ComputeACLIP: fastly.ToPointer("73.49.184.42"),
		})
		if entry != nil {
			t.Errorf("unexpected lookup result: got %+v, expected 'nil'", entry)
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestClient_Create_validation(t *testing.T) {
	_, err := Create(fastly.TestClient, &CreateInput{
		Name: nil,
	})
	if err != fastly.ErrMissingName {
		t.Errorf("expected ErrMissingName: got %s", err)
	}
}

func TestClient_Describe_validation(t *testing.T) {
	_, err := Describe(fastly.TestClient, &DescribeInput{
		ComputeACLID: nil,
	})
	if err != fastly.ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %s", err)
	}
}

func TestClient_Delete_validation(t *testing.T) {
	err := Delete(fastly.TestClient, &DeleteInput{
		ComputeACLID: nil,
	})
	if err != fastly.ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %s", err)
	}
}

func TestClient_Lookup_validation(t *testing.T) {
	var err error
	_, err = Lookup(fastly.TestClient, &LookupInput{
		ComputeACLID: nil,
		ComputeACLIP: fastly.ToPointer("1.2.3.4"),
	})
	if err != fastly.ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %s", err)
	}

	_, err = Lookup(fastly.TestClient, &LookupInput{
		ComputeACLID: fastly.ToPointer("foo"),
		ComputeACLIP: nil,
	})
	if err != fastly.ErrMissingComputeACLIP {
		t.Errorf("expected ErrMissingComputeACLIP: got %s", err)
	}
}

func TestClient_ListEntries_validation(t *testing.T) {
	_, err := ListEntries(fastly.TestClient, &ListEntriesInput{
		ComputeACLID: nil,
	})
	if err != fastly.ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %s", err)
	}
}

func TestClient_Update_validation(t *testing.T) {
	err := Update(fastly.TestClient, &UpdateInput{
		ComputeACLID: nil,
	})
	if err != fastly.ErrMissingComputeACLID {
		t.Errorf("expected ErrMissingComputeACLID: got %s", err)
	}
}
