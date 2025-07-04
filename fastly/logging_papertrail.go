package fastly

import (
	"context"
	"strconv"
	"time"
)

// Papertrail represents a papertrail response from the Fastly API.
type Papertrail struct {
	Address           *string    `mapstructure:"address"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	Placement         *string    `mapstructure:"placement"`
	Port              *int       `mapstructure:"port"`
	ProcessingRegion  *string    `mapstructure:"log_processing_region"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// ListPapertrailsInput is used as input to the ListPapertrails function.
type ListPapertrailsInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListPapertrails retrieves all resources.
func (c *Client) ListPapertrails(i *ListPapertrailsInput) ([]*Papertrail, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "papertrail")
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ps []*Papertrail
	if err := DecodeBodyMap(resp.Body, &ps); err != nil {
		return nil, err
	}
	return ps, nil
}

// CreatePapertrailInput is used as input to the CreatePapertrail function.
type CreatePapertrailInput struct {
	// Address is a hostname or IPv4 address.
	Address *string `url:"address,omitempty"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Port is the port number.
	Port *int `url:"port,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to Papertrail.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// CreatePapertrail creates a new resource.
func (c *Client) CreatePapertrail(i *CreatePapertrailInput) (*Papertrail, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "papertrail")
	resp, err := c.PostForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p *Papertrail
	if err := DecodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// GetPapertrailInput is used as input to the GetPapertrail function.
type GetPapertrailInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Name is the name of the papertrail to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetPapertrail retrieves the specified resource.
func (c *Client) GetPapertrail(i *GetPapertrailInput) (*Papertrail, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "papertrail", i.Name)
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p *Papertrail
	if err := DecodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// UpdatePapertrailInput is used as input to the UpdatePapertrail function.
type UpdatePapertrailInput struct {
	// Address is a hostname or IPv4 address.
	Address *string `url:"address,omitempty"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the papertrail to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Port is the port number.
	Port *int `url:"port,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to Papertrail.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// UpdatePapertrail updates the specified resource.
func (c *Client) UpdatePapertrail(i *UpdatePapertrailInput) (*Papertrail, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "papertrail", i.Name)
	resp, err := c.PutForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p *Papertrail
	if err := DecodeBodyMap(resp.Body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

// DeletePapertrailInput is the input parameter to DeletePapertrail.
type DeletePapertrailInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Name is the name of the papertrail to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeletePapertrail deletes the specified resource.
func (c *Client) DeletePapertrail(i *DeletePapertrailInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "papertrail", i.Name)
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
		return ErrNotOK
	}
	return nil
}
