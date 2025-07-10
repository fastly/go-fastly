package jira

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

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
