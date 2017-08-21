package fastly

import "fmt"

// SetData represents information about a configuration_set.
type SetData struct {
	Type string `mapstructure:"type"`
	ID   string `mapstructure:"id"`
}

// FirewallConfigurationSet is the data inside of a configuration_set relationship.
type FirewallConfigurationSet struct {
	Data SetData `mapstructure:"data"`
}

// FirewallRelationships holds the relationships of a firewall object (currently just `configuration_set`).
type FirewallRelationships struct {
	ConfigurationSet FirewallConfigurationSet `mapstructure:"configuration_set"`
}

// FirewallAttributes holds the attributes of a firewall object.
type FirewallAttributes struct {
	PrefetchCondition string `mapstructure:"prefetch_condition"`
	Response          string `mapstructure:"response"`
	LastPush          string `mapstructure:"last_push"`
	Version           string `mapstructure:"version"`
}

// FirewallData is the information about a firewall object.
type FirewallData struct {
	ID            string                `mapstructure:"id"`
	Type          string                `mapstructure:"type"`
	Attributes    FirewallAttributes    `mapstructure:"attributes"`
	Relationships FirewallRelationships `mapstructure:"relationships"`
}

// FirewallLinks specifies first, last, and previous (if applicable) page links.
type FirewallLinks struct {
	Last     string `mapstructure:"last"`
	First    string `mapstructure:"first"`
	Previous string `mapstructure:"previous"`
}

// FirewallObjects represents a response from Fastly's API for listing firewall objects.
type FirewallObjects struct {
	// Data lists all of the firewall objects.
	Data []FirewallData `mapstructure:"data"`

	// Included lists relationships.
	Included []FirewallData `mapstructure:"included"`

	// Lists links in relation to your current page.
	Links FirewallLinks `mapstructure:"links"`
}

// GetFirewallObjectsInput is the input needed to list firewall objects.
type GetFirewallObjectsInput struct {
	// Service is the ID of the service (required).
	Service string `form:"field,omitempty"`

	// Version is the specific configuration version (required).
	Version int `form:"field,omitempty"`
}

// GetFirewallObjects lists all firewall objects for a service and version.
func (c *Client) GetFirewallObjects(i *GetFirewallObjectsInput) (*FirewallObjects, error) {
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

	var fo *FirewallObjects
	if err := decodeJSON(&fo, resp.Body); err != nil {
		return nil, err
	}

	return fo, nil
}

// GetFirewallObjectInput is the input needed to get a firewall object.
type GetFirewallObjectInput struct {
	// Service is the ID of the service (required).
	Service string `form:"field,omitempty"`

	// Version is the specific configuration version (required).
	Version int `form:"field,omitempty"`

	// WafID is the sepcific WAF ID for the firewall object you want to get (required).
	WafID string `form:"field,omitempty"`
}

// FirewallObject represents a firewall object returned from Fastly's API.
type FirewallObject struct {
	// Data contains the requested firewall object.
	Data FirewallData `mapstructure:"data"`
}

// GetFirewallObject returns a specific firewall object.
func (c *Client) GetFirewallObject(i *GetFirewallObjectInput) (*FirewallObject, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.WafID == "" {
		return nil, ErrMissingWafID
	}

	path := fmt.Sprintf("/service/%s/version/%d/wafs/%s", i.Service, i.Version, i.WafID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var fo *FirewallObject
	if err := decodeJSON(&fo, resp.Body); err != nil {
		return nil, err
	}

	return fo, nil
}
