package fastly

import (
	"strconv"
	"time"
)

const (
	// SnippetTypeInit sets the type to init.
	SnippetTypeInit SnippetType = "init"

	// SnippetTypeRecv sets the type to recv.
	SnippetTypeRecv SnippetType = "recv"

	// SnippetTypeHash sets the type to hash.
	SnippetTypeHash SnippetType = "hash"

	// SnippetTypeHit sets the type to hit.
	SnippetTypeHit SnippetType = "hit"

	// SnippetTypeMiss sets the type to miss.
	SnippetTypeMiss SnippetType = "miss"

	// SnippetTypePass sets the type to pass.
	SnippetTypePass SnippetType = "pass"

	// SnippetTypeFetch sets the type to fetch.
	SnippetTypeFetch SnippetType = "fetch"

	// SnippetTypeError sets the type to error.
	SnippetTypeError SnippetType = "error"

	// SnippetTypeDeliver sets the type to deliver.
	SnippetTypeDeliver SnippetType = "deliver"

	// SnippetTypeLog sets the type to log.
	SnippetTypeLog SnippetType = "log"

	// SnippetTypeNone sets the type to none.
	SnippetTypeNone SnippetType = "none"
)

// SnippetType is the type of VCL Snippet.
type SnippetType string

// Snippet is the Fastly Snippet object.
type Snippet struct {
	Content        *string      `mapstructure:"content"`
	CreatedAt      *time.Time   `mapstructure:"created_at"`
	DeletedAt      *time.Time   `mapstructure:"deleted_at"`
	Dynamic        *int         `mapstructure:"dynamic"`
	SnippetID      *string      `mapstructure:"id"`
	Name           *string      `mapstructure:"name"`
	Priority       *int         `mapstructure:"priority"`
	ServiceID      *string      `mapstructure:"service_id"`
	ServiceVersion *int         `mapstructure:"version"`
	Type           *SnippetType `mapstructure:"type"`
	UpdatedAt      *time.Time   `mapstructure:"updated_at"`
}

// CreateSnippetInput is the input for CreateSnippet.
type CreateSnippetInput struct {
	// Content is the VCL code that specifies exactly what the snippet does.
	Content *string `url:"content,omitempty"`
	// Dynamic sets the snippet version to regular (0) or dynamic (1).
	Dynamic *int `url:"dynamic,omitempty"`
	// Name is the name for the snippet (required).
	Name *string `url:"name,omitempty"`
	// Priority determines the ordering for multiple snippets. Lower numbers execute first.
	Priority *int `url:"priority,omitempty"`
	// ServiceID is the ID of the service to add the snippet to (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the editable configuration version (required).
	ServiceVersion int `url:"-"`
	// Type is the location in generated VCL where the snippet should be placed.
	Type *SnippetType `url:"type,omitempty"`
}

// CreateSnippet creates a new resource.
func (c *Client) CreateSnippet(i *CreateSnippetInput) (*Snippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "snippet")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippet *Snippet
	if err := DecodeBodyMap(resp.Body, &snippet); err != nil {
		return nil, err
	}
	return snippet, err
}

// UpdateSnippetInput is the input for UpdateSnippet.
type UpdateSnippetInput struct {
	// Content is the VCL code that specifies exactly what the snippet does.
	Content *string `url:"content,omitempty"`
	// Name is the name for the snippet (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Priority determines the ordering for multiple snippets. Lower numbers execute first.
	Priority *int `url:"priority,omitempty"`
	// ServiceID is the ID of the service to add the snippet to (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the editable configuration version (required).
	ServiceVersion int `url:"-"`
	// Type is the location in generated VCL where the snippet should be placed.
	Type *SnippetType `url:"type,omitempty"`
}

// UpdateSnippet updates the specified resource.
func (c *Client) UpdateSnippet(i *UpdateSnippetInput) (*Snippet, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "snippet", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippet *Snippet
	if err := DecodeBodyMap(resp.Body, &snippet); err != nil {
		return nil, err
	}
	return snippet, err
}

// DynamicSnippet is the object returned when updating or retrieving a Dynamic
// Snippet.
type DynamicSnippet struct {
	Content   *string    `mapstructure:"content"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	ServiceID *string    `mapstructure:"service_id"`
	SnippetID *string    `mapstructure:"snippet_id"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
}

// UpdateDynamicSnippetInput is the input for UpdateDynamicSnippet.
type UpdateDynamicSnippetInput struct {
	// Content is the VCL code that specifies exactly what the snippet does.
	Content *string `url:"content,omitempty"`
	// SnippetID is the SnippetID of the Snippet to modify (required)
	SnippetID string `url:"-"`
	// ServiceID is the ID of the Service to add the snippet to (required).
	ServiceID string `url:"-"`
}

// UpdateDynamicSnippet updates the specified resource.
func (c *Client) UpdateDynamicSnippet(i *UpdateDynamicSnippetInput) (*DynamicSnippet, error) {
	if i.SnippetID == "" {
		return nil, ErrMissingSnippetID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "snippet", i.SnippetID)

	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateSnippet *DynamicSnippet
	if err := DecodeBodyMap(resp.Body, &updateSnippet); err != nil {
		return nil, err
	}
	return updateSnippet, err
}

// DeleteSnippetInput is the input parameter to the DeleteSnippet function.
type DeleteSnippetInput struct {
	// Name is the Name of the Snippet to Delete (required).
	Name string
	// ServiceID is the ID of the Service to add the snippet to (required).
	ServiceID string
	// ServiceVersion is the editable configuration version (required).
	ServiceVersion int
}

// DeleteSnippet deletes the specified resource.
func (c *Client) DeleteSnippet(i *DeleteSnippetInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "snippet", i.Name)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
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

// ListSnippets retrieves all resources.
//
// Content is not displayed for Dynmanic Snippets due to them being
// versionless, use the GetDynamicSnippet function to show current content.
func (c *Client) ListSnippets(i *ListSnippetsInput) ([]*Snippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "snippet")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippets []*Snippet
	if err := DecodeBodyMap(resp.Body, &snippets); err != nil {
		return nil, err
	}
	return snippets, nil
}

// GetSnippetInput is used as input to the GetSnippet function.
type GetSnippetInput struct {
	// Name is the name of the Snippet to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetSnippet retrieves the specified resource.
//
// Dynamic Snippets will not show content due to them being versionless, use
// GetDynamicSnippet to see content.
func (c *Client) GetSnippet(i *GetSnippetInput) (*Snippet, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "snippet", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippet *Snippet
	if err := DecodeBodyMap(resp.Body, &snippet); err != nil {
		return nil, err
	}
	return snippet, nil
}

// GetDynamicSnippetInput is used as input to the GetDynamicSnippet function.
type GetDynamicSnippetInput struct {
	// SnippetID is the SnippetID of the Snippet to fetch (required).
	SnippetID string
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// GetDynamicSnippet retrieves the specified resource.
//
// This will show the current content associated with a Dynamic Snippet.
func (c *Client) GetDynamicSnippet(i *GetDynamicSnippetInput) (*DynamicSnippet, error) {
	if i.SnippetID == "" {
		return nil, ErrMissingSnippetID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "snippet", i.SnippetID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippet *DynamicSnippet
	if err := DecodeBodyMap(resp.Body, &snippet); err != nil {
		return nil, err
	}
	return snippet, nil
}
