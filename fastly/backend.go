package fastly

import (
	"fmt"
	"strconv"
	"time"
)

// Backend represents a backend response from the Fastly API.
type Backend struct {
	Address             *string    `mapstructure:"address"`
	AutoLoadbalance     *bool      `mapstructure:"auto_loadbalance"`
	BetweenBytesTimeout *int       `mapstructure:"between_bytes_timeout"`
	Comment             *string    `mapstructure:"comment"`
	ConnectTimeout      *int       `mapstructure:"connect_timeout"`
	CreatedAt           *time.Time `mapstructure:"created_at"`
	DeletedAt           *time.Time `mapstructure:"deleted_at"`
	ErrorThreshold      *int       `mapstructure:"error_threshold"`
	FirstByteTimeout    *int       `mapstructure:"first_byte_timeout"`
	HealthCheck         *string    `mapstructure:"healthcheck"`
	Hostname            *string    `mapstructure:"hostname"`
	KeepAliveTime       *int       `mapstructure:"keepalive_time"`
	MaxConn             *int       `mapstructure:"max_conn"`
	MaxTLSVersion       *string    `mapstructure:"max_tls_version"`
	MinTLSVersion       *string    `mapstructure:"min_tls_version"`
	Name                *string    `mapstructure:"name"`
	OverrideHost        *string    `mapstructure:"override_host"`
	Port                *int       `mapstructure:"port"`
	RequestCondition    *string    `mapstructure:"request_condition"`
	ShareKey            *string    `mapstructure:"share_key"`
	SSLCACert           *string    `mapstructure:"ssl_ca_cert"`
	SSLCertHostname     *string    `mapstructure:"ssl_cert_hostname"`
	SSLCheckCert        *bool      `mapstructure:"ssl_check_cert"`
	SSLCiphers          *string    `mapstructure:"ssl_ciphers"`
	SSLClientCert       *string    `mapstructure:"ssl_client_cert"`
	SSLClientKey        *string    `mapstructure:"ssl_client_key"`
	SSLSNIHostname      *string    `mapstructure:"ssl_sni_hostname"`
	ServiceID           *string    `mapstructure:"service_id"`
	ServiceVersion      *int       `mapstructure:"version"`
	Shield              *string    `mapstructure:"shield"`
	TCPKeepAliveEnable  *bool      `mapstructure:"tcp_keepalive_enable"`
	TCPKeepAliveIntvl   *int       `mapstructure:"tcp_keepalive_interval"`
	TCPKeepAliveProbes  *int       `mapstructure:"tcp_keepalive_probes"`
	TCPKeepAliveTime    *int       `mapstructure:"tcp_keepalive_time"`
	UpdatedAt           *time.Time `mapstructure:"updated_at"`
	UseSSL              *bool      `mapstructure:"use_ssl"`
	Weight              *int       `mapstructure:"weight"`
}

// ListBackendsInput is used as input to the ListBackends function.
type ListBackendsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListBackends retrieves all resources.
func (c *Client) ListBackends(i *ListBackendsInput) ([]*Backend, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "backend")

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bs []*Backend
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	return bs, nil
}

// CreateBackendInput is used as input to the CreateBackend function.
type CreateBackendInput struct {
	// Address is a hostname, IPv4, or IPv6 address for the backend.
	Address *string `url:"address,omitempty"`
	// AutoLoadbalance is whether or not this backend should be automatically load balanced.
	AutoLoadbalance *Compatibool `url:"auto_loadbalance,omitempty"`
	// BetweenBytesTimeout is the maximum duration in milliseconds that Fastly will wait while receiving no data on a download from a backend.
	BetweenBytesTimeout *int `url:"between_bytes_timeout,omitempty"`
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// ConnectTimeout is the maximum duration in milliseconds to wait for a connection to this backend to be established.
	ConnectTimeout *int `url:"connect_timeout,omitempty"`
	// ErrorThreshold is the number of errors to allow before the Backend is marked as down.
	ErrorThreshold *int `url:"error_threshold,omitempty"`
	// FirstByteTimeout is how long to wait for the first bytes in milliseconds.
	FirstByteTimeout *int `url:"first_byte_timeout,omitempty"`
	// HealthCheck is the name of the healthcheck to use with this backend.
	HealthCheck *string `url:"healthcheck,omitempty"`
	// KeepAliveTime is how long in seconds to keep a persistent connection to the backend between requests.
	KeepAliveTime *int `url:"keepalive_time,omitempty"`
	// MaxConn is the maximum number of concurrent connections this backend will accept.
	MaxConn *int `url:"max_conn,omitempty"`
	// MaxTLSVersion is the maximum allowed TLS version on SSL connections to this backend.
	MaxTLSVersion *string `url:"max_tls_version,omitempty"`
	// MinTLSVersion is the minimum allowed TLS version on SSL connections to this backend.
	MinTLSVersion *string `url:"min_tls_version,omitempty"`
	// Name is the name of the backend.
	Name *string `url:"name,omitempty"`
	// OverrideHost is, if set, will replace the client-supplied HTTP Host header on connections to this backend.
	OverrideHost *string `url:"override_host,omitempty"`
	// Port is the port on which the backend server is listening for connections from Fastly.
	Port *int `url:"port,omitempty"`
	// RequestCondition is the name of a Condition, which if satisfied, will select this backend during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// ShareKey is a value that when shared across backends will enable those backends to share the same health check.
	ShareKey *string `url:"share_key,omitempty"`
	// SSLCACert is a CA certificate attached to origin.
	SSLCACert *string `url:"ssl_ca_cert,omitempty"`
	// SSLCertHostname is an overrides ssl_hostname, but only for cert verification.
	SSLCertHostname *string `url:"ssl_cert_hostname,omitempty"`
	// SSLCheckCert forces being strict on checking SSL certs.
	SSLCheckCert *Compatibool `url:"ssl_check_cert,omitempty"`
	// SSLCiphers is a list of OpenSSL ciphers to support for connections to this origin.
	SSLCiphers *string `url:"ssl_ciphers,omitempty"`
	// SSLClientCert is a client certificate attached to origin.
	SSLClientCert *string `url:"ssl_client_cert,omitempty"`
	// SSLClientKey is a client key attached to origin.
	SSLClientKey *string `url:"ssl_client_key,omitempty"`
	// SSLSNIHostname overrides ssl_hostname, but only for SNI in the handshake. Does not affect cert validation at all.
	SSLSNIHostname *string `url:"ssl_sni_hostname,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Shield is an identifier of the POP to use as a shield.
	Shield *string `url:"shield,omitempty"`
	// TCPKeepAliveEnable is whether to enable TCP keepalives to this backend.
	TCPKeepAliveEnable *bool `url:"tcp_keepalive_enable,omitempty"`
	// TCPKeepAliveIntvl is how long to wait between each TCP keepalive probe sent to this backend.
	TCPKeepAliveIntvl *int `url:"tcp_keepalive_interval,omitempty"`
	// TCPKeepAliveProbes is how many unacknowledged TCP keepalive probes to send to this backend before it's considered dead.
	TCPKeepAliveProbes *int `url:"tcp_keepalive_probes,omitempty"`
	// TCPKeepAliveTime is how long to wait after the last sent data before sending TCP keepalive probes.
	TCPKeepAliveTime *int `url:"tcp_keepalive_time,omitempty"`
	// UseSSL indicates whether or not to require TLS for connections to this backend.
	UseSSL *Compatibool `url:"use_ssl,omitempty"`
	// Weight is the weight used to load balance this backend against others.
	Weight *int `url:"weight,omitempty"`
}

// CreateBackend creates a new resource.
func (c *Client) CreateBackend(i *CreateBackendInput) (*Backend, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "backend")

	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Backend
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetBackendInput is used as input to the GetBackend function.
type GetBackendInput struct {
	// Name is the name of the backend to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetBackend retrieves the specified resource.
func (c *Client) GetBackend(i *GetBackendInput) (*Backend, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "backend", i.Name)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Backend
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateBackendInput is used as input to the UpdateBackend function.
type UpdateBackendInput struct {
	// Address is a hostname, IPv4, or IPv6 address for the backend.
	Address *string `url:"address,omitempty"`
	// AutoLoadbalance is whether or not this backend should be automatically load balanced.
	AutoLoadbalance *Compatibool `url:"auto_loadbalance,omitempty"`
	// BetweenBytesTimeout is the maximum duration in milliseconds that Fastly will wait while receiving no data on a download from a backend.
	BetweenBytesTimeout *int `url:"between_bytes_timeout,omitempty"`
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// ConnectTimeout is the maximum duration in milliseconds to wait for a connection to this backend to be established.
	ConnectTimeout *int `url:"connect_timeout,omitempty"`
	// ErrorThreshold is the number of errors to allow before the Backend is marked as down.
	ErrorThreshold *int `url:"error_threshold,omitempty"`
	// FirstByteTimeout is how long to wait for the first bytes in milliseconds.
	FirstByteTimeout *int `url:"first_byte_timeout,omitempty"`
	// HealthCheck is the name of the healthcheck to use with this backend.
	HealthCheck *string `url:"healthcheck,omitempty"`
	// KeepAliveTime is how long in seconds to keep a persistent connection to the backend between requests.
	KeepAliveTime *int `url:"keepalive_time,omitempty"`
	// MaxConn is the maximum number of concurrent connections this backend will accept.
	MaxConn *int `url:"max_conn,omitempty"`
	// MaxTLSVersion is the maximum allowed TLS version on SSL connections to this backend.
	MaxTLSVersion *string `url:"max_tls_version,omitempty"`
	// MinTLSVersion is the minimum allowed TLS version on SSL connections to this backend.
	MinTLSVersion *string `url:"min_tls_version,omitempty"`
	// Name is the name of the backend to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// OverrideHost is, if set, will replace the client-supplied HTTP Host header on connections to this backend.
	OverrideHost *string `url:"override_host,omitempty"`
	// Port is the port on which the backend server is listening for connections from Fastly.
	Port *int `url:"port,omitempty"`
	// RequestCondition is the name of a Condition, which if satisfied, will select this backend during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// ShareKey is a value that when shared across backends will enable those backends to share the same health check.
	ShareKey *string `url:"share_key,omitempty"`
	// SSLCACert is a CA certificate attached to origin.
	SSLCACert *string `url:"ssl_ca_cert,omitempty"`
	// SSLCertHostname is an overrides ssl_hostname, but only for cert verification.
	SSLCertHostname *string `url:"ssl_cert_hostname,omitempty"`
	// SSLCheckCert forces being strict on checking SSL certs.
	SSLCheckCert *Compatibool `url:"ssl_check_cert,omitempty"`
	// SSLCiphers is a list of OpenSSL ciphers to support for connections to this origin.
	SSLCiphers *string `url:"ssl_ciphers,omitempty"`
	// SSLClientCert is a client certificate attached to origin.
	SSLClientCert *string `url:"ssl_client_cert,omitempty"`
	// SSLClientKey is a client key attached to origin.
	SSLClientKey *string `url:"ssl_client_key,omitempty"`
	// SSLSNIHostname overrides ssl_hostname, but only for SNI in the handshake. Does not affect cert validation at all.
	SSLSNIHostname *string `url:"ssl_sni_hostname,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Shield is an identifier of the POP to use as a shield.
	Shield *string `url:"shield,omitempty"`
	// TCPKeepAliveEnable is whether to enable TCP keepalives to this backend.
	TCPKeepAliveEnable *bool `url:"tcp_keepalive_enable,omitempty"`
	// TCPKeepAliveIntvl is how long to wait between each TCP keepalive probe sent to this backend.
	TCPKeepAliveIntvl *int `url:"tcp_keepalive_interval,omitempty"`
	// TCPKeepAliveProbes is how many unacknowledged TCP keepalive probes to send to this backend before it's considered dead.
	TCPKeepAliveProbes *int `url:"tcp_keepalive_probes,omitempty"`
	// TCPKeepAliveTime is how long to wait after the last sent data before sending TCP keepalive probes.
	TCPKeepAliveTime *int `url:"tcp_keepalive_time,omitempty"`
	// UseSSL indicates whether or not to require TLS for connections to this backend.
	UseSSL *Compatibool `url:"use_ssl,omitempty"`
	// Weight is the weight used to load balance this backend against others.
	Weight *int `url:"weight,omitempty"`
}

// UpdateBackend updates the specified resource.
func (c *Client) UpdateBackend(i *UpdateBackendInput) (*Backend, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "backend", i.Name)

	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Backend
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteBackendInput is the input parameter to DeleteBackend.
type DeleteBackendInput struct {
	// Name is the name of the backend to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteBackend deletes the specified resource.
func (c *Client) DeleteBackend(i *DeleteBackendInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "backend", i.Name)

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
		return fmt.Errorf("not ok")
	}
	return nil
}
