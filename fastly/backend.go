package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Backend represents a backend response from the Fastly API.
type Backend struct {
	Address             string     `mapstructure:"address"`
	AutoLoadbalance     bool       `mapstructure:"auto_loadbalance"`
	BetweenBytesTimeout uint       `mapstructure:"between_bytes_timeout"`
	Comment             string     `mapstructure:"comment"`
	ConnectTimeout      uint       `mapstructure:"connect_timeout"`
	CreatedAt           *time.Time `mapstructure:"created_at"`
	DeletedAt           *time.Time `mapstructure:"deleted_at"`
	ErrorThreshold      uint       `mapstructure:"error_threshold"`
	FirstByteTimeout    uint       `mapstructure:"first_byte_timeout"`
	HealthCheck         string     `mapstructure:"healthcheck"`
	Hostname            string     `mapstructure:"hostname"`
	MaxConn             uint       `mapstructure:"max_conn"`
	MaxTLSVersion       string     `mapstructure:"max_tls_version"`
	MinTLSVersion       string     `mapstructure:"min_tls_version"`
	Name                string     `mapstructure:"name"`
	OverrideHost        string     `mapstructure:"override_host"`
	Port                uint       `mapstructure:"port"`
	RequestCondition    string     `mapstructure:"request_condition"`
	SSLCACert           string     `mapstructure:"ssl_ca_cert"`
	SSLCertHostname     string     `mapstructure:"ssl_cert_hostname"`
	SSLCheckCert        bool       `mapstructure:"ssl_check_cert"`
	SSLCiphers          string     `mapstructure:"ssl_ciphers"`
	SSLClientCert       string     `mapstructure:"ssl_client_cert"`
	SSLClientKey        string     `mapstructure:"ssl_client_key"`
	SSLHostname         string     `mapstructure:"ssl_hostname"`
	SSLSNIHostname      string     `mapstructure:"ssl_sni_hostname"`
	ServiceID           string     `mapstructure:"service_id"`
	ServiceVersion      int        `mapstructure:"version"`
	Shield              string     `mapstructure:"shield"`
	UpdatedAt           *time.Time `mapstructure:"updated_at"`
	UseSSL              bool       `mapstructure:"use_ssl"`
	Weight              uint       `mapstructure:"weight"`
}

// backendsByName is a sortable list of backends.
type backendsByName []*Backend

// Len, Swap, and Less implement the sortable interface.
func (s backendsByName) Len() int {
	return len(s)
}

func (s backendsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s backendsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListBackendsInput is used as input to the ListBackends function.
type ListBackendsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListBackends returns the list of backends for the configuration version.
func (c *Client) ListBackends(i *ListBackendsInput) ([]*Backend, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/backend", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bs []*Backend
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	sort.Stable(backendsByName(bs))
	return bs, nil
}

// CreateBackendInput is used as input to the CreateBackend function.
type CreateBackendInput struct {
	Address             string       `url:"address,omitempty"`
	AutoLoadbalance     Compatibool  `url:"auto_loadbalance,omitempty"`
	BetweenBytesTimeout *uint        `url:"between_bytes_timeout,omitempty"`
	Comment             string       `url:"comment,omitempty"`
	ConnectTimeout      *uint        `url:"connect_timeout,omitempty"`
	ErrorThreshold      *uint        `url:"error_threshold,omitempty"`
	FirstByteTimeout    *uint        `url:"first_byte_timeout,omitempty"`
	HealthCheck         string       `url:"healthcheck,omitempty"`
	MaxConn             *uint        `url:"max_conn,omitempty"`
	MaxTLSVersion       string       `url:"max_tls_version,omitempty"`
	MinTLSVersion       string       `url:"min_tls_version,omitempty"`
	Name                string       `url:"name,omitempty"`
	OverrideHost        string       `url:"override_host,omitempty"`
	Port                uint         `url:"port,omitempty"`
	RequestCondition    string       `url:"request_condition,omitempty"`
	SSLCACert           string       `url:"ssl_ca_cert,omitempty"`
	SSLCertHostname     string       `url:"ssl_cert_hostname,omitempty"`
	SSLCheckCert        *Compatibool `url:"ssl_check_cert,omitempty"`
	SSLCiphers          string       `url:"ssl_ciphers,omitempty"`
	SSLClientCert       string       `url:"ssl_client_cert,omitempty"`
	SSLClientKey        string       `url:"ssl_client_key,omitempty"`
	SSLHostname         string       `url:"ssl_hostname,omitempty"`
	SSLSNIHostname      string       `url:"ssl_sni_hostname,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Shield         string      `url:"shield,omitempty"`
	UseSSL         Compatibool `url:"use_ssl,omitempty"`
	Weight         *uint       `url:"weight,omitempty"`
}

// CreateBackend creates a new resource.
func (c *Client) CreateBackend(i *CreateBackendInput) (*Backend, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/backend", i.ServiceID, i.ServiceVersion)
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
	// Name is the name of the backend to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetBackend gets the backend configuration with the given parameters.
func (c *Client) GetBackend(i *GetBackendInput) (*Backend, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/backend/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
	Address             *string      `url:"address,omitempty"`
	AutoLoadbalance     *Compatibool `url:"auto_loadbalance,omitempty"`
	BetweenBytesTimeout *uint        `url:"between_bytes_timeout,omitempty"`
	Comment             *string      `url:"comment,omitempty"`
	ConnectTimeout      *uint        `url:"connect_timeout,omitempty"`
	ErrorThreshold      *uint        `url:"error_threshold,omitempty"`
	FirstByteTimeout    *uint        `url:"first_byte_timeout,omitempty"`
	HealthCheck         *string      `url:"healthcheck,omitempty"`
	MaxConn             *uint        `url:"max_conn,omitempty"`
	MaxTLSVersion       *string      `url:"max_tls_version,omitempty"`
	MinTLSVersion       *string      `url:"min_tls_version,omitempty"`
	// Name is the name of the backend to update.
	Name             string
	NewName          *string      `url:"name,omitempty"`
	OverrideHost     *string      `url:"override_host,omitempty"`
	Port             *uint        `url:"port,omitempty"`
	RequestCondition *string      `url:"request_condition,omitempty"`
	SSLCACert        *string      `url:"ssl_ca_cert,omitempty"`
	SSLCertHostname  *string      `url:"ssl_cert_hostname,omitempty"`
	SSLCheckCert     *Compatibool `url:"ssl_check_cert,omitempty"`
	SSLCiphers       *string      `url:"ssl_ciphers,omitempty"`
	SSLClientCert    *string      `url:"ssl_client_cert,omitempty"`
	SSLClientKey     *string      `url:"ssl_client_key,omitempty"`
	SSLSNIHostname   *string      `url:"ssl_sni_hostname,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Shield         *string      `url:"shield,omitempty"`
	UseSSL         *Compatibool `url:"use_ssl,omitempty"`
	Weight         *uint        `url:"weight,omitempty"`
}

// UpdateBackend updates a specific backend.
func (c *Client) UpdateBackend(i *UpdateBackendInput) (*Backend, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/backend/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/backend/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
