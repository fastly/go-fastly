package fastly

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// HTTPS represents an HTTPS Logging response from the Fastly API.
type HTTPS struct {
	CompressionCodec  *string    `mapstructure:"compression_codec"`
	ContentType       *string    `mapstructure:"content_type"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	GzipLevel         *int       `mapstructure:"gzip_level"`
	HeaderName        *string    `mapstructure:"header_name"`
	HeaderValue       *string    `mapstructure:"header_value"`
	JSONFormat        *string    `mapstructure:"json_format"`
	MessageType       *string    `mapstructure:"message_type"`
	Method            *string    `mapstructure:"method"`
	Name              *string    `mapstructure:"name"`
	Period            *int       `mapstructure:"period"`
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
	URL               *string    `mapstructure:"url"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// ListHTTPSInput is used as input to the ListHTTPS function.
type ListHTTPSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListHTTPS retrieves all resources.
func (c *Client) ListHTTPS(ctx context.Context, i *ListHTTPSInput) ([]*HTTPS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "https")
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var https []*HTTPS
	if err := DecodeBodyMap(resp.Body, &https); err != nil {
		return nil, err
	}
	return https, nil
}

// CreateHTTPSInput is used as input to the CreateHTTPS function.
type CreateHTTPSInput struct {
	// CompressionCodec is the codec used for compressing your logs (zstd, snappy, gzip).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// ContentType is the content type of the header sent with the request.
	ContentType *string `url:"content_type,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// HeaderName is the name of the custom header sent with the request.
	HeaderName *string `url:"header_name,omitempty"`
	// HeaderValue is the value of the custom header sent with the request.
	HeaderValue *string `url:"header_value,omitempty"`
	// JSONFormat enforces valid JSON formatting for log entries (0: disabled, 1: array of JSON, 2: newline delimited JSON).
	JSONFormat *string `url:"json_format,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Method is the HTTP method used for request (POST, PUT).
	Method *string `url:"method,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Period is the time interval (in seconds) for log flushing, defaults to 5 seconds.
	Period *int `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to the HTTPS server.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// RequestMaxBytes is the maximum number of bytes sent in one request. Defaults 0 (100MB).
	RequestMaxBytes *int `url:"request_max_bytes,omitempty"`
	// RequestMaxEntries is the maximum number of logs sent in one request. Defaults 0 (10k).
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
	// URL is the URL to send logs to. Must use HTTPS
	URL *string `url:"url,omitempty"`
}

// CreateHTTPS creates a new resource.
func (c *Client) CreateHTTPS(ctx context.Context, i *CreateHTTPSInput) (*HTTPS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "https")
	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var https *HTTPS
	if err := DecodeBodyMap(resp.Body, &https); err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", https)
	return https, nil
}

// GetHTTPSInput is used as input to the GetHTTPS function.
type GetHTTPSInput struct {
	// Name is the name of the HTTPS endpoint to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetHTTPS retrieves the specified resource.
func (c *Client) GetHTTPS(ctx context.Context, i *GetHTTPSInput) (*HTTPS, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "https", i.Name)
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *HTTPS
	if err := DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}

	return h, nil
}

// UpdateHTTPSInput is the input parameter to the UpdateHTTPS function.
type UpdateHTTPSInput struct {
	// CompressionCodec is the codec used for compressing your logs. Valid values are zstd, snappy, and gzip.
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// ContentType is the content type of the header sent with the request.
	ContentType *string `url:"content_type,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// HeaderName is the name of the custom header sent with the request.
	HeaderName *string `url:"header_name,omitempty"`
	// HeaderValue is the value of the custom header sent with the request.
	HeaderValue *string `url:"header_value,omitempty"`
	// JSONFormat enforces valid JSON formatting for log entries (0: disabled, 1: array of JSON, 2: newline delimited JSON).
	JSONFormat *string `url:"json_format,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Method is the HTTP method used for request (POST, PUT).
	Method *string `url:"method,omitempty"`
	// Name is the name of the HTTPS endpoint to fetch (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Period is the time interval (in seconds) for log flushing, defaults to 5 seconds.
	Period *int `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to the HTTPS server.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// RequestMaxBytes is the maximum number of bytes sent in one request. Defaults 0 (100MB).
	RequestMaxBytes *int `url:"request_max_bytes,omitempty"`
	// RequestMaxEntries is the maximum number of logs sent in one request. Defaults 0 (10k).
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
	// URL is the URL to send logs to. Must use HTTPS
	URL *string `url:"url,omitempty"`
}

// UpdateHTTPS updates the specified resource.
func (c *Client) UpdateHTTPS(ctx context.Context, i *UpdateHTTPSInput) (*HTTPS, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "https", i.Name)
	resp, err := c.PutForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *HTTPS
	if err := DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// DeleteHTTPSInput is the input parameter to the DeleteHTTPS function.
type DeleteHTTPSInput struct {
	// Name is the name of the HTTPS endpoint to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteHTTPS deletes the specified resource.
func (c *Client) DeleteHTTPS(ctx context.Context, i *DeleteHTTPSInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "https", i.Name)
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
