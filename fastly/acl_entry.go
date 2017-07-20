package fastly

import (
	"fmt"
	"sort"
)

type AclEntry struct {
	ServiceID string `mapstructure:"service_id"`
	AclID     string `mapstructure:"acl_id"`

	ID      string `mapstructure:"id"`
	IP      string `mapstructure:"ip"`
	Subnet  string `mapstructure:"subnet"`
	Negated bool   `mapstructure:"negated"`
	Comment string `mapstructure:"comment"`
}

// entriesById is a sortable list of Acl entries.
type entriesById []*AclEntry

// Len, Swap, and Less implements the sortable interface.
func (s entriesById) Len() int      { return len(s) }
func (s entriesById) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s entriesById) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// ListAclEntriesInput is the input parameter to ListAclEntries function.
type ListAclEntriesInput struct {
	Service string
	Acl     string
}

// ListAclEntries return a list of entries for an Acl
func (c *Client) ListAclEntries(i *ListAclEntriesInput) ([]*AclEntry, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Acl == "" {
		return nil, ErrMissingAcl
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entries", i.Service, i.Acl)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var es []*AclEntry
	if err := decodeJSON(&es, resp.Body); err != nil {
		return nil, err
	}

	sort.Stable(entriesById(es))

	return es, nil
}

// GetAclEntryInput is the input parameter to GetAclEntry function.
type GetAclEntryInput struct {
	Service string
	Acl     string
	ID      string
}

// GetAclEntry returns a single Acl entry based on its ID.
func (c *Client) GetAclEntry(i *GetAclEntryInput) (*AclEntry, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Acl == "" {
		return nil, ErrMissingAcl
	}

	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entry/%s", i.Service, i.Acl, i.ID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var e *AclEntry
	if err := decodeJSON(&e, resp.Body); err != nil {
		return nil, err
	}

	return e, nil
}

// CreateAclEntryInput the input parameter to CreateAclEntry function.
type CreateAclEntryInput struct {
	// Required fields
	Service string
	Acl     string
	IP      string `form:"ip"`

	// Optional fields
	Subnet  string `form:"subnet,omitempty"`
	Negated bool   `form:"negated,omitempty"`
	Comment string `form:"comment,omitempty"`
}

// CreateAclEntry creates and returns a new Acl entry.
func (c *Client) CreateAclEntry(i *CreateAclEntryInput) (*AclEntry, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Acl == "" {
		return nil, ErrMissingAcl
	}

	if i.IP == "" {
		return nil, ErrMissingIP
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entry", i.Service, i.Acl)

	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var e *AclEntry
	if err := decodeJSON(&e, resp.Body); err != nil {
		return nil, err
	}

	return e, nil
}

// DeleteAclEntryInput the input parameter to DeleteAclEntry function.
type DeleteAclEntryInput struct {
	Service string
	Acl     string
	ID      string
}

// DeleteAclEntry deletes an entry from an Acl based on its ID
func (c *Client) DeleteAclEntry(i *DeleteAclEntryInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Acl == "" {
		return ErrMissingAcl
	}

	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entry/%s", i.Service, i.Acl, i.ID)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeJSON(&r, resp.Body); err != nil {
		return err
	}

	if !r.Ok() {
		return fmt.Errorf("Not OK")
	}

	return nil

}

// UpdateAclEntryInput is the input parameter to UpdateAclEntry function.
type UpdateAclEntryInput struct {
	// Required fields
	Service string
	Acl     string
	ID      string

	// Optional fields
	IP      string `form:"ip,omitempty"`
	Subnet  string `form:"subnet,omitempty"`
	Negated bool   `form:"negated,omitempty"`
	Comment string `form:"comment,omitempty"`
}

// UpdateAclEntry updates an Acl entry
func (c *Client) UpdateAclEntry(i *UpdateAclEntryInput) (*AclEntry, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Acl == "" {
		return nil, ErrMissingAcl
	}

	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entry/%s", i.Service, i.Acl, i.ID)

	resp, err := c.RequestForm("PATCH", path, i, nil)
	if err != nil {
		return nil, err
	}

	var e *AclEntry
	if err := decodeJSON(&e, resp.Body); err != nil {
		return nil, err
	}

	return e, nil
}
