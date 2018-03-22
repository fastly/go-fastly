package fastly

import (
	"fmt"
	"sort"
	"time"
)

// VCLSnippet represents a response about VCLSnippet from the Fastly API.
type VCLSnippet struct {
	Service string `mapstructure:"service_id"`
	Version int    `mapstructure:"version"`

	Name      string     `mapstructure:"name"`
	Type      string     `mapstructure:"type"`
	Content   string     `mapstructure:"content"`
	Priority  int        `mapstructure:"priority"`
	Dynamic   int        `mapstructure:"dynamic"`
	ID        string     `mapstructure:"id"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
}

// vclSnippetsByName is a sortable list of VCLSnippets.
type vclSnippetsByName []*VCLSnippet

// Len, Swap, and Less implement the sortable interface.
func (s vclSnippetsByName) Len() int {
	return len(s)
}
func (s vclSnippetsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s vclSnippetsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListVCLSnippetsInput is used as input to the ListVCLSnippets function.
type ListVCLSnippetsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListVCLSnippets returns the list all VCL snippets for a particular service and version.
func (c *Client) ListVCLSnippets(i *ListVCLSnippetsInput) ([]*VCLSnippet, error) {
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

	var vclSnippets []*VCLSnippet
	if err := decodeJSON(&vclSnippets, resp.Body); err != nil {
		return nil, err
	}
	sort.Stable(vclSnippetsByName(vclSnippets))
	return vclSnippets, nil
}

// GetVCLSnippetInput is used as input to the GetVCL function.
type GetVCLSnippetInput struct {
	// Service is the ID of the service. Required.
	Service string

	// Either Version & Name OR just ID need to additionally required.
	// Service is the version of the service. Required only if Name is set and ID is not.
	Version int

	// Name is the name of the VCL Snippet to fetch.
	Name string

	// ID is the id of the VCL Snippet to fetch
	ID string
}

// hasID determines whether input has ID property.
func hasID(i *GetVCLSnippetInput) bool {
	return (i.ID != "")
}

// hasSnippetNameAndServiceVersion determines whether input has Name and Version properties.
func hasSnippetNameAndServiceVersion(i *GetVCLSnippetInput) bool {
	return (i.Name != "" && i.Version != 0)
}

// GetVCLSnippet gets a VCL snippet with the given parameters.
func (c *Client) GetVCLSnippet(i *GetVCLSnippetInput) (*VCLSnippet, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	// Have required params?
	if !hasID(i) && !hasSnippetNameAndServiceVersion(i) {
		if i.Version == 0 {
			return nil, ErrMissingVersion
		}
		if i.Name == "" {
			return nil, ErrMissingName
		}
		if i.ID == "" {
			return nil, ErrMissingID
		}
	}

	var path string
	// Maybe get a single snippet for a particular service and version.
	if hasSnippetNameAndServiceVersion(i) {
		path = fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.Service, i.Version, i.Name)
		// Maybe get a single dynamic snippet for a particular service.
	} else if hasID(i) {
		path = fmt.Sprintf("/service/%s/snippet/%s", i.Service, i.ID)
	}
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var vclSnippet *VCLSnippet
	if err := decodeJSON(&vclSnippet, resp.Body); err != nil {
		return nil, err
	}
	return vclSnippet, nil
}

// GetVCLSnippetByID gets a VCL snippet by service name and snippet ID.
func (c *Client) GetVCLSnippetByID(service string, id string) (*VCLSnippet, error) {
	return c.GetVCLSnippet(&GetVCLSnippetInput{
		Service: service,
		ID:      id,
	})
}

// CreateVCLSnippetInput is used as input to the CreateVCLSnippet function.
type CreateVCLSnippetInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name of VCL snippet. Required.
	Name string `form:"name"`

	// Content of VCL snippet.
	// Default: ""
	Content string `form:"content"`

	// Whether VCL snippet is dynamic or not.
	// Default: 0
	Dynamic int `form:"dynamic"`

	// Type of VCL Snippet. This can be init or a sub-routine. Either: init, recv, hit, miss, pass, fetch, error, deliver, log, or none.
	// Default: none
	Type string `form:"type"`

	// Priority of VCL Snippet.
	// Default: 100
	Priority int `form:"priority"`

	// Whether to Activate
	// Default: false
	Activate bool `form:"activate"`
}

// setDefaults sets the defaults for the CreateVCLSnippetInput.
func (self *CreateVCLSnippetInput) setDefaults() {
	if self.Type == "" {
		self.Type = "none"
	}
	if self.Priority == 0 {
		self.Priority = 100
	}
}

// CreateVCLSnippet creates a VCL snippet for a particular service and version.
func (c *Client) CreateVCLSnippet(i *CreateVCLSnippetInput) (*VCLSnippet, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	i.setDefaults()
	path := fmt.Sprintf("/service/%s/version/%d/snippet", i.Service, i.Version)

	// Create a new Version

	// Create Snippet
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	var vclSnippet *VCLSnippet

	if err := decodeJSON(&vclSnippet, resp.Body); err != nil {
		return nil, err
	}

	// Validate
	_, err = c.ValidateVCL(&ValidateVCLInput{
		Service: i.Service,
		Version: i.Version,
	})

	if err != nil {
		// Remove Dynamic Snippet
		c.DeleteVCLSnippet(&DeleteVCLSnippetInput{
			Service: i.Service,
			Version: i.Version,
			Name:    i.Name,
		})
		return nil, err
	}

	// To Activate, Need the Main VCL
	vcls, err := c.ListVCLs(&ListVCLsInput{
		Service: i.Service,
		Version: i.Version,
	})
	for _, v := range vcls {
		if v.Main {
			// Activate
			c.ActivateVCL(&ActivateVCLInput{
				Service: i.Service,
				Version: i.Version,
				Name:    v.Name,
			})
		}
	}

	return vclSnippet, nil
}

// UpdateDynamicVCLSnippetInput is used as input to the UpdateVCL function.
type UpdateDynamicVCLSnippetInput struct {
	// Service is the ID of the service (required).
	Service string

	// ID is the ID of the VCL Snippet to update (required).
	ID string

	// New Name of VCL snippet.
	Name string `form:"name,omitempty"`

	// Content of VCL snippet.
	// Default: ""
	Content string `form:"content,omitempty"`

	// Whether VCL snippet is dynamic or not.
	// Default: 0
	Dynamic int `form:"dynamic,omitempty"`

	// Type of VCL Snippet. This can be init or a sub-routine. Either: init, recv, hit, miss, pass, fetch, error, deliver, log, or none.
	// Default: none
	Type string `form:"type,omitempty"`

	// Priority of VCL Snippet.
	// Default: 100
	Priority int `form:"priority,omitempty"`

	// Whether to Activate
	// Default: false
	Activate bool `form:"activate,omitempty"`
}

// UpdateDynamicVCLSnippet updates a dynamic snippet for a particular service.
func (c *Client) UpdateDynamicVCLSnippet(i *UpdateDynamicVCLSnippetInput) (*VCLSnippet, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.ID == "" {
		return nil, ErrMissingName
	}

	// Update a dynamic snippet for a particular service.
	path := fmt.Sprintf("/service/%s/snippet/%s", i.Service, i.ID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var vclSnippet *VCLSnippet
	if err := decodeJSON(&vclSnippet, resp.Body); err != nil {
		return nil, err
	}
	return vclSnippet, nil
}

// UpdateVCLSnippetInput is used as input to the UpdateVCLSnippetName function.
type UpdateVCLSnippetInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int

	// Old Name of VCL snippet. Required.
	OldName string `form:"name,omitempty"`

	// New Name of VCL snippet. Required.
	NewName string `form:"name,omitempty"`
}

// UpdateVCLSnippetName updates a dynamic snippet for a particular service.
func (c *Client) UpdateVCLSnippetName(i *UpdateVCLSnippetInput) (*VCLSnippet, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.OldName == "" || i.NewName == "" {
		return nil, ErrMissingName
	}

	// Update a specific snippet for a particular service and version
	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.Service, i.Version, i.OldName)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var vclSnippet *VCLSnippet

	if err := decodeJSON(&vclSnippet, resp.Body); err != nil {
		return nil, err
	}

	return vclSnippet, nil
}

// DeleteVCLSnippetInput is the input parameter to DeleteVCLSnippet.
type DeleteVCLSnippetInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the VCL to delete (required).
	Name string
}

// DeleteVCLSnippet deletes a specific snippet for a particular service and version.
func (c *Client) DeleteVCLSnippet(i *DeleteVCLSnippetInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.Service, i.Version, i.Name)
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
