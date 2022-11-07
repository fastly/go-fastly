package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Kinesis represents a Kinesis response from the Fastly API.
type Kinesis struct {
	AccessKey         string     `mapstructure:"access_key"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     int        `mapstructure:"format_version"`
	IAMRole           string     `mapstructure:"iam_role"`
	Name              string     `mapstructure:"name"`
	Placement         string     `mapstructure:"placement"`
	Region            string     `mapstructure:"region"`
	ResponseCondition string     `mapstructure:"response_condition"`
	SecretKey         string     `mapstructure:"secret_key"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	StreamName        string     `mapstructure:"topic"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// kinesisByName is a sortable list of Kinesis.
type kinesisByName []*Kinesis

// Len implement the sortable interface.
func (s kinesisByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s kinesisByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s kinesisByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListKinesisInput is used as input to the ListKinesis function.
type ListKinesisInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListKinesis retrieves all resources.
func (c *Client) ListKinesis(i *ListKinesisInput) ([]*Kinesis, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kinesis", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var kineses []*Kinesis
	if err := decodeBodyMap(resp.Body, &kineses); err != nil {
		return nil, err
	}
	sort.Stable(kinesisByName(kineses))
	return kineses, nil
}

// CreateKinesisInput is used as input to the CreateKinesis function.
type CreateKinesisInput struct {
	// AccessKey is the access key associated with the target Amazon Kinesis stream. Not required if iam_role is specified.
	AccessKey *string `url:"access_key,omitempty"`
	// Format is a Fastly log format string. Must produce valid JSON that Kinesis can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// IAMRole is the ARN for an IAM role granting Fastly access to the target Amazon Kinesis stream.
	IAMRole *string `url:"iam_role,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Region is a named set of AWS resources that's in the same geographical area.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is the secret key associated with the target Amazon Kinesis stream. Not required if iam_role is specified.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// StreamName is the Amazon Kinesis stream to send logs to.
	StreamName *string `url:"topic,omitempty"`
}

// CreateKinesis creates a new resource.
func (c *Client) CreateKinesis(i *CreateKinesisInput) (*Kinesis, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kinesis", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var kinesis *Kinesis
	if err := decodeBodyMap(resp.Body, &kinesis); err != nil {
		return nil, err
	}
	return kinesis, nil
}

// GetKinesisInput is used as input to the GetKinesis function.
type GetKinesisInput struct {
	// Name is the name of the Kinesis logging object to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetKinesis retrieves the specified resource.
func (c *Client) GetKinesis(i *GetKinesisInput) (*Kinesis, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kinesis/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var kinesis *Kinesis
	if err := decodeBodyMap(resp.Body, &kinesis); err != nil {
		return nil, err
	}
	return kinesis, nil
}

// UpdateKinesisInput is used as input to the UpdateKinesis function.
type UpdateKinesisInput struct {
	// AccessKey is the access key associated with the target Amazon Kinesis stream. Not required if iam_role is specified.
	AccessKey *string `url:"access_key,omitempty"`
	// Format is a Fastly log format string. Must produce valid JSON that Kinesis can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// IAMRole is the ARN for an IAM role granting Fastly access to the target Amazon Kinesis stream.
	IAMRole *string `url:"iam_role,omitempty"`
	// Name is the name of the Kinesis logging object to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Region is a named set of AWS resources that's in the same geographical area.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is the secret key associated with the target Amazon Kinesis stream. Not required if iam_role is specified.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// StreamName is the Amazon Kinesis stream to send logs to.
	StreamName *string `url:"topic,omitempty"`
}

// UpdateKinesis updates the specified resource.
func (c *Client) UpdateKinesis(i *UpdateKinesisInput) (*Kinesis, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kinesis/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var kinesis *Kinesis
	if err := decodeBodyMap(resp.Body, &kinesis); err != nil {
		return nil, err
	}
	return kinesis, nil
}

// DeleteKinesisInput is the input parameter to DeleteKinesis.
type DeleteKinesisInput struct {
	// Name is the name of the Kinesis logging object to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteKinesis deletes the specified resource.
func (c *Client) DeleteKinesis(i *DeleteKinesisInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kinesis/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
