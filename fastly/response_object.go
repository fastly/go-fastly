package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// ResponseObject represents a response object response from the Fastly API.
type ResponseObject struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name             string     `mapstructure:"name"`
	Status           uint       `mapstructure:"status"`
	Response         string     `mapstructure:"response"`
	Content          string     `mapstructure:"content"`
	ContentType      string     `mapstructure:"content_type"`
	RequestCondition string     `mapstructure:"request_condition"`
	CacheCondition   string     `mapstructure:"cache_condition"`
	CreatedAt        *time.Time `mapstructure:"created_at"`
	UpdatedAt        *time.Time `mapstructure:"updated_at"`
	DeletedAt        *time.Time `mapstructure:"deleted_at"`
}

// responseObjectsByName is a sortable list of response objects.
type responseObjectsByName []*ResponseObject

// Len, Swap, and Less implement the sortable interface.
func (s responseObjectsByName) Len() int      { return len(s) }
func (s responseObjectsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s responseObjectsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListResponseObjectsInput is used as input to the ListResponseObjects
// function.
type ListResponseObjectsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListResponseObjects returns the list of response objects for the
// configuration version.
func (c *Client) ListResponseObjects(i *ListResponseObjectsInput) ([]*ResponseObject, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/response_object", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bs []*ResponseObject
	if err := decodeBodyMap(resp.Body, &bs); err != nil {
		return nil, err
	}
	sort.Stable(responseObjectsByName(bs))
	return bs, nil
}

// CreateResponseObjectInput is used as input to the CreateResponseObject
// function.
type CreateResponseObjectInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name             string `url:"name,omitempty"`
	Status           *uint  `url:"status,omitempty"`
	Response         string `url:"response,omitempty"`
	Content          string `url:"content,omitempty"`
	ContentType      string `url:"content_type,omitempty"`
	RequestCondition string `url:"request_condition,omitempty"`
	CacheCondition   string `url:"cache_condition,omitempty"`
}

// CreateResponseObject creates a new Fastly response object.
func (c *Client) CreateResponseObject(i *CreateResponseObjectInput) (*ResponseObject, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/response_object", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *ResponseObject
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetResponseObjectInput is used as input to the GetResponseObject function.
type GetResponseObjectInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the response object to fetch.
	Name string
}

// GetResponseObject gets the response object configuration with the given
// parameters.
func (c *Client) GetResponseObject(i *GetResponseObjectInput) (*ResponseObject, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/response_object/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *ResponseObject
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateResponseObjectInput is used as input to the UpdateResponseObject
// function.
type UpdateResponseObjectInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the response object to update.
	Name string

	NewName          *string `url:"name,omitempty"`
	Status           *uint   `url:"status,omitempty"`
	Response         *string `url:"response,omitempty"`
	Content          *string `url:"content,omitempty"`
	ContentType      *string `url:"content_type,omitempty"`
	RequestCondition *string `url:"request_condition,omitempty"`
	CacheCondition   *string `url:"cache_condition,omitempty"`
}

// UpdateResponseObject updates a specific response object.
func (c *Client) UpdateResponseObject(i *UpdateResponseObjectInput) (*ResponseObject, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/response_object/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *ResponseObject
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteResponseObjectInput is the input parameter to DeleteResponseObject.
type DeleteResponseObjectInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the response object to delete (required).
	Name string
}

// DeleteResponseObject deletes the given response object version.
func (c *Client) DeleteResponseObject(i *DeleteResponseObjectInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/response_object/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
