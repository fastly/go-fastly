package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Sumologic represents a sumologic response from the Fastly API.
type Sumologic struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Address           string     `mapstructure:"address"`
	URL               string     `mapstructure:"url"`
	Format            string     `mapstructure:"format"`
	ResponseCondition string     `mapstructure:"response_condition"`
	MessageType       string     `mapstructure:"message_type"`
	FormatVersion     int        `mapstructure:"format_version"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Placement         string     `mapstructure:"placement"`
}

// sumologicsByName is a sortable list of sumologics.
type sumologicsByName []*Sumologic

// Len, Swap, and Less implement the sortable interface.
func (s sumologicsByName) Len() int      { return len(s) }
func (s sumologicsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sumologicsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListSumologicsInput is used as input to the ListSumologics function.
type ListSumologicsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListSumologics returns the list of sumologics for the configuration version.
func (c *Client) ListSumologics(i *ListSumologicsInput) ([]*Sumologic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ss []*Sumologic
	if err := decodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	sort.Stable(sumologicsByName(ss))
	return ss, nil
}

// CreateSumologicInput is used as input to the CreateSumologic function.
type CreateSumologicInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `form:"name,omitempty"`
	Address           string `form:"address,omitempty"`
	URL               string `form:"url,omitempty"`
	Format            string `form:"format,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	MessageType       string `form:"message_type,omitempty"`
	FormatVersion     int    `form:"format_version,omitempty"`
	Placement         string `form:"placement,omitempty"`
}

// CreateSumologic creates a new Fastly sumologic.
func (c *Client) CreateSumologic(i *CreateSumologicInput) (*Sumologic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Sumologic
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSumologicInput is used as input to the GetSumologic function.
type GetSumologicInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the sumologic to fetch.
	Name string
}

// GetSumologic gets the sumologic configuration with the given parameters.
func (c *Client) GetSumologic(i *GetSumologicInput) (*Sumologic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var s *Sumologic
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateSumologicInput is used as input to the UpdateSumologic function.
type UpdateSumologicInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the sumologic to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Address           *string `form:"address,omitempty"`
	URL               *string `form:"url,omitempty"`
	Format            *string `form:"format,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	MessageType       *string `form:"message_type,omitempty"`
	FormatVersion     *int    `form:"format_version,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateSumologic updates a specific sumologic.
func (c *Client) UpdateSumologic(i *UpdateSumologicInput) (*Sumologic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Sumologic
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSumologicInput is the input parameter to DeleteSumologic.
type DeleteSumologicInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the sumologic to delete (required).
	Name string
}

// DeleteSumologic deletes the given sumologic version.
func (c *Client) DeleteSumologic(i *DeleteSumologicInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
