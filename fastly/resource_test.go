package fastly

import (
	"testing"
)

func TestClient_Resources(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "resources/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create object-store resource we want to link to via Resource API.
	var o *ObjectStore
	record(t, "resources/create-object-store", func(c *Client) {
		o, err = c.CreateObjectStore(&CreateObjectStoreInput{
			Name: "test-object-store",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure object-store resource is deleted
	defer func() {
		record(t, "resources/cleanup-object-store", func(c *Client) {
			_ = c.DeleteObjectStore(&DeleteObjectStoreInput{
				ID: o.ID,
			})
		})
	}()

	// Create
	var r *Resource
	record(t, "resources/create", func(c *Client) {
		r, err = c.CreateResource(&CreateResourceInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           String("test-resource"),
			ResourceID:     String(o.ID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "resources/cleanup", func(c *Client) {
			// NOTE: The API documentation is confusing here because they name the
			// parameter `resource_id` but they actually mean (as far as their data model
			// is concerned) the `id` field. `resource_id`, from the API perspective, is
			// referring to the resource you're creating a link to (e.g. an object store).
			_ = c.DeleteResource(&DeleteResourceInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				ResourceID:     r.ID,
			})
		})
	}()

	if r.Name != "test-resource" {
		t.Errorf("bad name: %q", r.Name)
	}

	// List
	var rs []*Resource
	record(t, "resources/list", func(c *Client) {
		rs, err = c.ListResources(&ListResourcesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) < 1 {
		t.Errorf("bad resources: %v", rs)
	}

	// Get
	var gr *Resource
	record(t, "resources/get", func(c *Client) {
		// NOTE: The API documentation is confusing here because they name the
		// parameter `resource_id` but they actually mean (as far as their data model
		// is concerned) the `id` field. `resource_id`, from the API perspective, is
		// referring to the resource you're creating a link to (e.g. an object store).
		gr, err = c.GetResource(&GetResourceInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			ResourceID:     r.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Name != gr.Name {
		t.Errorf("bad name: %q (%q)", r.Name, gr.Name)
	}

	// Update
	var ur *Resource
	record(t, "resources/update", func(c *Client) {
		// NOTE: The API documentation is confusing here because they name the
		// parameter `resource_id` but they actually mean (as far as their data model
		// is concerned) the `id` field. `resource_id`, from the API perspective, is
		// referring to the resource you're creating a link to (e.g. an object store).
		ur, err = c.UpdateResource(&UpdateResourceInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			ResourceID:     r.ID,
			Name:           String("new-test-resource"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ur.Name != "new-test-resource" {
		t.Errorf("bad name: %q", ur.Name)
	}

	// Delete
	record(t, "resources/delete", func(c *Client) {
		// NOTE: The API documentation is confusing here because they name the
		// parameter `resource_id` but they actually mean (as far as their data model
		// is concerned) the `id` field. `resource_id`, from the API perspective, is
		// referring to the resource you're creating a link to (e.g. an object store).
		err = c.DeleteResource(&DeleteResourceInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			ResourceID:     ur.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListResources_validation(t *testing.T) {
	var err error
	_, err = testClient.ListResources(&ListResourcesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListResources(&ListResourcesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateResource_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateResource(&CreateResourceInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateResource(&CreateResourceInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetResource_validation(t *testing.T) {
	var err error

	_, err = testClient.GetResource(&GetResourceInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingResourceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetResource(&GetResourceInput{
		ResourceID:     "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetResource(&GetResourceInput{
		ResourceID: "test",
		ServiceID:  "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateResource_validation(t *testing.T) {
	var err error

	_, err = testClient.UpdateResource(&UpdateResourceInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateResource(&UpdateResourceInput{
		ResourceID:     "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateResource(&UpdateResourceInput{
		ResourceID: "test",
		ServiceID:  "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteResource_validation(t *testing.T) {
	var err error

	err = testClient.DeleteResource(&DeleteResourceInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingResourceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteResource(&DeleteResourceInput{
		ResourceID:     "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteResource(&DeleteResourceInput{
		ResourceID: "test",
		ServiceID:  "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
