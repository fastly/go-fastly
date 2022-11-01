package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Sumologic represents a sumologic response from the Fastly API.
type Sumologic struct {
	Address           string     `mapstructure:"address"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     int        `mapstructure:"format_version"`
	MessageType       string     `mapstructure:"message_type"`
	Name              string     `mapstructure:"name"`
	Placement         string     `mapstructure:"placement"`
	ResponseCondition string     `mapstructure:"response_condition"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	URL               string     `mapstructure:"url"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// sumologicsByName is a sortable list of sumologics.
type sumologicsByName []*Sumologic

// Len implement the sortable interface.
func (s sumologicsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s sumologicsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s sumologicsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListSumologicsInput is used as input to the ListSumologics function.
type ListSumologicsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListSumologics returns the list of sumologics for the configuration version.
func (c *Client) ListSumologics(i *ListSumologicsInput) ([]*Sumologic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ss []*Sumologic
	if err := decodeBodyMap(resp.Body, &ss); err != nil {
		return nil, err
	}
	sort.Stable(sumologicsByName(ss))
	return ss, nil
}

// CreateSumologicInput is used as input to the CreateSumologic function.
type CreateSumologicInput struct {
	Address           string `url:"address,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     int    `url:"format_version,omitempty"`
	MessageType       string `url:"message_type,omitempty"`
	Name              string `url:"name,omitempty"`
	Placement         string `url:"placement,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	URL            string `url:"url,omitempty"`
}

// CreateSumologic creates a new resource.
func (c *Client) CreateSumologic(i *CreateSumologicInput) (*Sumologic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Sumologic
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSumologicInput is used as input to the GetSumologic function.
type GetSumologicInput struct {
	// Name is the name of the sumologic to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetSumologic gets the sumologic configuration with the given parameters.
func (c *Client) GetSumologic(i *GetSumologicInput) (*Sumologic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Sumologic
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateSumologicInput is used as input to the UpdateSumologic function.
type UpdateSumologicInput struct {
	Address       *string `url:"address,omitempty"`
	Format        *string `url:"format,omitempty"`
	FormatVersion *int    `url:"format_version,omitempty"`
	MessageType   *string `url:"message_type,omitempty"`
	// Name is the name of the sumologic to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	URL            *string `url:"url,omitempty"`
}

// UpdateSumologic updates a specific sumologic.
func (c *Client) UpdateSumologic(i *UpdateSumologicInput) (*Sumologic, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Sumologic
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSumologicInput is the input parameter to DeleteSumologic.
type DeleteSumologicInput struct {
	// Name is the name of the sumologic to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteSumologic deletes the specified resource.
func (c *Client) DeleteSumologic(i *DeleteSumologicInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/sumologic/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
