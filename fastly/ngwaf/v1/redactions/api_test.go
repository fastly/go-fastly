package redactions

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

func TestClient_Redactions(t *testing.T) {
	t.Parallel()
	var err error
	testField := "somefield"
	testType := "request_parameter"
	var redactionID string

	// Create a test redaction.
	var redaction *Redaction

	fastly.Record(t, "create_redaction", func(c *fastly.Client) {
		redaction, err = Create(c, &CreateInput{
			Field:       fastly.ToPointer(testField),
			Type:        fastly.ToPointer(testType),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if redaction.Type != testType {
		t.Errorf("unexpected redaction type: got %q, expected %q", redaction.Type, testType)
	}
	if redaction.Field != testField {
		t.Errorf("unexpected redaction field: got %q, expected %q", redaction.Type, testField)
	}
	redactionID = redaction.RedactionID

	// Ensure we delete the test redaction at the end.
	defer func() {
		fastly.Record(t, "delete_redaction", func(c *fastly.Client) {
			err = Delete(c, &DeleteInput{
				RedactionID: fastly.ToPointer(redactionID),
				WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			})
		})
		if err != nil {
			t.Errorf("error during redaction cleanup: %v", err)
		}
	}()

	// Get the test redaction.
	var getTestRedaction *Redaction
	fastly.Record(t, "get_redaction", func(c *fastly.Client) {
		getTestRedaction, err = Get(c, &GetInput{
			RedactionID: fastly.ToPointer(redactionID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if getTestRedaction.Type != testType {
		t.Errorf("unexpected redaction type: got %q, expected %q", getTestRedaction.Type, testType)
	}
	if getTestRedaction.Field != testField {
		t.Errorf("unexpected redaction field: got %q, expected %q", getTestRedaction.Type, testField)
	}

	// Update the test redaction.
	const updatedRedactionType = "response_header"
	const updatedRedactionField string = "some-other-field"

	var updatedRedaction *Redaction
	fastly.Record(t, "update_redaction", func(c *fastly.Client) {
		updatedRedaction, err = Update(c, &UpdateInput{
			Field:       fastly.ToPointer(updatedRedactionField),
			RedactionID: fastly.ToPointer(redactionID),
			Type:        fastly.ToPointer(updatedRedactionType),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedRedaction.Type != updatedRedactionType {
		t.Errorf("unexpected redaction type: got %q, expected %q", updatedRedaction.Type, updatedRedactionType)
	}
	if updatedRedaction.Field != updatedRedactionField {
		t.Errorf("unexpected redaction field: got %q, expected %q", updatedRedaction.Type, updatedRedactionField)
	}

	// List the redactions for the test workspace and check the updated one is the only entry
	var listedRedactions *Redactions
	fastly.Record(t, "list_redaction", func(c *fastly.Client) {
		listedRedactions, err = List(c, &ListInput{
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(listedRedactions.Data) != 1 {
		t.Errorf("unexpected redaction list length: got %q, expected %q", len(listedRedactions.Data), 1)
	}
	if listedRedactions.Data[0].Type != updatedRedactionType {
		t.Errorf("unexpected redaction type: got %q, expected %q", listedRedactions.Data[0].Type, updatedRedactionType)
	}
	if listedRedactions.Data[0].Field != updatedRedactionField {
		t.Errorf("unexpected redaction field: got %q, expected %q", listedRedactions.Data[0].Type, updatedRedactionField)
	}
}

func TestClient_GetRedaction_validation(t *testing.T) {
	var err error
	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(fastly.TestClient, &GetInput{
		RedactionID: nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingRedactionID) {
		t.Errorf("expected ErrMissingRedactionID: got %s", err)
	}
}

func TestClient_CreateRedaction_validation(t *testing.T) {
	var err error
	_, err = Create(fastly.TestClient, &CreateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Field:       nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingField) {
		t.Errorf("expected ErrMissingField: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Field:       fastly.ToPointer("somefield"),
		Type:        nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingType) {
		t.Errorf("expected ErrMissingType: got %s", err)
	}
}

func TestClient_UpdateRedaction_validation(t *testing.T) {
	var err error
	_, err = Update(fastly.TestClient, &UpdateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Update(fastly.TestClient, &UpdateInput{
		RedactionID: nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingRedactionID) {
		t.Errorf("expected ErrMissingRedactionID: got %s", err)
	}
	_, err = Update(fastly.TestClient, &UpdateInput{
		RedactionID: fastly.ToPointer("someID"),
		Field:       nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingField) {
		t.Errorf("expected ErrMissingField: got %s", err)
	}
	_, err = Update(fastly.TestClient, &UpdateInput{
		RedactionID: fastly.ToPointer("someID"),
		Field:       fastly.ToPointer("somefield"),
		Type:        nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingType) {
		t.Errorf("expected ErrMissingType: got %s", err)
	}
}

func TestClient_DeleteRedaction_validation(t *testing.T) {
	var err error
	err = Delete(fastly.TestClient, &DeleteInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	err = Delete(fastly.TestClient, &DeleteInput{
		RedactionID: nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingRedactionID) {
		t.Errorf("expected ErrMissingRedactionID: got %s", err)
	}
}
