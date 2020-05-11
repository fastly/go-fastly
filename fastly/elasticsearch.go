package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Elasticsearch represents an Elasticsearch Logging response from the Fastly API.
type Elasticsearch struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Format            string     `mapstructure:"format"`
	Index             string     `mapstructure:"index"`
	URL               string     `mapstructure:"url"`
	Pipeline          string     `mapstructure:"pipeline"`
	User              string     `mapstructure:"user"`
	Password          string     `mapstructure:"password"`
	RequestMaxEntries uint       `mapstructure:"request_max_entries"`
	RequestMaxBytes   uint       `mapstructure:"request_max_bytes"`
	Placement         string     `mapstructure:"placement"`
	TLSCACert         string     `mapstructure:"tls_ca_cert"`
	TLSClientCert     string     `mapstructure:"tls_client_cert"`
	TLSClientKey      string     `mapstructure:"tls_client_key"`
	TLSHostname       string     `mapstructure:"tls_hostname"`
	FormatVersion     uint       `mapstructure:"format_version"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// elasticsearchByName is a sortable list of Elasticsearch logs.
type elasticsearchByName []*Elasticsearch

// Len, Swap, and Less implement the sortable interface.
func (s elasticsearchByName) Len() int      { return len(s) }
func (s elasticsearchByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s elasticsearchByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListElasticsearchInput is used as input to the ListElasticsearch function.
type ListElasticsearchInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListElasticsearch returns the list of Elasticsearch logs for the configuration version.
func (c *Client) ListElasticsearch(i *ListElasticsearchInput) ([]*Elasticsearch, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var elasticsearch []*Elasticsearch
	if err := decodeBodyMap(resp.Body, &elasticsearch); err != nil {
		return nil, err
	}
	sort.Stable(elasticsearchByName(elasticsearch))
	return elasticsearch, nil
}

// CreateElasticsearchInput is used as input to the CreateElasticsearch function.
type CreateElasticsearchInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	Name              *string `form:"name,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Format            *string `form:"format,omitempty"`
	Index             *string `form:"index,omitempty"`
	URL               *string `form:"url,omitempty"`
	Pipeline          *string `form:"pipeline,omitempty"`
	User              *string `form:"user,omitempty"`
	Password          *string `form:"password,omitempty"`
	RequestMaxEntries *uint   `form:"request_max_entries,omitempty"`
	RequestMaxBytes   *uint   `form:"request_max_bytes,omitempty"`
	Placement         *string `form:"placement,omitempty"`
	TLSCACert         *string `form:"tls_ca_cert,omitempty"`
	TLSClientCert     *string `form:"tls_client_cert,omitempty"`
	TLSClientKey      *string `form:"tls_client_key,omitempty"`
	TLSHostname       *string `form:"tls_hostname,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
}

// CreateElasticsearch creates a new Fastly Elasticsearch logging endpoint.
func (c *Client) CreateElasticsearch(i *CreateElasticsearchInput) (*Elasticsearch, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var elasticsearch *Elasticsearch
	if err := decodeBodyMap(resp.Body, &elasticsearch); err != nil {
		return nil, err
	}
	return elasticsearch, nil
}

// GetElasticsearchInput is used as input to the GetElasticsearch function.
type GetElasticsearchInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the Elasticsearch endpoint to fetch.
	Name string
}

func (c *Client) GetElasticsearch(i *GetElasticsearchInput) (*Elasticsearch, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var es *Elasticsearch
	if err := decodeBodyMap(resp.Body, &es); err != nil {
		return nil, err
	}

	return es, nil
}

type UpdateElasticsearchInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the Elasticsearch endpoint to fetch.
	Name string

	NewName           *string `form:"name,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Format            *string `form:"format,omitempty"`
	Index             *string `form:"index,omitempty"`
	URL               *string `form:"url,omitempty"`
	Pipeline          *string `form:"pipeline,omitempty"`
	User              *string `form:"user,omitempty"`
	Password          *string `form:"password,omitempty"`
	RequestMaxEntries *uint   `form:"request_max_entries,omitempty"`
	RequestMaxBytes   *uint   `form:"request_max_bytes,omitempty"`
	Placement         *string `form:"placement,omitempty"`
	TLSCACert         *string `form:"tls_ca_cert,omitempty"`
	TLSClientCert     *string `form:"tls_client_cert,omitempty"`
	TLSClientKey      *string `form:"tls_client_key,omitempty"`
	TLSHostname       *string `form:"tls_hostname,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
}

func (c *Client) UpdateElasticsearch(i *UpdateElasticsearchInput) (*Elasticsearch, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var es *Elasticsearch
	if err := decodeBodyMap(resp.Body, &es); err != nil {
		return nil, err
	}
	return es, nil
}

type DeleteElasticsearchInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the Elasticsearch endpoint to fetch.
	Name string
}

func (c *Client) DeleteElasticsearch(i *DeleteElasticsearchInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch/%s", i.Service, i.Version, url.PathEscape(i.Name))
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
