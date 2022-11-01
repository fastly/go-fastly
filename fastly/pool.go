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
	Comment          string     `mapstructure:"comment"`
	ConnectTimeout   uint       `mapstructure:"connect_timeout"`
	CreatedAt        *time.Time `mapstructure:"created_at"`
	DeletedAt        *time.Time `mapstructure:"deleted_at"`
	FirstByteTimeout uint       `mapstructure:"first_byte_timeout"`
	Healthcheck      string     `mapstructure:"healthcheck"`
	ID               string     `mapstructure:"id"`
	MaxConnDefault   uint       `mapstructure:"max_conn_default"`
	MaxTLSVersion    string     `mapstructure:"max_tls_version"`
	MinTLSVersion    string     `mapstructure:"min_tls_version"`
	Name             string     `mapstructure:"name"`
	OverrideHost     string     `mapstructure:"override_host"`
	Quorum           uint       `mapstructure:"quorum"`
	RequestCondition string     `mapstructure:"request_condition"`
	ServiceID        string     `mapstructure:"service_id"`
	ServiceVersion   int        `mapstructure:"version"`
	Shield           string     `mapstructure:"shield"`
	TLSCACert        string     `mapstructure:"tls_ca_cert"`
	TLSCertHostname  string     `mapstructure:"tls_cert_hostname"`
	TLSCheckCert     bool       `mapstructure:"tls_check_cert"`
	TLSCiphers       string     `mapstructure:"tls_ciphers"`
	TLSClientCert    string     `mapstructure:"tls_client_cert"`
	TLSClientKey     string     `mapstructure:"tls_client_key"`
	TLSSNIHostname   string     `mapstructure:"tls_sni_hostname"`
	Type             PoolType   `mapstructure:"type"`
	UpdatedAt        *time.Time `mapstructure:"updated_at"`
	UseTLS           bool       `mapstructure:"use_tls"`
}

// poolsByName is a sortable list of pools.
type poolsByName []*Pool

// Len implement the sortable interface.
func (s poolsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s poolsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
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
	defer resp.Body.Close()

	var ps []*Pool
	if err := decodeBodyMap(resp.Body, &ps); err != nil {
		return nil, err
	}
	sort.Stable(poolsByName(ps))
	return ps, nil
}

// CreatePoolInput is used as input to the CreatePool function.
type CreatePoolInput struct {
	Comment          string `url:"comment,omitempty"`
	ConnectTimeout   uint   `url:"connect_timeout,omitempty"`
	FirstByteTimeout uint   `url:"first_byte_timeout,omitempty"`
	Healthcheck      string `url:"healthcheck,omitempty"`
	MaxConnDefault   uint   `url:"max_conn_default,omitempty"`
	MaxTLSVersion    string `url:"max_tls_version,omitempty"`
	MinTLSVersion    string `url:"min_tls_version,omitempty"`
	// Name is the name of the pool to create (required).
	Name             string `url:"name"`
	OverrideHost     string `url:"override_host,omitempty"`
	Quorum           uint   `url:"quorum,omitempty"`
	RequestCondition string `url:"request_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion  int
	Shield          string      `url:"shield,omitempty"`
	TLSCACert       string      `url:"tls_ca_cert,omitempty"`
	TLSCertHostname string      `url:"tls_cert_hostname,omitempty"`
	TLSCheckCert    Compatibool `url:"tls_check_cert,omitempty"`
	TLSCiphers      string      `url:"tls_ciphers,omitempty"`
	TLSClientCert   string      `url:"tls_client_cert,omitempty"`
	TLSClientKey    string      `url:"tls_client_key,omitempty"`
	TLSSNIHostname  string      `url:"tls_sni_hostname,omitempty"`
	Type            PoolType    `url:"type,omitempty"`
	UseTLS          Compatibool `url:"use_tls,omitempty"`
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
	defer resp.Body.Close()

	var p *Pool
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// GetPoolInput is used as input to the GetPool function.
type GetPoolInput struct {
	// Name is the name of the pool of interest (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
	defer resp.Body.Close()

	var p *Pool
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// UpdatePoolInput is used as input to the UpdatePool function.
type UpdatePoolInput struct {
	Comment          *string `url:"comment,omitempty"`
	ConnectTimeout   *uint   `url:"connect_timeout,omitempty"`
	FirstByteTimeout *uint   `url:"first_byte_timeout,omitempty"`
	Healthcheck      *string `url:"healthcheck,omitempty"`
	MaxConnDefault   *uint   `url:"max_conn_default,omitempty"`
	MaxTLSVersion    *string `url:"max_tls_version,omitempty"`
	MinTLSVersion    *string `url:"min_tls_version,omitempty"`
	// Name is the name of the pool to update (required).
	Name             string
	NewName          *string `url:"name,omitempty"`
	OverrideHost     *string `url:"override_host,omitempty"`
	Quorum           *uint   `url:"quorum,omitempty"`
	RequestCondition *string `url:"request_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion  int
	Shield          *string      `url:"shield,omitempty"`
	TLSCACert       *string      `url:"tls_ca_cert,omitempty"`
	TLSCertHostname *string      `url:"tls_cert_hostname,omitempty"`
	TLSCheckCert    *Compatibool `url:"tls_check_cert,omitempty"`
	TLSCiphers      *string      `url:"tls_ciphers,omitempty"`
	TLSClientCert   *string      `url:"tls_client_cert,omitempty"`
	TLSClientKey    *string      `url:"tls_client_key,omitempty"`
	TLSSNIHostname  *string      `url:"tls_sni_hostname,omitempty"`
	Type            *PoolType    `url:"type,omitempty"`
	UseTLS          *Compatibool `url:"use_tls,omitempty"`
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
	defer resp.Body.Close()

	var p *Pool
	if err := decodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// DeletePoolInput is used as input to the DeletePool function.
type DeletePoolInput struct {
	// Name is the name of the pool to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
