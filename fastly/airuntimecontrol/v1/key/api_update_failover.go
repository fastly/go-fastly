package key

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateFailoverInput specifies the information needed for the
// UpdateFailover() function to perform the operation.
type UpdateFailoverInput struct {
	// KeyID is the ID of the virtual key (required).
	KeyID *string `json:"-"`
	// Targets is the ordered list of provider connections requests are sent to.
	// The first entry is the primary connection; subsequent entries are
	// fallbacks, tried in order.
	Targets []FailoverTarget `json:"targets,omitempty"`
	// TimeoutEnabled is whether timeouts trigger failover.
	TimeoutEnabled *bool `json:"timeout_enabled,omitempty"`
	// RateLimit429Enabled is whether 429 rate-limit responses trigger failover.
	RateLimit429Enabled *bool `json:"rate_limit_429_enabled,omitempty"`
	// ServerError5xxEnabled is whether 5xx server errors trigger failover.
	ServerError5xxEnabled *bool `json:"server_error_5xx_enabled,omitempty"`
}

// UpdateFailover sets the failover policy for a virtual key.
func UpdateFailover(ctx context.Context, c *fastly.Client, i *UpdateFailoverInput) (*Failover, error) {
	if i.KeyID == nil {
		return nil, fastly.ErrMissingKeyID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "keys", *i.KeyID, "failover")

	resp, err := c.PutJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var f *Failover
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return f, nil
}
