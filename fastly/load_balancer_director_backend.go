package fastly

import (
	"context"
	"strconv"
	"time"
)

// DirectorBackend is the relationship between a director and a backend in the
// Fastly API.
type DirectorBackend struct {
	Backend        *string    `mapstructure:"backend_name"`
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	Director       *string    `mapstructure:"director_name"`
	ServiceID      *string    `mapstructure:"service_id"`
	ServiceVersion *int       `mapstructure:"version"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}

// CreateDirectorBackendInput is used as input to the CreateDirectorBackend
// function.
type CreateDirectorBackendInput struct {
	// Backend is the name of the backend (required).
	Backend string
	// Director is the name of the director (required).
	Director string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// CreateDirectorBackend creates a new resource.
func (c *Client) CreateDirectorBackend(ctx context.Context, i *CreateDirectorBackendInput) (*DirectorBackend, error) {
	if i.Backend == "" {
		return nil, ErrMissingBackend
	}
	if i.Director == "" {
		return nil, ErrMissingDirector
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "director", i.Director, "backend", i.Backend)

	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DirectorBackend
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetDirectorBackendInput is used as input to the GetDirectorBackend function.
type GetDirectorBackendInput struct {
	// Backend is the name of the backend (required).
	Backend string
	// Director is the name of the director (required).
	Director string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetDirectorBackend retrieves the specified resource.
func (c *Client) GetDirectorBackend(ctx context.Context, i *GetDirectorBackendInput) (*DirectorBackend, error) {
	if i.Backend == "" {
		return nil, ErrMissingBackend
	}
	if i.Director == "" {
		return nil, ErrMissingDirector
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "director", i.Director, "backend", i.Backend)

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DirectorBackend
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteDirectorBackendInput is the input parameter to DeleteDirectorBackend.
type DeleteDirectorBackendInput struct {
	// Backend is the name of the backend (required).
	Backend string
	// Director is the name of the director (required).
	Director string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteDirectorBackend deletes the specified resource.
func (c *Client) DeleteDirectorBackend(ctx context.Context, i *DeleteDirectorBackendInput) error {
	if i.Backend == "" {
		return ErrMissingBackend
	}
	if i.Director == "" {
		return ErrMissingDirector
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "director", i.Director, "backend", i.Backend)

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
