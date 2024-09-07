package fastly

import (
	"fmt"
	"net/url"
	"time"
)

// Honeycomb represents a honeycomb response from the Fastly API.
type Honeycomb struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	Dataset           *string    `mapstructure:"dataset"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	Placement         *string    `mapstructure:"placement"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	Token             *string    `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// ListHoneycombsInput is used as input to the ListHoneycombs function.
type ListHoneycombsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListHoneycombs retrieves all resources.
func (c *Client) ListHoneycombs(i *ListHoneycombsInput) ([]*Honeycomb, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var hs []*Honeycomb
	if err := decodeBodyMap(resp.Body, &hs); err != nil {
		return nil, err
	}
	return hs, nil
}

// CreateHoneycombInput is used as input to the CreateHoneycomb function.
type CreateHoneycombInput struct {
	// Dataset is the Honeycomb Dataset you want to log to.
	Dataset *string `url:"dataset,omitempty"`
	// Format is a Fastly log format string. Must produce valid JSON that Honeycomb can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Token is the Write Key from the Account page of your Honeycomb account.
	Token *string `url:"token,omitempty"`
}

// CreateHoneycomb creates a new resource.
func (c *Client) CreateHoneycomb(i *CreateHoneycombInput) (*Honeycomb, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetHoneycombInput is used as input to the GetHoneycomb function.
type GetHoneycombInput struct {
	// Name is the name for the real-time logging configuration (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetHoneycomb retrieves the specified resource.
func (c *Client) GetHoneycomb(i *GetHoneycombInput) (*Honeycomb, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// UpdateHoneycombInput is used as input to the UpdateHoneycomb function.
type UpdateHoneycombInput struct {
	// Dataset is the Honeycomb Dataset you want to log to.
	Dataset *string `url:"dataset,omitempty"`
	// Format is a Fastly log format string. Must produce valid JSON that Honeycomb can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the honeycomb to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Token is the Write Key from the Account page of your Honeycomb account.
	Token *string `url:"token,omitempty"`
}

// UpdateHoneycomb updates the specified resource.
func (c *Client) UpdateHoneycomb(i *UpdateHoneycombInput) (*Honeycomb, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// DeleteHoneycombInput is the input parameter to DeleteHoneycomb.
type DeleteHoneycombInput struct {
	// Name is the name of the honeycomb to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteHoneycomb deletes the specified resource.
func (c *Client) DeleteHoneycomb(i *DeleteHoneycombInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/honeycomb/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
		return ErrStatusNotOk
	}
	return nil
}
