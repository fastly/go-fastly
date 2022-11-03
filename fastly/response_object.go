package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// ResponseObject represents a response object response from the Fastly API.
type ResponseObject struct {
	CacheCondition   string     `mapstructure:"cache_condition"`
	Content          string     `mapstructure:"content"`
	ContentType      string     `mapstructure:"content_type"`
	CreatedAt        *time.Time `mapstructure:"created_at"`
	DeletedAt        *time.Time `mapstructure:"deleted_at"`
	Name             string     `mapstructure:"name"`
	RequestCondition string     `mapstructure:"request_condition"`
	Response         string     `mapstructure:"response"`
	ServiceID        string     `mapstructure:"service_id"`
	ServiceVersion   int        `mapstructure:"version"`
	Status           uint       `mapstructure:"status"`
	UpdatedAt        *time.Time `mapstructure:"updated_at"`
}

// responseObjectsByName is a sortable list of response objects.
type responseObjectsByName []*ResponseObject

// Len implement the sortable interface.
func (s responseObjectsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s responseObjectsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
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

// ListResponseObjects retrieves all resources.
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
	// CacheCondition is the name of the cache condition controlling when this configuration applies.
	CacheCondition string `url:"cache_condition,omitempty"`
	// Content is the content to deliver for the response object, can be empty.
	Content string `url:"content,omitempty"`
	// ContentType is the MIME type of the content, can be empty.
	ContentType string `url:"content_type,omitempty"`
	// Name is the name for the request settings.
	Name string `url:"name,omitempty"`
	// RequestCondition is the condition which, if met, will select this configuration during a request.
	RequestCondition string `url:"request_condition,omitempty"`
	// Response is the HTTP response.
	Response string `url:"response,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// Status is the HTTP status code.
	Status *uint `url:"status,omitempty"`
}

// CreateResponseObject creates a new resource.
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
	// Name is the name of the response object to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetResponseObject retrieves the specified resource.
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
	// CacheCondition is the name of the cache condition controlling when this configuration applies.
	CacheCondition *string `url:"cache_condition,omitempty"`
	// Content is the content to deliver for the response object, can be empty.
	Content *string `url:"content,omitempty"`
	// ContentType is the MIME type of the content, can be empty.
	ContentType *string `url:"content_type,omitempty"`
	// Name is the name of the response object to update.
	Name string
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// RequestCondition is the condition which, if met, will select this configuration during a request.
	RequestCondition *string `url:"request_condition,omitempty"`
	// Response is the HTTP response.
	Response *string `url:"response,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// Status is the HTTP status code.
	Status *uint `url:"status,omitempty"`
}

// UpdateResponseObject updates the specified resource.
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
	// Name is the name of the response object to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteResponseObject deletes the specified resource.
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
