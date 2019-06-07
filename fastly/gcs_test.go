package fastly

import "testing"

func TestClient_GCSs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "gcses/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var gcs *GCS
	record(t, "gcses/create", func(c *Client) {
		gcs, err = c.CreateGCS(&CreateGCSInput{
			Service:         testServiceID,
			Version:         tv.Number,
			Name:            "test-gcs",
			Bucket:          "bucket",
			User:            "user",
			SecretKey:       "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9aQoqdHVA86oq\nTdRQ5HqwMfpiLBBMKNQcAJsO71RKNrDWwJJZiyYbvM4FOWRZFtRSdPIDgX0C0Wg1\nNnqWYvHDyA5Ug+T8kowiQDn56dU6Km2FWO4wnqZeA8q5G7rQVXlqdibuiP7FglHA\neURUzFsqyymXMUGrqDPqrHsVWC2E3NTJEb4QlywtrwI13qbhlvTx6/9oRfUjytXJ\nRuUIE5xL8yhRCagNr5ZW250aa+wBwu5DSCk5fDNr0eDuZjw84WHDll+mHxBFGV+X\nKJ5jCOmGumGqjVWZesJpNN1My3M9bsY9layNJJ0eiDeHDEi/yXhhO/mNEXhvhq/R\nfN0Jh2A3AgMBAAECggEAef+CEL5aF6/aVs0yh7fiXkKSp1ECXkud8ztgpEn63KJF\nXM1EdnBt50fA2xSQUeGmeEXi6+cngf0nRb8FToAEgLoGoOEjSJuLrzP3I8U9Fe3m\nBRG2uZI2Ti/bD0eRGEc1oSDhCpsqnkTGK1bwcD4AKpwY+c08Izh/2BOoY6McDoqh\ndQ89jzTuMtD4cNlnPiIrY9HbxoNjshK2ax1OaeXyYKZFG1TxqMFv5gA/G5+S3Cwr\nVG4fkAxYi5vdIK3b8jUXrTM/kpoTl+d3dlQ7rRZYf7KyT31/HtJ/GNzxFI6upzO7\niDNrrUOyeOPjWXdzUh9budv3j+6UfbYK7uZIoebHIQKBgQDykYX1L/akGaOC2tfS\njzCWUzPxGFYVZQ7zD1PM6UyouuS1KLURDEGk9RxqVzTPh/pYd8Ivnz3vOVski5Zt\n19ozLGxdtDhn122DcnVpfCdYzHBdAzPCzORenFohX+MhiX5fEotTlVi7wfOmzTP5\nhUCMSd/17bJrV4XMLhkdrMRBFQKBgQDH5fwV7o+Ej/ZfcdGIa3mAFazToPDzxhHU\nnwADxaxpNGKRU03XCaiYkykarLYdG6Rk+7dXUv8eLy+6Dcq1SWQtfCWKEor++IIp\n1RwWmFHfYriHGkmxSkkEkLFvL8old9xM5YWbEXc4QIXvnfR4BZxdyJHVzIDdbI2i\nFgcn17U3GwKBgDd1njMY7ENIuWHJt16k7m7wRwfwkH4DxQ89ieNn0+cgE/p3fC6R\nptCYWg7WMXThmhNwDi3lMrvnWTdZ0uL6XyEkHwKtmdfkIV3UZZPglv5uf6JEgSkg\nv3YCOXk3+y5HyWTjUIejtc334kVY1XFPThrFKTeJSSnRsP2l7IgkYBqhAoGAYGsr\nM3z1FrDF2nWw5odIfKJ30UAw2LRyB0eGH0uqhLgyzvwKcK2E98sLqYUi9llN6zOK\n1IEA8xM5hxl97AFxY4sdJEMbbi55wim7uZ5Q51nbvbbNUsmM/Lm6C/JWI8pzpVeU\nIR7EjYp50AE1WOsD6CyFQ0W35pWkn0jWvL4L938CgYEAlax4dLE5+YCG5DeLT1/N\nin6kQjl7JS4oCNlPEj6aRPLX2OJresI7oOn+uNatKvlVyfPm6rdxeHCmER1FpR7Q\nA/aNVjPeViKT/R4nK9ytsa+s/1IJVrwLFHJK3toGE660g5w3vKrCjWisMdP4yzzQ\nbv1KwcKoQbNVXwauH79JKc0=\n-----END PRIVATE KEY-----\n",
			Path:            "/path",
			Period:          12,
			GzipLevel:       9,
			FormatVersion:   2,
			Format:          "format",
			MessageType:     "blank",
			TimestampFormat: "%Y",
			Placement:       "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "gcses/cleanup", func(c *Client) {
			c.DeleteGCS(&DeleteGCSInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-gcs",
			})

			c.DeleteGCS(&DeleteGCSInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-gcs",
			})
		})
	}()

	if gcs.Name != "test-gcs" {
		t.Errorf("bad name: %q", gcs.Name)
	}
	if gcs.Bucket != "bucket" {
		t.Errorf("bad bucket: %q", gcs.Bucket)
	}
	if gcs.User != "user" {
		t.Errorf("bad user: %q", gcs.User)
	}
	if gcs.SecretKey != "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9aQoqdHVA86oq\nTdRQ5HqwMfpiLBBMKNQcAJsO71RKNrDWwJJZiyYbvM4FOWRZFtRSdPIDgX0C0Wg1\nNnqWYvHDyA5Ug+T8kowiQDn56dU6Km2FWO4wnqZeA8q5G7rQVXlqdibuiP7FglHA\neURUzFsqyymXMUGrqDPqrHsVWC2E3NTJEb4QlywtrwI13qbhlvTx6/9oRfUjytXJ\nRuUIE5xL8yhRCagNr5ZW250aa+wBwu5DSCk5fDNr0eDuZjw84WHDll+mHxBFGV+X\nKJ5jCOmGumGqjVWZesJpNN1My3M9bsY9layNJJ0eiDeHDEi/yXhhO/mNEXhvhq/R\nfN0Jh2A3AgMBAAECggEAef+CEL5aF6/aVs0yh7fiXkKSp1ECXkud8ztgpEn63KJF\nXM1EdnBt50fA2xSQUeGmeEXi6+cngf0nRb8FToAEgLoGoOEjSJuLrzP3I8U9Fe3m\nBRG2uZI2Ti/bD0eRGEc1oSDhCpsqnkTGK1bwcD4AKpwY+c08Izh/2BOoY6McDoqh\ndQ89jzTuMtD4cNlnPiIrY9HbxoNjshK2ax1OaeXyYKZFG1TxqMFv5gA/G5+S3Cwr\nVG4fkAxYi5vdIK3b8jUXrTM/kpoTl+d3dlQ7rRZYf7KyT31/HtJ/GNzxFI6upzO7\niDNrrUOyeOPjWXdzUh9budv3j+6UfbYK7uZIoebHIQKBgQDykYX1L/akGaOC2tfS\njzCWUzPxGFYVZQ7zD1PM6UyouuS1KLURDEGk9RxqVzTPh/pYd8Ivnz3vOVski5Zt\n19ozLGxdtDhn122DcnVpfCdYzHBdAzPCzORenFohX+MhiX5fEotTlVi7wfOmzTP5\nhUCMSd/17bJrV4XMLhkdrMRBFQKBgQDH5fwV7o+Ej/ZfcdGIa3mAFazToPDzxhHU\nnwADxaxpNGKRU03XCaiYkykarLYdG6Rk+7dXUv8eLy+6Dcq1SWQtfCWKEor++IIp\n1RwWmFHfYriHGkmxSkkEkLFvL8old9xM5YWbEXc4QIXvnfR4BZxdyJHVzIDdbI2i\nFgcn17U3GwKBgDd1njMY7ENIuWHJt16k7m7wRwfwkH4DxQ89ieNn0+cgE/p3fC6R\nptCYWg7WMXThmhNwDi3lMrvnWTdZ0uL6XyEkHwKtmdfkIV3UZZPglv5uf6JEgSkg\nv3YCOXk3+y5HyWTjUIejtc334kVY1XFPThrFKTeJSSnRsP2l7IgkYBqhAoGAYGsr\nM3z1FrDF2nWw5odIfKJ30UAw2LRyB0eGH0uqhLgyzvwKcK2E98sLqYUi9llN6zOK\n1IEA8xM5hxl97AFxY4sdJEMbbi55wim7uZ5Q51nbvbbNUsmM/Lm6C/JWI8pzpVeU\nIR7EjYp50AE1WOsD6CyFQ0W35pWkn0jWvL4L938CgYEAlax4dLE5+YCG5DeLT1/N\nin6kQjl7JS4oCNlPEj6aRPLX2OJresI7oOn+uNatKvlVyfPm6rdxeHCmER1FpR7Q\nA/aNVjPeViKT/R4nK9ytsa+s/1IJVrwLFHJK3toGE660g5w3vKrCjWisMdP4yzzQ\nbv1KwcKoQbNVXwauH79JKc0=\n-----END PRIVATE KEY-----\n" {
		t.Errorf("bad secret_key: %q", gcs.SecretKey)
	}
	if gcs.Path != "/path" {
		t.Errorf("bad path: %q", gcs.Path)
	}
	if gcs.Period != 12 {
		t.Errorf("bad period: %q", gcs.Period)
	}
	if gcs.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", gcs.GzipLevel)
	}
	if gcs.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", gcs.FormatVersion)
	}
	if gcs.Format != "format" {
		t.Errorf("bad format: %q", gcs.Format)
	}
	if gcs.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", gcs.TimestampFormat)
	}
	if gcs.MessageType != "blank" {
		t.Errorf("bad message_type: %q", gcs.MessageType)
	}
	if gcs.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", gcs.Placement)
	}

	// List
	var gcses []*GCS
	record(t, "gcses/list", func(c *Client) {
		gcses, err = c.ListGCSs(&ListGCSsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(gcses) < 1 {
		t.Errorf("bad gcses: %v", gcses)
	}

	// Get
	var ngcs *GCS
	record(t, "gcses/get", func(c *Client) {
		ngcs, err = c.GetGCS(&GetGCSInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-gcs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gcs.Name != ngcs.Name {
		t.Errorf("bad name: %q", gcs.Name)
	}
	if gcs.Bucket != ngcs.Bucket {
		t.Errorf("bad bucket: %q", gcs.Bucket)
	}
	if gcs.User != ngcs.User {
		t.Errorf("bad user: %q", gcs.User)
	}
	if gcs.SecretKey != ngcs.SecretKey {
		t.Errorf("bad secret_key: %q", gcs.SecretKey)
	}
	if gcs.Path != ngcs.Path {
		t.Errorf("bad path: %q", gcs.Path)
	}
	if gcs.Period != ngcs.Period {
		t.Errorf("bad period: %q", gcs.Period)
	}
	if gcs.GzipLevel != ngcs.GzipLevel {
		t.Errorf("bad gzip_level: %q", gcs.GzipLevel)
	}
	if gcs.FormatVersion != ngcs.FormatVersion {
		t.Errorf("bad format_version: %q", gcs.FormatVersion)
	}
	if gcs.Format != ngcs.Format {
		t.Errorf("bad format: %q", gcs.Format)
	}
	if gcs.TimestampFormat != ngcs.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", gcs.TimestampFormat)
	}
	if gcs.MessageType != ngcs.MessageType {
		t.Errorf("bad message_type: %q", gcs.MessageType)
	}
	if gcs.Placement != ngcs.Placement {
		t.Errorf("bad placement: %q", gcs.Placement)
	}

	// Update
	var ugcs *GCS
	record(t, "gcses/update", func(c *Client) {
		ugcs, err = c.UpdateGCS(&UpdateGCSInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-gcs",
			NewName: "new-test-gcs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ugcs.Name != "new-test-gcs" {
		t.Errorf("bad name: %q", ugcs.Name)
	}

	// Delete
	record(t, "gcses/delete", func(c *Client) {
		err = c.DeleteGCS(&DeleteGCSInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-gcs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListGCSs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListGCSs(&ListGCSsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListGCSs(&ListGCSsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateGCS(&CreateGCSInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateGCS(&CreateGCSInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.GetGCS(&GetGCSInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetGCS(&GetGCSInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetGCS(&GetGCSInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteGCS_validation(t *testing.T) {
	var err error
	err = testClient.DeleteGCS(&DeleteGCSInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteGCS(&DeleteGCSInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteGCS(&DeleteGCSInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
