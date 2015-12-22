package fastly

import "testing"

func TestClient_ResponseObjects(t *testing.T) {
	t.Parallel()

	tv := testVersion(t)

	// Create
	ro, err := testClient.CreateResponseObject(&CreateResponseObjectInput{
		Service:     testServiceID,
		Version:     tv.Number,
		Name:        "test-response-object",
		Status:      200,
		Response:    "Ok",
		Content:     "abcd",
		ContentType: "text/plain",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteResponseObject(&DeleteResponseObjectInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-response-object",
		})

		testClient.DeleteResponseObject(&DeleteResponseObjectInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-response-object",
		})
	}()

	if ro.Name != "test-response-object" {
		t.Errorf("bad name: %q", ro.Name)
	}
	if ro.Status != 200 {
		t.Errorf("bad status: %q", ro.Status)
	}
	if ro.Response != "Ok" {
		t.Errorf("bad response: %q", ro.Response)
	}
	if ro.Content != "abcd" {
		t.Errorf("bad content: %q", ro.Content)
	}
	if ro.ContentType != "text/plain" {
		t.Errorf("bad content_type: %q", ro.ContentType)
	}

	// List
	bs, err := testClient.ListResponseObjects(&ListResponseObjectsInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) < 1 {
		t.Errorf("bad response objects: %v", bs)
	}

	// Get
	nro, err := testClient.GetResponseObject(&GetResponseObjectInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-response-object",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ro.Name != nro.Name {
		t.Errorf("bad name: %q", ro.Name)
	}
	if ro.Status != nro.Status {
		t.Errorf("bad status: %q", ro.Status)
	}
	if ro.Response != nro.Response {
		t.Errorf("bad response: %q", ro.Response)
	}
	if ro.Content != nro.Content {
		t.Errorf("bad content: %q", ro.Content)
	}
	if ro.ContentType != nro.ContentType {
		t.Errorf("bad content_type: %q", ro.ContentType)
	}

	// Update
	uro, err := testClient.UpdateResponseObject(&UpdateResponseObjectInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-response-object",
		NewName: "new-test-response-object",
	})
	if err != nil {
		t.Fatal(err)
	}
	if uro.Name != "new-test-response-object" {
		t.Errorf("bad name: %q", uro.Name)
	}

	// Delete
	if err := testClient.DeleteResponseObject(&DeleteResponseObjectInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "new-test-response-object",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListResponseObjects_validation(t *testing.T) {
	var err error
	_, err = testClient.ListResponseObjects(&ListResponseObjectsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListResponseObjects(&ListResponseObjectsInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateResponseObject_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateResponseObject(&CreateResponseObjectInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateResponseObject(&CreateResponseObjectInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetResponseObject_validation(t *testing.T) {
	var err error
	_, err = testClient.GetResponseObject(&GetResponseObjectInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetResponseObject(&GetResponseObjectInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetResponseObject(&GetResponseObjectInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateResponseObject_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateResponseObject(&UpdateResponseObjectInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateResponseObject(&UpdateResponseObjectInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateResponseObject(&UpdateResponseObjectInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteResponseObject_validation(t *testing.T) {
	var err error
	err = testClient.DeleteResponseObject(&DeleteResponseObjectInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteResponseObject(&DeleteResponseObjectInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteResponseObject(&DeleteResponseObjectInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
