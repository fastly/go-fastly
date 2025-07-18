package fastly

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// BigQuery represents a BigQuery response from the Fastly API.
type BigQuery struct {
	AccountName       *string    `mapstructure:"account_name"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	Dataset           *string    `mapstructure:"dataset"`
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
	Table             *string    `mapstructure:"table"`
	Template          *string    `mapstructure:"template_suffix"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	User              *string    `mapstructure:"user"`
}

// ListBigQueriesInput is used as input to the ListBigQueries function.
type ListBigQueriesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListBigQueries retrieves all resources.
func (c *Client) ListBigQueries(ctx context.Context, i *ListBigQueriesInput) ([]*BigQuery, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "bigquery")
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bigQueries []*BigQuery
	if err := DecodeBodyMap(resp.Body, &bigQueries); err != nil {
		return nil, err
	}
	return bigQueries, nil
}

// CreateBigQueryInput is used as input to the CreateBigQuery function.
type CreateBigQueryInput struct {
	// AccountName is the name of the Google Cloud Platform service account associated with the target log collection service.
	AccountName *string `url:"account_name,omitempty"`
	// Dataset is your BigQuery dataset.
	Dataset *string `url:"dataset,omitempty"`
	// Format is a Fastly log format string. Must produce JSON that matches the schema of your BigQuery table.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the BigQuery logging object. Used as a primary key for API access.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to BigQuery.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// ProjectID is your Google Cloud Platform project ID.
	ProjectID *string `url:"project_id,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is your Google Cloud Platform account secret key. The private_key field in your service account authentication JSON. Not required if account_name is specified.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Table is your BigQuery table.
	Table *string `url:"table,omitempty"`
	// Template is a BigQuery table name suffix template.
	Template *string `url:"template_suffix,omitempty"`
	// User is your Google Cloud Platform service account email address. The client_email field in your service account authentication JSON. Not required if account_name is specified.
	User *string `url:"user,omitempty"`
}

// CreateBigQuery creates a new resource.
func (c *Client) CreateBigQuery(ctx context.Context, i *CreateBigQueryInput) (*BigQuery, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "bigquery")
	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bigQuery *BigQuery
	if err := DecodeBodyMap(resp.Body, &bigQuery); err != nil {
		return nil, err
	}
	return bigQuery, nil
}

// GetBigQueryInput is used as input to the GetBigQuery function.
type GetBigQueryInput struct {
	// Name is the name of the BigQuery to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetBigQuery retrieves the specified resource.
func (c *Client) GetBigQuery(ctx context.Context, i *GetBigQueryInput) (*BigQuery, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "bigquery", i.Name)
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bigQuery *BigQuery
	if err := DecodeBodyMap(resp.Body, &bigQuery); err != nil {
		return nil, err
	}
	return bigQuery, nil
}

// UpdateBigQueryInput is used as input to the UpdateBigQuery function.
type UpdateBigQueryInput struct {
	// AccountName is the name of the Google Cloud Platform service account associated with the target log collection service.
	AccountName *string `url:"account_name,omitempty"`
	// Dataset is your BigQuery dataset.
	Dataset *string `url:"dataset,omitempty"`
	// Format is a Fastly log format string. Must produce JSON that matches the schema of your BigQuery table.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the BigQuery to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to BigQuery.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// ProjectID is your Google Cloud Platform project ID.
	ProjectID *string `url:"project_id,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is your Google Cloud Platform account secret key. The private_key field in your service account authentication JSON. Not required if account_name is specified.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Table is your BigQuery table.
	Table *string `url:"table,omitempty"`
	// Template is a BigQuery table name suffix template.
	Template *string `url:"template_suffix,omitempty"`
	// User is your Google Cloud Platform service account email address. The client_email field in your service account authentication JSON. Not required if account_name is specified.
	User *string `url:"user,omitempty"`
}

// UpdateBigQuery updates the specified resource.
func (c *Client) UpdateBigQuery(ctx context.Context, i *UpdateBigQueryInput) (*BigQuery, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "bigquery", i.Name)
	resp, err := c.PutForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bigQuery *BigQuery
	if err := DecodeBodyMap(resp.Body, &bigQuery); err != nil {
		return nil, err
	}
	return bigQuery, nil
}

// DeleteBigQueryInput is the input parameter to DeleteBigQuery.
type DeleteBigQueryInput struct {
	// Name is the name of the BigQuery to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteBigQuery deletes the specified resource.
func (c *Client) DeleteBigQuery(ctx context.Context, i *DeleteBigQueryInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "bigquery", i.Name)
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
		return fmt.Errorf("not ok")
	}
	return nil
}
