package operations

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v14/fastly"
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

	// Create operation #1.
	var op1 *Operation
	fastly.Record(t, "create_operation", func(c *fastly.Client) {
		op1, err = Create(ctx, c, &CreateInput{
			ServiceID: fastly.ToPointer(serviceID),
			Method:    fastly.ToPointer("GET"),
			Domain:    fastly.ToPointer("example.com"),
			Path:      fastly.ToPointer("/test"),
			TagIDs:    []string{tag.ID},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, op1)
	require.NotEmpty(t, op1.ID)

	defer func() {
		fastly.Record(t, "delete_operation", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ServiceID:   fastly.ToPointer(serviceID),
				OperationID: fastly.ToPointer(op1.ID),
			})
		})
	}()

	// Create operation #2 (for pagination test).
	var op2 *Operation
	fastly.Record(t, "create_operation_2", func(c *fastly.Client) {
		op2, err = Create(ctx, c, &CreateInput{
			ServiceID: fastly.ToPointer(serviceID),
			Method:    fastly.ToPointer("GET"),
			Domain:    fastly.ToPointer("example.com"),
			Path:      fastly.ToPointer("/test-pagination"),
			TagIDs:    []string{tag.ID},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, op2)
	require.NotEmpty(t, op2.ID)

	defer func() {
		fastly.Record(t, "delete_operation_2", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ServiceID:   fastly.ToPointer(serviceID),
				OperationID: fastly.ToPointer(op2.ID),
			})
		})
	}()

	// Describe operation.
	var described *Operation
	fastly.Record(t, "describe_operation", func(c *fastly.Client) {
		described, err = Describe(ctx, c, &DescribeInput{
			ServiceID:   fastly.ToPointer(serviceID),
			OperationID: fastly.ToPointer(op1.ID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, described)
	require.Equal(t, op1.ID, described.ID)

	// Update operation.
	var updated *Operation
	fastly.Record(t, "update_operation", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			ServiceID:   fastly.ToPointer(serviceID),
			OperationID: fastly.ToPointer(op1.ID),
			Description: fastly.ToPointer("updated"),
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
		if item.ID == op1.ID {
			found = true
			break
		}
	}
	require.True(t, found, "expected created operation to appear in filtered list")

	// ---- Pagination test for operations ----
	limit := 1
	p := NewOperationPaginator(ctx, fastly.TestClient, &ListOperationsInput{
		ServiceID: fastly.ToPointer(serviceID),
		Method:    []string{"GET"},
		Domain:    []string{"example.com"},
		Limit:     &limit,
		Page:      fastly.ToPointer(0),
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
			ServiceID: fastly.ToPointer(serviceID),
			Method:    []string{"GET"},
			Domain:    []string{"example.com"},
			Path:      fastly.ToPointer("/test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, discoveredFiltered)

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

	if len(discoveredAny.Data) > 0 && discoveredAny.Data[0].ID != "" {
		discoveredID := discoveredAny.Data[0].ID

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
	}

	// Bulk create ops + bulk tags.
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
					ServiceID:   fastly.ToPointer(serviceID),
					OperationID: fastly.ToPointer(opID),
				})
			}
		})
	}()

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
}
