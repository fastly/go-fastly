package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// FTP represents an FTP logging response from the Fastly API.
type FTP struct {
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
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	Username          string     `mapstructure:"user"`
}

// ftpsByName is a sortable list of ftps.
type ftpsByName []*FTP

// Len implement the sortable interface.
func (s ftpsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s ftpsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
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

// ListFTPs retrieves all resources.
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
	defer resp.Body.Close()

	var ftps []*FTP
	if err := decodeBodyMap(resp.Body, &ftps); err != nil {
		return nil, err
	}
	sort.Stable(ftpsByName(ftps))
	return ftps, nil
}

// CreateFTPInput is used as input to the CreateFTP function.
type CreateFTPInput struct {
	// Address is an hostname or IPv4 address.
	Address string `url:"address,omitempty"`
	// CompressionCodec is the codec used for compressing your logs (zstd, snappy, gzip).
	CompressionCodec string `url:"compression_codec,omitempty"`
	// Format is a Fastly log format string.
	Format string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion uint `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel uint8 `url:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType string `url:"message_type,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name string `url:"name,omitempty"`
	// Password is the password for the server. For anonymous use an email address.
	Password string `url:"password,omitempty"`
	// Path is the path to upload log files to. If the path ends in / then it is treated as a directory.
	Path string `url:"path,omitempty"`
	// Period is how frequently log files are finalized so they can be available for reading (in seconds).
	Period uint `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement string `url:"placement,omitempty"`
	// Port is the port number.
	Port uint `url:"port,omitempty"`
	// PublicKey is a PGP public key that Fastly will use to encrypt your log files before writing them to disk.
	PublicKey string `url:"public_key,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// TimestampFormat is a timestamp format.
	TimestampFormat string `url:"timestamp_format,omitempty"`
	// Username is the username for the server. Can be anonymous.
	Username string `url:"user,omitempty"`
}

// CreateFTP creates a new resource.
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
	defer resp.Body.Close()

	var ftp *FTP
	if err := decodeBodyMap(resp.Body, &ftp); err != nil {
		return nil, err
	}
	return ftp, nil
}

// GetFTPInput is used as input to the GetFTP function.
type GetFTPInput struct {
	// Name is the name of the FTP to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetFTP retrieves the specified resource.
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
	defer resp.Body.Close()

	var b *FTP
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateFTPInput is used as input to the UpdateFTP function.
type UpdateFTPInput struct {
	// Address is an hostname or IPv4 address.
	Address *string `url:"address,omitempty"`
	// CompressionCodec is the codec used for compressing your logs (zstd, snappy, gzip).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *uint `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *uint8 `url:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name of the FTP to update.
	Name string
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Password is the password for the server. For anonymous use an email address.
	Password *string `url:"password,omitempty"`
	// Path is the path to upload log files to. If the path ends in / then it is treated as a directory.
	Path *string `url:"path,omitempty"`
	// Period is how frequently log files are finalized so they can be available for reading (in seconds).
	Period *uint `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Port is the port number.
	Port *uint `url:"port,omitempty"`
	// PublicKey is a PGP public key that Fastly will use to encrypt your log files before writing them to disk.
	PublicKey *string `url:"public_key,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
	// Username is the username for the server. Can be anonymous.
	Username *string `url:"user,omitempty"`
}

// UpdateFTP updates the specified resource.
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
	defer resp.Body.Close()

	var b *FTP
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteFTPInput is the input parameter to DeleteFTP.
type DeleteFTPInput struct {
	// Name is the name of the FTP to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteFTP deletes the specified resource.
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
