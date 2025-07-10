package webhook

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Global workspace value for the tests.
var testWorkspaceID = fastly.TestNGWAFWorkspaceID

func Test_WorkspaceAlerts(t *testing.T) {
	t.Parallel()

	var AlertID string
	var err error
	var WorkSpaceAlert *WorkspaceAlert
	testConfig := CreateConfig{
		Webhook: fastly.ToPointer("https://example.com/webhook"),
	}
	testDescription := "This is a test alert."
	testEvent := "flag"
	testType := IntegrationType

	// Create a workspace alert.
	fastly.Record(t, "create_workspacealerts", func(c *fastly.Client) {
		WorkSpaceAlert, err = Create(c, &CreateInput{
			Type:        fastly.ToPointer(testType),
			Config:      testConfig,
			Events:      []string{testEvent},
			Description: fastly.ToPointer(testDescription),
			WorkspaceID: &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if WorkSpaceAlert == nil {
		t.Fatal("expected workspace alert response, got nil")
	}
	AlertID = WorkSpaceAlert.ID

	// Ensure that we delete the test workspace alert after use.
	defer func() {
		fastly.Record(t, "delete_workspaceAlert", func(c *fastly.Client) {
			err = Delete(c, &DeleteInput{
				AlertID:     fastly.ToPointer(AlertID),
				WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			})
		})
		if err != nil {
			t.Errorf("error during workspace alert cleanup: %v", err)
		}
	}()

	// Get the test workspace alert.
	var getTestWorkspaceAlert *WorkspaceAlert
	fastly.Record(t, "get_workspacealert", func(c *fastly.Client) {
		getTestWorkspaceAlert, err = Get(c, &GetInput{
			AlertID:     fastly.ToPointer(AlertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *getTestWorkspaceAlert.Config.Webhook != *testConfig.Webhook {
		t.Errorf("unexpected workspace alert config webhook: got %v, expected %v", *getTestWorkspaceAlert.Config.Webhook, *testConfig.Webhook)
	}
	if getTestWorkspaceAlert.Type != testType {
		t.Errorf("unexpected workspace alert type: got %+v, expected %+v", getTestWorkspaceAlert.Type, testType)
	}
	if len(getTestWorkspaceAlert.Events) != 1 {
		t.Errorf("unexpected workspace alerts event length: got %d, expected %d", len(getTestWorkspaceAlert.Events), 1)
	}
	if getTestWorkspaceAlert.Events[0] != testEvent {
		t.Errorf("unexpected workspace alert events: got %+v, expected %+v", getTestWorkspaceAlert.Events[0], testEvent)
	}
	if getTestWorkspaceAlert.Description != testDescription {
		t.Errorf("unexpected workspace alert description: got %+v, expected %+v", getTestWorkspaceAlert.Description, testDescription)
	}

	// Get the signing key for the webhook alert.
	var signingKey *WorkspaceAlertsKey
	fastly.Record(t, "get_workspacealert_signing_key", func(c *fastly.Client) {
		signingKey, err = GetKey(c, &GetKeyInput{
			AlertID:     fastly.ToPointer(AlertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if signingKey == nil {
		t.Fatal("expected signing key response, got nil")
	}
	originalKey := signingKey.SigningKey

	// Rotate the signing key.
	var rotatedKey *WorkspaceAlertsKey
	fastly.Record(t, "rotate_workspacealert_signing_key", func(c *fastly.Client) {
		rotatedKey, err = RotateKey(c, &RotateKeyInput{
			AlertID:     fastly.ToPointer(AlertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if rotatedKey == nil {
		t.Fatal("expected rotated key response, got nil")
	}
	if rotatedKey.SigningKey == originalKey {
		t.Errorf("expected rotated key to be different from original key")
	}

	// Update the test workspace alert.
	updatedConfig := UpdateConfig{
		Webhook: fastly.ToPointer("https://updated.example.com/webhook"),
	}
	updatedEvent := "flag"
	var updateWorkspaceAlert *WorkspaceAlert
	fastly.Record(t, "update_workspacealert", func(c *fastly.Client) {
		updateWorkspaceAlert, err = Update(c, &UpdateInput{
			AlertID:     fastly.ToPointer(AlertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			Config:      updatedConfig,
			Events:      []string{updatedEvent},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *updateWorkspaceAlert.Config.Webhook != *updatedConfig.Webhook {
		t.Errorf("unexpected updated workspace alert config webhook: got %v, expected %v", *updateWorkspaceAlert.Config.Webhook, *updatedConfig.Webhook)
	}
	if len(updateWorkspaceAlert.Events) != 1 {
		t.Errorf("unexpected updated workspace alerts event length: got %d, expected %d", len(updateWorkspaceAlert.Events), 1)
	}
	if updateWorkspaceAlert.Events[0] != updatedEvent {
		t.Errorf("unexpected updated workspace alert events: got %+v, expected %+v", updateWorkspaceAlert.Events[0], updatedEvent)
	}

	// List the workspace alerts for the test workspace and check the updated one is the only entry.
	var WorkspaceAlerts *WorkspaceAlerts
	fastly.Record(t, "list_workspacealerts", func(c *fastly.Client) {
		WorkspaceAlerts, err = List(c, &ListInput{
			WorkspaceID: &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if WorkspaceAlerts == nil {
		t.Fatal("expected workspace alert response, got nil")
	}
	if len(WorkspaceAlerts.Data) != 1 {
		t.Errorf("unexpected workspace alerts list length: got %d, expected %d", len(WorkspaceAlerts.Data), 1)
	}
	// Validate the listed alert matches the updated values
	listedAlert := WorkspaceAlerts.Data[0]
	if *listedAlert.Config.Webhook != *updatedConfig.Webhook {
		t.Errorf("unexpected listed workspace alert config webhook: got %v, expected %v", *listedAlert.Config.Webhook, *updatedConfig.Webhook)
	}
	if len(listedAlert.Events) != 1 {
		t.Errorf("unexpected listed workspace alerts event length: got %d, expected %d", len(listedAlert.Events), 1)
	}
	if listedAlert.Events[0] != updatedEvent {
		t.Errorf("unexpected listed workspace alert events: got %+v, expected %+v", listedAlert.Events[0], updatedEvent)
	}
}

func TestClient_CreateWorkspaceAlert_validation(t *testing.T) {
	var err error
	_, err = Create(fastly.TestClient, &CreateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Type:        nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrInvalidConfigType) {
		t.Errorf("expected ErrInvalidConfigType: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer(IntegrationType),
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingConfig) {
		t.Errorf("expected ErrMissingConfig: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer(IntegrationType),
		Config:      CreateConfig{Webhook: fastly.ToPointer("https://example.com/webhook")},
		Events:      nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingEvents) {
		t.Errorf("expected ErrMissingEvents: got %s", err)
	}
}

func TestClient_GetWorkspaceAlert_validation(t *testing.T) {
	var err error
	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(fastly.TestClient, &GetInput{
		AlertID:     nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingAlertID) {
		t.Errorf("expected ErrMissingAlertID: got %s", err)
	}
}

func TestClient_UpdateWorkspaceAlert_validation(t *testing.T) {
	var err error
	_, err = Update(fastly.TestClient, &UpdateInput{
		AlertID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingAlertID) {
		t.Errorf("expected ErrMissingAlertID: got %s", err)
	}
	_, err = Update(fastly.TestClient, &UpdateInput{
		AlertID:     fastly.ToPointer("test-id"),
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}

func TestClient_DeleteWorkspaceAlert_validation(t *testing.T) {
	var err error
	err = Delete(fastly.TestClient, &DeleteInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	err = Delete(fastly.TestClient, &DeleteInput{
		AlertID:     nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingAlertID) {
		t.Errorf("expected ErrMissingAlertID: got %s", err)
	}
}

func TestClient_ListWorkspaceAlerts_validation(t *testing.T) {
	var err error
	_, err = List(fastly.TestClient, &ListInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}

func TestClient_GetKey_validation(t *testing.T) {
	var err error
	_, err = GetKey(fastly.TestClient, &GetKeyInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = GetKey(fastly.TestClient, &GetKeyInput{
		AlertID:     nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingAlertID) {
		t.Errorf("expected ErrMissingAlertID: got %s", err)
	}
}

func TestClient_RotateKey_validation(t *testing.T) {
	var err error
	_, err = RotateKey(fastly.TestClient, &RotateKeyInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = RotateKey(fastly.TestClient, &RotateKeyInput{
		AlertID:     nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingAlertID) {
		t.Errorf("expected ErrMissingAlertID: got %s", err)
	}
}