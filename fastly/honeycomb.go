package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Honeycomb represents a honeycomb response from the Fastly API.
type Honeycomb struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	Dataset           string     `mapstructure:"dataset"`
	Token             string     `mapstructure:"token"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// honeycombsByName is a sortable list of honeycombs.
type honeycombsByName []*Honeycomb

// Len, Swap, and Less implement the sortable interface.
func (h honeycombsByName) Len() int      { return len(h) }
func (h honeycombsByName) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h honeycombsByName) Less(i, j int) bool {
	return h[i].Name < h[j].Name
}

// ListHoneycombsInput is used as input to the ListHoneycombs function.
type ListHoneycombsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListHoneycombs returns the list of honeycombs for the configuration version.
func (c *Client) ListHoneycombs(i *ListHoneycombsInput) ([]*Honeycomb, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var hs []*Honeycomb
	if err := decodeBodyMap(resp.Body, &hs); err != nil {
		return nil, err
	}
	sort.Stable(honeycombsByName(hs))
	return hs, nil
}

// CreateHoneycombInput is used as input to the CreateHoneycomb function.
type CreateHoneycombInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `url:"name,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	Dataset           string `url:"dataset,omitempty"`
	Token             string `url:"token,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	Placement         string `url:"placement,omitempty"`
}

// CreateHoneycomb creates a new Fastly honeycomb.
func (c *Client) CreateHoneycomb(i *CreateHoneycombInput) (*Honeycomb, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Token == "" {
		return nil, ErrMissingToken
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetHoneycombInput is used as input to the GetHoneycomb function.
type GetHoneycombInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the honeycomb to fetch.
	Name string
}

// GetHoneycomb gets the honeycomb configuration with the given parameters.
func (c *Client) GetHoneycomb(i *GetHoneycombInput) (*Honeycomb, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// UpdateHoneycombInput is used as input to the UpdateHoneycomb function.
type UpdateHoneycombInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the honeycomb to update.
	Name string

	NewName           *string `url:"name,omitempty"`
	Format            *string `url:"format,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
	Dataset           *string `url:"dataset,omitempty"`
	Token             *string `url:"token,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	Placement         *string `url:"placement,omitempty"`
}

// UpdateHoneycomb updates a specific honeycomb.
func (c *Client) UpdateHoneycomb(i *UpdateHoneycombInput) (*Honeycomb, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	if i.Token != nil && *i.Token == "" {
		return nil, ErrTokenEmpty
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// DeleteHoneycombInput is the input parameter to DeleteHoneycomb.
type DeleteHoneycombInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the honeycomb to delete (required).
	Name string
}

// DeleteHoneycomb deletes the given honeycomb version.
func (c *Client) DeleteHoneycomb(i *DeleteHoneycombInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
