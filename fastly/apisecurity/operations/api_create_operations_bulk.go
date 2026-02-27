package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fastly/go-fastly/v13/fastly"
)

// BulkCreateOperationsInput specifies the information needed for the
// BulkCreateOperations() function to perform the operation.
type BulkCreateOperationsInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string `json:"-"`
	// Operations is the list of operations to create (required).
	Operations []OperationBulkCreateItem `json:"operations"`
}

// BulkCreateOperations creates multiple operations in a single request.
//
// The API returns HTTP 207 Multi-Status with per-item results.
func BulkCreateOperations(ctx context.Context, c *fastly.Client, i *BulkCreateOperationsInput) (*BulkCreateOperationsResponse, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if len(i.Operations) == 0 {
		return nil, fastly.NewFieldError("Operations")
	}

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "operations-bulk")

	body := &OperationBulkCreate{Operations: i.Operations}

	resp, err := c.PostJSON(ctx, path, body, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus {
		return nil, fastly.NewHTTPError(resp)
	}

	var out *BulkCreateOperationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return out, nil
}
