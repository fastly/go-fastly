package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Elasticsearch represents an Elasticsearch Logging response from the Fastly API.
type Elasticsearch struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListElasticsearch returns the list of Elasticsearch logs for the configuration version.
func (c *Client) ListElasticsearch(i *ListElasticsearchInput) ([]*Elasticsearch, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var elasticsearch []*Elasticsearch
	if err := decodeBodyMap(resp.Body, &elasticsearch); err != nil {
		return nil, err
	}
	sort.Stable(elasticsearchByName(elasticsearch))
	return elasticsearch, nil
}

// CreateElasticsearchInput is used as input to the CreateElasticsearch function.
type CreateElasticsearchInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `url:"name,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	Format            string `url:"format,omitempty"`
	Index             string `url:"index,omitempty"`
	URL               string `url:"url,omitempty"`
	Pipeline          string `url:"pipeline,omitempty"`
	User              string `url:"user,omitempty"`
	Password          string `url:"password,omitempty"`
	RequestMaxEntries uint   `url:"request_max_entries,omitempty"`
	RequestMaxBytes   uint   `url:"request_max_bytes,omitempty"`
	Placement         string `url:"placement,omitempty"`
	TLSCACert         string `url:"tls_ca_cert,omitempty"`
	TLSClientCert     string `url:"tls_client_cert,omitempty"`
	TLSClientKey      string `url:"tls_client_key,omitempty"`
	TLSHostname       string `url:"tls_hostname,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
}

// CreateElasticsearch creates a new Fastly Elasticsearch logging endpoint.
func (c *Client) CreateElasticsearch(i *CreateElasticsearchInput) (*Elasticsearch, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch", i.ServiceID, i.ServiceVersion)
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Elasticsearch endpoint to fetch.
	Name string
}

func (c *Client) GetElasticsearch(i *GetElasticsearchInput) (*Elasticsearch, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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

type UpdateElasticsearchInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Elasticsearch endpoint to fetch.
	Name string

	NewName           *string `url:"name,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	Format            *string `url:"format,omitempty"`
	Index             *string `url:"index,omitempty"`
	URL               *string `url:"url,omitempty"`
	Pipeline          *string `url:"pipeline,omitempty"`
	User              *string `url:"user,omitempty"`
	Password          *string `url:"password,omitempty"`
	RequestMaxEntries *uint   `url:"request_max_entries,omitempty"`
	RequestMaxBytes   *uint   `url:"request_max_bytes,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	TLSCACert         *string `url:"tls_ca_cert,omitempty"`
	TLSClientCert     *string `url:"tls_client_cert,omitempty"`
	TLSClientKey      *string `url:"tls_client_key,omitempty"`
	TLSHostname       *string `url:"tls_hostname,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
}

func (c *Client) UpdateElasticsearch(i *UpdateElasticsearchInput) (*Elasticsearch, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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

type DeleteElasticsearchInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Elasticsearch endpoint to fetch.
	Name string
}

func (c *Client) DeleteElasticsearch(i *DeleteElasticsearchInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/elasticsearch/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
