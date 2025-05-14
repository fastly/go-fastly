package workspaces

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

func TestClient_Workspace(t *testing.T) {
	t.Parallel()

	const wsName = "test-workspace"
	const wsDescription = "test-description"
	const wsMode = "log"
	const wsIPAnonymization = "hashed"

	wsAttackSignalThresholds := new(AttackSignalThresholdsInput)
	wsAttackSignalThresholds.OneMinute = fastly.ToPointer(10000)
	wsAttackSignalThresholds.TenMinutes = fastly.ToPointer(10000)
	wsAttackSignalThresholds.OneHour = fastly.ToPointer(10000)
	wsAttackSignalThresholds.Immediate = fastly.ToPointer(true)

	var wss *Workspaces
	var err error

	// List all workspaces.
	fastly.Record(t, "list_workspaces", func(c *fastly.Client) {
		wss, err = List(c, &ListInput{})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the test workspace we're going to create isn't among them.
	for _, ws := range wss.Data {
		if ws.Name == wsName {
			t.Errorf("found test workspace %q, aborting", wsName)
		}
	}

	// Create a test workspace.
	var ws *Workspace
	fastly.Record(t, "create_workspace", func(c *fastly.Client) {
		ws, err = Create(c, &CreateInput{
			Name:                   fastly.ToPointer(wsName),
			Description:            fastly.ToPointer(wsDescription),
			Mode:                   fastly.ToPointer(wsMode),
			IPAnonymization:        fastly.ToPointer(wsIPAnonymization),
			AttackSignalThresholds: wsAttackSignalThresholds,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ws.Name != wsName {
		t.Errorf("unexpected workspace name: got %q, expected %q", ws.Name, wsName)
	}
	if ws.Description != wsDescription {
		t.Errorf("unexpected workspace description: got %q, expected %q", ws.Description, wsDescription)
	}
	if ws.Mode != wsMode {
		t.Errorf("unexpected workspace mode: got %q, expected %q", ws.Mode, wsMode)
	}
	if ws.IPAnonymization != wsIPAnonymization {
		t.Errorf("unexpected workspace IP anonymization: got %q, expected %q", ws.IPAnonymization, wsIPAnonymization)
	}
	// TODO: Update the following assertions for the AttackSignalThresholds settings once CDTOOL-1072 has been resolved and deployed.
	// Note: At present, the API does not honor the AttackSignalThresholds settings when making a POST request to /ngwaf/v1/workspaces.
	if ws.AttackSignalThresholds.Immediate != false {
		t.Errorf("unexpected workspace attack signal thresholds immediate parameter: got %t, expected %t", ws.AttackSignalThresholds.Immediate, *wsAttackSignalThresholds.Immediate)
	}
	if ws.AttackSignalThresholds.OneMinute != 0 {
		t.Errorf("unexpected workspace attack signal thresholds one_minute parameter: got %v, expected %v", ws.AttackSignalThresholds.OneMinute, *wsAttackSignalThresholds.OneMinute)
	}
	if ws.AttackSignalThresholds.TenMinutes != 0 {
		t.Errorf("unexpected workspace attack signal thresholds ten_minutes parameter: got %v, expected %v", ws.AttackSignalThresholds.TenMinutes, *wsAttackSignalThresholds.TenMinutes)
	}
	if ws.AttackSignalThresholds.OneHour != 0 {
		t.Errorf("unexpected workspace attack signal thresholds one_hour parameter: got %v, expected %v", ws.AttackSignalThresholds.OneHour, *wsAttackSignalThresholds.OneHour)
	}

	// Ensure we delete the test workspace at the end.
	defer func() {
		fastly.Record(t, "delete_workspace", func(c *fastly.Client) {
			err = Delete(c, &DeleteInput{
				WorkspaceID: fastly.ToPointer(ws.WorkspaceID),
			})
		})
		if err != nil {
			t.Errorf("error during workspace cleanup: %v", err)
		}
	}()

	// Get the test workspace.
	var gws *Workspace
	fastly.Record(t, "get_workspace", func(c *fastly.Client) {
		gws, err = Get(c, &GetInput{
			WorkspaceID: fastly.ToPointer(ws.WorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gws.Name != ws.Name {
		t.Errorf("unexpected workspace name: got %q, expected %q", gws.Name, ws.Name)
	}
	if gws.Description != ws.Description {
		t.Errorf("unexpected workspace description: got %q, expected %q", gws.Description, ws.Description)
	}
	if gws.Mode != ws.Mode {
		t.Errorf("unexpected workspace mode: got %q, expected %q", gws.Mode, ws.Mode)
	}
	if gws.IPAnonymization != ws.IPAnonymization {
		t.Errorf("unexpected workspace IP anonymization: got %q, expected %q", gws.IPAnonymization, ws.IPAnonymization)
	}
	if gws.AttackSignalThresholds.Immediate != ws.AttackSignalThresholds.Immediate {
		t.Errorf("unexpected workspace attack signal thresholds immediate parameter: got %t, expected %t", gws.AttackSignalThresholds.Immediate, ws.AttackSignalThresholds.Immediate)
	}
	if gws.AttackSignalThresholds.OneMinute != ws.AttackSignalThresholds.OneMinute {
		t.Errorf("unexpected workspace attack signal thresholds one_minute parameter: got %v, expected %v", gws.AttackSignalThresholds.OneMinute, ws.AttackSignalThresholds.OneMinute)
	}
	if gws.AttackSignalThresholds.TenMinutes != ws.AttackSignalThresholds.TenMinutes {
		t.Errorf("unexpected workspace attack signal thresholds ten_minutes parameter: got %v, expected %v", gws.AttackSignalThresholds.TenMinutes, ws.AttackSignalThresholds.TenMinutes)
	}
	if gws.AttackSignalThresholds.OneHour != ws.AttackSignalThresholds.OneHour {
		t.Errorf("unexpected workspace attack signal thresholds one_hour parameter: got %v, expected %v", gws.AttackSignalThresholds.OneHour, ws.AttackSignalThresholds.OneHour)
	}

	// Update the test workspace.
	const uwsName = "test-workspace"
	const uwsDescription = "test-description"
	const uwsMode = "log"
	const uwsIPAnonymization = "hashed"

	uwsAttackSignalThresholds := new(AttackSignalThresholdsInput)
	uwsAttackSignalThresholds.OneMinute = fastly.ToPointer(10000)
	uwsAttackSignalThresholds.TenMinutes = fastly.ToPointer(10000)
	uwsAttackSignalThresholds.OneHour = fastly.ToPointer(10000)
	uwsAttackSignalThresholds.Immediate = fastly.ToPointer(true)

	var uws *Workspace
	fastly.Record(t, "update_workspace", func(c *fastly.Client) {
		uws, err = Update(c, &UpdateInput{
			WorkspaceID:            fastly.ToPointer(ws.WorkspaceID),
			Name:                   fastly.ToPointer(uwsName),
			Description:            fastly.ToPointer(uwsDescription),
			Mode:                   fastly.ToPointer(uwsMode),
			IPAnonymization:        fastly.ToPointer(uwsIPAnonymization),
			AttackSignalThresholds: uwsAttackSignalThresholds,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uws.Name != uwsName {
		t.Errorf("unexpected workspace name: got %q, expected %q", uws.Name, uwsName)
	}
	if uws.Description != uwsDescription {
		t.Errorf("unexpected workspace description: got %q, expected %q", uws.Description, uwsDescription)
	}
	if uws.Mode != uwsMode {
		t.Errorf("unexpected workspace mode: got %q, expected %q", uws.Mode, uwsMode)
	}
	if uws.IPAnonymization != uwsIPAnonymization {
		t.Errorf("unexpected workspace IP anonymization: got %q, expected %q", uws.IPAnonymization, uwsIPAnonymization)
	}
	if uws.AttackSignalThresholds.Immediate != *uwsAttackSignalThresholds.Immediate {
		t.Errorf("unexpected workspace attack signal thresholds immediate parameter: got %t, expected %t", uws.AttackSignalThresholds.Immediate, *uwsAttackSignalThresholds.Immediate)
	}
	if uws.AttackSignalThresholds.OneMinute != *uwsAttackSignalThresholds.OneMinute {
		t.Errorf("unexpected workspace attack signal thresholds one_minute parameter: got %v, expected %v", uws.AttackSignalThresholds.OneMinute, *uwsAttackSignalThresholds.OneMinute)
	}
	if uws.AttackSignalThresholds.TenMinutes != *uwsAttackSignalThresholds.TenMinutes {
		t.Errorf("unexpected workspace attack signal thresholds ten_minutes parameter: got %v, expected %v", uws.AttackSignalThresholds.TenMinutes, *uwsAttackSignalThresholds.TenMinutes)
	}
	if uws.AttackSignalThresholds.OneHour != *uwsAttackSignalThresholds.OneHour {
		t.Errorf("unexpected workspace attack signal thresholds one_hour parameter: got %v, expected %v", uws.AttackSignalThresholds.OneHour, *uwsAttackSignalThresholds.OneHour)
	}
}

func TestClient_GetWorkspace_validation(t *testing.T) {
	var err error
	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}

func TestClient_UpdateWorkspace_validation(t *testing.T) {
	var err error
	_, err = Update(fastly.TestClient, &UpdateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}

func TestClient_DeleteWorkspace_validation(t *testing.T) {
	err := Delete(fastly.TestClient, &DeleteInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}
