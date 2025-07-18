package lists

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/common"
)

const (
	listDescription string = "This is a list"
	listEntry       string = "entry1"
	listName        string = "List Name "
	listReferenceID string = "list-name-"
	listType        string = "string"
)

func TestClient_Lists_WorkspaceScope(t *testing.T) {
	runListsTest(t, common.ScopeTypeWorkspace, fastly.TestNGWAFWorkspaceID)
}

func TestClient_Lists_AccountScope(t *testing.T) {
	runListsTest(t, common.ScopeTypeAccount, "*") // assuming TestNGWAFAccountID exists
}

func runListsTest(t *testing.T, scopeType common.ScopeType, appliesToID string) {
	var err error
	listEntries := []string{listEntry}
	testListName := listName + string(scopeType)
	testListReferenceID := listReferenceID + string(scopeType)

	// Create a test list.
	var list *List
	fastly.Record(t, fmt.Sprintf("%s_create_list", scopeType), func(c *fastly.Client) {
		list, err = Create(context.TODO(), c, &CreateInput{
			Description: fastly.ToPointer(listDescription),
			Entries:     fastly.ToPointer(listEntries),
			Name:        fastly.ToPointer(testListName),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			Type: fastly.ToPointer(listType),
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
	if list.Name != testListName {
		t.Errorf("unexpected list name: got %q, expected %q", list.Name, testListName)
	}
	if !strings.Contains(list.ReferenceID, testListReferenceID) {
		t.Errorf("unexpected list reference ID: expected %q to contain %q", list.ReferenceID, testListReferenceID)
	}
	if list.Type != listType {
		t.Errorf("unexpected list type: got %q, expected %q", list.Type, listType)
	}
	if list.Scope.Type != string(scopeType) {
		t.Errorf("unexpected list scope: got %q, expected %q", list.Scope.Type, string(scopeType))
	}

	// Ensure we delete the test list at the end.
	defer func() {
		fastly.Record(t, fmt.Sprintf("%s_delete_list", scopeType), func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				ListID: fastly.ToPointer(list.ListID),
				Scope: &common.Scope{
					Type:      scopeType,
					AppliesTo: []string{appliesToID},
				},
			})
		})
		if err != nil {
			t.Errorf("error during list cleanup: %v", err)
		}
	}()

	// Get the test list.
	var getList *List
	fastly.Record(t, fmt.Sprintf("%s_get_list", scopeType), func(c *fastly.Client) {
		getList, err = Get(context.TODO(), c, &GetInput{
			ListID: fastly.ToPointer(list.ListID),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
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
	if getList.Name != testListName {
		t.Errorf("unexpected list name: got %q, expected %q", getList.Name, testListName)
	}
	if !strings.Contains(getList.ReferenceID, testListReferenceID) {
		t.Errorf("unexpected list reference ID: expected %q to contain %q", getList.ReferenceID, testListReferenceID)
	}
	if getList.Type != listType {
		t.Errorf("unexpected list type: got %q, expected %q", getList.Type, listType)
	}
	if getList.Scope.Type != string(scopeType) {
		t.Errorf("unexpected list scope: got %q, expected %q", getList.Scope.Type, string(scopeType))
	}

	// Update the test list.
	updateListDescription := "This is an updated list"
	updateListEntries := []string{
		"127.0.0.1",
	}

	var updateList *List
	fastly.Record(t, fmt.Sprintf("%s_update_list", scopeType), func(c *fastly.Client) {
		updateList, err = Update(context.TODO(), c, &UpdateInput{
			Description: fastly.ToPointer(updateListDescription),
			Entries:     fastly.ToPointer(updateListEntries),
			ListID:      fastly.ToPointer(list.ListID),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
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
	if updateList.Name != testListName {
		t.Errorf("unexpected list name: got %q, expected %q", updateList.Name, testListName)
	}
	if !strings.Contains(updateList.ReferenceID, testListReferenceID) {
		t.Errorf("unexpected list reference ID: expected %q to contain %q", updateList.ReferenceID, testListReferenceID)
	}
	if updateList.Type != listType {
		t.Errorf("unexpected list type: got %q, expected %q", updateList.Type, listType)
	}
	if updateList.Scope.Type != string(scopeType) {
		t.Errorf("unexpected list scope: got %q, expected %q", updateList.Scope.Type, string(scopeType))
	}

	var lists *Lists

	// List all lists.
	fastly.Record(t, fmt.Sprintf("%s_list_lists", scopeType), func(c *fastly.Client) {
		lists, err = ListLists(context.TODO(), c, &ListInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
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
	if listedList.Name != testListName {
		t.Errorf("unexpected list name: got %q, expected %q", listedList.Name, listName)
	}
	if !strings.Contains(listedList.ReferenceID, testListReferenceID) {
		t.Errorf("unexpected list reference ID: expected %q to contain %q", listedList.ReferenceID, listReferenceID)
	}
	if listedList.Type != listType {
		t.Errorf("unexpected list type: got %q, expected %q", listedList.Type, listType)
	}
	if listedList.Scope.Type != string(scopeType) {
		t.Errorf("unexpected list scope: got %q, expected %q", listedList.Scope.Type, string(scopeType))
	}
}

func TestClient_CreateList_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Entries: nil,
	})
	if !errors.Is(err, fastly.ErrMissingEntries) {
		t.Errorf("expected ErrmissingEntries: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Entries: fastly.ToPointer([]string{listEntry}),
		Name:    nil,
	})
	if !errors.Is(err, fastly.ErrMissingName) {
		t.Errorf("expected ErrMissingName: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Entries: fastly.ToPointer([]string{listEntry}),
		Name:    fastly.ToPointer(listName),
		Type:    nil,
	})
	if !errors.Is(err, fastly.ErrMissingType) {
		t.Errorf("expected ErrMissingType: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Entries: fastly.ToPointer([]string{listEntry}),
		Name:    fastly.ToPointer(listName),
		Type:    fastly.ToPointer(listType),
		Scope:   nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}

func TestClient_GetList_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		ListID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingListID) {
		t.Errorf("expected ErrMissingListID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		ListID: fastly.ToPointer("someID"),
		Scope:  nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}

func TestClient_UpdateList_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		ListID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingListID) {
		t.Errorf("expected ErrMissingListID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		ListID: fastly.ToPointer("someID"),
		Scope:  nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}

func TestClient_DeleteList_validation(t *testing.T) {
	var err error
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		ListID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingListID) {
		t.Errorf("expected ErrMissingListID: got %s", err)
	}
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		ListID: fastly.ToPointer("someID"),
		Scope:  nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}
