package fastly

import "fmt"

type Snippet struct {
	// Priority determines the ordering for multiple snippets. Lower numbers execute first.
	Priority int `mapstructure:"priority"`

	// Dynamic sets the snippet version to regular (0) or dynamic (1).
	Dynamic int `mapstructure:"dynamic"`

	// Name is the name for the snippet.
	Name string `mapstructure:"name"`

	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `mapstructure:"content"`

	// ID is the snippet ID
	ID string `mapstructure:"id"`

	// Type is the location in generated VCL where the snippet should be placed.
	Type string `mapstructure:"type"`

	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string `mapstructure:"service_id"`

	// Version is the editable version of the service.
	Version int `mapstructure:"version"`

	DeletedAt string `mapstructure:"deleted_at"`
	CreatedAt string `mapstructure:"created_at"`
	UpdatedAt string `mapstructure:"updated_at"`
}

type CreateSnippetInput struct {
	// Priority determines the ordering for multiple snippets. Lower numbers execute first.
	Priority int `form:"priority"`

	// Version is the editable version of the service
	Version int

	// Dynamic sets the snippet version to regular (0) or dynamic (1).
	Dynamic int `form:"dynamic"`

	// Name is the name for the snippet.
	Name string `form:"name"`

	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `form:"content"`

	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string

	// Type is the location in generated VCL where the snippet should be placed.
	Type string `form:"type"`
}

func (c *Client) CreateSnippet(i *CreateSnippetInput) (*Snippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet", i.ServiceID, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var snippet *Snippet
	if err := decodeJSON(&snippet, resp.Body); err != nil {
		return nil, err
	}
	return snippet, err
}

// DynamicSnippet is the object returned when updating a Dynamic Snippet
type DynamicSnippet struct {
	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string `mapstructure:"service_id"`

	// SnippetID is the ID of the Snippet to modify
	SnippetID string `mapstructure:"snippet_id"`

	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `mapstructure:"content"`

	CreatedAt string `mapstructure:"created_at"`
	UpdatedAt string `mapstructure:"updated_at"`
}

// UpdateSnippetInput is the input for updating a dynamic snippet
type UpdateSnippetInput struct {
	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string

	// SnippetID is the ID of the Snippet to modify
	SnippetID string

	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `form:"content"`
}

func (c *Client) UpdateSnippet(i *UpdateSnippetInput) (*DynamicSnippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingService
	}

	if i.SnippetID == "" {
		return nil, ErrMissingSnippetID
	}

	path := fmt.Sprintf("/service/%s/snippet/%s", i.ServiceID, i.SnippetID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var updateSnippet *DynamicSnippet
	if err := decodeJSON(&updateSnippet, resp.Body); err != nil {
		return nil, err
	}
	return updateSnippet, err
}

type DeleteSnippetInput struct {
	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string

	// SnippetName is the Name of the Snippet to Delete
	SnippetName string

	// Version is the editable version of the service
	Version int
}

func (c *Client) DeleteSnippet(i *DeleteSnippetInput) error {
	if i.ServiceID == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.SnippetName == "" {
		return ErrMissingSnippetName
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.ServiceID, i.Version, i.SnippetName)
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

// ListSnippetsInput is used as input to the ListSnippets function.
type ListSnippetsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListSnippets returns the list of Snippets for the configuration version.
func (c *Client) ListSnippets(i *ListSnippetsInput) ([]*Snippet, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var snippets []*Snippet
	if err := decodeJSON(&snippets, resp.Body); err != nil {
		return nil, err
	}
	return snippets, nil
}

// GetSnippetbyNameInput is used as input to the GetSnippet function.
type GetSnippetbyNameInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the Snippet to fetch.
	Name string
}

// GetSnippetByName gets the Snippet configuration with the given parameters.
func (c *Client) GetSnippetByName(i *GetSnippetbyNameInput) (*Snippet, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.Service, i.Version, i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var snippet *Snippet
	if err := decodeJSON(&snippet, resp.Body); err != nil {
		return nil, err
	}
	return snippet, nil
}

// GetSnippetByIDInput is used as input to the GetSnippet function.
type GetSnippetByIDInput struct {
	// Service is the ID of the service.
	Service string

	// ID is the ID of the Snippet to fetch.
	ID string
}

// GetSnippetByID gets the Snippet configuration with the given parameters.
func (c *Client) GetSnippetByID(i *GetSnippetByIDInput) (*DynamicSnippet, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.ID == "" {
		return nil, ErrMissingSnippetID
	}

	path := fmt.Sprintf("/service/%s/snippet/%s", i.Service, i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var snippet *DynamicSnippet
	if err := decodeJSON(&snippet, resp.Body); err != nil {
		return nil, err
	}
	return snippet, nil
}
