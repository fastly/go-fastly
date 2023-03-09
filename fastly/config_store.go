package fastly

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Config Store.
// A container that lets you store data in key-value pairs.
// https://developer.fastly.com/reference/api/services/resources/config-store/
// https://developer.fastly.com/reference/api/services/resources/config-store-item/

// ConfigStore represents a config store response from the Fastly API.
type ConfigStore struct {
	Name      string     `json:"name"`
	ID        string     `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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
func (c *Client) CreateConfigStore(i *CreateConfigStoreInput) (*ConfigStore, error) {
	path := "/resources/stores/config"
	resp, err := c.PostForm(path, i, &RequestOptions{
		Headers: map[string]string{
			// PostForm adds the appropriate Content-Type header.
			"Accept": "application/json",
		},
		Parallel: true,
	})
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
	// ID is the ID of the config store to delete (required).
	ID string
}

// DeleteConfigStore deletes the given config store version.
func (c *Client) DeleteConfigStore(i *DeleteConfigStoreInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/resources/stores/config/%s", i.ID)
	resp, err := c.Delete(path, &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Parallel: true,
	})
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
	// ID is the ID of the config store (required).
	ID string
}

// GetConfigStore returns the config store for the given input parameters.
func (c *Client) GetConfigStore(i *GetConfigStoreInput) (*ConfigStore, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/resources/stores/config/%s", i.ID)
	resp, err := c.Get(path, &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Parallel: true,
	})
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
	// ID is the ID of the config store (required).
	ID string
}

// GetConfigStoreMetadata returns the config store's metadata for the given input parameters.
func (c *Client) GetConfigStoreMetadata(i *GetConfigStoreMetadataInput) (*ConfigStoreMetadata, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/resources/stores/config/%s/info", i.ID)
	resp, err := c.Get(path, &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Parallel: true,
	})
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

// ListConfigStores returns a list of config stores sorted by name.
func (c *Client) ListConfigStores() ([]*ConfigStore, error) {
	path := "/resources/stores/config"
	resp, err := c.Get(path, &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Parallel: true,
	})
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
	// ID is the ID of the config store (required).
	ID string
}

// ListConfigStoreServices returns the list of services that are associated with
// a given config store.
func (c *Client) ListConfigStoreServices(i *ListConfigStoreServicesInput) ([]*Service, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/resources/stores/config/%s/services", i.ID)
	resp, err := c.Get(path, &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Parallel: true,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ss []*Service
	if err = json.NewDecoder(resp.Body).Decode(&ss); err != nil {
		return nil, err
	}

	byName := servicesByName(ss)
	sort.Sort(byName)

	return byName, nil
}

// UpdateConfigStoreInput is the input to UpdateConfigStore.
type UpdateConfigStoreInput struct {
	// ID is the ID of the config store to update (required).
	ID string

	// Name is the new name of the config store (required).
	Name string `url:"name"`
}

// UpdateConfigStore updates a specific config store.
func (c *Client) UpdateConfigStore(i *UpdateConfigStoreInput) (*ConfigStore, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/resources/stores/config/%s", i.ID)
	resp, err := c.PutForm(path, i, &RequestOptions{
		Headers: map[string]string{
			// PutForm adds the appropriate Content-Type header.
			"Accept": "application/json",
		},
		Parallel: true,
	})
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
