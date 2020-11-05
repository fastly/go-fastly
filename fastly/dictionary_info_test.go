package fastly

import (
	"testing"
)

func TestClient_GetDictionaryInfo(t *testing.T) {

	fixtureBase := "dictionary_info/"
	nameSuffix := "DictionaryInfo"

	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", testService.ID)

	testVersion := createTestVersion(t, fixtureBase+"version", testService.ID)

	testDictionary := createTestDictionary(t, fixtureBase+"dictionary", testService.ID, testVersion.Number, nameSuffix)
	defer deleteTestDictionary(t, testDictionary, fixtureBase+"delete_dictionary")

	var (
		err  error
		info *DictionaryInfo
	)

	record(t, fixtureBase+"create_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(&BatchModifyDictionaryItemsInput{
			ServiceID:    testService.ID,
			DictionaryID: testDictionary.ID,
			Items: []*BatchDictionaryItem{
				{
					Operation: CreateBatchOperation,
					ItemKey:   "test-dictionary-item-0",
					ItemValue: "value",
				},
				{
					Operation: CreateBatchOperation,
					ItemKey:   "test-dictionary-item-1",
					ItemValue: "value",
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, fixtureBase+"get", func(c *Client) {
		info, err = c.GetDictionaryInfo(&GetDictionaryInfoInput{
			ServiceID:      testService.ID,
			ServiceVersion: testVersion.Number,
			ID:             testDictionary.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if info.ItemCount != 2 {
		t.Errorf("bad item_count: %d", info.ItemCount)
	}
}

func TestClient_GetDictionaryInfo_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		ID:             "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
