package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

const (
	// SnippetTypeInit sets the type to init
	SnippetTypeInit SnippetType = "init"

	// SnippetTypeRecv sets the type to recv
	SnippetTypeRecv SnippetType = "recv"

	// SnippetTypeHash sets the type to hash
	SnippetTypeHash SnippetType = "hash"

	// SnippetTypeHit sets the type to hit
	SnippetTypeHit SnippetType = "hit"

	// SnippetTypeMiss sets the type to miss
	SnippetTypeMiss SnippetType = "miss"

	// SnippetTypePass sets the type to pass
	SnippetTypePass SnippetType = "pass"

	// SnippetTypeFetch sets the type to fetch
	SnippetTypeFetch SnippetType = "fetch"

	// SnippetTypeError sets the type to error
	SnippetTypeError SnippetType = "error"

	// SnippetTypeDeliver sets the type to deliver
	SnippetTypeDeliver SnippetType = "deliver"

	// SnippetTypeLog sets the type to log
	SnippetTypeLog SnippetType = "log"

	// SnippetTypeNone sets the type to none
	SnippetTypeNone SnippetType = "none"
)

// SnippetType is the type of VCL Snippet
type SnippetType string

// SnippetTypeToString is a helper function to get a pointer to string
func SnippetTypeToString(b string) *SnippetType {
	p := SnippetType(b)
	return &p
}

// Snippet is the Fastly Snippet object
type Snippet struct {
	Content        string      `mapstructure:"content"`
	CreatedAt      *time.Time  `mapstructure:"created_at"`
	DeletedAt      *time.Time  `mapstructure:"deleted_at"`
	Dynamic        int         `mapstructure:"dynamic"`
	ID             string      `mapstructure:"id"`
	Name           string      `mapstructure:"name"`
	Priority       int         `mapstructure:"priority"`
	ServiceID      string      `mapstructure:"service_id"`
	ServiceVersion int         `mapstructure:"version"`
	Type           SnippetType `mapstructure:"type"`
	UpdatedAt      *time.Time  `mapstructure:"updated_at"`
}

// CreateSnippetInput is the input for CreateSnippet
type CreateSnippetInput struct {
	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `url:"content"`
	// Dynamic sets the snippet version to regular (0) or dynamic (1).
	Dynamic int `url:"dynamic"`
	// Name is the name for the snippet.
	Name string `url:"name"`
	// Priority determines the ordering for multiple snippets. Lower numbers execute first.
	Priority *int `url:"priority,omitempty"`
	// ServiceID is the ID of the service to add the snippet to (required).
	ServiceID string
	// ServiceVersion is the editable configuration version (required).
	ServiceVersion int
	// Type is the location in generated VCL where the snippet should be placed.
	Type SnippetType `url:"type"`
}

// CreateSnippet creates a new resource.
func (c *Client) CreateSnippet(i *CreateSnippetInput) (*Snippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	// 0 = versioned snippet
	// 1 = dynamic snippet
	if i.Dynamic == 0 && i.Content == "" {
		return nil, ErrMissingContent
	}

	if i.Type == "" {
		return nil, ErrMissingType
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippet *Snippet
	if err := decodeBodyMap(resp.Body, &snippet); err != nil {
		return nil, err
	}
	return snippet, err
}

// UpdateSnippetInput is the input for UpdateSnippet
type UpdateSnippetInput struct {
	// Content is the VCL code that specifies exactly what the snippet does.
	Content *string `url:"content,omitempty"`
	// Name is the name for the snippet.
	Name    string
	NewName *string `url:"name,omitempty"`
	// Priority determines the ordering for multiple snippets. Lower numbers execute first.
	Priority *int `url:"priority,omitempty"`
	// ServiceID is the ID of the service to add the snippet to (required).
	ServiceID string
	// ServiceVersion is the editable configuration version (required).
	ServiceVersion int
	// Type is the location in generated VCL where the snippet should be placed.
	Type *SnippetType `url:"type,omitempty"`
}

// UpdateSnippet updates a snippet on a unlocked version
func (c *Client) UpdateSnippet(i *UpdateSnippetInput) (*Snippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippet *Snippet
	if err := decodeBodyMap(resp.Body, &snippet); err != nil {
		return nil, err
	}
	return snippet, err
}

// DynamicSnippet is the object returned when updating or retrieving a Dynamic Snippet
type DynamicSnippet struct {
	Content   string     `mapstructure:"content"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	ID        string     `mapstructure:"snippet_id"`
	ServiceID string     `mapstructure:"service_id"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
}

// UpdateDynamicSnippetInput is the input for UpdateDynamicSnippet
type UpdateDynamicSnippetInput struct {
	// Content is the VCL code that specifies exactly what the snippet does.
	Content *string `url:"content,omitempty"`
	// ID is the ID of the Snippet to modify
	ID string
	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string
}

// UpdateDynamicSnippet replaces the content of a Dynamic Snippet
func (c *Client) UpdateDynamicSnippet(i *UpdateDynamicSnippetInput) (*DynamicSnippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/snippet/%s", i.ServiceID, i.ID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateSnippet *DynamicSnippet
	if err := decodeBodyMap(resp.Body, &updateSnippet); err != nil {
		return nil, err
	}
	return updateSnippet, err
}

// DeleteSnippetInput is the input parameter to the DeleteSnippet function.
type DeleteSnippetInput struct {
	// Name is the Name of the Snippet to Delete
	Name string
	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string
	// ServiceVersion is the editable configuration version (required).
	ServiceVersion int
}

// DeleteSnippet deletes the specified resource.
func (c *Client) DeleteSnippet(i *DeleteSnippetInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
		return ErrNotOK
	}
	return nil
}

// ListSnippetsInput is used as input to the ListSnippets function.
type ListSnippetsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// snippetsByName is a sortable list of Snippets.
type snippetsByName []*Snippet

// Len implement the sortable interface.
func (s snippetsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s snippetsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s snippetsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListSnippets retrieves all resources.
// versionless, use the GetDynamicSnippet function to show current content.
func (c *Client) ListSnippets(i *ListSnippetsInput) ([]*Snippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippets []*Snippet
	if err := decodeBodyMap(resp.Body, &snippets); err != nil {
		return nil, err
	}
	sort.Stable(snippetsByName(snippets))
	return snippets, nil
}

// GetSnippetInput is used as input to the GetSnippet function.
type GetSnippetInput struct {
	// Name is the name of the Snippet to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetSnippet gets the Snippet configuration with the given parameters. Dynamic Snippets will not show content due to them
// being versionless, use GetDynamicSnippet to see content.
func (c *Client) GetSnippet(i *GetSnippetInput) (*Snippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippet *Snippet
	if err := decodeBodyMap(resp.Body, &snippet); err != nil {
		return nil, err
	}
	return snippet, nil
}

// GetDynamicSnippetInput is used as input to the GetDynamicSnippet function.
type GetDynamicSnippetInput struct {
	// ID is the ID of the Snippet to fetch.
	ID string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// GetDynamicSnippet gets the Snippet configuration with the given parameters. This will show the current content
// associated with a Dynamic Snippet.
func (c *Client) GetDynamicSnippet(i *GetDynamicSnippetInput) (*DynamicSnippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/snippet/%s", i.ServiceID, i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippet *DynamicSnippet
	if err := decodeBodyMap(resp.Body, &snippet); err != nil {
		return nil, err
	}
	return snippet, nil
}
