package fastly

import (
	"fmt"
	"time"

	"github.com/google/jsonapi"
)

// GetPlatformPrivateKeysInput is an input to the GetPlatformPrivateKeys function.
// Allowed values for the fields are described at https://docs.fastly.com/api/platform-tls
type GetPlatformPrivateKeyInput struct {
	ID string
}

// PlatformPrivateKeysResponse is a response from the service tls API endpoint
type PlatformPrivateKeysResponse struct {
	Status  string                `mapstructure:"status"`
	Meta    map[string]string     `mapstructure:"meta"`
	Message string                `mapstructure:"msg"`
	Data    []*PlatformPrivateKey `mapstructure:"data"`
}

// PlatformPrivateKeyResponse is a response from the service tls API endpoint
type PlatformPrivateKeyResponse struct {
	Status  string              `mapstructure:"status"`
	Meta    map[string]string   `mapstructure:"meta"`
	Message string              `mapstructure:"msg"`
	Data    *PlatformPrivateKey `mapstructure:"data"`
}

// PlatformPrivateKey .
type PlatformPrivateKey struct {
	ID            string                        `mapstructure:"id"`
	Type          string                        `mapstructure:"type"`
	Attributes    PlatformPrivateKeysAttributes `mapstructure:"attributes"`
	Comment       string                        `mapstructure:"comment"`
	CustomerID    string                        `mapstructure:"customer_id"`
	CreatedAt     *time.Time                    `mapstructure:"created_at"`
	UpdatedAt     *time.Time                    `mapstructure:"updated_at"`
	DeletedAt     *time.Time                    `mapstructure:"deleted_at"`
	ActiveVersion uint                          `mapstructure:"version"`
	Versions      []*Version                    `mapstructure:"versions"`
}

// PlatformPrivateKeysAttributes .
type PlatformPrivateKeysAttributes struct {
	KeyLength     string     `mapstructure:"key_length"`
	KeyType       string     `mapstructure:"key_type"`
	Name          string     `mapstructure:"name"`
	CreatedAt     *time.Time `mapstructure:"created_at"`
	Replace       bool       `mapstructure:"replace"`
	PublicKeySHA1 string     `mapstructure:"public_key_sha1"`
}

// GetPlatformPrivateKeys returns stats data based on GetPlatformPrivateKeysInput
func (c *Client) GetPlatformPrivateKeys() (*PlatformPrivateKeysResponse, error) {

	p := "/tls/private_keys"

	r, err := c.Get(p, &RequestOptions{
		Params: map[string]string{},
	})
	if err != nil {
		return nil, err
	}

	var ppkr *PlatformPrivateKeysResponse
	if err := decodeJSON(&ppkr, r.Body); err != nil {
		return nil, err
	}

	return ppkr, nil
}

// GetPlatformPrivateKey returns a specific private key
func (c *Client) GetPlatformPrivateKey(i *GetPlatformPrivateKeyInput) (*PlatformPrivateKeyResponse, error) {

	p := "/tls/private_keys"

	if i.ID != "" {
		p = fmt.Sprintf("%s/%s", p, i.ID)
	}

	r, err := c.Get(p, &RequestOptions{
		Params: map[string]string{},
	})
	if err != nil {
		return nil, err
	}

	var ppkr *PlatformPrivateKeyResponse
	if err := decodeJSON(&ppkr, r.Body); err != nil {
		return nil, err
	}

	return ppkr, nil
}

// CreatePlatformPrivateKeyInput is an input to the CreatePlatformPrivateKey function.
// Allowed values for the fields are described at https://docs.fastly.com/api/platform-tls
type CreatePlatformPrivateKeyInput struct {
	Data CreatePlatformPrivateKeyData `form:"data"`
}

// CreatePlatformPrivateKeyData .
type CreatePlatformPrivateKeyData struct {
	Type       string                             `form:"type"`
	Attributes CreatePlatformPrivateKeyAttributes `form:"attributes"`
}

// CreatePlatformPrivateKeyAttributes .
type CreatePlatformPrivateKeyAttributes struct {
	Key  string `form:"subnet,omitempty"`
	Name string `form:"subnet,omitempty"`
}

// CreatePlatformPrivateKey returns a specific private key
func (c *Client) CreatePlatformPrivateKey(i *CreatePlatformPrivateKeyInput) (*PlatformPrivateKeyResponse, error) {

	p := "/tls/private_keys"

	if i.Data.Attributes.Key == "" {
		return nil, ErrMissingKey
	}

	if i.Data.Attributes.Name == "" {
		return nil, ErrMissingName
	}

	r, err := c.PostForm(p, i, nil)
	if err != nil {
		return nil, err
	}

	var ppkr *PlatformPrivateKeyResponse
	if err := decodeJSON(&ppkr, r.Body); err != nil {
		return nil, err
	}

	return ppkr, nil
}

type DeletePrivateKeyInput struct {
	ID string
}

// DeletePrivateKey deletes a specific private key
func (c *Client) DeletePrivateKey(i *DeletePrivateKeyInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/private_keys/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}

type TLSConfigurations struct {
	Data []TLSConfiguration `json:"data"`
}

type TLSConfiguration struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type TLSDomains struct {
	Data []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"data"`
}

type BulkCertificateResponseAttributes struct {
	Data []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"data"`
}

type BulkCertificateResponsRelationships struct {
	TLSConfigurations TLSConfigurations `json:"tls_configurations"`
	TLSDomains        TLSDomains        `json:"tls_domains"`
}

type BulkCertificate struct {
	ID            string                              `json:"id"`
	Type          string                              `json:"type"`
	Attributes    BulkCertificateResponseAttributes   `json:"attributes"`
	Relationships BulkCertificateResponsRelationships `json:"relationships"`
}

type GetBulkCertificateResponse struct {
	Data BulkCertificate `json:"data"`
}

type GetBulkCertificatesResponse struct {
	Data []BulkCertificate `json:"data"`
}

// GetBulkCertificates returns certificate data based on GetBulkCertificatesResponse
func (c *Client) GetBulkCertificates() (*GetBulkCertificatesResponse, error) {

	p := "/tls/bulk/certificates"

	r, err := c.Get(p, &RequestOptions{
		Params: map[string]string{},
	})
	if err != nil {
		return nil, err
	}

	var gbcr *GetBulkCertificatesResponse
	if err := decodeJSON(&gbcr, r.Body); err != nil {
		return nil, err
	}

	return gbcr, nil
}

type GetBulkCertificateInput struct {
	ID string
}

// GetBulkCertificate returns a specific certificate
func (c *Client) GetBulkCertificate(i *GetBulkCertificateInput) (*GetBulkCertificateResponse, error) {

	p := "/tls/bulk/certificates"

	if i.ID != "" {
		p = fmt.Sprintf("%s/%s", p, i.ID)
	}

	r, err := c.Get(p, &RequestOptions{
		Params: map[string]string{},
	})
	if err != nil {
		return nil, err
	}

	var bcr *GetBulkCertificateResponse
	if err := decodeJSON(&bcr, r.Body); err != nil {
		return nil, err
	}

	return bcr, nil
}

type CreateBulkCertificatesInput struct {
	Data CreateBulkCertificatesData `json:"data"`
}

type CreateBulkCertificatesData struct {
	Type          string                              `json:"type"`
	Attributes    CreateBulkCertificatesAttributes    `json:"attributes"`
	Relationships CreateBulkCertificatesRelationships `json:"relationships"`
}

type CreateBulkCertificatesAttributes struct {
	CertBlob          string `json:"cert_blob"`
	IntermediatesBlob string `json:"intermediates_blob"`
}

type CreateBulkCertificatesRelationships struct {
	TLSConfigurations TLSConfigurations `json:"tls_configurations"`
}

// CreateBulkCertificate uploads a new certificate
func (c *Client) CreateBulkCertificate(i *CreateBulkCertificatesInput) (*GetBulkCertificateResponse, error) {

	p := "/tls/bulk/certificates"

	r, err := c.PostForm(p, i, nil)
	if err != nil {
		return nil, err
	}

	var bcr *GetBulkCertificateResponse
	if err := decodeJSON(&bcr, r.Body); err != nil {
		return nil, err
	}

	return bcr, nil
}

type UpdateBulkCertificateInput struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	CertBlob          string `json:"cert_blob"`
	IntermediatesBlob string `json:"intermediates_blob"`
}

// validate makes sure the UpdateBulkCertificateInput instance has all
// fields we need to request a change. Almost, but not quite, identical to
// UpdateBulkCertificateInput.validate()
func (i UpdateBulkCertificateInput) validate() error {
	if i.ID == "" {
		return ErrMissingID
	}
	if i.Type == "" {
		return ErrMissingType
	}
	if i.CertBlob == "" {
		return ErrMissingCertBlob
	}
	if i.IntermediatesBlob == "" {
		return ErrMissingIntermediatesBlob
	}
	return nil
}

// UpdateBulkCertificate replace a certificate with a newly reissued certificate
func (c *Client) UpdateBulkCertificate(i *UpdateBulkCertificateInput) (*GetBulkCertificateResponse, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.Type == "" {
		return nil, ErrMissingType
	}

	if i.CertBlob == "" {
		return nil, ErrMissingOWASPID
	}

	path := fmt.Sprintf("/tls/bulk/certificates/%s", i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var gbcr GetBulkCertificateResponse
	if err := jsonapi.UnmarshalPayload(resp.Body, &gbcr); err != nil {
		return nil, err
	}
	return &gbcr, nil
}

type DeleteBulkCertificateInput struct {
	ID string
}

// DeleteBulkCertificate deletes a specific private key
func (c *Client) DeleteBulkCertificate(i *DeleteBulkCertificateInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/bulk/certificates/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}
