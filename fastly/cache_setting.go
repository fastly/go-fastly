package fastly

import (
	"fmt"
	"net/url"
	"sort"
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
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name           string             `mapstructure:"name"`
	Action         CacheSettingAction `mapstructure:"action"`
	TTL            uint               `mapstructure:"ttl"`
	StaleTTL       uint               `mapstructure:"stale_ttl"`
	CacheCondition string             `mapstructure:"cache_condition"`
	CreatedAt      *time.Time         `mapstructure:"created_at"`
	UpdatedAt      *time.Time         `mapstructure:"updated_at"`
	DeletedAt      *time.Time         `mapstructure:"deleted_at"`
}

// cacheSettingsByName is a sortable list of cache settings.
type cacheSettingsByName []*CacheSetting

// Len, Swap, and Less implement the sortable interface.
func (s cacheSettingsByName) Len() int      { return len(s) }
func (s cacheSettingsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s cacheSettingsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListCacheSettingsInput is used as input to the ListCacheSettings function.
type ListCacheSettingsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListCacheSettings returns the list of cache settings for the configuration
// version.
func (c *Client) ListCacheSettings(i *ListCacheSettingsInput) ([]*CacheSetting, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/cache_settings", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var cs []*CacheSetting
	if err := decodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	sort.Stable(cacheSettingsByName(cs))
	return cs, nil
}

// CreateCacheSettingInput is used as input to the CreateCacheSetting function.
type CreateCacheSettingInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name           string             `url:"name,omitempty"`
	Action         CacheSettingAction `url:"action,omitempty"`
	TTL            uint               `url:"ttl,omitempty"`
	StaleTTL       uint               `url:"stale_ttl,omitempty"`
	CacheCondition string             `url:"cache_condition,omitempty"`
}

// CreateCacheSetting creates a new Fastly cache setting.
func (c *Client) CreateCacheSetting(i *CreateCacheSettingInput) (*CacheSetting, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/cache_settings", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var cs *CacheSetting
	if err := decodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// GetCacheSettingInput is used as input to the GetCacheSetting function.
type GetCacheSettingInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the cache setting to fetch.
	Name string
}

// GetCacheSetting gets the cache setting configuration with the given
// parameters.
func (c *Client) GetCacheSetting(i *GetCacheSettingInput) (*CacheSetting, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/cache_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var cs *CacheSetting
	if err := decodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// UpdateCacheSettingInput is used as input to the UpdateCacheSetting function.
type UpdateCacheSettingInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the cache setting to update.
	Name string

	NewName        *string            `url:"name,omitempty"`
	Action         CacheSettingAction `url:"action,omitempty"`
	TTL            *uint              `url:"ttl,omitempty"`
	StaleTTL       *uint              `url:"stale_ttl,omitempty"`
	CacheCondition *string            `url:"cache_condition,omitempty"`
}

// UpdateCacheSetting updates a specific cache setting.
func (c *Client) UpdateCacheSetting(i *UpdateCacheSettingInput) (*CacheSetting, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/cache_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var cs *CacheSetting
	if err := decodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// DeleteCacheSettingInput is the input parameter to DeleteCacheSetting.
type DeleteCacheSettingInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the cache setting to delete (required).
	Name string
}

// DeleteCacheSetting deletes the given cache setting version.
func (c *Client) DeleteCacheSetting(i *DeleteCacheSettingInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/cache_settings/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
