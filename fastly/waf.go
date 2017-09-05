package fastly

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/jsonapi"
)

// WAFConfigurationSet represents information about a configuration_set.
type WAFConfigurationSet struct {
	ID string `jsonapi:"primary,configuration_set"`
}

// WAF is the information about a firewall object.
type WAF struct {
	ID                string     `jsonapi:"primary,waf"`
	Version           int        `jsonapi:"attr,version"`
	PrefetchCondition string     `jsonapi:"attr,prefetch_condition"`
	Response          string     `jsonapi:"attr,response"`
	LastPush          *time.Time `jsonapi:"attr,last_push"`

	ConfigurationSet *WAFConfigurationSet `jsonapi:"relation,configuration_set"`
}

// wafType is used for reflection because JSONAPI wants to know what it's
// decoding into.
var wafType = reflect.TypeOf(new(WAF))

// ListWAFsInput is used as input to the ListWAFs function.
type ListWAFsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListWAFs returns the list of wafs for the configuration version.
func (c *Client) ListWAFs(i *ListWAFsInput) ([]*WAF, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, wafType)
	if err != nil {
		return nil, err
	}

	wafs := make([]*WAF, len(data))
	for i := range data {
		typed, ok := data[i].(*WAF)
		if !ok {
			return nil, fmt.Errorf("got back a non-WAF response")
		}
		wafs[i] = typed
	}
	return wafs, nil
}

// CreateWAFInput is used as input to the CreateWAF function.
type CreateWAFInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	ID                string `jsonapi:"primary,waf"`
	PrefetchCondition string `jsonapi:"attr,prefetch_condition,omitempty"`
	Response          string `jsonapi:"attr,response,omitempty"`
}

// CreateWAF creates a new Fastly WAF.
func (c *Client) CreateWAF(i *CreateWAFInput) (*WAF, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs", i.Service, i.Version)
	resp, err := c.PostJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// GetWAFInput is used as input to the GetWAF function.
type GetWAFInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// ID is the id of the WAF to get.
	ID string
}

func (c *Client) GetWAF(i *GetWAFInput) (*WAF, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs/%s", i.Service, i.Version, i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// UpdateWAFInput is used as input to the UpdateWAF function.
type UpdateWAFInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	ID                string `jsonapi:"primary,waf"`
	PrefetchCondition string `jsonapi:"attr,prefetch_condition,omitempty"`
	Response          string `jsonapi:"attr,response,omitempty"`
}

// UpdateWAF updates a specific WAF.
func (c *Client) UpdateWAF(i *UpdateWAFInput) (*WAF, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.ID == "" {
		return nil, ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs/%s", i.Service, i.Version, i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}

	var waf WAF
	if err := jsonapi.UnmarshalPayload(resp.Body, &waf); err != nil {
		return nil, err
	}
	return &waf, nil
}

// DeleteWAFInput is used as input to the GetWAF function.
type DeleteWAFInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// ID is the id of the WAF to delete.
	ID string
}

func (c *Client) DeleteWAF(i *DeleteWAFInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.ID == "" {
		return ErrMissingWAFID
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs/%s", i.Service, i.Version, i.ID)
	_, err := c.Delete(path, nil)
	return err
}
