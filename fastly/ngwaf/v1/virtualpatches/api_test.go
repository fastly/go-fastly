package virtualpatches

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Global value for the test workspace ID.
var testWorkspaceID string

// ID of Virtual Patch to test against.
const vpID = "CVE-2017-5638"

func TestVirtual_Patches(t *testing.T) {
	t.Parallel()

	const wsName = "!!vp-workspace"
	const wsDescription = "vp-workspace"
	const wsMode = "block"

	// Create a test workspace.
	var ws *Workspace
	var err error

	fastly.Record(t, "create_workspace", func(c *fastly.Client) {
		ws, err = Create(c, &CreateWorkspace{
			Name:        fastly.ToPointer(wsName),
			Description: fastly.ToPointer(wsDescription),
			Mode:        fastly.ToPointer(wsMode),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Store the test workspace ID globally.
	testWorkspaceID = ws.WorkspaceID

	// Ensure we delete the test workspace at the end.
	defer func() {
		fastly.Record(t, "delete_workspace", func(c *fastly.Client) {
			err = Delete(c, &DeleteInput{
				WorkspaceID: &testWorkspaceID,
			})
		})
		if err != nil {
			t.Errorf("error during workspace cleanup: %v", err)
		}
	}()

	// List all virtual patches.
	fastly.Record(t, "list_virtualpatches", func(c *fastly.Client) {
		_, err = List(c, &ListInput{
			WorkspaceID: &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get a virual patch.
	var _ *VirtualPatch
	fastly.Record(t, "get_virtualpatch", func(c *fastly.Client) {
		_, err = Get(c, &GetInput{
			VirtualPatchID: fastly.ToPointer(vpID),
			WorkspaceID:    &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Update the virtual patch
	const uvpDescription = "Apache Struts multipart/form remote execution"
	const uvpEnabled = true
	const uvpMode = "block"

	var uvp *VirtualPatch
	fastly.Record(t, "update_virtualpatch", func(c *fastly.Client) {
		uvp, err = Update(c, &UpdateInput{
			Enabled:        fastly.ToPointer(uvpEnabled),
			Mode:           fastly.ToPointer(uvpMode),
			VirtualPatchID: fastly.ToPointer(vpID),
			WorkspaceID:    &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uvp.Enabled != uvpEnabled {
		t.Errorf("unexpected virtual patch status: got %t, expected %t", uvp.Enabled, uvpEnabled)
	}
	if uvp.Mode != uvpMode {
		t.Errorf("unexpected virtual patch mode: got %q, expected %q", uvp.Mode, uvpMode)
	}
	if uvp.ID != vpID {
		t.Errorf("unexpected virtual identifier: got %q, expected %q", uvp.ID, vpID)
	}
	if uvp.Description != uvpDescription {
		t.Errorf("unexpected virtual description: got %q, expected %q", uvp.Description, uvpDescription)
	}
}

func TestClient_GetVirtualPatch_validation(t *testing.T) {
	var err error
	_, err = Get(fastly.TestClient, &GetInput{
		VirtualPatchID: nil,
		WorkspaceID:    &testWorkspaceID,
	})
	if !errors.Is(err, fastly.ErrMissingVirtualPatchID) {
		t.Errorf("expected ErrMissingVirtualPatchID: got %s", err)
	}

	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID:    nil,
		VirtualPatchID: fastly.ToPointer(vpID),
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}

func TestClient_ListVirtualPatch_validation(t *testing.T) {
	var err error
	_, err = List(fastly.TestClient, &ListInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}

func TestClient_UpdateVirtualPatch_validation(t *testing.T) {
	var err error
	_, err = Update(fastly.TestClient, &UpdateInput{
		VirtualPatchID: nil,
		WorkspaceID:    &testWorkspaceID,
	})
	if !errors.Is(err, fastly.ErrMissingVirtualPatchID) {
		t.Errorf("expected ErrMissingVirtualPatchID: got %s", err)
	}

	_, err = Update(fastly.TestClient, &UpdateInput{
		WorkspaceID:    nil,
		VirtualPatchID: fastly.ToPointer(vpID),
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}
