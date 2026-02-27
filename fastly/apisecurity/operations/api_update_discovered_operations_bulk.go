package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fastly/go-fastly/v13/fastly"
)

// BulkUpdateDiscoveredStatusInput specifies the information needed for the
// BulkUpdateDiscoveredStatus() function to perform the operation.
type BulkUpdateDiscoveredStatusInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string `json:"-"`
	// OperationIDs is the list of discovered operation IDs to update (required).
	OperationIDs []string `json:"operation_ids"`
	// Status is the new status to apply (required).
	Status *string `json:"status"`
}

// BulkUpdateDiscoveredStatus updates the status of multiple discovered operations in a single request.
func BulkUpdateDiscoveredStatus(ctx context.Context, c *fastly.Client, i *BulkUpdateDiscoveredStatusInput) (*BulkOperationResultsResponse, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.Status == nil || *i.Status == "" {
		return nil, fastly.ErrMissingStatus
	}
	if len(i.OperationIDs) == 0 {
		return nil, fastly.NewFieldError("OperationIDs")
	}

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "discovered-operations-bulk")

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus {
		return nil, fastly.NewHTTPError(resp)
	}

	var out *BulkOperationResultsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return out, nil
}
