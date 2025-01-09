package fastly

import (
	"bufio"
	"io"
	"net/http"
	"strconv"
)

// https://www.fastly.com/documentation/reference/api/acls/acls/

// ComputeACL represents a compute ACL response from the Fastly API.
type ComputeACL struct {
	Name         string `mapstructure:"name"`
	ComputeACLID string `mapstructure:"id"`
}

// CreateComputeACLInput is the input to the CreateComputeACL function.
type CreateComputeACLInput struct {
	// Name is the name of the compute ACL to create (required).
	Name string `json:"name"`
}

// CreateComputeACL creates a new resource.
func (c *Client) CreateComputeACL(i *CreateComputeACLInput) (*ComputeACL, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}

	const path = "/resources/acls"
	resp, err := c.PostJSON(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var acl *ComputeACL
	if err := DecodeBodyMap(resp.Body, &acl); err != nil {
		return nil, err
	}
	return acl, nil
}

// ListComputeACLsResponse retrieves all resources.
type ListComputeACLsResponse struct {
	// Data is the list of returned cumpute ACLs.
	Data []ComputeACL
	// Meta is the information for total compute ACLs.
	Meta map[string]string
}

// ListComputeACLs retrieves all compute ACLs.
func (c *Client) ListComputeACLs() (*ListComputeACLsResponse, error) {
	const path = "/resources/acls"

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var acls *ListComputeACLsResponse
	if err := DecodeBodyMap(resp.Body, &acls); err != nil {
		return nil, err
	}
	return acls, nil
}

// DescribeComputeACLInput is the input to the DescribeComputeACL function.
type DescribeComputeACLInput struct {
	// ACL Identifier (UUID). Required.
	ComputeACLID string
}

// DescribeComputeACL describes the specified resource.
func (c *Client) DescribeComputeACL(i *DescribeComputeACLInput) (*ComputeACL, error) {
	if i.ComputeACLID == "" {
		return nil, ErrMissingComputeACLID
	}

	path := ToSafeURL("resources", "acls", i.ComputeACLID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var acl *ComputeACL
	if err := DecodeBodyMap(resp.Body, &acl); err != nil {
		return nil, err
	}
	return acl, nil
}

// DeleteComputeACLInput is the input to the DeleteComputeACL function.
type DeleteComputeACLInput struct {
	// ACL Identifier (UUID). Required.
	ComputeACLID string
}

// DeleteComputeACL deletes the specified resource.
func (c *Client) DeleteComputeACL(i *DeleteComputeACLInput) error {
	if i.ComputeACLID == "" {
		return ErrMissingComputeACLID
	}

	path := ToSafeURL("resources", "acls", i.ComputeACLID)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewHTTPError(resp)
	}

	return nil
}

// ComputeACLEntry represents a compute ACL entry response from the Fastly API.
type ComputeACLEntry struct {
	Prefix string `mapstructure:"prefix"`
	Action string `mapstructure:"action"`
}

// ComputeACLLookupInput is the input to the ComputeACLLookup function.
type ComputeACLLookupInput struct {
	// ACL Identifier (UUID). Required.
	ComputeACLID string
	// Valid IPv4 or IPv6 address. Required.
	ComputeACLIP string
}

// ComputeACLLookup finds a matching ACL entry for an IP address.
func (c *Client) ComputeACLLookup(i *ComputeACLLookupInput) (*ComputeACLEntry, error) {
	if i.ComputeACLID == "" {
		return nil, ErrMissingComputeACLID
	}
	if i.ComputeACLIP == "" {
		return nil, ErrMissingComputeACLIP
	}

	path := ToSafeURL("resources", "acls", i.ComputeACLID, "entry", i.ComputeACLIP)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var entry *ComputeACLEntry
	if err := DecodeBodyMap(resp.Body, &entry); err != nil {
		return nil, err
	}
	return entry, nil
}

// ListComputeACLEntriesInput is used as an input to the ListComputeACLEntries function.
type ListComputeACLEntriesInput struct {
	// ACL Identifier (UUID). Required.
	ComputeACLID string
	// Cursor is used for paginating through results.
	Cursor string
	// Limit is the maximum number of entries included the response.
	Limit int
}

func (l *ListComputeACLEntriesInput) formatFilters() map[string]string {
	if l == nil {
		return nil
	}

	m := make(map[string]string)

	if l.Limit != 0 {
		m["limit"] = strconv.Itoa(l.Limit)
	}

	if l.Cursor != "" {
		m["cursor"] = l.Cursor
	}

	return m
}

// ListComputeACLEntriesResponse retrieves all entries of a compute ACL.
type ListComputeACLEntriesResponse struct {
	// Entries is the list of Compute ACL entries.
	Entries []ComputeACLEntry
	// Meta is the information for pagination.
	Meta map[string]string
}

// ListComputeACLEntries retrieves all entries of a compute ACL.
func (c *Client) ListComputeACLEntries(i *ListComputeACLEntriesInput) (*ListComputeACLEntriesResponse, error) {
	if i.ComputeACLID == "" {
		return nil, ErrMissingComputeACLID
	}

	path := ToSafeURL("resources", "acls", i.ComputeACLID, "entries")

	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var entries *ListComputeACLEntriesResponse
	if err := DecodeBodyMap(resp.Body, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// ListComputeACLEntriesPaginator is the opaque type for a ListComputeACLEntries call with pagination.
type ListComputeACLEntriesPaginator struct {
	client   *Client
	cursor   string
	err      error
	finished bool
	input    *ListComputeACLEntriesInput
	entries  []ComputeACLEntry
}

// NewListComputeACLEntriesPaginator returns a new paginator for the provided ListComputeACLEntriesInput.
func (c *Client) NewListComputeACLEntriesPaginator(i *ListComputeACLEntriesInput) PaginatorComputeACLEntries {
	return &ListComputeACLEntriesPaginator{
		client: c,
		input:  i,
	}
}

// Next advances the paginator.
func (l *ListComputeACLEntriesPaginator) Next() bool {
	if l.finished {
		l.entries = nil
		return false
	}

	l.input.Cursor = l.cursor
	o, err := l.client.ListComputeACLEntries(l.input)
	if err != nil {
		l.err = err
		l.finished = true
		return false
	}

	l.entries = o.Entries
	if next := o.Meta["next_cursor"]; next == "" {
		l.finished = true
	} else {
		l.cursor = next
	}

	return true
}

// Err returns any error from the paginator.
func (l *ListComputeACLEntriesPaginator) Err() error {
	return l.err
}

// Entries returns the current set of entries retrieved by the paginator.
func (l *ListComputeACLEntriesPaginator) Entries() []ComputeACLEntry {
	return l.entries
}

// BatchModifyComputeACLEntriesInput is the input to the BatchModifyComputeACLEntries function.
type BatchModifyComputeACLEntriesInput struct {
	// ACL Identifier (UUID). Required.
	ComputeACLID string
	// Body is the HTTP request body containing a JSON object of entries. Required.
	Body io.Reader
}

// BatchModifyComputeACLEntries streams an entries JSON object into a compute ACL.
// NOTE: We wrap the io.Reader with *bufio.Reader to handle large streams.
func (c *Client) BatchModifyComputeACLEntries(i *BatchModifyComputeACLEntriesInput) error {
	if i.ComputeACLID == "" {
		return ErrMissingComputeACLID
	}

	path := ToSafeURL("resources", "acls", i.ComputeACLID, "entries")

	resp, err := c.Patch(path, &RequestOptions{
		Body: bufio.NewReader(i.Body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Parallel: true,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = checkResp(resp, err)
	return err
}
