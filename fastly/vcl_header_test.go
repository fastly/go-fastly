package fastly

import (
	"errors"
	"testing"
)

func TestClient_Headers(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "headers/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var h *Header
	Record(t, "headers/create", func(c *Client) {
		h, err = c.CreateHeader(&CreateHeaderInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-header"),
			Action:         ToPointer(HeaderActionSet),
			IgnoreIfSet:    ToPointer(Compatibool(false)),
			Type:           ToPointer(HeaderTypeRequest),
			Destination:    ToPointer("http.foo"),
			Source:         ToPointer("client.ip"),
			Regex:          ToPointer("foobar"),
			Substitution:   ToPointer("123"),
			Priority:       ToPointer(50),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "headers/cleanup", func(c *Client) {
			_ = c.DeleteHeader(&DeleteHeaderInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-header",
			})

			_ = c.DeleteHeader(&DeleteHeaderInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-header",
			})
		})
	}()

	if *h.Name != "test-header" {
		t.Errorf("bad name: %q", *h.Name)
	}
	if *h.Action != HeaderActionSet {
		t.Errorf("bad header_action_set: %q", *h.Action)
	}
	if *h.IgnoreIfSet {
		t.Errorf("bad ignore_if_set: %t", *h.IgnoreIfSet)
	}
	if *h.Type != HeaderTypeRequest {
		t.Errorf("bad type: %q", *h.Type)
	}
	if *h.Destination != "http.foo" {
		t.Errorf("bad destination: %q", *h.Destination)
	}
	if *h.Source != "client.ip" {
		t.Errorf("bad source: %q", *h.Source)
	}
	if *h.Regex != "foobar" {
		t.Errorf("bad regex: %q", *h.Regex)
	}
	if *h.Substitution != "123" {
		t.Errorf("bad substitution: %q", *h.Substitution)
	}
	if *h.Priority != 50 {
		t.Errorf("bad priority: %d", *h.Priority)
	}

	// List
	var hs []*Header
	Record(t, "headers/list", func(c *Client) {
		hs, err = c.ListHeaders(&ListHeadersInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hs) < 1 {
		t.Errorf("bad headers: %v", hs)
	}

	// Get
	var nh *Header
	Record(t, "headers/get", func(c *Client) {
		nh, err = c.GetHeader(&GetHeaderInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-header",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *h.Name != *nh.Name {
		t.Errorf("bad name: %q (%q)", *h.Name, *nh.Name)
	}
	if *h.Action != *nh.Action {
		t.Errorf("bad header_action_set: %q", *h.Action)
	}
	if *h.IgnoreIfSet != *nh.IgnoreIfSet {
		t.Errorf("bad ignore_if_set: %t", *h.IgnoreIfSet)
	}
	if *h.Type != *nh.Type {
		t.Errorf("bad type: %q", *h.Type)
	}
	if *h.Destination != *nh.Destination {
		t.Errorf("bad destination: %q", *h.Destination)
	}
	if *h.Source != *nh.Source {
		t.Errorf("bad source: %q", *h.Source)
	}
	if *h.Regex != *nh.Regex {
		t.Errorf("bad regex: %q", *h.Regex)
	}
	if *h.Substitution != *nh.Substitution {
		t.Errorf("bad substitution: %q", *h.Substitution)
	}
	if *h.Priority != *nh.Priority {
		t.Errorf("bad priority: %d", *h.Priority)
	}

	// Update
	var uh *Header
	Record(t, "headers/update", func(c *Client) {
		uh, err = c.UpdateHeader(&UpdateHeaderInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-header",
			NewName:        ToPointer("new-test-header"),
			Action:         ToPointer(HeaderActionAppend),
			Type:           ToPointer(HeaderTypeFetch),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uh.Name != "new-test-header" {
		t.Errorf("bad name: %q", *uh.Name)
	}

	// Delete
	Record(t, "headers/delete", func(c *Client) {
		err = c.DeleteHeader(&DeleteHeaderInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-header",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListHeaders_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListHeaders(&ListHeadersInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListHeaders(&ListHeadersInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHeader_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateHeader(&CreateHeaderInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateHeader(&CreateHeaderInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHeader_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetHeader(&GetHeaderInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHeader(&GetHeaderInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetHeader(&GetHeaderInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHeader_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateHeader(&UpdateHeaderInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHeader(&UpdateHeaderInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateHeader(&UpdateHeaderInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHeader_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteHeader(&DeleteHeaderInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHeader(&DeleteHeaderInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteHeader(&DeleteHeaderInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
