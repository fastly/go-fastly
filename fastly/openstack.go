package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Openstack represents a Openstack response from the Fastly API.
type Openstack struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	User              string     `mapstructure:"user"`
	AccessKey         string     `mapstructure:"access_key"`
	BucketName        string     `mapstructure:"bucket_name"`
	URL               string     `mapstructure:"url"`
	Path              string     `mapstructure:"path"`
	Placement         string     `mapstructure:"placement"`
	Period            uint       `mapstructure:"period"`
	GzipLevel         uint       `mapstructure:"gzip_level"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	MessageType       string     `mapstructure:"message_type"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	PublicKey         string     `mapstructure:"public_key"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// openstacksByName is a sortable list of Openstack.
type openstacksByName []*Openstack

// Len, Swap, and Less implement the sortable interface.
func (o openstacksByName) Len() int      { return len(o) }
func (o openstacksByName) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o openstacksByName) Less(i, j int) bool {
	return o[i].Name < o[j].Name
}

// ListOpenstackInput is used as input to the ListOpenstack function.
type ListOpenstackInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListOpenstack returns the list of Openstack for the configuration version.
func (c *Client) ListOpenstack(i *ListOpenstackInput) ([]*Openstack, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/openstack", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var openstacks []*Openstack
	if err := decodeBodyMap(resp.Body, &openstacks); err != nil {
		return nil, err
	}
	sort.Stable(openstacksByName(openstacks))
	return openstacks, nil
}

// CreateOpenstackInput is used as input to the CreateOpenstack function.
type CreateOpenstackInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `form:"name,omitempty"`
	User              string `form:"user,omitempty"`
	AccessKey         string `form:"access_key,omitempty"`
	BucketName        string `form:"bucket_name,omitempty"`
	URL               string `form:"url,omitempty"`
	Path              string `form:"path,omitempty"`
	Placement         string `form:"placement,omitempty"`
	Period            uint   `form:"period,omitempty"`
	GzipLevel         uint   `form:"gzip_level,omitempty"`
	Format            string `form:"format,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	MessageType       string `form:"message_type,omitempty"`
	TimestampFormat   string `form:"timestamp_format,omitempty"`
	PublicKey         string `form:"public_key,omitempty"`
}

// CreateOpenstack creates a new Fastly Openstack.
func (c *Client) CreateOpenstack(i *CreateOpenstackInput) (*Openstack, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/openstack", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var openstack *Openstack
	if err := decodeBodyMap(resp.Body, &openstack); err != nil {
		return nil, err
	}
	return openstack, nil
}

// GetOpenstackInput is used as input to the GetOpenstack function.
type GetOpenstackInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Openstack to fetch.
	Name string
}

// GetOpenstack gets the Openstack configuration with the given parameters.
func (c *Client) GetOpenstack(i *GetOpenstackInput) (*Openstack, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/openstack/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var openstack *Openstack
	if err := decodeBodyMap(resp.Body, &openstack); err != nil {
		return nil, err
	}
	return openstack, nil
}

// UpdateOpenstackInput is used as input to the UpdateOpenstack function.
type UpdateOpenstackInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Openstack to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	AccessKey         *string `form:"access_key,omitempty"`
	BucketName        *string `form:"bucket_name,omitempty"`
	URL               *string `form:"url,omitempty"`
	User              *string `form:"user,omitempty"`
	Path              *string `form:"path,omitempty"`
	Placement         *string `form:"placement,omitempty"`
	Period            *uint   `form:"period,omitempty"`
	GzipLevel         *uint   `form:"gzip_level,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	MessageType       *string `form:"message_type,omitempty"`
	TimestampFormat   *string `form:"timestamp_format,omitempty"`
	PublicKey         *string `form:"public_key,omitempty"`
}

// UpdateOpenstack updates a specific Openstack.
func (c *Client) UpdateOpenstack(i *UpdateOpenstackInput) (*Openstack, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/openstack/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var openstack *Openstack
	if err := decodeBodyMap(resp.Body, &openstack); err != nil {
		return nil, err
	}
	return openstack, nil
}

// DeleteOpenstackInput is the input parameter to DeleteOpenstack.
type DeleteOpenstackInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Openstack to delete (required).
	Name string
}

// DeleteOpenstack deletes the given Openstack version.
func (c *Client) DeleteOpenstack(i *DeleteOpenstackInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/openstack/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
