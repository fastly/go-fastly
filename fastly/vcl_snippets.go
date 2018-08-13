package fastly

import "fmt"

type Snippet struct {
	// Priority determines the ordering for multiple snippets. Lower numbers execute first.
	Priority int `json:"priority"`

	// Dynamic sets the snippet version to regular (0) or dynamic (1).
	Dynamic int `json:"dynamic"`

	// Name is the name for the snippet.
	Name string `json:"name"`

	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `json:"content"`

	// ID is the snippet ID
	ID string `json:"id"`

	// Type is the location in generated VCL where the snippet should be placed.
	Type string `json:"type"`

	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string `json:"service_id"`

	// Version is the editable version of the service.
	Version int `json:"version"`

	DeletedAt string `json:"deleted_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateSnippetInput struct {
	// Priority determines the ordering for multiple snippets. Lower numbers execute first.
	Priority int `form:"priority"`

	// Version is the editable version of the service
	Version int `json:"version"`

	// Dynamic sets the snippet version to regular (0) or dynamic (1).
	Dynamic int `form:"dynamic"`

	// Name is the name for the snippet.
	Name string `form:"name"`

	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `form:"content"`

	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string `json:"service_id"`

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

// UpdateSnippet is the object returned when updating a Dynamic Snippet
type UpdateSnippet struct {
	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string `json:"service_id"`

	// SnippetID is the ID of the Snippet to modify
	SnippetID string `json:"snippet_id"`

	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `json:"content"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UpdateSnippetInput is the input for updating a dynamic snippet
type UpdateSnippetInput struct {
	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string `json:"service_id"`

	// SnippetID is the ID of the Snippet to modify
	SnippetID string `json:"snippet_id"`

	// Content is the VCL code that specifies exactly what the snippet does.
	Content string `form:"content"`
}

func (c *Client) UpdateSnippet(i *UpdateSnippetInput) (*UpdateSnippet, error) {
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

	var updateSnippet *UpdateSnippet
	if err := decodeJSON(&updateSnippet, resp.Body); err != nil {
		return nil, err
	}
	return updateSnippet, err
}

type DeleteSnippetInput struct {
	// ServiceID is the ID of the Service to add the snippet to.
	ServiceID string `json:"service_id"`

	// SnippetName is the Name of the Snippet to Delete
	SnippetName string `json:"snippet"`

	// Version is the editable version of the service
	Version int `json:"version"`
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
