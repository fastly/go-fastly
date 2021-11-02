package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Splunk represents a splunk response from the Fastly API.
type Splunk struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	URL               string     `mapstructure:"url"`
	RequestMaxEntries uint       `mapstructure:"request_max_entries"`
	RequestMaxBytes   uint       `mapstructure:"request_max_bytes"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	Token             string     `mapstructure:"token"`
	UseTLS            bool       `mapstructure:"use_tls"`
	TLSCACert         string     `mapstructure:"tls_ca_cert"`
	TLSHostname       string     `mapstructure:"tls_hostname"`
	TLSClientCert     string     `mapstructure:"tls_client_cert"`
	TLSClientKey      string     `mapstructure:"tls_client_key"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// splunkByName is a sortable list of splunks.
type splunkByName []*Splunk

// Len, Swap, and Less implement the sortable interface.
func (s splunkByName) Len() int      { return len(s) }
func (s splunkByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s splunkByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListSplunksInput is used as input to the ListSplunks function.
type ListSplunksInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListSplunks returns the list of splunks for the configuration version.
func (c *Client) ListSplunks(i *ListSplunksInput) ([]*Splunk, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/splunk", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ss []*Splunk
	if err := decodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	sort.Stable(splunkByName(ss))
	return ss, nil
}

// CreateSplunkInput is used as input to the CreateSplunk function.
type CreateSplunkInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string      `url:"name,omitempty"`
	URL               string      `url:"url,omitempty"`
	RequestMaxEntries uint        `url:"request_max_entries,omitempty"`
	RequestMaxBytes   uint        `url:"request_max_bytes,omitempty"`
	Format            string      `url:"format,omitempty"`
	FormatVersion     uint        `url:"format_version,omitempty"`
	ResponseCondition string      `url:"response_condition,omitempty"`
	Placement         string      `url:"placement,omitempty"`
	Token             string      `url:"token,omitempty"`
	UseTLS            Compatibool `url:"use_tls,omitempty"`
	TLSCACert         string      `url:"tls_ca_cert,omitempty"`
	TLSHostname       string      `url:"tls_hostname,omitempty"`
	TLSClientCert     string      `url:"tls_client_cert,omitempty"`
	TLSClientKey      string      `url:"tls_client_key,omitempty"`
}

// CreateSplunk creates a new Fastly splunk.
func (c *Client) CreateSplunk(i *CreateSplunkInput) (*Splunk, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/splunk", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Splunk
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSplunkInput is used as input to the GetSplunk function.
type GetSplunkInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the splunk to fetch.
	Name string
}

// GetSplunk gets the splunk configuration with the given parameters.
func (c *Client) GetSplunk(i *GetSplunkInput) (*Splunk, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/splunk/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var s *Splunk
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateSplunkInput is used as input to the UpdateSplunk function.
type UpdateSplunkInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the splunk to update.
	Name string

	NewName           *string      `url:"name,omitempty"`
	URL               *string      `url:"url,omitempty"`
	RequestMaxEntries *uint        `url:"request_max_entries,omitempty"`
	RequestMaxBytes   *uint        `url:"request_max_bytes,omitempty"`
	Format            *string      `url:"format,omitempty"`
	FormatVersion     *uint        `url:"format_version,omitempty"`
	ResponseCondition *string      `url:"response_condition,omitempty"`
	Placement         *string      `url:"placement,omitempty"`
	Token             *string      `url:"token,omitempty"`
	UseTLS            *Compatibool `url:"use_tls,omitempty"`
	TLSCACert         *string      `url:"tls_ca_cert,omitempty"`
	TLSHostname       *string      `url:"tls_hostname,omitempty"`
	TLSClientCert     *string      `url:"tls_client_cert,omitempty"`
	TLSClientKey      *string      `url:"tls_client_key,omitempty"`
}

// UpdateSplunk updates a specific splunk.
func (c *Client) UpdateSplunk(i *UpdateSplunkInput) (*Splunk, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/splunk/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Splunk
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSplunkInput is the input parameter to DeleteSplunk.
type DeleteSplunkInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the splunk to delete (required).
	Name string
}

// DeleteSplunk deletes the given splunk version.
func (c *Client) DeleteSplunk(i *DeleteSplunkInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/splunk/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
