package status

import (
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Scope determines the depth of availability checking.
type Scope string

const (
	// ScopePrecise checks domain registry-level availability.
	ScopePrecise Scope = "precise"
	// ScopeEstimate checks DNS and aftermarket level availability.
	ScopeEstimate Scope = "estimate"
)

// Input specifies the parameters for a domain status check request.
type Input struct {
	// Domain is the domain name being checked for availability.
	Domain string
	// Scope determines the availability check to perform, defaulting to precise (optional).
	Scope *Scope
}

// Get performs a domain status check for a given domain.
func Get(c *fastly.Client, i *Input) (*Status, error) {
	if i.Domain == "" {
		return nil, fastly.ErrMissingDomain
	}

	ro := &fastly.RequestOptions{
		Params: map[string]string{
			"domain": i.Domain,
		},
	}

	if i.Scope != nil {
		ro.Params["scope"] = string(*i.Scope)
	}

	path := fastly.ToSafeURL("domains", "v1", "tools", "status")
	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status *Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return status, err
}
