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
	CreatedAt    *time.Time `mapstructure:"created_at"`
	DeletedAt    *time.Time `mapstructure:"deleted_at"`
	DictionaryID string     `mapstructure:"dictionary_id"`
	ItemKey      string     `mapstructure:"item_key"`
	ItemValue    string     `mapstructure:"item_value"`
	ServiceID    string     `mapstructure:"service_id"`
	UpdatedAt    *time.Time `mapstructure:"updated_at"`
}

// dictionaryItemsByKey is a sortable list of dictionary items.
type dictionaryItemsByKey []*DictionaryItem

// Len implement the sortable interface.
func (s dictionaryItemsByKey) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s dictionaryItemsByKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s dictionaryItemsByKey) Less(i, j int) bool {
	return s[i].ItemKey < s[j].ItemKey
}

// ListDictionaryItemsInput is used as input to the ListDictionaryItems function.
type ListDictionaryItemsInput struct {
	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string
	// Direction is the direction in which to sort results.
	Direction string
	// Page is the current page.
	Page int
	// PerPage is the number of records per page.
	PerPage int
	// ServiceID is the ID of the service (required).
	ServiceID string
	// Sort is the field on which to sort.
	Sort string
}

// ListDictionaryItems retrieves all resources.
func (c *Client) ListDictionaryItems(i *ListDictionaryItemsInput) ([]*DictionaryItem, error) {
	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/items", i.ServiceID, i.DictionaryID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bs []*DictionaryItem
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	sort.Stable(dictionaryItemsByKey(bs))
	return bs, nil
}

// ListDictionaryItemsPaginator implements the PaginatorDictionaryItems interface.
type ListDictionaryItemsPaginator struct {
	CurrentPage int
	LastPage    int
	NextPage    int

	// Private
	client   *Client
	consumed bool
	options  *ListDictionaryItemsInput
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

	// page is not specified, fetch from the beginning
	if i.Page <= 0 && p.CurrentPage == 0 {
		p.CurrentPage = 1
	} else {
		// page is specified, fetch from a given page
		if !p.consumed {
			p.CurrentPage = i.Page
		} else {
			p.CurrentPage = p.CurrentPage + 1
		}
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
	defer resp.Body.Close()

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
	// ItemKey is the dictionary item key, maximum 256 characters.
	ItemKey string `url:"item_key,omitempty"`
	// ItemValue is the dictionary item value, maximum 8000 characters.
	ItemValue string `url:"item_value,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string `url:"-"`
}

// CreateDictionaryItem creates a new resource.
func (c *Client) CreateDictionaryItem(i *CreateDictionaryItemInput) (*DictionaryItem, error) {
	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/item", i.ServiceID, i.DictionaryID)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DictionaryItem
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// CreateDictionaryItems creates a new resource.
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
	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string
	// ItemKey is the name of the dictionary item to fetch (required).
	ItemKey string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// GetDictionaryItem retrieves the specified resource.
func (c *Client) GetDictionaryItem(i *GetDictionaryItemInput) (*DictionaryItem, error) {
	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}
	if i.ItemKey == "" {
		return nil, ErrMissingItemKey
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/item/%s", i.ServiceID, i.DictionaryID, url.PathEscape(i.ItemKey))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DictionaryItem
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateDictionaryItemInput is used as input to the UpdateDictionaryItem function.
type UpdateDictionaryItemInput struct {
	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string
	// ItemKey is the name of the dictionary item to fetch (required).
	ItemKey string
	// ItemValue is the new value of the dictionary item.
	ItemValue string `url:"item_value,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// UpdateDictionaryItem updates the specified resource.
func (c *Client) UpdateDictionaryItem(i *UpdateDictionaryItemInput) (*DictionaryItem, error) {
	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}
	if i.ItemKey == "" {
		return nil, ErrMissingItemKey
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/item/%s", i.ServiceID, i.DictionaryID, url.PathEscape(i.ItemKey))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DictionaryItem
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// BatchModifyDictionaryItemsInput is the input parameter to the
// BatchModifyDictionaryItems function.
type BatchModifyDictionaryItemsInput struct {
	// DictionaryID is the ID of the dictionary to modify items for (required).
	DictionaryID string `json:"-"`
	// Items is a list of dictionary items.
	Items []*BatchDictionaryItem `json:"items"`
	// ServiceID is the ID of the service (required).
	ServiceID string `json:"-"`
}

// BatchDictionaryItem represents a dictionary item.
type BatchDictionaryItem struct {
	// ItemKey is an item key (maximum 256 characters).
	ItemKey string `json:"item_key"`
	// ItemValue is an item value (maximum 8000 characters).
	ItemValue string `json:"item_value"`
	// Operation is a batching operation variant.
	Operation BatchOperation `json:"op"`
}

// BatchModifyDictionaryItems bulk updates dictionary items.
func (c *Client) BatchModifyDictionaryItems(i *BatchModifyDictionaryItemsInput) error {
	if i.DictionaryID == "" {
		return ErrMissingDictionaryID
	}
	if len(i.Items) > BatchModifyMaximumOperations {
		return ErrMaxExceededItems
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/dictionary/%s/items", i.ServiceID, i.DictionaryID)
	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var batchModifyResult map[string]string
	return decodeBodyMap(resp.Body, &batchModifyResult)
}

// DeleteDictionaryItemInput is the input parameter to DeleteDictionaryItem.
type DeleteDictionaryItemInput struct {
	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string
	// ItemKey is the name of the dictionary item to delete (required).
	ItemKey string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// DeleteDictionaryItem deletes the specified resource.
func (c *Client) DeleteDictionaryItem(i *DeleteDictionaryItemInput) error {
	if i.DictionaryID == "" {
		return ErrMissingDictionaryID
	}
	if i.ItemKey == "" {
		return ErrMissingItemKey
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
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
