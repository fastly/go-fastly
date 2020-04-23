package fastly

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/jsonapi"
)

// GetPlatformPrivateKeyInput is an input to the GetPlatformPrivateKey function.
// Allowed values for the fields are described at https://docs.fastly.com/api/platform-tls.
type GetPlatformPrivateKeyInput struct {
	ID string
}

// PlatformPrivateKey .
type PlatformPrivateKey struct {
	KeyLength     int    `jsonapi:"attr,key_length"`
	KeyType       string `jsonapi:"attr,key_type"`
	Name          string `jsonapi:"attr,name"`
	CreatedAt     string `jsonapi:"attr,created_at"`
	Replace       bool   `jsonapi:"attr,replace"`
	PublicKeySHA1 string `jsonapi:"attr,public_key_sha1"`
}

// ListPrivateKeysInput .
type ListPrivateKeysInput struct {
	PageNumber  int    // The page index for pagination.
	PageSize    int    // The number of keys per page.
	FilterInUse string // Limit the returned keys to those without any matching TLS certificates.
}

// formatFilters converts user input into query parameters for filtering
func (i *ListPrivateKeysInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[in_use]": i.FilterInUse,
		"page[size]":     i.PageSize,
		"page[number]":   i.PageNumber,
	}
	// NOTE: This setup means we will not be able to send the zero value
	// of any of these filters. It doesn't appear we would need to at present.
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
		case "[]int":
			// convert ints to strings
			toStrings := []string{}
			values := value.([]int)
			for _, i := range values {
				toStrings = append(toStrings, strconv.Itoa(i))
			}
			// concat strings
			if len(values) > 0 {
				result[key] = strings.Join(toStrings, ",")
			}
		}
	}
	return result
}

// ListPlatformPrivateKeys list all TLS private keys.
func (c *Client) ListPlatformPrivateKeys(i *ListPrivateKeysInput) ([]*PlatformPrivateKey, error) {

	p := "/tls/private_keys"
	filters := &RequestOptions{Params: i.formatFilters()}

	r, err := c.Get(p, filters)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(PlatformPrivateKey)))
	if err != nil {
		return nil, err
	}

	ppk := make([]*PlatformPrivateKey, len(data))
	for i := range data {
		typed, ok := data[i].(*PlatformPrivateKey)
		if !ok {
			return nil, fmt.Errorf("got back a non-PlatformPrivateKey response")
		}
		ppk[i] = typed
	}

	return ppk, nil
}

// GetPlatformPrivateKey show a TLS private key.
func (c *Client) GetPlatformPrivateKey(i *GetPlatformPrivateKeyInput) (*PlatformPrivateKey, error) {

	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := fmt.Sprintf("/tls/private_keys/%s", i.ID)

	r, err := c.Get(p, nil)
	if err != nil {
		return nil, err
	}

	var ppk PlatformPrivateKey
	if err := jsonapi.UnmarshalPayload(r.Body, &ppk); err != nil {
		return nil, err
	}

	return &ppk, nil
}

// CreatePlatformPrivateKeyInput .
type CreatePlatformPrivateKeyInput struct {
	Key  string `jsonapi:"attr,key,omitempty"`
	Name string `jsonapi:"attr,name,omitempty"`
}

// CreatePlatformPrivateKey create a TLS private key.
func (c *Client) CreatePlatformPrivateKey(i *CreatePlatformPrivateKeyInput) (*PlatformPrivateKey, error) {

	p := "/tls/private_keys"

	if i.Key == "" {
		return nil, ErrMissingKey
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	r, err := c.PostJSONAPI(p, i, nil)
	if err != nil {
		return nil, err
	}

	var ppk PlatformPrivateKey
	if err := jsonapi.UnmarshalPayload(r.Body, &ppk); err != nil {
		return nil, err
	}

	return &ppk, nil
}

// DeletePrivateKeyInput used for deleting a private key.
type DeletePrivateKeyInput struct {
	ID string
}

// DeletePrivateKey destroy a TLS private key. only private keys not already matched to any certificates can be deleted.
func (c *Client) DeletePrivateKey(i *DeletePrivateKeyInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/private_keys/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}

// BulkCertificate .
type BulkCertificate struct {
	ID                string              `jsonapi:"primary,tls_bulk_certificate"`
	NotAfter          string              `jsonapi:"attr,not_after"`
	NotBefore         string              `jsonapi:"attr,not_before"`
	CreatedAt         string              `jsonapi:"attr,created_at"`
	UpdatedAt         string              `jsonapi:"attr,updated_at"`
	Replace           bool                `jsonapi:"attr,replace"`
	TLSConfigurations []*TLSConfiguration `jsonapi:"relation,tls_configurations,tls_configuration"`
	TLSDomains        []*TLSDomain        `jsonapi:"relation,tls_domains,tls_domain"`
}

// TLSConfiguration .
type TLSConfiguration struct {
	ID   string `jsonapi:"primary,tls_configuration"`
	Type string `jsonapi:"attr,type"`
}

// TLSDomain .
type TLSDomain struct {
	ID   string `jsonapi:"primary,tls_domain"`
	Type string `jsonapi:"attr,type"`
}

// ListBulkCertificatesInput .
type ListBulkCertificatesInput struct {
	PageNumber int // The page index for pagination.
	PageSize   int // The number of keys per page.
	// `tls_domains.id` filter seems to work where `tls_domain.id` does not, documentation wrong?
	// https://docs.fastly.com/api/platform-tls#tls_bulk_certificates_81cc5acbf847f71ecd4068ed58bfc5c5
	FilterTLSDomainIDMatch  string // Filter certificates by their matching, fully-qualified domain name. Returns all partial matches. Must provide a value longer than 3 characters.
	FilterTLSDomainsIDMatch string // Filter certificates by their matching, fully-qualified domain name. Returns all partial matches. Must provide a value longer than 3 characters.
	Sort                    string // The order in which to list certificates. Valid values are created_at, not_before, not_after. May precede any value with a - for descending.
}

// formatFilters converts user input into query parameters for filtering
func (i *ListBulkCertificatesInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[tls_domain.id][match]":  i.FilterTLSDomainIDMatch,
		"filter[tls_domains.id][match]": i.FilterTLSDomainsIDMatch,
		"page[size]":                    i.PageSize,
		"page[number]":                  i.PageNumber,
		"sort":                          i.Sort,
	}
	// NOTE: This setup means we will not be able to send the zero value
	// of any of these filters. It doesn't appear we would need to at present.
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
		case "[]int":
			// convert ints to strings
			toStrings := []string{}
			values := value.([]int)
			for _, i := range values {
				toStrings = append(toStrings, strconv.Itoa(i))
			}
			// concat strings
			if len(values) > 0 {
				result[key] = strings.Join(toStrings, ",")
			}
		}
	}
	return result
}

// ListBulkCertificates list all certificates.
func (c *Client) ListBulkCertificates(i *ListBulkCertificatesInput) ([]*BulkCertificate, error) {

	p := "/tls/bulk/certificates"
	filters := &RequestOptions{Params: i.formatFilters()}

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

// GetBulkCertificateInput .
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

// CreateBulkCertificateInput .
type CreateBulkCertificateInput struct {
	CertBlob          string              `jsonapi:"attr,cert_blob"`
	IntermediatesBlob string              `jsonapi:"attr,intermediates_blob"`
	TLSConfigurations []*TLSConfiguration `jsonapi:"relation,tls_configurations,tls_configuration"`
}

// CreateBulkCertificate create a TLS private key.
func (c *Client) CreateBulkCertificate(i *CreateBulkCertificateInput) (*BulkCertificate, error) {

	if i.CertBlob == "" {
		return nil, ErrMissingKey
	}
	if i.IntermediatesBlob == "" {
		return nil, ErrMissingName
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

// UpdateBulkCertificateInput .
type UpdateBulkCertificateInput struct {
	ID                string `jsonapi:"attr,id"`
	CertBlob          string `jsonapi:"attr,cert_blob"`
	IntermediatesBlob string `jsonapi:"attr,intermediates_blob"`
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
