package fastly

import (
	"errors"
	"testing"
)

func TestClient_DictionaryItems(t *testing.T) {
	fixtureBase := "dictionary_items/"
	nameSuffix := "DictionaryItems"

	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"version", *testService.ServiceID)

	testDictionary := createTestDictionary(t, fixtureBase+"dictionary", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestDictionary(t, testDictionary, fixtureBase+"delete_dictionary")

	// Create
	var err error
	var createdDictionaryItem *DictionaryItem
	Record(t, fixtureBase+"create", func(c *Client) {
		createdDictionaryItem, err = c.CreateDictionaryItem(&CreateDictionaryItemInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
			ItemKey:      ToPointer("test-dictionary-item"),
			ItemValue:    ToPointer("value"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteDictionaryItem(&DeleteDictionaryItemInput{
				ServiceID:    *testService.ServiceID,
				DictionaryID: *testDictionary.DictionaryID,
				ItemKey:      "test-dictionary-item",
			})
		})
	}()

	if *createdDictionaryItem.ItemKey != "test-dictionary-item" {
		t.Errorf("bad item_key: %q", *createdDictionaryItem.ItemKey)
	}
	if *createdDictionaryItem.ItemValue != "value" {
		t.Errorf("bad item_value: %q", *createdDictionaryItem.ItemValue)
	}

	// List
	var dictionaryItems []*DictionaryItem
	Record(t, fixtureBase+"list", func(c *Client) {
		dictionaryItems, err = c.ListDictionaryItems(&ListDictionaryItemsInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(dictionaryItems) < 1 {
		t.Errorf("bad dictionary items: %v", dictionaryItems)
	}

	// List with paginator
	var dictionaryItems2 []*DictionaryItem
	var paginator *ListPaginator[DictionaryItem]
	Record(t, fixtureBase+"list2", func(c *Client) {
		paginator = c.GetDictionaryItems(&GetDictionaryItemsInput{
			DictionaryID: *testDictionary.DictionaryID,
			Direction:    ToPointer("ascend"),
			PerPage:      ToPointer(50),
			ServiceID:    *testService.ServiceID,
			Sort:         ToPointer("item_key"),
		})

		for paginator.HasNext() {
			data, err := paginator.GetNext()
			if err != nil {
				t.Errorf("Bad paginator (remaining: %d): %s", paginator.Remaining(), err)
				return
			}
			dictionaryItems2 = append(dictionaryItems2, data...)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(dictionaryItems2) != 1 {
		t.Errorf("Bad items: %v", dictionaryItems2)
	}
	if paginator.HasNext() {
		t.Errorf("Bad paginator (remaining: %v)", paginator.Remaining())
	}

	// Get
	var retrievedDictionaryItem *DictionaryItem
	Record(t, fixtureBase+"get", func(c *Client) {
		retrievedDictionaryItem, err = c.GetDictionaryItem(&GetDictionaryItemInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
			ItemKey:      "test-dictionary-item",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *retrievedDictionaryItem.ItemKey != "test-dictionary-item" {
		t.Errorf("bad item_key: %q", *retrievedDictionaryItem.ItemKey)
	}
	if *retrievedDictionaryItem.ItemValue != "value" {
		t.Errorf("bad item_value: %q", *retrievedDictionaryItem.ItemValue)
	}

	// Update
	var updatedDictionaryItem *DictionaryItem
	Record(t, fixtureBase+"update", func(c *Client) {
		updatedDictionaryItem, err = c.UpdateDictionaryItem(&UpdateDictionaryItemInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
			ItemKey:      "test-dictionary-item",
			ItemValue:    "new-value",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *updatedDictionaryItem.ItemValue != "new-value" {
		t.Errorf("bad item_value: %q", *updatedDictionaryItem.ItemValue)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteDictionaryItem(&DeleteDictionaryItemInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
			ItemKey:      "test-dictionary-item",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDictionaryItems_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListDictionaryItems(&ListDictionaryItemsInput{})
	if !errors.Is(err, ErrMissingDictionaryID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListDictionaryItems(&ListDictionaryItemsInput{
		DictionaryID: "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDictionaryItem_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateDictionaryItem(&CreateDictionaryItemInput{})
	if !errors.Is(err, ErrMissingDictionaryID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateDictionaryItem(&CreateDictionaryItemInput{
		DictionaryID: "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDictionaryItem_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetDictionaryItem(&GetDictionaryItemInput{})
	if !errors.Is(err, ErrMissingDictionaryID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDictionaryItem(&GetDictionaryItemInput{
		DictionaryID: "test",
	})
	if !errors.Is(err, ErrMissingItemKey) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDictionaryItem(&GetDictionaryItemInput{
		DictionaryID: "test",
		ItemKey:      "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDictionaryItem_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{})
	if !errors.Is(err, ErrMissingDictionaryID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		DictionaryID: "test",
	})
	if !errors.Is(err, ErrMissingItemKey) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		DictionaryID: "test",
		ItemKey:      "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDictionaryItem_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{})
	if !errors.Is(err, ErrMissingDictionaryID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		DictionaryID: "test",
	})
	if !errors.Is(err, ErrMissingItemKey) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		DictionaryID: "test",
		ItemKey:      "123",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_BatchModifyDictionaryItem_validation(t *testing.T) {
	var err error

	err = TestClient.BatchModifyDictionaryItems(&BatchModifyDictionaryItemsInput{})
	if !errors.Is(err, ErrMissingDictionaryID) {
		t.Errorf("bad error: %s", err)
	}

	oversizedDictionaryItems := make([]*BatchDictionaryItem, BatchModifyMaximumOperations+1)
	err = TestClient.BatchModifyDictionaryItems(&BatchModifyDictionaryItemsInput{
		DictionaryID: "bar",
		Items:        oversizedDictionaryItems,
	})
	if !errors.Is(err, ErrMaxExceededItems) {
		t.Errorf("bad error: %s", err)
	}

	validDictionaryItems := make([]*BatchDictionaryItem, BatchModifyMaximumOperations)
	err = TestClient.BatchModifyDictionaryItems(&BatchModifyDictionaryItemsInput{
		DictionaryID: "bar",
		Items:        validDictionaryItems,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}
