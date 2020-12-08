package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// FTP represents an FTP logging response from the Fastly API.
type FTP struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Address           string     `mapstructure:"address"`
	Port              uint       `mapstructure:"port"`
	Username          string     `mapstructure:"user"`
	Password          string     `mapstructure:"password"`
	PublicKey         string     `mapstructure:"public_key"`
	Path              string     `mapstructure:"path"`
	Period            uint       `mapstructure:"period"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	GzipLevel         uint8      `mapstructure:"gzip_level"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	MessageType       string     `mapstructure:"message_type"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// ftpsByName is a sortable list of ftps.
type ftpsByName []*FTP

// Len, Swap, and Less implement the sortable interface.
func (s ftpsByName) Len() int      { return len(s) }
func (s ftpsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ftpsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListFTPsInput is used as input to the ListFTPs function.
type ListFTPsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListFTPs returns the list of ftps for the configuration version.
func (c *Client) ListFTPs(i *ListFTPsInput) ([]*FTP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/ftp", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ftps []*FTP
	if err := decodeBodyMap(resp.Body, &ftps); err != nil {
		return nil, err
	}
	sort.Stable(ftpsByName(ftps))
	return ftps, nil
}

// CreateFTPInput is used as input to the CreateFTP function.
type CreateFTPInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `form:"name,omitempty"`
	Address           string `form:"address,omitempty"`
	Port              uint   `form:"port,omitempty"`
	Username          string `form:"user,omitempty"`
	Password          string `form:"password,omitempty"`
	PublicKey         string `form:"public_key,omitempty"`
	Path              string `form:"path,omitempty"`
	Period            uint   `form:"period,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	CompressionCodec  string `form:"compression_codec,omitempty"`
	GzipLevel         uint8  `form:"gzip_level,omitempty"`
	Format            string `form:"format,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	MessageType       string `form:"message_type,omitempty"`
	TimestampFormat   string `form:"timestamp_format,omitempty"`
	Placement         string `form:"placement,omitempty"`
}

// CreateFTP creates a new Fastly FTP.
func (c *Client) CreateFTP(i *CreateFTPInput) (*FTP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/ftp", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var ftp *FTP
	if err := decodeBodyMap(resp.Body, &ftp); err != nil {
		return nil, err
	}
	return ftp, nil
}

// GetFTPInput is used as input to the GetFTP function.
type GetFTPInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the FTP to fetch.
	Name string
}

// GetFTP gets the FTP configuration with the given parameters.
func (c *Client) GetFTP(i *GetFTPInput) (*FTP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/ftp/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var b *FTP
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateFTPInput is used as input to the UpdateFTP function.
type UpdateFTPInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the FTP to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Address           *string `form:"address,omitempty"`
	Port              *uint   `form:"port,omitempty"`
	PublicKey         *string `form:"public_key,omitempty"`
	Username          *string `form:"user,omitempty"`
	Password          *string `form:"password,omitempty"`
	Path              *string `form:"path,omitempty"`
	Period            *uint   `form:"period,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	CompressionCodec  *string `form:"compression_codec,omitempty"`
	GzipLevel         *uint8  `form:"gzip_level,omitempty"`
	Format            *string `form:"format,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	MessageType       *string `form:"message_type,omitempty"`
	TimestampFormat   *string `form:"timestamp_format,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateFTP updates a specific FTP.
func (c *Client) UpdateFTP(i *UpdateFTPInput) (*FTP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/ftp/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *FTP
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteFTPInput is the input parameter to DeleteFTP.
type DeleteFTPInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the FTP to delete (required).
	Name string
}

// DeleteFTP deletes the given FTP version.
func (c *Client) DeleteFTP(i *DeleteFTPInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/ftp/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
