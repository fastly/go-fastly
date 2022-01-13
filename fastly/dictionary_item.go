package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/peterhellberg/link"
)

// DictionaryItem represents a dictionary item response from the Fastly API.
type DictionaryItem struct {
	ServiceID    string `mapstructure:"service_id"`
	DictionaryID string `mapstructure:"dictionary_id"`
	ItemKey      string `mapstructure:"item_key"`

	ItemValue string     `mapstructure:"item_value"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
}

// dictionaryItemsByKey is a sortable list of dictionary items.
type dictionaryItemsByKey []*DictionaryItem

// Len, Swap, and Less implement the sortable interface.
func (s dictionaryItemsByKey) Len() int      { return len(s) }
func (s dictionaryItemsByKey) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s dictionaryItemsByKey) Less(i, j int) bool {
	return s[i].ItemKey < s[j].ItemKey
}

// ListDictionaryItemsInput is used as input to the ListDictionaryItems function.
type ListDictionaryItemsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string
	Direction    string
	PerPage      int
	Page         int
	Sort         string
}

// ListDictionaryItems returns a list of items for a dictionary
func (c *Client) ListDictionaryItems(i *ListDictionaryItemsInput) ([]*DictionaryItem, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/items", i.ServiceID, i.DictionaryID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var bs []*DictionaryItem
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	sort.Stable(dictionaryItemsByKey(bs))
	return bs, nil
}

type ListDictionaryItemsPaginator struct {
	consumed    bool
	CurrentPage int
	NextPage    int
	LastPage    int
	client      *Client
	options     *ListDictionaryItemsInput
}

// HasNext returns a boolean indicating whether more pages are available
func (p *ListDictionaryItemsPaginator) HasNext() bool {
	return !p.consumed || p.Remaining() != 0
}

// Remaining returns the remaining page count
func (p *ListDictionaryItemsPaginator) Remaining() int {
	if p.LastPage == 0 {
		return 0
	}
	return p.LastPage - p.CurrentPage
}

// GetNext retrieves data in the next page
func (p *ListDictionaryItemsPaginator) GetNext() ([]*DictionaryItem, error) {
	return p.client.listDictionaryItemsWithPage(p.options, p)
}

// NewListDictionaryItemsPaginator returns a new paginator
func (c *Client) NewListDictionaryItemsPaginator(i *ListDictionaryItemsInput) PaginatorDictionaryItems {
	return &ListDictionaryItemsPaginator{
		client:  c,
		options: i,
	}
}

// listDictionaryItemsWithPage returns a list of items for a dictionary of a given page
func (c *Client) listDictionaryItemsWithPage(i *ListDictionaryItemsInput, p *ListDictionaryItemsPaginator) ([]*DictionaryItem, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}

	var perPage int
	const maxPerPage = 100
	if i.PerPage <= 0 {
		perPage = maxPerPage
	} else {
		perPage = i.PerPage
	}

	if i.Page <= 0 && p.CurrentPage == 0 {
		p.CurrentPage = 1
	} else {
		p.CurrentPage = p.CurrentPage + 1
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/items", i.ServiceID, i.DictionaryID)
	requestOptions := &RequestOptions{
		Params: map[string]string{
			"per_page": strconv.Itoa(perPage),
			"page":     strconv.Itoa(p.CurrentPage),
		},
	}

	if i.Direction != "" {
		requestOptions.Params["direction"] = i.Direction
	}
	if i.Sort != "" {
		requestOptions.Params["sort"] = i.Sort
	}

	resp, err := c.Get(path, requestOptions)
	if err != nil {
		return nil, err
	}

	for _, l := range link.ParseResponse(resp) {
		// indicates the Link response header contained the next page instruction
		if l.Rel == "next" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.NextPage, _ = strconv.Atoi(query["page"][0])
		}
		// indicates the Link response header contained the last page instruction
		if l.Rel == "last" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.LastPage, _ = strconv.Atoi(query["page"][0])
		}
	}

	p.consumed = true

	var bs []*DictionaryItem
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	sort.Stable(dictionaryItemsByKey(bs))

	return bs, nil
}

// CreateDictionaryItemInput is used as input to the CreateDictionaryItem function.
type CreateDictionaryItemInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string

	ItemKey   string `url:"item_key,omitempty"`
	ItemValue string `url:"item_value,omitempty"`
}

// CreateDictionaryItem creates a new Fastly dictionary item.
func (c *Client) CreateDictionaryItem(i *CreateDictionaryItemInput) (*DictionaryItem, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/item", i.ServiceID, i.DictionaryID)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *DictionaryItem
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// CreateDictionaryItems creates new Fastly dictionary items from a slice.
func (c *Client) CreateDictionaryItems(i []CreateDictionaryItemInput) ([]DictionaryItem, error) {

	var b []DictionaryItem
	for _, cdii := range i {
		di, err := c.CreateDictionaryItem(&cdii)
		if err != nil {
			return nil, err
		}
		b = append(b, *di)
	}
	return b, nil
}

// GetDictionaryItemInput is used as input to the GetDictionaryItem function.
type GetDictionaryItemInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string

	// ItemKey is the name of the dictionary item to fetch.
	ItemKey string
}

// GetDictionaryItem gets the dictionary item with the given parameters.
func (c *Client) GetDictionaryItem(i *GetDictionaryItemInput) (*DictionaryItem, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}

	if i.ItemKey == "" {
		return nil, ErrMissingItemKey
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/item/%s", i.ServiceID, i.DictionaryID, url.PathEscape(i.ItemKey))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var b *DictionaryItem
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateDictionaryItemInput is used as input to the UpdateDictionaryItem function.
type UpdateDictionaryItemInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string

	// ItemKey is the name of the dictionary item to fetch (required).
	ItemKey string

	// ItemValue is the new value of the dictionary item (required).
	ItemValue string `url:"item_value"`
}

// UpdateDictionaryItem updates a specific dictionary item.
func (c *Client) UpdateDictionaryItem(i *UpdateDictionaryItemInput) (*DictionaryItem, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}

	if i.ItemKey == "" {
		return nil, ErrMissingItemKey
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/item/%s", i.ServiceID, i.DictionaryID, url.PathEscape(i.ItemKey))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *DictionaryItem
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

type BatchModifyDictionaryItemsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string `json:"-"`

	// DictionaryID is the ID of the dictionary to modify items for (required).
	DictionaryID string `json:"-"`

	Items []*BatchDictionaryItem `json:"items"`
}

type BatchDictionaryItem struct {
	Operation BatchOperation `json:"op"`
	ItemKey   string         `json:"item_key"`
	ItemValue string         `json:"item_value"`
}

func (c *Client) BatchModifyDictionaryItems(i *BatchModifyDictionaryItemsInput) error {

	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.DictionaryID == "" {
		return ErrMissingDictionaryID
	}

	if len(i.Items) > BatchModifyMaximumOperations {
		return ErrMaxExceededItems
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/items", i.ServiceID, i.DictionaryID)
	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return err
	}

	var batchModifyResult map[string]string
	if err := decodeBodyMap(resp.Body, &batchModifyResult); err != nil {
		return err
	}

	return nil
}

// DeleteDictionaryItemInput is the input parameter to DeleteDictionaryItem.
type DeleteDictionaryItemInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string

	// ItemKey is the name of the dictionary item to delete.
	ItemKey string
}

// DeleteDictionaryItem deletes the given dictionary item.
func (c *Client) DeleteDictionaryItem(i *DeleteDictionaryItemInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.DictionaryID == "" {
		return ErrMissingDictionaryID
	}

	if i.ItemKey == "" {
		return ErrMissingItemKey
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/item/%s", i.ServiceID, i.DictionaryID, url.PathEscape(i.ItemKey))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Unlike other endpoints, the dictionary endpoint does not return a status
	// response - it just returns a 200 OK.
	return nil
}
