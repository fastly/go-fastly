package fastly

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// CustomTLSConfiguration represents a TLS configuration response from the Fastly API.
type CustomTLSConfiguration struct {
	ID            string       `jsonapi:"primary,tls_configuration"`
	DNSRecords    []*DNSRecord `jsonapi:"relation,dns_records"`
	Bulk          bool         `jsonapi:"attr,bulk"`
	Default       bool         `jsonapi:"attr,default"`
	HTTPProtocols []string     `jsonapi:"attr,http_protocols"`
	Name          string       `jsonapi:"attr,name"`
	TLSProtocols  []string     `jsonapi:"attr,tls_protocols"`
	CreatedAt     *time.Time   `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt     *time.Time   `jsonapi:"attr,updated_at,iso8601"`
}

// DNSRecord is a child of CustomTLSConfiguration
type DNSRecord struct {
	ID         string `jsonapi:"primary,dns_record"`
	RecordType string `jsonapi:"attr,record_type"`
	Region     string `jsonapi:"attr,region"`
}

// ListCustomTLSConfigurationsInput is used as input to the ListCustomTLSConfigurationsInput function.
type ListCustomTLSConfigurationsInput struct {
	FilterBulk bool   // Whether or not to only include bulk=true configurations
	Include    string // Include related objects. Optional, comma-separated values. Permitted values: dns_records.
	PageNumber int    // The page index for pagination.
	PageSize   int    // The number of keys per page.
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListCustomTLSConfigurationsInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[bulk]": i.FilterBulk,
		"include":      i.Include,
		"page[size]":   i.PageSize,
		"page[number]": i.PageNumber,
	}

	for key, value := range pairings {
		switch t := reflect.TypeOf(value).String(); t {
		case "string":
			if value != "" {
				result[key] = value.(string)
			}
		case "int":
			if value != 0 {
				result[key] = strconv.Itoa(value.(int))
			}
		}
	}

	return result
}

// ListCustomTLSConfigurations list all TLS configurations.
func (c *Client) ListCustomTLSConfigurations(i *ListCustomTLSConfigurationsInput) ([]*CustomTLSConfiguration, error) {
	p := "/tls/configurations"
	ro := &RequestOptions{
		Params: i.formatFilters(),
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the filters don't work
		},
	}

	r, err := c.Get(p, ro)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(CustomTLSConfiguration)))
	if err != nil {
		return nil, err
	}

	con := make([]*CustomTLSConfiguration, len(data))
	for i := range data {
		typed, ok := data[i].(*CustomTLSConfiguration)
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T", data[i])
		}
		con[i] = typed
	}

	return con, nil
}

// GetCustomTLSConfigurationInput is used as input to the GetCustomTLSConfiguration function.
type GetCustomTLSConfigurationInput struct {
	ID      string
	Include string // Include related objects. Optional, comma-separated values. Permitted values: dns_records.
}

// GetCustomTLSConfiguration returns a single TLS configuration.
func (c *Client) GetCustomTLSConfiguration(i *GetCustomTLSConfigurationInput) (*CustomTLSConfiguration, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := fmt.Sprintf("/tls/configurations/%s", i.ID)

	ro := &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the params don't work
		},
	}

	if i.Include != "" {
		ro.Params = map[string]string{"include": i.Include}
	}

	r, err := c.Get(p, ro)
	if err != nil {
		return nil, err
	}

	var con CustomTLSConfiguration
	if err := jsonapi.UnmarshalPayload(r.Body, &con); err != nil {
		return nil, err
	}

	return &con, nil
}

// UpdateCustomTLSConfigurationInput is used as input to the UpdateCustomTLSConfiguration function.
type UpdateCustomTLSConfigurationInput struct {
	ID   string
	Name string `jsonapi:"attr,name"`
}

// UpdateCustomTLSConfiguration can only be used to change the name of the configuration
func (c *Client) UpdateCustomTLSConfiguration(i *UpdateCustomTLSConfigurationInput) (*CustomTLSConfiguration, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/tls/configurations/%s", i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var con CustomTLSConfiguration
	if err := jsonapi.UnmarshalPayload(resp.Body, &con); err != nil {
		return nil, err
	}
	return &con, nil
}
