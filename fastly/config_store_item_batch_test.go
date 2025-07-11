package fastly

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_ConfigStoreBatch(t *testing.T) {
	t.Parallel()

	cs := createConfigStoreForBatch(t)

	var err error

	Record(t, "config_store_batch/create_and_upsert", func(c *Client) {
		err = c.BatchModifyConfigStoreItems(context.TODO(), &BatchModifyConfigStoreItemsInput{
			StoreID: cs.StoreID,
			Items: []*BatchConfigStoreItem{
				{
					ItemKey:   "key1",
					ItemValue: "value1",
					Operation: CreateBatchOperation,
				},
				{
					ItemKey:   "key2",
					ItemValue: "value2",
					Operation: UpsertBatchOperation,
				},
			},
		})
	})
	if err != nil {
		t.Fatalf("error batch creating config store items: %v", err)
	}

	Record(t, "config_store_batch/update_and_upsert", func(c *Client) {
		err = c.BatchModifyConfigStoreItems(context.TODO(), &BatchModifyConfigStoreItemsInput{
			StoreID: cs.StoreID,
			Items: []*BatchConfigStoreItem{
				{
					ItemKey:   "key1",
					ItemValue: "value2",
					Operation: UpdateBatchOperation,
				},
				{
					ItemKey:   "key2",
					ItemValue: "value3",
					Operation: UpsertBatchOperation,
				},
			},
		})
	})
	if err != nil {
		t.Fatalf("error batch updating config store items: %v", err)
	}

	Record(t, "config_store_batch/delete", func(c *Client) {
		err = c.BatchModifyConfigStoreItems(context.TODO(), &BatchModifyConfigStoreItemsInput{
			StoreID: cs.StoreID,
			Items: []*BatchConfigStoreItem{
				{
					ItemKey:   "key1",
					Operation: DeleteBatchOperation,
				},
				{
					ItemKey:   "key2",
					Operation: DeleteBatchOperation,
				},
			},
		})
	})
	if err != nil {
		t.Fatalf("error batch deleting config store items: %v", err)
	}
}

func createConfigStoreForBatch(t *testing.T) *ConfigStore {
	t.Helper()

	var (
		cs  *ConfigStore
		err error
	)
	Record(t, fmt.Sprintf("config_store_batch/%s/create_store", t.Name()), func(c *Client) {
		cs, err = c.CreateConfigStore(context.TODO(), &CreateConfigStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating config store: %v", err)
	}

	// Ensure Config Store is cleaned up.
	t.Cleanup(func() {
		Record(t, fmt.Sprintf("config_store_batch/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteConfigStore(context.TODO(), &DeleteConfigStoreInput{
				StoreID: cs.StoreID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting config store %q: %v", cs.StoreID, err)
		}
	})

	return cs
}
