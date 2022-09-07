package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Logshuttle represents a logshuttle response from the Fastly API.
type Logshuttle struct {
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

// logshuttlesByName is a sortable list of logshuttles.
type logshuttlesByName []*Logshuttle

// Len, Swap, and Less implement the sortable interface.
func (l logshuttlesByName) Len() int      { return len(l) }
func (l logshuttlesByName) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l logshuttlesByName) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

// ListLogshuttlesInput is used as input to the ListLogshuttles function.
type ListLogshuttlesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListLogshuttles returns the list of logshuttles for the configuration version.
func (c *Client) ListLogshuttles(i *ListLogshuttlesInput) ([]*Logshuttle, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logshuttle", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ls []*Logshuttle
	if err := decodeBodyMap(resp.Body, &ls); err != nil {
		return nil, err
	}
	sort.Stable(logshuttlesByName(ls))
	return ls, nil
}

// CreateLogshuttleInput is used as input to the CreateLogshuttle function.
type CreateLogshuttleInput struct {
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

// CreateLogshuttle creates a new Fastly logshuttle.
func (c *Client) CreateLogshuttle(i *CreateLogshuttleInput) (*Logshuttle, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logshuttle", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Logshuttle
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// GetLogshuttleInput is used as input to the GetLogshuttle function.
type GetLogshuttleInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the logshuttle to fetch.
	Name string
}

// GetLogshuttle gets the logshuttle configuration with the given parameters.
func (c *Client) GetLogshuttle(i *GetLogshuttleInput) (*Logshuttle, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logshuttle/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Logshuttle
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// UpdateLogshuttleInput is used as input to the UpdateLogshuttle function.
type UpdateLogshuttleInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the logshuttle to update.
	Name string

	NewName           *string `url:"name,omitempty"`
	Format            *string `url:"format,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
	URL               *string `url:"url,omitempty"`
	Token             *string `url:"token,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	Placement         *string `url:"placement,omitempty"`
}

// UpdateLogshuttle updates a specific logshuttle.
func (c *Client) UpdateLogshuttle(i *UpdateLogshuttleInput) (*Logshuttle, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logshuttle/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Logshuttle
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// DeleteLogshuttleInput is the input parameter to DeleteLogshuttle.
type DeleteLogshuttleInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the logshuttle to delete (required).
	Name string
}

// DeleteLogshuttle deletes the given logshuttle version.
func (c *Client) DeleteLogshuttle(i *DeleteLogshuttleInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logshuttle/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
