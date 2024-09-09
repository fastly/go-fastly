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
	Bulk          bool         `jsonapi:"attr,bulk"`
	CreatedAt     *time.Time   `jsonapi:"attr,created_at,iso8601"`
	DNSRecords    []*DNSRecord `jsonapi:"relation,dns_records"`
	Default       bool         `jsonapi:"attr,default"`
	HTTPProtocols []string     `jsonapi:"attr,http_protocols"`
	ID            string       `jsonapi:"primary,tls_configuration"`
	Name          string       `jsonapi:"attr,name"`
	TLSProtocols  []string     `jsonapi:"attr,tls_protocols"`
	UpdatedAt     *time.Time   `jsonapi:"attr,updated_at,iso8601"`
}

// DNSRecord is a child of CustomTLSConfiguration.
type DNSRecord struct {
	ID         string `jsonapi:"primary,dns_record"`
	RecordType string `jsonapi:"attr,record_type"`
	Region     string `jsonapi:"attr,region"`
}

// ListCustomTLSConfigurationsInput is used as input to the ListCustomTLSConfigurationsInput function.
type ListCustomTLSConfigurationsInput struct {
	// FilterBulk is whether or not to only include bulk=true configurations
	FilterBulk bool
	// Include captures related objects. Optional, comma-separated values. Permitted values: dns_records.
	Include string
	// PageNumber is the page index for pagination.
	PageNumber int
	// PageSize is the number of keys per page.
	PageSize int
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListCustomTLSConfigurationsInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"filter[bulk]": i.FilterBulk,
		"include":      i.Include,
		"page[size]":   i.PageSize,
		"page[number]": i.PageNumber,
	}

	for key, value := range pairings {
		switch t := reflect.TypeOf(value).String(); t {
		case "string":
			if value != "" {
				v, _ := value.(string) // type assert to avoid runtime panic (v will have zero value for its type)
				result[key] = v
			}
		case "int":
			if value != 0 {
				v, _ := value.(int) // type assert to avoid runtime panic (v will have zero value for its type)
				result[key] = strconv.Itoa(v)
			}
		}
	}

	return result
}

// ListCustomTLSConfigurations retrieves all resources.
func (c *Client) ListCustomTLSConfigurations(i *ListCustomTLSConfigurationsInput) ([]*CustomTLSConfiguration, error) {
	path := "/tls/configurations"
	ro := &RequestOptions{
		Params: i.formatFilters(),
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the filters don't work
		},
	}

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CustomTLSConfiguration)))
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
	// ID is an alphanumeric string identifying a TLS configuration.
	ID string
	// Include captures related objects. Optional, comma-separated values. Permitted values: dns_records.
	Include string
}

// GetCustomTLSConfiguration retrieves the specified resource.
func (c *Client) GetCustomTLSConfiguration(i *GetCustomTLSConfigurationInput) (*CustomTLSConfiguration, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("tls", "configurations", i.ID)

	ro := &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the params don't work
		},
	}

	if i.Include != "" {
		ro.Params = map[string]string{"include": i.Include}
	}

	resp, err := c.Get(path, ro)
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

// UpdateCustomTLSConfigurationInput is used as input to the UpdateCustomTLSConfiguration function.
type UpdateCustomTLSConfigurationInput struct {
	// ID is an alphanumeric string identifying a TLS configuration.
	ID string
	// Name is a custom name for your TLS configuration.
	Name string `jsonapi:"attr,name"`
}

// UpdateCustomTLSConfiguration updates the specified resource.
func (c *Client) UpdateCustomTLSConfiguration(i *UpdateCustomTLSConfigurationInput) (*CustomTLSConfiguration, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := ToSafeURL("tls", "configurations", i.ID)

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
