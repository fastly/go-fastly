package operations

import (
	"context"
	"testing"

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
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: fastly.ToPointer("svc"),
		Method:    nil,
		Domain:    fastly.ToPointer("example.com"),
		Path:      fastly.ToPointer("/test"),
	})
	if err != fastly.ErrMissingMethod {
		t.Errorf("expected ErrMissingMethod: got %v", err)
	}

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: fastly.ToPointer("svc"),
		Method:    fastly.ToPointer("GET"),
		Domain:    nil,
		Path:      fastly.ToPointer("/test"),
	})
	if err != fastly.ErrMissingDomain {
		t.Errorf("expected ErrMissingDomain: got %v", err)
	}

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		ServiceID: fastly.ToPointer("svc"),
		Method:    fastly.ToPointer("GET"),
		Domain:    fastly.ToPointer("example.com"),
		Path:      nil,
	})
	if err != fastly.ErrMissingPath {
		t.Errorf("expected ErrMissingPath: got %v", err)
	}
}

func TestClient_Describe_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Describe(ctx, fastly.TestClient, &DescribeInput{
		ServiceID:   nil,
		OperationID: fastly.ToPointer("op"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	_, err = Describe(ctx, fastly.TestClient, &DescribeInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: nil,
	})
	if err != fastly.ErrMissingID {
		t.Errorf("expected ErrMissingID: got %v", err)
	}
}

func TestClient_Delete_validation(t *testing.T) {
	ctx := context.TODO()

	err := Delete(ctx, fastly.TestClient, &DeleteInput{
		ServiceID:   nil,
		OperationID: fastly.ToPointer("op"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	err = Delete(ctx, fastly.TestClient, &DeleteInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: nil,
	})
	if err != fastly.ErrMissingID {
		t.Errorf("expected ErrMissingID: got %v", err)
	}
}

func TestClient_Update_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Update(ctx, fastly.TestClient, &UpdateInput{
		ServiceID:   nil,
		OperationID: fastly.ToPointer("op"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: nil,
	})
	if err != fastly.ErrMissingID {
		t.Errorf("expected ErrMissingID: got %v", err)
	}
}

func TestClient_ListOperations_validation(t *testing.T) {
	_, err := ListOperations(context.TODO(), fastly.TestClient, &ListOperationsInput{
		ServiceID: nil,
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}
}

func TestClient_ListDiscovered_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := ListDiscovered(ctx, fastly.TestClient, &ListDiscoveredInput{
		ServiceID: nil,
		Status:    fastly.ToPointer("SAVED"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	_, err = ListDiscovered(ctx, fastly.TestClient, &ListDiscoveredInput{
		ServiceID: fastly.ToPointer("svc"),
		Status:    nil,
	})
	if err != fastly.ErrMissingStatus {
		t.Errorf("expected ErrMissingStatus: got %v", err)
	}
}

func TestClient_Tags_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := CreateTag(ctx, fastly.TestClient, &CreateTagInput{
		ServiceID: nil,
		Name:      fastly.ToPointer("name"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	_, err = CreateTag(ctx, fastly.TestClient, &CreateTagInput{
		ServiceID: fastly.ToPointer("svc"),
		Name:      nil,
	})
	if err != fastly.ErrMissingName {
		t.Errorf("expected ErrMissingName: got %v", err)
	}

	_, err = DescribeTag(ctx, fastly.TestClient, &DescribeTagInput{
		ServiceID: nil,
		TagID:     fastly.ToPointer("tag"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	_, err = DescribeTag(ctx, fastly.TestClient, &DescribeTagInput{
		ServiceID: fastly.ToPointer("svc"),
		TagID:     nil,
	})
	if err != fastly.ErrMissingID {
		t.Errorf("expected ErrMissingID: got %v", err)
	}

	_, err = UpdateTag(ctx, fastly.TestClient, &UpdateTagInput{
		ServiceID: nil,
		TagID:     fastly.ToPointer("tag"),
		Name:      fastly.ToPointer("tag-name"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	_, err = UpdateTag(ctx, fastly.TestClient, &UpdateTagInput{
		ServiceID: fastly.ToPointer("svc"),
		TagID:     nil,
		Name:      fastly.ToPointer("tag-name"),
	})
	if err != fastly.ErrMissingID {
		t.Errorf("expected ErrMissingID: got %v", err)
	}

	err = DeleteTag(ctx, fastly.TestClient, &DeleteTagInput{
		ServiceID: nil,
		TagID:     fastly.ToPointer("tag"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	err = DeleteTag(ctx, fastly.TestClient, &DeleteTagInput{
		ServiceID: fastly.ToPointer("svc"),
		TagID:     nil,
	})
	if err != fastly.ErrMissingID {
		t.Errorf("expected ErrMissingID: got %v", err)
	}

	_, err = ListTags(ctx, fastly.TestClient, &ListTagsInput{ServiceID: nil})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}
}

func TestClient_UpdateDiscoveredStatus_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := UpdateDiscoveredStatus(ctx, fastly.TestClient, &UpdateDiscoveredStatusInput{
		ServiceID:   nil,
		OperationID: fastly.ToPointer("op"),
		Status:      fastly.ToPointer("IGNORED"),
	})
	if err != fastly.ErrMissingServiceID {
		t.Errorf("expected ErrMissingServiceID: got %v", err)
	}

	_, err = UpdateDiscoveredStatus(ctx, fastly.TestClient, &UpdateDiscoveredStatusInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: nil,
		Status:      fastly.ToPointer("IGNORED"),
	})
	if err != fastly.ErrMissingID {
		t.Errorf("expected ErrMissingID: got %v", err)
	}

	_, err = UpdateDiscoveredStatus(ctx, fastly.TestClient, &UpdateDiscoveredStatusInput{
		ServiceID:   fastly.ToPointer("svc"),
		OperationID: fastly.ToPointer("op"),
		Status:      nil,
	})
	if err != fastly.ErrMissingStatus {
		t.Errorf("expected ErrMissingStatus: got %v", err)
	}
}
