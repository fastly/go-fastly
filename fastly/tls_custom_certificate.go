package fastly

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// CustomTLSCertificate represents a custom certificate. Uses common TLSDomain type from BulkCertificate.
type CustomTLSCertificate struct {
	CreatedAt          *time.Time   `jsonapi:"attr,created_at,iso8601"`
	Domains            []*TLSDomain `jsonapi:"relation,tls_domains"`
	ID                 string       `jsonapi:"primary,tls_certificate"`
	IssuedTo           string       `jsonapi:"attr,issued_to"`
	Issuer             string       `jsonapi:"attr,issuer"`
	Name               string       `jsonapi:"attr,name"`
	NotAfter           *time.Time   `jsonapi:"attr,not_after,iso8601"`
	NotBefore          *time.Time   `jsonapi:"attr,not_before,iso8601"`
	Replace            bool         `jsonapi:"attr,replace"`
	SerialNumber       string       `jsonapi:"attr,serial_number"`
	SignatureAlgorithm string       `jsonapi:"attr,signature_algorithm"`
	UpdatedAt          *time.Time   `jsonapi:"attr,updated_at,iso8601"`
}

// ListCustomTLSCertificatesInput is used as input to the Client.ListCustomTLSCertificates function.
type ListCustomTLSCertificatesInput struct {
	// FilterInUse limits the returned certificates to those currently using Fastly to terminate TLS (that is, certificates associated with an activation). Permitted values: true, false.
	FilterInUse *bool
	// FilterNotAfter limits the returned certificates to those that expire prior to the specified date in UTC. Accepts parameters: lte (e.g., filter[not_after][lte]=2020-05-05).
	FilterNotAfter string
	// FilterTLSDomainsID limits the returned certificates to those that include the specific domain.
	FilterTLSDomainsID string
	// Include captures related objects. Optional, comma-separated values. Permitted values: tls_activations.
	Include string
	// PageNumber is the page index for pagination.
	PageNumber int
	// PageSize is the number of keys per page.
	PageSize int
	// Sort is the order in which to list certificates. Valid values are created_at, not_before, not_after. May precede any value with a - for descending.
	Sort string
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListCustomTLSCertificatesInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"filter[in_use]":         i.FilterInUse,
		"filter[not_after]":      i.FilterNotAfter,
		"filter[tls_domains.id]": i.FilterTLSDomainsID,
		"include":                i.Include,
		"page[size]":             i.PageSize,
		"page[number]":           i.PageNumber,
		"sort":                   i.Sort,
	}

	for key, value := range pairings {
		switch t := value.(type) {
		case string:
			if t != "" {
				result[key] = t
			}
		case int:
			if t != 0 {
				result[key] = strconv.Itoa(t)
			}
		case *bool:
			if t != nil {
				result[key] = strconv.FormatBool(*t)
			}
		}
	}

	return result
}

// ListCustomTLSCertificates retrieves all resources.
func (c *Client) ListCustomTLSCertificates(i *ListCustomTLSCertificatesInput) ([]*CustomTLSCertificate, error) {
	path := "/tls/certificates"
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

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CustomTLSCertificate)))
	if err != nil {
		return nil, err
	}

	cc := make([]*CustomTLSCertificate, len(data))
	for i := range data {
		typed, ok := data[i].(*CustomTLSCertificate)
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T", data[i])
		}
		cc[i] = typed
	}

	return cc, nil
}

// GetCustomTLSCertificateInput is used as input to the GetCustomTLSCertificate function.
type GetCustomTLSCertificateInput struct {
	// ID is an alphanumeric string identifying a TLS certificate.
	ID string
}

// GetCustomTLSCertificate retrieves the specified resource.
func (c *Client) GetCustomTLSCertificate(i *GetCustomTLSCertificateInput) (*CustomTLSCertificate, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("tls", "certificates", i.ID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cc CustomTLSCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &cc); err != nil {
		return nil, err
	}

	return &cc, nil
}

// CreateCustomTLSCertificateInput is used as input to the CreateCustomTLSCertificate function.
type CreateCustomTLSCertificateInput struct {
	// CertBlob is the PEM-formatted certificate blob.
	CertBlob string `jsonapi:"attr,cert_blob"`
	// ID is an alphanumeric string identifying a TLS certificate.
	ID string `jsonapi:"primary,tls_certificate"` // ID value does not need to be set.
	// Name is a customizable name for your certificate.
	Name string `jsonapi:"attr,name,omitempty"`
}

// CreateCustomTLSCertificate creates a new resource.
func (c *Client) CreateCustomTLSCertificate(i *CreateCustomTLSCertificateInput) (*CustomTLSCertificate, error) {
	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}

	path := "/tls/certificates"

	resp, err := c.PostJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cc CustomTLSCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &cc); err != nil {
		return nil, err
	}

	return &cc, nil
}

// UpdateCustomTLSCertificateInput is used as input to the UpdateCustomTLSCertificate function.
type UpdateCustomTLSCertificateInput struct {
	// CertBlob is the PEM-formatted certificate blob.
	CertBlob string `jsonapi:"attr,cert_blob"`
	// ID is an alphanumeric string identifying a TLS certificate.
	ID string `jsonapi:"primary,tls_certificate"`
	// Name is a customizable name for your certificate.
	Name string `jsonapi:"attr,name,omitempty"`
}

// UpdateCustomTLSCertificate updates the specified resource.
//
// By using this endpoint, the original certificate will cease to be used for future TLS handshakes.
// Thus, only SAN entries that appear in the replacement certificate will become TLS enabled.
// Any SAN entries that are missing in the replacement certificate will become disabled.
func (c *Client) UpdateCustomTLSCertificate(i *UpdateCustomTLSCertificateInput) (*CustomTLSCertificate, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}

	path := ToSafeURL("tls", "certificates", i.ID)

	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cc CustomTLSCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &cc); err != nil {
		return nil, err
	}
	return &cc, nil
}

// DeleteCustomTLSCertificateInput used for deleting a certificate.
type DeleteCustomTLSCertificateInput struct {
	// ID is an alphanumeric string identifying a TLS certificate.
	ID string
}

// DeleteCustomTLSCertificate deletes the specified resource.
func (c *Client) DeleteCustomTLSCertificate(i *DeleteCustomTLSCertificateInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := ToSafeURL("tls", "certificates", i.ID)

	_, err := c.Delete(path, nil)
	return err
}
