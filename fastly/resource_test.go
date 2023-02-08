package fastly

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
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

	// NOTE: This doesn't have to match the actual object-store name.
	// This is an opportunity for you to use an 'alias' for your object store.
	// So your service will now refer to the object-store using this name.
	const objectStoreNameForServiceLinking = "test-object-store-name-for-linking"

	// Create
	var r *Resource
	record(t, "resources/create", func(c *Client) {
		r, err = c.CreateResource(&CreateResourceInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           String(objectStoreNameForServiceLinking),
			ResourceID:     String(o.ID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "resources/cleanup", func(c *Client) {
			_ = c.DeleteResource(&DeleteResourceInput{
				ID:             r.ID,
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
			})
		})
	}()

	if r.Name != objectStoreNameForServiceLinking {
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
		gr, err = c.GetResource(&GetResourceInput{
			ID:             r.ID,
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
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
		ur, err = c.UpdateResource(&UpdateResourceInput{
			ID:             r.ID,
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           String("new-object-store-alias-for-my-service"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ur.Name != "new-object-store-alias-for-my-service" {
		t.Errorf("bad name: %q", ur.Name)
	}

	// Delete
	record(t, "resources/delete", func(c *Client) {
		err = c.DeleteResource(&DeleteResourceInput{
			ID:             ur.ID,
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
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
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetResource(&GetResourceInput{
		ID:             "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetResource(&GetResourceInput{
		ID:        "test",
		ServiceID: "foo",
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
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateResource(&UpdateResourceInput{
		ID:             "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateResource(&UpdateResourceInput{
		ID:        "test",
		ServiceID: "foo",
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
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteResource(&DeleteResourceInput{
		ID:             "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteResource(&DeleteResourceInput{
		ID:        "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestResourceJSONRoundtrip(t *testing.T) {
	now := time.Now()
	r := Resource{
		CreatedAt:      &now,
		DeletedAt:      &now,
		HREF:           "the/href",
		ID:             "the-id",
		Name:           "the-name",
		ResourceID:     "the-resource-id",
		ResourceType:   "the-resource-type",
		ServiceID:      "the-service-id",
		ServiceVersion: "the-service-version",
		UpdatedAt:      &now,
	}

	// Ensure that decode(encode(resource)) == resource.

	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		t.Fatal(err)
	}

	encoded := out.String()
	t.Logf("Encoded:\n%s", encoded)

	var decoded Resource
	if err := decodeBodyMap(&out, &decoded); err != nil {
		t.Fatal(err)
	}
	t.Logf("Decoded:\n%#v", decoded)

	if got, want := decoded.HREF, r.HREF; got != want {
		t.Errorf("HREF: got %q, want %q", got, want)
	}
	if got, want := decoded.ID, r.ID; got != want {
		t.Errorf("ID: got %q, want %q", got, want)
	}
	if got, want := decoded.Name, r.Name; got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
	if got, want := decoded.ResourceID, r.ResourceID; got != want {
		t.Errorf("ResourceID: got %q, want %q", got, want)
	}
	if got, want := decoded.ResourceType, r.ResourceType; got != want {
		t.Errorf("ResourceType: got %q, want %q", got, want)
	}
	if got, want := decoded.ServiceID, r.ServiceID; got != want {
		t.Errorf("ServiceID: got %q, want %q", got, want)
	}
	if got, want := decoded.ServiceVersion, r.ServiceVersion; got != want {
		t.Errorf("ServiceVersion: got %q, want %q", got, want)
	}

	if got, want := decoded.CreatedAt, r.CreatedAt; got == nil || !got.Equal(*want) {
		t.Errorf("CreatedAt: got %s, want %s", got, want)
	}
	if got, want := decoded.DeletedAt, r.DeletedAt; got == nil || !got.Equal(*want) {
		t.Errorf("DeletedAt: got %s, want %s", got, want)
	}
	if got, want := decoded.UpdatedAt, r.UpdatedAt; got == nil || !got.Equal(*want) {
		t.Errorf("UpdatedAt: got %s, want %s", got, want)
	}
}
