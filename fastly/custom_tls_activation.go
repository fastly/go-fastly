package fastly

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// TLSActivation represents a /tls/activations response.
type TLSActivation struct {
	ID            string                `jsonapi:"primary,tls_activation"`
	Configuration *TLSConfiguration     `jsonapi:"relation,tls_configuration"` // TLSConfiguration type shared with BulkCertificate
	Domain        *TLSDomain            `jsonapi:"relation,tls_domain"`        // TLSDomain type shared with BulkCertificate
	Certificate   *CustomTLSCertificate `jsonapi:"relation,tls_certificate"`
	CreatedAt     *time.Time            `jsonapi:"attr,created_at,iso8601"`
}

// ListTLSActivationsInput is used as input to the ListTLSActivations function.
type ListTLSActivationsInput struct {
	FilterTLSCertificateID   string // Limit the returned activations to a specific certificate.
	FilterTLSConfigurationID string // Limit the returned activations to a specific TLS configuration.
	FilterTLSDomainID        string // Limit the returned rules to a specific domain name.
	Include                  string // Include related objects. Optional, comma-separated values. Permitted values: tls_certificate, tls_configuration, and tls_domain.
	PageNumber               int    // The page index for pagination.
	PageSize                 int    // The number of activations per page.
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListTLSActivationsInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[tls_certificate.id]":   i.FilterTLSCertificateID,
		"filter[tls_configuration.id]": i.FilterTLSConfigurationID,
		"filter[tls_domain.id]":        i.FilterTLSDomainID,
		"include":                      i.Include,
		"page[number]":                 i.PageNumber,
		"page[size]":                   i.PageSize,
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

// ListTLSActivations list all activations.
func (c *Client) ListTLSActivations(i *ListTLSActivationsInput) ([]*TLSActivation, error) {
	p := "/tls/activations"
	filters := &RequestOptions{
		Params: i.formatFilters(),
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the filters don't work
		},
	}

	r, err := c.Get(p, filters)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(TLSActivation)))
	if err != nil {
		return nil, err
	}

	a := make([]*TLSActivation, len(data))
	for i := range data {
		typed, ok := data[i].(*TLSActivation)
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T", data[i])
		}
		a[i] = typed
	}

	return a, nil
}

// GetTLSActivationInput is used as input to the GetTLSActivation function.
type GetTLSActivationInput struct {
	ID      string
	Include *string // Include related objects. Optional, comma-separated values. Permitted values: tls_certificate, tls_configuration, and tls_domain.
}

// GetTLSActivation retrieve a single activation.
func (c *Client) GetTLSActivation(i *GetTLSActivationInput) (*TLSActivation, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := fmt.Sprintf("/tls/activations/%s", i.ID)

	ro := &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the params don't work
		},
	}

	if i.Include != nil {
		ro.Params = map[string]string{"include": *i.Include}
	}

	r, err := c.Get(p, ro)
	if err != nil {
		return nil, err
	}

	var a TLSActivation
	if err := jsonapi.UnmarshalPayload(r.Body, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

// CreateTLSActivationInput is used as input to the CreateTLSActivation function.
type CreateTLSActivationInput struct {
	ID            string                `jsonapi:"primary,tls_activation"`   // ID value does not need to be set.
	Certificate   *CustomTLSCertificate `jsonapi:"relation,tls_certificate"` // Only ID of CustomTLSCertificate needs to be set.
	Configuration *TLSConfiguration     `jsonapi:"relation,tls_configuration,omitempty"`
	Domain        *TLSDomain            `jsonapi:"relation,tls_domain"`
}

// CreateTLSActivation enable TLS for a domain using a custom certificate.
func (c *Client) CreateTLSActivation(i *CreateTLSActivationInput) (*TLSActivation, error) {
	if i.Certificate == nil {
		return nil, ErrMissingTLSCertificate
	}
	if i.Domain == nil {
		return nil, ErrMissingTLSDomain
	}

	p := "/tls/activations"

	r, err := c.PostJSONAPI(p, i, nil)
	if err != nil {
		return nil, err
	}

	var a TLSActivation
	if err := jsonapi.UnmarshalPayload(r.Body, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

// UpdateTLSActivationInput is used as input to the UpdateTLSActivation function.
type UpdateTLSActivationInput struct {
	ID          string                `jsonapi:"primary,tls_activation"`
	Certificate *CustomTLSCertificate `jsonapi:"relation,tls_certificate"` // Only ID of CustomTLSCertificate needs to be set.
}

// UpdateTLSActivation updates the certificate used to terminate TLS traffic for the domain associated with this TLS activation.
func (c *Client) UpdateTLSActivation(i *UpdateTLSActivationInput) (*TLSActivation, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}
	if i.Certificate == nil {
		return nil, ErrMissingTLSCertificate
	}

	path := fmt.Sprintf("/tls/activations/%s", i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var ta TLSActivation
	if err := jsonapi.UnmarshalPayload(resp.Body, &ta); err != nil {
		return nil, err
	}
	return &ta, nil
}

// DeleteTLSActivationInput used for deleting a certificate.
type DeleteTLSActivationInput struct {
	ID string
}

// DeleteTLSActivation destroy a certificate.
func (c *Client) DeleteTLSActivation(i *DeleteTLSActivationInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/activations/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}
