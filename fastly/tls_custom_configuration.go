package fastly

import (
	"context"
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
		"filter[bulk]":               i.FilterBulk,
		"include":                    i.Include,
		jsonapi.QueryParamPageSize:   i.PageSize,
		jsonapi.QueryParamPageNumber: i.PageNumber,
	}

	for key, value := range pairings {
		switch v := value.(type) {
		case string:
			if v != "" {
				result[key] = v
			}
		case int:
			if v != 0 {
				result[key] = strconv.Itoa(v)
			}
		}
	}

	return result
}

// ListCustomTLSConfigurations retrieves all resources.
func (c *Client) ListCustomTLSConfigurations(ctx context.Context, i *ListCustomTLSConfigurationsInput) ([]*CustomTLSConfiguration, error) {
	path := "/tls/configurations"
	requestOptions := CreateRequestOptions()
	requestOptions.Params = i.formatFilters()
	requestOptions.Headers["Accept"] = jsonapi.MediaType // this is required otherwise the filters don't work

	resp, err := c.Get(ctx, path, requestOptions)
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
func (c *Client) GetCustomTLSConfiguration(ctx context.Context, i *GetCustomTLSConfigurationInput) (*CustomTLSConfiguration, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("tls", "configurations", i.ID)

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = jsonapi.MediaType // this is required otherwise the filters don't work

	if i.Include != "" {
		requestOptions.Params["include"] = i.Include
	}

	resp, err := c.Get(ctx, path, requestOptions)
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
func (c *Client) UpdateCustomTLSConfiguration(ctx context.Context, i *UpdateCustomTLSConfigurationInput) (*CustomTLSConfiguration, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := ToSafeURL("tls", "configurations", i.ID)

	resp, err := c.PatchJSONAPI(ctx, path, i, CreateRequestOptions())
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
