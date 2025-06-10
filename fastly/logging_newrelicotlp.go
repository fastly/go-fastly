package fastly

import (
	"context"
	"strconv"
	"time"
)

// NewRelicOTLP represents a newrelic response from the Fastly API.
type NewRelicOTLP struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	Placement         *string    `mapstructure:"placement"`
	ProcessingRegion  *string    `mapstructure:"log_processing_region"`
	Region            *string    `mapstructure:"region"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	Token             *string    `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	URL               *string    `mapstructure:"url"`
}

// ListNewRelicOTLPInput is used as input to the ListNewRelicOTLP function.
type ListNewRelicOTLPInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListNewRelicOTLP returns the list of newrelic for the configuration version.
func (c *Client) ListNewRelicOTLP(i *ListNewRelicOTLPInput) ([]*NewRelicOTLP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "newrelicotlp")
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var n []*NewRelicOTLP
	if err := DecodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// CreateNewRelicOTLPInput is used as input to the CreateNewRelicOTLP function.
type CreateNewRelicOTLPInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
	// Format is a Fastly log format string. Must produce valid JSON that New Relic Logs can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to New Relic.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// Region is the region where logs are received and stored by New Relic.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Token is the Insert API key from the Account page of your New Relic account.
	Token *string `url:"token,omitempty"`
	// URL is the optional URL of a New Relic trace observer to send logs
	// to. Must be a New Relic domain name.
	URL *string `url:"url,omitempty"`
}

// CreateNewRelicOTLP creates a new Fastly newrelic.
func (c *Client) CreateNewRelicOTLP(i *CreateNewRelicOTLPInput) (*NewRelicOTLP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "newrelicotlp")
	resp, err := c.PostForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var n *NewRelicOTLP
	if err := DecodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// GetNewRelicOTLPInput is used as input to the GetNewRelicOTLP function.
type GetNewRelicOTLPInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Name is the name of the newrelic to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetNewRelicOTLP gets the newrelic configuration with the given parameters.
func (c *Client) GetNewRelicOTLP(i *GetNewRelicOTLPInput) (*NewRelicOTLP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}
	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "newrelicotlp", i.Name)
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var n *NewRelicOTLP
	if err := DecodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// UpdateNewRelicOTLPInput is used as input to the UpdateNewRelicOTLP function.
type UpdateNewRelicOTLPInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
	// Format is a Fastly log format string. Must produce valid JSON that New Relic Logs can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the newrelic to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the Fastly region where logs will be processed before streaming to the endpoint.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// Region is the region to which to stream logs.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Token is the Insert API key from the Account page of your New Relic account.
	Token *string `url:"token,omitempty"`
	// URL is the optional URL of a New Relic trace observer to send logs
	// to. Must be a New Relic domain name.
	URL *string `url:"url,omitempty"`
}

// UpdateNewRelicOTLP updates a specific newrelic.
func (c *Client) UpdateNewRelicOTLP(i *UpdateNewRelicOTLPInput) (*NewRelicOTLP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}
	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "newrelicotlp", i.Name)
	resp, err := c.PutForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var n *NewRelicOTLP
	if err := DecodeBodyMap(resp.Body, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// DeleteNewRelicOTLPInput is the input parameter to DeleteNewRelicOTLP.
type DeleteNewRelicOTLPInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Name is the name of the newrelicotlp to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteNewRelicOTLP deletes the given newrelic version.
func (c *Client) DeleteNewRelicOTLP(i *DeleteNewRelicOTLPInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}
	if i.Name == "" {
		return ErrMissingName
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "newrelicotlp", i.Name)
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
