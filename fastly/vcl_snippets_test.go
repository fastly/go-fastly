package fastly

import "testing"

func TestClient_VCLSnippets(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "vcl_snippets/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	vclSnippetName := "test-vcl-snippet"
	content := `
if ( req.url ) {
    set req.http.my-snippet-test-header = "true";
}
`

	// Create
	var vclSnippet *VCLSnippet
	record(t, "vcl_snippets/create", func(c *Client) {
		vclSnippet, err = c.CreateVCLSnippet(&CreateVCLSnippetInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    vclSnippetName,
			Content: content,
			Type:    "recv",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "vcl_snippets/cleanup", func(c *Client) {
			c.DeleteVCLSnippet(&DeleteVCLSnippetInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    vclSnippetName,
			})
		})
	}()

	if vclSnippet.Name != vclSnippetName {
		t.Errorf("bad name: %q", vclSnippet.Name)
	}
	if vclSnippet.Content != content {
		t.Errorf("bad content: %q", vclSnippet.Content)
	}

	// List
	var vclSnippets []*VCLSnippet
	record(t, "vcl_snippets/list", func(c *Client) {
		vclSnippets, err = c.ListVCLSnippets(&ListVCLSnippetsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(vclSnippets) < 1 {
		t.Errorf("bad vclSnippets: %v", vclSnippets)
	}

	// Get
	var nvclSnippet *VCLSnippet
	record(t, "vcl_snippets/get", func(c *Client) {
		nvclSnippet, err = c.GetVCLSnippet(&GetVCLSnippetInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    vclSnippetName,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if nvclSnippet.Name != nvclSnippet.Name {
		t.Errorf("bad name: %q", nvclSnippet.Name)
	}
	if nvclSnippet.Content != nvclSnippet.Content {
		t.Errorf("bad content: %q", nvclSnippet.Content)
	}

	// Update
	var uvclSnippet *VCLSnippet
	record(t, "vcl_snippets/update_dynamic", func(c *Client) {
		uvclSnippet, err = c.UpdateDynamicVCLSnippet(&UpdateDynamicVCLSnippetInput{
			Service: testServiceID,
			Name:    vclSnippetName,
			Content: "",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uvclSnippet.Name != vclSnippetName {
		t.Errorf("bad name: %q", uvclSnippet.Name)
	}

	// Update
	uvclSnippet = nil
	record(t, "vcl_snippets/update_regular", func(c *Client) {
		uvclSnippet, err = c.UpdateVCLSnippetName(&UpdateVCLSnippetInput{
			Service: testServiceID,
			Version: tv.Number,
			OldName: vclSnippetName,
			NewName: vclSnippetName + "-new",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uvclSnippet.Name != vclSnippetName+"-new" {
		t.Errorf("bad name: %q", uvclSnippet.Name)
	}

	// Delete
	record(t, "vcl_snippets/delete", func(c *Client) {
		err = c.DeleteVCLSnippet(&DeleteVCLSnippetInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    vclSnippetName,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListVCLSnippets_validation(t *testing.T) {
	var err error
	_, err = testClient.ListVCLs(&ListVCLsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListVCLs(&ListVCLsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateVCLSnippet_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateVCL(&CreateVCLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateVCL(&CreateVCLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetVCLSnippet_validation(t *testing.T) {
	var err error
	_, err = testClient.GetVCL(&GetVCLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVCL(&GetVCLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVCL(&GetVCLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetVCLSnippetByID_validation(t *testing.T) {
	var err error
	_, err = testClient.GetVCL(&GetVCLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVCL(&GetVCLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVCL(&GetVCLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDynamicVCLSnippet_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateVCLSnippet_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteVCLSnippet_validation(t *testing.T) {
	var err error
	err = testClient.DeleteVCL(&DeleteVCLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteVCL(&DeleteVCLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteVCL(&DeleteVCLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
