package fastly

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// BulkCertificate represents a bulk certificate.
type BulkCertificate struct {
	Configurations []*TLSConfiguration `jsonapi:"relation,tls_configurations,tls_configuration"`
	CreatedAt      *time.Time          `jsonapi:"attr,created_at,iso8601"`
	Domains        []*TLSDomain        `jsonapi:"relation,tls_domains,tls_domain"`
	ID             string              `jsonapi:"primary,tls_bulk_certificate"`
	NotAfter       *time.Time          `jsonapi:"attr,not_after,iso8601"`
	NotBefore      *time.Time          `jsonapi:"attr,not_before,iso8601"`
	Replace        bool                `jsonapi:"attr,replace"`
	UpdatedAt      *time.Time          `jsonapi:"attr,updated_at,iso8601"`
}

// TLSConfiguration represents the dedicated IP address pool that will be used to route traffic from the TLSDomain.
type TLSConfiguration struct {
	// ID is an alphanumeric string identifying a TLS configuration.
	ID string `jsonapi:"primary,tls_configuration"`
	// Type is a resource type (default: tls_configuration).
	Type string `jsonapi:"attr,type"`
}

// TLSDomain represents a domain (including wildcard domains) that is listed on a certificate's Subject Alternative Names (SAN) list.
type TLSDomain struct {
	// Activations is a list of TLS Activations.
	Activations []*TLSActivation `jsonapi:"relation,tls_activations,omitempty"`
	// Certificates is a list of Custom TLS Certificates.
	Certificates []*CustomTLSCertificate `jsonapi:"relation,tls_certificates,omitempty"`
	// ID is the domain name.
	ID string `jsonapi:"primary,tls_domain"`
	// Subscriptions is a list of TLS Subscriptions.
	Subscriptions []*TLSSubscription `jsonapi:"relation,tls_subscriptions,omitempty"`
	// Type is the resource type (default: tls_domain).
	Type string `jsonapi:"attr,type"`
}

// ListBulkCertificatesInput is used as input to the ListBulkCertificates function.
type ListBulkCertificatesInput struct {
	// FilterTLSDomainsIDMatch filters certificates by their matching, fully-qualified domain name. Returns all partial matches. Must provide a value longer than 3 characters.
	FilterTLSDomainsIDMatch string
	// PageNumber is the page index for pagination.
	PageNumber int
	// PageSize is the number of keys per page.
	PageSize int
	// Sort is the order in which to list certificates. Valid values are created_at, not_before, not_after. May precede any value with a - for descending.
	Sort string
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListBulkCertificatesInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"filter[tls_domains.id][match]": i.FilterTLSDomainsIDMatch,
		jsonapi.QueryParamPageSize:      i.PageSize,
		jsonapi.QueryParamPageNumber:    i.PageNumber,
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

// ListBulkCertificates retrieves all resources.
func (c *Client) ListBulkCertificates(ctx context.Context, i *ListBulkCertificatesInput) ([]*BulkCertificate, error) {
	path := "/tls/bulk/certificates"
	requestOptions := CreateRequestOptions()
	requestOptions.Params = i.formatFilters()
	requestOptions.Headers["Accept"] = jsonapi.MediaType // this is required otherwise the filters don't work

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(BulkCertificate)))
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
	// ID is an alphanumeric string identifying a TLS bulk certificate.
	ID string
}

// GetBulkCertificate retrieves the specified resource.
func (c *Client) GetBulkCertificate(ctx context.Context, i *GetBulkCertificateInput) (*BulkCertificate, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("tls", "bulk", "certificates", i.ID)

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bc BulkCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &bc); err != nil {
		return nil, err
	}

	return &bc, nil
}

// CreateBulkCertificateInput is used as input to the CreateBulkCertificate function.
type CreateBulkCertificateInput struct {
	// AllowUntrusted enables certificates that chain to untrusted roots.
	AllowUntrusted bool `jsonapi:"attr,allow_untrusted_root,omitempty"`
	// CertBlob is the PEM-formatted certificate blob.
	CertBlob string `jsonapi:"attr,cert_blob"`
	// Configurations is a list of TLS configurations.
	Configurations []*TLSConfiguration `jsonapi:"relation,tls_configurations,tls_configuration"`
	// IntermediatesBlob is the PEM-formatted chain of intermediate blobs.
	IntermediatesBlob string `jsonapi:"attr,intermediates_blob"`
}

// CreateBulkCertificate creates a new resource.
func (c *Client) CreateBulkCertificate(ctx context.Context, i *CreateBulkCertificateInput) (*BulkCertificate, error) {
	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}
	if i.IntermediatesBlob == "" {
		return nil, ErrMissingIntermediatesBlob
	}

	path := "/tls/bulk/certificates"

	resp, err := c.PostJSONAPI(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bc BulkCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &bc); err != nil {
		return nil, err
	}

	return &bc, nil
}

// UpdateBulkCertificateInput is used as input to the UpdateBulkCertificate function.
type UpdateBulkCertificateInput struct {
	// AllowUntrusted enables certificates that chain to untrusted roots.
	AllowUntrusted bool `jsonapi:"attr,allow_untrusted_root"`
	// CertBlob is the PEM-formatted certificate blob.
	CertBlob string `jsonapi:"attr,cert_blob"`
	// ID is an alphanumeric string identifying a TLS bulk certificate.
	ID string `jsonapi:"attr,id"`
	// IntermediatesBlob is the PEM-formatted chain of intermediate blobs.
	IntermediatesBlob string `jsonapi:"attr,intermediates_blob,omitempty"`
}

// UpdateBulkCertificate updates the specified resource.
//
// By using this endpoint, the original certificate will cease to be used for future TLS handshakes.
// Thus, only SAN entries that appear in the replacement certificate will become TLS enabled.
// Any SAN entries that are missing in the replacement certificate will become disabled.
func (c *Client) UpdateBulkCertificate(ctx context.Context, i *UpdateBulkCertificateInput) (*BulkCertificate, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.CertBlob == "" {
		return nil, ErrMissingCertBlob
	}

	if i.IntermediatesBlob == "" {
		return nil, ErrMissingIntermediatesBlob
	}

	path := ToSafeURL("tls", "bulk", "certificates", i.ID)

	resp, err := c.PatchJSONAPI(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bc BulkCertificate
	if err := jsonapi.UnmarshalPayload(resp.Body, &bc); err != nil {
		return nil, err
	}
	return &bc, nil
}

// DeleteBulkCertificateInput used for deleting a certificate.
type DeleteBulkCertificateInput struct {
	// ID is an alphanumeric string identifying a TLS bulk certificate.
	ID string
}

// DeleteBulkCertificate deletes the specified resource.
func (c *Client) DeleteBulkCertificate(ctx context.Context, i *DeleteBulkCertificateInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := ToSafeURL("tls", "bulk", "certificates", i.ID)

	ignored, err := c.Delete(ctx, path, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer ignored.Body.Close()
	return nil
}
