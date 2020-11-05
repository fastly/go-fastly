package fastly

import (
	"testing"
)

func TestClient_Package(t *testing.T) {

	fixtureBase := "package/"
	nameSuffix := "package"

	testService := createTestServiceWasm(t, fixtureBase+"service_create", nameSuffix)
	testVersion := createTestVersion(t, fixtureBase+"service_version", testService.ID)
	defer deleteTestService(t, fixtureBase+"service_delete", testService.ID)

	var testData = Package{
		Metadata: PackageMetadata{
			Name:        "wasm-test",
			Description: "Default package template used by the Fastly CLI for Rust-based Compute@Edge projects.",
			Language:    "rust",
			Size:        2015936,
			HashSum:     "f99485bd301e23f028474d26d398da525de17a372ae9e7026891d7f85361d2540d14b3b091929c3f170eade573595e20b3405a9e29651ede59915f2e1652f616",
		},
	}

	var wp *Package
	var err error

	// Update

	recordIgnoreBody(t, fixtureBase+"update", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      testService.ID,
			ServiceVersion: testVersion.Number,
			PackagePath:    "test_assets/package/valid.tar.gz",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if wp.ServiceID != testService.ID {
		t.Errorf("bad serviceID: %q != %q", wp.ID, testService.ID)
	}
	if wp.Version != testVersion.Number {
		t.Errorf("bad serviceID: %q != %q", wp.ID, testService.ID)
	}

	// Get
	record(t, fixtureBase+"get", func(c *Client) {
		wp, err = c.GetPackage(&GetPackageInput{
			ServiceID:      testService.ID,
			ServiceVersion: testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if wp.ServiceID != testService.ID {
		t.Errorf("bad serviceID: %q != %q", wp.ID, testService.ID)
	}
	if wp.Version != testVersion.Number {
		t.Errorf("bad serviceID: %q != %q", wp.ID, testService.ID)
	}

	if wp.Metadata.Name != testData.Metadata.Name {
		t.Errorf("bad package name: %q != %q", wp.Metadata.Name, testData.Metadata.Name)
	}
	if wp.Metadata.Description != testData.Metadata.Description {
		t.Errorf("bad package description: %q != %q", wp.Metadata.Description, testData.Metadata.Description)
	}
	if wp.Metadata.Size != testData.Metadata.Size {
		t.Errorf("bad package size: %q != %q", wp.Metadata.Size, testData.Metadata.Size)
	}
	if wp.Metadata.HashSum != testData.Metadata.HashSum {
		t.Errorf("bad package hashsum: %q != %q", wp.Metadata.HashSum, testData.Metadata.HashSum)
	}
	if wp.Metadata.Language != testData.Metadata.Language {
		t.Errorf("bad package language: %q != %q", wp.Metadata.Language, testData.Metadata.Language)
	}

	// Update with invalid package

	recordIgnoreBody(t, fixtureBase+"update_invalid", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      testService.ID,
			ServiceVersion: testVersion.Number,
			PackagePath:    "test_assets/package/invalid.tar.gz",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if wp.Metadata.Size > 0 || wp.Metadata.Language != "" || wp.Metadata.HashSum != "" || wp.Metadata.Description != "" ||
		wp.Metadata.Name != "" {
		t.Fatal("Invalid package upload completed rather than failed.")
	}

}

func TestClient_GetPackage_validation(t *testing.T) {
	var err error
	_, err = testClient.GetPackage(&GetPackageInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetPackage(&GetPackageInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePackage_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdatePackage(&UpdatePackageInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdatePackage(&UpdatePackageInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
