package fastly

import (
	"context"
	"strconv"
)

// Settings represents a backend response from the Fastly API.
type Settings struct {
	DefaultHost     *string `mapstructure:"general.default_host"`
	DefaultTTL      *uint   `mapstructure:"general.default_ttl"`
	ServiceID       *string `mapstructure:"service_id"`
	ServiceVersion  *int    `mapstructure:"version"`
	StaleIfError    *bool   `mapstructure:"general.stale_if_error"`
	StaleIfErrorTTL *uint   `mapstructure:"general.stale_if_error_ttl"`
}

// GetSettingsInput is used as input to the GetSettings function.
type GetSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetSettings retrieves the specified resource.
func (c *Client) GetSettings(ctx context.Context, i *GetSettingsInput) (*Settings, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "settings")
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Settings
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateSettingsInput is used as input to the UpdateSettings function.
type UpdateSettingsInput struct {
	// DefaultHost is the default host name for the version.
	DefaultHost *string `url:"general.default_host,omitempty"`
	// DefaultTTL is the default time-to-live (TTL) for the version.
	DefaultTTL *uint `url:"general.default_ttl"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// StaleIfError enables serving a stale object if there is an error.
	StaleIfError *bool `url:"general.stale_if_error,omitempty"`
	// StaleIfErrorTTL is the default time-to-live (TTL) for serving the stale object for the version.
	StaleIfErrorTTL *uint `url:"general.stale_if_error_ttl,omitempty"`
}

// UpdateSettings updates the specified resource.
func (c *Client) UpdateSettings(ctx context.Context, i *UpdateSettingsInput) (*Settings, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "settings")
	resp, err := c.PutForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Settings
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}
