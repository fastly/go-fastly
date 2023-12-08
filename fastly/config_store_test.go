package fastly

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestClient_CreateConfigStore(t *testing.T) {
	t.Parallel()

	var (
		cs  *ConfigStore
		err error
	)
	record(t, fmt.Sprintf("config_store/%s/create_store", t.Name()), func(c *Client) {
		cs, err = c.CreateConfigStore(&CreateConfigStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating config store: %v", err)
	}

	// Ensure Config Store is cleaned up.
	t.Cleanup(func() {
		record(t, fmt.Sprintf("config_store/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteConfigStore(&DeleteConfigStoreInput{
				ID: cs.ID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting config store %q: %v", cs.ID, err)
		}
	})

	if got := cs.ID; len(got) == 0 {
		t.Errorf("ID: got %q, want not empty", got)
	}
	if got, want := cs.Name, t.Name(); got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
	if got := cs.CreatedAt; got == nil || got.IsZero() {
		t.Errorf("CreatedAt: got %v, want non-zero value", got)
	}
}

func TestClient_DeleteConfigStore(t *testing.T) {
	t.Parallel()

	var (
		cs  *ConfigStore
		err error
	)
	record(t, fmt.Sprintf("config_store/%s/create_store", t.Name()), func(c *Client) {
		cs, err = c.CreateConfigStore(&CreateConfigStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating config store: %v", err)
	}

	record(t, fmt.Sprintf("config_store/%s/delete_store", t.Name()), func(c *Client) {
		err = c.DeleteConfigStore(&DeleteConfigStoreInput{
			ID: cs.ID,
		})
	})
	if err != nil {
		t.Fatalf("DeleteConfigStore: error: %v", err)
	}

	// Verify that deleting a non-existent store gives an error.

	record(t, fmt.Sprintf("config_store/%s/delete_store_404", t.Name()), func(c *Client) {
		err = c.DeleteConfigStore(&DeleteConfigStoreInput{
			ID: cs.ID,
		})
	})

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		if got, want := httpErr.StatusCode, http.StatusNotFound; got != want {
			t.Fatalf("DeleteConfigStore HTTPError.StatusCode: got %d, want %d", got, want)
		}
		t.Logf("DeleteConfigStore: %v", httpErr)
	} else {
		t.Fatalf("DeleteConfigStore error: got %v (%T), want HTTPError", err, err)
	}
}

func TestClient_GetConfigStore(t *testing.T) {
	t.Parallel()

	var (
		cs  *ConfigStore
		err error
	)
	record(t, fmt.Sprintf("config_store/%s/create_store", t.Name()), func(c *Client) {
		cs, err = c.CreateConfigStore(&CreateConfigStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating config store: %v", err)
	}

	// Ensure Config Stores are cleaned up.
	t.Cleanup(func() {
		record(t, fmt.Sprintf("config_store/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteConfigStore(&DeleteConfigStoreInput{
				ID: cs.ID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting config store %q: %v", cs.ID, err)
		}
	})

	var getResult *ConfigStore
	record(t, fmt.Sprintf("config_store/%s/get_store", t.Name()), func(c *Client) {
		getResult, err = c.GetConfigStore(&GetConfigStoreInput{
			ID: cs.ID,
		})
	})
	if err != nil {
		t.Fatalf("GetConfigStore: error: %v", err)
	}

	if got, want := getResult.ID, cs.ID; got != want {
		t.Errorf("ID: got %q, want %q", got, want)
	}
	if got, want := getResult.Name, cs.Name; got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}

	// Verify that getting a non-existent store gives an error.

	record(t, fmt.Sprintf("config_store/%s/get_store_404", t.Name()), func(c *Client) {
		getResult, err = c.GetConfigStore(&GetConfigStoreInput{
			ID: "DOES-NOT-EXIST",
		})
	})

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		if got, want := httpErr.StatusCode, http.StatusNotFound; got != want {
			t.Fatalf("GetConfigStore HTTPError.StatusCode: got %d, want %d", got, want)
		}
		t.Logf("GetConfigStore: %v", httpErr)
	} else {
		t.Fatalf("GetConfigStore error: got %v (%T), want HTTPError", err, err)
	}
}

func TestClient_GetConfigStoreMetadata(t *testing.T) {
	t.Parallel()

	var (
		cs  *ConfigStore
		err error
	)
	record(t, fmt.Sprintf("config_store/%s/create_store", t.Name()), func(c *Client) {
		cs, err = c.CreateConfigStore(&CreateConfigStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating config store: %v", err)
	}

	// Ensure Config Stores are cleaned up.
	t.Cleanup(func() {
		record(t, fmt.Sprintf("config_store/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteConfigStore(&DeleteConfigStoreInput{
				ID: cs.ID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting config store %q: %v", cs.ID, err)
		}
	})

	var metadataResult *ConfigStoreMetadata
	record(t, fmt.Sprintf("config_store/%s/get_store_metadata", t.Name()), func(c *Client) {
		metadataResult, err = c.GetConfigStoreMetadata(&GetConfigStoreMetadataInput{
			ID: cs.ID,
		})
	})
	if err != nil {
		t.Fatalf("GetConfigStoreMetadata: error: %v", err)
	}

	if got, want := metadataResult.ItemCount, 0; got != want {
		t.Errorf("ItemCount: got %d, want %d", got, want)
	}

	// Verify that getting a non-existent store gives an error.

	record(t, fmt.Sprintf("config_store/%s/get_store_metadata_404", t.Name()), func(c *Client) {
		metadataResult, err = c.GetConfigStoreMetadata(&GetConfigStoreMetadataInput{
			ID: "DOES-NOT-EXIST",
		})
	})

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		if got, want := httpErr.StatusCode, http.StatusNotFound; got != want {
			t.Fatalf("GetConfigStoreMetadata HTTPError.StatusCode: got %d, want %d", got, want)
		}
		t.Logf("GetConfigStoreMetadata: %v", httpErr)
	} else {
		t.Fatalf("GetConfigStoreMetadata error: got %v (%T), want HTTPError", err, err)
	}
}

func TestClient_ListConfigStores(t *testing.T) {
	var (
		css []*ConfigStore
		err error
	)

	// Verify list works when there are no stores.
	record(t, fmt.Sprintf("config_store/%s/empty", t.Name()), func(c *Client) {
		css, err = c.ListConfigStores(&ListConfigStoresInput{})
	})
	if err != nil {
		t.Fatalf("ListConfigStores: unexpected error: %v", err)
	}
	if got, want := len(css), 0; got != want {
		t.Fatalf("ListConfigStores: got %d entries, want %d\n%+v", got, want, css[0])
	}

	// Create stores to be able to list them.
	var (
		stores []*ConfigStore
		cs     *ConfigStore
	)
	record(t, fmt.Sprintf("config_store/%s/create_stores", t.Name()), func(c *Client) {
		for i := 0; i < 5; i++ {
			cs, err = c.CreateConfigStore(&CreateConfigStoreInput{
				Name: fmt.Sprintf("%s-%02d", t.Name(), i),
			})
			if err != nil {
				t.Fatalf("error creating config store: %v", err)
			}
			t.Log(cs)
			stores = append(stores, cs)
		}
	})
	// Ensure Config Stores are cleaned up.
	t.Cleanup(func() {
		record(t, fmt.Sprintf("config_store/%s/delete_stores", t.Name()), func(c *Client) {
			for _, cs := range stores {
				err = c.DeleteConfigStore(&DeleteConfigStoreInput{
					ID: cs.ID,
				})
				if err != nil {
					t.Fatalf("error deleting config store %q: %v", cs.ID, err)
				}
			}
		})
	})

	record(t, fmt.Sprintf("config_store/%s/list", t.Name()), func(c *Client) {
		css, err = c.ListConfigStores(&ListConfigStoresInput{})
	})

	if got, want := len(css), len(stores); got != want {
		t.Fatalf("ListConfigStores: got %d entries, want %d", got, want)
	}

	for i, cs := range css {
		if got, want := cs.ID, stores[i].ID; got != want {
			t.Errorf("ListConfigStores: index %d: ID: got %q, want %q", i, got, want)
		}
	}

	record(t, fmt.Sprintf("config_store/%s/list-with-name", t.Name()), func(c *Client) {
		css, err = c.ListConfigStores(&ListConfigStoresInput{Name: stores[0].Name})
	})

	if got, want := len(css), 1; got != want {
		t.Fatalf("ListConfigStores: got %d entries, want %d", got, want)
	}

	if got, want := css[0].ID, stores[0].ID; got != want {
		t.Errorf("ListConfigStores: index %d: ID: got %q, want %q", 0, got, want)
	}
}

func TestClient_ListConfigStoreServices(t *testing.T) {
	var (
		cs  *ConfigStore
		err error
	)
	record(t, fmt.Sprintf("config_store/%s/create_store", t.Name()), func(c *Client) {
		cs, err = c.CreateConfigStore(&CreateConfigStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating config store: %v", err)
	}

	// Ensure Config Stores are cleaned up.
	t.Cleanup(func() {
		record(t, fmt.Sprintf("config_store/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteConfigStore(&DeleteConfigStoreInput{
				ID: cs.ID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting config store %q: %v", cs.ID, err)
		}
	})

	var servicesResult []*Service
	record(t, fmt.Sprintf("config_store/%s/list_services", t.Name()), func(c *Client) {
		servicesResult, err = c.ListConfigStoreServices(&ListConfigStoreServicesInput{
			ID: cs.ID,
		})
	})

	if got, want := len(servicesResult), 0; got != want {
		t.Fatalf("ListConfigStoreServices: got %d entries, want %d", got, want)
	}
}

func TestClient_UpdateConfigStore(t *testing.T) {
	t.Parallel()

	var (
		cs  *ConfigStore
		err error
	)
	record(t, fmt.Sprintf("config_store/%s/create_store", t.Name()), func(c *Client) {
		cs, err = c.CreateConfigStore(&CreateConfigStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating config store: %v", err)
	}

	// Ensure Config Stores are cleaned up.
	t.Cleanup(func() {
		record(t, fmt.Sprintf("config_store/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteConfigStore(&DeleteConfigStoreInput{
				ID: cs.ID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting config store %q: %v", cs.ID, err)
		}
	})

	const updatedName = "UPDATED NAME!"
	var updateResult *ConfigStore
	record(t, fmt.Sprintf("config_store/%s/update_store", t.Name()), func(c *Client) {
		updateResult, err = c.UpdateConfigStore(&UpdateConfigStoreInput{
			ID:   cs.ID,
			Name: updatedName,
		})
	})
	if err != nil {
		t.Fatalf("UpdateConfigStore: error: %v", err)
	}

	if got, want := updateResult.ID, cs.ID; got != want {
		t.Errorf("ID: got %q, want %q", got, want)
	}
	if got, want := updateResult.Name, updatedName; got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
}
