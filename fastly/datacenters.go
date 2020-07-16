package fastly

// Coordinates represent the location of a datacenter.
type Coordinates struct {
	Latitude   float64 `mapstructure:"latitude"`
	Longtitude float64 `mapstructure:"longitude"`
	X          float64 `mapstructure:"x"`
	Y          float64 `mapstructure:"y"`
}

// Datacenter is a list of Datacenters returned by the Fastly API.
type Datacenter struct {
	Code        string      `mapstructure:"code"`
	Coordinates Coordinates `mapstructure:"coordinates"`
	Group       string      `mapstructure:"group"`
	Name        string      `mapstructure:"name"`
	Shield      string      `mapstructure:"shield"`
}

// AllDatacenters returns the lists of datacenters for Fastly's network.
func (c *Client) AllDatacenters() (datacenters []Datacenter, err error) {
	resp, err := c.Get("/datacenters", nil)
	if err != nil {
		return nil, err
	}
	var m []Datacenter
	if err := decodeBodyMap(resp.Body, &m); err != nil {
		return nil, err
	}

	return m, nil
}
