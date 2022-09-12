package fastly

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestClient_ObjectStore(t *testing.T) {
	t.Parallel()

	const createStoreName = "object-store-test-store"

	// List
	var objectStoreListResp1 *ListObjectStoresResponse
	var err error
	record(t, "object_store/list-store", func(c *Client) {
		objectStoreListResp1, err = c.ListObjectStores(nil)
	})
	if err != nil {
		t.Fatal(err)
	}

	// make sure our test store we're going to create isn't among them
	for _, store := range objectStoreListResp1.Data {
		if store.Name == createStoreName {
			t.Errorf("Found test store %q, aborting", createStoreName)
		}
	}

	// Create
	var objectStore *ObjectStore
	input := &CreateObjectStoreInput{
		Name: createStoreName,
	}
	record(t, "object_store/create-store", func(c *Client) {
		objectStore, err = c.CreateObjectStore(input)
	})
	if err != nil {
		t.Fatal(err)
	}

	if objectStore.Name != createStoreName {
		t.Errorf("CreateObjectStores: unexpected store name: got %q, want %q", objectStore.Name, createStoreName)
	}

	// ensure we delete it
	defer func() {
		record(t, "object_store/cleanup", func(c *Client) {
			// first delete all the keys in it
			p := c.NewListObjectStoreKeysPaginator(&ListObjectStoreKeysInput{ID: objectStore.ID})
			for p.Next() {
				keys := p.Keys()
				sort.Strings(keys)
				for _, key := range keys {
					err = c.DeleteObjectStoreKey(&DeleteObjectStoreKeyInput{ID: objectStore.ID, Key: key})
					if err != nil {
						t.Errorf("error during key cleanup: %v", err)
					}
				}
			}
			if err := p.Err(); err != nil {
				t.Errorf("error during cleanup pagination: %v", err)
			}

			err = c.DeleteObjectStore(&DeleteObjectStoreInput{ID: objectStore.ID})
			if err != nil {
				t.Errorf("error during store cleanup: %v", err)
			}
		})
	}()

	// fetch the newly created store and verify it matches
	var getObjectStoreResponse *ObjectStore
	record(t, "object_store/get-store", func(c *Client) {
		getObjectStoreResponse, err = c.GetObjectStore(&GetObjectStoreInput{ID: objectStore.ID})
	})
	if err != nil {
		t.Fatal(err)
	}

	if getObjectStoreResponse.Name != objectStore.Name || getObjectStoreResponse.ID != objectStore.ID {
		t.Errorf("error fetching info for store %q: got %q, want %q", createStoreName, getObjectStoreResponse, objectStore)

	}

	// create a bunch of keys in our object store
	keys := []string{"apple", "banana", "carrot", "dragonfruit", "eggplant"}

	record(t, "object_store/create-keys", func(c *Client) {
		for i, key := range keys {
			err := c.InsertObjectStoreKey(&InsertObjectStoreKeyInput{ID: objectStore.ID, Key: key, Value: key + strconv.Itoa(i)})
			if err != nil {
				t.Errorf("error inserting key %q: %v", key, err)
			}
		}
	})

	record(t, "object_store/check-keys", func(c *Client) {
		for i, key := range keys {
			got, err := c.GetObjectStoreKey(&GetObjectStoreKeyInput{ID: objectStore.ID, Key: key})
			if err != nil {
				t.Errorf("error fetching key %q: %v", key, err)
			}
			want := key + strconv.Itoa(i)
			if got != want {
				t.Errorf("mismatch fetching key %q: got %q, want %q", key, got, want)
			}
		}
	})

	// fetch the keys
	var objectStoreListKeys *ListObjectStoreKeysResponse
	record(t, "object_store/list-keys", func(c *Client) {
		objectStoreListKeys, err = c.ListObjectStoreKeys(&ListObjectStoreKeysInput{ID: objectStore.ID})
	})

	if err != nil {
		t.Errorf("error listing keys: %v", err)
	}

	sort.Strings(objectStoreListKeys.Data)
	if !reflect.DeepEqual(keys, objectStoreListKeys.Data) {
		t.Errorf("mismatch listing keys: got %q, want %q", objectStoreListKeys.Data, keys)
	}

	record(t, "object_store/list-keys-pagination", func(c *Client) {
		p := c.NewListObjectStoreKeysPaginator(&ListObjectStoreKeysInput{ID: objectStore.ID, Limit: 4})
		var listed []string
		expected := []int{4, 1}
		var page int
		for p.Next() {
			keys := p.Keys()
			if len(keys) != expected[page] {
				t.Errorf("wrong number of keys for page %d: got %d, want %d", page, len(keys), expected[page])
			}
			page++
			listed = append(listed, keys...)
		}
		if err := p.Err(); err != nil {
			t.Errorf("error during keys pagination: %v", err)
		}
		sort.Strings(listed)
		if !reflect.DeepEqual(keys, listed) {
			t.Errorf("mismatch listing paginated keys: got %q, want %q", objectStoreListKeys.Data, keys)
		}
	})

	// create more stores
	// list
	// list with pagination
	// cleanup

}
