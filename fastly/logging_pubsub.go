package fastly

import (
	"context"
	"strconv"
	"time"
)

// Pubsub represents an Pubsub logging response from the Fastly API.
type Pubsub struct {
	AccountName       *string    `mapstructure:"account_name"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	Placement         *string    `mapstructure:"placement"`
	ProcessingRegion  *string    `mapstructure:"log_processing_region"`
	ProjectID         *string    `mapstructure:"project_id"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	SecretKey         *string    `mapstructure:"secret_key"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	Topic             *string    `mapstructure:"topic"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	User              *string    `mapstructure:"user"`
}

// ListPubsubsInput is used as input to the ListPubsubs function.
type ListPubsubsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListPubsubs retrieves all resources.
func (c *Client) ListPubsubs(ctx context.Context, i *ListPubsubsInput) ([]*Pubsub, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "pubsub")
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pubsubs []*Pubsub
	if err := DecodeBodyMap(resp.Body, &pubsubs); err != nil {
		return nil, err
	}
	return pubsubs, nil
}

// CreatePubsubInput is used as input to the CreatePubsub function.
type CreatePubsubInput struct {
	// AccountName is the name of the Google Cloud Platform service account associated with the target log collection service. Not required if user and secret_key are provided.
	AccountName *string `url:"account_name,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to Google Cloud Pub/Sub.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// ProjectID is your Google Cloud Platform project ID. Required.
	ProjectID *string `url:"project_id,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is your Google Cloud Platform account secret key. The private_key field in your service account authentication JSON. Not required if account_name is specified.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Topic is the Google Cloud Pub/Sub topic to which logs will be published.
	Topic *string `url:"topic,omitempty"`
	// User is your Google Cloud Platform service account email address. The client_email field in your service account authentication JSON. Not required if account_name is specified.
	User *string `url:"user,omitempty"`
}

// CreatePubsub creates a new resource.
func (c *Client) CreatePubsub(ctx context.Context, i *CreatePubsubInput) (*Pubsub, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "pubsub")
	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pubsub *Pubsub
	if err := DecodeBodyMap(resp.Body, &pubsub); err != nil {
		return nil, err
	}
	return pubsub, nil
}

// GetPubsubInput is used as input to the GetPubsub function.
type GetPubsubInput struct {
	// Name is the name of the Pubsub to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetPubsub retrieves the specified resource.
func (c *Client) GetPubsub(ctx context.Context, i *GetPubsubInput) (*Pubsub, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "pubsub", i.Name)
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Pubsub
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdatePubsubInput is used as input to the UpdatePubsub function.
type UpdatePubsubInput struct {
	// AccountName is the name of the Google Cloud Platform service account associated with the target log collection service. Not required if user and secret_key are provided.
	AccountName *string `url:"account_name,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the Pubsub to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to Google Cloud Pub/Sub.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// ProjectID is your Google Cloud Platform project ID. Required.
	ProjectID *string `url:"project_id,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is your Google Cloud Platform account secret key. The private_key field in your service account authentication JSON. Not required if account_name is specified.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Topic is the Google Cloud Pub/Sub topic to which logs will be published.
	Topic *string `url:"topic,omitempty"`
	// User is your Google Cloud Platform service account email address. The client_email field in your service account authentication JSON. Not required if account_name is specified.
	User *string `url:"user,omitempty"`
}

// UpdatePubsub updates the specified resource.
func (c *Client) UpdatePubsub(ctx context.Context, i *UpdatePubsubInput) (*Pubsub, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "pubsub", i.Name)
	resp, err := c.PutForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Pubsub
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeletePubsubInput is the input parameter to DeletePubsub.
type DeletePubsubInput struct {
	// Name is the name of the Pubsub to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeletePubsub deletes the specified resource.
func (c *Client) DeletePubsub(ctx context.Context, i *DeletePubsubInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "pubsub", i.Name)
	resp, err := c.Delete(ctx, path, CreateRequestOptions())
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
