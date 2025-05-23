package virtualpatches

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

func TestVirtual_Patches(t *testing.T) {
	t.Parallel()

	const vpDescription = "Apache Struts multipart/form remote execution"
	const vpEnabled = "true"
	const vpID = "CVE-2017-5638"
	const vpMode = "block"
	const vpWorkspaceID = "S7ql67y0WTogAAWhhvMlF7"

	var err error

	// List all virtual patches.
	fastly.Record(t, "list_virtualpatches", func(c *fastly.Client) {
		_, err = List(c, &ListInput{
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get a virual patch.
	var gvp *VirtualPatch
	fastly.Record(t, "get_virtualpatch", func(c *fastly.Client) {
		gvp, err = Get(c, &GetInput{
			VirtualPatchID: fastly.ToPointer(vpID),
			WorkspaceID:    fastly.ToPointer(vpWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gvp.Description != vpDescription {
		t.Errorf("unexpected workspace description: got %q, expected %q", gvp.Description, vpDescription)
	}
	if gvp.Enabled != vpEnabled {
		t.Errorf("unexpected virtual patch status: got %q, expected %q", gvp.Enabled, vpEnabled)
	}
	if gvp.ID != vpID {
		t.Errorf("unexpected virtual patch ID: got  %q, expected %q", gvp.ID, vpID)
	}
	if gvp.Mode != vpMode {
		t.Errorf("unexpected virtual patch mode: got  %q, expected %q", gvp.Mode, vpMode)
	}

	// Update the virtual patch
	const uvpDescription = "Apache Struts multipart/form remote execution"
	const uvpEnabled = "false"
	const uvpID = "CVE-2017-5638"
	const uvpMode = "log"
	const uvpWorkspaceID = "S7ql67y0WTogAAWhhvMlF7"

	var uvp *VirtualPatch
	fastly.Record(t, "update_virtualpatch", func(c *fastly.Client) {
		uvp, err = Update(c, &UpdateInput{
			Enabled:        fastly.ToPointer(uvp.Enabled),
			Mode:           fastly.ToPointer(uvpMode),
			VirtualPatchID: fastly.ToPointer(uvpID),
			WorkspaceID:    fastly.ToPointer(uvpWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uvp.Enabled != uvpEnabled {
		t.Errorf("unexpected virtual patch status: got %q, expected %q", uvp.Enabled, uvpEnabled)
	}
	if uvp.Mode != uvpMode {
		t.Errorf("unexpected virtual patch mode: got %q, expected %q", uvp.Mode, uvpMode)
	}
	if uvp.ID != uvpID {
		t.Errorf("unexpected virtual identifier: got %q, expected %q", uvp.ID, uvpID)
	}
	if uvp.Description != uvpDescription {
		t.Errorf("unexpected virtual description: got %q, expected %q", uvp.Description, uvpDescription)
	}
}

func TestClient_GetVirtualPatch_validation(t *testing.T) {
	var err error
	_, err = Get(fastly.TestClient, &GetInput{
		VirtualPatchID: nil,
		WorkspaceID:    fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingVirtualPatchID) {
		t.Errorf("expected ErrMissingVirtualPatchID: got %s", err)
	}

	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID:    nil,
		VirtualPatchID: fastly.ToPointer(fastly.TestingNGWAFVirtualPatchID),
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
		WorkspaceID:    fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingVirtualPatchID) {
		t.Errorf("expected ErrMissingVirtualPatchID: got %s", err)
	}

	_, err = Update(fastly.TestClient, &UpdateInput{
		WorkspaceID:    nil,
		VirtualPatchID: fastly.ToPointer(fastly.TestingNGWAFVirtualPatchID),
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}
