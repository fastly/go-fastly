package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Papertrail represents a papertrail response from the Fastly API.
type Papertrail struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Address           string     `mapstructure:"address"`
	Port              uint       `mapstructure:"port"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Placement         string     `mapstructure:"placement"`
}

// papertrailsByName is a sortable list of papertrails.
type papertrailsByName []*Papertrail

// Len, Swap, and Less implement the sortable interface.
func (s papertrailsByName) Len() int      { return len(s) }
func (s papertrailsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s papertrailsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListPapertrailsInput is used as input to the ListPapertrails function.
type ListPapertrailsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListPapertrails returns the list of papertrails for the configuration version.
func (c *Client) ListPapertrails(i *ListPapertrailsInput) ([]*Papertrail, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/papertrail", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ps []*Papertrail
	if err := decodeBodyMap(resp.Body, &ps); err != nil {
		return nil, err
	}
	sort.Stable(papertrailsByName(ps))
	return ps, nil
}

// CreatePapertrailInput is used as input to the CreatePapertrail function.
type CreatePapertrailInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string     `url:"name,omitempty"`
	Address           string     `url:"address,omitempty"`
	Port              uint       `url:"port,omitempty"`
	FormatVersion     uint       `url:"format_version,omitempty"`
	Format            string     `url:"format,omitempty"`
	ResponseCondition string     `url:"response_condition,omitempty"`
	CreatedAt         *time.Time `url:"created_at,omitempty"`
	UpdatedAt         *time.Time `url:"updated_at,omitempty"`
	DeletedAt         *time.Time `url:"deleted_at,omitempty"`
	Placement         string     `url:"placement,omitempty"`
}

// CreatePapertrail creates a new Fastly papertrail.
func (c *Client) CreatePapertrail(i *CreatePapertrailInput) (*Papertrail, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/papertrail", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p *Papertrail
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// GetPapertrailInput is used as input to the GetPapertrail function.
type GetPapertrailInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the papertrail to fetch.
	Name string
}

// GetPapertrail gets the papertrail configuration with the given parameters.
func (c *Client) GetPapertrail(i *GetPapertrailInput) (*Papertrail, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/papertrail/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p *Papertrail
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// UpdatePapertrailInput is used as input to the UpdatePapertrail function.
type UpdatePapertrailInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the papertrail to update.
	Name string

	NewName           *string    `url:"name,omitempty"`
	Address           *string    `url:"address,omitempty"`
	Port              *uint      `url:"port,omitempty"`
	FormatVersion     *uint      `url:"format_version,omitempty"`
	Format            *string    `url:"format,omitempty"`
	ResponseCondition *string    `url:"response_condition,omitempty"`
	CreatedAt         *time.Time `url:"created_at,omitempty"`
	UpdatedAt         *time.Time `url:"updated_at,omitempty"`
	DeletedAt         *time.Time `url:"deleted_at,omitempty"`
	Placement         *string    `url:"placement,omitempty"`
}

// UpdatePapertrail updates a specific papertrail.
func (c *Client) UpdatePapertrail(i *UpdatePapertrailInput) (*Papertrail, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/papertrail/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p *Papertrail
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// DeletePapertrailInput is the input parameter to DeletePapertrail.
type DeletePapertrailInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the papertrail to delete (required).
	Name string
}

// DeletePapertrail deletes the given papertrail version.
func (c *Client) DeletePapertrail(i *DeletePapertrailInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/papertrail/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
