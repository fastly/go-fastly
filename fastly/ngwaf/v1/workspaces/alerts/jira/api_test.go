package jira

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
	var Alert *Alert
	testConfig := CreateConfig{
		Host:      fastly.ToPointer("https://example.exampleJira.net"),
		Key:       fastly.ToPointer("ATATT3xFfGF0a1b2c3d4e5f6789012345678901234567890"),
		Project:   fastly.ToPointer("TEST"),
		Username:  fastly.ToPointer("testuser"),
		IssueType: fastly.ToPointer("Bug"),
	}
	testDescription := "This is a test alert."
	testEvent := "flag"
	testType := IntegrationType

	// Create a workspace alert.
	fastly.Record(t, "create_workspacealerts", func(c *fastly.Client) {
		Alert, err = Create(c, &CreateInput{
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
	if Alert == nil {
		t.Fatal("expected workspace alert response, got nil")
	}
	AlertID = Alert.ID

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
	var getTestAlert *Alert
	fastly.Record(t, "get_workspacealert", func(c *fastly.Client) {
		getTestAlert, err = Get(c, &GetInput{
			AlertID:     fastly.ToPointer(AlertID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})

	})
	if err != nil {
		t.Fatal(err)
	}
	if *getTestAlert.Config.Host != *testConfig.Host {
		t.Errorf("unexpected workspace alert config host: got %v, expected %v", *getTestAlert.Config.Host, *testConfig.Host)
	}
	if *getTestAlert.Config.Key != *testConfig.Key {
		t.Errorf("unexpected workspace alert config key: got %v, expected %v", *getTestAlert.Config.Key, *testConfig.Key)
	}
	if *getTestAlert.Config.Project != *testConfig.Project {
		t.Errorf("unexpected workspace alert config project: got %v, expected %v", *getTestAlert.Config.Project, *testConfig.Project)
	}
	if *getTestAlert.Config.Username != *testConfig.Username {
		t.Errorf("unexpected workspace alert config username: got %v, expected %v", *getTestAlert.Config.Username, *testConfig.Username)
	}
	if *getTestAlert.Config.IssueType != *testConfig.IssueType {
		t.Errorf("unexpected workspace alert config issue type: got %v, expected %v", *getTestAlert.Config.IssueType, *testConfig.IssueType)
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
		Host:      fastly.ToPointer("https://example.exampleJira.net"),
		Key:       fastly.ToPointer("ATATT3xFfGF0b2c3d4e5f6789012345678901234567891"),
		Project:   fastly.ToPointer("UPDATED"),
		Username:  fastly.ToPointer("updateduser"),
		IssueType: fastly.ToPointer("Bug"),
	}
	updatedEvent := "flag"
	var updateAlert *Alert
	fastly.Record(t, "update_workspacealert", func(c *fastly.Client) {
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
	if *updateAlert.Config.Host != *updatedConfig.Host {
		t.Errorf("unexpected updated workspace alert config host: got %v, expected %v", *updateAlert.Config.Host, *updatedConfig.Host)
	}
	if *updateAlert.Config.Key != *updatedConfig.Key {
		t.Errorf("unexpected updated workspace alert config key: got %v, expected %v", *updateAlert.Config.Key, *updatedConfig.Key)
	}
	if *updateAlert.Config.Project != *updatedConfig.Project {
		t.Errorf("unexpected updated workspace alert config project: got %v, expected %v", *updateAlert.Config.Project, *updatedConfig.Project)
	}
	if *updateAlert.Config.Username != *updatedConfig.Username {
		t.Errorf("unexpected updated workspace alert config username: got %v, expected %v", *updateAlert.Config.Username, *updatedConfig.Username)
	}
	if *updateAlert.Config.IssueType != *updatedConfig.IssueType {
		t.Errorf("unexpected updated workspace alert config issue type: got %v, expected %v", *updateAlert.Config.IssueType, *updatedConfig.IssueType)
	}
	if len(updateAlert.Events) != 1 {
		t.Errorf("unexpected updated workspace alerts event length: got %d, expected %d", len(updateAlert.Events), 1)
	}
	if updateAlert.Events[0] != updatedEvent {
		t.Errorf("unexpected updated workspace alert events: got %+v, expected %+v", updateAlert.Events[0], updatedEvent)
	}

	// List the workspace alerts for the test workspace and check the updated one is the only entry.
	var Alerts *Alerts
	fastly.Record(t, "list_workspacealerts", func(c *fastly.Client) {
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
	if *listedAlert.Config.Host != *updatedConfig.Host {
		t.Errorf("unexpected listed workspace alert config host: got %v, expected %v", *listedAlert.Config.Host, *updatedConfig.Host)
	}
	if *listedAlert.Config.Key != *updatedConfig.Key {
		t.Errorf("unexpected listed workspace alert config key: got %v, expected %v", *listedAlert.Config.Key, *updatedConfig.Key)
	}
	if *listedAlert.Config.Project != *updatedConfig.Project {
		t.Errorf("unexpected listed workspace alert config project: got %v, expected %v", *listedAlert.Config.Project, *updatedConfig.Project)
	}
	if *listedAlert.Config.Username != *updatedConfig.Username {
		t.Errorf("unexpected listed workspace alert config username: got %v, expected %v", *listedAlert.Config.Username, *updatedConfig.Username)
	}
	if *listedAlert.Config.IssueType != *updatedConfig.IssueType {
		t.Errorf("unexpected listed workspace alert config issue type: got %v, expected %v", *listedAlert.Config.IssueType, *updatedConfig.IssueType)
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
		Config:      CreateConfig{Host: fastly.ToPointer("example.atlassian.net"), Key: fastly.ToPointer("test"), Project: fastly.ToPointer("TEST"), Username: fastly.ToPointer("user")},
		Events:      nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingEvents) {
		t.Errorf("expected ErrMissingEvents: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer(IntegrationType),
		Config:      CreateConfig{Key: fastly.ToPointer("test"), Project: fastly.ToPointer("TEST"), Username: fastly.ToPointer("user")},
		Events:      []string{"flag"},
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingHost) {
		t.Errorf("expected ErrMissingHost: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer(IntegrationType),
		Config:      CreateConfig{Host: fastly.ToPointer("example.atlassian.net"), Project: fastly.ToPointer("TEST"), Username: fastly.ToPointer("user")},
		Events:      []string{"flag"},
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingKey) {
		t.Errorf("expected ErrMissingKey: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer(IntegrationType),
		Config:      CreateConfig{Host: fastly.ToPointer("example.atlassian.net"), Key: fastly.ToPointer("test"), Username: fastly.ToPointer("user")},
		Events:      []string{"flag"},
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingProject) {
		t.Errorf("expected ErrMissingProject: got %s", err)
	}
	_, err = Create(fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer(IntegrationType),
		Config:      CreateConfig{Host: fastly.ToPointer("example.atlassian.net"), Key: fastly.ToPointer("test"), Project: fastly.ToPointer("TEST")},
		Events:      []string{"flag"},
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingUserName) {
		t.Errorf("expected ErrMissingUserName: got %s", err)
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
