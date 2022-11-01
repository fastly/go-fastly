package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Dictionary represents a dictionary response from the Fastly API.
type Dictionary struct {
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	ID             string     `mapstructure:"id"`
	Name           string     `mapstructure:"name"`
	ServiceID      string     `mapstructure:"service_id"`
	ServiceVersion int        `mapstructure:"version"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
	WriteOnly      bool       `mapstructure:"write_only"`
}

// dictionariesByName is a sortable list of dictionaries.
type dictionariesByName []*Dictionary

// Len implement the sortable interface.
func (s dictionariesByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s dictionariesByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s dictionariesByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListDictionariesInput is used as input to the ListDictionaries function.
type ListDictionariesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDictionaries returns the list of dictionaries for the configuration version.
func (c *Client) ListDictionaries(i *ListDictionariesInput) ([]*Dictionary, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/dictionary", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bs []*Dictionary
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	sort.Stable(dictionariesByName(bs))
	return bs, nil
}

// CreateDictionaryInput is used as input to the CreateDictionary function.
type CreateDictionaryInput struct {
	Name string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	WriteOnly      Compatibool `url:"write_only,omitempty"`
}

// CreateDictionary creates a new resource.
func (c *Client) CreateDictionary(i *CreateDictionaryInput) (*Dictionary, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/dictionary", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Dictionary
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetDictionaryInput is used as input to the GetDictionary function.
type GetDictionaryInput struct {
	// Name is the name of the dictionary to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetDictionary gets the dictionary configuration with the given parameters.
func (c *Client) GetDictionary(i *GetDictionaryInput) (*Dictionary, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/dictionary/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Dictionary
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateDictionaryInput is used as input to the UpdateDictionary function.
type UpdateDictionaryInput struct {
	// Name is the name of the dictionary to update.
	Name    string
	NewName *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	WriteOnly      *Compatibool `url:"write_only,omitempty"`
}

// UpdateDictionary updates a specific dictionary.
func (c *Client) UpdateDictionary(i *UpdateDictionaryInput) (*Dictionary, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/dictionary/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Dictionary
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteDictionaryInput is the input parameter to DeleteDictionary.
type DeleteDictionaryInput struct {
	// Name is the name of the dictionary to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteDictionary deletes the specified resource.
func (c *Client) DeleteDictionary(i *DeleteDictionaryInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/dictionary/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Unlike other endpoints, the dictionary endpoint does not return a status
	// response - it just returns a 200 OK.
	return nil
}
