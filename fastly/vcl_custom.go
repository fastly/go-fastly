package fastly

import (
	"strconv"
	"time"
)

// VCL represents a response about VCL from the Fastly API.
type VCL struct {
	Content        *string    `mapstructure:"content"`
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	Main           *bool      `mapstructure:"main"`
	Name           *string    `mapstructure:"name"`
	ServiceID      *string    `mapstructure:"service_id"`
	ServiceVersion *int       `mapstructure:"version"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}

// ListVCLsInput is used as input to the ListVCLs function.
type ListVCLsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListVCLs retrieves all resources.
func (c *Client) ListVCLs(i *ListVCLsInput) ([]*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "vcl")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vcls []*VCL
	if err := decodeBodyMap(resp.Body, &vcls); err != nil {
		return nil, err
	}
	return vcls, nil
}

// GetVCLInput is used as input to the GetVCL function.
type GetVCLInput struct {
	// Name is the name of the VCL to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetVCL retrieves the specified resource.
func (c *Client) GetVCL(i *GetVCLInput) (*VCL, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "vcl", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

// GetGeneratedVCL retrieves the specified resource.
func (c *Client) GetGeneratedVCL(i *GetGeneratedVCLInput) (*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "generated_vcl")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// CreateVCLInput is used as input to the CreateVCL function.
type CreateVCLInput struct {
	// Content is the VCL code to be included.
	Content *string `url:"content,omitempty"`
	// Main is set to true when this is the main VCL, otherwise false.
	Main *bool `url:"main,omitempty"`
	// Name is the name of this VCL.
	Name *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// CreateVCL creates a new resource.
func (c *Client) CreateVCL(i *CreateVCLInput) (*VCL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "vcl")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// UpdateVCLInput is used as input to the UpdateVCL function.
type UpdateVCLInput struct {
	// Content is the VCL code to be included.
	Content *string `url:"content,omitempty"`
	// Name is the name of the VCL to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// UpdateVCL updates the specified resource.
func (c *Client) UpdateVCL(i *UpdateVCLInput) (*VCL, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "vcl", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// ActivateVCLInput is used as input to the ActivateVCL function.
type ActivateVCLInput struct {
	// Name is the name of the VCL to mark as main (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ActivateVCL creates a new Fastly VCL.
func (c *Client) ActivateVCL(i *ActivateVCLInput) (*VCL, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "vcl", i.Name, "main")
	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vcl *VCL
	if err := decodeBodyMap(resp.Body, &vcl); err != nil {
		return nil, err
	}
	return vcl, nil
}

// DeleteVCLInput is the input parameter to DeleteVCL.
type DeleteVCLInput struct {
	// Name is the name of the VCL to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteVCL deletes the specified resource.
func (c *Client) DeleteVCL(i *DeleteVCLInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "vcl", i.Name)
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
