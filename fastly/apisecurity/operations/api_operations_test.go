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

	// Create a tag to associate with operations.
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
			Method:    []string{"GET"},
			Domain:    []string{"example.com"},
			Path:      fastly.ToPointer("/test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, ops)

	found := false
	for _, item := range ops.Data {
		if item.ID == op.ID {
			found = true
			break
		}
	}
	require.True(t, found, "expected created operation to appear in filtered list")

	// List discovered operations with filters.
	// This list can legitimately be empty (discovered operations depend on traffic).
	var discoveredFiltered *DiscoveredOperations
	fastly.Record(t, "list_discovered_operations", func(c *fastly.Client) {
		discoveredFiltered, err = ListDiscovered(ctx, c, &ListDiscoveredInput{
			ServiceID: fastly.ToPointer(serviceID),
			Method:    []string{"GET"},
			Domain:    []string{"example.com"},
			Path:      fastly.ToPointer("/test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, discoveredFiltered)

	// Fetch at least one discovered operation ID (unfiltered) so we can exercise
	// the discovered status update endpoints.
	var discoveredAny *DiscoveredOperations
	fastly.Record(t, "list_discovered_operations_any", func(c *fastly.Client) {
		discoveredAny, err = ListDiscovered(ctx, c, &ListDiscoveredInput{
			ServiceID: fastly.ToPointer(serviceID),
			Limit:     fastly.ToPointer(1),
			Page:      fastly.ToPointer(0),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, discoveredAny)

	if len(discoveredAny.Data) == 0 || discoveredAny.Data[0].ID == "" {
		// We still keep the rest of the test coverage (operations + bulk ops),
		// but we can't exercise the discovered status update endpoints without
		// a discovered operation ID.
		t.Skip("no discovered operations available to exercise discovered status update endpoints")
	}

	discoveredID := discoveredAny.Data[0].ID

	// Single update discovered operation status.
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

	// Bulk update discovered operation status.
	var bulkUpdatedDiscovered *BulkOperationResultsResponse
	fastly.Record(t, "bulk_update_discovered_operation_status", func(c *fastly.Client) {
		bulkUpdatedDiscovered, err = BulkUpdateDiscoveredStatus(ctx, c, &BulkUpdateDiscoveredStatusInput{
			ServiceID:    fastly.ToPointer(serviceID),
			OperationIDs: []string{discoveredID},
			Status:       fastly.ToPointer("IGNORED"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, bulkUpdatedDiscovered)
	require.GreaterOrEqual(t, len(bulkUpdatedDiscovered.Data), 1)

	// Exercise bulk operation endpoints.
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
				},
			},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, bulkCreated)
	require.GreaterOrEqual(t, len(bulkCreated.Data), 1)

	var createdIDs []string
	for _, r := range bulkCreated.Data {
		if r.Operation != nil && r.Operation.ID != "" {
			createdIDs = append(createdIDs, r.Operation.ID)
		}
	}
	require.NotEmpty(t, createdIDs, "expected at least one bulk-created operation id")

	// Cleanup bulk-created operations in a single recording to avoid fixture sprawl.
	defer func() {
		if len(createdIDs) == 0 {
			return
		}

		fastly.Record(t, "delete_bulk_operations", func(c *fastly.Client) {
			for _, opID := range createdIDs {
				_ = Delete(ctx, c, &DeleteInput{
					ServiceID:   fastly.ToPointer(serviceID),
					OperationID: fastly.ToPointer(opID),
				})
			}
		})
	}()

	// Bulk add tags to operations.
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
	require.GreaterOrEqual(t, len(bulkTagged.Data), 1)
}
