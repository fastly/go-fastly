package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Scalyr represents a scalyr response from the Fastly API.
type Scalyr struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	Token             string     `mapstructure:"token"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// scalyrByName is a sortable list of scalyrs.
type scalyrsByName []*Scalyr

// Len, Swap, and Less implement the sortable interface.
func (s scalyrsByName) Len() int      { return len(s) }
func (s scalyrsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s scalyrsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListScalyrsInput is used as input to the ListScalyrs function.
type ListScalyrsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListScalyrs returns the list of scalyrs for the configuration version.
func (c *Client) ListScalyrs(i *ListScalyrsInput) ([]*Scalyr, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ss []*Scalyr
	if err := decodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	sort.Stable(scalyrsByName(ss))
	return ss, nil
}

// CreateScalyrInput is used as input to the CreateScalyr function.
type CreateScalyrInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	Name              *string `form:"name,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	Token             *string `form:"token,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// CreateScalyr creates a new Fastly scalyr.
func (c *Client) CreateScalyr(i *CreateScalyrInput) (*Scalyr, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Scalyr
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetScalyrInput is used as input to the GetScalyr function.
type GetScalyrInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the scalyr to fetch.
	Name string
}

// GetScalyr gets the scalyr configuration with the given parameters.
func (c *Client) GetScalyr(i *GetScalyrInput) (*Scalyr, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var s *Scalyr
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateScalyrInput is used as input to the UpdateScalyr function.
type UpdateScalyrInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the scalyr to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	Token             *string `form:"token,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateScalyr updates a specific scalyr.
func (c *Client) UpdateScalyr(i *UpdateScalyrInput) (*Scalyr, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Scalyr
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteScalyrInput is the input parameter to DeleteScalyr.
type DeleteScalyrInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the scalyr to delete (required).
	Name string
}

// DeleteScalyr deletes the given scalyr version.
func (c *Client) DeleteScalyr(i *DeleteScalyrInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr/%s", i.Service, i.Version, url.PathEscape(i.Name))
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
