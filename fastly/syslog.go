package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Syslog represents a syslog response from the Fastly API.
type Syslog struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Address           string     `mapstructure:"address"`
	Hostname          string     `mapstructure:"hostname"`
	Port              uint       `mapstructure:"port"`
	UseTLS            bool       `mapstructure:"use_tls"`
	IPV4              string     `mapstructure:"ipv4"`
	TLSCACert         string     `mapstructure:"tls_ca_cert"`
	TLSHostname       string     `mapstructure:"tls_hostname"`
	TLSClientCert     string     `mapstructure:"tls_client_cert"`
	TLSClientKey      string     `mapstructure:"tls_client_key"`
	Token             string     `mapstructure:"token"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	MessageType       string     `mapstructure:"message_type"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// syslogsByName is a sortable list of syslogs.
type syslogsByName []*Syslog

// Len, Swap, and Less implement the sortable interface.
func (s syslogsByName) Len() int      { return len(s) }
func (s syslogsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s syslogsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListSyslogsInput is used as input to the ListSyslogs function.
type ListSyslogsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListSyslogs returns the list of syslogs for the configuration version.
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
	sort.Stable(syslogsByName(ss))
	return ss, nil
}

// CreateSyslogInput is used as input to the CreateSyslog function.
type CreateSyslogInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string      `url:"name,omitempty"`
	Address           string      `url:"address,omitempty"`
	Hostname          string      `url:"hostname,omitempty"`
	Port              uint        `url:"port,omitempty"`
	UseTLS            Compatibool `url:"use_tls,omitempty"`
	IPV4              string      `url:"ipv4,omitempty"`
	TLSCACert         string      `url:"tls_ca_cert,omitempty"`
	TLSHostname       string      `url:"tls_hostname,omitempty"`
	TLSClientCert     string      `url:"tls_client_cert,omitempty"`
	TLSClientKey      string      `url:"tls_client_key,omitempty"`
	Token             string      `url:"token,omitempty"`
	Format            string      `url:"format,omitempty"`
	FormatVersion     uint        `url:"format_version,omitempty"`
	MessageType       string      `url:"message_type,omitempty"`
	ResponseCondition string      `url:"response_condition,omitempty"`
	Placement         string      `url:"placement,omitempty"`
}

// CreateSyslog creates a new Fastly syslog.
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the syslog to fetch.
	Name string
}

// GetSyslog gets the syslog configuration with the given parameters.
func (c *Client) GetSyslog(i *GetSyslogInput) (*Syslog, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the syslog to update.
	Name string

	NewName           *string      `url:"name,omitempty"`
	Address           *string      `url:"address,omitempty"`
	Hostname          *string      `url:"hostname,omitempty"`
	Port              *uint        `url:"port,omitempty"`
	UseTLS            *Compatibool `url:"use_tls,omitempty"`
	IPV4              *string      `url:"ipv4,omitempty"`
	TLSCACert         *string      `url:"tls_ca_cert,omitempty"`
	TLSHostname       *string      `url:"tls_hostname,omitempty"`
	TLSClientCert     *string      `url:"tls_client_cert,omitempty"`
	TLSClientKey      *string      `url:"tls_client_key,omitempty"`
	Token             *string      `url:"token,omitempty"`
	Format            *string      `url:"format,omitempty"`
	FormatVersion     *uint        `url:"format_version,omitempty"`
	MessageType       *string      `url:"message_type,omitempty"`
	ResponseCondition *string      `url:"response_condition,omitempty"`
	Placement         *string      `url:"placement,omitempty"`
}

// UpdateSyslog updates a specific syslog.
func (c *Client) UpdateSyslog(i *UpdateSyslogInput) (*Syslog, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the syslog to delete (required).
	Name string
}

// DeleteSyslog deletes the given syslog version.
func (c *Client) DeleteSyslog(i *DeleteSyslogInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
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
