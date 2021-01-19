package fastly

import (
	"testing"
)

func TestClient_Dictionaries(t *testing.T) {
	t.Parallel()

	fixtureBase := "dictionaries/"

	testVersion := createTestVersion(t, fixtureBase+"version", testServiceID)

	// Create
	var err error
	var d *Dictionary
	record(t, fixtureBase+"create", func(c *Client) {
		d, err = c.CreateDictionary(&CreateDictionaryInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           "test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteDictionary(&DeleteDictionaryInput{
				ServiceID:      testServiceID,
				ServiceVersion: testVersion.Number,
				Name:           "test_dictionary",
			})

			c.DeleteDictionary(&DeleteDictionaryInput{
				ServiceID:      testServiceID,
				ServiceVersion: testVersion.Number,
				Name:           "new_test_dictionary",
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
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
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
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           "test_dictionary",
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
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           "test_dictionary",
			NewName:        String("new_test_dictionary"),
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
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           "new_test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDictionaries_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDictionaries(&ListDictionariesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDictionaries(&ListDictionariesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDictionary_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDictionary(&CreateDictionaryInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDictionary(&CreateDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDictionary_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDictionary(&GetDictionaryInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionary(&GetDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionary(&GetDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDictionary_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDictionary(&UpdateDictionaryInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDictionary(&UpdateDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDictionary(&UpdateDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDictionary_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDictionary(&DeleteDictionaryInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDictionary(&DeleteDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDictionary(&DeleteDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
