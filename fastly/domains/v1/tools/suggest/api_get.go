package suggest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// GetInput specifies the various parameters for performing real-time queries against the known zones database.
type GetInput struct {
	// Query are the term(s) to search against.
	Query string
	// Defaults is a comma-separated list of default zones to include in the search results response.
	Defaults *string
	// Keywords is a comma-separated list of keywords for seeding the results.
	// Example: a new gTLD like `kitchen`, or a related keyword like `vegan`.
	// Helpful for search result relevance (e.g. from a targeted ad click, a user profile, etc.)
	Keywords *string
	// Location overrides the IP location detection for country-code zones, with a two-character country code.
	Location *string
	// Vendor is the domain name of a specific registrar or vendor.
	// Vendor is used to filter results by the zones supported by the vendor.
	Vendor *string
}

// Get returns a list of domain suggestions matching the query criteria.
func Get(ctx context.Context, c *fastly.Client, g *GetInput) (*Suggestions, error) {
	if g.Query == "" {
		return nil, fastly.ErrMissingDomainQuery
	}

	ro := fastly.CreateRequestOptions()
	ro.Params["query"] = g.Query

	if g.Defaults != nil {
		ro.Params["defaults"] = *g.Defaults
	}

	if g.Keywords != nil {
		ro.Params["keywords"] = *g.Keywords
	}

	if g.Location != nil {
		ro.Params["location"] = *g.Location
	}

	if g.Vendor != nil {
		ro.Params["vendor"] = *g.Vendor
	}

	path := fastly.ToSafeURL("domains", "v1", "tools", "suggest")
	resp, err := c.Get(ctx, path, ro)
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
