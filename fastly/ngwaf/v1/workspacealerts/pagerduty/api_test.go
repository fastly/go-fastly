package pagerduty

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
		Key: fastly.ToPointer("a1b2c3d4e5f6789012345678901234567"),
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
	if *getTestWorkspaceAlert.Config.Key != *testConfig.Key {
		t.Errorf("unexpected workspace alert config key: got %v, expected %v", *getTestWorkspaceAlert.Config.Key, *testConfig.Key)
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

	// Update the test workspace alert.
	updatedConfig := UpdateConfig{
		Key: fastly.ToPointer("b2c3d4e5f6789012345678901234567a"),
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
	if *updateWorkspaceAlert.Config.Key != *updatedConfig.Key {
		t.Errorf("unexpected updated workspace alert config key: got %v, expected %v", *updateWorkspaceAlert.Config.Key, *updatedConfig.Key)
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
	if *listedAlert.Config.Key != *updatedConfig.Key {
		t.Errorf("unexpected listed workspace alert config key: got %v, expected %v", *listedAlert.Config.Key, *updatedConfig.Key)
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
		Config:      CreateConfig{Key: fastly.ToPointer("test-key")},
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