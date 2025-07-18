package fastly

import (
	"context"
	"strconv"
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

// Director represents a director response from the Fastly API.
type Director struct {
	Backends       []string      `mapstructure:"backends"`
	Capacity       *int          `mapstructure:"capacity"`
	Comment        *string       `mapstructure:"comment"`
	CreatedAt      *time.Time    `mapstructure:"created_at"`
	DeletedAt      *time.Time    `mapstructure:"deleted_at"`
	Name           *string       `mapstructure:"name"`
	Quorum         *int          `mapstructure:"quorum"`
	Retries        *int          `mapstructure:"retries"`
	ServiceID      *string       `mapstructure:"service_id"`
	ServiceVersion *int          `mapstructure:"version"`
	Shield         *string       `mapstructure:"shield"`
	Type           *DirectorType `mapstructure:"type"`
	UpdatedAt      *time.Time    `mapstructure:"updated_at"`
}

// ListDirectorsInput is used as input to the ListDirectors function.
type ListDirectorsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDirectors retrieves all resources.
func (c *Client) ListDirectors(ctx context.Context, i *ListDirectorsInput) ([]*Director, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "director")
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ds []*Director
	if err := DecodeBodyMap(resp.Body, &ds); err != nil {
		return nil, err
	}
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
func (c *Client) CreateDirector(ctx context.Context, i *CreateDirectorInput) (*Director, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "director")
	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Director
	if err := DecodeBodyMap(resp.Body, &d); err != nil {
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
func (c *Client) GetDirector(ctx context.Context, i *GetDirectorInput) (*Director, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "director", i.Name)
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Director
	if err := DecodeBodyMap(resp.Body, &d); err != nil {
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
	Type *DirectorType `url:"type,omitempty"`
}

// UpdateDirector updates the specified resource.
func (c *Client) UpdateDirector(ctx context.Context, i *UpdateDirectorInput) (*Director, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "director", i.Name)
	resp, err := c.PutForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Director
	if err := DecodeBodyMap(resp.Body, &d); err != nil {
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
func (c *Client) DeleteDirector(ctx context.Context, i *DeleteDirectorInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "director", i.Name)
	resp, err := c.Delete(ctx, path, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
