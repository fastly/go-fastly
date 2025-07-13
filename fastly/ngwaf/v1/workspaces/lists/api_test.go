package lists

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v11/fastly"
)

const (
	listDescription string = "This is a list"
	listEntry       string = "entry1"
	listName        string = "List Name"
	listReferenceID string = "site.list-name"
	listType        string = "string"
)

func TestClient_List(t *testing.T) {
	var err error
	listEntries := []string{listEntry}

	// Create a test list.
	var list *List
	fastly.Record(t, "create_list", func(c *fastly.Client) {
		list, err = Create(context.TODO(), c, &CreateInput{
			Description: fastly.ToPointer(listDescription),
			Entries:     fastly.ToPointer(listEntries),
			Name:        fastly.ToPointer(listName),
			Type:        fastly.ToPointer(listType),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if list.Description != listDescription {
		t.Errorf("unexpected list description: got %q, expected %q", list.Description, listDescription)
	}
	if len(list.Entries) != len(listEntries) {
		t.Errorf("unexpected list entry length: got %q, expected %q", len(list.Entries), len(listEntries))
	}
	if list.Entries[0] != listEntries[0] {
		t.Errorf("unexpected list entry: got %q, expected %q", list.Entries[0], listEntries[0])
	}
	if list.Name != listName {
		t.Errorf("unexpected list name: got %q, expected %q", list.Name, listName)
	}
	if list.ReferenceID != listReferenceID {
		t.Errorf("unexpected list reference ID: got %q, expected %q", list.ReferenceID, listReferenceID)
	}
	if list.Type != listType {
		t.Errorf("unexpected list type: got %q, expected %q", list.Type, listType)
	}
	// Scope should always be workspace
	if list.Scope.Type != "workspace" {
		t.Errorf("unexpected list scope: got %q, expected %q", list.Scope.Type, "workspace")
	}

	// Ensure we delete the test list at the end.
	defer func() {
		fastly.Record(t, "delete_list", func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				ListID:      fastly.ToPointer(list.ListID),
				WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			})
		})
		if err != nil {
			t.Errorf("error during list cleanup: %v", err)
		}
	}()

	// Get the test list.
	var getList *List
	fastly.Record(t, "get_list", func(c *fastly.Client) {
		getList, err = Get(context.TODO(), c, &GetInput{
			ListID:      fastly.ToPointer(list.ListID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if getList.Description != listDescription {
		t.Errorf("unexpected list description: got %q, expected %q", getList.Description, listDescription)
	}
	if len(getList.Entries) != len(listEntries) {
		t.Errorf("unexpected list entry length: got %q, expected %q", len(getList.Entries), len(listEntries))
	}
	if getList.Entries[0] != listEntries[0] {
		t.Errorf("unexpected list entry: got %q, expected %q", getList.Entries[0], listEntries[0])
	}
	if getList.Name != listName {
		t.Errorf("unexpected list name: got %q, expected %q", getList.Name, listName)
	}
	if getList.ReferenceID != listReferenceID {
		t.Errorf("unexpected list reference ID: got %q, expected %q", getList.ReferenceID, listReferenceID)
	}
	if getList.Type != listType {
		t.Errorf("unexpected list type: got %q, expected %q", getList.Type, listType)
	}
	// Scope should always be workspace
	if getList.Scope.Type != "workspace" {
		t.Errorf("unexpected list scope: got %q, expected %q", getList.Scope.Type, "workspace")
	}

	// Update the test list.
	updateListDescription := "This is an updated list"
	updateListEntries := []string{
		"127.0.0.1",
	}

	var updateList *List
	fastly.Record(t, "update_list", func(c *fastly.Client) {
		updateList, err = Update(context.TODO(), c, &UpdateInput{
			Description: fastly.ToPointer(updateListDescription),
			Entries:     fastly.ToPointer(updateListEntries),
			ListID:      fastly.ToPointer(list.ListID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updateList.Description != updateListDescription {
		t.Errorf("unexpected list description: got %q, expected %q", updateList.Description, updateListDescription)
	}
	if len(updateList.Entries) != len(updateListEntries) {
		t.Errorf("unexpected list entry length: got %q, expected %q", len(updateList.Entries), len(updateListEntries))
	}
	if updateList.Entries[0] != updateListEntries[0] {
		t.Errorf("unexpected list entry: got %q, expected %q", updateList.Entries[0], updateListEntries[0])
	}
	if updateList.Name != listName {
		t.Errorf("unexpected list name: got %q, expected %q", updateList.Name, listName)
	}
	if updateList.ReferenceID != listReferenceID {
		t.Errorf("unexpected list reference ID: got %q, expected %q", updateList.ReferenceID, listReferenceID)
	}
	if updateList.Type != listType {
		t.Errorf("unexpected list type: got %q, expected %q", updateList.Type, listType)
	}
	// Scope should always be workspace
	if updateList.Scope.Type != "workspace" {
		t.Errorf("unexpected list scope: got %q, expected %q", updateList.Scope.Type, "workspace")
	}

	var lists *Lists

	// List all lists.
	fastly.Record(t, "list_lists", func(c *fastly.Client) {
		lists, err = ListLists(context.TODO(), c, &ListInput{
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lists.Data) != 1 {
		t.Errorf("unexpected lists length: got %q, expected %q", len(lists.Data), 1)
	}
	listedList := &lists.Data[0]
	if listedList.Description != updateListDescription {
		t.Errorf("unexpected list description: got %q, expected %q", listedList.Description, updateListDescription)
	}
	if len(listedList.Entries) != len(updateListEntries) {
		t.Errorf("unexpected list entry length: got %q, expected %q", len(listedList.Entries), len(updateListEntries))
	}
	if listedList.Entries[0] != updateListEntries[0] {
		t.Errorf("unexpected list entry: got %q, expected %q", listedList.Entries[0], updateListEntries[0])
	}
	if listedList.Name != listName {
		t.Errorf("unexpected list name: got %q, expected %q", listedList.Name, listName)
	}
	if listedList.ReferenceID != listReferenceID {
		t.Errorf("unexpected list reference ID: got %q, expected %q", listedList.ReferenceID, listReferenceID)
	}
	if listedList.Type != listType {
		t.Errorf("unexpected list type: got %q, expected %q", listedList.Type, listType)
	}
	// Scope should always be workspace
	if listedList.Scope.Type != "workspace" {
		t.Errorf("unexpected list scope: got %q, expected %q", listedList.Scope.Type, "workspace")
	}
}

func TestClient_CreateList_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		Entries:     nil,
	})
	if !errors.Is(err, fastly.ErrMissingEntries) {
		t.Errorf("expected ErrmissingEntries: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		Entries:     fastly.ToPointer([]string{listEntry}),
		Name:        nil,
	})
	if !errors.Is(err, fastly.ErrMissingName) {
		t.Errorf("expected ErrMissingName: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		Entries:     fastly.ToPointer([]string{listEntry}),
		Name:        fastly.ToPointer(listName),
		Type:        nil,
	})
	if !errors.Is(err, fastly.ErrMissingType) {
		t.Errorf("expected ErrMissingType: got %s", err)
	}
}

func TestClient_GetList_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		ListID:      nil,
	})
	if !errors.Is(err, fastly.ErrMissingListID) {
		t.Errorf("expected ErrMissingListID: got %s", err)
	}
}

func TestClient_UpdateList_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		ListID:      nil,
	})
	if !errors.Is(err, fastly.ErrMissingListID) {
		t.Errorf("expected ErrMissingListID: got %s", err)
	}
}

func TestClient_DeleteList_validation(t *testing.T) {
	var err error
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		ListID:      nil,
	})
	if !errors.Is(err, fastly.ErrMissingListID) {
		t.Errorf("expected ErrMissingListID: got %s", err)
	}
}
