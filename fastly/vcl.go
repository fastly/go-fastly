package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// VCL represents a response about VCL from the Fastly API.
type VCL struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name      string     `mapstructure:"name"`
	Main      bool       `mapstructure:"main"`
	Content   string     `mapstructure:"content"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
}

// vclsByName is a sortable list of VCLs.
type vclsByName []*VCL

// Len, Swap, and Less implement the sortable interface.
func (s vclsByName) Len() int      { return len(s) }
func (s vclsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s vclsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListVCLsInput is used as input to the ListVCLs function.
type ListVCLsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListVCLs returns the list of VCLs for the configuration version.
func (c *Client) ListVCLs(i *ListVCLsInput) ([]*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/vcl", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var vcls []*VCL
	if err := decodeBodyMap(resp.Body, &vcls); err != nil {
		return nil, err
	}
	sort.Stable(vclsByName(vcls))
	return vcls, nil
}

// GetVCLInput is used as input to the GetVCL function.
type GetVCLInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the VCL to fetch.
	Name string
}

// GetVCL gets the VCL configuration with the given parameters.
func (c *Client) GetVCL(i *GetVCLInput) (*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/vcl/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// GetGeneratedVCLInput is used as input to the GetGeneratedVCL function.
type GetGeneratedVCLInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetGeneratedVCL gets the VCL configuration with the given parameters.
func (c *Client) GetGeneratedVCL(i *GetGeneratedVCLInput) (*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/generated_vcl", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// CreateVCLInput is used as input to the CreateVCL function.
type CreateVCLInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name    string `form:"name,omitempty"`
	Content string `form:"content,omitempty"`
	Main    bool   `form:"main,omitempty"`
}

// CreateVCL creates a new Fastly VCL.
func (c *Client) CreateVCL(i *CreateVCLInput) (*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/vcl", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// UpdateVCLInput is used as input to the UpdateVCL function.
type UpdateVCLInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the VCL to update (required).
	Name string

	NewName *string `form:"name,omitempty"`
	Content *string `form:"content,omitempty"`
}

// UpdateVCL creates a new Fastly VCL.
func (c *Client) UpdateVCL(i *UpdateVCLInput) (*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/vcl/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// ActivateVCLInput is used as input to the ActivateVCL function.
type ActivateVCLInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the VCL to mark as main (required).
	Name string
}

// ActivateVCL creates a new Fastly VCL.
func (c *Client) ActivateVCL(i *ActivateVCLInput) (*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/vcl/%s/main", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// DeleteVCLInput is the input parameter to DeleteVCL.
type DeleteVCLInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the VCL to delete (required).
	Name string
}

// DeleteVCL deletes the given VCL version.
func (c *Client) DeleteVCL(i *DeleteVCLInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/vcl/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
