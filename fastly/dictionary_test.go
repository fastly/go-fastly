package fastly

import "testing"

func TestClient_Dictionaries(t *testing.T) {
	t.Parallel()

	fixtureBase := "dictionaries/"

	testVersion := createTestVersion(t, fixtureBase+"version", testServiceID)

	// Create
	var err error
	var d *Dictionary
	record(t, fixtureBase+"create", func(c *Client) {
		d, err = c.CreateDictionary(&CreateDictionaryInput{
			Service: testServiceID,
			Version: testVersion.Number,
			Name:    "test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteDictionary(&DeleteDictionaryInput{
				Service: testServiceID,
				Version: testVersion.Number,
				Name:    "test_dictionary",
			})

			c.DeleteDictionary(&DeleteDictionaryInput{
				Service: testServiceID,
				Version: testVersion.Number,
				Name:    "new_test_dictionary",
			})
		})
	}()

	if d.Name != "test_dictionary" {
		t.Errorf("bad name: %q", d.Name)
	}

	// List
	var ds []*Dictionary
	record(t, fixtureBase+"list", func(c *Client) {
		ds, err = c.ListDictionaries(&ListDictionariesInput{
			Service: testServiceID,
			Version: testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) < 1 {
		t.Errorf("bad dictionaries: %v", ds)
	}

	// Get
	var nd *Dictionary
	record(t, fixtureBase+"get", func(c *Client) {
		nd, err = c.GetDictionary(&GetDictionaryInput{
			Service: testServiceID,
			Version: testVersion.Number,
			Name:    "test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != nd.Name {
		t.Errorf("bad name: %q (%q)", d.Name, nd.Name)
	}

	// Update
	var ud *Dictionary
	record(t, fixtureBase+"update", func(c *Client) {
		ud, err = c.UpdateDictionary(&UpdateDictionaryInput{
			Service: testServiceID,
			Version: testVersion.Number,
			Name:    "test_dictionary",
			NewName: "new_test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.Name != "new_test_dictionary" {
		t.Errorf("bad name: %q", ud.Name)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteDictionary(&DeleteDictionaryInput{
			Service: testServiceID,
			Version: testVersion.Number,
			Name:    "new_test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDictionaries_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDictionaries(&ListDictionariesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDictionaries(&ListDictionariesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDictionary_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDictionary(&CreateDictionaryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDictionary(&CreateDictionaryInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDictionary_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDictionary(&GetDictionaryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionary(&GetDictionaryInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionary(&GetDictionaryInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDictionary_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDictionary(&UpdateDictionaryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDictionary(&UpdateDictionaryInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDictionary(&UpdateDictionaryInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDictionary_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDictionary(&DeleteDictionaryInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDictionary(&DeleteDictionaryInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDictionary(&DeleteDictionaryInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
