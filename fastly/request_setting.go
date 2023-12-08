package fastly

import (
	"fmt"
	"net/url"
	"time"
)

const (
	// RequestSettingActionUnset unsets the action.
	RequestSettingActionUnset RequestSettingAction = ""

	// RequestSettingActionLookup sets request handling to lookup via the cache.
	RequestSettingActionLookup RequestSettingAction = "lookup"

	// RequestSettingActionPass sets request handling to pass the cache.
	RequestSettingActionPass RequestSettingAction = "pass"
)

// RequestSettingAction is a type of request setting action.
type RequestSettingAction string

const (
	// RequestSettingXFFClear clears any X-Forwarded-For headers.
	RequestSettingXFFClear RequestSettingXFF = "clear"

	// RequestSettingXFFLeave leaves any X-Forwarded-For headers untouched.
	RequestSettingXFFLeave RequestSettingXFF = "leave"

	// RequestSettingXFFAppend adds Fastly X-Forwarded-For headers.
	RequestSettingXFFAppend RequestSettingXFF = "append"

	// RequestSettingXFFAppendAll appends all Fastly X-Forwarded-For headers.
	RequestSettingXFFAppendAll RequestSettingXFF = "append_all"

	// RequestSettingXFFOverwrite clears any X-Forwarded-For headers and replaces
	// with Fastly ones.
	RequestSettingXFFOverwrite RequestSettingXFF = "overwrite"
)

// RequestSettingXFF is a type of X-Forwarded-For value to set.
type RequestSettingXFF string

// RequestSetting represents a request setting response from the Fastly API.
type RequestSetting struct {
	Action           *RequestSettingAction `mapstructure:"action"`
	BypassBusyWait   *bool                 `mapstructure:"bypass_busy_wait"`
	CreatedAt        *time.Time            `mapstructure:"created_at"`
	DefaultHost      *string               `mapstructure:"default_host"`
	DeletedAt        *time.Time            `mapstructure:"deleted_at"`
	ForceMiss        *bool                 `mapstructure:"force_miss"`
	ForceSSL         *bool                 `mapstructure:"force_ssl"`
	GeoHeaders       *bool                 `mapstructure:"geo_headers"`
	HashKeys         *string               `mapstructure:"hash_keys"`
	MaxStaleAge      *int                  `mapstructure:"max_stale_age"`
	Name             *string               `mapstructure:"name"`
	RequestCondition *string               `mapstructure:"request_condition"`
	ServiceID        *string               `mapstructure:"service_id"`
	ServiceVersion   *int                  `mapstructure:"version"`
	TimerSupport     *bool                 `mapstructure:"timer_support"`
	UpdatedAt        *time.Time            `mapstructure:"updated_at"`
	XForwardedFor    *RequestSettingXFF    `mapstructure:"xff"`
}

// ListRequestSettingsInput is used as input to the ListRequestSettings
// function.
type ListRequestSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListRequestSettings retrieves all resources.
func (c *Client) ListRequestSettings(i *ListRequestSettingsInput) ([]*RequestSetting, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bs []*RequestSetting
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	return bs, nil
}

// CreateRequestSettingInput is used as input to the CreateRequestSetting
// function.
type CreateRequestSettingInput struct {
	// Action allows you to terminate request handling and immediately perform an action.
	Action *RequestSettingAction `url:"action,omitempty"`
	// BypassBusyWait disables collapsed forwarding, so you don't wait for other objects to origin.
	BypassBusyWait *Compatibool `url:"bypass_busy_wait,omitempty"`
	// DefaultHost sets the host header.
	DefaultHost *string `url:"default_host,omitempty"`
	// ForceMiss allows you to force a cache miss for the request. Replaces the item in the cache if the content is cacheable.
	ForceMiss *Compatibool `url:"force_miss,omitempty"`
	// ForceSSL forces the request use SSL (redirects a non-SSL to SSL).
	ForceSSL *Compatibool `url:"force_ssl,omitempty"`
	// GeoHeaders injects Fastly-Geo-Country, Fastly-Geo-City, and Fastly-Geo-Region into the request headers.
	GeoHeaders *Compatibool `url:"geo_headers,omitempty"`
	// HashKeys is a comma separated list of varnish request object fields that should be in the hash key.
	HashKeys *string `url:"hash_keys,omitempty"`
	// MaxStaleAge is how old an object is allowed to be to serve stale-if-error or stale-while-revalidate.
	MaxStaleAge *int `url:"max_stale_age,omitempty"`
	// Name is the name for the request settings.
	Name *string `url:"name,omitempty"`
	// RequestCondition is the condition which, if met, will select this configuration during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimerSupport injects the X-Timer info into the request for viewing origin fetch durations.
	TimerSupport *Compatibool `url:"timer_support,omitempty"`
	// XForwardedFor determines header value (clear, leave, append, append_all, overwrite)
	XForwardedFor *RequestSettingXFF `url:"xff,omitempty"`
}

// CreateRequestSetting creates a new resource.
func (c *Client) CreateRequestSetting(i *CreateRequestSettingInput) (*RequestSetting, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *RequestSetting
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetRequestSettingInput is used as input to the GetRequestSetting function.
type GetRequestSettingInput struct {
	// Name is the name of the request settings to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetRequestSetting retrieves the specified resource.
func (c *Client) GetRequestSetting(i *GetRequestSettingInput) (*RequestSetting, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *RequestSetting
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateRequestSettingInput is used as input to the UpdateRequestSetting
// function.
type UpdateRequestSettingInput struct {
	// Action allows you to terminate request handling and immediately perform an action.
	Action *RequestSettingAction `url:"action,omitempty"`
	// BypassBusyWait disables collapsed forwarding, so you don't wait for other objects to origin.
	BypassBusyWait *Compatibool `url:"bypass_busy_wait,omitempty"`
	// DefaultHost sets the host header.
	DefaultHost *string `url:"default_host,omitempty"`
	// ForceMiss allows you to force a cache miss for the request. Replaces the item in the cache if the content is cacheable.
	ForceMiss *Compatibool `url:"force_miss,omitempty"`
	// ForceSSL forces the request use SSL (redirects a non-SSL to SSL).
	ForceSSL *Compatibool `url:"force_ssl,omitempty"`
	// GeoHeaders injects Fastly-Geo-Country, Fastly-Geo-City, and Fastly-Geo-Region into the request headers.
	GeoHeaders *Compatibool `url:"geo_headers,omitempty"`
	// HashKeys is a comma separated list of varnish request object fields that should be in the hash key.
	HashKeys *string `url:"hash_keys,omitempty"`
	// MaxStaleAge is how old an object is allowed to be to serve stale-if-error or stale-while-revalidate.
	MaxStaleAge *int `url:"max_stale_age,omitempty"`
	// Name is the name of the request settings to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// RequestCondition is the condition which, if met, will select this configuration during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimerSupport injects the X-Timer info into the request for viewing origin fetch durations.
	TimerSupport *Compatibool `url:"timer_support,omitempty"`
	// XForwardedFor determines header value (clear, leave, append, append_all, overwrite)
	XForwardedFor *RequestSettingXFF `url:"xff,omitempty"`
}

// UpdateRequestSetting updates the specified resource.
func (c *Client) UpdateRequestSetting(i *UpdateRequestSettingInput) (*RequestSetting, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *RequestSetting
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteRequestSettingInput is the input parameter to DeleteRequestSetting.
type DeleteRequestSettingInput struct {
	// Name is the name of the request settings to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteRequestSetting deletes the specified resource.
func (c *Client) DeleteRequestSetting(i *DeleteRequestSettingInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
