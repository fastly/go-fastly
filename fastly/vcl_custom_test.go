package fastly

import (
	"errors"
	"testing"
)

func TestClient_VCLs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "vcls/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	content := `
backend default {
  .host = "127.0.0.1";
  .port = "9092";
}

sub vcl_recv {
  set req.backend = default;

  if (req.url.path ~ "(1|2)") {
    // ...
  }
}

sub vcl_hash {
  set req.hash += req.url;
  set req.hash += req.http.host;
  set req.hash += "0";
}
`

	// Create
	var vcl *VCL
	Record(t, "vcls/create", func(c *Client) {
		vcl, err = c.CreateVCL(&CreateVCLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-vcl"),
			Content:        ToPointer(content),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "vcls/cleanup", func(c *Client) {
			_ = c.DeleteVCL(&DeleteVCLInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-vcl",
			})

			_ = c.DeleteVCL(&DeleteVCLInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-vcl",
			})
		})
	}()

	if *vcl.Name != "test-vcl" {
		t.Errorf("bad name: %q", *vcl.Name)
	}
	if *vcl.Content != content {
		t.Errorf("bad content: %q", *vcl.Content)
	}

	// List
	var vcls []*VCL
	Record(t, "vcls/list", func(c *Client) {
		vcls, err = c.ListVCLs(&ListVCLsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "vcls/get", func(c *Client) {
		nvcl, err = c.GetVCL(&GetVCLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-vcl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *vcl.Name != *nvcl.Name {
		t.Errorf("bad name: %q", *vcl.Name)
	}
	if *vcl.Content != *nvcl.Content {
		t.Errorf("bad address: %q", *vcl.Content)
	}

	// Update
	var uvcl *VCL
	Record(t, "vcls/update", func(c *Client) {
		uvcl, err = c.UpdateVCL(&UpdateVCLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-vcl",
			NewName:        ToPointer("new-test-vcl"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uvcl.Name != "new-test-vcl" {
		t.Errorf("bad name: %q", *uvcl.Name)
	}

	// Activate
	var avcl *VCL
	Record(t, "vcls/activate", func(c *Client) {
		avcl, err = c.ActivateVCL(&ActivateVCLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-vcl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !*avcl.Main {
		t.Errorf("bad main: %t", *avcl.Main)
	}

	// Delete
	Record(t, "vcls/delete", func(c *Client) {
		err = c.DeleteVCL(&DeleteVCLInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-vcl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListVCLs_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListVCLs(&ListVCLsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListVCLs(&ListVCLsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateVCL_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateVCL(&CreateVCLInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateVCL(&CreateVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetVCL_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetVCL(&GetVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetVCL(&GetVCLInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetVCL(&GetVCLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateVCL_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateVCL(&UpdateVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateVCL(&UpdateVCLInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateVCL(&UpdateVCLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ActivateVCL_validation(t *testing.T) {
	var err error

	_, err = TestClient.ActivateVCL(&ActivateVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ActivateVCL(&ActivateVCLInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ActivateVCL(&ActivateVCLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteVCL_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteVCL(&DeleteVCLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteVCL(&DeleteVCLInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteVCL(&DeleteVCLInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
