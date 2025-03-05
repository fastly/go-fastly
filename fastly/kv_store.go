package fastly

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// https://www.fastly.com/documentation/reference/api/services/resources/kv-store/
// https://www.fastly.com/documentation/reference/api/services/resources/kv-store-item/

// KVStore represents an KV Store response from the Fastly API.
type KVStore struct {
	CreatedAt *time.Time `mapstructure:"created_at"`
	Name      string     `mapstructure:"name"`
	StoreID   string     `mapstructure:"id"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
}

// CreateKVStoreInput is used as an input to the CreateKVStore function.
type CreateKVStoreInput struct {
	// Name is the name of the store to create (required).
	Name string `json:"name"`
	// Location is the regional location of the store (optional).
	Location string `json:"-"`
}

// CreateKVStore creates a new resource.
func (c *Client) CreateKVStore(i *CreateKVStoreInput) (*KVStore, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}

	ro := &RequestOptions{
		Params: map[string]string{},
	}
	if i.Location != "" {
		ro.Params["location"] = i.Location
	}

	const path = "/resources/stores/kv"
	resp, err := c.PostJSON(path, i, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var store *KVStore
	if err := DecodeBodyMap(resp.Body, &store); err != nil {
		return nil, err
	}
	return store, nil
}

// ListKVStoresInput is used as an input to the ListKVStores function.
type ListKVStoresInput struct {
	// Cursor is used for paginating through results.
	Cursor string
	// Limit is the maximum number of items included the response.
	Limit int
}

func (l *ListKVStoresInput) formatFilters() map[string]string {
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

// ListKVStoresResponse retrieves all resources.
type ListKVStoresResponse struct {
	// Data is the list of returned kv stores
	Data []KVStore
	// Meta is the information for pagination
	Meta map[string]string
}

// ListKVStores retrieves all resources.
func (c *Client) ListKVStores(i *ListKVStoresInput) (*ListKVStoresResponse, error) {
	const path = "/resources/stores/kv"

	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output *ListKVStoresResponse
	if err := DecodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// ListKVStoresPaginator is the opaque type for a ListKVStores call with pagination.
type ListKVStoresPaginator struct {
	client   *Client
	cursor   string // == "" if no more pages
	err      error
	finished bool
	input    *ListKVStoresInput
	stores   []KVStore // stored response from previous api call
}

// NewListKVStoresPaginator creates a new paginator for the given ListKVStoresInput.
func (c *Client) NewListKVStoresPaginator(i *ListKVStoresInput) *ListKVStoresPaginator {
	return &ListKVStoresPaginator{
		client: c,
		input:  i,
	}
}

// Next advances the paginator and fetches the next set of kv stores.
func (l *ListKVStoresPaginator) Next() bool {
	if l.finished {
		l.stores = nil
		return false
	}

	l.input.Cursor = l.cursor
	o, err := l.client.ListKVStores(l.input)
	if err != nil {
		l.err = err
		l.finished = true
		return false
	}

	l.stores = o.Data
	if next := o.Meta["next_cursor"]; next == "" {
		l.finished = true
	} else {
		l.cursor = next
	}

	return true
}

// Stores returns the current partial list of kv stores.
func (l *ListKVStoresPaginator) Stores() []KVStore {
	return l.stores
}

// Err returns any error from the pagination.
func (l *ListKVStoresPaginator) Err() error {
	return l.err
}

// GetKVStoreInput is the input to the GetKVStore function.
type GetKVStoreInput struct {
	// StoreID is the StoreID of the store to fetch (required).
	StoreID string
}

// GetKVStore retrieves the specified resource.
func (c *Client) GetKVStore(i *GetKVStoreInput) (*KVStore, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "kv", i.StoreID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output *KVStore
	if err := DecodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// DeleteKVStoreInput is the input to the DeleteKVStore function.
type DeleteKVStoreInput struct {
	// StoreID is the StoreID of the kv store to delete (required).
	StoreID string
}

// DeleteKVStore deletes the specified resource.
func (c *Client) DeleteKVStore(i *DeleteKVStoreInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "kv", i.StoreID)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}

// Consistency is a base for the different consistency variants.
type Consistency int64

func (c Consistency) String() string {
	switch c {
	case ConsistencyEventual:
		return "eventual"
	case ConsistencyUndefined, ConsistencyStrong:
		return "strong"
	}
	return "strong" // default
}

const (
	ConsistencyUndefined Consistency = iota
	ConsistencyEventual
	ConsistencyStrong
)

// ListKVStoreKeysInput is the input to the ListKVStoreKeys function.
type ListKVStoreKeysInput struct {
	// Consistency determines accuracy of results (values: eventual, strong). i.e. 'eventual' uses caching to improve performance (default: strong)
	Consistency Consistency
	// Prefix limits the results to keys which begin with the specified string.
	Prefix string
	// Cursor is used for paginating through results.
	Cursor string
	// StoreID is the StoreID of the kv store to list keys for (required).
	StoreID string
	// Limit is the maximum number of items included the response.
	Limit int
}

func (l *ListKVStoreKeysInput) formatFilters() map[string]string {
	if l == nil {
		return nil
	}

	m := make(map[string]string)
	m["consistency"] = l.Consistency.String()

	if l.Limit != 0 {
		m["limit"] = strconv.Itoa(l.Limit)
	}

	if l.Cursor != "" {
		m["cursor"] = l.Cursor
	}

	if l.Prefix != "" {
		m["prefix"] = l.Prefix
	}

	return m
}

// ListKVStoreKeysResponse retrieves all resources.
type ListKVStoreKeysResponse struct {
	// Data is the list of keys
	Data []string
	// Meta is the information for pagination
	Meta map[string]string
}

// ListKVStoreKeys retrieves all resources.
func (c *Client) ListKVStoreKeys(i *ListKVStoreKeysInput) (*ListKVStoreKeysResponse, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "kv", i.StoreID, "keys")

	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output *ListKVStoreKeysResponse
	if err := DecodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// ListKVStoreKeysPaginator is the opaque type for a ListKVStoreKeys calls with pagination.
type ListKVStoreKeysPaginator struct {
	client   *Client
	cursor   string // == "" if no more pages
	err      error
	finished bool
	input    *ListKVStoreKeysInput
	keys     []string // stored response from previous api call
}

// NewListKVStoreKeysPaginator returns a new paginator for the provided LitKVStoreKeysInput.
func (c *Client) NewListKVStoreKeysPaginator(i *ListKVStoreKeysInput) PaginatorKVStoreEntries {
	return &ListKVStoreKeysPaginator{
		client: c,
		input:  i,
	}
}

// Next advances the paginator.
func (l *ListKVStoreKeysPaginator) Next() bool {
	if l.finished {
		l.keys = nil
		return false
	}

	l.input.Cursor = l.cursor
	o, err := l.client.ListKVStoreKeys(l.input)
	if err != nil {
		l.err = err
		l.finished = true
		return false
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
func (l *ListKVStoreKeysPaginator) Err() error {
	return l.err
}

// Keys returns the current set of keys retrieved by the paginator.
func (l *ListKVStoreKeysPaginator) Keys() []string {
	return l.keys
}

// GetKVStoreKeyInput is the input to the GetKVStoreKey function.
type GetKVStoreKeyInput struct {
	// StoreID is the StoreID of the kv store (required).
	StoreID string
	// Key is the key to fetch (required).
	Key string
}

// GetKVStoreKey retrieves the specified resource.
func (c *Client) GetKVStoreKey(i *GetKVStoreKeyInput) (string, error) {
	if i.StoreID == "" {
		return "", ErrMissingStoreID
	}
	if i.Key == "" {
		return "", ErrMissingKey
	}

	path := ToSafeURL("resources", "stores", "kv", i.StoreID, "keys", i.Key)

	resp, err := c.Get(path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// GetKVStoreItemInput is the input to the GetKVStoreItem function.
type GetKVStoreItemInput struct {
	// StoreID is the StoreID of the kv store (required).
	StoreID string
	// Key is the key of the item to fetch (required).
	Key string
}

// GetKVStoreItemOutput is the output of the GetKVStoreItem function.
type GetKVStoreItemOutput struct {
	// Value is the value stored in the item.
	Value io.ReadCloser
	// Metadata is the metadata stored in the item, if any.
	Metadata string
	// Generation is the generation marker of the item.
	Generation uint64
}

func (o *GetKVStoreItemOutput) ValueAsBytes() ([]byte, error) {
	defer o.Value.Close()

	return io.ReadAll(o.Value)
}

func (o *GetKVStoreItemOutput) ValueAsString() (string, error) {
	if result, err := o.ValueAsBytes(); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

// GetKVStoreItem retrieves the specified item.
func (c *Client) GetKVStoreItem(i *GetKVStoreItemInput) (GetKVStoreItemOutput, error) {
	if i.StoreID == "" {
		return GetKVStoreItemOutput{}, ErrMissingStoreID
	}
	if i.Key == "" {
		return GetKVStoreItemOutput{}, ErrMissingKey
	}

	path := ToSafeURL("resources", "stores", "kv", i.StoreID, "keys", i.Key)

	resp, err := c.Get(path, nil)
	if err != nil {
		return GetKVStoreItemOutput{}, err
	}

	output := GetKVStoreItemOutput{Value: resp.Body}

	output.Metadata = resp.Header.Get("metadata")

	output.Generation, err = strconv.ParseUint(resp.Header.Get("generation"), 10, 64)
	if err != nil {
		resp.Body.Close()
		return GetKVStoreItemOutput{}, err
	}

	return output, nil
}

// LengthReader represents a type that can be read and exposes its length.
type LengthReader interface {
	io.Reader
	Len() int
}

// FileLengthReader allows an os.File type to be passed as a LengthReader to the
// InsertKVStoreKeyInput.Body field.
func FileLengthReader(f *os.File) (LengthReader, error) {
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &fileLenReader{
		f:   f,
		len: int(s.Size()),
	}, nil
}

type fileLenReader struct {
	f   *os.File
	len int
}

func (f *fileLenReader) Read(p []byte) (int, error) {
	return f.f.Read(p)
}

func (f *fileLenReader) Len() int {
	return f.len
}

// InsertKVStoreKeyInput is the input to the InsertKVStoreKey function.
type InsertKVStoreKeyInput struct {
	// Body is the value to insert and will be streamed to the endpoint.
	// This is for users who are passing very large files.
	// Otherwise use the 'Value' field instead.
	Body LengthReader
	// StoreID is the StoreID of the kv store (required).
	StoreID string
	// Key is the key to add (required).
	Key string
	// Value is the value to insert (ignored if Body is set).
	Value string
	// IfGenerationMatch specifies a 'generation marker' value
	// which must match the value on the specified key for the
	// deletion to proceed.
	IfGenerationMatch uint64
	// Add specifies that the operation must fail if the key
	// already exists.
	Add bool
	// Append specifies that the provided Body or Value will be
	// appended to the key's existing value, if any.
	Append bool
	// Prepend specifies that the provided Body or Value will be
	// prepended to the key's existing value, if any.
	Prepend bool
	// BackgroundFetch specifies that the new value for the key
	// does not need to be immediately visible to other users of
	// the store.
	BackgroundFetch bool
	// Metadata is a string which will be stored alongside the
	// key's value. This is specified as a pointer-to-string so
	// that existing metadata can be removed by specifying an
	// empty string.
	Metadata *string
	// TimeToLiveSec specifies the number of seconds (from the
	// completion of the insert/update operation) that the key
	// should be retrievable.
	TimeToLiveSec int
}

// InsertKVStoreKey inserts a key/value pair into an kv store.
func (c *Client) InsertKVStoreKey(i *InsertKVStoreKeyInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}
	if i.Key == "" {
		return ErrMissingKey
	}

	ro := RequestOptions{
		Parallel: true, // This will allow the Fastly CLI to make bulk inserts.
		Params:   map[string]string{},
		Headers:  map[string]string{},
	}

	if i.IfGenerationMatch != 0 {
		ro.Headers["if-generation-match"] = strconv.FormatUint(i.IfGenerationMatch, 10)
	}

	if i.Add {
		ro.Params["add"] = "true"
	}

	if i.Append {
		ro.Params["append"] = "true"
	}

	if i.Prepend {
		ro.Params["prepend"] = "true"
	}

	if i.BackgroundFetch {
		ro.Params["background_fetch"] = "true"
	}

	if i.Metadata != nil {
		ro.Headers["metadata"] = *i.Metadata
	}

	if i.TimeToLiveSec != 0 {
		ro.Headers["time_to_live_sec"] = strconv.Itoa(i.TimeToLiveSec)
	}

	if i.Body != nil {
		ro.Body = bufio.NewReader(i.Body)
		ro.BodyLength = int64(i.Body.Len())
	} else {
		ro.Body = strings.NewReader(i.Value)
		ro.BodyLength = int64(len(i.Value))
	}

	path := ToSafeURL("resources", "stores", "kv", i.StoreID, "keys", i.Key)

	resp, err := c.Put(path, &ro)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ignored, err := checkResp(resp, err)
	if err != nil {
		return err
	}
	defer ignored.Body.Close()
	return nil
}

// DeleteKVStoreKeyInput is the input to the DeleteKVStoreKey function.
type DeleteKVStoreKeyInput struct {
	// StoreID is the StoreID of the kv store (required).
	StoreID string
	// Key is the key to delete (required).
	Key string
	// Force is a flag to ignore a failure if the specified key
	// was not found.
	Force bool
	// IfGenerationMatch specifies a 'generation marker' value
	// which must match the value on the specified key for the
	// deletion to proceed.
	IfGenerationMatch uint64
}

// DeleteKVStoreKey deletes the specified resource.
func (c *Client) DeleteKVStoreKey(i *DeleteKVStoreKeyInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}
	if i.Key == "" {
		return ErrMissingKey
	}

	ro := RequestOptions{
		Parallel: true, // This will allow the Fastly CLI to make bulk deletes.
		Params:   map[string]string{},
		Headers:  map[string]string{},
	}

	if i.Force {
		ro.Params["force"] = "true"
	}

	if i.IfGenerationMatch != 0 {
		ro.Headers["if-generation-match"] = strconv.FormatUint(i.IfGenerationMatch, 10)
	}

	path := ToSafeURL("resources", "stores", "kv", i.StoreID, "keys", i.Key)

	resp, err := c.Delete(path, &ro)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}

// BatchModifyKVStoreKeyInput is the input to the BatchModifyKVStoreKey function.
type BatchModifyKVStoreKeyInput struct {
	// StoreID is the StoreID of the kv store (required).
	StoreID string
	// Body is the HTTP request body containing a collection of JSON objects
	// separated by a new line. {"key": "example","value": "<base64-encoded>"}
	// (required).
	Body io.Reader
}

// BatchModifyKVStoreKey streams key/value JSON objects into an kv store.
// NOTE: We wrap the io.Reader with *bufio.Reader to handle large streams.
func (c *Client) BatchModifyKVStoreKey(i *BatchModifyKVStoreKeyInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}

	path := ToSafeURL("resources", "stores", "kv", i.StoreID, "batch")

	resp, err := c.Put(path, &RequestOptions{
		Body: bufio.NewReader(i.Body),
		Headers: map[string]string{
			"Content-Type": "application/x-ndjson",
		},
		Parallel: true,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ignored, err := checkResp(resp, err)
	if err != nil {
		return err
	}
	defer ignored.Body.Close()
	return nil
}
