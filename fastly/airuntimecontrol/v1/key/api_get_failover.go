package key

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetFailoverInput specifies the information needed for the GetFailover()
// function to perform the operation.
type GetFailoverInput struct {
	// KeyID is the ID of the virtual key (required).
	KeyID *string
}

// GetFailover retrieves the failover policy for a virtual key, including the
// ordered list of provider connection targets and the conditions that trigger
// failover.
func GetFailover(ctx context.Context, c *fastly.Client, i *GetFailoverInput) (*Failover, error) {
	if i.KeyID == nil {
		return nil, fastly.ErrMissingKeyID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "keys", *i.KeyID, "failover")

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
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
