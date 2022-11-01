package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Condition represents a condition response from the Fastly API.
type Condition struct {
	Comment        string     `mapstructure:"comment"`
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	Name           string     `mapstructure:"name"`
	Priority       int        `mapstructure:"priority"`
	ServiceID      string     `mapstructure:"service_id"`
	ServiceVersion int        `mapstructure:"version"`
	Statement      string     `mapstructure:"statement"`
	Type           string     `mapstructure:"type"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}

// conditionsByName is a sortable list of conditions.
type conditionsByName []*Condition

// Len implement the sortable interface.
func (s conditionsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s conditionsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s conditionsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListConditionsInput is used as input to the ListConditions function.
type ListConditionsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListConditions retrieves all resources.
func (c *Client) ListConditions(i *ListConditionsInput) ([]*Condition, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/condition", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs []*Condition
	if err := decodeBodyMap(resp.Body, &cs); err != nil {
		return nil, err
	}
	sort.Stable(conditionsByName(cs))
	return cs, nil
}

// CreateConditionInput is used as input to the CreateCondition function.
type CreateConditionInput struct {
	Name     string `url:"name,omitempty"`
	Priority *int   `url:"priority,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Statement      string `url:"statement,omitempty"`
	Type           string `url:"type,omitempty"`
}

// CreateCondition creates a new resource.
func (c *Client) CreateCondition(i *CreateConditionInput) (*Condition, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/condition", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var co *Condition
	if err := decodeBodyMap(resp.Body, &co); err != nil {
		return nil, err
	}
	return co, nil
}

// GetConditionInput is used as input to the GetCondition function.
type GetConditionInput struct {
	// Name is the name of the condition to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetCondition gets the condition configuration with the given parameters.
func (c *Client) GetCondition(i *GetConditionInput) (*Condition, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/condition/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var co *Condition
	if err := decodeBodyMap(resp.Body, &co); err != nil {
		return nil, err
	}
	return co, nil
}

// UpdateConditionInput is used as input to the UpdateCondition function.
type UpdateConditionInput struct {
	Comment *string `url:"comment,omitempty"`
	// Name is the name of the condition to update.
	Name     string
	Priority *int `url:"priority,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Statement      *string `url:"statement,omitempty"`
	Type           *string `url:"type,omitempty"`
}

// UpdateCondition updates the specified resource.
func (c *Client) UpdateCondition(i *UpdateConditionInput) (*Condition, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/condition/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var co *Condition
	if err := decodeBodyMap(resp.Body, &co); err != nil {
		return nil, err
	}
	return co, nil
}

// DeleteConditionInput is the input parameter to DeleteCondition.
type DeleteConditionInput struct {
	// Name is the name of the condition to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteCondition deletes the specified resource.
func (c *Client) DeleteCondition(i *DeleteConditionInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/condition/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
