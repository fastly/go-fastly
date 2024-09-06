package fastly

import (
	"strconv"
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

// Pool represents a pool response from the Fastly API.
type Pool struct {
	Comment          *string    `mapstructure:"comment"`
	ConnectTimeout   *int       `mapstructure:"connect_timeout"`
	CreatedAt        *time.Time `mapstructure:"created_at"`
	DeletedAt        *time.Time `mapstructure:"deleted_at"`
	FirstByteTimeout *int       `mapstructure:"first_byte_timeout"`
	Healthcheck      *string    `mapstructure:"healthcheck"`
	MaxConnDefault   *int       `mapstructure:"max_conn_default"`
	MaxTLSVersion    *string    `mapstructure:"max_tls_version"`
	MinTLSVersion    *string    `mapstructure:"min_tls_version"`
	Name             *string    `mapstructure:"name"`
	OverrideHost     *string    `mapstructure:"override_host"`
	PoolID           *string    `mapstructure:"id"`
	Quorum           *int       `mapstructure:"quorum"`
	RequestCondition *string    `mapstructure:"request_condition"`
	ServiceID        *string    `mapstructure:"service_id"`
	ServiceVersion   *int       `mapstructure:"version"`
	Shield           *string    `mapstructure:"shield"`
	TLSCACert        *string    `mapstructure:"tls_ca_cert"`
	TLSCertHostname  *string    `mapstructure:"tls_cert_hostname"`
	TLSCheckCert     *bool      `mapstructure:"tls_check_cert"`
	TLSCiphers       *string    `mapstructure:"tls_ciphers"`
	TLSClientCert    *string    `mapstructure:"tls_client_cert"`
	TLSClientKey     *string    `mapstructure:"tls_client_key"`
	TLSSNIHostname   *string    `mapstructure:"tls_sni_hostname"`
	Type             *PoolType  `mapstructure:"type"`
	UpdatedAt        *time.Time `mapstructure:"updated_at"`
	UseTLS           *bool      `mapstructure:"use_tls"`
}

// ListPoolsInput is used as input to the ListPools function.
type ListPoolsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListPools retrieves all resources.
func (c *Client) ListPools(i *ListPoolsInput) ([]*Pool, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "pool")

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ps []*Pool
	if err := decodeBodyMap(resp.Body, &ps); err != nil {
		return nil, err
	}
	return ps, nil
}

// CreatePoolInput is used as input to the CreatePool function.
type CreatePoolInput struct {
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// ConnectTimeout is how long to wait for a timeout in milliseconds.
	ConnectTimeout *int `url:"connect_timeout,omitempty"`
	// FirstByteTimeout is how long to wait for the first byte in milliseconds.
	FirstByteTimeout *int `url:"first_byte_timeout,omitempty"`
	// Healthcheck is the name of the healthcheck to use with this pool.
	Healthcheck *string `url:"healthcheck,omitempty"`
	// MaxConnDefault is the maximum number of connections.
	MaxConnDefault *int `url:"max_conn_default,omitempty"`
	// MaxTLSVersion is the maximum allowed TLS version on connections to this server.
	MaxTLSVersion *string `url:"max_tls_version,omitempty"`
	// MinTLSVersion is the minimum allowed TLS version on connections to this server.
	MinTLSVersion *string `url:"min_tls_version,omitempty"`
	// Name is the name of the pool to create (required).
	Name *string `url:"name,omitempty"`
	// OverrideHost is the hostname to override the Host header.
	OverrideHost *string `url:"override_host,omitempty"`
	// Quorum is the percentage of capacity (0-100) that needs to be operationally available for a pool to be considered up.
	Quorum *int `url:"quorum,omitempty"`
	// RequestCondition is the condition which, if met, will select this configuration during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Shield is the selected POP to serve as a shield for the servers.
	Shield *string `url:"shield,omitempty"`
	// TLSCACert is a secure certificate to authenticate a server with. Must be in PEM format.
	TLSCACert *string `url:"tls_ca_cert,omitempty"`
	// TLSCertHostname is the hostname used to verify a server's certificate.
	TLSCertHostname *string `url:"tls_cert_hostname,omitempty"`
	// TLSCheckCert forces strict checking of TLS certs.
	TLSCheckCert *Compatibool `url:"tls_check_cert,omitempty"`
	// TLSCiphers is a list of OpenSSL ciphers (see the openssl.org manpages for details).
	TLSCiphers *string `url:"tls_ciphers,omitempty"`
	// TLSClientCert is the client certificate used to make authenticated requests. Must be in PEM format.
	TLSClientCert *string `url:"tls_client_cert,omitempty"`
	// TLSClientKey is the client private key used to make authenticated requests. Must be in PEM format.
	TLSClientKey *string `url:"tls_client_key,omitempty"`
	// TLSSNIHostname is the SNI hostname.
	TLSSNIHostname *string `url:"tls_sni_hostname,omitempty"`
	// Type is what type of load balance group to use (random, hash, client).
	Type *PoolType `url:"type,omitempty"`
	// UseTLS indicates whether to use TLS.
	UseTLS *Compatibool `url:"use_tls,omitempty"`
}

// CreatePool creates a new resource.
func (c *Client) CreatePool(i *CreatePoolInput) (*Pool, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "pool")
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

// GetPool retrieves the specified resource.
func (c *Client) GetPool(i *GetPoolInput) (*Pool, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "pool", i.Name)

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
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// ConnectTimeout is how long to wait for a timeout in milliseconds.
	ConnectTimeout *int `url:"connect_timeout,omitempty"`
	// FirstByteTimeout is how long to wait for the first byte in milliseconds.
	FirstByteTimeout *int `url:"first_byte_timeout,omitempty"`
	// Healthcheck is the name of the healthcheck to use with this pool.
	Healthcheck *string `url:"healthcheck,omitempty"`
	// MaxConnDefault is the maximum number of connections.
	MaxConnDefault *int `url:"max_conn_default,omitempty"`
	// MaxTLSVersion is the maximum allowed TLS version on connections to this server.
	MaxTLSVersion *string `url:"max_tls_version,omitempty"`
	// MinTLSVersion is the minimum allowed TLS version on connections to this server.
	MinTLSVersion *string `url:"min_tls_version,omitempty"`
	// Name is the name of the pool to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// OverrideHost is the hostname to override the Host header.
	OverrideHost *string `url:"override_host,omitempty"`
	// Quorum is the percentage of capacity (0-100) that needs to be operationally available for a pool to be considered up.
	Quorum *int `url:"quorum,omitempty"`
	// RequestCondition is the condition which, if met, will select this configuration during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Shield is the selected POP to serve as a shield for the servers.
	Shield *string `url:"shield,omitempty"`
	// TLSCACert is a secure certificate to authenticate a server with. Must be in PEM format.
	TLSCACert *string `url:"tls_ca_cert,omitempty"`
	// TLSCertHostname is the hostname used to verify a server's certificate.
	TLSCertHostname *string `url:"tls_cert_hostname,omitempty"`
	// TLSCheckCert forces strict checking of TLS certs.
	TLSCheckCert *Compatibool `url:"tls_check_cert,omitempty"`
	// TLSCiphers is a list of OpenSSL ciphers (see the openssl.org manpages for details).
	TLSCiphers *string `url:"tls_ciphers,omitempty"`
	// TLSClientCert is the client certificate used to make authenticated requests. Must be in PEM format.
	TLSClientCert *string `url:"tls_client_cert,omitempty"`
	// TLSClientKey is the client private key used to make authenticated requests. Must be in PEM format.
	TLSClientKey *string `url:"tls_client_key,omitempty"`
	// TLSSNIHostname is the SNI hostname.
	TLSSNIHostname *string `url:"tls_sni_hostname,omitempty"`
	// Type is what type of load balance group to use (random, hash, client).
	Type *PoolType `url:"type,omitempty"`
	// UseTLS indicates whether to use TLS.
	UseTLS *Compatibool `url:"use_tls,omitempty"`
}

// UpdatePool updates the specified resource.
func (c *Client) UpdatePool(i *UpdatePoolInput) (*Pool, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "pool", i.Name)

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

// DeletePool deletes the specified resource.
func (c *Client) DeletePool(i *DeletePoolInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "pool", i.Name)

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
