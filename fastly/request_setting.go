package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

const (
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
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name             string               `mapstructure:"name"`
	ForceMiss        bool                 `mapstructure:"force_miss"`
	ForceSSL         bool                 `mapstructure:"force_ssl"`
	Action           RequestSettingAction `mapstructure:"action"`
	BypassBusyWait   bool                 `mapstructure:"bypass_busy_wait"`
	MaxStaleAge      uint                 `mapstructure:"max_stale_age"`
	HashKeys         string               `mapstructure:"hash_keys"`
	XForwardedFor    RequestSettingXFF    `mapstructure:"xff"`
	TimerSupport     bool                 `mapstructure:"timer_support"`
	GeoHeaders       bool                 `mapstructure:"geo_headers"`
	DefaultHost      string               `mapstructure:"default_host"`
	RequestCondition string               `mapstructure:"request_condition"`
	CreatedAt        *time.Time           `mapstructure:"created_at"`
	UpdatedAt        *time.Time           `mapstructure:"updated_at"`
	DeletedAt        *time.Time           `mapstructure:"deleted_at"`
}

// requestSettingsByName is a sortable list of request settings.
type requestSettingsByName []*RequestSetting

// Len, Swap, and Less implement the sortable interface.
func (s requestSettingsByName) Len() int      { return len(s) }
func (s requestSettingsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s requestSettingsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListRequestSettingsInput is used as input to the ListRequestSettings
// function.
type ListRequestSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListRequestSettings returns the list of request settings for the
// configuration version.
func (c *Client) ListRequestSettings(i *ListRequestSettingsInput) ([]*RequestSetting, error) {
	if i.ServiceID == "" {
		return nil, NewFieldError("ServiceID")
	}

	if i.ServiceVersion == 0 {
		return nil, NewFieldError("ServiceVersion")
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var bs []*RequestSetting
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	sort.Stable(requestSettingsByName(bs))
	return bs, nil
}

// CreateRequestSettingInput is used as input to the CreateRequestSetting
// function.
type CreateRequestSettingInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name             string               `form:"name,omitempty"`
	ForceMiss        Compatibool          `form:"force_miss,omitempty"`
	ForceSSL         Compatibool          `form:"force_ssl,omitempty"`
	Action           RequestSettingAction `form:"action,omitempty"`
	BypassBusyWait   Compatibool          `form:"bypass_busy_wait,omitempty"`
	MaxStaleAge      uint                 `form:"max_stale_age,omitempty"`
	HashKeys         string               `form:"hash_keys,omitempty"`
	XForwardedFor    RequestSettingXFF    `form:"xff,omitempty"`
	TimerSupport     Compatibool          `form:"timer_support,omitempty"`
	GeoHeaders       Compatibool          `form:"geo_headers,omitempty"`
	DefaultHost      string               `form:"default_host,omitempty"`
	RequestCondition string               `form:"request_condition,omitempty"`
}

// CreateRequestSetting creates a new Fastly request settings.
func (c *Client) CreateRequestSetting(i *CreateRequestSettingInput) (*RequestSetting, error) {
	if i.ServiceID == "" {
		return nil, NewFieldError("ServiceID")
	}

	if i.ServiceVersion == 0 {
		return nil, NewFieldError("ServiceVersion")
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *RequestSetting
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetRequestSettingInput is used as input to the GetRequestSetting function.
type GetRequestSettingInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the request settings to fetch.
	Name string
}

// GetRequestSetting gets the request settings configuration with the given
// parameters.
func (c *Client) GetRequestSetting(i *GetRequestSettingInput) (*RequestSetting, error) {
	if i.ServiceID == "" {
		return nil, NewFieldError("ServiceID")
	}

	if i.ServiceVersion == 0 {
		return nil, NewFieldError("ServiceVersion")
	}

	if i.Name == "" {
		return nil, NewFieldError("Name")
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var b *RequestSetting
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateRequestSettingInput is used as input to the UpdateRequestSetting
// function.
type UpdateRequestSettingInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the request settings to update.
	Name string

	NewName          *string              `form:"name,omitempty"`
	ForceMiss        *Compatibool         `form:"force_miss,omitempty"`
	ForceSSL         *Compatibool         `form:"force_ssl,omitempty"`
	Action           RequestSettingAction `form:"action,omitempty"`
	BypassBusyWait   *Compatibool         `form:"bypass_busy_wait,omitempty"`
	MaxStaleAge      *uint                `form:"max_stale_age,omitempty"`
	HashKeys         *string              `form:"hash_keys,omitempty"`
	XForwardedFor    RequestSettingXFF    `form:"xff,omitempty"`
	TimerSupport     *Compatibool         `form:"timer_support,omitempty"`
	GeoHeaders       *Compatibool         `form:"geo_headers,omitempty"`
	DefaultHost      *string              `form:"default_host,omitempty"`
	RequestCondition *string              `form:"request_condition,omitempty"`
}

// UpdateRequestSetting updates a specific request settings.
func (c *Client) UpdateRequestSetting(i *UpdateRequestSettingInput) (*RequestSetting, error) {
	if i.ServiceID == "" {
		return nil, NewFieldError("ServiceID")
	}

	if i.ServiceVersion == 0 {
		return nil, NewFieldError("ServiceVersion")
	}

	if i.Name == "" {
		return nil, NewFieldError("Name")
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *RequestSetting
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteRequestSettingInput is the input parameter to DeleteRequestSetting.
type DeleteRequestSettingInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the request settings to delete (required).
	Name string
}

// DeleteRequestSetting deletes the given request settings version.
func (c *Client) DeleteRequestSetting(i *DeleteRequestSettingInput) error {
	if i.ServiceID == "" {
		return NewFieldError("ServiceID")
	}

	if i.ServiceVersion == 0 {
		return NewFieldError("ServiceVersion")
	}

	if i.Name == "" {
		return NewFieldError("Name")
	}

	path := fmt.Sprintf("/service/%s/version/%d/request_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
