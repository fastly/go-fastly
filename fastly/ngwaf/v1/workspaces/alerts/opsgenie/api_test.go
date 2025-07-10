package opsgenie

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Global workspace value for the tests.
var testWorkspaceID = fastly.TestNGWAFWorkspaceID

func Test_Alerts(t *testing.T) {
	t.Parallel()

	var AlertID string
	var err error
	var WorkSpaceAlert *Alert
	testConfig := CreateConfig{
		Key: fastly.ToPointer("123456789"),
	}
	testDescription := "This is a test alert."
	testEvent := "flag"
	testType := IntegrationType

	// Create a workspace alert.
	fastly.Record(t, "create_alert", func(c *fastly.Client) {
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
		fastly.Record(t, "delete_alert", func(c *fastly.Client) {
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
	var getTestAlert *Alert
	fastly.Record(t, "get_alert", func(c *fastly.Client) {
		getTestAlert, err = Get(c, &GetInput{
			AlertID:     fastly.ToPointer(AlertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *getTestAlert.Config.Key != *testConfig.Key {
		t.Errorf("unexpected workspace alert config key: got %v, expected %v", *getTestAlert.Config.Key, *testConfig.Key)
	}
	if getTestAlert.Type != testType {
		t.Errorf("unexpected workspace alert type: got %+v, expected %+v", getTestAlert.Type, testType)
	}
	if len(getTestAlert.Events) != 1 {
		t.Errorf("unexpected workspace alerts event length: got %d, expected %d", len(getTestAlert.Events), 1)
	}
	if getTestAlert.Events[0] != testEvent {
		t.Errorf("unexpected workspace alert events: got %+v, expected %+v", getTestAlert.Events[0], testEvent)
	}
	if getTestAlert.Description != testDescription {
		t.Errorf("unexpected workspace alert description: got %+v, expected %+v", getTestAlert.Description, testDescription)
	}

	// Update the test workspace alert.
	updatedConfig := UpdateConfig{
		Key: fastly.ToPointer("987654321"),
	}
	updatedEvent := "flag"
	var updateAlert *Alert
	fastly.Record(t, "update_alert", func(c *fastly.Client) {
		updateAlert, err = Update(c, &UpdateInput{
			AlertID:     fastly.ToPointer(AlertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			Config:      updatedConfig,
			Events:      []string{updatedEvent},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *updateAlert.Config.Key != *updatedConfig.Key {
		t.Errorf("unexpected updated workspace alert config key: got %v, expected %v", *updateAlert.Config.Key, *updatedConfig.Key)
	}
	if len(updateAlert.Events) != 1 {
		t.Errorf("unexpected updated workspace alerts event length: got %d, expected %d", len(updateAlert.Events), 1)
	}
	if updateAlert.Events[0] != updatedEvent {
		t.Errorf("unexpected updated workspace alert events: got %+v, expected %+v", updateAlert.Events[0], updatedEvent)
	}

	// List the workspace alerts for the test workspace and check the updated one is the only entry.
	var Alerts *Alerts
	fastly.Record(t, "list_alerts", func(c *fastly.Client) {
		Alerts, err = List(c, &ListInput{
			WorkspaceID: &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if Alerts == nil {
		t.Fatal("expected workspace alert response, got nil")
	}
	if len(Alerts.Data) != 1 {
		t.Errorf("unexpected workspace alerts list length: got %d, expected %d", len(Alerts.Data), 1)
	}
	// Validate the listed alert matches the updated values
	listedAlert := Alerts.Data[0]
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

func TestClient_CreateAlert_validation(t *testing.T) {
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
		Config:      CreateConfig{Key: fastly.ToPointer("111222333")},
		Events:      nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingEvents) {
		t.Errorf("expected ErrMissingEvents: got %s", err)
	}
}

func TestClient_GetAlert_validation(t *testing.T) {
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

func TestClient_UpdateAlert_validation(t *testing.T) {
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

func TestClient_DeleteAlert_validation(t *testing.T) {
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

func TestClient_ListAlerts_validation(t *testing.T) {
	var err error
	_, err = List(fastly.TestClient, &ListInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}