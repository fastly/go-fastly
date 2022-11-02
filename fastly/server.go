package fastly

import (
	"fmt"
	"sort"
	"time"
)

// Server represents a server response from the Fastly API.
type Server struct {
	Address      string     `mapstructure:"address"`
	Comment      string     `mapstructure:"comment"`
	CreatedAt    *time.Time `mapstructure:"created_at"`
	DeletedAt    *time.Time `mapstructure:"deleted_at"`
	Disabled     bool       `mapstructure:"disabled"`
	ID           string     `mapstructure:"id"`
	MaxConn      uint       `mapstructure:"max_conn"`
	OverrideHost string     `mapstructure:"override_host"`
	PoolID       string     `mapstructure:"pool_id"`
	Port         uint       `mapstructure:"port"`
	ServiceID    string     `mapstructure:"service_id"`
	UpdatedAt    *time.Time `mapstructure:"updated_at"`
	Weight       uint       `mapstructure:"weight"`
}

// serversByAddress is a sortable list of servers.
type serversByAddress []*Server

// Len implement the sortable interface.
func (s serversByAddress) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s serversByAddress) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s serversByAddress) Less(i, j int) bool {
	return s[i].Address < s[j].Address
}

// ListServersInput is used as input to the ListServers function.
type ListServersInput struct {
	// PoolID is the ID of the pool (required).
	PoolID string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// ListServers retrieves all resources.
func (c *Client) ListServers(i *ListServersInput) ([]*Server, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.PoolID == "" {
		return nil, ErrMissingPoolID
	}

	path := fmt.Sprintf("/service/%s/pool/%s/servers", i.ServiceID, i.PoolID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ss []*Server
	if err := decodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	sort.Stable(serversByAddress(ss))
	return ss, nil
}

// CreateServerInput is used as input to the CreateServer function.
type CreateServerInput struct {
	// Address is the hostname or IP of the origin server (required).
	Address      string `url:"address"`
	Comment      string `url:"comment,omitempty"`
	Disabled     bool   `url:"disabled,omitempty"`
	MaxConn      uint   `url:"max_conn,omitempty"`
	OverrideHost string `url:"override_host,omitempty"`
	// PoolID is the ID of the pool (required).
	PoolID string
	Port   uint `url:"port,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	Weight    uint `url:"weight,omitempty"`
}

// CreateServer creates a new resource.
// Servers are versionless resources that are associated with a Pool.
func (c *Client) CreateServer(i *CreateServerInput) (*Server, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.PoolID == "" {
		return nil, ErrMissingPoolID
	}

	if i.Address == "" {
		return nil, ErrMissingAddress
	}

	path := fmt.Sprintf("/service/%s/pool/%s/server", i.ServiceID, i.PoolID)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Server
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetServerInput is used as input to the GetServer function.
type GetServerInput struct {
	// PoolID is the ID of the pool (required).
	PoolID string
	Server string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// GetServer retrieves the specified resource.
func (c *Client) GetServer(i *GetServerInput) (*Server, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.PoolID == "" {
		return nil, ErrMissingPoolID
	}

	if i.Server == "" {
		return nil, ErrMissingServer
	}

	path := fmt.Sprintf("/service/%s/pool/%s/server/%s", i.ServiceID, i.PoolID, i.Server)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Server
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateServerInput is used as input to the UpdateServer function.
type UpdateServerInput struct {
	Address      *string `url:"address,omitempty"`
	Comment      *string `url:"comment,omitempty"`
	Disabled     *bool   `url:"disabled,omitempty"`
	MaxConn      *uint   `url:"max_conn,omitempty"`
	OverrideHost *string `url:"override_host,omitempty"`
	// PoolID is the ID of the pool (required).
	PoolID string
	Port   *uint `url:"port,omitempty"`
	Server string
	// ServiceID is the ID of the service (required).
	ServiceID string
	Weight    *uint `url:"weight,omitempty"`
}

// UpdateServer updates the specified resource.
func (c *Client) UpdateServer(i *UpdateServerInput) (*Server, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.PoolID == "" {
		return nil, ErrMissingPoolID
	}

	if i.Server == "" {
		return nil, ErrMissingServer
	}

	path := fmt.Sprintf("/service/%s/pool/%s/server/%s", i.ServiceID, i.PoolID, i.Server)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Server
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteServerInput is used as input to the DeleteServer function.
type DeleteServerInput struct {
	// PoolID is the ID of the pool (required).
	PoolID string
	Server string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// DeleteServer deletes the specified resource.
func (c *Client) DeleteServer(i *DeleteServerInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.PoolID == "" {
		return ErrMissingPoolID
	}

	if i.Server == "" {
		return ErrMissingServer
	}

	path := fmt.Sprintf("/service/%s/pool/%s/server/%s", i.ServiceID, i.PoolID, i.Server)
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
