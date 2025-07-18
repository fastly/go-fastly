package fastly

import "context"

// Coordinates represent the location of a datacenter.
type Coordinates struct {
	Latitude  *float64 `mapstructure:"latitude"`
	Longitude *float64 `mapstructure:"longitude"`
	X         *float64 `mapstructure:"x"`
	Y         *float64 `mapstructure:"y"`
}

// Datacenter is a list of Datacenters returned by the Fastly API.
type Datacenter struct {
	Code        *string      `mapstructure:"code"`
	Coordinates *Coordinates `mapstructure:"coordinates"`
	Group       *string      `mapstructure:"group"`
	Name        *string      `mapstructure:"name"`
	Shield      *string      `mapstructure:"shield"`
}

// AllDatacenters returns the lists of datacenters for Fastly's network.
func (c *Client) AllDatacenters(ctx context.Context) (datacenters []Datacenter, err error) {
	resp, err := c.Get(ctx, "/datacenters", CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m []Datacenter
	if err := DecodeBodyMap(resp.Body, &m); err != nil {
		return nil, err
	}

	return m, nil
}
