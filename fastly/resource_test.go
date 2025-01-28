package fastly

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestClient_Resources(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "resources/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create kv-store resource we want to link to via Resource API.
	var o *KVStore
	Record(t, "resources/create-kv-store", func(c *Client) {
		o, err = c.CreateKVStore(&CreateKVStoreInput{
			Name: "test-kv-store",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure kv-store resource is deleted
	defer func() {
		Record(t, "resources/cleanup-kv-store", func(c *Client) {
			_ = c.DeleteKVStore(&DeleteKVStoreInput{
				StoreID: o.StoreID,
			})
		})
	}()

	// NOTE: This doesn't have to match the actual kv-store name.
	// This is an opportunity for you to use an 'alias' for your kv store.
	// So your service will now refer to the kv-store using this name.
	const kvStoreNameForServiceLinking = "test-kv-store-name-for-linking"

	// Create
	var r *Resource
	Record(t, "resources/create", func(c *Client) {
		r, err = c.CreateResource(&CreateResourceInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer(kvStoreNameForServiceLinking),
			ResourceID:     ToPointer(o.StoreID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "resources/cleanup", func(c *Client) {
			_ = c.DeleteResource(&DeleteResourceInput{
				ResourceID:     *r.LinkID,
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
			})
		})
	}()

	if *r.Name != kvStoreNameForServiceLinking {
		t.Errorf("bad name: %q", *r.Name)
	}

	// List
	var rs []*Resource
	Record(t, "resources/list", func(c *Client) {
		rs, err = c.ListResources(&ListResourcesInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "resources/get", func(c *Client) {
		gr, err = c.GetResource(&GetResourceInput{
			ResourceID:     *r.LinkID,
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *r.Name != *gr.Name {
		t.Errorf("bad name: %q (%q)", *r.Name, *gr.Name)
	}

	// Update
	var ur *Resource
	Record(t, "resources/update", func(c *Client) {
		ur, err = c.UpdateResource(&UpdateResourceInput{
			ResourceID:     *r.LinkID,
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("new-kv-store-alias-for-my-service"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ur.Name != "new-kv-store-alias-for-my-service" {
		t.Errorf("bad name: %q", *ur.Name)
	}

	// Delete
	Record(t, "resources/delete", func(c *Client) {
		err = c.DeleteResource(&DeleteResourceInput{
			ResourceID:     *ur.LinkID,
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListResources_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListResources(&ListResourcesInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListResources(&ListResourcesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateResource_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateResource(&CreateResourceInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateResource(&CreateResourceInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetResource_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetResource(&GetResourceInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingResourceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetResource(&GetResourceInput{
		ResourceID:     "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetResource(&GetResourceInput{
		ResourceID: "test",
		ServiceID:  "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateResource_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateResource(&UpdateResourceInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingResourceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateResource(&UpdateResourceInput{
		ResourceID:     "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateResource(&UpdateResourceInput{
		ResourceID: "test",
		ServiceID:  "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteResource_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteResource(&DeleteResourceInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingResourceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteResource(&DeleteResourceInput{
		ResourceID:     "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteResource(&DeleteResourceInput{
		ResourceID: "test",
		ServiceID:  "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestResourceJSONRoundtrip(t *testing.T) {
	now := time.Now()
	r := Resource{
		CreatedAt:      &now,
		DeletedAt:      &now,
		HREF:           ToPointer("the/href"),
		LinkID:         ToPointer("the-id"),
		Name:           ToPointer("the-name"),
		ResourceID:     ToPointer("the-resource-id"),
		ResourceType:   ToPointer("the-resource-type"),
		ServiceID:      ToPointer("the-service-id"),
		ServiceVersion: ToPointer(1),
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
	if err := DecodeBodyMap(&out, &decoded); err != nil {
		t.Fatal(err)
	}
	t.Logf("Decoded:\n%#v", decoded)

	if got, want := *decoded.HREF, *r.HREF; got != want {
		t.Errorf("HREF: got %q, want %q", got, want)
	}
	if got, want := *decoded.LinkID, *r.LinkID; got != want {
		t.Errorf("ID: got %q, want %q", got, want)
	}
	if got, want := *decoded.Name, *r.Name; got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
	if got, want := *decoded.ResourceID, *r.ResourceID; got != want {
		t.Errorf("ResourceID: got %q, want %q", got, want)
	}
	if got, want := *decoded.ResourceType, *r.ResourceType; got != want {
		t.Errorf("ResourceType: got %q, want %q", got, want)
	}
	if got, want := *decoded.ServiceID, *r.ServiceID; got != want {
		t.Errorf("ServiceID: got %q, want %q", got, want)
	}
	if got, want := *decoded.ServiceVersion, *r.ServiceVersion; got != want {
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
