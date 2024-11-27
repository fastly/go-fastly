package fastly

import (
	"strconv"
	"time"
)

const (
	// CacheSettingActionCache sets the cache to cache.
	CacheSettingActionCache CacheSettingAction = "cache"

	// CacheSettingActionPass sets the cache to pass through.
	CacheSettingActionPass CacheSettingAction = "pass"

	// CacheSettingActionRestart sets the cache to restart the request.
	CacheSettingActionRestart CacheSettingAction = "restart"
)

// CacheSettingAction is the type of cache action.
type CacheSettingAction string

// CacheSetting represents a response from Fastly's API for cache settings.
type CacheSetting struct {
	Action         *CacheSettingAction `mapstructure:"action"`
	CacheCondition *string             `mapstructure:"cache_condition"`
	CreatedAt      *time.Time          `mapstructure:"created_at"`
	DeletedAt      *time.Time          `mapstructure:"deleted_at"`
	Name           *string             `mapstructure:"name"`
	ServiceID      *string             `mapstructure:"service_id"`
	ServiceVersion *int                `mapstructure:"version"`
	StaleTTL       *int                `mapstructure:"stale_ttl"`
	TTL            *int                `mapstructure:"ttl"`
	UpdatedAt      *time.Time          `mapstructure:"updated_at"`
}

// ListCacheSettingsInput is used as input to the ListCacheSettings function.
type ListCacheSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListCacheSettings retrieves all resources.
func (c *Client) ListCacheSettings(i *ListCacheSettingsInput) ([]*CacheSetting, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "cache_settings")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs []*CacheSetting
	if err := DecodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// CreateCacheSettingInput is used as input to the CreateCacheSetting function.
type CreateCacheSettingInput struct {
	// Action determines vcl_fetch behaviour (pass, cache, restart).
	Action *CacheSettingAction `url:"action,omitempty"`
	// CacheCondition is name of the cache condition controlling when this configuration applies.
	CacheCondition *string `url:"cache_condition,omitempty"`
	// Name is the name for the cache settings object.
	Name *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// StaleTTL is the maximum time in seconds to continue to use a stale version of the object if future requests to your backend server fail (also known as 'stale if error').
	StaleTTL *int `url:"stale_ttl,omitempty"`
	// TTL is the maximum time to consider the object fresh in the cache (the cache 'time to live').
	TTL *int `url:"ttl,omitempty"`
}

// CreateCacheSetting creates a new resource.
func (c *Client) CreateCacheSetting(i *CreateCacheSettingInput) (*CacheSetting, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "cache_settings")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs *CacheSetting
	if err := DecodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// GetCacheSettingInput is used as input to the GetCacheSetting function.
type GetCacheSettingInput struct {
	// Name is the name of the cache setting to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetCacheSetting retrieves the specified resource.
func (c *Client) GetCacheSetting(i *GetCacheSettingInput) (*CacheSetting, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "cache_settings", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs *CacheSetting
	if err := DecodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// UpdateCacheSettingInput is used as input to the UpdateCacheSetting function.
type UpdateCacheSettingInput struct {
	// Action determines vcl_fetch behaviour (pass, cache, restart).
	Action *CacheSettingAction `url:"action,omitempty"`
	// CacheCondition is name of the cache condition controlling when this configuration applies.
	CacheCondition *string `url:"cache_condition,omitempty"`
	// Name is the name of the cache setting to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// StaleTTL is the maximum time in seconds to continue to use a stale version of the object if future requests to your backend server fail (also known as 'stale if error').
	StaleTTL *int `url:"stale_ttl,omitempty"`
	// TTL is the maximum time to consider the object fresh in the cache (the cache 'time to live').
	TTL *int `url:"ttl,omitempty"`
}

// UpdateCacheSetting updates the specified resource.
func (c *Client) UpdateCacheSetting(i *UpdateCacheSettingInput) (*CacheSetting, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "cache_settings", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs *CacheSetting
	if err := DecodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// DeleteCacheSettingInput is the input parameter to DeleteCacheSetting.
type DeleteCacheSettingInput struct {
	// Name is the name of the cache setting to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteCacheSetting deletes the specified resource.
func (c *Client) DeleteCacheSetting(i *DeleteCacheSettingInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "cache_settings", i.Name)
	resp, err := c.Delete(path, nil)
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
