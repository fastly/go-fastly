package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// ACL represents a server response from the Fastly API.
type ACL struct {
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	ID             string     `mapstructure:"id"`
	Name           string     `mapstructure:"name"`
	ServiceID      string     `mapstructure:"service_id"`
	ServiceVersion int        `mapstructure:"version"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}

// ACLsByName is a sortable list of ACLs.
type ACLsByName []*ACL

// Len implements the sortable interface.
func (s ACLsByName) Len() int {
	return len(s)
}

// Swap implements the sortable interface.
func (s ACLsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements the sortable interface.
func (s ACLsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListACLsInput is used as input to the ListACLs function.
type ListACLsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListACLs retrieves all resources.
func (c *Client) ListACLs(i *ListACLsInput) ([]*ACL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var as []*ACL
	if err := decodeBodyMap(resp.Body, &as); err != nil {
		return nil, err
	}
	sort.Stable(ACLsByName(as))
	return as, nil
}

// CreateACLInput is used as input to the CreateACL function.
type CreateACLInput struct {
	// Name is the name of the ACL to create (required)
	Name *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// CreateACL creates a new resource.
func (c *Client) CreateACL(i *CreateACLInput) (*ACL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *ACL
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteACLInput is the input parameter to DeleteACL function.
type DeleteACLInput struct {
	// Name is the name of the ACL to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteACL deletes the specified resource.
func (c *Client) DeleteACL(i *DeleteACLInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf("not ok")
	}
	return nil
}

// GetACLInput is the input parameter to GetACL function.
type GetACLInput struct {
	// Name is the name of the ACL to get (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetACL retrieves the specified resource.
func (c *Client) GetACL(i *GetACLInput) (*ACL, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *ACL
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateACLInput is the input parameter to UpdateACL function.
type UpdateACLInput struct {
	// Name is the name of the ACL to update (required).
	Name string
	// NewName is the new name of the ACL to update (required).
	NewName *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// UpdateACL updates the specified resource.
func (c *Client) UpdateACL(i *UpdateACLInput) (*ACL, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *ACL
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}

	return a, nil
}
