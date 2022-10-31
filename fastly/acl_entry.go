package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/peterhellberg/link"
)

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

// entriesById is a sortable list of ACL entries.
type entriesById []*ACLEntry

// Len, Swap, and Less implements the sortable interface.
func (s entriesById) Len() int      { return len(s) }
func (s entriesById) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s entriesById) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// ListACLEntriesInput is the input parameter to ListACLEntries function.
type ListACLEntriesInput struct {
	ACLID     string
	Direction string
	Page      int
	PerPage   int
	ServiceID string
	Sort      string
}

// ListACLEntries return a list of entries for an ACL
func (c *Client) ListACLEntries(i *ListACLEntriesInput) ([]*ACLEntry, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entries", i.ServiceID, i.ACLID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var es []*ACLEntry
	if err := decodeBodyMap(resp.Body, &es); err != nil {
		return nil, err
	}

	sort.Stable(entriesById(es))

	return es, nil
}

type ListAclEntriesPaginator struct {
	CurrentPage int
	LastPage    int
	NextPage    int

	// Private
	client   *Client
	consumed bool
	options  *ListACLEntriesInput
}

// HasNext returns a boolean indicating whether more pages are available
func (p *ListAclEntriesPaginator) HasNext() bool {
	return !p.consumed || p.Remaining() != 0
}

// Remaining returns the remaining page count
func (p *ListAclEntriesPaginator) Remaining() int {
	if p.LastPage == 0 {
		return 0
	}
	return p.LastPage - p.CurrentPage
}

// GetNext retrieves data in the next page
func (p *ListAclEntriesPaginator) GetNext() ([]*ACLEntry, error) {
	return p.client.listACLEntriesWithPage(p.options, p)
}

// NewListACLEntriesPaginator returns a new paginator
func (c *Client) NewListACLEntriesPaginator(i *ListACLEntriesInput) PaginatorACLEntries {
	return &ListAclEntriesPaginator{
		client:  c,
		options: i,
	}
}

// listACLEntriesWithPage return a list of entries for an ACL of a given page
func (c *Client) listACLEntriesWithPage(i *ListACLEntriesInput, p *ListAclEntriesPaginator) ([]*ACLEntry, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}

	var perPage int
	const maxPerPage = 100
	if i.PerPage <= 0 {
		perPage = maxPerPage
	} else {
		perPage = i.PerPage
	}

	// page is not specified, fetch from the beginning
	if i.Page <= 0 && p.CurrentPage == 0 {
		p.CurrentPage = 1
	} else {
		// page is specified, fetch from a given page
		if !p.consumed {
			p.CurrentPage = i.Page
		} else {
			p.CurrentPage = p.CurrentPage + 1
		}
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entries", i.ServiceID, i.ACLID)
	requestOptions := &RequestOptions{
		Params: map[string]string{
			"per_page": strconv.Itoa(perPage),
			"page":     strconv.Itoa(p.CurrentPage),
		},
	}

	if i.Direction != "" {
		requestOptions.Params["direction"] = i.Direction
	}
	if i.Sort != "" {
		requestOptions.Params["sort"] = i.Sort
	}

	resp, err := c.Get(path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	for _, l := range link.ParseResponse(resp) {
		// indicates the Link response header contained the next page instruction
		if l.Rel == "next" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.NextPage, _ = strconv.Atoi(query["page"][0])
		}
		// indicates the Link response header contained the last page instruction
		if l.Rel == "last" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.LastPage, _ = strconv.Atoi(query["page"][0])
		}
	}

	p.consumed = true

	var es []*ACLEntry
	if err := decodeBodyMap(resp.Body, &es); err != nil {
		return nil, err
	}

	sort.Stable(entriesById(es))

	return es, nil
}

// GetACLEntryInput is the input parameter to GetACLEntry function.
type GetACLEntryInput struct {
	ACLID     string
	ID        string
	ServiceID string
}

// GetACLEntry returns a single ACL entry based on its ID.
func (c *Client) GetACLEntry(i *GetACLEntryInput) (*ACLEntry, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}

	if i.ID == "" {
		return nil, ErrMissingID
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

// CreateACLEntryInput the input parameter to CreateACLEntry function.
type CreateACLEntryInput struct {
	ACLID     string
	IP        string `url:"ip"`
	ServiceID string
	Comment   string      `url:"comment,omitempty"`
	Negated   Compatibool `url:"negated,omitempty"`
	Subnet    int         `url:"subnet,omitempty"`
}

// CreateACLEntry creates and returns a new ACL entry.
func (c *Client) CreateACLEntry(i *CreateACLEntryInput) (*ACLEntry, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}

	if i.IP == "" {
		return nil, ErrMissingIP
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
	ACLID     string
	ID        string
	ServiceID string
}

// DeleteACLEntry deletes an entry from an ACL based on its ID
func (c *Client) DeleteACLEntry(i *DeleteACLEntryInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ACLID == "" {
		return ErrMissingACLID
	}

	if i.ID == "" {
		return ErrMissingID
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
	ACLID     string
	Comment   *string `url:"comment,omitempty"`
	ID        string
	IP        *string      `url:"ip,omitempty"`
	Negated   *Compatibool `url:"negated,omitempty"`
	ServiceID string
	Subnet    *int `url:"subnet,omitempty"`
}

// UpdateACLEntry updates an ACL entry
func (c *Client) UpdateACLEntry(i *UpdateACLEntryInput) (*ACLEntry, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ACLID == "" {
		return nil, ErrMissingACLID
	}

	if i.ID == "" {
		return nil, ErrMissingID
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

type BatchModifyACLEntriesInput struct {
	ACLID     string           `json:"-"`
	Entries   []*BatchACLEntry `json:"entries"`
	ServiceID string           `json:"-"`
}

type BatchACLEntry struct {
	Comment   *string        `json:"comment,omitempty"`
	ID        *string        `json:"id,omitempty"`
	IP        *string        `json:"ip,omitempty"`
	Negated   *Compatibool   `json:"negated,omitempty"`
	Operation BatchOperation `json:"op"`
	Subnet    *int           `json:"subnet,omitempty"`
}

func (c *Client) BatchModifyACLEntries(i *BatchModifyACLEntriesInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ACLID == "" {
		return ErrMissingACLID
	}

	if len(i.Entries) > BatchModifyMaximumOperations {
		return ErrMaxExceededEntries
	}

	path := fmt.Sprintf("/service/%s/acl/%s/entries", i.ServiceID, i.ACLID)
	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var batchModifyResult map[string]string
	if err := decodeBodyMap(resp.Body, &batchModifyResult); err != nil {
		return err
	}

	return nil
}
