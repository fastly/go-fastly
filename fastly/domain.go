package fastly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"time"
)

// Domain represents the the domain name Fastly will serve content for.
type Domain struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name      string     `mapstructure:"name"`
	Comment   string     `mapstructure:"comment"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
}

// domainsByName is a sortable list of backends.
type domainsByName []*Domain

// Len, Swap, and Less implement the sortable interface.
func (s domainsByName) Len() int      { return len(s) }
func (s domainsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s domainsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListDomainsInput is used as input to the ListDomains function.
type ListDomainsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDomains returns the list of domains for this Service.
func (c *Client) ListDomains(i *ListDomainsInput) ([]*Domain, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/domain", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ds []*Domain
	if err := decodeBodyMap(resp.Body, &ds); err != nil {
		return nil, err
	}
	sort.Stable(domainsByName(ds))
	return ds, nil
}

// CreateDomainInput is used as input to the CreateDomain function.
type CreateDomainInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the domain that the service will respond to (required).
	Name string `form:"name"`

	// Comment is a personal, freeform descriptive note.
	Comment string `form:"comment,omitempty"`
}

// CreateDomain creates a new domain with the given information.
func (c *Client) CreateDomain(i *CreateDomainInput) (*Domain, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/domain", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var d *Domain
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetDomainInput is used as input to the GetDomain function.
type GetDomainInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the domain to fetch.
	Name string `form:"name"`
}

// GetDomain retrieves information about the given domain name.
func (c *Client) GetDomain(i *GetDomainInput) (*Domain, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/domain/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var d *Domain
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// UpdateDomainInput is used as input to the UpdateDomain function.
type UpdateDomainInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the domain that the service will respond to (required).
	Name string

	// NewName is the updated name of the domain
	NewName *string `form:"name,omitempty"`

	// Comment is a personal, freeform descriptive note.
	Comment *string `form:"comment,omitempty"`
}

// UpdateDomain updates a single domain for the current service. The only allowed
// parameters are `Name` and `Comment`.
func (c *Client) UpdateDomain(i *UpdateDomainInput) (*Domain, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	if i.NewName == nil && i.Comment == nil {
		return nil, ErrMissingOptionalNameComment
	}

	path := fmt.Sprintf("/service/%s/version/%d/domain/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var d *Domain
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteDomainInput is used as input to the DeleteDomain function.
type DeleteDomainInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the domain that the service will respond to (required).
	Name string `form:"name"`
}

// DeleteDomain removes a single domain by the given name.
func (c *Client) DeleteDomain(i *DeleteDomainInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/domain/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	_, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	return nil
}

// ValidateDomainInput is used as input to the ValidateDomain function.
type ValidateDomainInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the domain to validate.
	Name string `form:"name"`
}

// ValidateDomain validates the given domain.
func (c *Client) ValidateDomain(i *ValidateDomainInput) (*DomainValidationResult, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/domain/%s/check", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
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
	Metadata DomainMetadata
	CName    string
	Valid    bool
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
	ServiceID      string `json:"service_id"`
	ServiceVersion int    `json:"version"`

	Name      string     `json:"name"`
	Comment   string     `json:"comment"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ValidateAllDomainsInput is used as input to the ValidateAllDomains function.
type ValidateAllDomainsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ValidateAllDomains validates the given domain.
func (c *Client) ValidateAllDomains(i *ValidateAllDomainsInput) (results []*DomainValidationResult, err error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/domain/check_all", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
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
