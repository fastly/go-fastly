package fastly

import (
	"context"
	"testing"
)

func TestClient_BatchModifyDictionaryItems_Create(t *testing.T) {
	fixtureBase := "dictionary_items_batch/create/"
	nameSuffix := "BatchModifyDictionaryItems_Create"

	// Given: a test service with a dictionary and a batch of create operations,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"create_version", *testService.ServiceID)

	testDictionary := createTestDictionary(t, fixtureBase+"create_dictionary", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestDictionary(t, testDictionary, fixtureBase+"delete_dictionary")

	batchCreateOperations := &BatchModifyDictionaryItemsInput{
		ServiceID:    *testService.ServiceID,
		DictionaryID: *testDictionary.DictionaryID,
		Items: []*BatchDictionaryItem{
			{
				Operation: ToPointer(CreateBatchOperation),
				ItemKey:   ToPointer("key1"),
				ItemValue: ToPointer("val1"),
			},
			{
				Operation: ToPointer(CreateBatchOperation),
				ItemKey:   ToPointer("key2"),
				ItemValue: ToPointer("val2"),
			},
		},
	}

	// When: I execute the batch create operations against the Fastly API,
	var err error
	Record(t, fixtureBase+"create_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(context.TODO(), batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the created dictionary items.
	var actualDictionaryItems []*DictionaryItem
	Record(t, fixtureBase+"list_after_create", func(c *Client) {
		actualDictionaryItems, err = c.ListDictionaryItems(context.TODO(), &ListDictionaryItemsInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	actualNumberOfDictItems := len(actualDictionaryItems)
	expectedNumberOfDictItems := len(batchCreateOperations.Items)
	if actualNumberOfDictItems != expectedNumberOfDictItems {
		t.Errorf("Incorrect number of dictionary items returned, expected: %d, got %d", expectedNumberOfDictItems, actualNumberOfDictItems)
	}

	for i, item := range actualDictionaryItems {
		actualItemKey := item.ItemKey
		expectedItemKey := batchCreateOperations.Items[i].ItemKey
		if *actualItemKey != *expectedItemKey {
			t.Errorf("First ItemKey did not match, expected %s, got %s", *expectedItemKey, *actualItemKey)
		}

		actualItemValue := item.ItemValue
		expectedItemValue := batchCreateOperations.Items[i].ItemValue
		if *actualItemValue != *expectedItemValue {
			t.Errorf("First ItemValue did not match, expected %s, got %s", *expectedItemValue, *actualItemValue)
		}
	}
}

func TestClient_BatchModifyDictionaryItems_Delete(t *testing.T) {
	fixtureBase := "dictionary_items_batch/delete/"
	nameSuffix := "BatchModifyDictionaryItems_Delete"

	// Given: a test service with a dictionary and dictionary items,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"create_version", *testService.ServiceID)

	testDictionary := createTestDictionary(t, fixtureBase+"create_dictionary", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestDictionary(t, testDictionary, fixtureBase+"delete_dictionary")

	batchCreateOperations := &BatchModifyDictionaryItemsInput{
		ServiceID:    *testService.ServiceID,
		DictionaryID: *testDictionary.DictionaryID,
		Items: []*BatchDictionaryItem{
			{
				Operation: ToPointer(CreateBatchOperation),
				ItemKey:   ToPointer("key1"),
				ItemValue: ToPointer("val1"),
			},
			{
				Operation: ToPointer(CreateBatchOperation),
				ItemKey:   ToPointer("key2"),
				ItemValue: ToPointer("val2"),
			},
		},
	}

	var err error
	Record(t, fixtureBase+"create_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(context.TODO(), batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// When: I execute the batch delete operations against the Fastly API,
	batchDeleteOperations := &BatchModifyDictionaryItemsInput{
		ServiceID:    *testService.ServiceID,
		DictionaryID: *testDictionary.DictionaryID,
		Items: []*BatchDictionaryItem{
			{
				Operation: ToPointer(DeleteBatchOperation),
				ItemKey:   ToPointer("key2"),
				ItemValue: ToPointer("val2"),
			},
		},
	}

	Record(t, fixtureBase+"delete_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(context.TODO(), batchDeleteOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list a single dictionary item.
	var actualDictionaryItems []*DictionaryItem
	Record(t, fixtureBase+"list_after_delete", func(client *Client) {
		actualDictionaryItems, err = client.ListDictionaryItems(context.TODO(), &ListDictionaryItemsInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	actualNumberOfDictItems := len(actualDictionaryItems)
	expectedNumberOfDictItems := len(batchDeleteOperations.Items)
	if actualNumberOfDictItems != expectedNumberOfDictItems {
		t.Errorf("Incorrect number of dictionary items returned, expected: %d, got %d", expectedNumberOfDictItems, actualNumberOfDictItems)
	}
}

func TestClient_BatchModifyDictionaryItems_Update(t *testing.T) {
	fixtureBase := "dictionary_items_batch/update/"
	nameSuffix := "BatchModifyDictionaryItems_Update"

	// Given: a test service with a dictionary and dictionary items,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"create_version", *testService.ServiceID)

	testDictionary := createTestDictionary(t, fixtureBase+"create_dictionary", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestDictionary(t, testDictionary, fixtureBase+"delete_dictionary")

	batchCreateOperations := &BatchModifyDictionaryItemsInput{
		ServiceID:    *testService.ServiceID,
		DictionaryID: *testDictionary.DictionaryID,
		Items: []*BatchDictionaryItem{
			{
				Operation: ToPointer(CreateBatchOperation),
				ItemKey:   ToPointer("key1"),
				ItemValue: ToPointer("val1"),
			},
			{
				Operation: ToPointer(CreateBatchOperation),
				ItemKey:   ToPointer("key2"),
				ItemValue: ToPointer("val2"),
			},
		},
	}

	var err error
	Record(t, fixtureBase+"create_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(context.TODO(), batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// When: I execute the batch update operations against the Fastly API,
	batchUpdateOperations := &BatchModifyDictionaryItemsInput{
		ServiceID:    *testService.ServiceID,
		DictionaryID: *testDictionary.DictionaryID,
		Items: []*BatchDictionaryItem{
			{
				Operation: ToPointer(UpdateBatchOperation),
				ItemKey:   ToPointer("key2"),
				ItemValue: ToPointer("val2Updated"),
			},
		},
	}

	Record(t, fixtureBase+"update_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(context.TODO(), batchUpdateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the dictionary items with modifications applied to a single item.
	var actualDictionaryItems []*DictionaryItem
	Record(t, fixtureBase+"list_after_update", func(c *Client) {
		actualDictionaryItems, err = c.ListDictionaryItems(context.TODO(), &ListDictionaryItemsInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	actualNumberOfDictItems := len(actualDictionaryItems)
	expectedNumberOfDictItems := len(batchCreateOperations.Items)
	if actualNumberOfDictItems != expectedNumberOfDictItems {
		t.Errorf("Incorrect number of dictionary items returned, expected: %d, got %d", expectedNumberOfDictItems, actualNumberOfDictItems)
	}

	actualItemKey := actualDictionaryItems[0].ItemKey
	expectedItemKey := batchCreateOperations.Items[0].ItemKey

	if *actualItemKey != *expectedItemKey {
		t.Errorf("First ItemKey did not match, expected %s, got %s", *expectedItemKey, *actualItemKey)
	}

	actualItemValue := actualDictionaryItems[0].ItemValue
	expectedItemValue := batchCreateOperations.Items[0].ItemValue

	// Confirm the second dictionary item contains the modifications.
	if *actualItemValue != *expectedItemValue {
		t.Errorf("First ItemValue did not match, expected %s, got %s", *expectedItemValue, *actualItemValue)
	}

	actualItemKey = actualDictionaryItems[1].ItemKey
	expectedItemKey = batchUpdateOperations.Items[0].ItemKey

	if *actualItemKey != *expectedItemKey {
		t.Errorf("Second ItemKey did not match, expected %s, got %s", *expectedItemKey, *actualItemKey)
	}

	actualItemValue = actualDictionaryItems[1].ItemValue
	expectedItemValue = batchUpdateOperations.Items[0].ItemValue

	if *actualItemValue != *expectedItemValue {
		t.Errorf("Second ItemValue did not match, expected %s, got %s", *expectedItemValue, *actualItemValue)
	}
}

func TestClient_BatchModifyDictionaryItems_Upsert(t *testing.T) {
	fixtureBase := "dictionary_items_batch/upsert/"
	nameSuffix := "BatchModifyDictionaryItems_Upsert"

	// Given: a test service with a dictionary and dictionary items,
	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := CreateTestVersion(t, fixtureBase+"create_version", *testService.ServiceID)

	testDictionary := createTestDictionary(t, fixtureBase+"create_dictionary", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestDictionary(t, testDictionary, fixtureBase+"delete_dictionary")

	batchCreateOperations := &BatchModifyDictionaryItemsInput{
		ServiceID:    *testService.ServiceID,
		DictionaryID: *testDictionary.DictionaryID,
		Items: []*BatchDictionaryItem{
			{
				Operation: ToPointer(CreateBatchOperation),
				ItemKey:   ToPointer("key1"),
				ItemValue: ToPointer("val1"),
			},
		},
	}

	var err error
	Record(t, fixtureBase+"create_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(context.TODO(), batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// When: I execute the batch upsert operations against the Fastly API
	batchUpsertOperations := &BatchModifyDictionaryItemsInput{
		ServiceID:    *testService.ServiceID,
		DictionaryID: *testDictionary.DictionaryID,
		Items: []*BatchDictionaryItem{
			{
				Operation: ToPointer(UpsertBatchOperation),
				ItemKey:   ToPointer("key1"),
				ItemValue: ToPointer("val1Updated"),
			},
			{
				Operation: ToPointer(UpsertBatchOperation),
				ItemKey:   ToPointer("key2"),
				ItemValue: ToPointer("val2"),
			},
		},
	}

	Record(t, fixtureBase+"upsert_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(context.TODO(), batchUpsertOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the dictionary items with the modification present.
	var actualDictionaryItems []*DictionaryItem
	Record(t, fixtureBase+"list_after_upsert", func(c *Client) {
		actualDictionaryItems, err = c.ListDictionaryItems(context.TODO(), &ListDictionaryItemsInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	actualNumberOfDictItems := len(actualDictionaryItems)
	expectedNumberOfDictItems := len(batchUpsertOperations.Items)
	if actualNumberOfDictItems != expectedNumberOfDictItems {
		t.Errorf("Incorrect number of dictionary items returned, expected: %d, got %d", expectedNumberOfDictItems, actualNumberOfDictItems)
	}

	for i, item := range actualDictionaryItems {
		actualItemKey := item.ItemKey
		expectedItemKey := batchUpsertOperations.Items[i].ItemKey
		if *actualItemKey != *expectedItemKey {
			t.Errorf("First ItemKey did not match, expected %s, got %s", *expectedItemKey, *actualItemKey)
		}

		actualItemValue := item.ItemValue
		expectedItemValue := batchUpsertOperations.Items[i].ItemValue
		if *actualItemValue != *expectedItemValue {
			t.Errorf("First ItemValue did not match, expected %s, got %s", *expectedItemValue, *actualItemValue)
		}
	}
}
