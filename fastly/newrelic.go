package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// NewRelic represents a newrelic response from the Fastly API.
type NewRelic struct {
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

// newrelicByName is a sortable list of newrelic.
type newrelicByName []*NewRelic

// Len, Swap, and Less implement the sortable interface.
func (s newrelicByName) Len() int      { return len(s) }
func (s newrelicByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s newrelicByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListNewRelicInput is used as input to the ListNewRelic function.
type ListNewRelicInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListNewRelic returns the list of newrelic for the configuration version.
func (c *Client) ListNewRelic(i *ListNewRelicInput) ([]*NewRelic, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var n []*NewRelic
	if err := decodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	sort.Stable(newrelicByName(n))
	return n, nil
}

// CreateNewRelicInput is used as input to the CreateNewRelic function.
type CreateNewRelicInput struct {
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

// CreateNewRelic creates a new Fastly newrelic.
func (c *Client) CreateNewRelic(i *CreateNewRelicInput) (*NewRelic, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var n *NewRelic
	if err := decodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// GetNewRelicInput is used as input to the GetNewRelic function.
type GetNewRelicInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the newrelic to fetch.
	Name string
}

// GetNewRelic gets the newrelic configuration with the given parameters.
func (c *Client) GetNewRelic(i *GetNewRelicInput) (*NewRelic, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var n *NewRelic
	if err := decodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// UpdateNewRelicInput is used as input to the UpdateNewRelic function.
type UpdateNewRelicInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the newrelic to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Token             *string `form:"token,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateNewRelic updates a specific newrelic.
func (c *Client) UpdateNewRelic(i *UpdateNewRelicInput) (*NewRelic, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var n *NewRelic
	if err := decodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// DeleteNewRelicInput is the input parameter to DeleteNewRelic.
type DeleteNewRelicInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the newrelic to delete (required).
	Name string
}

// DeleteNewRelic deletes the given newrelic version.
func (c *Client) DeleteNewRelic(i *DeleteNewRelicInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic/%s", i.Service, i.Version, url.PathEscape(i.Name))
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
