package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Kafka represents a kafka response from the Fastly API.
type Kafka struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Brokers           string     `mapstructure:"brokers"`
	Topic             string     `mapstructure:"topic"`
	RequiredACKs      string     `mapstructure:"required_acks"`
	UseTLS            bool       `mapstructure:"use_tls"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	TLSCACert         string     `mapstructure:"tls_ca_cert"`
	TLSHostname       string     `mapstructure:"tls_hostname"`
	TLSClientCert     string     `mapstructure:"tls_client_cert"`
	TLSClientKey      string     `mapstructure:"tls_client_key"`
	ParseLogKeyvals   bool       `mapstructure:"parse_log_keyvals"`
	RequestMaxBytes   uint       `mapstructure:"request_max_bytes"`
	AuthMethod        string     `mapstructure:"auth_method"`
	User              string     `mapstructure:"user"`
	Password          string     `mapstructure:"password"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// kafkaByName is a sortable list of kafkas.
type kafkasByName []*Kafka

// Len, Swap, and Less implement the sortable interface.
func (s kafkasByName) Len() int      { return len(s) }
func (s kafkasByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s kafkasByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListKafkasInput is used as input to the ListKafkas function.
type ListKafkasInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListKafkas returns the list of kafkas for the configuration version.
func (c *Client) ListKafkas(i *ListKafkasInput) ([]*Kafka, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var k []*Kafka
	if err := decodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	sort.Stable(kafkasByName(k))
	return k, nil
}

// CreateKafkaInput is used as input to the CreateKafka function.
type CreateKafkaInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string      `url:"name,omitempty"`
	Brokers           string      `url:"brokers,omitempty"`
	Topic             string      `url:"topic,omitempty"`
	RequiredACKs      string      `url:"required_acks,omitempty"`
	UseTLS            Compatibool `url:"use_tls,omitempty"`
	CompressionCodec  string      `url:"compression_codec,omitempty"`
	Format            string      `url:"format,omitempty"`
	FormatVersion     uint        `url:"format_version,omitempty"`
	ResponseCondition string      `url:"response_condition,omitempty"`
	Placement         string      `url:"placement,omitempty"`
	TLSCACert         string      `url:"tls_ca_cert,omitempty"`
	TLSHostname       string      `url:"tls_hostname,omitempty"`
	TLSClientCert     string      `url:"tls_client_cert,omitempty"`
	TLSClientKey      string      `url:"tls_client_key,omitempty"`
	ParseLogKeyvals   Compatibool `url:"parse_log_keyvals,omitempty"`
	RequestMaxBytes   uint        `url:"request_max_bytes,omitempty"`
	AuthMethod        string      `url:"auth_method,omitempty"`
	User              string      `url:"user,omitempty"`
	Password          string      `url:"password,omitempty"`
}

// CreateKafka creates a new Fastly kafka.
func (c *Client) CreateKafka(i *CreateKafkaInput) (*Kafka, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var k *Kafka
	if err := decodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// GetKafkaInput is used as input to the GetKafka function.
type GetKafkaInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the kafka to fetch.
	Name string
}

// GetKafka gets the kafka configuration with the given parameters.
func (c *Client) GetKafka(i *GetKafkaInput) (*Kafka, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var k *Kafka
	if err := decodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// UpdateKafkaInput is used as input to the UpdateKafka function.
type UpdateKafkaInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the kafka to update.
	Name string

	NewName           *string      `url:"name,omitempty"`
	Brokers           *string      `url:"brokers,omitempty"`
	Topic             *string      `url:"topic,omitempty"`
	RequiredACKs      *string      `url:"required_acks,omitempty"`
	UseTLS            *Compatibool `url:"use_tls,omitempty"`
	CompressionCodec  *string      `url:"compression_codec,omitempty"`
	Format            *string      `url:"format,omitempty"`
	FormatVersion     *uint        `url:"format_version,omitempty"`
	ResponseCondition *string      `url:"response_condition,omitempty"`
	Placement         *string      `url:"placement,omitempty"`
	TLSCACert         *string      `url:"tls_ca_cert,omitempty"`
	TLSHostname       *string      `url:"tls_hostname,omitempty"`
	TLSClientCert     *string      `url:"tls_client_cert,omitempty"`
	TLSClientKey      *string      `url:"tls_client_key,omitempty"`
	ParseLogKeyvals   *Compatibool `url:"parse_log_keyvals,omitempty"`
	RequestMaxBytes   *uint        `url:"request_max_bytes,omitempty"`
	AuthMethod        *string      `url:"auth_method,omitempty"`
	User              *string      `url:"user,omitempty"`
	Password          *string      `url:"password,omitempty"`
}

// UpdateKafka updates a specific kafka.
func (c *Client) UpdateKafka(i *UpdateKafkaInput) (*Kafka, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var k *Kafka
	if err := decodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// DeleteKafkaInput is the input parameter to DeleteKafka.
type DeleteKafkaInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the kafka to delete (required).
	Name string
}

// DeleteKafka deletes the given kafka version.
func (c *Client) DeleteKafka(i *DeleteKafkaInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
