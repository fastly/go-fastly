package fastly

import (
	"testing"
)

func TestClient_Wasm(t *testing.T) {

	fixtureBase := "wasm_package/"
	nameSuffix := "wasmPackage"

	testService := createTestServiceWasm(t, fixtureBase+"service_create", nameSuffix)
	testVersion := createTestVersion(t, fixtureBase+"service_version", testService.ID)
	defer deleteTestService(t, fixtureBase+"service_delete", testService.ID)

	var testData = WasmPackage{
		Metadata: WasmPackageMetadata{
			Name:        "wasm-test",
			Description: "Default package template used by the Fastly CLI for Rust-based Compute@Edge projects.",
			Language:    "rust",
			Size:        2179443,
			HashSum:     "119a5ac6ec8a6dc29107b87de3bba5d98b6e89827310dd01f4e8cbf89234125bf4f7b81603ac7df42ce0205d047ac2caa53ec0cc52839112aed786cb7f27ac92",
		},
	}

	var wp *WasmPackage
	var err error

	// Update

	record(t, fixtureBase+"update", func(c *Client) {
		wp, err = c.UpdateWasmPackage(&UpdateWasmPackageInput{
			Service:     testService.ID,
			Version:     testVersion.Number,
			PackagePath: "test_assets/wasm/valid.tar.gz",
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
		wp, err = c.GetWasmPackage(&GetWasmPackageInput{
			Service: testService.ID,
			Version: testVersion.Number,
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

	record(t, fixtureBase+"update_invalid", func(c *Client) {
		wp, err = c.UpdateWasmPackage(&UpdateWasmPackageInput{
			Service:     testService.ID,
			Version:     testVersion.Number,
			PackagePath: "test_assets/wasm/invalid.tar.gz",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if wp.Metadata.Size > 0 || wp.Metadata.Language!="" || wp.Metadata.HashSum!="" || wp.Metadata.Description!="" ||
	   wp.Metadata.Name!="" {
		t.Fatal("Invalid package upload completed rather than failed.")
	}

}



func TestClient_GetWasmPackage_validation(t *testing.T) {
	var err error
	_, err = testClient.GetWasmPackage(&GetWasmPackageInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetWasmPackage(&GetWasmPackageInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWasmPackage_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateWasmPackage(&UpdateWasmPackageInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateWasmPackage(&UpdateWasmPackageInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}