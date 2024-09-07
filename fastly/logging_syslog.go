package fastly

import (
	"fmt"
	"net/url"
	"time"
)

// Syslog represents a syslog response from the Fastly API.
type Syslog struct {
	Address           *string    `mapstructure:"address"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Hostname          *string    `mapstructure:"hostname"`
	IPV4              *string    `mapstructure:"ipv4"`
	MessageType       *string    `mapstructure:"message_type"`
	Name              *string    `mapstructure:"name"`
	Placement         *string    `mapstructure:"placement"`
	Port              *int       `mapstructure:"port"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	TLSCACert         *string    `mapstructure:"tls_ca_cert"`
	TLSClientCert     *string    `mapstructure:"tls_client_cert"`
	TLSClientKey      *string    `mapstructure:"tls_client_key"`
	TLSHostname       *string    `mapstructure:"tls_hostname"`
	Token             *string    `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	UseTLS            *bool      `mapstructure:"use_tls"`
}

// ListSyslogsInput is used as input to the ListSyslogs function.
type ListSyslogsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListSyslogs retrieves all resources.
func (c *Client) ListSyslogs(i *ListSyslogsInput) ([]*Syslog, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/syslog", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ss []*Syslog
	if err := decodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	return ss, nil
}

// CreateSyslogInput is used as input to the CreateSyslog function.
type CreateSyslogInput struct {
	// Address is a hostname or IPv4 address.
	Address *string `url:"address,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Hostname is the hostname used for the syslog endpoint.
	Hostname *string `url:"hostname,omitempty"`
	// IPV4 is the IPv4 address used for the syslog endpoint.
	IPV4 *string `url:"ipv4,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Port is the port number.
	Port *int `url:"port,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TLSCACert is a secure certificate to authenticate a server with. Must be in PEM format.
	TLSCACert *string `url:"tls_ca_cert,omitempty"`
	// TLSClientCert is the client certificate used to make authenticated requests. Must be in PEM format.
	TLSClientCert *string `url:"tls_client_cert,omitempty"`
	// TLSClientKey is the client private key used to make authenticated requests. Must be in PEM format.
	TLSClientKey *string `url:"tls_client_key,omitempty"`
	// TLSHostname is the hostname to verify the server's certificate. This should be one of the Subject Alternative Name (SAN) fields for the certificate. Common Names (CN) are not supported.
	TLSHostname *string `url:"tls_hostname,omitempty"`
	// Token is whether to prepend each message with a specific token.
	Token *string `url:"token,omitempty"`
	// UseTLS is whether to use TLS (0: do not use, 1: use).
	UseTLS *Compatibool `url:"use_tls,omitempty"`
}

// CreateSyslog creates a new resource.
func (c *Client) CreateSyslog(i *CreateSyslogInput) (*Syslog, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/syslog", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Syslog
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSyslogInput is used as input to the GetSyslog function.
type GetSyslogInput struct {
	// Name is the name of the syslog to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetSyslog retrieves the specified resource.
func (c *Client) GetSyslog(i *GetSyslogInput) (*Syslog, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/syslog/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Syslog
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateSyslogInput is used as input to the UpdateSyslog function.
type UpdateSyslogInput struct {
	// Address is a hostname or IPv4 address.
	Address *string `url:"address,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Hostname is the hostname used for the syslog endpoint.
	Hostname *string `url:"hostname,omitempty"`
	// IPV4 is the IPv4 address used for the syslog endpoint.
	IPV4 *string `url:"ipv4,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name of the syslog to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Port is the port number.
	Port *int `url:"port,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TLSCACert is a secure certificate to authenticate a server with. Must be in PEM format.
	TLSCACert *string `url:"tls_ca_cert,omitempty"`
	// TLSClientCert is the client certificate used to make authenticated requests. Must be in PEM format.
	TLSClientCert *string `url:"tls_client_cert,omitempty"`
	// TLSClientKey is the client private key used to make authenticated requests. Must be in PEM format.
	TLSClientKey *string `url:"tls_client_key,omitempty"`
	// TLSHostname is the hostname to verify the server's certificate. This should be one of the Subject Alternative Name (SAN) fields for the certificate. Common Names (CN) are not supported.
	TLSHostname *string `url:"tls_hostname,omitempty"`
	// Token is whether to prepend each message with a specific token.
	Token *string `url:"token,omitempty"`
	// UseTLS is whether to use TLS (0: do not use, 1: use).
	UseTLS *Compatibool `url:"use_tls,omitempty"`
}

// UpdateSyslog updates the specified resource.
func (c *Client) UpdateSyslog(i *UpdateSyslogInput) (*Syslog, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/syslog/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Syslog
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSyslogInput is the input parameter to DeleteSyslog.
type DeleteSyslogInput struct {
	// Name is the name of the syslog to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteSyslog deletes the specified resource.
func (c *Client) DeleteSyslog(i *DeleteSyslogInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/syslog/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
