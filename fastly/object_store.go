package fastly

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// https://developer.fastly.com/reference/api/object-store/

type ObjectStore struct {
	Name      string     `mapstructure:"name"`
	ID        string     `mapstructure:"id"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
}

type CreateObjectStoreInput struct {
	Name string `json:"name"`
}

func (c *Client) CreateObjectStore(i *CreateObjectStoreInput) (*ObjectStore, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}

	const path = "/resources/stores/object"
	resp, err := c.PostJSON(path, i, nil)
	if err != nil {
		return nil, err
	}

	var store *ObjectStore
	if err := decodeBodyMap(resp.Body, &store); err != nil {
		return nil, err
	}
	return store, nil
}

type ListObjectStoresResponse struct {
	Data []ObjectStore
	Meta map[string]string
}

func (c *Client) ListObjectStores() (*ListObjectStoresResponse, error) {
	const path = "/resources/stores/object"
	// TODO(dgryski): How do I handle pagination?
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var output *ListObjectStoresResponse
	if err := decodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

type GetObjectStoreInput struct {
	ID string
}

func (c *Client) GetObjectStore(i *GetObjectStoreInput) (*ObjectStore, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := "/resources/stores/object/" + i.ID
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var output *ObjectStore
	if err := decodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

type DeleteObjectStoreInput struct {
	ID string
}

func (c *Client) DeleteObjectStore(i *DeleteObjectStoreInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := "/resources/stores/object/" + i.ID
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}

type ListObjectStoreKeysInput struct {
	ID string
}

type ListObjectStoreKeysResponse struct {
	Data []string
	Meta map[string]string
}

func (c *Client) ListObjectStoreKeys(i *ListObjectStoreKeysInput) ([]string, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := "/resources/stores/object/" + i.ID + "/keys"
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var output *ListObjectStoreKeysResponse
	if err := decodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output.Data, nil
}

// missing example requests from API reference

type GetObjectStoreKeyInput struct {
	ID  string
	Key string
}

func (c *Client) GetObjectStoreKey(i *GetObjectStoreKeyInput) (string, error) {
	if i.ID == "" {
		return "", ErrMissingID
	}

	if i.Key == "" {
		return "", ErrMissingKey
	}

	path := "/resources/stores/object/" + i.ID + "/keys/" + i.Key
	resp, err := c.Get(path, nil)
	if err != nil {
		return "", err
	}

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

type InsertObjectStoreKeyInput struct {
	ID    string
	Key   string
	Value string // TODO(dgryski): io.Reader or string?
}

func (c *Client) InsertObjectStoreKey(i *InsertObjectStoreKeyInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	if i.Key == "" {
		return ErrMissingKey
	}

	path := "/resources/stores/object/" + i.ID + "/keys/" + i.Key
	resp, err := c.Put(path, &RequestOptions{Body: io.NopCloser(strings.NewReader(i.Value))})
	if err != nil {
		return err
	}

	_, err = checkResp(resp, err)
	return err
}

type DeleteObjectStoreKeyInput struct {
	ID  string
	Key string
}

func (c *Client) DeleteObjectStoreKey(i *DeleteObjectStoreKeyInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	if i.Key == "" {
		return ErrMissingKey
	}

	path := "/resources/stores/object/" + i.ID + "/keys/" + i.Key
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}
