package fastly

import (
	"os"
	"testing"
)

func TestClient_Package(t *testing.T) {
	fixtureBase := "package/"
	nameSuffix := "package"

	testService := createTestServiceWasm(t, fixtureBase+"service_create", nameSuffix)
	testVersion := createTestVersion(t, fixtureBase+"service_version", *testService.ID)
	defer deleteTestService(t, fixtureBase+"service_delete", *testService.ID)

	testData := Package{
		Metadata: &PackageMetadata{
			Description: ToPointer("Default package template used by the Fastly CLI for Rust-based Compute@Edge projects."),
			HashSum:     ToPointer("f99485bd301e23f028474d26d398da525de17a372ae9e7026891d7f85361d2540d14b3b091929c3f170eade573595e20b3405a9e29651ede59915f2e1652f616"),
			Language:    ToPointer("rust"),
			Name:        ToPointer("wasm-test"),
			Size:        ToPointer(int64(2015936)),
			FilesHash:   ToPointer("a763d3c88968ebc17691900d3c14306762296df8e47a1c2d7661cee0e0c5aa6d4c082a7c128d6e719fe333b73b46fe3ae32694716ccd2efa21f5d9f049ceec6d"),
		},
	}

	var wp *Package
	var err error

	// Update with valid package file path

	recordIgnoreBody(t, fixtureBase+"update", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      *testService.ID,
			ServiceVersion: testVersion.Number,
			PackagePath:    ToPointer("test_assets/package/valid.tar.gz"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *wp.ServiceID != *testService.ID {
		t.Errorf("bad serviceID: %q != %q", *wp.ID, *testService.ID)
	}
	if *wp.ServiceVersion != testVersion.Number {
		t.Errorf("bad serviceVersion: %d != %d", *wp.ServiceVersion, testVersion.Number)
	}

	// Get
	record(t, fixtureBase+"get", func(c *Client) {
		wp, err = c.GetPackage(&GetPackageInput{
			ServiceID:      *testService.ID,
			ServiceVersion: testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *wp.ServiceID != *testService.ID {
		t.Errorf("bad serviceID: %q != %q", *wp.ID, *testService.ID)
	}
	if *wp.ServiceVersion != testVersion.Number {
		t.Errorf("bad serviceVersion: %d != %d", wp.ServiceVersion, testVersion.Number)
	}

	if *wp.Metadata.Name != *testData.Metadata.Name {
		t.Errorf("bad package name: %q != %q", *wp.Metadata.Name, *testData.Metadata.Name)
	}
	if *wp.Metadata.Description != *testData.Metadata.Description {
		t.Errorf("bad package description: %q != %q", *wp.Metadata.Description, *testData.Metadata.Description)
	}
	if *wp.Metadata.Size != *testData.Metadata.Size {
		t.Errorf("bad package size: %q != %q", *wp.Metadata.Size, *testData.Metadata.Size)
	}
	if *wp.Metadata.HashSum != *testData.Metadata.HashSum {
		t.Errorf("bad package hashsum: %q != %q", *wp.Metadata.HashSum, *testData.Metadata.HashSum)
	}
	if *wp.Metadata.FilesHash != *testData.Metadata.FilesHash {
		t.Errorf("bad package files_hash: %q != %q", *wp.Metadata.FilesHash, *testData.Metadata.FilesHash)
	}
	if *wp.Metadata.Language != *testData.Metadata.Language {
		t.Errorf("bad package language: %q != %q", *wp.Metadata.Language, *testData.Metadata.Language)
	}

	// Update with valid package bytes

	validPackageContent, _ := os.ReadFile("test_assets/package/valid.tar.gz")
	recordIgnoreBody(t, fixtureBase+"update", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      *testService.ID,
			ServiceVersion: testVersion.Number,
			PackageContent: validPackageContent,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *wp.ServiceID != *testService.ID {
		t.Errorf("bad serviceID: %q != %q", *wp.ID, *testService.ID)
	}
	if *wp.ServiceVersion != testVersion.Number {
		t.Errorf("bad serviceVersion: %d != %d", *wp.ServiceVersion, testVersion.Number)
	}

	// Get
	record(t, fixtureBase+"get", func(c *Client) {
		wp, err = c.GetPackage(&GetPackageInput{
			ServiceID:      *testService.ID,
			ServiceVersion: testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *wp.ServiceID != *testService.ID {
		t.Errorf("bad serviceID: %q != %q", *wp.ID, *testService.ID)
	}
	if *wp.ServiceVersion != testVersion.Number {
		t.Errorf("bad serviceVersion: %d != %d", wp.ServiceVersion, testVersion.Number)
	}

	if *wp.Metadata.Name != *testData.Metadata.Name {
		t.Errorf("bad package name: %q != %q", *wp.Metadata.Name, *testData.Metadata.Name)
	}
	if *wp.Metadata.Description != *testData.Metadata.Description {
		t.Errorf("bad package description: %q != %q", *wp.Metadata.Description, *testData.Metadata.Description)
	}
	if *wp.Metadata.Size != *testData.Metadata.Size {
		t.Errorf("bad package size: %q != %q", *wp.Metadata.Size, *testData.Metadata.Size)
	}
	if *wp.Metadata.HashSum != *testData.Metadata.HashSum {
		t.Errorf("bad package hashsum: %q != %q", *wp.Metadata.HashSum, *testData.Metadata.HashSum)
	}
	if *wp.Metadata.Language != *testData.Metadata.Language {
		t.Errorf("bad package language: %q != %q", *wp.Metadata.Language, *testData.Metadata.Language)
	}

	// Update with invalid package file path

	recordIgnoreBody(t, fixtureBase+"update_invalid", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      *testService.ID,
			ServiceVersion: testVersion.Number,
			PackagePath:    ToPointer("test_assets/package/invalid.tar.gz"),
		})
	})
	if err == nil && (wp.Metadata.Size != nil || wp.Metadata.Language != nil || wp.Metadata.HashSum != nil || wp.Metadata.Description != nil || wp.Metadata.Name != nil) {
		t.Fatal("Invalid package upload completed rather than failed.")
	}

	// Update with invalid package bytes

	invalidPackageContent, _ := os.ReadFile("test_assets/package/invalid.tar.gz")
	recordIgnoreBody(t, fixtureBase+"update_invalid", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      *testService.ID,
			ServiceVersion: testVersion.Number,
			PackageContent: invalidPackageContent,
		})
	})
	if err == nil && (wp.Metadata.Size != nil || wp.Metadata.Language != nil || wp.Metadata.HashSum != nil || wp.Metadata.Description != nil || wp.Metadata.Name != nil) {
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
