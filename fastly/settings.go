package fastly

import "fmt"

// Settings represents a backend response from the Fastly API.
type Settings struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	DefaultTTL      uint   `mapstructure:"general.default_ttl"`
	DefaultHost     string `mapstructure:"general.default_host"`
	StaleIfError    bool   `mapstructure:"general.stale_if_error"`
	StaleIfErrorTTL uint   `mapstructure:"general.stale_if_error_ttl"`
}

// GetSettingsInput is used as input to the GetSettings function.
type GetSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetSettings gets the backend configuration with the given parameters.
func (c *Client) GetSettings(i *GetSettingsInput) (*Settings, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/settings", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var b *Settings
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateSettingsInput is used as input to the UpdateSettings function.
type UpdateSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	DefaultTTL      uint    `form:"general.default_ttl"`
	DefaultHost     *string `form:"general.default_host,omitempty"`
	StaleIfError    *bool   `form:"general.stale_if_error,omitempty"`
	StaleIfErrorTTL *uint   `form:"general.stale_if_error_ttl,omitempty"`
}

// UpdateSettings updates a specific backend.
func (c *Client) UpdateSettings(i *UpdateSettingsInput) (*Settings, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/settings", i.ServiceID, i.ServiceVersion)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *Settings
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}
