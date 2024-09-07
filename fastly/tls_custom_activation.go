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
	Certificate          *CustomTLSCertificate    `jsonapi:"relation,tls_certificate"`
	Configuration        *TLSConfiguration        `jsonapi:"relation,tls_configuration"` // TLSConfiguration type shared with BulkCertificate
	CreatedAt            *time.Time               `jsonapi:"attr,created_at,iso8601"`
	Domain               *TLSDomain               `jsonapi:"relation,tls_domain"` // TLSDomain type shared with BulkCertificate
	ID                   string                   `jsonapi:"primary,tls_activation"`
	MutualAuthentication *TLSMutualAuthentication `jsonapi:"relation,mutual_authentication"`
}

// ListTLSActivationsInput is used as input to the ListTLSActivations function.
type ListTLSActivationsInput struct {
	// FilterTLSCertificateID limits the returned activations to a specific certificate.
	FilterTLSCertificateID string
	// FilterTLSConfigurationID limits the returned activations to a specific TLS configuration.
	FilterTLSConfigurationID string
	// FilterTLSDomainID limits the returned rules to a specific domain name.
	FilterTLSDomainID string
	// Include captures related objects. Optional, comma-separated values. Permitted values: tls_certificate, tls_configuration, and tls_domain.
	Include string
	// PageNumber is the page index for pagination.
	PageNumber int
	// PageSize is the number of activations per page.
	PageSize int
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListTLSActivationsInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
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

// ListTLSActivations retrieves all resources.
func (c *Client) ListTLSActivations(i *ListTLSActivationsInput) ([]*TLSActivation, error) {
	path := "/tls/activations"
	filters := &RequestOptions{
		Params: i.formatFilters(),
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the filters don't work
		},
	}

	resp, err := c.Get(path, filters)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(TLSActivation)))
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
	// ID is an alphanumeric string identifying a TLS activation.
	ID string
	// Include related objects. Optional, comma-separated values. Permitted values: tls_certificate, tls_configuration, and tls_domain.
	Include *string
}

// GetTLSActivation retrieves the specified resource.
func (c *Client) GetTLSActivation(i *GetTLSActivationInput) (*TLSActivation, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("tls", "activations", i.ID)

	ro := &RequestOptions{
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the params don't work
		},
	}

	if i.Include != nil {
		ro.Params = map[string]string{"include": *i.Include}
	}

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a TLSActivation
	if err := jsonapi.UnmarshalPayload(resp.Body, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

// CreateTLSActivationInput is used as input to the CreateTLSActivation function.
type CreateTLSActivationInput struct {
	// Certificate is an alphanumeric string identifying a TLS certificate.
	Certificate *CustomTLSCertificate `jsonapi:"relation,tls_certificate"` // Only ID of CustomTLSCertificate needs to be set.
	// Configuration is an alphanumeric string identifying a TLS configuration.
	Configuration *TLSConfiguration `jsonapi:"relation,tls_configuration,omitempty"`
	// Domain is the domain name.
	Domain *TLSDomain `jsonapi:"relation,tls_domain"`
	// ID is an aphanumeric string identifying a TLS activation.
	ID string `jsonapi:"primary,tls_activation"` // ID value does not need to be set.
}

// CreateTLSActivation creates a new resource.
func (c *Client) CreateTLSActivation(i *CreateTLSActivationInput) (*TLSActivation, error) {
	if i.Certificate == nil {
		return nil, ErrMissingTLSCertificate
	}
	if i.Domain == nil {
		return nil, ErrMissingTLSDomain
	}

	path := "/tls/activations"

	resp, err := c.PostJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a TLSActivation
	if err := jsonapi.UnmarshalPayload(resp.Body, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

// UpdateTLSActivationInput is used as input to the UpdateTLSActivation function.
type UpdateTLSActivationInput struct {
	// Certificate is an alphanumeric string identifying a TLS certificate.
	Certificate *CustomTLSCertificate `jsonapi:"relation,tls_certificate"` // Only ID of CustomTLSCertificate needs to be set.
	// ID is an aphanumeric string identifying a TLS activation.
	ID string `jsonapi:"primary,tls_activation"`
	// MutualAuthentication is an alphanumeric string identifying a mutual authentication.
	MutualAuthentication *TLSMutualAuthentication `jsonapi:"relation,mutual_authentication"` // Only ID of TLSMutualAuthentication needs to be set.
}

// UpdateTLSActivation updates the specified resource.
func (c *Client) UpdateTLSActivation(i *UpdateTLSActivationInput) (*TLSActivation, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}
	if i.Certificate == nil && i.MutualAuthentication == nil {
		return nil, ErrMissingCertificateMTLS
	}

	path := ToSafeURL("tls", "activations", i.ID)

	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ta TLSActivation
	if err := jsonapi.UnmarshalPayload(resp.Body, &ta); err != nil {
		return nil, err
	}
	return &ta, nil
}

// DeleteTLSActivationInput used for deleting a certificate.
type DeleteTLSActivationInput struct {
	// ID is an alphanumeric string identifying a TLS activation.
	ID string
}

// DeleteTLSActivation deletes the specified resource.
func (c *Client) DeleteTLSActivation(i *DeleteTLSActivationInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := ToSafeURL("tls", "activations", i.ID)

	_, err := c.Delete(path, nil)
	return err
}
