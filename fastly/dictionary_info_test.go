package fastly

import (
	"testing"
)

func TestClient_GetDictionaryInfo(t *testing.T) {
	fixtureBase := "dictionary_info/"
	nameSuffix := "DictionaryInfo"

	testService := createTestService(t, fixtureBase+"create_service", nameSuffix)
	defer deleteTestService(t, fixtureBase+"delete_service", *testService.ServiceID)

	testVersion := createTestVersion(t, fixtureBase+"version", *testService.ServiceID)

	testDictionary := createTestDictionary(t, fixtureBase+"dictionary", *testService.ServiceID, *testVersion.Number, nameSuffix)
	defer deleteTestDictionary(t, testDictionary, fixtureBase+"delete_dictionary")

	var (
		err  error
		info *DictionaryInfo
	)

	record(t, fixtureBase+"create_dictionary_items", func(c *Client) {
		err = c.BatchModifyDictionaryItems(&BatchModifyDictionaryItemsInput{
			ServiceID:    *testService.ServiceID,
			DictionaryID: *testDictionary.DictionaryID,
			Items: []*BatchDictionaryItem{
				{
					Operation: ToPointer(CreateBatchOperation),
					ItemKey:   ToPointer("test-dictionary-item-0"),
					ItemValue: ToPointer("value"),
				},
				{
					Operation: ToPointer(CreateBatchOperation),
					ItemKey:   ToPointer("test-dictionary-item-1"),
					ItemValue: ToPointer("value"),
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, fixtureBase+"get", func(c *Client) {
		info, err = c.GetDictionaryInfo(&GetDictionaryInfoInput{
			ServiceID:      *testService.ServiceID,
			ServiceVersion: *testVersion.Number,
			DictID:         *testDictionary.DictionaryID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *info.ItemCount != 2 {
		t.Errorf("bad item_count: %d", *info.ItemCount)
	}
}

func TestClient_GetDictionaryInfo_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{})
	if err != ErrMissingDictID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{
		DictID: "123",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{
		DictID:         "123",
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
