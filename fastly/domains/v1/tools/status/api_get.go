package status

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Scope determines the depth of availability checking.
type Scope string

const (
	// ScopeEstimate checks DNS and aftermarket level availability.
	ScopeEstimate Scope = "estimate"
)

// GetInput specifies the parameters for a domain status check request.
type GetInput struct {
	// Domain is the domain name being checked for availability.
	Domain string
	// Scope determines the availability check to perform.
	// Scope defaults to a precise status check, specify ScopeEstimate for an estimated check.
	Scope *Scope
}

// Get performs a domain status check for a given domain.
func Get(ctx context.Context, c *fastly.Client, g *GetInput) (*Status, error) {
	if g.Domain == "" {
		return nil, fastly.ErrMissingDomain
	}

	ro := fastly.CreateRequestOptions()
	ro.Params["domain"] = g.Domain

	if g.Scope != nil {
		ro.Params["scope"] = string(*g.Scope)
	}

	path := fastly.ToSafeURL("domains", "v1", "tools", "status")
	resp, err := c.Get(ctx, path, ro)
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
