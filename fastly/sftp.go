package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// SFTP represents an SFTP logging response from the Fastly API.
type SFTP struct {
	Address           string     `mapstructure:"address"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	GzipLevel         uint8      `mapstructure:"gzip_level"`
	MessageType       string     `mapstructure:"message_type"`
	Name              string     `mapstructure:"name"`
	Password          string     `mapstructure:"password"`
	Path              string     `mapstructure:"path"`
	Period            uint       `mapstructure:"period"`
	Placement         string     `mapstructure:"placement"`
	Port              uint       `mapstructure:"port"`
	PublicKey         string     `mapstructure:"public_key"`
	ResponseCondition string     `mapstructure:"response_condition"`
	SSHKnownHosts     string     `mapstructure:"ssh_known_hosts"`
	SecretKey         string     `mapstructure:"secret_key"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	User              string     `mapstructure:"user"`
}

// sftpsByName is a sortable list of sftps.
type sftpsByName []*SFTP

// Len implement the sortable interface.
func (s sftpsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s sftpsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
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
	defer resp.Body.Close()

	var sftps []*SFTP
	if err := decodeBodyMap(resp.Body, &sftps); err != nil {
		return nil, err
	}
	sort.Stable(sftpsByName(sftps))
	return sftps, nil
}

// CreateSFTPInput is used as input to the CreateSFTP function.
type CreateSFTPInput struct {
	Address           string `url:"address,omitempty"`
	CompressionCodec  string `url:"compression_codec,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	GzipLevel         uint8  `url:"gzip_level,omitempty"`
	MessageType       string `url:"message_type,omitempty"`
	Name              string `url:"name,omitempty"`
	Password          string `url:"password,omitempty"`
	Path              string `url:"path,omitempty"`
	Period            uint   `url:"period,omitempty"`
	Placement         string `url:"placement,omitempty"`
	Port              uint   `url:"port,omitempty"`
	PublicKey         string `url:"public_key,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	SSHKnownHosts     string `url:"ssh_known_hosts,omitempty"`
	SecretKey         string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion  int
	TimestampFormat string `url:"timestamp_format,omitempty"`
	User            string `url:"user,omitempty"`
}

// CreateSFTP creates a new resource.
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
	defer resp.Body.Close()

	var ftp *SFTP
	if err := decodeBodyMap(resp.Body, &ftp); err != nil {
		return nil, err
	}
	return ftp, nil
}

// GetSFTPInput is used as input to the GetSFTP function.
type GetSFTPInput struct {
	// Name is the name of the SFTP to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
	defer resp.Body.Close()

	var b *SFTP
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateSFTPInput is used as input to the UpdateSFTP function.
type UpdateSFTPInput struct {
	Address          *string `url:"address,omitempty"`
	CompressionCodec *string `url:"compression_codec,omitempty"`
	Format           *string `url:"format,omitempty"`
	FormatVersion    *uint   `url:"format_version,omitempty"`
	GzipLevel        *uint8  `url:"gzip_level,omitempty"`
	MessageType      *string `url:"message_type,omitempty"`
	// Name is the name of the SFTP to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Password          *string `url:"password,omitempty"`
	Path              *string `url:"path,omitempty"`
	Period            *uint   `url:"period,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	Port              *uint   `url:"port,omitempty"`
	PublicKey         *string `url:"public_key,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	SSHKnownHosts     *string `url:"ssh_known_hosts,omitempty"`
	SecretKey         *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion  int
	TimestampFormat *string `url:"timestamp_format,omitempty"`
	User            *string `url:"user,omitempty"`
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
	defer resp.Body.Close()

	var b *SFTP
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteSFTPInput is the input parameter to DeleteSFTP.
type DeleteSFTPInput struct {
	// Name is the name of the SFTP to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
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
