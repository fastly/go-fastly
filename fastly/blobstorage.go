package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// BlobStorage represents a blob storage response from the Fastly API.
type BlobStorage struct {
	AccountName       string     `mapstructure:"account_name"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	Container         string     `mapstructure:"container"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	FileMaxBytes      uint       `mapstructure:"file_max_bytes"`
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
	SASToken          string     `mapstructure:"sas_token"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// blobStorageByName is a sortable list of blob storages.
type blobStorageByName []*BlobStorage

// Len implement the sortable interface.
func (s blobStorageByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s blobStorageByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s blobStorageByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListBlobStoragesInput is used as input to the ListBlobStorages function.
type ListBlobStoragesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListBlobStorages retrieves all resources.
func (c *Client) ListBlobStorages(i *ListBlobStoragesInput) ([]*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var as []*BlobStorage
	if err := decodeBodyMap(resp.Body, &as); err != nil {
		return nil, err
	}
	sort.Stable(blobStorageByName(as))
	return as, nil
}

// CreateBlobStorageInput is used as input to the CreateBlobStorage function.
type CreateBlobStorageInput struct {
	AccountName       string `url:"account_name,omitempty"`
	CompressionCodec  string `url:"compression_codec,omitempty"`
	Container         string `url:"container,omitempty"`
	FileMaxBytes      uint   `url:"file_max_bytes,omitempty"`
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
	SASToken          string `url:"sas_token,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion  int
	TimestampFormat string `url:"timestamp_format,omitempty"`
}

// CreateBlobStorage creates a new resource.
func (c *Client) CreateBlobStorage(i *CreateBlobStorageInput) (*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// GetBlobStorageInput is used as input to the GetBlobStorage function.
type GetBlobStorageInput struct {
	// Name is the name of the blob storage to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetBlobStorage gets the blob storage configuration with the given parameters.
func (c *Client) GetBlobStorage(i *GetBlobStorageInput) (*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateBlobStorageInput is used as input to the UpdateBlobStorage function.
type UpdateBlobStorageInput struct {
	AccountName      *string `url:"account_name,omitempty"`
	CompressionCodec *string `url:"compression_codec,omitempty"`
	Container        *string `url:"container,omitempty"`
	FileMaxBytes     *uint   `url:"file_max_bytes,omitempty"`
	Format           *string `url:"format,omitempty"`
	FormatVersion    *uint   `url:"format_version,omitempty"`
	GzipLevel        *uint8  `url:"gzip_level,omitempty"`
	MessageType      *string `url:"message_type,omitempty"`
	// Name is the name of the blob storage to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Path              *string `url:"path,omitempty"`
	Period            *uint   `url:"period,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	PublicKey         *string `url:"public_key,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	SASToken          *string `url:"sas_token,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion  int
	TimestampFormat *string `url:"timestamp_format,omitempty"`
}

// UpdateBlobStorage updates a specific blob storage.
func (c *Client) UpdateBlobStorage(i *UpdateBlobStorageInput) (*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteBlobStorageInput is the input parameter to DeleteBlobStorage.
type DeleteBlobStorageInput struct {
	// Name is the name of the blob storage to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteBlobStorage deletes the specified resource.
func (c *Client) DeleteBlobStorage(i *DeleteBlobStorageInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
