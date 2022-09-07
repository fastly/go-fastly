package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Heroku represents a heroku response from the Fastly API.
type Heroku struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListHerokus returns the list of herokus for the configuration version.
func (c *Client) ListHerokus(i *ListHerokusInput) ([]*Heroku, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku", i.ServiceID, i.ServiceVersion)
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `url:"name,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	URL               string `url:"url,omitempty"`
	Token             string `url:"token,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	Placement         string `url:"placement,omitempty"`
}

// CreateHeroku creates a new Fastly heroku.
func (c *Client) CreateHeroku(i *CreateHerokuInput) (*Heroku, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Heroku
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetHerokuInput is used as input to the GetHeroku function.
type GetHerokuInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the heroku to fetch.
	Name string
}

// GetHeroku gets the heroku configuration with the given parameters.
func (c *Client) GetHeroku(i *GetHerokuInput) (*Heroku, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Heroku
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// UpdateHerokuInput is used as input to the UpdateHeroku function.
type UpdateHerokuInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the heroku to update.
	Name string

	NewName           *string `url:"name,omitempty"`
	Format            *string `url:"format,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
	URL               *string `url:"url,omitempty"`
	Token             *string `url:"token,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	Placement         *string `url:"placement,omitempty"`
}

// UpdateHeroku updates a specific heroku.
func (c *Client) UpdateHeroku(i *UpdateHerokuInput) (*Heroku, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Heroku
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// DeleteHerokuInput is the input parameter to DeleteHeroku.
type DeleteHerokuInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the heroku to delete (required).
	Name string
}

// DeleteHeroku deletes the given heroku version.
func (c *Client) DeleteHeroku(i *DeleteHerokuInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/heroku/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
