package fastly

import (
	"testing"
)

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
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-vcl",
			Content:        content,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "vcls/cleanup", func(c *Client) {
			c.DeleteVCL(&DeleteVCLInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-vcl",
			})

			c.DeleteVCL(&DeleteVCLInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-vcl",
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
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
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
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-vcl",
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
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-vcl",
			NewName:        String("new-test-vcl"),
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
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-vcl",
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
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-vcl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListVCLs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListVCLs(&ListVCLsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListVCLs(&ListVCLsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateVCL_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateVCL(&CreateVCLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateVCL(&CreateVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetVCL_validation(t *testing.T) {
	var err error
	_, err = testClient.GetVCL(&GetVCLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVCL(&GetVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVCL(&GetVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateVCL_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ActivateVCL_validation(t *testing.T) {
	var err error
	_, err = testClient.ActivateVCL(&ActivateVCLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ActivateVCL(&ActivateVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ActivateVCL(&ActivateVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteVCL_validation(t *testing.T) {
	var err error
	err = testClient.DeleteVCL(&DeleteVCLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteVCL(&DeleteVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteVCL(&DeleteVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
