package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// AzureBlob represents an AzureBlob logging response from the Fastly API.
type AzureBlob struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Container         string     `mapstructure:"container"`
	AccountName       string     `mapstructure:"account_name"`
	SASToken          string     `mapstructure:"sas_token"`
	Path              string     `mapstructure:"path"`
	Period            uint       `mapstructure:"period"`
	GzipLevel         uint8      `mapstructure:"gzip_level"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	MessageType       string     `mapstructure:"message_type"`
	ResponseCondition string     `mapstructure:"response_condition"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// azureblobsByName is a sortable list of azureblobs.
type azureblobsByName []*AzureBlob

// Len, Swap, and Less implement the sortable interface.
func (s azureblobsByName) Len() int      { return len(s) }
func (s azureblobsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s azureblobsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListAzureBlobsInput is used as input to the ListAzureBlobs function.
type ListAzureBlobsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListAzureBlobs returns the list of azureblobs for the configuration version.
func (c *Client) ListAzureBlobs(i *ListAzureBlobsInput) ([]*AzureBlob, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var azureblobs []*AzureBlob
	if err := decodeBodyMap(resp.Body, &azureblobs); err != nil {
		return nil, err
	}
	sort.Stable(azureblobsByName(azureblobs))
	return azureblobs, nil
}

// CreateAzureBlobInput is used as input to the CreateAzureBlob function.
type CreateAzureBlobInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	Name              string `form:"name,omitempty"`
	Container         string `form:"container"`
	AccountName       string `form:"account_name"`
	SASToken          string `form:"sas_token"`
	Path              string `form:"path,omitempty"`
	Period            uint   `form:"period,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	GzipLevel         uint8  `form:"gzip_level,omitempty"`
	Format            string `form:"format,omitempty"`
	MessageType       string `form:"message_type,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	TimestampFormat   string `form:"timestamp_format,omitempty"`
	Placement         string `form:"placement,omitempty"`
}

// CreateAzureBlob creates a new Fastly AzureBlob.
func (c *Client) CreateAzureBlob(i *CreateAzureBlobInput) (*AzureBlob, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var azureblob *AzureBlob
	if err := decodeBodyMap(resp.Body, &azureblob); err != nil {
		return nil, err
	}
	return azureblob, nil
}

// GetAzureBlobInput is used as input to the GetAzureBlob function.
type GetAzureBlobInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the AzureBlob to fetch.
	Name string
}

// GetAzureBlob gets the AzureBlob configuration with the given parameters.
func (c *Client) GetAzureBlob(i *GetAzureBlobInput) (*AzureBlob, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var b *AzureBlob
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateAzureBlobInput is used as input to the UpdateAzureBlob function.
type UpdateAzureBlobInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the AzureBlob to update.
	Name string

	NewName           string `form:"name,omitempty"`
	Container         string `form:"container"`
	AccountName       string `form:"account_name"`
	SASToken          string `form:"sas_token"`
	Path              string `form:"path,omitempty"`
	Period            uint   `form:"period,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	GzipLevel         uint8  `form:"gzip_level,omitempty"`
	Format            string `form:"format,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	TimestampFormat   string `form:"timestamp_format,omitempty"`
	Placement         string `form:"placement,omitempty"`
}

// UpdateAzureBlob updates a specific AzureBlob.
func (c *Client) UpdateAzureBlob(i *UpdateAzureBlobInput) (*AzureBlob, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *AzureBlob
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteAzureBlobInput is the input parameter to DeleteAzureBlob.
type DeleteAzureBlobInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the AzureBlob to delete (required).
	Name string
}

// DeleteAzureBlob deletes the given AzureBlob version.
func (c *Client) DeleteAzureBlob(i *DeleteAzureBlobInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf("Not Ok")
	}
	return nil
}
