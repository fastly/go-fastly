package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Heroku represents a heroku response from the Fastly API.
type Heroku struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	URL               string     `mapstructure:"url"`
	Token             string     `mapstructure:"token"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// herokusByName is a sortable list of herokus.
type herokusByName []*Heroku

// Len, Swap, and Less implement the sortable interface.
func (h herokusByName) Len() int      { return len(h) }
func (h herokusByName) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h herokusByName) Less(i, j int) bool {
	return h[i].Name < h[j].Name
}

// ListHerokusInput is used as input to the ListHerokus function.
type ListHerokusInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListHerokus returns the list of herokus for the configuration version.
func (c *Client) ListHerokus(i *ListHerokusInput) ([]*Heroku, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var hs []*Heroku
	if err := decodeBodyMap(resp.Body, &hs); err != nil {
		return nil, err
	}
	sort.Stable(herokusByName(hs))
	return hs, nil
}

// CreateHerokuInput is used as input to the CreateHeroku function.
type CreateHerokuInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	Name              *string `form:"name,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	URL               *string `form:"url,omitempty"`
	Token             *string `form:"token,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// CreateHeroku creates a new Fastly heroku.
func (c *Client) CreateHeroku(i *CreateHerokuInput) (*Heroku, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var h *Heroku
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetHerokuInput is used as input to the GetHeroku function.
type GetHerokuInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the heroku to fetch.
	Name string
}

// GetHeroku gets the heroku configuration with the given parameters.
func (c *Client) GetHeroku(i *GetHerokuInput) (*Heroku, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var h *Heroku
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// UpdateHerokuInput is used as input to the UpdateHeroku function.
type UpdateHerokuInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the heroku to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	URL               *string `form:"url,omitempty"`
	Token             *string `form:"token,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateHeroku updates a specific heroku.
func (c *Client) UpdateHeroku(i *UpdateHerokuInput) (*Heroku, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var h *Heroku
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// DeleteHerokuInput is the input parameter to DeleteHeroku.
type DeleteHerokuInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the heroku to delete (required).
	Name string
}

// DeleteHeroku deletes the given heroku version.
func (c *Client) DeleteHeroku(i *DeleteHerokuInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
