package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// SFTP represents an SFTP logging response from the Fastly API.
type SFTP struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Address           string     `mapstructure:"address"`
	Port              uint       `mapstructure:"port"`
	User              string     `mapstructure:"user"`
	Password          string     `mapstructure:"password"`
	PublicKey         string     `mapstructure:"public_key"`
	SecretKey         string     `mapstructure:"secret_key"`
	SSHKnownHosts     string     `mapstructure:"ssh_known_hosts"`
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

// sftpsByName is a sortable list of sftps.
type sftpsByName []*SFTP

// Len, Swap, and Less implement the sortable interface.
func (s sftpsByName) Len() int      { return len(s) }
func (s sftpsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sftpsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListSFTPsInput is used as input to the ListSFTPs function.
type ListSFTPsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListSFTPs returns the list of sftps for the configuration version.
func (c *Client) ListSFTPs(i *ListSFTPsInput) ([]*SFTP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sftp", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var sftps []*SFTP
	if err := decodeBodyMap(resp.Body, &sftps); err != nil {
		return nil, err
	}
	sort.Stable(sftpsByName(sftps))
	return sftps, nil
}

// CreateSFTPInput is used as input to the CreateSFTP function.
type CreateSFTPInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `form:"name,omitempty"`
	Address           string `form:"address,omitempty"`
	Port              uint   `form:"port,omitempty"`
	User              string `form:"user,omitempty"`
	Password          string `form:"password,omitempty"`
	PublicKey         string `form:"public_key,omitempty"`
	SecretKey         string `form:"secret_key,omitempty"`
	SSHKnownHosts     string `form:"ssh_known_hosts,omitempty"`
	Path              string `form:"path,omitempty"`
	Period            uint   `form:"period,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	CompressionCodec  string `form:"compression_codec,omitempty"`
	GzipLevel         uint   `form:"gzip_level,omitempty"`
	Format            string `form:"format,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	TimestampFormat   string `form:"timestamp_format,omitempty"`
	MessageType       string `form:"message_type,omitempty"`
	Placement         string `form:"placement,omitempty"`
}

// CreateSFTP creates a new Fastly SFTP.
func (c *Client) CreateSFTP(i *CreateSFTPInput) (*SFTP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sftp", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var ftp *SFTP
	if err := decodeBodyMap(resp.Body, &ftp); err != nil {
		return nil, err
	}
	return ftp, nil
}

// GetSFTPInput is used as input to the GetSFTP function.
type GetSFTPInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the SFTP to fetch.
	Name string
}

// GetSFTP gets the SFTP configuration with the given parameters.
func (c *Client) GetSFTP(i *GetSFTPInput) (*SFTP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sftp/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var b *SFTP
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateSFTPInput is used as input to the UpdateSFTP function.
type UpdateSFTPInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the SFTP to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Address           *string `form:"address,omitempty"`
	Port              *uint   `form:"port,omitempty"`
	PublicKey         *string `form:"public_key,omitempty"`
	SecretKey         *string `form:"secret_key,omitempty"`
	SSHKnownHosts     *string `form:"ssh_known_hosts,omitempty"`
	User              *string `form:"user,omitempty"`
	Password          *string `form:"password,omitempty"`
	Path              *string `form:"path,omitempty"`
	Period            *uint   `form:"period,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	CompressionCodec  *string `form:"compression_codec,omitempty"`
	GzipLevel         *uint   `form:"gzip_level,omitempty"`
	Format            *string `form:"format,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	TimestampFormat   *string `form:"timestamp_format,omitempty"`
	MessageType       *string `form:"message_type,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateSFTP updates a specific SFTP.
func (c *Client) UpdateSFTP(i *UpdateSFTPInput) (*SFTP, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sftp/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *SFTP
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteSFTPInput is the input parameter to DeleteSFTP.
type DeleteSFTPInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the SFTP to delete (required).
	Name string
}

// DeleteSFTP deletes the given SFTP version.
func (c *Client) DeleteSFTP(i *DeleteSFTPInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sftp/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
