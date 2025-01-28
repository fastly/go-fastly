package fastly

import (
	"errors"
	"testing"
)

func TestClient_Dictionaries(t *testing.T) {
	t.Parallel()

	fixtureBase := "dictionaries/"

	testVersion := CreateTestVersion(t, fixtureBase+"version", TestDeliveryServiceID)

	// Create
	var err error
	var d *Dictionary
	Record(t, fixtureBase+"create", func(c *Client) {
		d, err = c.CreateDictionary(&CreateDictionaryInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Name:           ToPointer("test_dictionary"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteDictionary(&DeleteDictionaryInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				Name:           "test_dictionary",
			})

			_ = c.DeleteDictionary(&DeleteDictionaryInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *testVersion.Number,
				Name:           "new_test_dictionary",
			})
		})
	}()

	if *d.Name != "test_dictionary" {
		t.Errorf("bad name: %q", *d.Name)
	}

	// List
	var ds []*Dictionary
	Record(t, fixtureBase+"list", func(c *Client) {
		ds, err = c.ListDictionaries(&ListDictionariesInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
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
	Record(t, fixtureBase+"get", func(c *Client) {
		nd, err = c.GetDictionary(&GetDictionaryInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Name:           "test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *d.Name != *nd.Name {
		t.Errorf("bad name: %q (%q)", *d.Name, *nd.Name)
	}

	// Update
	var ud *Dictionary
	Record(t, fixtureBase+"update", func(c *Client) {
		ud, err = c.UpdateDictionary(&UpdateDictionaryInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Name:           "test_dictionary",
			NewName:        ToPointer("new_test_dictionary"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ud.Name != "new_test_dictionary" {
		t.Errorf("bad name: %q", *ud.Name)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteDictionary(&DeleteDictionaryInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *testVersion.Number,
			Name:           "new_test_dictionary",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDictionaries_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListDictionaries(&ListDictionariesInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListDictionaries(&ListDictionariesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDictionary_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateDictionary(&CreateDictionaryInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateDictionary(&CreateDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDictionary_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetDictionary(&GetDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDictionary(&GetDictionaryInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDictionary(&GetDictionaryInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDictionary_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateDictionary(&UpdateDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDictionary(&UpdateDictionaryInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDictionary(&UpdateDictionaryInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDictionary_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteDictionary(&DeleteDictionaryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDictionary(&DeleteDictionaryInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDictionary(&DeleteDictionaryInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
