package fastly

import (
	"errors"
	"net/http"
	"testing"
)

func TestClient_ResponseObjects(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "response_objects/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var ro *ResponseObject
	Record(t, "response_objects/create", func(c *Client) {
		ro, err = c.CreateResponseObject(&CreateResponseObjectInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-response-object"),
			Status:         ToPointer(http.StatusOK),
			Response:       ToPointer("Ok"),
			Content:        ToPointer("abcd"),
			ContentType:    ToPointer("text/plain"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "response_objects/cleanup", func(c *Client) {
			_ = c.DeleteResponseObject(&DeleteResponseObjectInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-response-object",
			})

			_ = c.DeleteResponseObject(&DeleteResponseObjectInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-response-object",
			})
		})
	}()

	if *ro.Name != "test-response-object" {
		t.Errorf("bad name: %q", *ro.Name)
	}
	if *ro.Status != http.StatusOK {
		t.Errorf("bad status: %q", *ro.Status)
	}
	if *ro.Response != "Ok" {
		t.Errorf("bad response: %q", *ro.Response)
	}
	if *ro.Content != "abcd" {
		t.Errorf("bad content: %q", *ro.Content)
	}
	if *ro.ContentType != "text/plain" {
		t.Errorf("bad content_type: %q", *ro.ContentType)
	}

	// List
	var ros []*ResponseObject
	Record(t, "response_objects/list", func(c *Client) {
		ros, err = c.ListResponseObjects(&ListResponseObjectsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "response_objects/get", func(c *Client) {
		nro, err = c.GetResponseObject(&GetResponseObjectInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-response-object",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ro.Name != *nro.Name {
		t.Errorf("bad name: %q", *ro.Name)
	}
	if *ro.Status != *nro.Status {
		t.Errorf("bad status: %q", *ro.Status)
	}
	if *ro.Response != *nro.Response {
		t.Errorf("bad response: %q", *ro.Response)
	}
	if *ro.Content != *nro.Content {
		t.Errorf("bad content: %q", *ro.Content)
	}
	if *ro.ContentType != *nro.ContentType {
		t.Errorf("bad content_type: %q", *ro.ContentType)
	}

	// Update
	var uro *ResponseObject
	Record(t, "response_objects/update", func(c *Client) {
		uro, err = c.UpdateResponseObject(&UpdateResponseObjectInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-response-object",
			NewName:        ToPointer("new-test-response-object"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uro.Name != "new-test-response-object" {
		t.Errorf("bad name: %q", *uro.Name)
	}

	// Delete
	Record(t, "response_objects/delete", func(c *Client) {
		err = c.DeleteResponseObject(&DeleteResponseObjectInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-response-object",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListResponseObjects_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListResponseObjects(&ListResponseObjectsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListResponseObjects(&ListResponseObjectsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateResponseObject_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateResponseObject(&CreateResponseObjectInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateResponseObject(&CreateResponseObjectInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetResponseObject_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetResponseObject(&GetResponseObjectInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetResponseObject(&GetResponseObjectInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetResponseObject(&GetResponseObjectInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateResponseObject_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateResponseObject(&UpdateResponseObjectInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateResponseObject(&UpdateResponseObjectInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateResponseObject(&UpdateResponseObjectInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteResponseObject_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteResponseObject(&DeleteResponseObjectInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteResponseObject(&DeleteResponseObjectInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteResponseObject(&DeleteResponseObjectInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
