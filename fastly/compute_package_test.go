package fastly

import (
	"errors"
	"os"
	"testing"
)

func TestClient_Package(t *testing.T) {
	fixtureBase := "package/"
	nameSuffix := "package"

	testService := createTestServiceWasm(t, fixtureBase+"service_create", nameSuffix)
	testVersion := CreateTestVersion(t, fixtureBase+"service_version", *testService.ServiceID)
	defer deleteTestService(t, fixtureBase+"service_delete", *testService.ServiceID)

	testData := Package{
		Metadata: &PackageMetadata{
			ClonedFrom:  ToPointer("https://github.com/fastly/compute-starter-kit-rust-empty"),
			Description: ToPointer("An empty starter kit project template."),
			FilesHash:   ToPointer("75ff1cf4d953ff2242bb38e4a01b04503622baf4b7dc540256f4dd5fc89df5aed7fea115adab0b71caa79f6483bb846ac0d4f4f937885fb03ee35d2dfafba6f3"),
			HashSum:     ToPointer("ecc068efcd4071d36d6460152dcc50461b649f01f28589917540a69d33e9b1477decb5f5a9a2f6c269d83a13827502c1fe1f2efc3bdd6beadaabd23e22eb84fd"),
			Language:    ToPointer("rust"),
			Name:        ToPointer("test-package"),
			Size:        ToPointer(int64(1540845)),
		},
	}

	var wp *Package
	var err error

	// Update with valid package file path

	RecordIgnoreBody(t, fixtureBase+"update", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      *testService.ServiceID,
			ServiceVersion: *testVersion.Number,
			PackagePath:    ToPointer("test_assets/package/valid.tar.gz"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *wp.ServiceID != *testService.ServiceID {
		t.Errorf("bad serviceID: %q != %q", *wp.PackageID, *testService.ServiceID)
	}
	if *wp.ServiceVersion != *testVersion.Number {
		t.Errorf("bad serviceVersion: %d != %d", *wp.ServiceVersion, testVersion.Number)
	}

	// Get
	Record(t, fixtureBase+"get", func(c *Client) {
		wp, err = c.GetPackage(&GetPackageInput{
			ServiceID:      *testService.ServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *wp.ServiceID != *testService.ServiceID {
		t.Errorf("bad serviceID: %q != %q", *wp.PackageID, *testService.ServiceID)
	}
	if *wp.ServiceVersion != *testVersion.Number {
		t.Errorf("bad serviceVersion: %d != %d", wp.ServiceVersion, testVersion.Number)
	}

	if *wp.Metadata.Name != *testData.Metadata.Name {
		t.Errorf("bad package name: %q != %q", *wp.Metadata.Name, *testData.Metadata.Name)
	}
	if *wp.Metadata.Description != *testData.Metadata.Description {
		t.Errorf("bad package description: %q != %q", *wp.Metadata.Description, *testData.Metadata.Description)
	}
	if *wp.Metadata.Size != *testData.Metadata.Size {
		t.Errorf("bad package size: %d != %d", *wp.Metadata.Size, *testData.Metadata.Size)
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
	if ToValue(wp.Metadata.ClonedFrom) != ToValue(testData.Metadata.ClonedFrom) {
		t.Errorf("bad package cloned_from: %q != %q", ToValue(wp.Metadata.ClonedFrom), ToValue(testData.Metadata.ClonedFrom))
	}

	// Update with valid package bytes

	validPackageContent, _ := os.ReadFile("test_assets/package/valid.tar.gz")
	RecordIgnoreBody(t, fixtureBase+"update", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      *testService.ServiceID,
			ServiceVersion: *testVersion.Number,
			PackageContent: validPackageContent,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *wp.ServiceID != *testService.ServiceID {
		t.Errorf("bad serviceID: %q != %q", *wp.PackageID, *testService.ServiceID)
	}
	if *wp.ServiceVersion != *testVersion.Number {
		t.Errorf("bad serviceVersion: %d != %d", *wp.ServiceVersion, testVersion.Number)
	}

	// Get
	Record(t, fixtureBase+"get", func(c *Client) {
		wp, err = c.GetPackage(&GetPackageInput{
			ServiceID:      *testService.ServiceID,
			ServiceVersion: *testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *wp.ServiceID != *testService.ServiceID {
		t.Errorf("bad serviceID: %q != %q", *wp.PackageID, *testService.ServiceID)
	}
	if *wp.ServiceVersion != *testVersion.Number {
		t.Errorf("bad serviceVersion: %d != %d", wp.ServiceVersion, testVersion.Number)
	}

	if *wp.Metadata.Name != *testData.Metadata.Name {
		t.Errorf("bad package name: %q != %q", *wp.Metadata.Name, *testData.Metadata.Name)
	}
	if *wp.Metadata.Description != *testData.Metadata.Description {
		t.Errorf("bad package description: %q != %q", *wp.Metadata.Description, *testData.Metadata.Description)
	}
	if *wp.Metadata.Size != *testData.Metadata.Size {
		t.Errorf("bad package size: %d != %d", *wp.Metadata.Size, *testData.Metadata.Size)
	}
	if *wp.Metadata.HashSum != *testData.Metadata.HashSum {
		t.Errorf("bad package hashsum: %q != %q", *wp.Metadata.HashSum, *testData.Metadata.HashSum)
	}
	if *wp.Metadata.Language != *testData.Metadata.Language {
		t.Errorf("bad package language: %q != %q", *wp.Metadata.Language, *testData.Metadata.Language)
	}
	if ToValue(wp.Metadata.ClonedFrom) != ToValue(testData.Metadata.ClonedFrom) {
		t.Errorf("bad package cloned_from: %q != %q", ToValue(wp.Metadata.ClonedFrom), ToValue(testData.Metadata.ClonedFrom))
	}

	// Update with invalid package file path

	RecordIgnoreBody(t, fixtureBase+"update_invalid", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      *testService.ServiceID,
			ServiceVersion: *testVersion.Number,
			PackagePath:    ToPointer("test_assets/package/invalid.tar.gz"),
		})
	})
	if err == nil && (wp.Metadata.Size != nil || wp.Metadata.Language != nil || wp.Metadata.HashSum != nil || wp.Metadata.Description != nil || wp.Metadata.Name != nil) {
		t.Fatal("Invalid package upload completed rather than failed.")
	}

	// Update with invalid package bytes

	invalidPackageContent, _ := os.ReadFile("test_assets/package/invalid.tar.gz")
	RecordIgnoreBody(t, fixtureBase+"update_invalid", func(c *Client) {
		wp, err = c.UpdatePackage(&UpdatePackageInput{
			ServiceID:      *testService.ServiceID,
			ServiceVersion: *testVersion.Number,
			PackageContent: invalidPackageContent,
		})
	})
	if err == nil && (wp.Metadata.Size != nil || wp.Metadata.Language != nil || wp.Metadata.HashSum != nil || wp.Metadata.Description != nil || wp.Metadata.Name != nil) {
		t.Fatal("Invalid package upload completed rather than failed.")
	}
}

func TestClient_GetPackage_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetPackage(&GetPackageInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetPackage(&GetPackageInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePackage_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdatePackage(&UpdatePackageInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdatePackage(&UpdatePackageInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
