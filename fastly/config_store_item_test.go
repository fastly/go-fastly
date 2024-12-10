package fastly

import (
	"fmt"
	"testing"
)

func TestClient_CreateConfigStoreItem(t *testing.T) {
	t.Parallel()

	cs := createConfigStore(t)

	var (
		item *ConfigStoreItem
		err  error
	)
	const value = "testing 123"

	Record(t, fmt.Sprintf("config_store_item/%s/create_item", t.Name()), func(c *Client) {
		item, err = c.CreateConfigStoreItem(&CreateConfigStoreItemInput{
			StoreID: cs.StoreID,
			Key:     t.Name(),
			Value:   value,
		})
	})
	if err != nil {
		t.Fatalf("error creating config store item: %v", err)
	}

	if got, want := item.StoreID, cs.StoreID; got != want {
		t.Errorf("StoreID: got %q, want %q", got, want)
	}
	if got, want := item.Key, t.Name(); got != want {
		t.Errorf("Key: got %q, want %q", got, want)
	}
	if got, want := item.Value, value; got != want {
		t.Errorf("Value: got %q, want %q", got, want)
	}
	if got := item.CreatedAt; got == nil || got.IsZero() {
		t.Errorf("CreatedAt: got %v, want non-zero value", got)
	}
}

func TestClient_DeleteConfigStoreItem(t *testing.T) {
	t.Parallel()

	cs := createConfigStore(t)

	var err error

	Record(t, fmt.Sprintf("config_store_item/%s/create_item", t.Name()), func(c *Client) {
		_, err = c.CreateConfigStoreItem(&CreateConfigStoreItemInput{
			StoreID: cs.StoreID,
			Key:     t.Name(),
			Value:   "delete me",
		})
	})
	if err != nil {
		t.Fatalf("error creating config store item: %v", err)
	}

	Record(t, fmt.Sprintf("config_store_item/%s/delete_item", t.Name()), func(c *Client) {
		err = c.DeleteConfigStoreItem(&DeleteConfigStoreItemInput{
			StoreID: cs.StoreID,
			Key:     t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error deleting config store item: %v", err)
	}
}

func TestClient_GetConfigStoreItem(t *testing.T) {
	t.Parallel()

	cs := createConfigStore(t)

	var (
		gotItem *ConfigStoreItem
		item    *ConfigStoreItem
		err     error
	)

	Record(t, fmt.Sprintf("config_store_item/%s/create_item", t.Name()), func(c *Client) {
		item, err = c.CreateConfigStoreItem(&CreateConfigStoreItemInput{
			StoreID: cs.StoreID,
			Key:     t.Name(),
			Value:   "get me",
		})
	})
	if err != nil {
		t.Fatalf("error creating config store item: %v", err)
	}

	Record(t, fmt.Sprintf("config_store_item/%s/get_item", t.Name()), func(c *Client) {
		gotItem, err = c.GetConfigStoreItem(&GetConfigStoreItemInput{
			StoreID: cs.StoreID,
			Key:     t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error getting config store item: %v", err)
	}

	if got, want := gotItem.StoreID, item.StoreID; got != want {
		t.Errorf("StoreID: got %q, want %q", got, want)
	}
	if got, want := gotItem.Key, item.Key; got != want {
		t.Errorf("Key: got %q, want %q", got, want)
	}
	if got, want := gotItem.Value, item.Value; got != want {
		t.Errorf("Value: got %q, want %q", got, want)
	}
	if got, want := gotItem.CreatedAt, item.CreatedAt; got == nil || !got.Equal(*want) {
		t.Errorf("CreatedAt: got %v, want %v", got, want)
	}
}

func TestClient_ListConfigStoreItems(t *testing.T) {
	t.Parallel()

	cs := createConfigStore(t)

	var (
		keys = make([]string, 5)
		err  error
	)

	for i := range keys {
		keys[i] = fmt.Sprintf("%s-key-%02d", t.Name(), i)
	}

	Record(t, fmt.Sprintf("config_store_item/%s/create_items", t.Name()), func(c *Client) {
		for i, key := range keys {
			_, err = c.CreateConfigStoreItem(&CreateConfigStoreItemInput{
				StoreID: cs.StoreID,
				Key:     key,
				Value:   fmt.Sprintf("value %02d", i),
			})
			if err != nil {
				break
			}
		}
	})
	if err != nil {
		t.Fatalf("error creating config store item: %v", err)
	}

	var list []*ConfigStoreItem
	Record(t, fmt.Sprintf("config_store_item/%s/list_items", t.Name()), func(c *Client) {
		list, err = c.ListConfigStoreItems(&ListConfigStoreItemsInput{
			StoreID: cs.StoreID,
		})
	})
	if err != nil {
		t.Fatalf("error listing config store items: %v", err)
	}

	if got, want := len(list), len(keys); got != want {
		t.Fatalf("got %d items, want %d", got, want)
	}

	for i, gotItem := range list {
		if got, want := gotItem.StoreID, cs.StoreID; got != want {
			t.Errorf("StoreID: got %q, want %q", got, want)
		}
		if got, want := gotItem.Key, keys[i]; got != want {
			t.Errorf("Key: got %q, want %q", got, want)
		}
	}
}

func TestClient_UpdateConfigStoreItem(t *testing.T) {
	t.Parallel()

	cs := createConfigStore(t)

	var (
		gotItem *ConfigStoreItem
		err     error
	)

	const newValue = "I'm a new value"
	Record(t, fmt.Sprintf("config_store_item/%s/create_and_update_item", t.Name()), func(c *Client) {
		_, err = c.CreateConfigStoreItem(&CreateConfigStoreItemInput{
			StoreID: cs.StoreID,
			Key:     t.Name(),
			Value:   "OLD VALUE",
		})
		if err != nil {
			return
		}
		_, err = c.UpdateConfigStoreItem(&UpdateConfigStoreItemInput{
			StoreID: cs.StoreID,
			Key:     t.Name(),
			Value:   newValue,
		})
		if err != nil {
			return
		}
		gotItem, err = c.GetConfigStoreItem(&GetConfigStoreItemInput{
			StoreID: cs.StoreID,
			Key:     t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error updating config store item: %v", err)
	}

	if got, want := gotItem.Value, newValue; got != want {
		t.Errorf("Value: got %q, want %q", got, want)
	}

	// Upsert.

	const upsertValue = "i was upserted"
	Record(t, fmt.Sprintf("config_store_item/%s/upsert_item", t.Name()), func(c *Client) {
		gotItem, err = c.UpdateConfigStoreItem(&UpdateConfigStoreItemInput{
			Upsert:  true,
			StoreID: cs.StoreID,
			Key:     t.Name() + "-upsert",
			Value:   upsertValue,
		})
		if err != nil {
			return
		}
	})

	if err != nil {
		t.Fatalf("error upserting config store item: %v", err)
	}

	if got, want := gotItem.Value, upsertValue; got != want {
		t.Errorf("Value: got %q, want %q", got, want)
	}
}

func createConfigStore(t *testing.T) *ConfigStore {
	t.Helper()

	var (
		cs  *ConfigStore
		err error
	)
	Record(t, fmt.Sprintf("config_store_item/%s/create_store", t.Name()), func(c *Client) {
		cs, err = c.CreateConfigStore(&CreateConfigStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating config store: %v", err)
	}

	// Ensure Config Store is cleaned up.
	t.Cleanup(func() {
		Record(t, fmt.Sprintf("config_store_item/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteConfigStore(&DeleteConfigStoreInput{
				StoreID: cs.StoreID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting config store %q: %v", cs.StoreID, err)
		}
	})

	return cs
}
