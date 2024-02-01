package fastly

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestClient_KVStore(t *testing.T) {
	t.Parallel()

	const createStoreName = "kv-store-test-store"

	// List
	var kvStoreListResp1 *ListKVStoresResponse
	var err error
	record(t, "kv_store/list-store", func(c *Client) {
		kvStoreListResp1, err = c.ListKVStores(nil)
	})
	if err != nil {
		t.Fatal(err)
	}

	// make sure our test store we're going to create isn't among them
	for _, store := range kvStoreListResp1.Data {
		if store.Name == createStoreName {
			t.Errorf("Found test store %q, aborting", createStoreName)
		}
	}

	// Create
	var kvStore *KVStore
	input := &CreateKVStoreInput{
		Name: createStoreName,
	}
	record(t, "kv_store/create-store", func(c *Client) {
		kvStore, err = c.CreateKVStore(input)
	})
	if err != nil {
		t.Fatal(err)
	}

	if kvStore.Name != createStoreName {
		t.Errorf("CreateKVStores: unexpected store name: got %q, want %q", kvStore.Name, createStoreName)
	}

	// ensure we delete it
	defer func() {
		record(t, "kv_store/cleanup", func(c *Client) {
			// first delete all the keys in it
			p := c.NewListKVStoreKeysPaginator(&ListKVStoreKeysInput{
				Consistency: ConsistencyEventual,
				StoreID:     kvStore.StoreID,
			})
			for p.Next() {
				keys := p.Keys()
				sort.Strings(keys)
				for _, key := range keys {
					err = c.DeleteKVStoreKey(&DeleteKVStoreKeyInput{StoreID: kvStore.StoreID, Key: key})
					if err != nil {
						t.Errorf("error during key cleanup: %v", err)
					}
				}
			}
			if err := p.Err(); err != nil {
				t.Errorf("error during cleanup pagination: %v", err)
			}

			err = c.DeleteKVStore(&DeleteKVStoreInput{StoreID: kvStore.StoreID})
			if err != nil {
				t.Errorf("error during store cleanup: %v", err)
			}
		})
	}()

	// fetch the newly created store and verify it matches
	var getKVStoreResponse *KVStore
	record(t, "kv_store/get-store", func(c *Client) {
		getKVStoreResponse, err = c.GetKVStore(&GetKVStoreInput{StoreID: kvStore.StoreID})
	})
	if err != nil {
		t.Fatal(err)
	}

	if getKVStoreResponse.Name != kvStore.Name || getKVStoreResponse.StoreID != kvStore.StoreID {
		t.Errorf("error fetching info for store %q: got %q, want %q", createStoreName, getKVStoreResponse, kvStore)
	}

	// create a bunch of keys in our kv store
	keys := []string{"apple", "banana", "carrot", "dragonfruit", "eggplant"}

	record(t, "kv_store/create-keys", func(c *Client) {
		for i, key := range keys {
			err := c.InsertKVStoreKey(&InsertKVStoreKeyInput{StoreID: kvStore.StoreID, Key: key, Value: key + strconv.Itoa(i)})
			if err != nil {
				t.Errorf("error inserting key %q: %v", key, err)
			}
		}
	})

	record(t, "kv_store/check-keys", func(c *Client) {
		for i, key := range keys {
			got, err := c.GetKVStoreKey(&GetKVStoreKeyInput{StoreID: kvStore.StoreID, Key: key})
			if err != nil {
				t.Errorf("error fetching key %q: %v", key, err)
			}
			want := key + strconv.Itoa(i)
			if got != want {
				t.Errorf("mismatch fetching key %q: got %q, want %q", key, got, want)
			}
		}
	})

	record(t, "kv_store/batch-create-keys", func(c *Client) {
		keys := `{"key":"batch-1","value":"VkFMVUU="}
    {"key":"batch-2","value":"VkFMVUU="}`
		err := c.BatchModifyKVStoreKey(&BatchModifyKVStoreKeyInput{
			StoreID: kvStore.StoreID,
			Body:    strings.NewReader(keys),
		})
		if err != nil {
			t.Errorf("error inserting keys %q: %v", keys, err)
		}
	})

	allKeys := []string{"batch-1", "batch-2"}
	allKeys = append(allKeys, keys...)
	sort.Strings(allKeys)

	// fetch all keys and validate they match our input data
	var kvStoreListKeys *ListKVStoreKeysResponse
	record(t, "kv_store/list-keys", func(c *Client) {
		kvStoreListKeys, err = c.ListKVStoreKeys(&ListKVStoreKeysInput{
			Consistency: ConsistencyStrong,
			StoreID:     kvStore.StoreID,
		})
	})

	if err != nil {
		t.Errorf("error listing keys: %v", err)
	}

	sort.Strings(kvStoreListKeys.Data)
	if !reflect.DeepEqual(allKeys, kvStoreListKeys.Data) {
		t.Errorf("mismatch listing keys: got %q, want %q", kvStoreListKeys.Data, allKeys)
	}

	record(t, "kv_store/list-keys-pagination", func(c *Client) {
		p := c.NewListKVStoreKeysPaginator(&ListKVStoreKeysInput{
			StoreID: kvStore.StoreID,
			Limit:   4,
		})
		var listed []string
		expected := []int{4, 3}
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
		if !reflect.DeepEqual(allKeys, listed) {
			t.Errorf("mismatch listing paginated keys: got %q, want %q", kvStoreListKeys.Data, allKeys)
		}
	})
}
