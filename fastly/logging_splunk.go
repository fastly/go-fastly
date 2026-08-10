package fastly

import (
	"context"
	"strconv"
	"time"
)

// Splunk represents a splunk response from the Fastly API.
type Splunk struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	Placement         *string    `mapstructure:"placement"`
	ProcessingRegion  *string    `mapstructure:"log_processing_region"`
	RequestMaxBytes   *int       `mapstructure:"request_max_bytes"`
	RequestMaxEntries *int       `mapstructure:"request_max_entries"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	TLSCACert         *string    `mapstructure:"tls_ca_cert"`
	TLSClientCert     *string    `mapstructure:"tls_client_cert"`
	TLSClientKey      *string    `mapstructure:"tls_client_key"`
	TLSHostname       *string    `mapstructure:"tls_hostname"`
	Token             *string    `mapstructure:"token"`
	URL               *string    `mapstructure:"url"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	UseTLS            *bool      `mapstructure:"use_tls"`
}

// ListSplunksInput is used as input to the ListSplunks function.
type ListSplunksInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListSplunks retrieves all resources.
func (c *Client) ListSplunks(ctx context.Context, i *ListSplunksInput) ([]*Splunk, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "splunk")
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ss []*Splunk
	if err := DecodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	return ss, nil
}

// CreateSplunkInput is used as input to the CreateSplunk function.
type CreateSplunkInput struct {
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to Splunk.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// RequestMaxBytes is the maximum number of bytes sent in one request. Defaults 0 for unbounded.
	RequestMaxBytes *int `url:"request_max_bytes,omitempty"`
	// RequestMaxEntries is the maximum number of logs sent in one request. Defaults 0 for unbounded.
	RequestMaxEntries *int `url:"request_max_entries,omitempty"`
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
	// Token is a Splunk token for use in posting logs over HTTP to your collector.
	Token *string `url:"token,omitempty"`
	// URL is the URL to post logs to.
	URL *string `url:"url,omitempty"`
	// UseTLS is whether to use TLS (0: do not use, 1: use).
	UseTLS *Compatibool `url:"use_tls,omitempty"`
}

// CreateSplunk creates a new resource.
func (c *Client) CreateSplunk(ctx context.Context, i *CreateSplunkInput) (*Splunk, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "splunk")
	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Splunk
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSplunkInput is used as input to the GetSplunk function.
type GetSplunkInput struct {
	// Name is the name of the splunk to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetSplunk retrieves the specified resource.
func (c *Client) GetSplunk(ctx context.Context, i *GetSplunkInput) (*Splunk, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "splunk", i.Name)
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Splunk
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateSplunkInput is used as input to the UpdateSplunk function.
type UpdateSplunkInput struct {
	// Format is a Fastly log format string.
	Format *string `json:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `json:"format_version,omitempty"`
	// Name is the name of the splunk to update (required).
	Name string `json:"-"`
	// NewName is the new name for the resource.
	NewName *string `json:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed. Use
	// NullValue[string]() to reset the endpoint to automatic placement.
	Placement *Nullable[string] `json:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to Splunk.
	ProcessingRegion *string `json:"log_processing_region,omitempty"`
	// RequestMaxBytes is the maximum number of bytes sent in one request. Defaults 0 for unbounded.
	RequestMaxBytes *int `json:"request_max_bytes,omitempty"`
	// RequestMaxEntries is the maximum number of logs sent in one request. Defaults 0 for unbounded.
	RequestMaxEntries *int `json:"request_max_entries,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `json:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `json:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `json:"-"`
	// TLSCACert is a secure certificate to authenticate a server with. Must be in PEM format.
	// Use NullValue[string]() to clear the certificate.
	TLSCACert *Nullable[string] `json:"tls_ca_cert,omitempty"`
	// TLSClientCert is the client certificate used to make authenticated requests. Must be in PEM format.
	// Use NullValue[string]() to clear the certificate.
	TLSClientCert *Nullable[string] `json:"tls_client_cert,omitempty"`
	// TLSClientKey is the client private key used to make authenticated requests. Must be in PEM format.
	// Use NullValue[string]() to clear the key.
	TLSClientKey *Nullable[string] `json:"tls_client_key,omitempty"`
	// TLSHostname is the hostname to verify the server's certificate. This should be one of the Subject Alternative Name (SAN) fields for the certificate. Common Names (CN) are not supported.
	// Use NullValue[string]() to clear the hostname.
	TLSHostname *Nullable[string] `json:"tls_hostname,omitempty"`
	// Token is a Splunk token for use in posting logs over HTTP to your collector.
	Token *string `json:"token,omitempty"`
	// URL is the URL to post logs to.
	URL *string `json:"url,omitempty"`
	// UseTLS is whether to use TLS (0: do not use, 1: use).
	UseTLS *Compatibool `json:"use_tls,omitempty"`
}

// UpdateSplunk updates the specified resource.
func (c *Client) UpdateSplunk(ctx context.Context, i *UpdateSplunkInput) (*Splunk, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "splunk", i.Name)
	resp, err := c.PutJSON(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Splunk
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
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

// DeleteSplunk deletes the specified resource.
func (c *Client) DeleteSplunk(ctx context.Context, i *DeleteSplunkInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "splunk", i.Name)
	resp, err := c.Delete(ctx, path, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
