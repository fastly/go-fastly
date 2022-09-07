package fastly

import "fmt"

// Diff represents a diff of two versions as a response from the Fastly API.
type Diff struct {
	Format string `mapstructure:"format"`
	From   int    `mapstructure:"from"`
	To     int    `mapstructure:"to"`
	Diff   string `mapstructure:"diff"`
}

// GetDiffInput is used as input to the GetDiff function.
type GetDiffInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// From is the version to diff from. This can either be a string indicating a
	// positive number (e.g. "1") or a negative number from "-1" down ("-1" is the
	// latest version).
	From int

	// To is the version to diff up to. The same rules for From apply.
	To int

	// Format is an optional field to specify the format with which the diff will
	// be returned. Acceptable values are "text" (default), "html", or
	// "html_simple".
	Format string
}

// GetDiff returns the diff of the given versions.
func (c *Client) GetDiff(i *GetDiffInput) (*Diff, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.From == 0 {
		return nil, ErrMissingFrom
	}

	if i.To == 0 {
		return nil, ErrMissingTo
	}

	path := fmt.Sprintf("service/%s/diff/from/%d/to/%d", i.ServiceID, i.From, i.To)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Diff
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}
