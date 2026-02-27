package operations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v13/fastly"
)

func TestClient_Operations(t *testing.T) {
	ctx := context.TODO()

	serviceID := fastly.TestDeliveryServiceID

	var err error

	// Create a tag to associate with the operation.
	const tagName = "test-tag-operations"
	var tag *OperationTag
	fastly.Record(t, "create_tag_for_operation", func(c *fastly.Client) {
		tag, err = CreateTag(ctx, c, &CreateTagInput{
			ServiceID:   fastly.ToPointer(serviceID),
			Name:        fastly.ToPointer(tagName),
			Description: fastly.ToPointer("go-fastly test tag"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, tag)
	require.NotEmpty(t, tag.ID)
	require.Equal(t, tagName, tag.Name)

	defer func() {
		fastly.Record(t, "delete_tag_for_operation", func(c *fastly.Client) {
			_ = DeleteTag(ctx, c, &DeleteTagInput{
				ServiceID: fastly.ToPointer(serviceID),
				TagID:     fastly.ToPointer(tag.ID),
			})
		})
	}()

	// Create operation.
	var op *Operation
	fastly.Record(t, "create_operation", func(c *fastly.Client) {
		op, err = Create(ctx, c, &CreateInput{
			ServiceID: fastly.ToPointer(serviceID),
			Method:    fastly.ToPointer("GET"),
			Domain:    fastly.ToPointer("example.com"),
			Path:      fastly.ToPointer("/test"),
			TagIDs:    []string{tag.ID},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, op)
	require.NotEmpty(t, op.ID)
	require.Equal(t, "GET", op.Method)
	require.Equal(t, "example.com", op.Domain)
	require.Equal(t, "/test", op.Path)
	require.Contains(t, op.TagIDs, tag.ID)

	defer func() {
		fastly.Record(t, "delete_operation", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ServiceID:   fastly.ToPointer(serviceID),
				OperationID: fastly.ToPointer(op.ID),
			})
		})
	}()

	// Describe operation.
	var described *Operation
	fastly.Record(t, "describe_operation", func(c *fastly.Client) {
		described, err = Describe(ctx, c, &DescribeInput{
			ServiceID:   fastly.ToPointer(serviceID),
			OperationID: fastly.ToPointer(op.ID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, described)
	require.Equal(t, op.ID, described.ID)
	require.Equal(t, op.Method, described.Method)
	require.Equal(t, op.Domain, described.Domain)
	require.Equal(t, op.Path, described.Path)

	// Update operation.
	var updated *Operation
	fastly.Record(t, "update_operation", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			ServiceID:   fastly.ToPointer(serviceID),
			OperationID: fastly.ToPointer(op.ID),
			Description: fastly.ToPointer("updated"),
			TagIDs:      []string{tag.ID},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, op.ID, updated.ID)
	require.Equal(t, "updated", updated.Description)
	require.Contains(t, updated.TagIDs, tag.ID)

	// List operations (with filters).
	var ops *Operations
	fastly.Record(t, "list_operations", func(c *fastly.Client) {
		ops, err = ListOperations(ctx, c, &ListOperationsInput{
			ServiceID: fastly.ToPointer(serviceID),
			Method:    fastly.ToPointer("GET"),
			Domain:    fastly.ToPointer("example.com"),
			Path:      fastly.ToPointer("/test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, ops)

	// Avoid brittle assertions: just ensure our created operation exists in the
	// returned data when filtering exactly.
	found := false
	for _, item := range ops.Data {
		if item.ID == op.ID {
			found = true
			break
		}
	}
	require.True(t, found, "expected created operation to appear in filtered list")

	// List discovered operations (status is required for now; includes filters).
	var discovered *DiscoveredOperations
	fastly.Record(t, "list_discovered_operations", func(c *fastly.Client) {
		discovered, err = ListDiscovered(ctx, c, &ListDiscoveredInput{
			ServiceID: fastly.ToPointer(serviceID),
			Status:    fastly.ToPointer("SAVED"),
			Method:    fastly.ToPointer("GET"),
			Domain:    fastly.ToPointer("example.com"),
			Path:      fastly.ToPointer("/test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, discovered)

	// Exercise discovered status update endpoints only if we have an ID.
	if len(discovered.Data) > 0 && discovered.Data[0].ID != "" {
		discoveredID := discovered.Data[0].ID

		var singleUpdated *DiscoveredOperation
		fastly.Record(t, "update_discovered_operation_status", func(c *fastly.Client) {
			singleUpdated, err = UpdateDiscoveredStatus(ctx, c, &UpdateDiscoveredStatusInput{
				ServiceID:   fastly.ToPointer(serviceID),
				OperationID: fastly.ToPointer(discoveredID),
				Status:      fastly.ToPointer("IGNORED"),
			})
		})
		require.NoError(t, err)
		require.NotNil(t, singleUpdated)
		require.Equal(t, discoveredID, singleUpdated.ID)

		var bulkUpdated *BulkOperationResultsResponse
		fastly.Record(t, "bulk_update_discovered_operation_status", func(c *fastly.Client) {
			bulkUpdated, err = BulkUpdateDiscoveredStatus(ctx, c, &BulkUpdateDiscoveredStatusInput{
				ServiceID:    fastly.ToPointer(serviceID),
				OperationIDs: []string{discoveredID},
				Status:       fastly.ToPointer("DISCOVERED"),
			})
		})
		require.NoError(t, err)
		require.NotNil(t, bulkUpdated)
		require.GreaterOrEqual(t, len(bulkUpdated.Data), 1)
	}

	// NOTE: Bulk operation endpoints are not live yet (confirmed internally).
	// Commented out until /operations-bulk and /operations-bulk-tags are deployed.
	/*
		// Exercise new bulk operation endpoints.
		var bulkCreated *BulkCreateOperationsResponse
		fastly.Record(t, "bulk_create_operations", func(c *fastly.Client) {
			bulkCreated, err = BulkCreateOperations(ctx, c, &BulkCreateOperationsInput{
				ServiceID: fastly.ToPointer(serviceID),
				Operations: []OperationBulkCreateItem{
					{
						Method:      fastly.ToPointer("GET"),
						Domain:      fastly.ToPointer("example.com"),
						Path:        fastly.ToPointer("/bulk-test-1"),
						Description: fastly.ToPointer("bulk test 1"),
						TagIDs:      []string{tag.ID},
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

		if len(createdIDs) > 0 {
			var bulkTagged *BulkOperationResultsResponse
			fastly.Record(t, "bulk_add_tags_to_operations", func(c *fastly.Client) {
				bulkTagged, err = BulkAddTags(ctx, c, &BulkAddTagsInput{
					ServiceID:    fastly.ToPointer(serviceID),
					OperationIDs: createdIDs,
					TagIDs:       []string{tag.ID},
				})
			})
			require.NoError(t, err)
			require.NotNil(t, bulkTagged)

			// Cleanup created operations.
			for _, id := range createdIDs {
				id := id
				fastly.Record(t, "delete_bulk_operation_"+id, func(c *fastly.Client) {
					_ = Delete(ctx, c, &DeleteInput{
						ServiceID:   fastly.ToPointer(serviceID),
						OperationID: fastly.ToPointer(id),
					})
				})
			}
		}
	*/
}
