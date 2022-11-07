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
type DirectorType int

// DirectorTypePtr returns pointer to DirectorType.
func DirectorTypePtr(t DirectorType) *DirectorType {
	dt := DirectorType(t)
	return &dt
}

// Director represents a director response from the Fastly API.
type Director struct {
	Backends       []string     `mapstructure:"backends"`
	Capacity       int          `mapstructure:"capacity"`
	Comment        string       `mapstructure:"comment"`
	CreatedAt      *time.Time   `mapstructure:"created_at"`
	DeletedAt      *time.Time   `mapstructure:"deleted_at"`
	Name           string       `mapstructure:"name"`
	Quorum         int          `mapstructure:"quorum"`
	Retries        int          `mapstructure:"retries"`
	ServiceID      string       `mapstructure:"service_id"`
	ServiceVersion int          `mapstructure:"version"`
	Shield         string       `mapstructure:"shield"`
	Type           DirectorType `mapstructure:"type"`
	UpdatedAt      *time.Time   `mapstructure:"updated_at"`
}

// directorsByName is a sortable list of directors.
type directorsByName []*Director

// Len implement the sortable interface.
func (s directorsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s directorsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
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

// ListDirectors retrieves all resources.
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
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// Name is the name for the Director.
	Name *string `url:"name,omitempty"`
	// Quorum is the percentage of capacity that needs to be up for a director to be considered up. 0 to 100.
	Quorum *int `url:"quorum,omitempty"`
	// Retries is how many backends to search if it fails.
	Retries *int `url:"retries,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Shield is selected POP to serve as a shield for the backends.
	Shield *string `url:"shield,omitempty"`
	// Type is what type of load balance group to use (random, hash, client).
	Type *DirectorType `url:"type,omitempty"`
}

// CreateDirector creates a new resource.
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
	// Name is the name of the director to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetDirector retrieves the specified resource.
func (c *Client) GetDirector(i *GetDirectorInput) (*Director, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
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
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// Name is the name of the director to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Quorum is the percentage of capacity that needs to be up for a director to be considered up. 0 to 100.
	Quorum *int `url:"quorum,omitempty"`
	// Retries is how many backends to search if it fails.
	Retries *int `url:"retries,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Shield is selected POP to serve as a shield for the backends.
	Shield *string `url:"shield,omitempty"`
	// Type is what type of load balance group to use (random, hash, client).
	Type DirectorType `url:"type,omitempty"`
}

// UpdateDirector updates the specified resource.
func (c *Client) UpdateDirector(i *UpdateDirectorInput) (*Director, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
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
	// Name is the name of the director to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteDirector deletes the specified resource.
func (c *Client) DeleteDirector(i *DeleteDirectorInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
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
