package fastly

import "testing"

func TestClient_VCLs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "vcls/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	content := `
backend default {
  .host = "127.0.0.1";
  .port = "9092";
}

sub vcl_recv {
  set req.backend = default;
}

sub vcl_hash {
  set req.hash += req.url;
  set req.hash += req.http.host;
  set req.hash += "0";
}
`

	// Create
	var vcl *VCL
	record(t, "vcls/create", func(c *Client) {
		vcl, err = c.CreateVCL(&CreateVCLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-vcl",
			Content: content,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "vcls/cleanup", func(c *Client) {
			c.DeleteVCL(&DeleteVCLInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-vcl",
			})

			c.DeleteVCL(&DeleteVCLInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-vcl",
			})
		})
	}()

	if vcl.Name != "test-vcl" {
		t.Errorf("bad name: %q", vcl.Name)
	}
	if vcl.Content != content {
		t.Errorf("bad content: %q", vcl.Content)
	}

	// List
	var vcls []*VCL
	record(t, "vcls/list", func(c *Client) {
		vcls, err = c.ListVCLs(&ListVCLsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(vcls) < 1 {
		t.Errorf("bad vcls: %v", vcls)
	}

	// Get
	var nvcl *VCL
	record(t, "vcls/get", func(c *Client) {
		nvcl, err = c.GetVCL(&GetVCLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-vcl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if vcl.Name != nvcl.Name {
		t.Errorf("bad name: %q", vcl.Name)
	}
	if vcl.Content != nvcl.Content {
		t.Errorf("bad address: %q", vcl.Content)
	}

	// Update
	var uvcl *VCL
	record(t, "vcls/update", func(c *Client) {
		uvcl, err = c.UpdateVCL(&UpdateVCLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-vcl",
			NewName: "new-test-vcl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uvcl.Name != "new-test-vcl" {
		t.Errorf("bad name: %q", uvcl.Name)
	}

	// Activate
	var avcl *VCL
	record(t, "vcls/activate", func(c *Client) {
		avcl, err = c.ActivateVCL(&ActivateVCLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-vcl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if avcl.Main != true {
		t.Errorf("bad main: %t", avcl.Main)
	}

	// Delete
	record(t, "vcls/delete", func(c *Client) {
		err = c.DeleteVCL(&DeleteVCLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-vcl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListVCLs_validation(t *testing.T) {
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

func TestClient_CreateVCL_validation(t *testing.T) {
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

func TestClient_GetVCL_validation(t *testing.T) {
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

func TestClient_UpdateVCL_validation(t *testing.T) {
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

func TestClient_ActivateVCL_validation(t *testing.T) {
	var err error
	_, err = testClient.ActivateVCL(&ActivateVCLInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ActivateVCL(&ActivateVCLInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ActivateVCL(&ActivateVCLInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteVCL_validation(t *testing.T) {
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
