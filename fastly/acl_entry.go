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
	ServiceID string `mapstructure:"service_id"`
	ACLID     string `mapstructure:"acl_id"`

	ID        string     `mapstructure:"id"`
	IP        string     `mapstructure:"ip"`
	Subnet    *int       `mapstructure:"subnet"`
	Negated   bool       `mapstructure:"negated"`
	Comment   string     `mapstructure:"comment"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
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
	ServiceID string
	ACLID     string
	Direction string
	PerPage   int
	Page      int
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
	consumed    bool
	CurrentPage int
	NextPage    int
	LastPage    int
	client      *Client
	options     *ListACLEntriesInput
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
	ServiceID string
	ACLID     string
	ID        string
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
	// Required fields
	ServiceID string
	ACLID     string
	IP        string `url:"ip"`

	// Optional fields
	Subnet  int         `url:"subnet,omitempty"`
	Negated Compatibool `url:"negated,omitempty"`
	Comment string      `url:"comment,omitempty"`
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
	ServiceID string
	ACLID     string
	ID        string
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
	// Required fields
	ServiceID string
	ACLID     string
	ID        string

	// Optional fields
	IP      *string      `url:"ip,omitempty"`
	Subnet  *int         `url:"subnet,omitempty"`
	Negated *Compatibool `url:"negated,omitempty"`
	Comment *string      `url:"comment,omitempty"`
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
	ServiceID string `json:"-"`
	ACLID     string `json:"-"`

	Entries []*BatchACLEntry `json:"entries"`
}

type BatchACLEntry struct {
	Operation BatchOperation `json:"op"`
	ID        *string        `json:"id,omitempty"`
	IP        *string        `json:"ip,omitempty"`
	Subnet    *int           `json:"subnet,omitempty"`
	Negated   *Compatibool   `json:"negated,omitempty"`
	Comment   *string        `json:"comment,omitempty"`
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
