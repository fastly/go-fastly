package fastly

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

// ACLEntriesPath is exposed primarily for use by the generic Paginator.
// See ./paginator.go for details.
const ACLEntriesPath = "/service/%s/acl/%s/entries"

// ACLEntry represents a server response from the Fastly API.
type ACLEntry struct {
	ACLID     string     `mapstructure:"acl_id"`
	Comment   string     `mapstructure:"comment"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
	ID        string     `mapstructure:"id"`
	IP        string     `mapstructure:"ip"`
	Negated   bool       `mapstructure:"negated"`
	ServiceID string     `mapstructure:"service_id"`
	Subnet    *int       `mapstructure:"subnet"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
}

// entriesByID is a sortable list of ACL entries.
type entriesByID []*ACLEntry

// Len implements the sortable interface.
func (s entriesByID) Len() int {
	return len(s)
}

// Swap implements the sortable interface.
func (s entriesByID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements the sortable interface.
func (s entriesByID) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// ListACLEntriesInput is the input parameter to ListACLEntries function.
type ListACLEntriesInput struct {
	// ACLID is an alphanumeric string identifying a ACL (required).
	ACLID string
	// Direction is the direction in which to sort results.
	Direction string
	// Page is the current page.
	Page int
	// PerPage is the number of records per page.
	PerPage int
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string
	// Sort is the field on which to sort.
	Sort string
}

func (l *ListACLEntriesInput) formatFilters() map[string]string {
	m := make(map[string]string)

	if l.Direction != "" {
		m["direction"] = l.Direction
	}
	if l.Page != 0 {
		m["page"] = strconv.Itoa(l.Page)
	}
	if l.PerPage != 0 {
		m["per_page"] = strconv.Itoa(l.PerPage)
	}
	if l.Sort != "" {
		m["sort"] = l.Sort
	}

	return m
}

// ListACLEntries retrieves all resources.
func (c *Client) ListACLEntries(i *ListACLEntriesInput) ([]*ACLEntry, error) {
	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf(ACLEntriesPath, i.ServiceID, i.ACLID)

	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var es []*ACLEntry
	if err := decodeBodyMap(resp.Body, &es); err != nil {
		return nil, err
	}

	sort.Stable(entriesByID(es))

	return es, nil
}

// GetACLEntryInput is the input parameter to GetACLEntry function.
type GetACLEntryInput struct {
	// ACLID is an alphanumeric string identifying an ACL Entry (required).
	ACLID string
	// ID is an alphanumeric string identifying an ACL Entry (required).
	ID string
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string
}

// GetACLEntry retrieves the specified resource.
func (c *Client) GetACLEntry(i *GetACLEntryInput) (*ACLEntry, error) {
	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}
	if i.ID == "" {
		return nil, ErrMissingID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entry/%s", i.ServiceID, i.ACLID, i.ID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *ACLEntry
	if err := decodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}

	return e, nil
}

// CreateACLEntryInput is the input parameter to the CreateACLEntry function.
type CreateACLEntryInput struct {
	// ACLID is an alphanumeric string identifying a ACL (required).
	ACLID string `url:"-"`
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// IP is an IP address.
	IP *string `url:"ip,omitempty"`
	// Negated is whether to negate the match. Useful primarily when creating individual exceptions to larger subnets.
	Negated *Compatibool `url:"negated,omitempty"`
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string `url:"-"`
	// Subnet is a number of bits for the subnet mask applied to the IP address.
	Subnet *int `url:"subnet,omitempty"`
}

// CreateACLEntry creates a new resource.
func (c *Client) CreateACLEntry(i *CreateACLEntryInput) (*ACLEntry, error) {
	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entry", i.ServiceID, i.ACLID)

	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *ACLEntry
	if err := decodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}

	return e, nil
}

// DeleteACLEntryInput the input parameter to DeleteACLEntry function.
type DeleteACLEntryInput struct {
	// ACLID is an alphanumeric string identifying a ACL (required).
	ACLID string
	// ID is an alphanumeric string identifying an ACL Entry (required).
	ID string
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string
}

// DeleteACLEntry deletes the specified resource.
func (c *Client) DeleteACLEntry(i *DeleteACLEntryInput) error {
	if i.ACLID == "" {
		return ErrMissingACLID
	}
	if i.ID == "" {
		return ErrMissingID
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entry/%s", i.ServiceID, i.ACLID, i.ID)

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

// UpdateACLEntryInput is the input parameter to UpdateACLEntry function.
type UpdateACLEntryInput struct {
	// ACLID is an alphanumeric string identifying a ACL (required).
	ACLID string `url:"-"`
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// ID is an alphanumeric string identifying an ACL Entry (required).
	ID string `url:"-"`
	// IP is an IP address.
	IP *string `url:"ip,omitempty"`
	// Negated is whether to negate the match. Useful primarily when creating individual exceptions to larger subnets.
	Negated *Compatibool `url:"negated,omitempty"`
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string `url:"-"`
	// Subnet is a number of bits for the subnet mask applied to the IP address.
	Subnet *int `url:"subnet,omitempty"`
}

// UpdateACLEntry updates the specified resource.
func (c *Client) UpdateACLEntry(i *UpdateACLEntryInput) (*ACLEntry, error) {
	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}
	if i.ID == "" {
		return nil, ErrMissingID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entry/%s", i.ServiceID, i.ACLID, i.ID)

	resp, err := c.RequestForm("PATCH", path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *ACLEntry
	if err := decodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}

	return e, nil
}

// BatchModifyACLEntriesInput is the input parameter to the
// BatchModifyACLEntries function.
type BatchModifyACLEntriesInput struct {
	// ACLID is an alphanumeric string identifying a ACL (required).
	ACLID string `json:"-"`
	// Entries is a list of ACL entries.
	Entries []*BatchACLEntry `json:"entries"`
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string `json:"-"`
}

// BatchACLEntry represents a single ACL entry.
type BatchACLEntry struct {
	// Comment is a freeform descriptive note.
	Comment *string `json:"comment,omitempty"`
	// ID is an alphanumeric string identifying an ACL Entry.
	ID *string `json:"id,omitempty"`
	// IP is an IP address.
	IP *string `json:"ip,omitempty"`
	// Negated is whether to negate the match. Useful primarily when creating individual exceptions to larger subnets.
	Negated *Compatibool `json:"negated,omitempty"`
	// Operation is a batching operation variant.
	Operation BatchOperation `json:"op"`
	// Subnet is the number of bits for the subnet mask applied to the IP address.
	Subnet *int `json:"subnet,omitempty"`
}

// BatchModifyACLEntries updates the specified resources.
func (c *Client) BatchModifyACLEntries(i *BatchModifyACLEntriesInput) error {
	if i.ACLID == "" {
		return ErrMissingACLID
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if len(i.Entries) > BatchModifyMaximumOperations {
		return ErrMaxExceededEntries
	}

	path := fmt.Sprintf(ACLEntriesPath, i.ServiceID, i.ACLID)
	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var batchModifyResult map[string]string

	return decodeBodyMap(resp.Body, &batchModifyResult)
}
