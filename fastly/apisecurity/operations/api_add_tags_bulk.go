package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fastly/go-fastly/v13/fastly"
)

// BulkAddTagsInput specifies the information needed for the BulkAddTags() function
// to perform the operation.
type BulkAddTagsInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string `json:"-"`
	// OperationIDs is the list of operation IDs to add tags to (required).
	OperationIDs []string `json:"operation_ids"`
	// TagIDs is the list of tag IDs to add (required).
	TagIDs []string `json:"tag_ids"`
}

// BulkAddTags adds tags to multiple operations in a single request.
//
// The API returns HTTP 207 Multi-Status with a per-item result list.
func BulkAddTags(ctx context.Context, c *fastly.Client, i *BulkAddTagsInput) (*BulkOperationResultsResponse, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if len(i.OperationIDs) == 0 {
		return nil, fastly.NewFieldError("OperationIDs")
	}
	if len(i.TagIDs) == 0 {
		return nil, fastly.NewFieldError("TagIDs")
	}

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "operations-bulk-tags")

	body := &OperationBulkAddTags{
		OperationIDs: i.OperationIDs,
		TagIDs:       i.TagIDs,
	}

	resp, err := c.PostJSON(ctx, path, body, fastly.CreateRequestOptions())
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
