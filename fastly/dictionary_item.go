package fastly

import (
	"fmt"
	"time"
)

// DictionaryItem represents a dictionary item response from the Fastly API.
type DictionaryItem struct {
	CreatedAt    *time.Time `mapstructure:"created_at"`
	DeletedAt    *time.Time `mapstructure:"deleted_at"`
	DictionaryID *string    `mapstructure:"dictionary_id"`
	ItemKey      *string    `mapstructure:"item_key"`
	ItemValue    *string    `mapstructure:"item_value"`
	ServiceID    *string    `mapstructure:"service_id"`
	UpdatedAt    *time.Time `mapstructure:"updated_at"`
}

// GetDictionaryItemsInput is used as input to the GetDictionaryItems function.
type GetDictionaryItemsInput struct {
	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string
	// Direction is the direction in which to sort results.
	Direction *string
	// Page is the current page.
	Page *int
	// PerPage is the number of records per page.
	PerPage *int
	// ServiceID is the ID of the service (required).
	ServiceID string
	// Sort is the field on which to sort.
	Sort *string
}

// GetDictionaryItems returns a ListPaginator for paginating through the resources.
func (c *Client) GetDictionaryItems(i *GetDictionaryItemsInput) *ListPaginator[DictionaryItem] {
	input := ListOpts{}
	if i.Direction != nil {
		input.Direction = *i.Direction
	}
	if i.Sort != nil {
		input.Sort = *i.Sort
	}
	if i.Page != nil {
		input.Page = *i.Page
	}
	if i.PerPage != nil {
		input.PerPage = *i.PerPage
	}

	path := ToSafeURL("service", i.ServiceID, "dictionary", i.DictionaryID, "items")

	return NewPaginator[DictionaryItem](c, input, path)
}

// ListDictionaryItemsInput is used as input to the ListDictionaryItems function.
type ListDictionaryItemsInput struct {
	// DictionaryID is the ID of the dictionary to retrieve items for (required).
	DictionaryID string
	// Direction is the direction in which to sort results.
	Direction *string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// Sort is the field on which to sort.
	Sort *string
}

// ListDictionaryItems retrieves all resources. Not suitable for large
// collections.
func (c *Client) ListDictionaryItems(i *ListDictionaryItemsInput) ([]*DictionaryItem, error) {
	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	p := c.GetDictionaryItems(&GetDictionaryItemsInput{
		DictionaryID: i.DictionaryID,
		Direction:    i.Direction,
		ServiceID:    i.ServiceID,
		Sort:         i.Sort,
	})
	var results []*DictionaryItem
	for p.HasNext() {
		data, err := p.GetNext()
		if err != nil {
			return nil, fmt.Errorf("failed to get next page (remaining: %d): %s", p.Remaining(), err)
		}
		results = append(results, data...)
	}
	return results, nil
}

// CreateDictionaryItemInput is used as input to the CreateDictionaryItem function.
type CreateDictionaryItemInput struct {
	// ItemKey is the dictionary item key, maximum 256 characters.
	ItemKey *string `url:"item_key,omitempty"`
	// ItemValue is the dictionary item value, maximum 8000 characters.
	ItemValue *string `url:"item_value,omitempty"`
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

	path := ToSafeURL("service", i.ServiceID, "dictionary", i.DictionaryID, "item")

	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DictionaryItem
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// CreateDictionaryItems creates a new resource.
func (c *Client) CreateDictionaryItems(i []CreateDictionaryItemInput) ([]DictionaryItem, error) {
	var b []DictionaryItem
	for _, cdii := range i {
		cdii := cdii // it's unlikely the underlying value will have changed but we avoid a gosec warning this way (ref: https://bit.ly/go-range-bug)
		ptr, err := c.CreateDictionaryItem(&cdii)
		if err != nil {
			return nil, err
		}
		if ptr == nil {
			return nil, fmt.Errorf("error: unexpected nil pointer")
		}
		b = append(b, *ptr)
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

	path := ToSafeURL("service", i.ServiceID, "dictionary", i.DictionaryID, "item", i.ItemKey)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DictionaryItem
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
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

	path := ToSafeURL("service", i.ServiceID, "dictionary", i.DictionaryID, "item", i.ItemKey)

	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DictionaryItem
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
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
	ItemKey *string `json:"item_key"`
	// ItemValue is an item value (maximum 8000 characters).
	ItemValue *string `json:"item_value"`
	// Operation is a batching operation variant.
	Operation *BatchOperation `json:"op"`
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

	path := ToSafeURL("service", i.ServiceID, "dictionary", i.DictionaryID, "items")

	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var batchModifyResult map[string]string
	return DecodeBodyMap(resp.Body, &batchModifyResult)
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

	path := ToSafeURL("service", i.ServiceID, "dictionary", i.DictionaryID, "item", i.ItemKey)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Unlike other endpoints, the dictionary endpoint does not return a status
	// response - it just returns a 200 OK.
	return nil
}
