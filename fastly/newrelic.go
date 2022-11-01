package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// NewRelic represents a newrelic response from the Fastly API.
type NewRelic struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	Name              string     `mapstructure:"name"`
	Placement         string     `mapstructure:"placement"`
	Region            string     `mapstructure:"region"`
	ResponseCondition string     `mapstructure:"response_condition"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	Token             string     `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// newrelicByName is a sortable list of newrelic.
type newrelicByName []*NewRelic

// Len implement the sortable interface.
func (s newrelicByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s newrelicByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s newrelicByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListNewRelicInput is used as input to the ListNewRelic function.
type ListNewRelicInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListNewRelic returns the list of newrelic for the configuration version.
func (c *Client) ListNewRelic(i *ListNewRelicInput) ([]*NewRelic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var n []*NewRelic
	if err := decodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	sort.Stable(newrelicByName(n))
	return n, nil
}

// CreateNewRelicInput is used as input to the CreateNewRelic function.
type CreateNewRelicInput struct {
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	Name              string `url:"name,omitempty"`
	Placement         string `url:"placement,omitempty"`
	Region            string `url:"region,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Token          string `url:"token,omitempty"`
}

// CreateNewRelic creates a new resource.
func (c *Client) CreateNewRelic(i *CreateNewRelicInput) (*NewRelic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Token == "" {
		return nil, ErrMissingToken
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var n *NewRelic
	if err := decodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// GetNewRelicInput is used as input to the GetNewRelic function.
type GetNewRelicInput struct {
	// Name is the name of the newrelic to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetNewRelic gets the newrelic configuration with the given parameters.
func (c *Client) GetNewRelic(i *GetNewRelicInput) (*NewRelic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var n *NewRelic
	if err := decodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// UpdateNewRelicInput is used as input to the UpdateNewRelic function.
type UpdateNewRelicInput struct {
	Format        *string `url:"format,omitempty"`
	FormatVersion *uint   `url:"format_version,omitempty"`
	// Name is the name of the newrelic to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	Region            *string `url:"region,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Token          *string `url:"token,omitempty"`
}

// UpdateNewRelic updates a specific newrelic.
func (c *Client) UpdateNewRelic(i *UpdateNewRelicInput) (*NewRelic, error) {
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

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var n *NewRelic
	if err := decodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// DeleteNewRelicInput is the input parameter to DeleteNewRelic.
type DeleteNewRelicInput struct {
	// Name is the name of the newrelic to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteNewRelic deletes the specified resource.
func (c *Client) DeleteNewRelic(i *DeleteNewRelicInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/newrelic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
