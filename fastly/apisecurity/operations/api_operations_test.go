package operations

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_Operations(t *testing.T) {
	ctx := context.TODO()

	serviceID := fastly.TestDeliveryServiceID

	var err error

	// Create a tag to associate with operations.
	const tagName = "test-tag-operations"
	var tag *OperationTag
	fastly.Record(t, "create_tag_for_operation", func(c *fastly.Client) {
		tag, err = CreateTag(ctx, c, &CreateTagInput{
			ServiceID:   new(serviceID),
			Name:        new(tagName),
			Description: new("go-fastly test tag"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, tag)
	require.NotEmpty(t, tag.ID)
	require.Equal(t, tagName, tag.Name)

	defer func() {
		fastly.Record(t, "delete_tag_for_operation", func(c *fastly.Client) {
			_ = DeleteTag(ctx, c, &DeleteTagInput{
				ServiceID: new(serviceID),
				TagID:     new(tag.ID),
			})
		})
	}()

	// Create operation #1.
	var op1 *Operation
	fastly.Record(t, "create_operation", func(c *fastly.Client) {
		op1, err = Create(ctx, c, &CreateInput{
			ServiceID: new(serviceID),
			Method:    new("GET"),
			Domain:    new("example.com"),
			Path:      new("/test"),
			TagIDs:    []string{tag.ID},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, op1)
	require.NotEmpty(t, op1.ID)

	defer func() {
		fastly.Record(t, "delete_operation", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ServiceID:   new(serviceID),
				OperationID: new(op1.ID),
			})
		})
	}()

	// Create operation #2 (for pagination test).
	var op2 *Operation
	fastly.Record(t, "create_operation_2", func(c *fastly.Client) {
		op2, err = Create(ctx, c, &CreateInput{
			ServiceID: new(serviceID),
			Method:    new("GET"),
			Domain:    new("example.com"),
			Path:      new("/test-pagination"),
			TagIDs:    []string{tag.ID},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, op2)
	require.NotEmpty(t, op2.ID)

	defer func() {
		fastly.Record(t, "delete_operation_2", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ServiceID:   new(serviceID),
				OperationID: new(op2.ID),
			})
		})
	}()

	// Describe operation.
	var described *Operation
	fastly.Record(t, "describe_operation", func(c *fastly.Client) {
		described, err = Describe(ctx, c, &DescribeInput{
			ServiceID:   new(serviceID),
			OperationID: new(op1.ID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, described)
	require.Equal(t, op1.ID, described.ID)

	// Update operation.
	var updated *Operation
	fastly.Record(t, "update_operation", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			ServiceID:   new(serviceID),
			OperationID: new(op1.ID),
			Description: new("updated"),
			TagIDs:      []string{tag.ID},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, op1.ID, updated.ID)
	require.Equal(t, "updated", updated.Description)

	// List operations (with filters) - existing behavior.
	var ops *Operations
	fastly.Record(t, "list_operations", func(c *fastly.Client) {
		ops, err = ListOperations(ctx, c, &ListOperationsInput{
			ServiceID: new(serviceID),
			Method:    []string{"GET"},
			Domain:    []string{"example.com"},
			Path:      new("/test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, ops)

	found := false
	for _, item := range ops.Data {
		if item.ID == op1.ID {
			found = true
			break
		}
	}
	require.True(t, found, "expected created operation to appear in filtered list")

	// ---- Pagination test for operations ----
	limit := 1
	p := NewOperationPaginator(ctx, fastly.TestClient, &ListOperationsInput{
		ServiceID: new(serviceID),
		Method:    []string{"GET"},
		Domain:    []string{"example.com"},
		Limit:     &limit,
		Page:      new(0),
	})

	var collected []Operation
	for i := 0; i < 5 && p.HasNext(); i++ {
		recName := fmt.Sprintf("list_operations_page_%d", i)

		var pageData []Operation
		fastly.Record(t, recName, func(c *fastly.Client) {
			p.SetClient(c)
			pageData, err = p.GetNext()
		})
		require.NoError(t, err)
		collected = append(collected, pageData...)
	}

	seen1 := false
	seen2 := false
	for _, it := range collected {
		if it.ID == op1.ID {
			seen1 = true
		}
		if it.ID == op2.ID {
			seen2 = true
		}
	}
	require.True(t, seen1, "expected to paginate and find op1")
	require.True(t, seen2, "expected to paginate and find op2")

	// ---- Discovered operations list + status update endpoints ----
	var discoveredFiltered *DiscoveredOperations
	fastly.Record(t, "list_discovered_operations", func(c *fastly.Client) {
		discoveredFiltered, err = ListDiscovered(ctx, c, &ListDiscoveredInput{
			ServiceID: new(serviceID),
			Method:    []string{"GET"},
			Domain:    []string{"example.com"},
			Path:      new("/test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, discoveredFiltered)

	var discoveredAny *DiscoveredOperations
	fastly.Record(t, "list_discovered_operations_any", func(c *fastly.Client) {
		discoveredAny, err = ListDiscovered(ctx, c, &ListDiscoveredInput{
			ServiceID: new(serviceID),
			Limit:     new(1),
			Page:      new(0),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, discoveredAny)

	if len(discoveredAny.Data) > 0 && discoveredAny.Data[0].ID != "" {
		discoveredID := discoveredAny.Data[0].ID

		var singleUpdated *DiscoveredOperation
		fastly.Record(t, "update_discovered_operation_status", func(c *fastly.Client) {
			singleUpdated, err = UpdateDiscoveredStatus(ctx, c, &UpdateDiscoveredStatusInput{
				ServiceID:   new(serviceID),
				OperationID: new(discoveredID),
				Status:      new("IGNORED"),
			})
		})
		require.NoError(t, err)
		require.NotNil(t, singleUpdated)

		var bulkUpdatedDiscovered *BulkOperationResultsResponse
		fastly.Record(t, "bulk_update_discovered_operation_status", func(c *fastly.Client) {
			bulkUpdatedDiscovered, err = BulkUpdateDiscoveredStatus(ctx, c, &BulkUpdateDiscoveredStatusInput{
				ServiceID:    new(serviceID),
				OperationIDs: []string{discoveredID},
				Status:       new("IGNORED"),
			})
		})
		require.NoError(t, err)
		require.NotNil(t, bulkUpdatedDiscovered)
	}

	// Bulk create ops + bulk tags.
	var bulkCreated *BulkCreateOperationsResponse
	fastly.Record(t, "bulk_create_operations", func(c *fastly.Client) {
		bulkCreated, err = BulkCreateOperations(ctx, c, &BulkCreateOperationsInput{
			ServiceID: new(serviceID),
			Operations: []OperationBulkCreateItem{
				{
					Method:      new("GET"),
					Domain:      new("example.com"),
					Path:        new("/bulk-test-1"),
					Description: new("bulk test 1"),
				},
			},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, bulkCreated)

	var createdIDs []string
	for _, r := range bulkCreated.Data {
		if r.Operation != nil && r.Operation.ID != "" {
			createdIDs = append(createdIDs, r.Operation.ID)
		}
	}
	require.NotEmpty(t, createdIDs)

	defer func() {
		if len(createdIDs) == 0 {
			return
		}
		fastly.Record(t, "delete_bulk_operations", func(c *fastly.Client) {
			for _, opID := range createdIDs {
				_ = Delete(ctx, c, &DeleteInput{
					ServiceID:   new(serviceID),
					OperationID: new(opID),
				})
			}
		})
	}()

	var bulkTagged *BulkOperationResultsResponse
	fastly.Record(t, "bulk_add_tags_to_operations", func(c *fastly.Client) {
		bulkTagged, err = BulkAddTags(ctx, c, &BulkAddTagsInput{
			ServiceID:    new(serviceID),
			OperationIDs: createdIDs,
			TagIDs:       []string{tag.ID},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, bulkTagged)
}
