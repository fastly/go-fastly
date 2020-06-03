package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Kafka represents a kafka response from the Fastly API.
type Kafka struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

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
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListKafkas returns the list of kafkas for the configuration version.
func (c *Client) ListKafkas(i *ListKafkasInput) ([]*Kafka, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var k []*Kafka
	if err := decodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	sort.Stable(kafkasByName(k))
	return k, nil
}

// CreateKafkaInput is used as input to the CreateKafka function.
type CreateKafkaInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	Name              *string      `form:"name,omitempty"`
	Brokers           *string      `form:"brokers,omitempty"`
	Topic             *string      `form:"topic,omitempty"`
	RequiredACKs      *string      `form:"required_acks,omitempty"`
	UseTLS            *Compatibool `form:"use_tls,omitempty"`
	CompressionCodec  *string      `form:"compression_codec,omitempty"`
	Format            *string      `form:"format,omitempty"`
	FormatVersion     *uint        `form:"format_version,omitempty"`
	ResponseCondition *string      `form:"response_condition,omitempty"`
	Placement         *string      `form:"placement,omitempty"`
	TLSCACert         *string      `form:"tls_ca_cert,omitempty"`
	TLSHostname       *string      `form:"tls_hostname,omitempty"`
	TLSClientCert     *string      `form:"tls_client_cert,omitempty"`
	TLSClientKey      *string      `form:"tls_client_key,omitempty"`
}

// CreateKafka creates a new Fastly kafka.
func (c *Client) CreateKafka(i *CreateKafkaInput) (*Kafka, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var k *Kafka
	if err := decodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// GetKafkaInput is used as input to the GetKafka function.
type GetKafkaInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the kafka to fetch.
	Name string
}

// GetKafka gets the kafka configuration with the given parameters.
func (c *Client) GetKafka(i *GetKafkaInput) (*Kafka, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var k *Kafka
	if err := decodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// UpdateKafkaInput is used as input to the UpdateKafka function.
type UpdateKafkaInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the kafka to update.
	Name string

	NewName           *string      `form:"name,omitempty"`
	Brokers           *string      `form:"brokers,omitempty"`
	Topic             *string      `form:"topic,omitempty"`
	RequiredACKs      *string      `form:"required_acks,omitempty"`
	UseTLS            *Compatibool `form:"use_tls,omitempty"`
	CompressionCodec  *string      `form:"compression_codec,omitempty"`
	Format            *string      `form:"format,omitempty"`
	FormatVersion     *uint        `form:"format_version,omitempty"`
	ResponseCondition *string      `form:"response_condition,omitempty"`
	Placement         *string      `form:"placement,omitempty"`
	TLSCACert         *string      `form:"tls_ca_cert,omitempty"`
	TLSHostname       *string      `form:"tls_hostname,omitempty"`
	TLSClientCert     *string      `form:"tls_client_cert,omitempty"`
	TLSClientKey      *string      `form:"tls_client_key,omitempty"`
}

// UpdateKafka updates a specific kafka.
func (c *Client) UpdateKafka(i *UpdateKafkaInput) (*Kafka, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var k *Kafka
	if err := decodeBodyMap(resp.Body, &k); err != nil {
		return nil, err
	}
	return k, nil
}

// DeleteKafkaInput is the input parameter to DeleteKafka.
type DeleteKafkaInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the kafka to delete (required).
	Name string
}

// DeleteKafka deletes the given kafka version.
func (c *Client) DeleteKafka(i *DeleteKafkaInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kafka/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
