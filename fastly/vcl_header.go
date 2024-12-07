package fastly

import (
	"strconv"
	"time"
)

const (
	// HeaderActionSet is a header action that sets or resets a header.
	HeaderActionSet HeaderAction = "set"

	// HeaderActionAppend is a header action that appends to an existing header.
	HeaderActionAppend HeaderAction = "append"

	// HeaderActionDelete is a header action that deletes a header.
	HeaderActionDelete HeaderAction = "delete"

	// HeaderActionRegex is a header action that performs a single regex
	// replacement on a header.
	HeaderActionRegex HeaderAction = "regex"

	// HeaderActionRegexRepeat is a header action that performs a global regex
	// replacement on a header.
	HeaderActionRegexRepeat HeaderAction = "regex_repeat"
)

// HeaderAction is a type of header action.
type HeaderAction string

const (
	// HeaderTypeRequest is a header type that performs on the request before
	// lookups.
	HeaderTypeRequest HeaderType = "request"

	// HeaderTypeFetch is a header type that performs on the request to the origin
	// server.
	HeaderTypeFetch HeaderType = "fetch"

	// HeaderTypeCache is a header type that performs on the response before it's
	// store in the cache.
	HeaderTypeCache HeaderType = "cache"

	// HeaderTypeResponse is a header type that performs on the response before
	// delivering to the client.
	HeaderTypeResponse HeaderType = "response"
)

// HeaderType is a type of header.
type HeaderType string

// Header represents a header response from the Fastly API.
type Header struct {
	Action            *HeaderAction `mapstructure:"action"`
	CacheCondition    *string       `mapstructure:"cache_condition"`
	CreatedAt         *time.Time    `mapstructure:"created_at"`
	DeletedAt         *time.Time    `mapstructure:"deleted_at"`
	Destination       *string       `mapstructure:"dst"`
	IgnoreIfSet       *bool         `mapstructure:"ignore_if_set"`
	Name              *string       `mapstructure:"name"`
	Priority          *int          `mapstructure:"priority"`
	Regex             *string       `mapstructure:"regex"`
	RequestCondition  *string       `mapstructure:"request_condition"`
	ResponseCondition *string       `mapstructure:"response_condition"`
	ServiceID         *string       `mapstructure:"service_id"`
	ServiceVersion    *int          `mapstructure:"version"`
	Source            *string       `mapstructure:"src"`
	Substitution      *string       `mapstructure:"substitution"`
	Type              *HeaderType   `mapstructure:"type"`
	UpdatedAt         *time.Time    `mapstructure:"updated_at"`
}

// ListHeadersInput is used as input to the ListHeaders function.
type ListHeadersInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListHeaders retrieves all resources.
func (c *Client) ListHeaders(i *ListHeadersInput) ([]*Header, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "header")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bs []*Header
	if err := DecodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	return bs, nil
}

// CreateHeaderInput is used as input to the CreateHeader function.
type CreateHeaderInput struct {
	// Action accepts a string value (set, append, delete, regex, regex_repeat).
	Action *HeaderAction `url:"action,omitempty"`
	// CacheCondition is the name of the cache condition controlling when this configuration applies.
	CacheCondition *string `url:"cache_condition,omitempty"`
	// Destination is the header to set.
	Destination *string `url:"dst,omitempty"`
	// IgnoreIfSet prevents adding the header if it is added already. Only applies to 'set' action.
	IgnoreIfSet *Compatibool `url:"ignore_if_set,omitempty"`
	// Name is a handle to refer to this Header object.
	Name *string `url:"name,omitempty"`
	// Priority determines execution order. Lower numbers execute first.
	Priority *int `url:"priority,omitempty"`
	// Regex is the regular expression to use. Only applies to regex and regex_repeat actions.
	Regex *string `url:"regex,omitempty"`
	// RequestCondition is the condition which, if met, will select this configuration during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// ResponseCondition is an optional name of a response condition to apply.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Source is a variable to be used as a source for the header content. Does not apply to delete action.
	Source *string `url:"src,omitempty"`
	// Substitution is a value to substitute in place of regular expression. Only applies to regex and regex_repeat actions.
	Substitution *string `url:"substitution,omitempty"`
	// Type is a type of header (request, cache, response).
	Type *HeaderType `url:"type,omitempty"`
}

// CreateHeader creates a new resource.
func (c *Client) CreateHeader(i *CreateHeaderInput) (*Header, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "header")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Header
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetHeaderInput is used as input to the GetHeader function.
type GetHeaderInput struct {
	// Name is the name of the header to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetHeader retrieves the specified resource.
func (c *Client) GetHeader(i *GetHeaderInput) (*Header, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "header", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Header
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateHeaderInput is used as input to the UpdateHeader function.
type UpdateHeaderInput struct {
	// Action accepts a string value (set, append, delete, regex, regex_repeat).
	Action *HeaderAction `url:"action,omitempty"`
	// CacheCondition is the name of the cache condition controlling when this configuration applies.
	CacheCondition *string `url:"cache_condition,omitempty"`
	// Destination is the header to set.
	Destination *string `url:"dst,omitempty"`
	// IgnoreIfSet prevents adding the header if it is added already. Only applies to 'set' action.
	IgnoreIfSet *Compatibool `url:"ignore_if_set,omitempty"`
	// Name is the name of the header to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Priority determines execution order. Lower numbers execute first.
	Priority *int `url:"priority,omitempty"`
	// Regex is the regular expression to use. Only applies to regex and regex_repeat actions.
	Regex *string `url:"regex,omitempty"`
	// RequestCondition is the condition which, if met, will select this configuration during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// ResponseCondition is an optional name of a response condition to apply.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Source is a variable to be used as a source for the header content. Does not apply to delete action.
	Source *string `url:"src,omitempty"`
	// Substitution is a value to substitute in place of regular expression. Only applies to regex and regex_repeat actions.
	Substitution *string `url:"substitution,omitempty"`
	// Type is a type of header (request, cache, response).
	Type *HeaderType `url:"type,omitempty"`
}

// UpdateHeader updates the specified resource.
func (c *Client) UpdateHeader(i *UpdateHeaderInput) (*Header, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "header", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Header
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteHeaderInput is the input parameter to DeleteHeader.
type DeleteHeaderInput struct {
	// Name is the name of the header to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteHeader deletes the specified resource.
func (c *Client) DeleteHeader(i *DeleteHeaderInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "header", i.Name)
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
