package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Loggly represents a loggly response from the Fastly API.
type Loggly struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Token             string     `mapstructure:"token"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// logglyByName is a sortable list of loggly.
type logglyByName []*Loggly

// Len, Swap, and Less implement the sortable interface.
func (s logglyByName) Len() int      { return len(s) }
func (s logglyByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s logglyByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListLogglyInput is used as input to the ListLoggly function.
type ListLogglyInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListLoggly returns the list of loggly for the configuration version.
func (c *Client) ListLoggly(i *ListLogglyInput) ([]*Loggly, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/loggly", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ls []*Loggly
	if err := decodeBodyMap(resp.Body, &ls); err != nil {
		return nil, err
	}
	sort.Stable(logglyByName(ls))
	return ls, nil
}

// CreateLogglyInput is used as input to the CreateLoggly function.
type CreateLogglyInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	Name              *string `form:"name,omitempty"`
	Token             *string `form:"token,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// CreateLoggly creates a new Fastly loggly.
func (c *Client) CreateLoggly(i *CreateLogglyInput) (*Loggly, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/loggly", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var l *Loggly
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// GetLogglyInput is used as input to the GetLoggly function.
type GetLogglyInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the loggly to fetch.
	Name string
}

// GetLoggly gets the loggly configuration with the given parameters.
func (c *Client) GetLoggly(i *GetLogglyInput) (*Loggly, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/loggly/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var l *Loggly
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// UpdateLogglyInput is used as input to the UpdateLoggly function.
type UpdateLogglyInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the loggly to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Token             *string `form:"token,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateLoggly updates a specific loggly.
func (c *Client) UpdateLoggly(i *UpdateLogglyInput) (*Loggly, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/loggly/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var l *Loggly
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// DeleteLogglyInput is the input parameter to DeleteLoggly.
type DeleteLogglyInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the loggly to delete (required).
	Name string
}

// DeleteLoggly deletes the given loggly version.
func (c *Client) DeleteLoggly(i *DeleteLogglyInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/loggly/%s", i.Service, i.Version, url.PathEscape(i.Name))
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
