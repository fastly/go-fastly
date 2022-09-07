package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Pubsub represents an Pubsub logging response from the Fastly API.
type Pubsub struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Topic             string     `mapstructure:"topic"`
	User              string     `mapstructure:"user"`
	SecretKey         string     `mapstructure:"secret_key"`
	ProjectID         string     `mapstructure:"project_id"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// pubsubsByName is a sortable list of pubsubs.
type pubsubsByName []*Pubsub

// Len, Swap, and Less implement the sortable interface.
func (s pubsubsByName) Len() int      { return len(s) }
func (s pubsubsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s pubsubsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListPubsubsInput is used as input to the ListPubsubs function.
type ListPubsubsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListPubsubs returns the list of pubsubs for the configuration version.
func (c *Client) ListPubsubs(i *ListPubsubsInput) ([]*Pubsub, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/pubsub", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pubsubs []*Pubsub
	if err := decodeBodyMap(resp.Body, &pubsubs); err != nil {
		return nil, err
	}
	sort.Stable(pubsubsByName(pubsubs))
	return pubsubs, nil
}

// CreatePubsubInput is used as input to the CreatePubsub function.
type CreatePubsubInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `url:"name,omitempty"`
	Topic             string `url:"topic,omitempty"`
	User              string `url:"user,omitempty"`
	SecretKey         string `url:"secret_key,omitempty"`
	ProjectID         string `url:"project_id,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	Format            string `url:"format,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	Placement         string `url:"placement,omitempty"`
}

// CreatePubsub creates a new Fastly Pubsub.
func (c *Client) CreatePubsub(i *CreatePubsubInput) (*Pubsub, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/pubsub", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pubsub *Pubsub
	if err := decodeBodyMap(resp.Body, &pubsub); err != nil {
		return nil, err
	}
	return pubsub, nil
}

// GetPubsubInput is used as input to the GetPubsub function.
type GetPubsubInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Pubsub to fetch.
	Name string
}

// GetPubsub gets the Pubsub configuration with the given parameters.
func (c *Client) GetPubsub(i *GetPubsubInput) (*Pubsub, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/pubsub/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Pubsub
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdatePubsubInput is used as input to the UpdatePubsub function.
type UpdatePubsubInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Pubsub to update.
	Name string

	NewName           *string `url:"name,omitempty"`
	Topic             *string `url:"topic,omitempty"`
	User              *string `url:"user,omitempty"`
	SecretKey         *string `url:"secret_key,omitempty"`
	ProjectID         *string `url:"project_id,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
	Format            *string `url:"format,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	Placement         *string `url:"placement,omitempty"`
}

// UpdatePubsub updates a specific Pubsub.
func (c *Client) UpdatePubsub(i *UpdatePubsubInput) (*Pubsub, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/pubsub/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Pubsub
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeletePubsubInput is the input parameter to DeletePubsub.
type DeletePubsubInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Pubsub to delete (required).
	Name string
}

// DeletePubsub deletes the given Pubsub version.
func (c *Client) DeletePubsub(i *DeletePubsubInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/pubsub/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
