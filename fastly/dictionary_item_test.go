package fastly

import "testing"

func createTestDictionary(t *testing.T) *Dictionary {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "dictionary_items/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	var d *Dictionary
	record(t, "dictionary_items/dictionary", func(c *Client) {
		d, err = c.CreateDictionary(&CreateDictionaryInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func deleteTestDictionary(d *Dictionary, t *testing.T) {
	var err error
	record(t, "dictionary_items/delete_dictionary", func(c *Client) {
		err = c.DeleteDictionary(&DeleteDictionaryInput{
			Service: d.ServiceID,
			Version: d.Version,
			Name:    "test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DictionaryItems(t *testing.T) {
	td := createTestDictionary(t)
	defer deleteTestDictionary(td, t)

	// Create
	var err error
	var d *DictionaryItem
	record(t, "dictionary_items/create", func(c *Client) {
		d, err = c.CreateDictionaryItem(&CreateDictionaryItemInput{
			Service:    testServiceID,
			Dictionary: td.ID,
			ItemKey:    "test-dictionary-item",
			ItemValue:  "value",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "dictionary_items/cleanup", func(c *Client) {
			c.DeleteDictionaryItem(&DeleteDictionaryItemInput{
				Service:    testServiceID,
				Dictionary: td.ID,
				ItemKey:    "test-dictionary-item",
			})
		})
	}()

	if d.ItemKey != "test-dictionary-item" {
		t.Errorf("bad item_key: %q", d.ItemKey)
	}
	if d.ItemValue != "value" {
		t.Errorf("bad item_value: %q", d.ItemValue)
	}

	// List
	var ds []*DictionaryItem
	record(t, "dictionary_items/list", func(c *Client) {
		ds, err = c.ListDictionaryItems(&ListDictionaryItemsInput{
			Service:    testServiceID,
			Dictionary: td.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) < 1 {
		t.Errorf("bad dictionary items: %v", ds)
	}

	// Get
	var nd *DictionaryItem
	record(t, "dictionary_items/get", func(c *Client) {
		nd, err = c.GetDictionaryItem(&GetDictionaryItemInput{
			Service:    testServiceID,
			Dictionary: td.ID,
			ItemKey:    "test-dictionary-item",
		})
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
	var ud *DictionaryItem
	record(t, "dictionary_items/update", func(c *Client) {
		ud, err = c.UpdateDictionaryItem(&UpdateDictionaryItemInput{
			Service:    testServiceID,
			Dictionary: td.ID,
			ItemKey:    "test-dictionary-item",
			ItemValue:  "new-value",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.ItemValue != "new-value" {
		t.Errorf("bad item_value: %q", ud.ItemValue)
	}

	// Delete
	record(t, "dictionary_items/delete", func(c *Client) {
		err = c.DeleteDictionaryItem(&DeleteDictionaryItemInput{
			Service:    testServiceID,
			Dictionary: td.ID,
			ItemKey:    "test-dictionary-item",
		})
	})
	if err != nil {
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
