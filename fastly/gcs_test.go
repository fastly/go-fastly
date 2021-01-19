package fastly

import (
	"testing"
)

func TestClient_GCSs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "gcses/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var gcsCreateResp1, gcsCreateResp2, gcsCreateResp3 *GCS
	record(t, "gcses/create", func(c *Client) {
		gcsCreateResp1, err = c.CreateGCS(&CreateGCSInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-gcs",
			Bucket:           "bucket",
			User:             "user",
			SecretKey:        "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9aQoqdHVA86oq\nTdRQ5HqwMfpiLBBMKNQcAJsO71RKNrDWwJJZiyYbvM4FOWRZFtRSdPIDgX0C0Wg1\nNnqWYvHDyA5Ug+T8kowiQDn56dU6Km2FWO4wnqZeA8q5G7rQVXlqdibuiP7FglHA\neURUzFsqyymXMUGrqDPqrHsVWC2E3NTJEb4QlywtrwI13qbhlvTx6/9oRfUjytXJ\nRuUIE5xL8yhRCagNr5ZW250aa+wBwu5DSCk5fDNr0eDuZjw84WHDll+mHxBFGV+X\nKJ5jCOmGumGqjVWZesJpNN1My3M9bsY9layNJJ0eiDeHDEi/yXhhO/mNEXhvhq/R\nfN0Jh2A3AgMBAAECggEAef+CEL5aF6/aVs0yh7fiXkKSp1ECXkud8ztgpEn63KJF\nXM1EdnBt50fA2xSQUeGmeEXi6+cngf0nRb8FToAEgLoGoOEjSJuLrzP3I8U9Fe3m\nBRG2uZI2Ti/bD0eRGEc1oSDhCpsqnkTGK1bwcD4AKpwY+c08Izh/2BOoY6McDoqh\ndQ89jzTuMtD4cNlnPiIrY9HbxoNjshK2ax1OaeXyYKZFG1TxqMFv5gA/G5+S3Cwr\nVG4fkAxYi5vdIK3b8jUXrTM/kpoTl+d3dlQ7rRZYf7KyT31/HtJ/GNzxFI6upzO7\niDNrrUOyeOPjWXdzUh9budv3j+6UfbYK7uZIoebHIQKBgQDykYX1L/akGaOC2tfS\njzCWUzPxGFYVZQ7zD1PM6UyouuS1KLURDEGk9RxqVzTPh/pYd8Ivnz3vOVski5Zt\n19ozLGxdtDhn122DcnVpfCdYzHBdAzPCzORenFohX+MhiX5fEotTlVi7wfOmzTP5\nhUCMSd/17bJrV4XMLhkdrMRBFQKBgQDH5fwV7o+Ej/ZfcdGIa3mAFazToPDzxhHU\nnwADxaxpNGKRU03XCaiYkykarLYdG6Rk+7dXUv8eLy+6Dcq1SWQtfCWKEor++IIp\n1RwWmFHfYriHGkmxSkkEkLFvL8old9xM5YWbEXc4QIXvnfR4BZxdyJHVzIDdbI2i\nFgcn17U3GwKBgDd1njMY7ENIuWHJt16k7m7wRwfwkH4DxQ89ieNn0+cgE/p3fC6R\nptCYWg7WMXThmhNwDi3lMrvnWTdZ0uL6XyEkHwKtmdfkIV3UZZPglv5uf6JEgSkg\nv3YCOXk3+y5HyWTjUIejtc334kVY1XFPThrFKTeJSSnRsP2l7IgkYBqhAoGAYGsr\nM3z1FrDF2nWw5odIfKJ30UAw2LRyB0eGH0uqhLgyzvwKcK2E98sLqYUi9llN6zOK\n1IEA8xM5hxl97AFxY4sdJEMbbi55wim7uZ5Q51nbvbbNUsmM/Lm6C/JWI8pzpVeU\nIR7EjYp50AE1WOsD6CyFQ0W35pWkn0jWvL4L938CgYEAlax4dLE5+YCG5DeLT1/N\nin6kQjl7JS4oCNlPEj6aRPLX2OJresI7oOn+uNatKvlVyfPm6rdxeHCmER1FpR7Q\nA/aNVjPeViKT/R4nK9ytsa+s/1IJVrwLFHJK3toGE660g5w3vKrCjWisMdP4yzzQ\nbv1KwcKoQbNVXwauH79JKc0=\n-----END PRIVATE KEY-----\n",
			Path:             "/path",
			Period:           12,
			CompressionCodec: "snappy",
			FormatVersion:    2,
			Format:           "format",
			MessageType:      "blank",
			TimestampFormat:  "%Y",
			Placement:        "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "gcses/create2", func(c *Client) {
		gcsCreateResp2, err = c.CreateGCS(&CreateGCSInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-gcs-2",
			Bucket:           "bucket",
			User:             "user",
			SecretKey:        "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9aQoqdHVA86oq\nTdRQ5HqwMfpiLBBMKNQcAJsO71RKNrDWwJJZiyYbvM4FOWRZFtRSdPIDgX0C0Wg1\nNnqWYvHDyA5Ug+T8kowiQDn56dU6Km2FWO4wnqZeA8q5G7rQVXlqdibuiP7FglHA\neURUzFsqyymXMUGrqDPqrHsVWC2E3NTJEb4QlywtrwI13qbhlvTx6/9oRfUjytXJ\nRuUIE5xL8yhRCagNr5ZW250aa+wBwu5DSCk5fDNr0eDuZjw84WHDll+mHxBFGV+X\nKJ5jCOmGumGqjVWZesJpNN1My3M9bsY9layNJJ0eiDeHDEi/yXhhO/mNEXhvhq/R\nfN0Jh2A3AgMBAAECggEAef+CEL5aF6/aVs0yh7fiXkKSp1ECXkud8ztgpEn63KJF\nXM1EdnBt50fA2xSQUeGmeEXi6+cngf0nRb8FToAEgLoGoOEjSJuLrzP3I8U9Fe3m\nBRG2uZI2Ti/bD0eRGEc1oSDhCpsqnkTGK1bwcD4AKpwY+c08Izh/2BOoY6McDoqh\ndQ89jzTuMtD4cNlnPiIrY9HbxoNjshK2ax1OaeXyYKZFG1TxqMFv5gA/G5+S3Cwr\nVG4fkAxYi5vdIK3b8jUXrTM/kpoTl+d3dlQ7rRZYf7KyT31/HtJ/GNzxFI6upzO7\niDNrrUOyeOPjWXdzUh9budv3j+6UfbYK7uZIoebHIQKBgQDykYX1L/akGaOC2tfS\njzCWUzPxGFYVZQ7zD1PM6UyouuS1KLURDEGk9RxqVzTPh/pYd8Ivnz3vOVski5Zt\n19ozLGxdtDhn122DcnVpfCdYzHBdAzPCzORenFohX+MhiX5fEotTlVi7wfOmzTP5\nhUCMSd/17bJrV4XMLhkdrMRBFQKBgQDH5fwV7o+Ej/ZfcdGIa3mAFazToPDzxhHU\nnwADxaxpNGKRU03XCaiYkykarLYdG6Rk+7dXUv8eLy+6Dcq1SWQtfCWKEor++IIp\n1RwWmFHfYriHGkmxSkkEkLFvL8old9xM5YWbEXc4QIXvnfR4BZxdyJHVzIDdbI2i\nFgcn17U3GwKBgDd1njMY7ENIuWHJt16k7m7wRwfwkH4DxQ89ieNn0+cgE/p3fC6R\nptCYWg7WMXThmhNwDi3lMrvnWTdZ0uL6XyEkHwKtmdfkIV3UZZPglv5uf6JEgSkg\nv3YCOXk3+y5HyWTjUIejtc334kVY1XFPThrFKTeJSSnRsP2l7IgkYBqhAoGAYGsr\nM3z1FrDF2nWw5odIfKJ30UAw2LRyB0eGH0uqhLgyzvwKcK2E98sLqYUi9llN6zOK\n1IEA8xM5hxl97AFxY4sdJEMbbi55wim7uZ5Q51nbvbbNUsmM/Lm6C/JWI8pzpVeU\nIR7EjYp50AE1WOsD6CyFQ0W35pWkn0jWvL4L938CgYEAlax4dLE5+YCG5DeLT1/N\nin6kQjl7JS4oCNlPEj6aRPLX2OJresI7oOn+uNatKvlVyfPm6rdxeHCmER1FpR7Q\nA/aNVjPeViKT/R4nK9ytsa+s/1IJVrwLFHJK3toGE660g5w3vKrCjWisMdP4yzzQ\nbv1KwcKoQbNVXwauH79JKc0=\n-----END PRIVATE KEY-----\n",
			Path:             "/path",
			Period:           12,
			CompressionCodec: "snappy",
			GzipLevel:        8,
			FormatVersion:    2,
			Format:           "format",
			MessageType:      "blank",
			TimestampFormat:  "%Y",
			Placement:        "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "gcses/create3", func(c *Client) {
		gcsCreateResp3, err = c.CreateGCS(&CreateGCSInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-gcs-3",
			Bucket:           "bucket",
			User:             "user",
			SecretKey:        "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9aQoqdHVA86oq\nTdRQ5HqwMfpiLBBMKNQcAJsO71RKNrDWwJJZiyYbvM4FOWRZFtRSdPIDgX0C0Wg1\nNnqWYvHDyA5Ug+T8kowiQDn56dU6Km2FWO4wnqZeA8q5G7rQVXlqdibuiP7FglHA\neURUzFsqyymXMUGrqDPqrHsVWC2E3NTJEb4QlywtrwI13qbhlvTx6/9oRfUjytXJ\nRuUIE5xL8yhRCagNr5ZW250aa+wBwu5DSCk5fDNr0eDuZjw84WHDll+mHxBFGV+X\nKJ5jCOmGumGqjVWZesJpNN1My3M9bsY9layNJJ0eiDeHDEi/yXhhO/mNEXhvhq/R\nfN0Jh2A3AgMBAAECggEAef+CEL5aF6/aVs0yh7fiXkKSp1ECXkud8ztgpEn63KJF\nXM1EdnBt50fA2xSQUeGmeEXi6+cngf0nRb8FToAEgLoGoOEjSJuLrzP3I8U9Fe3m\nBRG2uZI2Ti/bD0eRGEc1oSDhCpsqnkTGK1bwcD4AKpwY+c08Izh/2BOoY6McDoqh\ndQ89jzTuMtD4cNlnPiIrY9HbxoNjshK2ax1OaeXyYKZFG1TxqMFv5gA/G5+S3Cwr\nVG4fkAxYi5vdIK3b8jUXrTM/kpoTl+d3dlQ7rRZYf7KyT31/HtJ/GNzxFI6upzO7\niDNrrUOyeOPjWXdzUh9budv3j+6UfbYK7uZIoebHIQKBgQDykYX1L/akGaOC2tfS\njzCWUzPxGFYVZQ7zD1PM6UyouuS1KLURDEGk9RxqVzTPh/pYd8Ivnz3vOVski5Zt\n19ozLGxdtDhn122DcnVpfCdYzHBdAzPCzORenFohX+MhiX5fEotTlVi7wfOmzTP5\nhUCMSd/17bJrV4XMLhkdrMRBFQKBgQDH5fwV7o+Ej/ZfcdGIa3mAFazToPDzxhHU\nnwADxaxpNGKRU03XCaiYkykarLYdG6Rk+7dXUv8eLy+6Dcq1SWQtfCWKEor++IIp\n1RwWmFHfYriHGkmxSkkEkLFvL8old9xM5YWbEXc4QIXvnfR4BZxdyJHVzIDdbI2i\nFgcn17U3GwKBgDd1njMY7ENIuWHJt16k7m7wRwfwkH4DxQ89ieNn0+cgE/p3fC6R\nptCYWg7WMXThmhNwDi3lMrvnWTdZ0uL6XyEkHwKtmdfkIV3UZZPglv5uf6JEgSkg\nv3YCOXk3+y5HyWTjUIejtc334kVY1XFPThrFKTeJSSnRsP2l7IgkYBqhAoGAYGsr\nM3z1FrDF2nWw5odIfKJ30UAw2LRyB0eGH0uqhLgyzvwKcK2E98sLqYUi9llN6zOK\n1IEA8xM5hxl97AFxY4sdJEMbbi55wim7uZ5Q51nbvbbNUsmM/Lm6C/JWI8pzpVeU\nIR7EjYp50AE1WOsD6CyFQ0W35pWkn0jWvL4L938CgYEAlax4dLE5+YCG5DeLT1/N\nin6kQjl7JS4oCNlPEj6aRPLX2OJresI7oOn+uNatKvlVyfPm6rdxeHCmER1FpR7Q\nA/aNVjPeViKT/R4nK9ytsa+s/1IJVrwLFHJK3toGE660g5w3vKrCjWisMdP4yzzQ\nbv1KwcKoQbNVXwauH79JKc0=\n-----END PRIVATE KEY-----\n",
			Path:             "/path",
			Period:           12,
			CompressionCodec: "snappy",
			FormatVersion:    2,
			Format:           "format",
			MessageType:      "blank",
			TimestampFormat:  "%Y",
			Placement:        "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "gcses/cleanup", func(c *Client) {
			c.DeleteGCS(&DeleteGCSInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-gcs",
			})

			c.DeleteGCS(&DeleteGCSInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-gcs-2",
			})

			c.DeleteGCS(&DeleteGCSInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-gcs-3",
			})

			c.DeleteGCS(&DeleteGCSInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-gcs",
			})
		})
	}()

	if gcsCreateResp1.Name != "test-gcs" {
		t.Errorf("bad name: %q", gcsCreateResp1.Name)
	}
	if gcsCreateResp1.Bucket != "bucket" {
		t.Errorf("bad bucket: %q", gcsCreateResp1.Bucket)
	}
	if gcsCreateResp1.User != "user" {
		t.Errorf("bad user: %q", gcsCreateResp1.User)
	}
	if gcsCreateResp1.SecretKey != "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9aQoqdHVA86oq\nTdRQ5HqwMfpiLBBMKNQcAJsO71RKNrDWwJJZiyYbvM4FOWRZFtRSdPIDgX0C0Wg1\nNnqWYvHDyA5Ug+T8kowiQDn56dU6Km2FWO4wnqZeA8q5G7rQVXlqdibuiP7FglHA\neURUzFsqyymXMUGrqDPqrHsVWC2E3NTJEb4QlywtrwI13qbhlvTx6/9oRfUjytXJ\nRuUIE5xL8yhRCagNr5ZW250aa+wBwu5DSCk5fDNr0eDuZjw84WHDll+mHxBFGV+X\nKJ5jCOmGumGqjVWZesJpNN1My3M9bsY9layNJJ0eiDeHDEi/yXhhO/mNEXhvhq/R\nfN0Jh2A3AgMBAAECggEAef+CEL5aF6/aVs0yh7fiXkKSp1ECXkud8ztgpEn63KJF\nXM1EdnBt50fA2xSQUeGmeEXi6+cngf0nRb8FToAEgLoGoOEjSJuLrzP3I8U9Fe3m\nBRG2uZI2Ti/bD0eRGEc1oSDhCpsqnkTGK1bwcD4AKpwY+c08Izh/2BOoY6McDoqh\ndQ89jzTuMtD4cNlnPiIrY9HbxoNjshK2ax1OaeXyYKZFG1TxqMFv5gA/G5+S3Cwr\nVG4fkAxYi5vdIK3b8jUXrTM/kpoTl+d3dlQ7rRZYf7KyT31/HtJ/GNzxFI6upzO7\niDNrrUOyeOPjWXdzUh9budv3j+6UfbYK7uZIoebHIQKBgQDykYX1L/akGaOC2tfS\njzCWUzPxGFYVZQ7zD1PM6UyouuS1KLURDEGk9RxqVzTPh/pYd8Ivnz3vOVski5Zt\n19ozLGxdtDhn122DcnVpfCdYzHBdAzPCzORenFohX+MhiX5fEotTlVi7wfOmzTP5\nhUCMSd/17bJrV4XMLhkdrMRBFQKBgQDH5fwV7o+Ej/ZfcdGIa3mAFazToPDzxhHU\nnwADxaxpNGKRU03XCaiYkykarLYdG6Rk+7dXUv8eLy+6Dcq1SWQtfCWKEor++IIp\n1RwWmFHfYriHGkmxSkkEkLFvL8old9xM5YWbEXc4QIXvnfR4BZxdyJHVzIDdbI2i\nFgcn17U3GwKBgDd1njMY7ENIuWHJt16k7m7wRwfwkH4DxQ89ieNn0+cgE/p3fC6R\nptCYWg7WMXThmhNwDi3lMrvnWTdZ0uL6XyEkHwKtmdfkIV3UZZPglv5uf6JEgSkg\nv3YCOXk3+y5HyWTjUIejtc334kVY1XFPThrFKTeJSSnRsP2l7IgkYBqhAoGAYGsr\nM3z1FrDF2nWw5odIfKJ30UAw2LRyB0eGH0uqhLgyzvwKcK2E98sLqYUi9llN6zOK\n1IEA8xM5hxl97AFxY4sdJEMbbi55wim7uZ5Q51nbvbbNUsmM/Lm6C/JWI8pzpVeU\nIR7EjYp50AE1WOsD6CyFQ0W35pWkn0jWvL4L938CgYEAlax4dLE5+YCG5DeLT1/N\nin6kQjl7JS4oCNlPEj6aRPLX2OJresI7oOn+uNatKvlVyfPm6rdxeHCmER1FpR7Q\nA/aNVjPeViKT/R4nK9ytsa+s/1IJVrwLFHJK3toGE660g5w3vKrCjWisMdP4yzzQ\nbv1KwcKoQbNVXwauH79JKc0=\n-----END PRIVATE KEY-----\n" {
		t.Errorf("bad secret_key: %q", gcsCreateResp1.SecretKey)
	}
	if gcsCreateResp1.Path != "/path" {
		t.Errorf("bad path: %q", gcsCreateResp1.Path)
	}
	if gcsCreateResp1.Period != 12 {
		t.Errorf("bad period: %q", gcsCreateResp1.Period)
	}
	if gcsCreateResp1.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", gcsCreateResp1.CompressionCodec)
	}
	if gcsCreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", gcsCreateResp1.GzipLevel)
	}
	if gcsCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", gcsCreateResp1.FormatVersion)
	}
	if gcsCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", gcsCreateResp1.Format)
	}
	if gcsCreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", gcsCreateResp1.TimestampFormat)
	}
	if gcsCreateResp1.MessageType != "blank" {
		t.Errorf("bad message_type: %q", gcsCreateResp1.MessageType)
	}
	if gcsCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", gcsCreateResp1.Placement)
	}

	if gcsCreateResp2.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", gcsCreateResp1.CompressionCodec)
	}
	if gcsCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", gcsCreateResp1.GzipLevel)
	}

	if gcsCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", gcsCreateResp1.CompressionCodec)
	}
	if gcsCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", gcsCreateResp1.GzipLevel)
	}

	// List
	var gcses []*GCS
	record(t, "gcses/list", func(c *Client) {
		gcses, err = c.ListGCSs(&ListGCSsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(gcses) < 1 {
		t.Errorf("bad gcses: %v", gcses)
	}

	// Get
	var gcsGetResp *GCS
	record(t, "gcses/get", func(c *Client) {
		gcsGetResp, err = c.GetGCS(&GetGCSInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-gcs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gcsCreateResp1.Name != gcsGetResp.Name {
		t.Errorf("bad name: %q", gcsCreateResp1.Name)
	}
	if gcsCreateResp1.Bucket != gcsGetResp.Bucket {
		t.Errorf("bad bucket: %q", gcsCreateResp1.Bucket)
	}
	if gcsCreateResp1.User != gcsGetResp.User {
		t.Errorf("bad user: %q", gcsCreateResp1.User)
	}
	if gcsCreateResp1.SecretKey != gcsGetResp.SecretKey {
		t.Errorf("bad secret_key: %q", gcsCreateResp1.SecretKey)
	}
	if gcsCreateResp1.Path != gcsGetResp.Path {
		t.Errorf("bad path: %q", gcsCreateResp1.Path)
	}
	if gcsCreateResp1.Period != gcsGetResp.Period {
		t.Errorf("bad period: %q", gcsCreateResp1.Period)
	}
	if gcsCreateResp1.CompressionCodec != gcsGetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", gcsCreateResp1.CompressionCodec)
	}
	if gcsCreateResp1.GzipLevel != gcsGetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", gcsCreateResp1.GzipLevel)
	}
	if gcsCreateResp1.FormatVersion != gcsGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", gcsCreateResp1.FormatVersion)
	}
	if gcsCreateResp1.Format != gcsGetResp.Format {
		t.Errorf("bad format: %q", gcsCreateResp1.Format)
	}
	if gcsCreateResp1.TimestampFormat != gcsGetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", gcsCreateResp1.TimestampFormat)
	}
	if gcsCreateResp1.MessageType != gcsGetResp.MessageType {
		t.Errorf("bad message_type: %q", gcsCreateResp1.MessageType)
	}
	if gcsCreateResp1.Placement != gcsGetResp.Placement {
		t.Errorf("bad placement: %q", gcsCreateResp1.Placement)
	}

	// Update
	var gcsUpdateResp1, gcsUpdateResp2, gcsUpdateResp3 *GCS
	record(t, "gcses/update", func(c *Client) {
		gcsUpdateResp1, err = c.UpdateGCS(&UpdateGCSInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-gcs",
			NewName:        String("new-test-gcs"),
			MessageType:    String("classic"),
			GzipLevel:      Uint8(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "gcses/update2", func(c *Client) {
		gcsUpdateResp2, err = c.UpdateGCS(&UpdateGCSInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-gcs-2",
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "gcses/update3", func(c *Client) {
		gcsUpdateResp3, err = c.UpdateGCS(&UpdateGCSInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-gcs-3",
			GzipLevel:      Uint8(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if gcsUpdateResp1.Name != "new-test-gcs" {
		t.Errorf("bad name: %q", gcsUpdateResp1.Name)
	}
	if gcsUpdateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", gcsUpdateResp1.MessageType)
	}
	if gcsUpdateResp1.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", gcsUpdateResp1.CompressionCodec)
	}
	if gcsUpdateResp1.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", gcsUpdateResp1.GzipLevel)
	}

	if gcsUpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", gcsUpdateResp2.CompressionCodec)
	}
	if gcsUpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", gcsUpdateResp2.GzipLevel)
	}

	if gcsUpdateResp3.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", gcsUpdateResp3.CompressionCodec)
	}
	if gcsUpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", gcsUpdateResp3.GzipLevel)
	}

	// Delete
	record(t, "gcses/delete", func(c *Client) {
		err = c.DeleteGCS(&DeleteGCSInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-gcs",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListGCSs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListGCSs(&ListGCSsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListGCSs(&ListGCSsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateGCS(&CreateGCSInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateGCS(&CreateGCSInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.GetGCS(&GetGCSInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetGCS(&GetGCSInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetGCS(&GetGCSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateGCS_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateGCS(&UpdateGCSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteGCS_validation(t *testing.T) {
	var err error
	err = testClient.DeleteGCS(&DeleteGCSInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteGCS(&DeleteGCSInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteGCS(&DeleteGCSInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
