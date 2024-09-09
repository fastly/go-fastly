package fastly

import (
	"strconv"
	"time"
)

// Elasticsearch represents an Elasticsearch Logging response from the Fastly API.
type Elasticsearch struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Index             *string    `mapstructure:"index"`
	Name              *string    `mapstructure:"name"`
	Password          *string    `mapstructure:"password"`
	Pipeline          *string    `mapstructure:"pipeline"`
	Placement         *string    `mapstructure:"placement"`
	RequestMaxBytes   *int       `mapstructure:"request_max_bytes"`
	RequestMaxEntries *int       `mapstructure:"request_max_entries"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	TLSCACert         *string    `mapstructure:"tls_ca_cert"`
	TLSClientCert     *string    `mapstructure:"tls_client_cert"`
	TLSClientKey      *string    `mapstructure:"tls_client_key"`
	TLSHostname       *string    `mapstructure:"tls_hostname"`
	URL               *string    `mapstructure:"url"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	User              *string    `mapstructure:"user"`
}

// ListElasticsearchInput is used as input to the ListElasticsearch function.
type ListElasticsearchInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListElasticsearch retrieves all resources.
func (c *Client) ListElasticsearch(i *ListElasticsearchInput) ([]*Elasticsearch, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "elasticsearch")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var elasticsearch []*Elasticsearch
	if err := decodeBodyMap(resp.Body, &elasticsearch); err != nil {
		return nil, err
	}
	return elasticsearch, nil
}

// CreateElasticsearchInput is used as input to the CreateElasticsearch function.
type CreateElasticsearchInput struct {
	// Format is a Fastly log format string. Must produce valid JSON that Elasticsearch can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Index is the name of the Elasticsearch index to send documents (logs) to.
	Index *string `url:"index,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Password is basic Auth password.
	Password *string `url:"password,omitempty"`
	// Pipeline is the ID of the Elasticsearch ingest pipeline to apply pre-process transformations to before indexing.
	Pipeline *string `url:"pipeline,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// RequestMaxBytes is the maximum number of bytes sent in one request.
	RequestMaxBytes *int `url:"request_max_bytes,omitempty"`
	// RequestMaxEntries is the maximum number of logs sent in one request.
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
	// URL is the URL to stream logs to. Must use HTTPS.
	URL *string `url:"url,omitempty"`
	// User is basic Auth username.
	User *string `url:"user,omitempty"`
}

// CreateElasticsearch creates a new resource.
func (c *Client) CreateElasticsearch(i *CreateElasticsearchInput) (*Elasticsearch, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "elasticsearch")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var elasticsearch *Elasticsearch
	if err := decodeBodyMap(resp.Body, &elasticsearch); err != nil {
		return nil, err
	}
	return elasticsearch, nil
}

// GetElasticsearchInput is used as input to the GetElasticsearch function.
type GetElasticsearchInput struct {
	// Name is the name of the Elasticsearch endpoint to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetElasticsearch retrieves the specified resource.
func (c *Client) GetElasticsearch(i *GetElasticsearchInput) (*Elasticsearch, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "elasticsearch", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var es *Elasticsearch
	if err := decodeBodyMap(resp.Body, &es); err != nil {
		return nil, err
	}

	return es, nil
}

// UpdateElasticsearchInput is the input parameter to the UpdateElasticsearch
// function.
type UpdateElasticsearchInput struct {
	// Format is a Fastly log format string. Must produce valid JSON that Elasticsearch can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Index is the name of the Elasticsearch index to send documents (logs) to.
	Index *string `url:"index,omitempty"`
	// Name is the name of the Elasticsearch endpoint to fetch (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Password is basic Auth password.
	Password *string `url:"password,omitempty"`
	// Pipeline is the ID of the Elasticsearch ingest pipeline to apply pre-process transformations to before indexing.
	Pipeline *string `url:"pipeline,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// RequestMaxBytes is the maximum number of bytes sent in one request.
	RequestMaxBytes *int `url:"request_max_bytes,omitempty"`
	// RequestMaxEntries is the maximum number of logs sent in one request.
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
	// URL is the URL to stream logs to. Must use HTTPS.
	URL *string `url:"url,omitempty"`
	// User is basic Auth username.
	User *string `url:"user,omitempty"`
}

// UpdateElasticsearch updates the specified resource.
func (c *Client) UpdateElasticsearch(i *UpdateElasticsearchInput) (*Elasticsearch, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "elasticsearch", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var es *Elasticsearch
	if err := decodeBodyMap(resp.Body, &es); err != nil {
		return nil, err
	}
	return es, nil
}

// DeleteElasticsearchInput is the input parameter to the DeleteElasticsearch
// function.
type DeleteElasticsearchInput struct {
	// Name is the name of the Elasticsearch endpoint to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteElasticsearch deletes the specified resource.
func (c *Client) DeleteElasticsearch(i *DeleteElasticsearchInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "elasticsearch", i.Name)
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
		return ErrStatusNotOk
	}
	return nil
}
