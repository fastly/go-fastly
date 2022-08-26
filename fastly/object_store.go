package fastly

import (
	"io"
	"net/http"
	"strconv"
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

type ListObjectStoresInput struct {
	Cursor string
	Limit  int
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

type ListObjectStoresResponse struct {
	Data []ObjectStore
	Meta map[string]string
}

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

type ListObjectStoresPaginator struct {
	cursor   string        // == "" if no more pages
	stores   []ObjectStore // stored response from previous api call
	finished bool
	err      error

	client *Client
	input  *ListObjectStoresInput
}

func (c *Client) NewListObjectStoresPaginator(i *ListObjectStoresInput) *ListObjectStoresPaginator {
	return &ListObjectStoresPaginator{
		client: c,
		input:  i,
	}
}

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

func (l *ListObjectStoresPaginator) Stores() []ObjectStore {
	return l.stores
}

func (l *ListObjectStoresPaginator) Err() error {
	return l.err
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
	ID     string
	Limit  int
	Cursor string
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

type ListObjectStoreKeysResponse struct {
	Data []string
	Meta map[string]string
}

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

type ListObjectsStoreKeysPaginator struct {
	cursor   string   // == "" if no more pages
	keys     []string // stored response from previous api call
	finished bool
	err      error

	client *Client
	input  *ListObjectStoreKeysInput
}

func (c *Client) NewListObjectStoreKeysPaginator(i *ListObjectStoreKeysInput) *ListObjectsStoreKeysPaginator {
	return &ListObjectsStoreKeysPaginator{
		client: c,
		input:  i,
	}
}

func (l *ListObjectsStoreKeysPaginator) Next() bool {
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

func (l *ListObjectsStoreKeysPaginator) Err() error {
	return l.err
}

func (l *ListObjectsStoreKeysPaginator) Keys() []string {
	return l.keys
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
