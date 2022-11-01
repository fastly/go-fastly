package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Openstack represents a Openstack response from the Fastly API.
type Openstack struct {
	AccessKey         string     `mapstructure:"access_key"`
	BucketName        string     `mapstructure:"bucket_name"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	GzipLevel         uint8      `mapstructure:"gzip_level"`
	MessageType       string     `mapstructure:"message_type"`
	Name              string     `mapstructure:"name"`
	Path              string     `mapstructure:"path"`
	Period            uint       `mapstructure:"period"`
	Placement         string     `mapstructure:"placement"`
	PublicKey         string     `mapstructure:"public_key"`
	ResponseCondition string     `mapstructure:"response_condition"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	URL               string     `mapstructure:"url"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	User              string     `mapstructure:"user"`
}

// openstacksByName is a sortable list of Openstack.
type openstacksByName []*Openstack

// Len implement the sortable interface.
func (o openstacksByName) Len() int {
	return len(o)
}

// Swap implement the sortable interface.
func (o openstacksByName) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Less implement the sortable interface.
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
	defer resp.Body.Close()

	var openstacks []*Openstack
	if err := decodeBodyMap(resp.Body, &openstacks); err != nil {
		return nil, err
	}
	sort.Stable(openstacksByName(openstacks))
	return openstacks, nil
}

// CreateOpenstackInput is used as input to the CreateOpenstack function.
type CreateOpenstackInput struct {
	AccessKey         string `url:"access_key,omitempty"`
	BucketName        string `url:"bucket_name,omitempty"`
	CompressionCodec  string `url:"compression_codec,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	GzipLevel         uint8  `url:"gzip_level,omitempty"`
	MessageType       string `url:"message_type,omitempty"`
	Name              string `url:"name,omitempty"`
	Path              string `url:"path,omitempty"`
	Period            uint   `url:"period,omitempty"`
	Placement         string `url:"placement,omitempty"`
	PublicKey         string `url:"public_key,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion  int
	TimestampFormat string `url:"timestamp_format,omitempty"`
	URL             string `url:"url,omitempty"`
	User            string `url:"user,omitempty"`
}

// CreateOpenstack creates a new resource.
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
	defer resp.Body.Close()

	var openstack *Openstack
	if err := decodeBodyMap(resp.Body, &openstack); err != nil {
		return nil, err
	}
	return openstack, nil
}

// GetOpenstackInput is used as input to the GetOpenstack function.
type GetOpenstackInput struct {
	// Name is the name of the Openstack to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
	defer resp.Body.Close()

	var openstack *Openstack
	if err := decodeBodyMap(resp.Body, &openstack); err != nil {
		return nil, err
	}
	return openstack, nil
}

// UpdateOpenstackInput is used as input to the UpdateOpenstack function.
type UpdateOpenstackInput struct {
	AccessKey        *string `url:"access_key,omitempty"`
	BucketName       *string `url:"bucket_name,omitempty"`
	CompressionCodec *string `url:"compression_codec,omitempty"`
	Format           *string `url:"format,omitempty"`
	FormatVersion    *uint   `url:"format_version,omitempty"`
	GzipLevel        *uint8  `url:"gzip_level,omitempty"`
	MessageType      *string `url:"message_type,omitempty"`
	// Name is the name of the Openstack to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Path              *string `url:"path,omitempty"`
	Period            *uint   `url:"period,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	PublicKey         *string `url:"public_key,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion  int
	TimestampFormat *string `url:"timestamp_format,omitempty"`
	URL             *string `url:"url,omitempty"`
	User            *string `url:"user,omitempty"`
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
	defer resp.Body.Close()

	var openstack *Openstack
	if err := decodeBodyMap(resp.Body, &openstack); err != nil {
		return nil, err
	}
	return openstack, nil
}

// DeleteOpenstackInput is the input parameter to DeleteOpenstack.
type DeleteOpenstackInput struct {
	// Name is the name of the Openstack to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteOpenstack deletes the specified resource.
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
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
