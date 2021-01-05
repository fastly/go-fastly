package fastly

import (
	"testing"
)

func TestClient_DictionaryItems(t *testing.T) {

	fixtureBase := "dictionary_items/"
	nameSuffix := "DictionaryItems"

	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", testService.ID)

	testVersion := createTestVersion(t, fixtureBase+"version", testService.ID)

	testDictionary := createTestDictionary(t, fixtureBase+"dictionary", testService.ID, testVersion.Number, nameSuffix)
	defer deleteTestDictionary(t, testDictionary, fixtureBase+"delete_dictionary")

	// Create
	var err error
	var createdDictionaryItem *DictionaryItem
	record(t, fixtureBase+"create", func(c *Client) {
		createdDictionaryItem, err = c.CreateDictionaryItem(&CreateDictionaryItemInput{
			ServiceID:    testService.ID,
			DictionaryID: testDictionary.ID,
			ItemKey:      "test-dictionary-item",
			ItemValue:    "value",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteDictionaryItem(&DeleteDictionaryItemInput{
				ServiceID:    testService.ID,
				DictionaryID: testDictionary.ID,
				ItemKey:      "test-dictionary-item",
			})
		})
	}()

	if createdDictionaryItem.ItemKey != "test-dictionary-item" {
		t.Errorf("bad item_key: %q", createdDictionaryItem.ItemKey)
	}
	if createdDictionaryItem.ItemValue != "value" {
		t.Errorf("bad item_value: %q", createdDictionaryItem.ItemValue)
	}

	// List
	var dictionaryItems []*DictionaryItem
	record(t, fixtureBase+"list", func(c *Client) {
		dictionaryItems, err = c.ListDictionaryItems(&ListDictionaryItemsInput{
			ServiceID:    testService.ID,
			DictionaryID: testDictionary.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(dictionaryItems) < 1 {
		t.Errorf("bad dictionary items: %v", dictionaryItems)
	}

	// Get
	var retrievedDictionaryItem *DictionaryItem
	record(t, fixtureBase+"get", func(c *Client) {
		retrievedDictionaryItem, err = c.GetDictionaryItem(&GetDictionaryItemInput{
			ServiceID:    testService.ID,
			DictionaryID: testDictionary.ID,
			ItemKey:      "test-dictionary-item",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if retrievedDictionaryItem.ItemKey != "test-dictionary-item" {
		t.Errorf("bad item_key: %q", retrievedDictionaryItem.ItemKey)
	}
	if retrievedDictionaryItem.ItemValue != "value" {
		t.Errorf("bad item_value: %q", retrievedDictionaryItem.ItemValue)
	}

	// Update
	var updatedDictionaryItem *DictionaryItem
	record(t, fixtureBase+"update", func(c *Client) {
		updatedDictionaryItem, err = c.UpdateDictionaryItem(&UpdateDictionaryItemInput{
			ServiceID:    testService.ID,
			DictionaryID: testDictionary.ID,
			ItemKey:      "test-dictionary-item",
			ItemValue:    "new-value",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedDictionaryItem.ItemValue != "new-value" {
		t.Errorf("bad item_value: %q", updatedDictionaryItem.ItemValue)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteDictionaryItem(&DeleteDictionaryItemInput{
			ServiceID:    testService.ID,
			DictionaryID: testDictionary.ID,
			ItemKey:      "test-dictionary-item",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDictionaryItems_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDictionaryItems(&ListDictionaryItemsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDictionaryItems(&ListDictionaryItemsInput{
		ServiceID:    "foo",
		DictionaryID: "",
	})
	if err != ErrMissingDictionaryID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDictionaryItem_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDictionaryItem(&CreateDictionaryItemInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDictionaryItem(&CreateDictionaryItemInput{
		ServiceID:    "foo",
		DictionaryID: "",
	})
	if err != ErrMissingDictionaryID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDictionaryItem_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDictionaryItem(&GetDictionaryItemInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryItem(&GetDictionaryItemInput{
		ServiceID:    "foo",
		DictionaryID: "",
	})
	if err != ErrMissingDictionaryID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryItem(&GetDictionaryItemInput{
		ServiceID:    "foo",
		DictionaryID: "test",
		ItemKey:      "",
	})
	if err != ErrMissingItemKey {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDictionaryItem_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		ServiceID:    "foo",
		DictionaryID: "",
	})
	if err != ErrMissingDictionaryID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		ServiceID:    "foo",
		DictionaryID: "test",
		ItemKey:      "",
	})
	if err != ErrMissingItemKey {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDictionaryItem_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		ServiceID:    "foo",
		DictionaryID: "",
	})
	if err != ErrMissingDictionaryID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		ServiceID:    "foo",
		DictionaryID: "test",
		ItemKey:      "",
	})
	if err != ErrMissingItemKey {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_BatchModifyDictionaryItem_validation(t *testing.T) {
	var err error
	err = testClient.BatchModifyDictionaryItems(&BatchModifyDictionaryItemsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}
	err = testClient.BatchModifyDictionaryItems(&BatchModifyDictionaryItemsInput{
		ServiceID:    "foo",
		DictionaryID: "",
	})
	if err != ErrMissingDictionaryID {
		t.Errorf("bad error: %s", err)
	}

	oversizedDictionaryItems := make([]*BatchDictionaryItem, BatchModifyMaximumOperations+1)
	err = testClient.BatchModifyDictionaryItems(&BatchModifyDictionaryItemsInput{
		ServiceID:    "foo",
		DictionaryID: "bar",
		Items:        oversizedDictionaryItems,
	})
	if err != ErrMaxExceededItems {
		t.Errorf("bad error: %s", err)
	}
}
