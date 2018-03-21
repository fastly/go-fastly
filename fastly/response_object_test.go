package fastly

import "testing"

func TestClient_ResponseObjects(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "response_objects/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var ro *ResponseObject
	record(t, "response_objects/create", func(c *Client) {
		ro, err = c.CreateResponseObject(&CreateResponseObjectInput{
			Service:     testServiceID,
			Version:     tv.Number,
			Name:        "test-response-object",
			Status:      200,
			Response:    "Ok",
			Content:     "abcd",
			ContentType: "text/plain",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "response_objects/cleanup", func(c *Client) {
			c.DeleteResponseObject(&DeleteResponseObjectInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-response-object",
			})

			c.DeleteResponseObject(&DeleteResponseObjectInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-response-object",
			})
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
	var ros []*ResponseObject
	record(t, "response_objects/list", func(c *Client) {
		ros, err = c.ListResponseObjects(&ListResponseObjectsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ros) < 1 {
		t.Errorf("bad response objects: %v", ros)
	}

	// Get
	var nro *ResponseObject
	record(t, "response_objects/get", func(c *Client) {
		nro, err = c.GetResponseObject(&GetResponseObjectInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-response-object",
		})
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
	var uro *ResponseObject
	record(t, "response_objects/update", func(c *Client) {
		uro, err = c.UpdateResponseObject(&UpdateResponseObjectInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-response-object",
			NewName: "new-test-response-object",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uro.Name != "new-test-response-object" {
		t.Errorf("bad name: %q", uro.Name)
	}

	// Delete
	record(t, "response_objects/delete", func(c *Client) {
		err = c.DeleteResponseObject(&DeleteResponseObjectInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-response-object",
		})
	})
	if err != nil {
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
		Version: 0,
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
		Version: 0,
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
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetResponseObject(&GetResponseObjectInput{
		Service: "foo",
		Version: 1,
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
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateResponseObject(&UpdateResponseObjectInput{
		Service: "foo",
		Version: 1,
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
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteResponseObject(&DeleteResponseObjectInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
