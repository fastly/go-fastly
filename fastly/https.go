package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// HTTPS represents an HTTPS Logging response from the Fastly API.
type HTTPS struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Format            string     `mapstructure:"format"`
	URL               string     `mapstructure:"url"`
	RequestMaxEntries uint       `mapstructure:"request_max_entries"`
	RequestMaxBytes   uint       `mapstructure:"request_max_bytes"`
	ContentType       string     `mapstructure:"content_type"`
	HeaderName        string     `mapstructure:"header_name"`
	HeaderValue       string     `mapstructure:"header_value"`
	Method            string     `mapstructure:"method"`
	JSONFormat        string     `mapstructure:"json_format"`
	Placement         string     `mapstructure:"placement"`
	TLSCACert         string     `mapstructure:"tls_ca_cert"`
	TLSClientCert     string     `mapstructure:"tls_client_cert"`
	TLSClientKey      string     `mapstructure:"tls_client_key"`
	TLSHostname       string     `mapstructure:"tls_hostname"`
	MessageType       string     `mapstructure:"message_type"`
	FormatVersion     uint       `mapstructure:"format_version"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// httpsByName is a sortable list of HTTPS logs.
type httpsByName []*HTTPS

// Len, Swap, and Less implement the sortable interface.
func (s httpsByName) Len() int      { return len(s) }
func (s httpsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s httpsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListHTTPSInput is used as input to the ListHTTPS function.
type ListHTTPSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListHTTPS returns the list of HTTPS logs for the configuration version.
func (c *Client) ListHTTPS(i *ListHTTPSInput) ([]*HTTPS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/https", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var https []*HTTPS
	if err := decodeBodyMap(resp.Body, &https); err != nil {
		return nil, err
	}
	sort.Stable(httpsByName(https))
	return https, nil
}

// CreateHTTPSInput is used as input to the CreateHTTPS function.
type CreateHTTPSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `url:"name,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	Format            string `url:"format,omitempty"`
	URL               string `url:"url,omitempty"`
	RequestMaxEntries uint   `url:"request_max_entries,omitempty"`
	RequestMaxBytes   uint   `url:"request_max_bytes,omitempty"`
	ContentType       string `url:"content_type,omitempty"`
	HeaderName        string `url:"header_name,omitempty"`
	HeaderValue       string `url:"header_value,omitempty"`
	Method            string `url:"method,omitempty"`
	JSONFormat        string `url:"json_format,omitempty"`
	Placement         string `url:"placement,omitempty"`
	TLSCACert         string `url:"tls_ca_cert,omitempty"`
	TLSClientCert     string `url:"tls_client_cert,omitempty"`
	TLSClientKey      string `url:"tls_client_key,omitempty"`
	TLSHostname       string `url:"tls_hostname,omitempty"`
	MessageType       string `url:"message_type,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
}

// CreateHTTPS creates a new Fastly HTTPS logging endpoint.
func (c *Client) CreateHTTPS(i *CreateHTTPSInput) (*HTTPS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/https", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var https *HTTPS
	if err := decodeBodyMap(resp.Body, &https); err != nil {
		return nil, err
	}
	return https, nil
}

// GetHTTPSInput is used as input to the GetHTTPS function.
type GetHTTPSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the HTTPS endpoint to fetch.
	Name string
}

func (c *Client) GetHTTPS(i *GetHTTPSInput) (*HTTPS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/https/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var h *HTTPS
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}

	return h, nil
}

type UpdateHTTPSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the HTTPS endpoint to fetch.
	Name string

	NewName           *string `url:"name,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	Format            *string `url:"format,omitempty"`
	URL               *string `url:"url,omitempty"`
	RequestMaxEntries *uint   `url:"request_max_entries,omitempty"`
	RequestMaxBytes   *uint   `url:"request_max_bytes,omitempty"`
	ContentType       *string `url:"content_type,omitempty"`
	HeaderName        *string `url:"header_name,omitempty"`
	HeaderValue       *string `url:"header_value,omitempty"`
	Method            *string `url:"method,omitempty"`
	JSONFormat        *string `url:"json_format,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	TLSCACert         *string `url:"tls_ca_cert,omitempty"`
	TLSClientCert     *string `url:"tls_client_cert,omitempty"`
	TLSClientKey      *string `url:"tls_client_key,omitempty"`
	TLSHostname       *string `url:"tls_hostname,omitempty"`
	MessageType       *string `url:"message_type,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
}

func (c *Client) UpdateHTTPS(i *UpdateHTTPSInput) (*HTTPS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/https/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var h *HTTPS
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

type DeleteHTTPSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the HTTPS endpoint to fetch.
	Name string
}

func (c *Client) DeleteHTTPS(i *DeleteHTTPSInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/https/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
