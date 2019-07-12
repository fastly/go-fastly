package fastly

import (
	"fmt"
	"testing"
	"time"
)

func createTestService(t *testing.T, serviceFixture string) *Service {

	var err error
	var service *Service

	record(t, serviceFixture, func(client *Client) {
		service, err = client.CreateService(&CreateServiceInput{
			Name: fmt.Sprintf("test_service_%d", time.Now().Unix()),
			Comment: "go-fastly client test",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	return service
}

func createTestVersion(t *testing.T, versionFixture string, serviceId string) *Version {

	var err error
	var version *Version

	record(t, versionFixture, func(client *Client) {
		testVersionLock.Lock()
		defer testVersionLock.Unlock()

		version, err = client.CreateVersion(&CreateVersionInput{
			Service: serviceId,
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	return version
}

func createTestDictionary(t *testing.T, dictionaryFixture string, serviceId string, version int) *Dictionary {

	var err error
	var dictionary *Dictionary

	record(t, dictionaryFixture, func(client *Client) {
		dictionary, err = client.CreateDictionary(&CreateDictionaryInput{
			Service: serviceId,
			Version: version,
			Name: fmt.Sprintf("test_dictionary_%d", time.Now().Unix()),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return dictionary
}

func deleteTestDictionary(t *testing.T, dictionary *Dictionary, deleteFixture string) {

	var err error

	record(t, deleteFixture, func(client *Client) {
		err = client.DeleteDictionary(&DeleteDictionaryInput{
			Service: dictionary.ServiceID,
			Version: dictionary.Version,
			Name:    dictionary.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestService(t *testing.T, cleanupFixture string, serviceId string){

	var err error

	record(t, cleanupFixture, func(client *Client) {
		err = client.DeleteService(&DeleteServiceInput{
			ID: serviceId,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_BatchModifyDictionaryItems_Create(t *testing.T) {

	fixtureBase := "dictionary_items_batch/create/"

	// Given: a test service with a dictionary and a batch of create operations,
	testService := createTestService(t, fixtureBase + "create_service")
	defer deleteTestService(t, fixtureBase +"delete_service", testService.ID)

	testVersion := createTestVersion(t,fixtureBase + "create_version", testService.ID)

	testDictionary := createTestDictionary(t, fixtureBase + "create_dictionary", testService.ID, testVersion.Number)
	defer deleteTestDictionary(t, testDictionary, fixtureBase + "delete_dictionary")

	batchCreateOperations := &BatchModifyDictionaryItemsInput {
		Service:    testService.ID,
		Dictionary: testDictionary.ID,
		Items: []BatchDictionaryItem{
			{
				Operation:Create,
				ItemKey: "key1",
				ItemValue: "val1",
			},
			{
				Operation:Create,
				ItemKey: "key2",
				ItemValue: "val2",
			},
		},
	}

	// When: I execute the batch create operations against the Fastly API,
	var err error
	record(t, fixtureBase + "create_dictionary_items", func(c *Client) {

		err = c.BatchModifyDictionaryItems(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the created dictionary items.
	var actualDictionaryItems []*DictionaryItem
	record(t, fixtureBase + "list_after_create", func(c *Client) {
		actualDictionaryItems, err = c.ListDictionaryItems(&ListDictionaryItemsInput{
			Service:    testService.ID,
			Dictionary: testDictionary.ID,
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
		if actualItemKey != expectedItemKey {
			t.Errorf("First ItemKey did not match, expected %s, got %s", expectedItemKey, actualItemKey)
		}

		actualItemValue := item.ItemValue
		expectedItemValue := batchCreateOperations.Items[i].ItemValue
		if actualItemValue != expectedItemValue {
			t.Errorf("First ItemValue did not match, expected %s, got %s", expectedItemValue, actualItemValue)
		}

	}

}

func TestClient_BatchModifyDictionaryItems_Delete(t *testing.T) {

	fixtureBase := "dictionary_items_batch/delete/"

	// Given: a test service with a dictionary and dictionary items,
	testService := createTestService(t, fixtureBase + "create_service")
	defer deleteTestService(t, fixtureBase + "delete_service", testService.ID)

	testVersion := createTestVersion(t,fixtureBase + "create_version", testService.ID)

	testDictionary := createTestDictionary(t, fixtureBase + "create_dictionary", testService.ID, testVersion.Number)
	defer deleteTestDictionary(t, testDictionary, fixtureBase + "delete_dictionary")

	batchCreateOperations := &BatchModifyDictionaryItemsInput {
		Service:    testService.ID,
		Dictionary: testDictionary.ID,
		Items: []BatchDictionaryItem{
			{
				Operation:Create,
				ItemKey: "key1",
				ItemValue: "val1",
			},
			{
				Operation:Create,
				ItemKey: "key2",
				ItemValue: "val2",
			},
		},
	}

	var err error
	record(t, fixtureBase + "create_dictionary_items", func(c *Client) {

		err = c.BatchModifyDictionaryItems(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// When: I execute the batch delete operations against the Fastly API,
	batchDeleteOperations := &BatchModifyDictionaryItemsInput {
		Service:    testService.ID,
		Dictionary: testDictionary.ID,
		Items: []BatchDictionaryItem{
			{
				Operation:Delete,
				ItemKey: "key2",
				ItemValue: "val2",
			},
		},
	}

	record(t, fixtureBase + "delete_dictionary_items", func(c *Client) {

		err = c.BatchModifyDictionaryItems(batchDeleteOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list a single dictionary item.
	var actualDictionaryItems []*DictionaryItem
	record(t, fixtureBase + "list_after_delete", func(client *Client) {
		actualDictionaryItems, err = client.ListDictionaryItems(&ListDictionaryItemsInput{
			Service:    testService.ID,
			Dictionary: testDictionary.ID,
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

	// Given: a test service with a dictionary and dictionary items,
	testService := createTestService(t, fixtureBase + "create_service")
	defer deleteTestService(t, fixtureBase + "delete_service", testService.ID)

	testVersion := createTestVersion(t,fixtureBase + "create_version", testService.ID)

	testDictionary := createTestDictionary(t, fixtureBase + "create_dictionary", testService.ID, testVersion.Number)
	defer deleteTestDictionary(t, testDictionary, fixtureBase + "delete_dictionary")

	batchCreateOperations := &BatchModifyDictionaryItemsInput {
		Service:    testService.ID,
		Dictionary: testDictionary.ID,
		Items: []BatchDictionaryItem{
			{
				Operation:Create,
				ItemKey: "key1",
				ItemValue: "val1",
			},
			{
				Operation:Create,
				ItemKey: "key2",
				ItemValue: "val2",
			},
		},
	}

	var err error
	record(t, fixtureBase + "create_dictionary_items", func(c *Client) {

		err = c.BatchModifyDictionaryItems(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// When: I execute the batch update operations against the Fastly API,
	batchUpdateOperations := &BatchModifyDictionaryItemsInput {
		Service:    testService.ID,
		Dictionary: testDictionary.ID,
		Items: []BatchDictionaryItem{
			{
				Operation:Update,
				ItemKey: "key2",
				ItemValue: "val2Updated",
			},
		},
	}

	record(t, fixtureBase + "update_dictionary_items", func(c *Client) {

		err = c.BatchModifyDictionaryItems(batchUpdateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then: I expect to be able to list all of the dictionary items with modifications applied to a single item.
	var actualDictionaryItems []*DictionaryItem
	record(t, fixtureBase + "list_after_update", func(c *Client) {
		actualDictionaryItems, err = c.ListDictionaryItems(&ListDictionaryItemsInput{
			Service:    testService.ID,
			Dictionary: testDictionary.ID,
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

	if actualItemKey != expectedItemKey {
		t.Errorf("First ItemKey did not match, expected %s, got %s", expectedItemKey, actualItemKey)
	}

	actualItemValue := actualDictionaryItems[0].ItemValue
	expectedItemValue := batchCreateOperations.Items[0].ItemValue


	// Confirm the second dictionary item contains the modifications.
	if actualItemValue != expectedItemValue {
		t.Errorf("First ItemValue did not match, expected %s, got %s", expectedItemValue, actualItemValue)
	}

	actualItemKey = actualDictionaryItems[1].ItemKey
	expectedItemKey = batchUpdateOperations.Items[0].ItemKey

	if actualItemKey != expectedItemKey {
		t.Errorf("Second ItemKey did not match, expected %s, got %s", expectedItemKey, actualItemKey)
	}

	actualItemValue = actualDictionaryItems[1].ItemValue
	expectedItemValue = batchUpdateOperations.Items[0].ItemValue


	if actualItemValue != expectedItemValue {
		t.Errorf("Second ItemValue did not match, expected %s, got %s", expectedItemValue, actualItemValue)
	}

}

func TestClient_BatchModifyDictionaryItems_Upsert(t *testing.T) {

	fixtureBase := "dictionary_items_batch/upsert/"

	// Given: a test service with a dictionary and dictionary items,
	testService := createTestService(t, fixtureBase + "create_service")
	defer deleteTestService(t, fixtureBase + "delete_service", testService.ID)

	testVersion := createTestVersion(t,fixtureBase + "create_version", testService.ID)

	testDictionary := createTestDictionary(t, fixtureBase + "create_dictionary", testService.ID, testVersion.Number)
	defer deleteTestDictionary(t, testDictionary, fixtureBase + "delete_dictionary")

	batchCreateOperations := &BatchModifyDictionaryItemsInput {
		Service:    testService.ID,
		Dictionary: testDictionary.ID,
		Items: []BatchDictionaryItem{
			{
				Operation:Create,
				ItemKey: "key1",
				ItemValue: "val1",
			},
		},
	}

	var err error
	record(t, fixtureBase + "create_dictionary_items", func(c *Client) {

		err = c.BatchModifyDictionaryItems(batchCreateOperations)
	})
	if err != nil {
		t.Fatal(err)
	}

	// When: I execute the batch upsert operations against the Fastly API
	batchUpsertOperations := &BatchModifyDictionaryItemsInput {
		Service:    testService.ID,
		Dictionary: testDictionary.ID,
		Items: []BatchDictionaryItem{
			{
				Operation:Upsert,
				ItemKey: "key1",
				ItemValue: "val1Updated",
			},
			{
				Operation:Upsert,
				ItemKey: "key2",
				ItemValue: "val2",
			},
		},
	}

	record(t, fixtureBase + "upsert_dictionary_items", func(c *Client) {

		err = c.BatchModifyDictionaryItems(batchUpsertOperations)
	})
	if err != nil {
		t.Fatal(err)
	}


	// Then: I expect to be able to list all of the dictionary items with the modification present.
	var actualDictionaryItems []*DictionaryItem
	record(t, fixtureBase + "list_after_upsert", func(c *Client) {
		actualDictionaryItems, err = c.ListDictionaryItems(&ListDictionaryItemsInput{
			Service:    testService.ID,
			Dictionary: testDictionary.ID,
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
		if actualItemKey != expectedItemKey {
			t.Errorf("First ItemKey did not match, expected %s, got %s", expectedItemKey, actualItemKey)
		}

		actualItemValue := item.ItemValue
		expectedItemValue := batchUpsertOperations.Items[i].ItemValue
		if actualItemValue != expectedItemValue {
			t.Errorf("First ItemValue did not match, expected %s, got %s", expectedItemValue, actualItemValue)
		}

	}

}