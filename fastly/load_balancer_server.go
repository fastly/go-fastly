package fastly

import (
	"context"
	"time"
)

// Server represents a server response from the Fastly API.
type Server struct {
	Address      *string    `mapstructure:"address"`
	Comment      *string    `mapstructure:"comment"`
	CreatedAt    *time.Time `mapstructure:"created_at"`
	DeletedAt    *time.Time `mapstructure:"deleted_at"`
	Disabled     *bool      `mapstructure:"disabled"`
	MaxConn      *int       `mapstructure:"max_conn"`
	OverrideHost *string    `mapstructure:"override_host"`
	PoolID       *string    `mapstructure:"pool_id"`
	Port         *int       `mapstructure:"port"`
	ServerID     *string    `mapstructure:"id"`
	ServiceID    *string    `mapstructure:"service_id"`
	UpdatedAt    *time.Time `mapstructure:"updated_at"`
	Weight       *int       `mapstructure:"weight"`
}

// ListServersInput is used as input to the ListServers function.
type ListServersInput struct {
	// PoolID is the ID of the pool (required).
	PoolID string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// ListServers retrieves all resources.
func (c *Client) ListServers(ctx context.Context, i *ListServersInput) ([]*Server, error) {
	if i.PoolID == "" {
		return nil, ErrMissingPoolID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "pool", i.PoolID, "servers")

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ss []*Server
	if err := DecodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	return ss, nil
}

// CreateServerInput is used as input to the CreateServer function.
type CreateServerInput struct {
	// Address is the hostname or IP of the origin server (required).
	Address *string `url:"address,omitempty"`
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// Disabled allows servers to be enabled and disabled in a pool.
	Disabled *bool `url:"disabled,omitempty"`
	// MaxConn is the maximum number of connections. If the value is 0, it inherits the value from pool's max_conn_default.
	MaxConn *int `url:"max_conn,omitempty"`
	// OverrideHost is the hostname to override the Host header.
	OverrideHost *string `url:"override_host,omitempty"`
	// PoolID is the ID of the pool (required).
	PoolID string `url:"-"`
	// Port is the port number.
	Port *int `url:"port,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// Weight is the weight (1-100) used to load balance this server against others.
	Weight *int `url:"weight,omitempty"`
}

// CreateServer creates a new resource.
// Servers are versionless resources that are associated with a Pool.
func (c *Client) CreateServer(ctx context.Context, i *CreateServerInput) (*Server, error) {
	if i.PoolID == "" {
		return nil, ErrMissingPoolID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "pool", i.PoolID, "server")

	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Server
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetServerInput is used as input to the GetServer function.
type GetServerInput struct {
	// PoolID is the ID of the pool (required).
	PoolID string
	// Server is an alphanumeric string identifying a Server (required).
	Server string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// GetServer retrieves the specified resource.
func (c *Client) GetServer(ctx context.Context, i *GetServerInput) (*Server, error) {
	if i.PoolID == "" {
		return nil, ErrMissingPoolID
	}
	if i.Server == "" {
		return nil, ErrMissingServer
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "pool", i.PoolID, "server", i.Server)

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Server
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateServerInput is used as input to the UpdateServer function.
type UpdateServerInput struct {
	// Address is the hostname or IP of the origin server (required).
	Address *string `url:"address,omitempty"`
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// Disabled allows servers to be enabled and disabled in a pool.
	Disabled *bool `url:"disabled,omitempty"`
	// MaxConn is the maximum number of connections. If the value is 0, it inherits the value from pool's max_conn_default.
	MaxConn *int `url:"max_conn,omitempty"`
	// OverrideHost is the hostname to override the Host header.
	OverrideHost *string `url:"override_host,omitempty"`
	// PoolID is the ID of the pool (required).
	PoolID string `url:"-"`
	// Port is the port number.
	Port *int `url:"port,omitempty"`
	// Server is an alphanumeric string identifying a Server (required).
	Server string `url:"-"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// Weight is the weight (1-100) used to load balance this server against others.
	Weight *int `url:"weight,omitempty"`
}

// UpdateServer updates the specified resource.
func (c *Client) UpdateServer(ctx context.Context, i *UpdateServerInput) (*Server, error) {
	if i.PoolID == "" {
		return nil, ErrMissingPoolID
	}
	if i.Server == "" {
		return nil, ErrMissingServer
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "pool", i.PoolID, "server", i.Server)

	resp, err := c.PutForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Server
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteServerInput is used as input to the DeleteServer function.
type DeleteServerInput struct {
	// PoolID is the ID of the pool (required).
	PoolID string
	// Server is an alphanumeric string identifying a Server (required).
	Server string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// DeleteServer deletes the specified resource.
func (c *Client) DeleteServer(ctx context.Context, i *DeleteServerInput) error {
	if i.PoolID == "" {
		return ErrMissingPoolID
	}
	if i.Server == "" {
		return ErrMissingServer
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "pool", i.PoolID, "server", i.Server)

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
