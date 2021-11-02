package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

const (
	// PoolTypeRandom is a pool that does random direction.
	PoolTypeRandom PoolType = "random"

	// PoolTypeHash is a pool that does hash direction.
	PoolTypeHash PoolType = "hash"

	// PoolTypeClient ins a pool that does client direction.
	PoolTypeClient PoolType = "client"
)

// PoolType is a type of pool.
type PoolType string

// PPoolType returns pointer to PoolType.
func PPoolType(t PoolType) *PoolType {
	pt := PoolType(t)
	return &pt
}

// Pool represents a pool response from the Fastly API.
type Pool struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	ID               string     `mapstructure:"id"`
	Name             string     `mapstructure:"name"`
	Comment          string     `mapstructure:"comment"`
	Shield           string     `mapstructure:"shield"`
	RequestCondition string     `mapstructure:"request_condition"`
	MaxConnDefault   uint       `mapstructure:"max_conn_default"`
	ConnectTimeout   uint       `mapstructure:"connect_timeout"`
	FirstByteTimeout uint       `mapstructure:"first_byte_timeout"`
	Quorum           uint       `mapstructure:"quorum"`
	UseTLS           bool       `mapstructure:"use_tls"`
	TLSCACert        string     `mapstructure:"tls_ca_cert"`
	TLSCiphers       string     `mapstructure:"tls_ciphers"`
	TLSClientKey     string     `mapstructure:"tls_client_key"`
	TLSClientCert    string     `mapstructure:"tls_client_cert"`
	TLSSNIHostname   string     `mapstructure:"tls_sni_hostname"`
	TLSCheckCert     bool       `mapstructure:"tls_check_cert"`
	TLSCertHostname  string     `mapstructure:"tls_cert_hostname"`
	MinTLSVersion    string     `mapstructure:"min_tls_version"`
	MaxTLSVersion    string     `mapstructure:"max_tls_version"`
	Healthcheck      string     `mapstructure:"healthcheck"`
	Type             PoolType   `mapstructure:"type"`
	OverrideHost     string     `mapstructure:"override_host"`
	CreatedAt        *time.Time `mapstructure:"created_at"`
	DeletedAt        *time.Time `mapstructure:"deleted_at"`
	UpdatedAt        *time.Time `mapstructure:"updated_at"`
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListPools lists all pools for a particular service and version.
func (c *Client) ListPools(i *ListPoolsInput) ([]*Pool, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ps []*Pool
	if err := decodeBodyMap(resp.Body, &ps); err != nil {
		return nil, err
	}
	sort.Stable(poolsByName(ps))
	return ps, nil
}

// CreatePoolInput is used as input to the CreatePool function.
type CreatePoolInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the pool to create (required).
	Name string `url:"name"`

	// Optional fields.
	Comment          string      `url:"comment,omitempty"`
	Shield           string      `url:"shield,omitempty"`
	RequestCondition string      `url:"request_condition,omitempty"`
	MaxConnDefault   uint        `url:"max_conn_default,omitempty"`
	ConnectTimeout   uint        `url:"connect_timeout,omitempty"`
	FirstByteTimeout uint        `url:"first_byte_timeout,omitempty"`
	Quorum           uint        `url:"quorum,omitempty"`
	UseTLS           Compatibool `url:"use_tls,omitempty"`
	TLSCACert        string      `url:"tls_ca_cert,omitempty"`
	TLSCiphers       string      `url:"tls_ciphers,omitempty"`
	TLSClientKey     string      `url:"tls_client_key,omitempty"`
	TLSClientCert    string      `url:"tls_client_cert,omitempty"`
	TLSSNIHostname   string      `url:"tls_sni_hostname,omitempty"`
	TLSCheckCert     Compatibool `url:"tls_check_cert,omitempty"`
	TLSCertHostname  string      `url:"tls_cert_hostname,omitempty"`
	MinTLSVersion    string      `url:"min_tls_version,omitempty"`
	MaxTLSVersion    string      `url:"max_tls_version,omitempty"`
	Healthcheck      string      `url:"healthcheck,omitempty"`
	Type             PoolType    `url:"type,omitempty"`
	OverrideHost     string      `url:"override_host,omitempty"`
}

// CreatePool creates a pool for a particular service and version.
func (c *Client) CreatePool(i *CreatePoolInput) (*Pool, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var p *Pool
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// GetPoolInput is used as input to the GetPool function.
type GetPoolInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the pool of interest (required).
	Name string
}

// GetPool gets a single pool for a particular service and version.
func (c *Client) GetPool(i *GetPoolInput) (*Pool, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var p *Pool
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// UpdatePoolInput is used as input to the UpdatePool function.
type UpdatePoolInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the pool to update (required).
	Name string

	// Optional fields.
	NewName          *string      `url:"name,omitempty"`
	Comment          *string      `url:"comment,omitempty"`
	Shield           *string      `url:"shield,omitempty"`
	RequestCondition *string      `url:"request_condition,omitempty"`
	MaxConnDefault   *uint        `url:"max_conn_default,omitempty"`
	ConnectTimeout   *uint        `url:"connect_timeout,omitempty"`
	FirstByteTimeout *uint        `url:"first_byte_timeout,omitempty"`
	Quorum           *uint        `url:"quorum,omitempty"`
	UseTLS           *Compatibool `url:"use_tls,omitempty"`
	TLSCACert        *string      `url:"tls_ca_cert,omitempty"`
	TLSCiphers       *string      `url:"tls_ciphers,omitempty"`
	TLSClientKey     *string      `url:"tls_client_key,omitempty"`
	TLSClientCert    *string      `url:"tls_client_cert,omitempty"`
	TLSSNIHostname   *string      `url:"tls_sni_hostname,omitempty"`
	TLSCheckCert     *Compatibool `url:"tls_check_cert,omitempty"`
	TLSCertHostname  *string      `url:"tls_cert_hostname,omitempty"`
	MinTLSVersion    *string      `url:"min_tls_version,omitempty"`
	MaxTLSVersion    *string      `url:"max_tls_version,omitempty"`
	Healthcheck      *string      `url:"healthcheck,omitempty"`
	Type             *PoolType    `url:"type,omitempty"`
	OverrideHost     *string      `url:"override_host,omitempty"`
}

// UpdatePool updates a specufic pool for a particular service and version.
func (c *Client) UpdatePool(i *UpdatePoolInput) (*Pool, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var p *Pool
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// DeletePoolInput is used as input to the DeletePool function.
type DeletePoolInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the pool to delete (required).
	Name string
}

// DeletePool deletes a specific pool for a particular service and version.
func (c *Client) DeletePool(i *DeletePoolInput) error {
	if i.ServiceID == "" {

		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/pool/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
