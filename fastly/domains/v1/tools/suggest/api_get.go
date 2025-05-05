package suggest

import (
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Input specifies the various search terms to perform real time queries
type Input struct {
	// Query are the Term(s) to search against
	Query string
	// Defaults is a comma-separated list of default zones to include in the search results response (optional).
	Defaults *string
	// Keywords is a comma-separated list of keywords for seeding the results (optional).
	// Example: a new gTLD like `kitchen`, or a related keyword like `vegan`.
	// Helpful for search result relevance (e.g. from a targeted ad click, a user profile, etc.)
	Keywords *string
	// Location overrides the IP location detection for country-code zones, with a two-character country code (optional).
	Location *string
	// Vendor is the domain name of a specific registrar or vendor,
	// for filtering results by the zones supported by the vendor.
	Vendor *string
}

func Get(c *fastly.Client, i *Input) (*Suggestions, error) {
	if i.Query == "" {
		return nil, fastly.ErrMissingDomainQuery
	}

	ro := &fastly.RequestOptions{
		Params: map[string]string{
			"query": i.Query,
		},
	}

	if i.Defaults != nil {
		ro.Params["defaults"] = *i.Defaults
	}

	if i.Keywords != nil {
		ro.Params["keywords"] = *i.Keywords
	}

	if i.Location != nil {
		ro.Params["location"] = *i.Location
	}

	if i.Vendor != nil {
		ro.Params["vendor"] = *i.Vendor
	}

	path := fastly.ToSafeURL("domains", "v1", "tools", "suggest")

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var suggestions *Suggestions
	if err := json.NewDecoder(resp.Body).Decode(&suggestions); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return suggestions, err
}
