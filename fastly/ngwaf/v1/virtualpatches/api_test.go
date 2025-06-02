package virtualpatches

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Global value for the test workspace ID.
var testWorkspaceID = fastly.TestNGWAFWorkspaceID

// ID of Virtual Patch to test against.
const vpID = "CVE-2017-5638"

func TestVirtual_Patches(t *testing.T) {
	t.Parallel()

	var err error
	var vps *VirtualPatches

	// List all virtual patches.
	fastly.Record(t, "list_virtualpatches", func(c *fastly.Client) {
		vps, err = List(c, &ListInput{
			WorkspaceID: &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if vps == nil {
		t.Fatal("expected VirtualPatches response, got nil")
	}

	// Sample a few listed virtual patches.
	expectedPatches := map[string]string{
		"CVE-2017-5638":  "Apache Struts multipart/form remote execution",
		"CVE-2021-26855": "Microsoft Exchange Server Remote Code Execution Vulnerability",
		"CVE-2017-7269":  "IIS 6.0 WebDAV buffer overflow",
	}

	// Create a map for quick lookup of listed virtual patches.
	returnedPatches := make(map[string]VirtualPatch)
	for _, patch := range vps.Data {
		returnedPatches[patch.ID] = patch
	}

	// Virtual Patch sample validation.
	for expectedID, expectedDescription := range expectedPatches {
		patch, found := returnedPatches[expectedID]
		if !found {
			t.Errorf("expected virtual patch %q not found in response", expectedID)
			continue
		}

		if patch.Description != expectedDescription {
			t.Errorf("virtual patch %q: unexpected description: got %q, expected %q",
				expectedID, patch.Description, expectedDescription)
		}
	}

	// Get a virual patch.
	var vp *VirtualPatch
	fastly.Record(t, "get_virtualpatch", func(c *fastly.Client) {
		vp, err = Get(c, &GetInput{
			VirtualPatchID: fastly.ToPointer(vpID),
			WorkspaceID:    &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if vp == nil {
		t.Fatal("expected VirtualPatch response, got nil")
	}
	if vp.ID != vpID {
		t.Errorf("unexpected virtual patch ID: got %q, expected %q", vp.ID, vpID)
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
