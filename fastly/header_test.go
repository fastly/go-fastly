package fastly

import "testing"

func TestClient_Headers(t *testing.T) {
	tv := testVersion(t)

	// Create
	h, err := testClient.CreateHeader(&CreateHeaderInput{
		Service:      testServiceID,
		Version:      tv.Number,
		Name:         "test-header",
		Action:       HeaderActionSet,
		IgnoreIfSet:  false,
		Type:         HeaderTypeRequest,
		Destination:  "http.foo",
		Source:       "client.ip",
		Regex:        "foobar",
		Substitution: "123",
		Priority:     50,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteHeader(&DeleteHeaderInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-header",
		})

		testClient.DeleteHeader(&DeleteHeaderInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-header",
		})
	}()

	if h.Name != "test-header" {
		t.Errorf("bad name: %q", h.Name)
	}
	if h.Action != HeaderActionSet {
		t.Errorf("bad header_action_set: %q", h.Action)
	}
	if h.IgnoreIfSet != false {
		t.Errorf("bad ignore_if_set: %b", h.IgnoreIfSet)
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
	hs, err := testClient.ListHeaders(&ListHeadersInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hs) < 1 {
		t.Errorf("bad headers: %v", hs)
	}

	// Get
	nh, err := testClient.GetHeader(&GetHeaderInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-header",
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
		t.Errorf("bad ignore_if_set: %b", h.IgnoreIfSet)
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
	uh, err := testClient.UpdateHeader(&UpdateHeaderInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-header",
		NewName: "new-test-header",
	})
	if err != nil {
		t.Fatal(err)
	}
	if uh.Name != "new-test-header" {
		t.Errorf("bad name: %q", uh.Name)
	}

	// Delete
	if err := testClient.DeleteHeader(&DeleteHeaderInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "new-test-header",
	}); err != nil {
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
		Version: "",
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
		Version: "",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHeader(&GetHeaderInput{
		Service: "foo",
		Version: "1",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHeader(&UpdateHeaderInput{
		Service: "foo",
		Version: "1",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHeader(&DeleteHeaderInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
