package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

const (
	// DirectorTypeRandom is a director that does random direction.
	DirectorTypeRandom DirectorType = 1

	// DirectorTypeRoundRobin is a director that does round-robin direction.
	DirectorTypeRoundRobin DirectorType = 2

	// DirectorTypeHash is a director that does hash direction.
	DirectorTypeHash DirectorType = 3

	// DirectorTypeClient is a director that does client direction.
	DirectorTypeClient DirectorType = 4
)

// DirectorType is a type of director.
type DirectorType uint8

// Director represents a director response from the Fastly API.
type Director struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name      string       `mapstructure:"name"`
	Backends  []string     `mapstructure:"backends"`
	Comment   string       `mapstructure:"comment"`
	Shield    string       `mapstructure:"shield"`
	Quorum    uint         `mapstructure:"quorum"`
	Type      DirectorType `mapstructure:"type"`
	Retries   uint         `mapstructure:"retries"`
	Capacity  uint         `mapstructure:"capacity"`
	CreatedAt *time.Time   `mapstructure:"created_at"`
	UpdatedAt *time.Time   `mapstructure:"updated_at"`
	DeletedAt *time.Time   `mapstructure:"deleted_at"`
}

// directorsByName is a sortable list of directors.
type directorsByName []*Director

// Len, Swap, and Less implement the sortable interface.
func (s directorsByName) Len() int      { return len(s) }
func (s directorsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s directorsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListDirectorsInput is used as input to the ListDirectors function.
type ListDirectorsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDirectors returns the list of directors for the configuration version.
func (c *Client) ListDirectors(i *ListDirectorsInput) ([]*Director, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/director", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ds []*Director
	if err := decodeBodyMap(resp.Body, &ds); err != nil {
		return nil, err
	}
	sort.Stable(directorsByName(ds))
	return ds, nil
}

// CreateDirectorInput is used as input to the CreateDirector function.
type CreateDirectorInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name     string       `url:"name,omitempty"`
	Comment  string       `url:"comment,omitempty"`
	Shield   string       `url:"shield,omitempty"`
	Quorum   *uint        `url:"quorum,omitempty"`
	Type     DirectorType `url:"type,omitempty"`
	Retries  *uint        `url:"retries,omitempty"`
	Capacity *uint        `url:"capacity,omitempty"`
}

// CreateDirector creates a new Fastly director.
func (c *Client) CreateDirector(i *CreateDirectorInput) (*Director, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/director", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Director
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetDirectorInput is used as input to the GetDirector function.
type GetDirectorInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the director to fetch.
	Name string
}

// GetDirector gets the director configuration with the given parameters.
func (c *Client) GetDirector(i *GetDirectorInput) (*Director, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/director/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Director
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// UpdateDirectorInput is used as input to the UpdateDirector function.
type UpdateDirectorInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the director to update.
	Name string

	NewName  *string      `url:"name,omitempty"`
	Comment  *string      `url:"comment,omitempty"`
	Shield   *string      `url:"shield,omitempty"`
	Quorum   *uint        `url:"quorum,omitempty"`
	Type     DirectorType `url:"type,omitempty"`
	Retries  *uint        `url:"retries,omitempty"`
	Capacity *uint        `url:"capacity,omitempty"`
}

// UpdateDirector updates a specific director.
func (c *Client) UpdateDirector(i *UpdateDirectorInput) (*Director, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/director/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Director
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteDirectorInput is the input parameter to DeleteDirector.
type DeleteDirectorInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the director to delete (required).
	Name string
}

// DeleteDirector deletes the given director version.
func (c *Client) DeleteDirector(i *DeleteDirectorInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/director/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
