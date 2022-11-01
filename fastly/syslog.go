package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Syslog represents a syslog response from the Fastly API.
type Syslog struct {
	Address           string     `mapstructure:"address"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	Hostname          string     `mapstructure:"hostname"`
	IPV4              string     `mapstructure:"ipv4"`
	MessageType       string     `mapstructure:"message_type"`
	Name              string     `mapstructure:"name"`
	Placement         string     `mapstructure:"placement"`
	Port              uint       `mapstructure:"port"`
	ResponseCondition string     `mapstructure:"response_condition"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	TLSCACert         string     `mapstructure:"tls_ca_cert"`
	TLSClientCert     string     `mapstructure:"tls_client_cert"`
	TLSClientKey      string     `mapstructure:"tls_client_key"`
	TLSHostname       string     `mapstructure:"tls_hostname"`
	Token             string     `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	UseTLS            bool       `mapstructure:"use_tls"`
}

// syslogsByName is a sortable list of syslogs.
type syslogsByName []*Syslog

// Len implement the sortable interface.
func (s syslogsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s syslogsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
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
	Address           string `url:"address,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	Hostname          string `url:"hostname,omitempty"`
	IPV4              string `url:"ipv4,omitempty"`
	MessageType       string `url:"message_type,omitempty"`
	Name              string `url:"name,omitempty"`
	Placement         string `url:"placement,omitempty"`
	Port              uint   `url:"port,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	TLSCACert      string      `url:"tls_ca_cert,omitempty"`
	TLSClientCert  string      `url:"tls_client_cert,omitempty"`
	TLSClientKey   string      `url:"tls_client_key,omitempty"`
	TLSHostname    string      `url:"tls_hostname,omitempty"`
	Token          string      `url:"token,omitempty"`
	UseTLS         Compatibool `url:"use_tls,omitempty"`
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
	// Name is the name of the syslog to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
	Address       *string `url:"address,omitempty"`
	Format        *string `url:"format,omitempty"`
	FormatVersion *uint   `url:"format_version,omitempty"`
	Hostname      *string `url:"hostname,omitempty"`
	IPV4          *string `url:"ipv4,omitempty"`
	MessageType   *string `url:"message_type,omitempty"`
	// Name is the name of the syslog to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	Port              *uint   `url:"port,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	TLSCACert      *string      `url:"tls_ca_cert,omitempty"`
	TLSClientCert  *string      `url:"tls_client_cert,omitempty"`
	TLSClientKey   *string      `url:"tls_client_key,omitempty"`
	TLSHostname    *string      `url:"tls_hostname,omitempty"`
	Token          *string      `url:"token,omitempty"`
	UseTLS         *Compatibool `url:"use_tls,omitempty"`
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
	// Name is the name of the syslog to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
