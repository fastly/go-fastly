package fastly

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://developer.fastly.com/reference/api/object-store/

// ObjectStore represents an Object Store response from the Fastly API.
type ObjectStore struct {
	CreatedAt *time.Time `mapstructure:"created_at"`
	ID        string     `mapstructure:"id"`
	Name      string     `mapstructure:"name"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
}

// CreateObjectStoreInput is used as an input to the CreateObjectStore function.
type CreateObjectStoreInput struct {
	// Name is the name of the store to create (required).
	Name string `json:"name"`
}

// CreateObjectStore creates a new resource.
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

// ListObjectStoresInput is used as an input to the ListObjectStores function.
type ListObjectStoresInput struct {
	// Cursor is used for paginating through results.
	Cursor string
	// Limit is the maximum number of items included the response.
	Limit int
}

func (l *ListObjectStoresInput) formatFilters() map[string]string {
	if l == nil {
		return nil
	}

	if l.Limit == 0 && l.Cursor == "" {
		return nil
	}

	m := make(map[string]string)

	if l.Limit != 0 {
		m["limit"] = strconv.Itoa(l.Limit)
	}

	if l.Cursor != "" {
		m["cursor"] = l.Cursor
	}

	return m
}

// ListObjectStoresResponse retrieves all resources.
type ListObjectStoresResponse struct {
	// Data is the list of returned object stores
	Data []ObjectStore
	// Meta is the information for pagination
	Meta map[string]string
}

// ListObjectStores retrieves all resources.
func (c *Client) ListObjectStores(i *ListObjectStoresInput) (*ListObjectStoresResponse, error) {
	const path = "/resources/stores/object"

	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}

	var output *ListObjectStoresResponse
	if err := decodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// ListObjectStoresPaginator is the opaque type for a ListObjectStores call with pagination.
type ListObjectStoresPaginator struct {
	client   *Client
	cursor   string // == "" if no more pages
	err      error
	finished bool
	input    *ListObjectStoresInput
	stores   []ObjectStore // stored response from previous api call
}

// NewListObjectStoresPaginator creates a new paginator for the given ListObjectStoresInput.
func (c *Client) NewListObjectStoresPaginator(i *ListObjectStoresInput) *ListObjectStoresPaginator {
	return &ListObjectStoresPaginator{
		client: c,
		input:  i,
	}
}

// Next advances the paginator and fetches the next set of object stores.
func (l *ListObjectStoresPaginator) Next() bool {
	if l.finished {
		l.stores = nil
		return false
	}

	l.input.Cursor = l.cursor
	o, err := l.client.ListObjectStores(l.input)
	if err != nil {
		l.err = err
		l.finished = true
	}

	l.stores = o.Data
	if next := o.Meta["next_cursor"]; next == "" {
		l.finished = true
	} else {
		l.cursor = next
	}

	return true
}

// Stores returns the current partial list of object stores.
func (l *ListObjectStoresPaginator) Stores() []ObjectStore {
	return l.stores
}

// Err returns any error from the pagination.
func (l *ListObjectStoresPaginator) Err() error {
	return l.err
}

// GetObjectStoreInput is the input to the GetObjectStore function.
type GetObjectStoreInput struct {
	// ID is the ID of the store to fetch (required).
	ID string
}

// GetObjectStore retrieves the specified resource.
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

// DeleteObjectStoreInput is the input to the DeleteObjectStore function.
type DeleteObjectStoreInput struct {
	// ID is the ID of the object store to delete (required).
	ID string
}

// DeleteObjectStore deletes the specified resource.
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

// ListObjectStoreKeysInput is the input to the ListObjectStoreKeys function.
type ListObjectStoreKeysInput struct {
	// Cursor is used for paginating through results.
	Cursor string
	// ID is the ID of the object store to list keys for (required).
	ID string
	// Limit is the maximum number of items included the response.
	Limit int
}

func (l *ListObjectStoreKeysInput) formatFilters() map[string]string {
	if l == nil {
		return nil
	}

	if l.Limit == 0 && l.Cursor == "" {
		return nil
	}

	m := make(map[string]string)

	if l.Limit != 0 {
		m["limit"] = strconv.Itoa(l.Limit)
	}

	if l.Cursor != "" {
		m["cursor"] = l.Cursor
	}

	return m
}

// ListObjectStoreKeysResponse retrieves all resources.
type ListObjectStoreKeysResponse struct {
	// Data is the list of keys
	Data []string
	// Meta is the information for pagination
	Meta map[string]string
}

// ListObjectStoreKeys retrieves all resources.
func (c *Client) ListObjectStoreKeys(i *ListObjectStoreKeysInput) (*ListObjectStoreKeysResponse, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := "/resources/stores/object/" + i.ID + "/keys"
	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}

	var output *ListObjectStoreKeysResponse
	if err := decodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// ListObjectStoreKeysPaginator is the opaque type for a ListObjectStoreKeys calls with pagination.
type ListObjectStoreKeysPaginator struct {
	client   *Client
	cursor   string // == "" if no more pages
	err      error
	finished bool
	input    *ListObjectStoreKeysInput
	keys     []string // stored response from previous api call
}

// NewListObjectStoreKeysPaginator returns a new paginator for the provided LitObjectStoreKeysInput.
func (c *Client) NewListObjectStoreKeysPaginator(i *ListObjectStoreKeysInput) *ListObjectStoreKeysPaginator {
	return &ListObjectStoreKeysPaginator{
		client: c,
		input:  i,
	}
}

// Next advanced the paginator.
func (l *ListObjectStoreKeysPaginator) Next() bool {
	if l.finished {
		l.keys = nil
		return false
	}

	l.input.Cursor = l.cursor
	o, err := l.client.ListObjectStoreKeys(l.input)
	if err != nil {
		l.err = err
		l.finished = true
	}

	l.keys = o.Data
	if next := o.Meta["next_cursor"]; next == "" {
		l.finished = true
	} else {
		l.cursor = next
	}

	return true
}

// Err returns any error from the paginator.
func (l *ListObjectStoreKeysPaginator) Err() error {
	return l.err
}

// Keys returns the current set of keys retrieved by the paginator.
func (l *ListObjectStoreKeysPaginator) Keys() []string {
	return l.keys
}

// GetObjectStoreKeyInput is the input to the GetObjectStoreKey function.
type GetObjectStoreKeyInput struct {
	// ID is the ID of the object store (required).
	ID string
	// Key is the key to fetch (required).
	Key string
}

// GetObjectStoreKey retrieves the specified resource.
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

// InsertObjectStoreKeyInput is the input to the InsertObjectStoreKey function.
type InsertObjectStoreKeyInput struct {
	// ID is the ID of the object store (required).
	ID string
	// Key is the key to add (required).
	Key string
	// Value is the value to insert (required).
	Value string
}

// InsertObjectStoreKey inserts a key/value pair into an object store.
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

// DeleteObjectStoreKeyInput is the input to the DeleteObjectStoreKey function.
type DeleteObjectStoreKeyInput struct {
	// ID is the ID of the object store (required).
	ID string
	// Key is the key to delete (required).
	Key string
}

// DeleteObjectStoreKey deletes the specified resource.
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
