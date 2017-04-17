package fastly

import "testing"

func TestClient_Headers(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "headers/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var h *Header
	record(t, "headers/create", func(c *Client) {
		h, err = c.CreateHeader(&CreateHeaderInput{
			Service:      testServiceID,
			Version:      tv.Number,
			Name:         "test-header",
			Action:       HeaderActionSet,
			IgnoreIfSet:  CBool(false),
			Type:         HeaderTypeRequest,
			Destination:  "http.foo",
			Source:       "client.ip",
			Regex:        "foobar",
			Substitution: "123",
			Priority:     50,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "headers/cleanup", func(c *Client) {
			c.DeleteHeader(&DeleteHeaderInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-header",
			})

			c.DeleteHeader(&DeleteHeaderInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-header",
			})
		})
	}()

	if h.Name != "test-header" {
		t.Errorf("bad name: %q", h.Name)
	}
	if h.Action != HeaderActionSet {
		t.Errorf("bad header_action_set: %q", h.Action)
	}
	if h.IgnoreIfSet != false {
		t.Errorf("bad ignore_if_set: %t", h.IgnoreIfSet)
	}
	if h.Type != HeaderTypeRequest {
		t.Errorf("bad type: %q", h.Type)
	}
	if h.Destination != "http.foo" {
		t.Errorf("bad destination: %q", h.Destination)
	}
	if h.Source != "client.ip" {
		t.Errorf("bad source: %q", h.Source)
	}
	if h.Regex != "foobar" {
		t.Errorf("bad regex: %q", h.Regex)
	}
	if h.Substitution != "123" {
		t.Errorf("bad substitution: %q", h.Substitution)
	}
	if h.Priority != 50 {
		t.Errorf("bad priority: %d", h.Priority)
	}

	// List
	var hs []*Header
	record(t, "headers/list", func(c *Client) {
		hs, err = c.ListHeaders(&ListHeadersInput{
			Service: testServiceID,
			Version: tv.Number,
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
	record(t, "headers/get", func(c *Client) {
		nh, err = c.GetHeader(&GetHeaderInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-header",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if h.Name != nh.Name {
		t.Errorf("bad name: %q (%q)", h.Name, nh.Name)
	}
	if h.Action != nh.Action {
		t.Errorf("bad header_action_set: %q", h.Action)
	}
	if h.IgnoreIfSet != nh.IgnoreIfSet {
		t.Errorf("bad ignore_if_set: %t", h.IgnoreIfSet)
	}
	if h.Type != nh.Type {
		t.Errorf("bad type: %q", h.Type)
	}
	if h.Destination != nh.Destination {
		t.Errorf("bad destination: %q", h.Destination)
	}
	if h.Source != nh.Source {
		t.Errorf("bad source: %q", h.Source)
	}
	if h.Regex != nh.Regex {
		t.Errorf("bad regex: %q", h.Regex)
	}
	if h.Substitution != nh.Substitution {
		t.Errorf("bad substitution: %q", h.Substitution)
	}
	if h.Priority != nh.Priority {
		t.Errorf("bad priority: %d", h.Priority)
	}

	// Update
	var uh *Header
	record(t, "headers/update", func(c *Client) {
		uh, err = c.UpdateHeader(&UpdateHeaderInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-header",
			NewName: "new-test-header",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uh.Name != "new-test-header" {
		t.Errorf("bad name: %q", uh.Name)
	}

	// Delete
	record(t, "headers/delete", func(c *Client) {
		err = c.DeleteHeader(&DeleteHeaderInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-header",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListHeaders_validation(t *testing.T) {
	var err error
	_, err = testClient.ListHeaders(&ListHeadersInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListHeaders(&ListHeadersInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHeader_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateHeader(&CreateHeaderInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateHeader(&CreateHeaderInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHeader_validation(t *testing.T) {
	var err error
	_, err = testClient.GetHeader(&GetHeaderInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHeader(&GetHeaderInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHeader(&GetHeaderInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHeader_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateHeader(&UpdateHeaderInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHeader(&UpdateHeaderInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHeader(&UpdateHeaderInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHeader_validation(t *testing.T) {
	var err error
	err = testClient.DeleteHeader(&DeleteHeaderInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHeader(&DeleteHeaderInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHeader(&DeleteHeaderInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
