package fastly

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

// Config Store Item.
// https://developer.fastly.com/reference/api/services/resources/config-store-item/

// ConfigStoreItem is a config store item response from the Fastly API.
type ConfigStoreItem struct {
	StoreID   string     `json:"store_id"`
	Key       string     `json:"item_key"`
	Value     string     `json:"item_value"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// CreateConfigStoreItemInput is the input to the CreateConfigStoreItem.
type CreateConfigStoreItemInput struct {
	// StoreID is the ID of the config store (required).
	StoreID string
	// Key is the item's name (required).
	Key string `url:"item_key"`
	// Value is the item's value (required).
	Value string `url:"item_value"`
}

// CreateConfigStoreItem creates a new Fastly config store item.
func (c *Client) CreateConfigStoreItem(ctx context.Context, i *CreateConfigStoreItemInput) (*ConfigStoreItem, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID, "item")

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.PostForm(ctx, path, i, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var csi *ConfigStoreItem
	if err = json.NewDecoder(resp.Body).Decode(&csi); err != nil {
		return nil, err
	}

	return csi, nil
}

// DeleteConfigStoreItemInput is the input to DeleteConfigStoreItem.
type DeleteConfigStoreItemInput struct {
	// StoreID is the ID of the item's config store (required).
	StoreID string
	// Key is the name of the config store item to delete (required).
	Key string
}

// DeleteConfigStoreItem deletes the given config store item.
func (c *Client) DeleteConfigStoreItem(ctx context.Context, i *DeleteConfigStoreItemInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}
	if i.Key == "" {
		return ErrMissingKey
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID, "item", i.Key)

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

// GetConfigStoreItemInput is the input to the GetConfigStoreItem.
type GetConfigStoreItemInput struct {
	// StoreID is the ID of the item's config store (required).
	StoreID string
	// Key is the name of the config store item to fetch (required).
	Key string
}

// GetConfigStoreItem gets the config store item with the given parameters.
func (c *Client) GetConfigStoreItem(ctx context.Context, i *GetConfigStoreItemInput) (*ConfigStoreItem, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}
	if i.Key == "" {
		return nil, ErrMissingKey
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID, "item", i.Key)

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var csi *ConfigStoreItem
	if err = json.NewDecoder(resp.Body).Decode(&csi); err != nil {
		return nil, err
	}

	return csi, nil
}

// ListConfigStoreItemsInput is the input to ListConfigStoreItems.
type ListConfigStoreItemsInput struct {
	// StoreID is the ID of the config store to retrieve items for (required).
	StoreID string
}

// ListConfigStoreItems returns a list of config store items for the given store.
func (c *Client) ListConfigStoreItems(ctx context.Context, i *ListConfigStoreItemsInput) ([]*ConfigStoreItem, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID, "items")

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var csi []*ConfigStoreItem
	if err = json.NewDecoder(resp.Body).Decode(&csi); err != nil {
		return nil, err
	}
	sort.Slice(csi, func(i, j int) bool {
		return csi[i].Key < csi[j].Key
	})

	return csi, nil
}

// UpdateConfigStoreItemInput is the input to the UpdateConfigStoreItem.
type UpdateConfigStoreItemInput struct {
	// Key is the name of the config store item to update (required).
	Key string
	// StoreID is the ID of the item's config store (required).
	StoreID string
	// Upsert, if true, will insert or update an item. Otherwise, update an item which must already exist.
	Upsert bool
	// Value is the new item's value (required).
	Value string `url:"item_value"`
}

// UpdateConfigStoreItem updates a specific config store item.
func (c *Client) UpdateConfigStoreItem(ctx context.Context, i *UpdateConfigStoreItemInput) (*ConfigStoreItem, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}
	if i.Key == "" {
		return nil, ErrMissingKey
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID, "item", i.Key)

	var httpMethod string
	if i.Upsert {
		// Insert or update an entry in a config store given a config store ID, item key, and item value.
		httpMethod = http.MethodPut
	} else {
		// Update an entry in a config store given a config store ID, item key, and item value.
		httpMethod = http.MethodPatch
	}

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = JSONMimeType
	requestOptions.Parallel = true

	resp, err := c.RequestForm(ctx, httpMethod, path, i, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var csi *ConfigStoreItem
	if err = json.NewDecoder(resp.Body).Decode(&csi); err != nil {
		return nil, err
	}

	return csi, nil
}

// BatchModifyConfigStoreItemsInput is the input parameter to the
// BatchModifyConfigStoreItems function.
type BatchModifyConfigStoreItemsInput struct {
	// Items is a list of Config Store items.
	Items []*BatchConfigStoreItem `json:"items"`
	// StoreID is the ID of the Config Store to modify items for (required).
	StoreID string `json:"-"`
}

// BatchDictionaryItem represents a dictionary item.
type BatchConfigStoreItem struct {
	// ItemKey is an item key (maximum 256 characters).
	ItemKey string `json:"item_key"`
	// ItemValue is an item value (maximum 8000 characters).
	ItemValue string `json:"item_value"`
	// Operation is a batching operation variant.
	Operation BatchOperation `json:"op"`
}

// BatchModifyConfigStoreItems bulk updates dictionary items.
func (c *Client) BatchModifyConfigStoreItems(ctx context.Context, i *BatchModifyConfigStoreItemsInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}
	if len(i.Items) > BatchModifyMaximumOperations {
		return ErrMaxExceededItems
	}

	path := ToSafeURL("resources", "stores", "config", i.StoreID, "items")

	resp, err := c.PatchJSON(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var batchModifyResult map[string]string
	return DecodeBodyMap(resp.Body, &batchModifyResult)
}
