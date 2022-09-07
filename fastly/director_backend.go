package fastly

import (
	"fmt"
	"net/url"
	"time"
)

// DirectorBackend is the relationship between a director and a backend in the
// Fastly API.
type DirectorBackend struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Director  string     `mapstructure:"director_name"`
	Backend   string     `mapstructure:"backend_name"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
}

// CreateDirectorBackendInput is used as input to the CreateDirectorBackend
// function.
type CreateDirectorBackendInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Director is the name of the director (required).
	Director string

	// Backend is the name of the backend (required).
	Backend string
}

// CreateDirectorBackend creates a new Fastly backend.
func (c *Client) CreateDirectorBackend(i *CreateDirectorBackendInput) (*DirectorBackend, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Director == "" {
		return nil, ErrMissingDirector
	}

	if i.Backend == "" {
		return nil, ErrMissingBackend
	}

	path := fmt.Sprintf("/service/%s/version/%d/director/%s/backend/%s",
		i.ServiceID, i.ServiceVersion, url.PathEscape(i.Director), url.PathEscape(i.Backend))

	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DirectorBackend
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetDirectorBackendInput is used as input to the GetDirectorBackend function.
type GetDirectorBackendInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Director is the name of the director (required).
	Director string

	// Backend is the name of the backend (required).
	Backend string
}

// GetDirectorBackend gets the backend configuration with the given parameters.
func (c *Client) GetDirectorBackend(i *GetDirectorBackendInput) (*DirectorBackend, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Director == "" {
		return nil, ErrMissingDirector
	}

	if i.Backend == "" {
		return nil, ErrMissingBackend
	}

	path := fmt.Sprintf("/service/%s/version/%d/director/%s/backend/%s",
		i.ServiceID, i.ServiceVersion, url.PathEscape(i.Director), url.PathEscape(i.Backend))

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DirectorBackend
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteDirectorBackendInput is the input parameter to DeleteDirectorBackend.
type DeleteDirectorBackendInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Director is the name of the director (required).
	Director string

	// Backend is the name of the backend (required).
	Backend string
}

// DeleteDirectorBackend deletes the given backend version.
func (c *Client) DeleteDirectorBackend(i *DeleteDirectorBackendInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Director == "" {
		return ErrMissingDirector
	}

	if i.Backend == "" {
		return ErrMissingBackend
	}

	path := fmt.Sprintf("/service/%s/version/%d/director/%s/backend/%s",
		i.ServiceID, i.ServiceVersion, url.PathEscape(i.Director), url.PathEscape(i.Backend))

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
