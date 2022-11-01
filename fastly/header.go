package fastly

import (
	"fmt"
	"net/url"
	"sort"
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

// PHeaderAction returns pointer to HeaderAction.
func PHeaderAction(t HeaderAction) *HeaderAction {
	ha := HeaderAction(t)
	return &ha
}

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

// PHeaderType returns pointer to HeaderType.
func PHeaderType(t HeaderType) *HeaderType {
	ht := HeaderType(t)
	return &ht
}

// Header represents a header response from the Fastly API.
type Header struct {
	Action            HeaderAction `mapstructure:"action"`
	CacheCondition    string       `mapstructure:"cache_condition"`
	CreatedAt         *time.Time   `mapstructure:"created_at"`
	DeletedAt         *time.Time   `mapstructure:"deleted_at"`
	Destination       string       `mapstructure:"dst"`
	IgnoreIfSet       bool         `mapstructure:"ignore_if_set"`
	Name              string       `mapstructure:"name"`
	Priority          uint         `mapstructure:"priority"`
	Regex             string       `mapstructure:"regex"`
	RequestCondition  string       `mapstructure:"request_condition"`
	ResponseCondition string       `mapstructure:"response_condition"`
	ServiceID         string       `mapstructure:"service_id"`
	ServiceVersion    int          `mapstructure:"version"`
	Source            string       `mapstructure:"src"`
	Substitution      string       `mapstructure:"substitution"`
	Type              HeaderType   `mapstructure:"type"`
	UpdatedAt         *time.Time   `mapstructure:"updated_at"`
}

// headersByName is a sortable list of headers.
type headersByName []*Header

// Len implement the sortable interface.
func (s headersByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s headersByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s headersByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListHeadersInput is used as input to the ListHeaders function.
type ListHeadersInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListHeaders returns the list of headers for the configuration version.
func (c *Client) ListHeaders(i *ListHeadersInput) ([]*Header, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/header", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var bs []*Header
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	sort.Stable(headersByName(bs))
	return bs, nil
}

// CreateHeaderInput is used as input to the CreateHeader function.
type CreateHeaderInput struct {
	Action            HeaderAction `url:"action,omitempty"`
	CacheCondition    string       `url:"cache_condition,omitempty"`
	Destination       string       `url:"dst,omitempty"`
	IgnoreIfSet       Compatibool  `url:"ignore_if_set,omitempty"`
	Name              string       `url:"name,omitempty"`
	Priority          *uint        `url:"priority,omitempty"`
	Regex             string       `url:"regex,omitempty"`
	RequestCondition  string       `url:"request_condition,omitempty"`
	ResponseCondition string       `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Source         string     `url:"src,omitempty"`
	Substitution   string     `url:"substitution,omitempty"`
	Type           HeaderType `url:"type,omitempty"`
}

// CreateHeader creates a new Fastly header.
func (c *Client) CreateHeader(i *CreateHeaderInput) (*Header, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/header", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Header
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetHeaderInput is used as input to the GetHeader function.
type GetHeaderInput struct {
	// Name is the name of the header to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetHeader gets the header configuration with the given parameters.
func (c *Client) GetHeader(i *GetHeaderInput) (*Header, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/header/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Header
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateHeaderInput is used as input to the UpdateHeader function.
type UpdateHeaderInput struct {
	Action         *HeaderAction `url:"action,omitempty"`
	CacheCondition *string       `url:"cache_condition,omitempty"`
	Destination    *string       `url:"dst,omitempty"`
	IgnoreIfSet    *Compatibool  `url:"ignore_if_set,omitempty"`
	// Name is the name of the header to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Priority          *uint   `url:"priority,omitempty"`
	Regex             *string `url:"regex,omitempty"`
	RequestCondition  *string `url:"request_condition,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Source         *string     `url:"src,omitempty"`
	Substitution   *string     `url:"substitution,omitempty"`
	Type           *HeaderType `url:"type,omitempty"`
}

// UpdateHeader updates a specific header.
func (c *Client) UpdateHeader(i *UpdateHeaderInput) (*Header, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/header/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Header
	if err := decodeBodyMap(resp.Body, &b); err != nil {
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

// DeleteHeader deletes the given header version.
func (c *Client) DeleteHeader(i *DeleteHeaderInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/header/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
