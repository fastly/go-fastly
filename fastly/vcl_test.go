package fastly

import "testing"

func TestClient_VCLs(t *testing.T) {
	tv := testVersion(t)

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
	vcl, err := testClient.CreateVCL(&CreateVCLInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-vcl",
		Content: content,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteVCL(&DeleteVCLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-vcl",
		})

		testClient.DeleteVCL(&DeleteVCLInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-vcl",
		})
	}()

	if vcl.Name != "test-vcl" {
		t.Errorf("bad name: %q", vcl.Name)
	}
	if vcl.Content != content {
		t.Errorf("bad content: %q", vcl.Content)
	}

	// List
	vcls, err := testClient.ListVCLs(&ListVCLsInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(vcls) < 1 {
		t.Errorf("bad vcls: %v", vcls)
	}

	// Get
	nvcl, err := testClient.GetVCL(&GetVCLInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-vcl",
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
	uvcl, err := testClient.UpdateVCL(&UpdateVCLInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-vcl",
		NewName: "new-test-vcl",
	})
	if err != nil {
		t.Fatal(err)
	}
	if uvcl.Name != "new-test-vcl" {
		t.Errorf("bad name: %q", uvcl.Name)
	}

	// Activate
	avcl, err := testClient.ActivateVCL(&ActivateVCLInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "new-test-vcl",
	})
	if err != nil {
		t.Fatal(err)
	}
	if avcl.Main != true {
		t.Errorf("bad main: %b", avcl.Main)
	}

	// Delete
	if err := testClient.DeleteVCL(&DeleteVCLInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "new-test-vcl",
	}); err != nil {
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
		Version: "",
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
		Version: "",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetVCL(&GetVCLInput{
		Service: "foo",
		Version: "1",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateVCL(&UpdateVCLInput{
		Service: "foo",
		Version: "1",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ActivateVCL(&ActivateVCLInput{
		Service: "foo",
		Version: "1",
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
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteVCL(&DeleteVCLInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
