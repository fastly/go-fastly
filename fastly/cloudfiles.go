package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Cloudfiles represents a Cloudfiles response from the Fastly API.
type Cloudfiles struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	User              string     `mapstructure:"user"`
	AccessKey         string     `mapstructure:"access_key"`
	BucketName        string     `mapstructure:"bucket_name"`
	Path              string     `mapstructure:"path"`
	Region            string     `mapstructure:"region"`
	Placement         string     `mapstructure:"placement"`
	Period            uint       `mapstructure:"period"`
	GzipLevel         uint       `mapstructure:"gzip_level"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	MessageType       string     `mapstructure:"message_type"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	PublicKey         string     `mapstructure:"public_key"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// cloudfilesByName is a sortable list of Cloudfiles.
type cloudfilesByName []*Cloudfiles

// Len, Swap, and Less implement the sortable interface.
func (c cloudfilesByName) Len() int      { return len(c) }
func (c cloudfilesByName) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c cloudfilesByName) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

// ListCloudfilesInput is used as input to the ListCloudfiles function.
type ListCloudfilesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListCloudfiles returns the list of Cloudfiles for the configuration version.
func (c *Client) ListCloudfiles(i *ListCloudfilesInput) ([]*Cloudfiles, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/cloudfiles", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var cloudfiles []*Cloudfiles
	if err := decodeBodyMap(resp.Body, &cloudfiles); err != nil {
		return nil, err
	}
	sort.Stable(cloudfilesByName(cloudfiles))
	return cloudfiles, nil
}

// CreateCloudfilesInput is used as input to the CreateCloudfiles function.
type CreateCloudfilesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `url:"name,omitempty"`
	User              string `url:"user,omitempty"`
	AccessKey         string `url:"access_key,omitempty"`
	BucketName        string `url:"bucket_name,omitempty"`
	Path              string `url:"path,omitempty"`
	Region            string `url:"region,omitempty"`
	Placement         string `url:"placement,omitempty"`
	Period            uint   `url:"period,omitempty"`
	GzipLevel         uint   `url:"gzip_level,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	MessageType       string `url:"message_type,omitempty"`
	TimestampFormat   string `url:"timestamp_format,omitempty"`
	PublicKey         string `url:"public_key,omitempty"`
	CompressionCodec  string `url:"compression_codec,omitempty"`
}

// CreateCloudfiles creates a new Fastly Cloudfiles.
func (c *Client) CreateCloudfiles(i *CreateCloudfilesInput) (*Cloudfiles, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/cloudfiles", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var cloudfiles *Cloudfiles
	if err := decodeBodyMap(resp.Body, &cloudfiles); err != nil {
		return nil, err
	}
	return cloudfiles, nil
}

// GetCloudfilesInput is used as input to the GetCloudfiles function.
type GetCloudfilesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Cloudfiles to fetch.
	Name string
}

// GetCloudfiles gets the Cloudfiles configuration with the given parameters.
func (c *Client) GetCloudfiles(i *GetCloudfilesInput) (*Cloudfiles, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/cloudfiles/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var cloudfiles *Cloudfiles
	if err := decodeBodyMap(resp.Body, &cloudfiles); err != nil {
		return nil, err
	}
	return cloudfiles, nil
}

// UpdateCloudfilesInput is used as input to the UpdateCloudfiles function.
type UpdateCloudfilesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Cloudfiles to update.
	Name string

	NewName           *string `url:"name,omitempty"`
	User              *string `url:"user,omitempty"`
	AccessKey         *string `url:"access_key,omitempty"`
	BucketName        *string `url:"bucket_name,omitempty"`
	Path              *string `url:"path,omitempty"`
	Region            *string `url:"region,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	Period            *uint   `url:"period,omitempty"`
	GzipLevel         *uint   `url:"gzip_level,omitempty"`
	Format            *string `url:"format,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	MessageType       *string `url:"message_type,omitempty"`
	TimestampFormat   *string `url:"timestamp_format,omitempty"`
	PublicKey         *string `url:"public_key,omitempty"`
	CompressionCodec  *string `url:"compression_codec,omitempty"`
}

// UpdateCloudfiles updates a specific Cloudfiles.
func (c *Client) UpdateCloudfiles(i *UpdateCloudfilesInput) (*Cloudfiles, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/cloudfiles/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var cloudfiles *Cloudfiles
	if err := decodeBodyMap(resp.Body, &cloudfiles); err != nil {
		return nil, err
	}
	return cloudfiles, nil
}

// DeleteCloudfilesInput is the input parameter to DeleteCloudfiles.
type DeleteCloudfilesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Cloudfiles to delete (required).
	Name string
}

// DeleteCloudfiles deletes the given Cloudfiles version.
func (c *Client) DeleteCloudfiles(i *DeleteCloudfilesInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/cloudfiles/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
