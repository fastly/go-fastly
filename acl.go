package fastly

import (
	"fmt"
	"sort"
)

type Acl struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Name string `mapstructure:"name"`
	ID   string `mapstructure:"id"`
}

// aclsByName is a sortable list of ACLs.
type aclsByName []*Acl

// Len, Swap, and Less implement the sortable interface.
func (s aclsByName) Len() int      { return len(s) }
func (s aclsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s aclsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListAclsInput is used as input to the ListAcls function.
type ListAclsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListAcls returns the list of ACLs for the configuration version.
func (c *Client) ListAcls(i *ListAclsInput) ([]*Acl, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var as []*Acl
	if err := decodeJSON(&as, resp.Body); err != nil {
		return nil, err
	}
	sort.Stable(aclsByName(as))
	return as, nil
}

// CreateAclInput is used as input to the CreateAcl function.
type CreateAclInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the ACL to create (required)
	Name string `form:"name"`
}

func (c *Client) CreateAcl(i *CreateAclInput) (*Acl, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var a *Acl
	if err := decodeJSON(&a, resp.Body); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteAclInput is the input parameter to DeleteAcl function.
type DeleteAclInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the acl to delete (required).
	Name string
}

// DeleteAcl deletes the given acl version.
func (c *Client) DeleteAcl(i *DeleteAclInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl/%s", i.Service, i.Version, i.Name)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeJSON(&r, resp.Body); err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf("Not Ok")
	}
	return nil
}

// GetAclInput is the input parameter to GetAcl function.
type GetAclInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the acl to get (required).
	Name string
}

// GetAcl gets the ACL configuration with the given parameters.
func (c *Client) GetAcl(i *GetAclInput) (*Acl, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl/%s", i.Service, i.Version, i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var a *Acl
	if err := decodeJSON(&a, resp.Body); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateAclInput is the input parameter to UpdateAcl function.
type UpdateAclInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the acl to update (required).
	Name string

	// NewName is the new name of the acl to update (required).
	NewName string `form:"name"`
}

// UpdateAcl updates the name of the ACL with the given parameters.
func (c *Client) UpdateAcl(i *UpdateAclInput) (*Acl, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	if i.NewName == "" {
		return nil, ErrMissingNewName
	}

	path := fmt.Sprintf("/service/%s/version/%d/acl/%s", i.Service, i.Version, i.Name)
	resp, err := c.PutForm(path, i, nil)

	if err != nil {
		return nil, err
	}

	var a *Acl
	if err := decodeJSON(&a, resp.Body); err != nil {
		return nil, err
	}

	return a, nil
}
