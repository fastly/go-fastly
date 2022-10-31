package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Scalyr represents a scalyr response from the Fastly API.
type Scalyr struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	Name              string     `mapstructure:"name"`
	Placement         string     `mapstructure:"placement"`
	Region            string     `mapstructure:"region"`
	ResponseCondition string     `mapstructure:"response_condition"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	Token             string     `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// scalyrByName is a sortable list of scalyrs.
type scalyrsByName []*Scalyr

// Len, Swap, and Less implement the sortable interface.
func (s scalyrsByName) Len() int      { return len(s) }
func (s scalyrsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s scalyrsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListScalyrsInput is used as input to the ListScalyrs function.
type ListScalyrsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListScalyrs returns the list of scalyrs for the configuration version.
func (c *Client) ListScalyrs(i *ListScalyrsInput) ([]*Scalyr, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ss []*Scalyr
	if err := decodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	sort.Stable(scalyrsByName(ss))
	return ss, nil
}

// CreateScalyrInput is used as input to the CreateScalyr function.
type CreateScalyrInput struct {
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	Name              string `url:"name,omitempty"`
	Placement         string `url:"placement,omitempty"`
	Region            string `url:"region,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Token          string `url:"token,omitempty"`
}

// CreateScalyr creates a new Fastly scalyr.
func (c *Client) CreateScalyr(i *CreateScalyrInput) (*Scalyr, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Token == "" {
		return nil, ErrMissingToken
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Scalyr
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetScalyrInput is used as input to the GetScalyr function.
type GetScalyrInput struct {
	// Name is the name of the scalyr to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetScalyr gets the scalyr configuration with the given parameters.
func (c *Client) GetScalyr(i *GetScalyrInput) (*Scalyr, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Scalyr
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateScalyrInput is used as input to the UpdateScalyr function.
type UpdateScalyrInput struct {
	Format        *string `url:"format,omitempty"`
	FormatVersion *uint   `url:"format_version,omitempty"`
	// Name is the name of the scalyr to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	Region            *string `url:"region,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Token          *string `url:"token,omitempty"`
}

// UpdateScalyr updates a specific scalyr.
func (c *Client) UpdateScalyr(i *UpdateScalyrInput) (*Scalyr, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	if i.Token != nil && *i.Token == "" {
		return nil, ErrTokenEmpty
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Scalyr
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteScalyrInput is the input parameter to DeleteScalyr.
type DeleteScalyrInput struct {
	// Name is the name of the scalyr to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteScalyr deletes the given scalyr version.
func (c *Client) DeleteScalyr(i *DeleteScalyrInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/scalyr/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
