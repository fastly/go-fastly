package fastly

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/jsonapi"
)

// CustomCertificate represents a custom certificate. Uses common TLSDomain type from CustomCertificate
type CustomCertificate struct {
	ID                 string       `jsonapi:"primary,tls_certificate"`
	CreatedAt          *time.Time   `jsonapi:"attr,created_at,iso8601"`
	IssuedTo           string       `jsonapi:"attr,issued_to"`
	Issuer             string       `jsonapi:"attr,issuer"`
	Name               string       `jsonapi:"attr,name"`
	NotAfter           *time.Time   `jsonapi:"attr,not_after,iso8601"`
	NotBefore          *time.Time   `jsonapi:"attr,not_before,iso8601"`
	Replace            bool         `jsonapi:"attr,replace"`
	SerialNumber       string       `jsonapi:"attr,serial_number"`
	SignatureAlgorithm string       `jsonapi:"attr,signature_algorithm"`
	UpdatedAt          *time.Time   `jsonapi:"attr,updated_at,iso8601"`
	TLSDomains         []*TLSDomain `jsonapi:"relation,tls_domains,tls_domain"` // TODO "type" is not populating
}

// ListCustomCertificatesInput is used as input to the ListCustomCertificatesInput function.
type ListCustomCertificatesInput struct {
	PageNumber              *uint   // The page index for pagination.
	PageSize                *uint   // The number of keys per page.
	FilterTLSDomainsIDMatch *string // Filter certificates by their matching, fully-qualified domain name. Returns all partial matches. Must provide a value longer than 3 characters.
	Sort                    *string // The order in which to list certificates. Valid values are created_at, not_before, not_after. May precede any value with a - for descending.
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListCustomCertificatesInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[tls_domains.id][match]": i.FilterTLSDomainsIDMatch,
		"page[size]":                    i.PageSize,
		"page[number]":                  i.PageNumber,
		"sort":                          i.Sort,
	}
	for key, value := range pairings {
		if !reflect.ValueOf(value).IsNil() {
			result[key] = fmt.Sprintf("%v", reflect.ValueOf(value).Elem())
		}
	}
	return result
}

// ListCustomCertificates list all certificates.
func (c *Client) ListCustomCertificates(i *ListCustomCertificatesInput) ([]*CustomCertificate, error) {

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

	data, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(CustomCertificate)))
	if err != nil {
		return nil, err
	}

	cc := make([]*CustomCertificate, len(data))
	for i := range data {
		typed, ok := data[i].(*CustomCertificate)
		if !ok {
			return nil, fmt.Errorf("got back a non-CustomCertificate response")
		}
		cc[i] = typed
	}

	return cc, nil
}

// GetCustomCertificateInput is used as input to the GetCustomCertificate function.
type GetCustomCertificateInput struct {
	ID string
}

func (c *Client) GetCustomCertificate(i *GetCustomCertificateInput) (*CustomCertificate, error) {

	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := fmt.Sprintf("/tls/certificates/%s", i.ID)

	r, err := c.Get(p, nil)
	if err != nil {
		return nil, err
	}

	var cc CustomCertificate
	if err := jsonapi.UnmarshalPayload(r.Body, &cc); err != nil {
		return nil, err
	}

	return &cc, nil
}

type CreateCustomCertificateInput struct {
	CertBlob string `jsonapi:"attr,cert_blob"`
	Name     string `jsonapi:"attr,name"`
}

func (c *Client) CreateCustomCertificate(i *CreateCustomCertificateInput) (*CustomCertificate, error) { // TODO Possible platform bug, returns "Internal Server Error" if "type" is blank

	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}
	if i.Name == "" {
		return nil, ErrMissingName
	}

	p := "/tls/certificates"

	r, err := c.PostJSONAPI(p, i, nil)
	if err != nil {
		return nil, err
	}

	var cc CustomCertificate
	if err := jsonapi.UnmarshalPayload(r.Body, &cc); err != nil {
		return nil, err
	}

	return &cc, nil
}

type UpdateCustomCertificateInput struct {
	ID       string `jsonapi:"attr,id"`
	CertBlob string `jsonapi:"attr,cert_blob"`
	Name     string `jsonapi:"attr,name"`
}

// UpdateCustomCertificate replace a certificate with a newly reissued certificate.
// By using this endpoint, the original certificate will cease to be used for future TLS handshakes.
// Thus, only SAN entries that appear in the replacement certificate will become TLS enabled.
// Any SAN entries that are missing in the replacement certificate will become disabled.
func (c *Client) UpdateCustomCertificate(i *UpdateCustomCertificateInput) (*CustomCertificate, error) { // TODO TEST THIS
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/tls/certificates/%s", i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var cc CustomCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &cc); err != nil {
		return nil, err
	}
	return &cc, nil
}

// DeleteCustomCertificateInput used for deleting a certificate.
type DeleteCustomCertificateInput struct {
	ID string
}

// DeleteCustomCertificate destroy a certificate. This disables TLS for all domains listed as SAN entries.
func (c *Client) DeleteCustomCertificate(i *DeleteCustomCertificateInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/certificates/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}
