package operations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v13/fastly"
)

func TestClient_Create_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: nil,
		Method:    fastly.ToPointer("GET"),
		Domain:    fastly.ToPointer("example.com"),
		Path:      fastly.ToPointer("/test"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: fastly.ToPointer("svc"),
		Method:    nil,
		Domain:    fastly.ToPointer("example.com"),
		Path:      fastly.ToPointer("/test"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingMethod)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: fastly.ToPointer("svc"),
		Method:    fastly.ToPointer("GET"),
		Domain:    nil,
		Path:      fastly.ToPointer("/test"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingDomain)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: fastly.ToPointer("svc"),
		Method:    fastly.ToPointer("GET"),
		Domain:    fastly.ToPointer("example.com"),
		Path:      nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingPath)
}

func TestClient_Describe_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Describe(ctx, fastly.TestClient, &DescribeInput{
		ServiceID:   nil,
		OperationID: fastly.ToPointer("op"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = Describe(ctx, fastly.TestClient, &DescribeInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_Delete_validation(t *testing.T) {
	ctx := context.TODO()

	err := Delete(ctx, fastly.TestClient, &DeleteInput{
		ServiceID:   nil,
		OperationID: fastly.ToPointer("op"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	err = Delete(ctx, fastly.TestClient, &DeleteInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_Update_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Update(ctx, fastly.TestClient, &UpdateInput{
		ServiceID:   nil,
		OperationID: fastly.ToPointer("op"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_ListOperations_validation(t *testing.T) {
	_, err := ListOperations(context.TODO(), fastly.TestClient, &ListOperationsInput{
		ServiceID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestClient_ListDiscovered_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := ListDiscovered(ctx, fastly.TestClient, &ListDiscoveredInput{
		ServiceID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestClient_Tags_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := CreateTag(ctx, fastly.TestClient, &CreateTagInput{
		ServiceID: nil,
		Name:      fastly.ToPointer("name"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = CreateTag(ctx, fastly.TestClient, &CreateTagInput{
		ServiceID: fastly.ToPointer("svc"),
		Name:      nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = DescribeTag(ctx, fastly.TestClient, &DescribeTagInput{
		ServiceID: nil,
		TagID:     fastly.ToPointer("tag"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = DescribeTag(ctx, fastly.TestClient, &DescribeTagInput{
		ServiceID: fastly.ToPointer("svc"),
		TagID:     nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	_, err = UpdateTag(ctx, fastly.TestClient, &UpdateTagInput{
		ServiceID: nil,
		TagID:     fastly.ToPointer("tag"),
		Name:      fastly.ToPointer("tag-name"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = UpdateTag(ctx, fastly.TestClient, &UpdateTagInput{
		ServiceID: fastly.ToPointer("svc"),
		TagID:     nil,
		Name:      fastly.ToPointer("tag-name"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	err = DeleteTag(ctx, fastly.TestClient, &DeleteTagInput{
		ServiceID: nil,
		TagID:     fastly.ToPointer("tag"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	err = DeleteTag(ctx, fastly.TestClient, &DeleteTagInput{
		ServiceID: fastly.ToPointer("svc"),
		TagID:     nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	_, err = ListTags(ctx, fastly.TestClient, &ListTagsInput{ServiceID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestClient_UpdateDiscoveredStatus_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := UpdateDiscoveredStatus(ctx, fastly.TestClient, &UpdateDiscoveredStatusInput{
		ServiceID:   nil,
		OperationID: fastly.ToPointer("op"),
		Status:      fastly.ToPointer("IGNORED"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = UpdateDiscoveredStatus(ctx, fastly.TestClient, &UpdateDiscoveredStatusInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: nil,
		Status:      fastly.ToPointer("IGNORED"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	_, err = UpdateDiscoveredStatus(ctx, fastly.TestClient, &UpdateDiscoveredStatusInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: fastly.ToPointer("op"),
		Status:      nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingStatus)
}

func TestClient_BulkUpdateDiscoveredStatus_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := BulkUpdateDiscoveredStatus(ctx, fastly.TestClient, &BulkUpdateDiscoveredStatusInput{
		ServiceID:    nil,
		OperationIDs: []string{"op"},
		Status:       fastly.ToPointer("IGNORED"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = BulkUpdateDiscoveredStatus(ctx, fastly.TestClient, &BulkUpdateDiscoveredStatusInput{
		ServiceID:    fastly.ToPointer("svc"),
		OperationIDs: []string{"op"},
		Status:       nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingStatus)

	_, err = BulkUpdateDiscoveredStatus(ctx, fastly.TestClient, &BulkUpdateDiscoveredStatusInput{
		ServiceID:    fastly.ToPointer("svc"),
		OperationIDs: nil,
		Status:       fastly.ToPointer("IGNORED"),
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "OperationIDs")
}

func TestClient_BulkCreateOperations_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := BulkCreateOperations(ctx, fastly.TestClient, &BulkCreateOperationsInput{
		ServiceID: nil,
		Operations: []OperationBulkCreateItem{
			{
				Method: fastly.ToPointer("GET"),
				Domain: fastly.ToPointer("example.com"),
				Path:   fastly.ToPointer("/x"),
			},
		},
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = BulkCreateOperations(ctx, fastly.TestClient, &BulkCreateOperationsInput{
		ServiceID:  fastly.ToPointer("svc"),
		Operations: nil,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "Operations")
}

func TestClient_BulkAddTags_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := BulkAddTags(ctx, fastly.TestClient, &BulkAddTagsInput{
		ServiceID:    nil,
		OperationIDs: []string{"op"},
		TagIDs:       []string{"tag"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = BulkAddTags(ctx, fastly.TestClient, &BulkAddTagsInput{
		ServiceID:    fastly.ToPointer("svc"),
		OperationIDs: nil,
		TagIDs:       []string{"tag"},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "OperationIDs")

	_, err = BulkAddTags(ctx, fastly.TestClient, &BulkAddTagsInput{
		ServiceID:    fastly.ToPointer("svc"),
		OperationIDs: []string{"op"},
		TagIDs:       nil,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "TagIDs")
}
