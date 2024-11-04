package fastly

import (
	"strconv"
	"time"
)

// GrafanaCloudLogs represents a GrafanaCloudLogs response from the Fastly API.
type GrafanaCloudLogs struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	MessageType       *string    `mapstructure:"message_type"`
	Placement         *string    `mapstructure:"placement"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`

	Index *string `mapstructure:"index"`
	Token *string `mapstructure:"token"`
	URL   *string `mapstructure:"url"`
	User  *string `mapstructure:"user"`
}

// ListGrafanaCloudLogsInput is used as input to the ListGrafanaCloudLogs function.
type ListGrafanaCloudLogsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListGrafanaCloudLogs retrieves all resources.
func (c *Client) ListGrafanaCloudLogs(i *ListGrafanaCloudLogsInput) ([]*GrafanaCloudLogs, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "grafanacloudlogs")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d []*GrafanaCloudLogs
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// CreateGrafanaCloudLogsInput is used as input to the CreateGrafanaCloudLogs function.
type CreateGrafanaCloudLogsInput struct {
	// Format is a Fastly log format string. Must produce valid JSON that GrafanaCloudLogs can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Index
	Index *string `url:"index,omitempty"`
	// Token is the API key from your GrafanaCloudLogs account.
	Token *string `url:"token,omitempty"`
	// Grafana User ID
	User *string `url:"user,omitempty"`
	// URL is the URL to stream logs to. Must use HTTPS.
	URL *string `url:"url,omitempty"`
}

// CreateGrafanaCloudLogs creates a new resource.
func (c *Client) CreateGrafanaCloudLogs(i *CreateGrafanaCloudLogsInput) (*GrafanaCloudLogs, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "grafanacloudlogs")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *GrafanaCloudLogs
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetGrafanaCloudLogsInput is used as input to the GetGrafanaCloudLogs function.
type GetGrafanaCloudLogsInput struct {
	// Name is the name of the GrafanaCloudLogs to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetGrafanaCloudLogs retrieves the specified resource.
func (c *Client) GetGrafanaCloudLogs(i *GetGrafanaCloudLogsInput) (*GrafanaCloudLogs, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "grafanacloudlogs", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *GrafanaCloudLogs
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// UpdateGrafanaCloudLogsInput is used as input to the UpdateGrafanaCloudLogs function.
type UpdateGrafanaCloudLogsInput struct {
	// Format is a Fastly log format string. Must produce valid JSON that GrafanaCloudLogs can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the GrafanaCloudLogs to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Grafana User ID
	User *string `url:"user,omitempty"`
	// Token is the API key from your GrafanaCloudLogs account.
	Token *string `url:"token,omitempty"`
	// URL is the URL to stream logs to. Must use HTTPS.
	URL *string `url:"url,omitempty"`
	// Index is the stream identifier
	Index *string `url:"index,omitempty"`
}

// UpdateGrafanaCloudLogs updates the specified resource.
func (c *Client) UpdateGrafanaCloudLogs(i *UpdateGrafanaCloudLogsInput) (*GrafanaCloudLogs, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "grafanacloudlogs", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *GrafanaCloudLogs
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteGrafanaCloudLogsInput is the input parameter to DeleteGrafanaCloudLogs.
type DeleteGrafanaCloudLogsInput struct {
	// Name is the name of the GrafanaCloudLogs to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteGrafanaCloudLogs deletes the specified resource.
func (c *Client) DeleteGrafanaCloudLogs(i *DeleteGrafanaCloudLogsInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "grafanacloudlogs", i.Name)
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
