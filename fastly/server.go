package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Server represents a server response from the Fastly API.
type Server struct {
	Id      string `mapstructure:"id"`
	Service string `mapstructure:"service_id"`
	Pool    string `mapstructure:"pool_id"`

	Weight       string `mapstructure:"weight"`
	MaxConn      string `mapstructure:"max_conn"`
	Port         string `mapstructure:"port"`
	Address      string `mapstructure:"address"`
	Comment      string `mapstructure:"comment"`
	Disabled     bool   `mapstructure:"disabled"`
	OverrideHost string `mapstructure:"override_host"`

	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
}

// serversById is a sortable list of servers.
type serversById []*Server

// Len, Swap, and Less implement the sortable interface.
func (s serversById) Len() int      { return len(s) }
func (s serversById) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s serversById) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

// ListServersInput is used as input to the ListServers function.
type ListServersInput struct {
	// Service is the ID of the service (required).
	Service string

	// Pool is the ID of the server pool (required).
	Pool string
}

// ListServers returns the list of Servers for the service configuration
// version.
func (c *Client) ListServers(i *ListServersInput) ([]*Server, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Pool == "" {
		return nil, ErrMissingPool
	}

	path := fmt.Sprintf("/service/%s/pool/%s/servers", i.Service, i.Pool)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var servers []*Server
	if err := decodeJSON(&servers, resp.Body); err != nil {
		return nil, err
	}
	sort.Stable(serversById(servers))
	return servers, nil
}

// CreateServerInput is used as input to the CreateServer function.
type CreateServerInput struct {
	// Service is the ID of the service.
	// Pool is the ID of the pool that the server is associated with
	// Both fields are required.
	Service string
	Pool    string

	Weight       string `form:"weight,omitempty"`
	MaxConn      string `form:"max_conn,omitempty"`
	Port         string `form:"port,omitempty"`
	Address      string `form:"address,omitempty"`
	Comment      string `form:"comment,omitempty"`
	Disabled     bool   `form:"disabled,omitempty"`
	OverrideHost string `form:"override_host,omitempty"`
}

// CreateServer creates a new Fastly health check.
func (c *Client) CreateServer(i *CreateServerInput) (*Server, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Pool == "" {
		return nil, ErrMissingPool
	}

	path := fmt.Sprintf("/service/%s/pool/%s/server", i.Service, i.Pool)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Server
	if err := decodeJSON(&s, resp.Body); err != nil {
		return nil, err
	}
	return s, nil
}

// GetServerInput is used as input to the GetServer function.
type GetServerInput struct {
	// Service is the ID of the service (required).
	Service string

	// Pool is the ID of the server pool (required).
	Pool string

	// Id of the server to fetch.
	Id string
}

// GetServer gets the Server configuration with the given parameters.
func (c *Client) GetServer(i *GetServerInput) (*Server, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Pool == "" {
		return nil, ErrMissingPool
	}

	if i.Id == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/pool/%s/server/%s", i.Service, i.Pool, url.PathEscape(i.Id))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var s *Server
	if err := decodeJSON(&s, resp.Body); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateServerInput is used as input to the UpdateServer function.
type UpdateServerInput struct {
	// ServiceiID is the ID of the service.
	// Pool is the ID of the pool that the server is associated with
	// Id of the server to update
	// These 3 fields are required.
	Service string
	Pool    string
	Id      string

	Weight       string `form:"weight,omitempty"`
	MaxConn      string `form:"max_conn,omitempty"`
	Port         string `form:"port,omitempty"`
	Address      string `form:"address,omitempty"`
	Comment      string `form:"comment,omitempty"`
	Disabled     bool   `form:"disabled,omitempty"`
	OverrideHost string `form:"override_host,omitempty"`
}

// UpdateServerInput updates a specific server configuration.
func (c *Client) UpdateServer(i *UpdateServerInput) (*Server, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Pool == "" {
		return nil, ErrMissingPool
	}

	if i.Id == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/pool/%s/server/%s", i.Service, i.Pool, url.PathEscape(i.Id))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Server
	if err := decodeJSON(&s, resp.Body); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteServerInput is the input parameter to DeleteServer.
type DeleteServerInput struct {
	// Service is the ID of the service.
	// Pool is the ID of the pool that the server is associated with
	// Server is the ID of the server to delete
	// These 3 fields are required.
	Service string
	Pool    string
	Id      string
}

// DeleteServer deletes the given health check.
func (c *Client) DeleteServer(i *DeleteServerInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Pool == "" {
		return ErrMissingPool
	}

	if i.Id == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/pool/%s/server/%s", i.Service, i.Pool, url.PathEscape(i.Id))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeJSON(&r, resp.Body); err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf("Not Ok")
	}
	return nil
}
