package fastly

import (
	"context"
	"strconv"
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
	ProcessingRegion  *string    `mapstructure:"log_processing_region"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	Token             *string    `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// ListHoneycombsInput is used as input to the ListHoneycombs function.
type ListHoneycombsInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
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

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "honeycomb")
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var hs []*Honeycomb
	if err := DecodeBodyMap(resp.Body, &hs); err != nil {
		return nil, err
	}
	return hs, nil
}

// CreateHoneycombInput is used as input to the CreateHoneycomb function.
type CreateHoneycombInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
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
	// ProcessingRegion is the Fastly region where logs will be processed before streaming to the endpoint.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
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

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "honeycomb")
	resp, err := c.PostForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetHoneycombInput is used as input to the GetHoneycomb function.
type GetHoneycombInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
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

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "honeycomb", i.Name)
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// UpdateHoneycombInput is used as input to the UpdateHoneycomb function.
type UpdateHoneycombInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
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
	// ProcessingRegion is the Fastly region where logs will be processed before streaming to the endpoint.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
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

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "honeycomb", i.Name)
	resp, err := c.PutForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *Honeycomb
	if err := DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// DeleteHoneycombInput is the input parameter to DeleteHoneycomb.
type DeleteHoneycombInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
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

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "honeycomb", i.Name)
	resp, err := c.Delete(path, CreateRequestOptions(i.Context))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
