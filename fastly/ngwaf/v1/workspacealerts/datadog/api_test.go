package datadog

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

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


