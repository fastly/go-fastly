package operations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_Create_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: nil,
		Method:    new("GET"),
		Domain:    new("example.com"),
		Path:      new("/test"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: new("svc"),
		Method:    nil,
		Domain:    new("example.com"),
		Path:      new("/test"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingMethod)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: new("svc"),
		Method:    new("GET"),
		Domain:    nil,
		Path:      new("/test"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingDomain)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: new("svc"),
		Method:    new("GET"),
		Domain:    new("example.com"),
		Path:      nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingPath)
}

func TestClient_Describe_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Describe(ctx, fastly.TestClient, &DescribeInput{
		ServiceID:   nil,
		OperationID: new("op"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = Describe(ctx, fastly.TestClient, &DescribeInput{
		ServiceID:   new("svc"),
		OperationID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_Delete_validation(t *testing.T) {
	ctx := context.TODO()

	err := Delete(ctx, fastly.TestClient, &DeleteInput{
		ServiceID:   nil,
		OperationID: new("op"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	err = Delete(ctx, fastly.TestClient, &DeleteInput{
		ServiceID:   new("svc"),
		OperationID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_Update_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Update(ctx, fastly.TestClient, &UpdateInput{
		ServiceID:   nil,
		OperationID: new("op"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{
		ServiceID:   new("svc"),
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
		Name:      new("name"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = CreateTag(ctx, fastly.TestClient, &CreateTagInput{
		ServiceID: new("svc"),
		Name:      nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = DescribeTag(ctx, fastly.TestClient, &DescribeTagInput{
		ServiceID: nil,
		TagID:     new("tag"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = DescribeTag(ctx, fastly.TestClient, &DescribeTagInput{
		ServiceID: new("svc"),
		TagID:     nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	_, err = UpdateTag(ctx, fastly.TestClient, &UpdateTagInput{
		ServiceID: nil,
		TagID:     new("tag"),
		Name:      new("tag-name"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = UpdateTag(ctx, fastly.TestClient, &UpdateTagInput{
		ServiceID: new("svc"),
		TagID:     nil,
		Name:      new("tag-name"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	err = DeleteTag(ctx, fastly.TestClient, &DeleteTagInput{
		ServiceID: nil,
		TagID:     new("tag"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	err = DeleteTag(ctx, fastly.TestClient, &DeleteTagInput{
		ServiceID: new("svc"),
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
		OperationID: new("op"),
		Status:      new("IGNORED"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = UpdateDiscoveredStatus(ctx, fastly.TestClient, &UpdateDiscoveredStatusInput{
		ServiceID:   new("svc"),
		OperationID: nil,
		Status:      new("IGNORED"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	_, err = UpdateDiscoveredStatus(ctx, fastly.TestClient, &UpdateDiscoveredStatusInput{
		ServiceID:   new("svc"),
		OperationID: new("op"),
		Status:      nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingStatus)
}

func TestClient_BulkUpdateDiscoveredStatus_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := BulkUpdateDiscoveredStatus(ctx, fastly.TestClient, &BulkUpdateDiscoveredStatusInput{
		ServiceID:    nil,
		OperationIDs: []string{"op"},
		Status:       new("IGNORED"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = BulkUpdateDiscoveredStatus(ctx, fastly.TestClient, &BulkUpdateDiscoveredStatusInput{
		ServiceID:    new("svc"),
		OperationIDs: []string{"op"},
		Status:       nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingStatus)

	_, err = BulkUpdateDiscoveredStatus(ctx, fastly.TestClient, &BulkUpdateDiscoveredStatusInput{
		ServiceID:    new("svc"),
		OperationIDs: nil,
		Status:       new("IGNORED"),
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
				Method: new("GET"),
				Domain: new("example.com"),
				Path:   new("/x"),
			},
		},
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceID)

	_, err = BulkCreateOperations(ctx, fastly.TestClient, &BulkCreateOperationsInput{
		ServiceID:  new("svc"),
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
		ServiceID:    new("svc"),
		OperationIDs: nil,
		TagIDs:       []string{"tag"},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "OperationIDs")

	_, err = BulkAddTags(ctx, fastly.TestClient, &BulkAddTagsInput{
		ServiceID:    new("svc"),
		OperationIDs: []string{"op"},
		TagIDs:       nil,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "TagIDs")
}
