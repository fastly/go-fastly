package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Pool represents a pool response from the Fastly API.
type Pool struct {
	Id               string `mapstructure:"id"`
	Service          string `mapstructure:"service_id"`
	Version          int    `mapstructure:"version"`

	Name             string `mapstructure:"name"`
	Shield           string `mapstructure:"shield"`
	OverrideHost     string `mapstructure:"override_host"`
	UseTls           bool   `mapstructure:"use_tls"`
	Type             string `mapstructure:"type"`
	RequestCondition string `mapstructure:"request_condition"`
	MaxConnDefault   uint   `mapstructure:"max_conn_default"`
	ConnectTimeout   uint   `mapstructure:"connect_timeout"`
	FirstByteTimeout uint   `mapstructure:"first_byte_timeout"`
	Quorum           uint   `mapstructure:"quorum"`
	TlsCaCert        string `mapstructure:"tls_ca_cert"`
	TlsCiphers       string `mapstructure:"tls_ciphers"`
	TlsClientKey     string `mapstructure:"tls_client_key"`
	TlsClientCert    string `mapstructure:"tls_client_cert"`
	TlsSniHostname   string `mapstructure:"tls_sni_hostname"`
	TlsCheckCert     string `mapstructure:"tls_check_cert"`
	TlsCertHostname  string `mapstructure:"tls_cert_hostname"`
	MinTlsVersion    string `mapstructure:"min_tls_version"`
	MaxTlsVersion    string `mapstructure:"max_tls_version"`
	HealthCheck      string `mapstructure:"healthcheck"`
	Comment          string `mapstructure:"comment"`
	CreatedAt        *time.Time `mapstructure:"created_at"`
	UpdatedAt        *time.Time `mapstructure:"updated_at"`
	DeletedAt        *time.Time `mapstructure:"deleted_at"`
}

// poolsByName is a sortable list of pools.
type poolsByName []*Pool

// Len, Swap, and Less implement the sortable interface.
func (s poolsByName) Len() int      { return len(s) }
func (s poolsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s poolsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListPoolsInput is used as input to the ListPools function.
type ListPoolsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListPools returns the list of Pools for the service.
// version.
func (c *Client) ListPools(i *ListPoolsInput) ([]*Pool, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var pools []*Pool
	if err := decodeJSON(&pools, resp.Body); err != nil {
		return nil, err
	}
	sort.Stable(poolsByName(pools))
	return pools, nil
}

// CreatePoolInput is used as input to the CreatePool function.
type CreatePoolInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

        // Name is the name of the pool to create.
	Name             string `form:"name"`

	Shield           string `form:"shield,omitempty"`
	OverrideHost     string `form:"override_host,omitempty"`
	UseTls           bool   `form:"use_tls,omitempty"`
	Type             string `form:"type,omitempty"`
	RequestCondition string `form:"request_condition,omitempty"`
	MaxConnDefault   uint   `form:"max_conn_default,omitempty"`
	ConnectTimeout   uint   `form:"connect_timeout,omitempty"`
	FirstByteTimeout uint   `form:"first_byte_timeout,omitempty"`
	Quorum           uint   `form:"quorum,omitempty"`
	TlsCaCert        string `form:"tls_ca_cert,omitempty"`
	TlsCiphers       string `form:"tls_ciphers,omitempty"`
	TlsClientKey     string `form:"tls_client_key,omitempty"`
	TlsClientCert    string `form:"tls_client_cert,omitempty"`
	TlsSniHostname   string `form:"tls_sni_hostname,omitempty"`
	TlsCheckCert     string `form:"tls_check_cert,omitempty"`
	TlsCertHostname  string `form:"tls_cert_hostname,omitempty"`
	MinTlsVersion    string `form:"min_tls_version,omitempty"`
	MaxTlsVersion    string `form:"max_tls_version,omitempty"`
	HealthCheck      string `form:"healthcheck,omitempty"`
	Comment          string `form:"comment,omitempty"`
}

// CreatePool creates a new Fastly dyanmic server pool.
func (c *Client) CreatePool(i *CreatePoolInput) (*Pool, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

        if i.Name == "" {
                return nil, ErrMissingName
        }

	path := fmt.Sprintf("/service/%s/version/%d/pool", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var pool *Pool
	if err := decodeJSON(&pool, resp.Body); err != nil {
		return nil, err
	}
	return pool, nil
}

// GetPoolInput is used as input to the GetPool function.
type GetPoolInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the pool to fetch.
	Name string
}

// GetPool gets the pool configuration with the given Service, Version and Name.
func (c *Client) GetPool(i *GetPoolInput) (*Pool, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var pool *Pool
	if err := decodeJSON(&pool, resp.Body); err != nil {
		return nil, err
	}
	return pool, nil
}

// UpdatePoolInput is used as input to the UpdatePool function.
type UpdatePoolInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the pool to update.
	Name string

        NewName          string `form:"name,omitempty"`
	Shield           string `form:"shield,omitempty"`
	OverrideHost     string `form:"override_host,omitempty"`
	UseTls           bool   `form:"use_tls,omitempty"`
	Type             string `form:"type,omitempty"`
	RequestCondition string `form:"request_condition,omitempty"`
	MaxConnDefault   uint   `form:"max_conn_default,omitempty"`
	ConnectTimeout   uint   `form:"connect_timeout,omitempty"`
	FirstByteTimeout uint   `form:"first_byte_timeout,omitempty"`
	Quorum           uint   `form:"quorum,omitempty"`
	TlsCaCert        string `form:"tls_ca_cert,omitempty"`
	TlsCiphers       string `form:"tls_ciphers,omitempty"`
	TlsClientKey     string `form:"tls_client_key,omitempty"`
	TlsClientCert    string `form:"tls_client_cert,omitempty"`
	TlsSniHostname   string `form:"tls_sni_hostname,omitempty"`
	TlsCheckCert     string `form:"tls_check_cert,omitempty"`
	TlsCertHostname  string `form:"tls_cert_hostname,omitempty"`
	MinTlsVersion    string `form:"min_tls_version,omitempty"`
	MaxTlsVersion    string `form:"max_tls_version,omitempty"`
	HealthCheck      string `form:"healthcheck,omitempty"`
	Comment          string `form:"comment,omitempty"`
}

// UpdatePool updates a specific pool.
func (c *Client) UpdatePool(i *UpdatePoolInput) (*Pool, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var pool *Pool
	if err := decodeJSON(&pool, resp.Body); err != nil {
		return nil, err
	}
	return pool, nil
}

// DeletePoolInput is the input parameter to DeletePool.
type DeletePoolInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the pool to delete (required).
	Name string
}

// DeletePool deletes the given pool.
func (c *Client) DeletePool(i *DeletePoolInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool/%s", i.Service, i.Version, url.PathEscape(i.Name))
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


