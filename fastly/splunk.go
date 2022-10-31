package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Splunk represents a splunk response from the Fastly API.
type Splunk struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	Name              string     `mapstructure:"name"`
	Placement         string     `mapstructure:"placement"`
	RequestMaxBytes   uint       `mapstructure:"request_max_bytes"`
	RequestMaxEntries uint       `mapstructure:"request_max_entries"`
	ResponseCondition string     `mapstructure:"response_condition"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	TLSCACert         string     `mapstructure:"tls_ca_cert"`
	TLSClientCert     string     `mapstructure:"tls_client_cert"`
	TLSClientKey      string     `mapstructure:"tls_client_key"`
	TLSHostname       string     `mapstructure:"tls_hostname"`
	Token             string     `mapstructure:"token"`
	URL               string     `mapstructure:"url"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	UseTLS            bool       `mapstructure:"use_tls"`
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
	defer resp.Body.Close()

	var ss []*Splunk
	if err := decodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	sort.Stable(splunkByName(ss))
	return ss, nil
}

// CreateSplunkInput is used as input to the CreateSplunk function.
type CreateSplunkInput struct {
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	Name              string `url:"name,omitempty"`
	Placement         string `url:"placement,omitempty"`
	RequestMaxBytes   uint   `url:"request_max_bytes,omitempty"`
	RequestMaxEntries uint   `url:"request_max_entries,omitempty"`
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
	URL            string      `url:"url,omitempty"`
	UseTLS         Compatibool `url:"use_tls,omitempty"`
}

// CreateSplunk creates a new Fastly splunk.
func (c *Client) CreateSplunk(i *CreateSplunkInput) (*Splunk, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Token == "" {
		return nil, ErrMissingToken
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/splunk", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Splunk
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSplunkInput is used as input to the GetSplunk function.
type GetSplunkInput struct {
	// Name is the name of the splunk to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
	defer resp.Body.Close()

	var s *Splunk
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateSplunkInput is used as input to the UpdateSplunk function.
type UpdateSplunkInput struct {
	Format        *string `url:"format,omitempty"`
	FormatVersion *uint   `url:"format_version,omitempty"`
	// Name is the name of the splunk to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	RequestMaxBytes   *uint   `url:"request_max_bytes,omitempty"`
	RequestMaxEntries *uint   `url:"request_max_entries,omitempty"`
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
	URL            *string      `url:"url,omitempty"`
	UseTLS         *Compatibool `url:"use_tls,omitempty"`
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

	if i.Token != nil && *i.Token == "" {
		return nil, ErrTokenEmpty
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/splunk/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Splunk
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSplunkInput is the input parameter to DeleteSplunk.
type DeleteSplunkInput struct {
	// Name is the name of the splunk to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
