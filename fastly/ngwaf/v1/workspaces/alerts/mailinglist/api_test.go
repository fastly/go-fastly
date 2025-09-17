package mailinglist

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v12/fastly"
)

// Global workspace value for the tests.
var testWorkspaceID = fastly.TestNGWAFWorkspaceID

func Test_Alerts(t *testing.T) {
	var alertID string
	var err error
	var alert *Alert
	testConfig := &CreateConfig{
		Address: fastly.ToPointer("test@example.com"),
	}
	testDescription := "This is a test alert."
	testEvent := "flag"
	testType := IntegrationType

	// Create a workspace alert.
	fastly.Record(t, "create_alert", func(c *fastly.Client) {
		alert, err = Create(context.TODO(), c, &CreateInput{
			Config:      testConfig,
			Events:      &[]string{testEvent},
			Description: fastly.ToPointer(testDescription),
			WorkspaceID: &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if alert == nil {
		t.Fatal("expected alert response, got nil")
	}
	if *alert.Config.Address != *testConfig.Address {
		t.Errorf("unexpected alert config address: got %v, expected %v", *alert.Config.Address, *testConfig.Address)
	}
	if alert.Type != testType {
		t.Errorf("unexpected alert type: got %+v, expected %+v", alert.Type, testType)
	}
	if len(alert.Events) != 1 {
		t.Errorf("unexpected alerts event length: got %d, expected %d", len(alert.Events), 1)
	}
	if alert.Events[0] != testEvent {
		t.Errorf("unexpected alert events: got %+v, expected %+v", alert.Events[0], testEvent)
	}
	if alert.Description != testDescription {
		t.Errorf("unexpected alert description: got %+v, expected %+v", alert.Description, testDescription)
	}
	alertID = alert.ID

	// Ensure that we delete the test workspace alert after use.
	defer func() {
		fastly.Record(t, "delete_alert", func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				AlertID:     fastly.ToPointer(alertID),
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
		getTestAlert, err = Get(context.TODO(), c, &GetInput{
			AlertID:     fastly.ToPointer(alertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *getTestAlert.Config.Address != *testConfig.Address {
		t.Errorf("unexpected workspace alert config address: got %v, expected %v", *getTestAlert.Config.Address, *testConfig.Address)
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
	updatedConfig := &UpdateConfig{
		Address: fastly.ToPointer("updated@example.com"),
	}
	updatedEvent := "flag"
	var updateAlert *Alert
	fastly.Record(t, "update_alert", func(c *fastly.Client) {
		updateAlert, err = Update(context.TODO(), c, &UpdateInput{
			AlertID:     fastly.ToPointer(alertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			Config:      updatedConfig,
			Events:      &[]string{updatedEvent},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *updateAlert.Config.Address != *updatedConfig.Address {
		t.Errorf("unexpected updated workspace alert config address: got %v, expected %v", *updateAlert.Config.Address, *updatedConfig.Address)
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
		Alerts, err = List(context.TODO(), c, &ListInput{
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
	if *listedAlert.Config.Address != *updatedConfig.Address {
		t.Errorf("unexpected listed workspace alert config address: got %v, expected %v", *listedAlert.Config.Address, *updatedConfig.Address)
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
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingConfig) {
		t.Errorf("expected ErrMissingConfig: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Config:      &CreateConfig{Address: fastly.ToPointer("test@example.com")},
		Events:      nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingEvents) {
		t.Errorf("expected ErrMissingEvents: got %s", err)
	}
}

func TestClient_GetAlert_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		AlertID:     nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingAlertID) {
		t.Errorf("expected ErrMissingAlertID: got %s", err)
	}
}

func TestClient_UpdateAlert_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		AlertID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingAlertID) {
		t.Errorf("expected ErrMissingAlertID: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		AlertID:     fastly.ToPointer("test-id"),
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}

func TestClient_DeleteAlert_validation(t *testing.T) {
	var err error
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		AlertID:     nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingAlertID) {
		t.Errorf("expected ErrMissingAlertID: got %s", err)
	}
}

func TestClient_ListAlerts_validation(t *testing.T) {
	var err error
	_, err = List(context.TODO(), fastly.TestClient, &ListInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}
