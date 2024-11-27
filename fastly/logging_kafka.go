package fastly

import (
	"strconv"
	"time"
)

// Kafka represents a kafka response from the Fastly API.
type Kafka struct {
	AuthMethod        *string    `mapstructure:"auth_method"`
	Brokers           *string    `mapstructure:"brokers"`
	CompressionCodec  *string    `mapstructure:"compression_codec"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	ParseLogKeyvals   *bool      `mapstructure:"parse_log_keyvals"`
	Password          *string    `mapstructure:"password"`
	Placement         *string    `mapstructure:"placement"`
	RequestMaxBytes   *int       `mapstructure:"request_max_bytes"`
	RequiredACKs      *string    `mapstructure:"required_acks"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	TLSCACert         *string    `mapstructure:"tls_ca_cert"`
	TLSClientCert     *string    `mapstructure:"tls_client_cert"`
	TLSClientKey      *string    `mapstructure:"tls_client_key"`
	TLSHostname       *string    `mapstructure:"tls_hostname"`
	Topic             *string    `mapstructure:"topic"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	UseTLS            *bool      `mapstructure:"use_tls"`
	User              *string    `mapstructure:"user"`
}

// ListKafkasInput is used as input to the ListKafkas function.
type ListKafkasInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListKafkas retrieves all resources.
func (c *Client) ListKafkas(i *ListKafkasInput) ([]*Kafka, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "kafka")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var k []*Kafka
	if err := DecodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// CreateKafkaInput is used as input to the CreateKafka function.
type CreateKafkaInput struct {
	// AuthMethod is the SASL authentication method (plain, scram-sha-256, scram-sha-512).
	AuthMethod *string `url:"auth_method,omitempty"`
	// Brokers is a comma-separated list of IP addresses or hostnames of Kafka brokers.
	Brokers *string `url:"brokers,omitempty"`
	// CompressionCodec is the codec used for compression of your logs (gzip, snappy, lz4, null).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// ParseLogKeyvals enables parsing of key=value tuples from the beginning of a logline, turning them into record headers.
	ParseLogKeyvals *Compatibool `url:"parse_log_keyvals,omitempty"`
	// Password is the SASL password.
	Password *string `url:"password,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// RequestMaxBytes is the maximum number of bytes sent in one request. Defaults 0 (no limit).
	RequestMaxBytes *int `url:"request_max_bytes,omitempty"`
	// RequiredACKs is the number of acknowledgements a leader must receive before a write is considered successful.
	RequiredACKs *string `url:"required_acks,omitempty"`
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
	// Topic is the Kafka topic to send logs to.
	Topic *string `url:"topic,omitempty"`
	// UseTLS is whether to use TLS (0: do not use, 1: use).
	UseTLS *Compatibool `url:"use_tls,omitempty"`
	// User is the SASL user.
	User *string `url:"user,omitempty"`
}

// CreateKafka creates a new resource.
func (c *Client) CreateKafka(i *CreateKafkaInput) (*Kafka, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "kafka")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var k *Kafka
	if err := DecodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// GetKafkaInput is used as input to the GetKafka function.
type GetKafkaInput struct {
	// Name is the name of the kafka to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetKafka retrieves the specified resource.
func (c *Client) GetKafka(i *GetKafkaInput) (*Kafka, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "kafka", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var k *Kafka
	if err := DecodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// UpdateKafkaInput is used as input to the UpdateKafka function.
type UpdateKafkaInput struct {
	// AuthMethod is the SASL authentication method (plain, scram-sha-256, scram-sha-512).
	AuthMethod *string `url:"auth_method,omitempty"`
	// Brokers is a comma-separated list of IP addresses or hostnames of Kafka brokers.
	Brokers *string `url:"brokers,omitempty"`
	// CompressionCodec is the codec used for compression of your logs (gzip, snappy, lz4, null).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the kafka to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// ParseLogKeyvals enables parsing of key=value tuples from the beginning of a logline, turning them into record headers.
	ParseLogKeyvals *Compatibool `url:"parse_log_keyvals,omitempty"`
	// Password is the SASL password.
	Password *string `url:"password,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// RequestMaxBytes is the maximum number of bytes sent in one request. Defaults 0 (no limit).
	RequestMaxBytes *int `url:"request_max_bytes,omitempty"`
	// RequiredACKs is the number of acknowledgements a leader must receive before a write is considered successful.
	RequiredACKs *string `url:"required_acks,omitempty"`
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
	// Topic is the Kafka topic to send logs to.
	Topic *string `url:"topic,omitempty"`
	// UseTLS is whether to use TLS (0: do not use, 1: use).
	UseTLS *Compatibool `url:"use_tls,omitempty"`
	// User is the SASL user.
	User *string `url:"user,omitempty"`
}

// UpdateKafka updates the specified resource.
func (c *Client) UpdateKafka(i *UpdateKafkaInput) (*Kafka, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "kafka", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var k *Kafka
	if err := DecodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// DeleteKafkaInput is the input parameter to DeleteKafka.
type DeleteKafkaInput struct {
	// Name is the name of the kafka to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteKafka deletes the specified resource.
func (c *Client) DeleteKafka(i *DeleteKafkaInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "kafka", i.Name)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
