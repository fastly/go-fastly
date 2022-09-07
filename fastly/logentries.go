package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Logentries represents a logentries response from the Fastly API.
type Logentries struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Port              uint       `mapstructure:"port"`
	UseTLS            bool       `mapstructure:"use_tls"`
	Token             string     `mapstructure:"token"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Region            string     `mapstructure:"region"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Placement         string     `mapstructure:"placement"`
}

// logentriesByName is a sortable list of logentries.
type logentriesByName []*Logentries

// Len, Swap, and Less implement the sortable interface.
func (s logentriesByName) Len() int      { return len(s) }
func (s logentriesByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s logentriesByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListLogentriesInput is used as input to the ListLogentries function.
type ListLogentriesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListLogentries returns the list of logentries for the configuration version.
func (c *Client) ListLogentries(i *ListLogentriesInput) ([]*Logentries, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logentries", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ls []*Logentries
	if err := decodeBodyMap(resp.Body, &ls); err != nil {
		return nil, err
	}
	sort.Stable(logentriesByName(ls))
	return ls, nil
}

// CreateLogentriesInput is used as input to the CreateLogentries function.
type CreateLogentriesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string      `url:"name,omitempty"`
	Port              uint        `url:"port,omitempty"`
	UseTLS            Compatibool `url:"use_tls,omitempty"`
	Token             string      `url:"token,omitempty"`
	Format            string      `url:"format,omitempty"`
	FormatVersion     uint        `url:"format_version,omitempty"`
	ResponseCondition string      `url:"response_condition,omitempty"`
	Region            string      `url:"region,omitempty"`
	Placement         string      `url:"placement,omitempty"`
}

// CreateLogentries creates a new Fastly logentries.
func (c *Client) CreateLogentries(i *CreateLogentriesInput) (*Logentries, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logentries", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Logentries
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// GetLogentriesInput is used as input to the GetLogentries function.
type GetLogentriesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the logentries to fetch.
	Name string
}

// GetLogentries gets the logentries configuration with the given parameters.
func (c *Client) GetLogentries(i *GetLogentriesInput) (*Logentries, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logentries/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Logentries
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// UpdateLogentriesInput is used as input to the UpdateLogentries function.
type UpdateLogentriesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the logentries to update.
	Name string

	NewName           *string      `url:"name,omitempty"`
	Port              *uint        `url:"port,omitempty"`
	UseTLS            *Compatibool `url:"use_tls,omitempty"`
	Token             *string      `url:"token,omitempty"`
	Format            *string      `url:"format,omitempty"`
	FormatVersion     *uint        `url:"format_version,omitempty"`
	ResponseCondition *string      `url:"response_condition,omitempty"`
	Region            *string      `url:"region,omitempty"`
	Placement         *string      `url:"placement,omitempty"`
}

// UpdateLogentries updates a specific logentries.
func (c *Client) UpdateLogentries(i *UpdateLogentriesInput) (*Logentries, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logentries/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Logentries
	if err := decodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// DeleteLogentriesInput is the input parameter to DeleteLogentries.
type DeleteLogentriesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the logentries to delete (required).
	Name string
}

// DeleteLogentries deletes the given logentries version.
func (c *Client) DeleteLogentries(i *DeleteLogentriesInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/logentries/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
		return ErrNotOK
	}
	return nil
}
