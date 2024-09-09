package fastly

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Domain represents the the domain name Fastly will serve content for.
type Domain struct {
	Comment        *string    `mapstructure:"comment"`
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	Name           *string    `mapstructure:"name"`
	ServiceID      *string    `mapstructure:"service_id"`
	ServiceVersion *int       `mapstructure:"version"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}

// ListDomainsInput is used as input to the ListDomains function.
type ListDomainsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDomains retrieves all resources.
func (c *Client) ListDomains(i *ListDomainsInput) ([]*Domain, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "domain")

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ds []*Domain
	if err := decodeBodyMap(resp.Body, &ds); err != nil {
		return nil, err
	}
	return ds, nil
}

// CreateDomainInput is used as input to the CreateDomain function.
type CreateDomainInput struct {
	// Comment is a personal, freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// Name is the name of the domain that the service will respond to.
	Name *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// CreateDomain creates a new resource.
func (c *Client) CreateDomain(i *CreateDomainInput) (*Domain, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "domain")

	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Domain
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetDomainInput is used as input to the GetDomain function.
type GetDomainInput struct {
	// Name is the name of the domain to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetDomain retrieves the specified resource.
func (c *Client) GetDomain(i *GetDomainInput) (*Domain, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "domain", i.Name)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Domain
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// UpdateDomainInput is used as input to the UpdateDomain function.
type UpdateDomainInput struct {
	// Comment is a personal, freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// Name is the name of the domain that the service will respond to (required).
	Name string `url:"-"`
	// NewName is the updated name of the domain
	NewName *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// UpdateDomain updates the specified resource.
func (c *Client) UpdateDomain(i *UpdateDomainInput) (*Domain, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "domain", i.Name)

	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Domain
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteDomainInput is used as input to the DeleteDomain function.
type DeleteDomainInput struct {
	// Name is the name of the domain that the service will respond to (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteDomain deletes the specified resource.
func (c *Client) DeleteDomain(i *DeleteDomainInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "domain", i.Name)

	_, err := c.Delete(path, nil)
	return err
}

// ValidateDomainInput is used as input to the ValidateDomain function.
type ValidateDomainInput struct {
	// Name is the name of the domain to validate (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ValidateDomain validates the specified resource.
func (c *Client) ValidateDomain(i *ValidateDomainInput) (*DomainValidationResult, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "domain", i.Name, "check")

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var d *DomainValidationResult
	err = json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// DomainValidationResult defines an idiomatic representation of the API
// response.
type DomainValidationResult struct {
	CName    *string
	Metadata *DomainMetadata
	Valid    *bool
}

// UnmarshalJSON works around the badly designed API response by coercing the
// raw data into a more appropriate data structure.
func (d *DomainValidationResult) UnmarshalJSON(data []byte) error {
	var tuple []json.RawMessage
	if err := json.Unmarshal(data, &tuple); err != nil {
		return fmt.Errorf("initial: %w", err)
	}

	if want, have := 3, len(tuple); want != have {
		return fmt.Errorf("unexpected array length: want %d, have %d", want, have)
	}

	if err := json.Unmarshal(tuple[0], &d.Metadata); err != nil {
		return fmt.Errorf("metadata: %w", err)
	}

	if err := json.Unmarshal(tuple[1], &d.CName); err != nil {
		return fmt.Errorf("name: %w", err)
	}

	if err := json.Unmarshal(tuple[2], &d.Valid); err != nil {
		return fmt.Errorf("valid: %w", err)
	}

	return nil
}

// DomainMetadata represents a domain name configured for a Fastly service.
type DomainMetadata struct {
	Comment        *string    `json:"comment"`
	CreatedAt      *time.Time `json:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	Name           *string    `json:"name"`
	ServiceID      *string    `json:"service_id"`
	ServiceVersion *int       `json:"version"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

// ValidateAllDomainsInput is used as input to the ValidateAllDomains function.
type ValidateAllDomainsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ValidateAllDomains validates the specified resource.
func (c *Client) ValidateAllDomains(i *ValidateAllDomainsInput) (results []*DomainValidationResult, err error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "domain", "check_all")

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tuple []json.RawMessage
	if err := json.Unmarshal(data, &tuple); err != nil {
		return nil, fmt.Errorf("initial: %w", err)
	}
	for _, t := range tuple {
		var d *DomainValidationResult
		err = json.Unmarshal(t, &d)
		if err != nil {
			return nil, err
		}
		results = append(results, d)
	}

	return results, nil
}
