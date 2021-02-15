package fastly

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// BulkCertificate represents a bulk certificate.
type BulkCertificate struct {
	ID             string              `jsonapi:"primary,tls_bulk_certificate"`
	Configurations []*TLSConfiguration `jsonapi:"relation,tls_configurations,tls_configuration"`
	Domains        []*TLSDomain        `jsonapi:"relation,tls_domains,tls_domain"`
	NotAfter       *time.Time          `jsonapi:"attr,not_after,iso8601"`
	NotBefore      *time.Time          `jsonapi:"attr,not_before,iso8601"`
	CreatedAt      *time.Time          `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt      *time.Time          `jsonapi:"attr,updated_at,iso8601"`
	Replace        bool                `jsonapi:"attr,replace"`
}

// TLSConfiguration represents the dedicated IP address pool that will be used to route traffic from the TLSDomain.
type TLSConfiguration struct {
	ID   string `jsonapi:"primary,tls_configuration"`
	Type string `jsonapi:"attr,type"`
}

// TLSDomain represents a domain (including wildcard domains) that is listed on a certificate's Subject Alternative Names (SAN) list.
type TLSDomain struct {
	ID            string                  `jsonapi:"primary,tls_domain"`
	Type          string                  `jsonapi:"attr,type"`
	Activations   []*TLSActivation        `jsonapi:"relation,tls_activations,omitempty"`
	Certificates  []*CustomTLSCertificate `jsonapi:"relation,tls_certificates,omitempty"`
	Subscriptions []*TLSSubscription      `jsonapi:"relation,tls_subscriptions,omitempty"`
}

// ListBulkCertificatesInput is used as input to the ListBulkCertificates function.
type ListBulkCertificatesInput struct {
	PageNumber              int    // The page index for pagination.
	PageSize                int    // The number of keys per page.
	FilterTLSDomainsIDMatch string // Filter certificates by their matching, fully-qualified domain name. Returns all partial matches. Must provide a value longer than 3 characters.
	Sort                    string // The order in which to list certificates. Valid values are created_at, not_before, not_after. May precede any value with a - for descending.
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListBulkCertificatesInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[tls_domains.id][match]": i.FilterTLSDomainsIDMatch,
		"page[size]":                    i.PageSize,
		"page[number]":                  i.PageNumber,
		"sort":                          i.Sort,
	}
	for key, value := range pairings {
		switch v := value.(type) {
		case int:
			if v != 0 {
				result[key] = strconv.Itoa(v)
			}
		case string:
			if v != "" {
				result[key] = v
			}
		}
	}
	return result
}

// ListBulkCertificates list all certificates.
func (c *Client) ListBulkCertificates(i *ListBulkCertificatesInput) ([]*BulkCertificate, error) {

	p := "/tls/bulk/certificates"
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

	data, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(BulkCertificate)))
	if err != nil {
		return nil, err
	}

	bc := make([]*BulkCertificate, len(data))
	for i := range data {
		typed, ok := data[i].(*BulkCertificate)
		if !ok {
			return nil, fmt.Errorf("got back a non-BulkCertificate response")
		}
		bc[i] = typed
	}

	return bc, nil
}

// GetBulkCertificateInput is used as input to the GetBulkCertificate function.
type GetBulkCertificateInput struct {
	ID string
}

// GetBulkCertificate retrieve a single certificate.
func (c *Client) GetBulkCertificate(i *GetBulkCertificateInput) (*BulkCertificate, error) {

	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := fmt.Sprintf("/tls/bulk/certificates/%s", i.ID)

	r, err := c.Get(p, nil)
	if err != nil {
		return nil, err
	}

	var bc BulkCertificate
	if err := jsonapi.UnmarshalPayload(r.Body, &bc); err != nil {
		return nil, err
	}

	return &bc, nil
}

// CreateBulkCertificateInput is used as input to the CreateBulkCertificate function.
type CreateBulkCertificateInput struct {
	CertBlob          string              `jsonapi:"attr,cert_blob"`
	IntermediatesBlob string              `jsonapi:"attr,intermediates_blob"`
	AllowUntrusted    bool                `jsonapi:"attr,allow_untrusted_root,omitempty"`
	Configurations    []*TLSConfiguration `jsonapi:"relation,tls_configurations,tls_configuration"`
}

// CreateBulkCertificate create a TLS private key.
func (c *Client) CreateBulkCertificate(i *CreateBulkCertificateInput) (*BulkCertificate, error) {

	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}
	if i.IntermediatesBlob == "" {
		return nil, ErrMissingIntermediatesBlob
	}

	p := "/tls/bulk/certificates"

	r, err := c.PostJSONAPI(p, i, nil)
	if err != nil {
		return nil, err
	}

	var bc BulkCertificate
	if err := jsonapi.UnmarshalPayload(r.Body, &bc); err != nil {
		return nil, err
	}

	return &bc, nil
}

// UpdateBulkCertificateInput is used as input to the UpdateBulkCertificate function.
type UpdateBulkCertificateInput struct {
	ID                string `jsonapi:"attr,id"`
	CertBlob          string `jsonapi:"attr,cert_blob"`
	IntermediatesBlob string `jsonapi:"attr,intermediates_blob,omitempty"`
	AllowUntrusted    bool   `jsonapi:"attr,allow_untrusted_root"`
}

// UpdateBulkCertificate replace a certificate with a newly reissued certificate.
// By using this endpoint, the original certificate will cease to be used for future TLS handshakes.
// Thus, only SAN entries that appear in the replacement certificate will become TLS enabled.
// Any SAN entries that are missing in the replacement certificate will become disabled.
func (c *Client) UpdateBulkCertificate(i *UpdateBulkCertificateInput) (*BulkCertificate, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}

	if i.IntermediatesBlob == "" {
		return nil, ErrMissingIntermediatesBlob
	}

	path := fmt.Sprintf("/tls/bulk/certificates/%s", i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var bc BulkCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &bc); err != nil {
		return nil, err
	}
	return &bc, nil
}

// DeleteBulkCertificateInput used for deleting a certificate.
type DeleteBulkCertificateInput struct {
	ID string
}

// DeleteBulkCertificate destroy a certificate. This disables TLS for all domains listed as SAN entries.
func (c *Client) DeleteBulkCertificate(i *DeleteBulkCertificateInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/bulk/certificates/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}
