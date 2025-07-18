package fastly

import (
	"context"
	"encoding/json"
	"sort"
	"time"
)

// Config Store.
// A container that lets you store data in key-value pairs.
// https://developer.fastly.com/reference/api/services/resources/config-store/
// https://developer.fastly.com/reference/api/services/resources/config-store-item/

// ConfigStore represents a config store response from the Fastly API.
type ConfigStore struct {
	CreatedAt *time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	StoreID   string     `json:"id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// ConfigStoreMetadata represents a config store metadata response from the Fastly API.
type ConfigStoreMetadata struct {
	// ItemCount is the number of items in a store.
	ItemCount int `json:"item_count"`
}

// CreateConfigStoreInput is the input to CreateConfigStore.
type CreateConfigStoreInput struct {
	// Name is the name of the store to create (required).
	Name string `url:"name"`
}

// CreateConfigStore creates a new Fastly config store.
func (c *Client) CreateConfigStore(ctx context.Context, i *CreateConfigStoreInput) (*ConfigStore, error) {
	path := "/resources/stores/config"

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.PostForm(ctx, path, i, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs *ConfigStore
	if err = json.NewDecoder(resp.Body).Decode(&cs); err != nil {
		return nil, err
	}

	return cs, nil
}

// DeleteConfigStoreInput is the input parameter to DeleteConfigStore.
type DeleteConfigStoreInput struct {
	// StoreID is the StoreID of the config store to delete (required).
	StoreID string
}

// DeleteConfigStore deletes the given config store version.
func (c *Client) DeleteConfigStore(ctx context.Context, i *DeleteConfigStoreInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID)

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.Delete(ctx, path, requestOptions)
	if err != nil {
		return err
	}

	// This endpoint returns a '200 Ok' on successful deletion,
	// which c.Delete verifies. The response body will be: '{"status":"ok"}'
	// on success, which we ignore.
	return resp.Body.Close()
}

// GetConfigStoreInput is the input to GetConfigStore.
type GetConfigStoreInput struct {
	// StoreID is the StoreID of the config store (required).
	StoreID string
}

// GetConfigStore returns the config store for the given input parameters.
func (c *Client) GetConfigStore(ctx context.Context, i *GetConfigStoreInput) (*ConfigStore, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID)

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs *ConfigStore
	if err = json.NewDecoder(resp.Body).Decode(&cs); err != nil {
		return nil, err
	}

	return cs, nil
}

// GetConfigStoreMetadataInput is the input to GetConfigStoreMetadata.
type GetConfigStoreMetadataInput struct {
	// StoreID is the StoreID of the config store (required).
	StoreID string
}

// GetConfigStoreMetadata returns the config store's metadata for the given input parameters.
func (c *Client) GetConfigStoreMetadata(ctx context.Context, i *GetConfigStoreMetadataInput) (*ConfigStoreMetadata, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID, "info")

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var csm *ConfigStoreMetadata
	if err = json.NewDecoder(resp.Body).Decode(&csm); err != nil {
		return nil, err
	}

	return csm, nil
}

// ListConfigStoreServicesInput is the input to ListConfigStoreServices.
type ListConfigStoresInput struct {
	// Name is the name of a config store.
	Name string
}

// ListConfigStores returns a list of config stores sorted by name.
func (c *Client) ListConfigStores(ctx context.Context, i *ListConfigStoresInput) ([]*ConfigStore, error) {
	path := "/resources/stores/config"

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	if i.Name != "" {
		requestOptions.Params["name"] = i.Name
	}

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var css []*ConfigStore
	if err = json.NewDecoder(resp.Body).Decode(&css); err != nil {
		return nil, err
	}
	sort.Slice(css, func(i, j int) bool {
		return css[i].Name < css[j].Name
	})

	return css, nil
}

// ListConfigStoreServicesInput is the input to ListConfigStoreServices.
type ListConfigStoreServicesInput struct {
	// StoreID is the StoreID of the config store (required).
	StoreID string
}

// ListConfigStoreServices returns the list of services that are associated with
// a given config store.
func (c *Client) ListConfigStoreServices(ctx context.Context, i *ListConfigStoreServicesInput) ([]*Service, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID, "services")

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ss []*Service
	if err = DecodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

// UpdateConfigStoreInput is the input to UpdateConfigStore.
type UpdateConfigStoreInput struct {
	// Name is the new name of the config store (required).
	Name string `url:"name"`
	// StoreID is the StoreID of the config store to update (required).
	StoreID string
}

// UpdateConfigStore updates a specific config store.
func (c *Client) UpdateConfigStore(ctx context.Context, i *UpdateConfigStoreInput) (*ConfigStore, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID)

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.PutForm(ctx, path, i, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs *ConfigStore
	if err = json.NewDecoder(resp.Body).Decode(&cs); err != nil {
		return nil, err
	}

	return cs, nil
}
