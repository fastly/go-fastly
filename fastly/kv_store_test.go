package fastly

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_KVStore(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	t.Parallel()

	var err error

	const createStoreName = "kv-store-test-store"

	// List
	var kvStoreListResp1 *ListKVStoresResponse
	Record(t, "kv_store/list-store", func(c *Client) {
		kvStoreListResp1, err = c.ListKVStores(nil)
	})
	require.NoError(err)
	require.NotContains(kvStoreListResp1.Data, createStoreName)

	// Create
	var kvStore *KVStore
	Record(t, "kv_store/create-store", func(c *Client) {
		kvStore, err = c.CreateKVStore(&CreateKVStoreInput{
			Name: createStoreName,
		})
	})
	require.NoError(err)
	require.Equal(createStoreName, kvStore.Name, "unexpected store name")

	// ensure we delete it
	defer func() {
		Record(t, "kv_store/cleanup", func(c *Client) {
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
					require.NoError(err, "key cleanup")
				}
			}
			require.NoError(p.Err(), "cleanup pagination")

			err = c.DeleteKVStore(&DeleteKVStoreInput{StoreID: kvStore.StoreID})
			require.NoError(err, "store cleanup")
		})
	}()

	// fetch the newly created store and verify it matches
	var getKVStoreResponse *KVStore
	Record(t, "kv_store/get-store", func(c *Client) {
		getKVStoreResponse, err = c.GetKVStore(&GetKVStoreInput{StoreID: kvStore.StoreID})
	})
	require.NoError(err)
	require.Equal(kvStore.Name, getKVStoreResponse.Name, "error fetching info for store")
	require.Equal(kvStore.StoreID, getKVStoreResponse.StoreID, "error fetching info for store")

	// create a bunch of keys in our kv store
	keys := []string{"apple", "banana", "carrot", "dragonfruit", "eggplant"}

	Record(t, "kv_store/create-keys", func(c *Client) {
		for i, key := range keys {
			err := c.InsertKVStoreKey(&InsertKVStoreKeyInput{StoreID: kvStore.StoreID, Key: key, Value: key + strconv.Itoa(i)})
			require.NoErrorf(err, "inserting key %q", key)
		}
	})

	Record(t, "kv_store/check-keys", func(c *Client) {
		for i, key := range keys {
			got, err := c.GetKVStoreKey(&GetKVStoreKeyInput{StoreID: kvStore.StoreID, Key: key})
			assert.NoErrorf(err, "fetching key %q", key)
			want := key + strconv.Itoa(i)
			assert.Equalf(want, got, "incorrect value key %q", key)
		}
	})

	Record(t, "kv_store/batch-create-keys", func(c *Client) {
		keys := `{"key":"batch-1","value":"VkFMVUU="}
    {"key":"batch-2","value":"VkFMVUU="}`
		err := c.BatchModifyKVStoreKey(&BatchModifyKVStoreKeyInput{
			StoreID: kvStore.StoreID,
			Body:    strings.NewReader(keys),
		})
		require.NoError(err, "batch inserting keys")
	})

	allKeys := []string{"batch-1", "batch-2"}
	allKeys = append(allKeys, keys...)

	// fetch all keys and validate they match our input data
	var kvStoreListKeys *ListKVStoreKeysResponse
	Record(t, "kv_store/list-keys", func(c *Client) {
		kvStoreListKeys, err = c.ListKVStoreKeys(&ListKVStoreKeysInput{
			Consistency: ConsistencyStrong,
			StoreID:     kvStore.StoreID,
		})
	})

	require.NoError(err, "listing keys")
	assert.ElementsMatch(allKeys, kvStoreListKeys.Data, "mismatch listing keys")

	Record(t, "kv_store/list-keys-pagination", func(c *Client) {
		p := c.NewListKVStoreKeysPaginator(&ListKVStoreKeysInput{
			StoreID: kvStore.StoreID,
			Limit:   4,
		})
		var listed []string
		expected := []int{4, 3}
		var page int
		for p.Next() {
			keys := p.Keys()
			assert.Equalf(expected[page], len(keys), "wrong number of keys for page %d", page)
			page++
			listed = append(listed, keys...)
		}
		assert.NoError(p.Err(), "keys pagination")
		assert.ElementsMatch(allKeys, listed, "mismatch listing paginated keys")
	})

	Record(t, "kv_store/list-keys-prefix", func(c *Client) {
		p := c.NewListKVStoreKeysPaginator(&ListKVStoreKeysInput{
			StoreID: kvStore.StoreID,
			Prefix:  "ba",
		})
		expectedKeys := []string{"banana", "batch-1", "batch-2"}
		var listed []string
		var page int
		for p.Next() {
			keys := p.Keys()
			assert.Equalf(len(expectedKeys), len(keys), "wrong number of keys for page %d", page)
			page++
			listed = append(listed, keys...)
		}
		assert.NoError(p.Err(), "keys pagination")
		assert.ElementsMatch(expectedKeys, listed, "mismatch listing prefixed keys")
	})

	testKey := "apple"
	expectedValue := "apple0"
	expectedMetadata := "meta"

	var item GetKVStoreItemOutput

	Record(t, "kv_store/get-item-round-1", func(c *Client) {
		item, err = c.GetKVStoreItem(&GetKVStoreItemInput{
			StoreID: kvStore.StoreID,
			Key:     testKey,
		})
		require.NoErrorf(err, "fetching key %q", testKey)
	})

	Record(t, "kv_store/insert-item-add-failure", func(c *Client) {
		err := c.InsertKVStoreKey(&InsertKVStoreKeyInput{
			StoreID: kvStore.StoreID,
			Key:     testKey,
			Add:     true,
		})
		assert.Errorf(err, "adding existing key %q should have failed", testKey)
	})

	Record(t, "kv_store/insert-item-prepend-set-metadata", func(c *Client) {
		err := c.InsertKVStoreKey(&InsertKVStoreKeyInput{
			StoreID:  kvStore.StoreID,
			Key:      testKey,
			Prepend:  true,
			Value:    "prefix",
			Metadata: ToPointer("meta"),
		})
		require.NoErrorf(err, "updating key %q", testKey)

		expectedValue = "prefix" + expectedValue
		expectedMetadata = "meta"
	})

	Record(t, "kv_store/get-item-round-2", func(c *Client) {
		updatedItem, err := c.GetKVStoreItem(&GetKVStoreItemInput{
			StoreID: kvStore.StoreID,
			Key:     testKey,
		})
		require.NoErrorf(err, "fetching key %q", testKey)

		updatedValue, err := updatedItem.ValueAsString()
		require.NoError(err, "reading updated value")

		assert.Equal(expectedValue, updatedValue, "mismatch of updated value")
		assert.Equal(expectedMetadata, updatedItem.Metadata, "mismatch of updated metadata")
		assert.NotEqual(item.Generation, updatedItem.Generation, "generation marker change")

		item = updatedItem
	})

	Record(t, "kv_store/insert-item-generation-match-failure", func(c *Client) {
		err := c.InsertKVStoreKey(&InsertKVStoreKeyInput{
			StoreID:           kvStore.StoreID,
			Key:               testKey,
			IfGenerationMatch: item.Generation + 1,
		})
		testName := "insert-item-generation-match-failure"
		var herr *HTTPError
		require.ErrorAs(err, &herr, testName)
		require.Truef(herr.IsPreconditionFailed(), testName+" expected HTTP 'Precondition Failed'")
	})

	Record(t, "kv_store/insert-item-append", func(c *Client) {
		err := c.InsertKVStoreKey(&InsertKVStoreKeyInput{
			StoreID: kvStore.StoreID,
			Key:     testKey,
			Append:  true,
			Value:   "suffix",
		})
		require.NoErrorf(err, "updating key %q", testKey)

		expectedValue = expectedValue + "suffix"
	})

	Record(t, "kv_store/get-item-round-3", func(c *Client) {
		updatedItem, err := c.GetKVStoreItem(&GetKVStoreItemInput{
			StoreID: kvStore.StoreID,
			Key:     testKey,
		})
		require.NoErrorf(err, "fetching key %q", testKey)

		updatedValue, err := updatedItem.ValueAsString()
		require.NoError(err, "reading updated value")

		assert.Equal(expectedValue, updatedValue, "mismatch of updated value")

		item = updatedItem
	})

	Record(t, "kv_store/delete-item-nonexistent-suppress-error", func(c *Client) {
		err := c.DeleteKVStoreKey(&DeleteKVStoreKeyInput{
			StoreID: kvStore.StoreID,
			Key:     testKey + "23",
			Force:   true,
		})
		assert.NoError(err, "error for deleting non-existent key should have been suppressed")
	})

	Record(t, "kv_store/delete-item-generation-match-failure", func(c *Client) {
		err := c.DeleteKVStoreKey(&DeleteKVStoreKeyInput{
			StoreID:           kvStore.StoreID,
			Key:               testKey,
			IfGenerationMatch: item.Generation + 1,
		})
		testName := "delete-item-generation-match-failure"
		var herr *HTTPError
		require.ErrorAs(err, &herr, testName)
		require.Truef(herr.IsPreconditionFailed(), testName+" expected HTTP 'Precondition Failed'")
	})
}

func TestClient_CreateKVStoresWithLocations(t *testing.T) {
	var (
		stores []*KVStore
		ks     *KVStore
		err    error
	)

	Record(t, fmt.Sprintf("kv_store/%s/create_stores", t.Name()), func(c *Client) {
		for _, location := range []string{"US", "EU", "ASIA", "AUS"} {
			ks, err = c.CreateKVStore(&CreateKVStoreInput{
				Name:     fmt.Sprintf("%s_%s", t.Name(), location),
				Location: location,
			})
			if err != nil {
				t.Fatalf("error creating kv store: %v", err)
			}

			if got := ks.StoreID; len(got) == 0 {
				t.Errorf("ID: got %q, want not empty", got)
			}
			if got, want := ks.Name, fmt.Sprintf("%s_%s", t.Name(), location); got != want {
				t.Errorf("Name: got %q, want %q", got, want)
			}

			stores = append(stores, ks)
		}
	})

	t.Cleanup(func() {
		Record(t, fmt.Sprintf("kv_store/%s/delete_stores", t.Name()), func(c *Client) {
			for _, ks := range stores {
				err = c.DeleteKVStore(&DeleteKVStoreInput{
					StoreID: ks.StoreID,
				})
				if err != nil {
					t.Fatalf("error deleting kv store %q: %v", ks.StoreID, err)
				}
			}
		})
	})
}
