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
	ID                 string       `jsonapi:"primary,tls_certificate"`
	IssuedTo           string       `jsonapi:"attr,issued_to"`
	Issuer             string       `jsonapi:"attr,issuer"`
	Name               string       `jsonapi:"attr,name"`
	NotAfter           *time.Time   `jsonapi:"attr,not_after,iso8601"`
	NotBefore          *time.Time   `jsonapi:"attr,not_before,iso8601"`
	Replace            bool         `jsonapi:"attr,replace"`
	SerialNumber       string       `jsonapi:"attr,serial_number"`
	SignatureAlgorithm string       `jsonapi:"attr,signature_algorithm"`
	TLSDomains         []*TLSDomain `jsonapi:"relation,tls_domains"`
	CreatedAt          *time.Time   `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt          *time.Time   `jsonapi:"attr,updated_at,iso8601"`
}

// ListCustomTLSCertificatesInput is used as input to the ListCustomTLSCertificatesInput function.
type ListCustomTLSCertificatesInput struct {
	FilterNotAfter     string // Limit the returned certificates to those that expire prior to the specified date in UTC. Accepts parameters: lte (e.g., filter[not_after][lte]=2020-05-05).
	FilterTLSDomainsID string // Limit the returned certificates to those that include the specific domain.
	Include            string // Include related objects. Optional, comma-separated values. Permitted values: tls_activations.
	PageNumber         int    // The page index for pagination.
	PageSize           int    // The number of keys per page.
	Sort               string // The order in which to list certificates. Valid values are created_at, not_before, not_after. May precede any value with a - for descending.
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListCustomTLSCertificatesInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[not_after]":      i.FilterNotAfter,
		"filter[tls_domains.id]": i.FilterTLSDomainsID,
		"include":                i.Include,
		"page[size]":             i.PageSize,
		"page[number]":           i.PageNumber,
		"sort":                   i.Sort,
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

// ListCustomTLSCertificates list all certificates.
func (c *Client) ListCustomTLSCertificates(i *ListCustomTLSCertificatesInput) ([]*CustomTLSCertificate, error) {
	p := "/tls/certificates"
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

	data, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(CustomTLSCertificate)))
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
	ID string
}

func (c *Client) GetCustomTLSCertificate(i *GetCustomTLSCertificateInput) (*CustomTLSCertificate, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := fmt.Sprintf("/tls/certificates/%s", i.ID)

	r, err := c.Get(p, nil)
	if err != nil {
		return nil, err
	}

	var cc CustomTLSCertificate
	if err := jsonapi.UnmarshalPayload(r.Body, &cc); err != nil {
		return nil, err
	}

	return &cc, nil
}

// CreateCustomTLSCertificateInput is used as input to the CreateCustomTLSCertificate function.
type CreateCustomTLSCertificateInput struct {
	ID       string `jsonapi:"primary,tls_certificate"` // ID value does not need to be set.
	CertBlob string `jsonapi:"attr,cert_blob"`
	Name     string `jsonapi:"attr,name,omitempty"`
}

// CreateCustomTLSCertificate creates a custom TLS certificate.
func (c *Client) CreateCustomTLSCertificate(i *CreateCustomTLSCertificateInput) (*CustomTLSCertificate, error) {
	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}

	p := "/tls/certificates"

	r, err := c.PostJSONAPI(p, i, nil)
	if err != nil {
		return nil, err
	}

	var cc CustomTLSCertificate
	if err := jsonapi.UnmarshalPayload(r.Body, &cc); err != nil {
		return nil, err
	}

	return &cc, nil
}

// UpdateCustomTLSCertificateInput is used as input to the UpdateCustomTLSCertificate function.
type UpdateCustomTLSCertificateInput struct {
	ID       string `jsonapi:"primary,tls_certificate"`
	CertBlob string `jsonapi:"attr,cert_blob"`
	Name     string `jsonapi:"attr,name,omitempty"`
}

// UpdateCustomTLSCertificate replace a certificate with a newly reissued certificate.
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

	path := fmt.Sprintf("/tls/certificates/%s", i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var cc CustomTLSCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &cc); err != nil {
		return nil, err
	}
	return &cc, nil
}

// DeleteCustomTLSCertificateInput used for deleting a certificate.
type DeleteCustomTLSCertificateInput struct {
	ID string
}

// DeleteCustomTLSCertificate destroy a certificate. This disables TLS for all domains listed as SAN entries.
func (c *Client) DeleteCustomTLSCertificate(i *DeleteCustomTLSCertificateInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/certificates/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}
