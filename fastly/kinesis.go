package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Kinesis represents a Kinesis response from the Fastly API.
type Kinesis struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	StreamName        string     `mapstructure:"topic"`
	Region            string     `mapstructure:"region"`
	AccessKey         string     `mapstructure:"access_key"`
	SecretKey         string     `mapstructure:"secret_key"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// kinesesByName is a sortable list of Kineses.
type kinesesByName []*Kinesis

// Len, Swap, and Less implement the sortable interface.
func (s kinesesByName) Len() int      { return len(s) }
func (s kinesesByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s kinesesByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListKinesesInput is used as input to the ListKineses function.
type ListKinesesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListKineses returns the list of Kineses for the configuration version.
func (c *Client) ListKineses(i *ListKinesesInput) ([]*Kinesis, error) {
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

	var kineses []*Kinesis
	if err := decodeBodyMap(resp.Body, &kineses); err != nil {
		return nil, err
	}
	sort.Stable(kinesesByName(kineses))
	return kineses, nil
}

// CreateKinesisInput is used as input to the CreateKinesis function.
type CreateKinesisInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `form:"name,omitempty"`
	StreamName        string `form:"topic,omitempty"`
	Region            string `form:"region,omitempty"`
	AccessKey         string `form:"access_key,omitempty"`
	SecretKey         string `form:"secret_key,omitempty"`
	Format            string `form:"format,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	Placement         string `form:"placement,omitempty"`
}

// CreateKinesis creates a new Fastly Kinesis.
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

	var kinesis *Kinesis
	if err := decodeBodyMap(resp.Body, &kinesis); err != nil {
		return nil, err
	}
	return kinesis, nil
}

// GetKinesisInput is used as input to the GetKinesis function.
type GetKinesisInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Kinesis logging object to fetch (required).
	Name string
}

// GetKinesis gets the Kinesis configuration with the given parameters.
func (c *Client) GetKinesis(i *GetKinesisInput) (*Kinesis, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kinesis/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var kinesis *Kinesis
	if err := decodeBodyMap(resp.Body, &kinesis); err != nil {
		return nil, err
	}
	return kinesis, nil
}

// UpdateKinesisInput is used as input to the UpdateKinesis function.
type UpdateKinesisInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Kinesis logging object to update (required).
	Name string

	NewName           *string `form:"name,omitempty"`
	StreamName        *string `form:"topic,omitempty"`
	Region            *string `form:"region,omitempty"`
	AccessKey         *string `form:"access_key,omitempty"`
	SecretKey         *string `form:"secret_key,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateKinesis updates a specific Kinesis.
func (c *Client) UpdateKinesis(i *UpdateKinesisInput) (*Kinesis, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kinesis/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var kinesis *Kinesis
	if err := decodeBodyMap(resp.Body, &kinesis); err != nil {
		return nil, err
	}
	return kinesis, nil
}

// DeleteKinesisInput is the input parameter to DeleteKinesis.
type DeleteKinesisInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Kinesis logging object to delete (required).
	Name string
}

// DeleteKinesis deletes the given Kinesis version.
func (c *Client) DeleteKinesis(i *DeleteKinesisInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/kinesis/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
