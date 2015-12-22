package fastly

import "testing"

func createTestDictionary(t *testing.T) *Dictionary {
	t.Parallel()

	tv := testVersion(t)

	d, err := testClient.CreateDictionary(&CreateDictionaryInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test_dictionary",
	})
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func deleteTestDictionary(d *Dictionary, t *testing.T) {
	if err := testClient.DeleteDictionary(&DeleteDictionaryInput{
		Service: d.ServiceID,
		Version: d.Version,
		Name:    "test_dictionary",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_DictionaryItems(t *testing.T) {
	td := createTestDictionary(t)
	defer deleteTestDictionary(td, t)

	// Create
	d, err := testClient.CreateDictionaryItem(&CreateDictionaryItemInput{
		Service:    testServiceID,
		Dictionary: td.ID,
		ItemKey:    "test-dictionary-item",
		ItemValue:  "value",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
			Service:    testServiceID,
			Dictionary: td.ID,
			ItemKey:    "test-dictionary-item",
		})
	}()

	if d.ItemKey != "test-dictionary-item" {
		t.Errorf("bad item_key: %q", d.ItemKey)
	}
	if d.ItemValue != "value" {
		t.Errorf("bad item_value: %q", d.ItemValue)
	}

	// List
	ds, err := testClient.ListDictionaryItems(&ListDictionaryItemsInput{
		Service:    testServiceID,
		Dictionary: td.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) < 1 {
		t.Errorf("bad dictionary items: %v", ds)
	}

	// Get
	nd, err := testClient.GetDictionaryItem(&GetDictionaryItemInput{
		Service:    testServiceID,
		Dictionary: td.ID,
		ItemKey:    "test-dictionary-item",
	})
	if err != nil {
		t.Fatal(err)
	}
	if nd.ItemKey != "test-dictionary-item" {
		t.Errorf("bad item_key: %q", nd.ItemKey)
	}
	if nd.ItemValue != "value" {
		t.Errorf("bad item_value: %q", nd.ItemValue)
	}

	// Update
	ud, err := testClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		Service:    testServiceID,
		Dictionary: td.ID,
		ItemKey:    "test-dictionary-item",
		ItemValue:  "new-value",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.ItemValue != "new-value" {
		t.Errorf("bad item_value: %q", ud.ItemValue)
	}

	// Delete
	if err := testClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		Service:    testServiceID,
		Dictionary: td.ID,
		ItemKey:    "test-dictionary-item",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDictionaryItems_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDictionaryItems(&ListDictionaryItemsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDictionaryItems(&ListDictionaryItemsInput{
		Service:    "foo",
		Dictionary: "",
	})
	if err != ErrMissingDictionary {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDictionaryItem_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDictionaryItem(&CreateDictionaryItemInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDictionaryItem(&CreateDictionaryItemInput{
		Service:    "foo",
		Dictionary: "",
	})
	if err != ErrMissingDictionary {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDictionaryItem_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDictionaryItem(&GetDictionaryItemInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryItem(&GetDictionaryItemInput{
		Service:    "foo",
		Dictionary: "",
	})
	if err != ErrMissingDictionary {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryItem(&GetDictionaryItemInput{
		Service:    "foo",
		Dictionary: "test",
		ItemKey:    "",
	})
	if err != ErrMissingItemKey {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDictionaryItem_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		Service:    "foo",
		Dictionary: "",
	})
	if err != ErrMissingDictionary {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDictionaryItem(&UpdateDictionaryItemInput{
		Service:    "foo",
		Dictionary: "test",
		ItemKey:    "",
	})
	if err != ErrMissingItemKey {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDictionaryItem_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		Service:    "foo",
		Dictionary: "",
	})
	if err != ErrMissingDictionary {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDictionaryItem(&DeleteDictionaryItemInput{
		Service:    "foo",
		Dictionary: "test",
		ItemKey:    "",
	})
	if err != ErrMissingItemKey {
		t.Errorf("bad error: %s", err)
	}
}
