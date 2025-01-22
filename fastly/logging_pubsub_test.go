package fastly

import (
	"errors"
	"strings"
	"testing"
)

func TestClient_Pubsubs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "pubsubs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var pubsub *Pubsub
	Record(t, "pubsubs/create", func(c *Client) {
		pubsub, err = c.CreatePubsub(&CreatePubsubInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-pubsub"),
			Topic:          ToPointer("topic"),
			User:           ToPointer("user"),
			AccountName:    ToPointer("service-account"),
			SecretKey:      ToPointer("-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9aQoqdHVA86oq\nTdRQ5HqwMfpiLBBMKNQcAJsO71RKNrDWwJJZiyYbvM4FOWRZFtRSdPIDgX0C0Wg1\nNnqWYvHDyA5Ug+T8kowiQDn56dU6Km2FWO4wnqZeA8q5G7rQVXlqdibuiP7FglHA\neURUzFsqyymXMUGrqDPqrHsVWC2E3NTJEb4QlywtrwI13qbhlvTx6/9oRfUjytXJ\nRuUIE5xL8yhRCagNr5ZW250aa+wBwu5DSCk5fDNr0eDuZjw84WHDll+mHxBFGV+X\nKJ5jCOmGumGqjVWZesJpNN1My3M9bsY9layNJJ0eiDeHDEi/yXhhO/mNEXhvhq/R\nfN0Jh2A3AgMBAAECggEAef+CEL5aF6/aVs0yh7fiXkKSp1ECXkud8ztgpEn63KJF\nXM1EdnBt50fA2xSQUeGmeEXi6+cngf0nRb8FToAEgLoGoOEjSJuLrzP3I8U9Fe3m\nBRG2uZI2Ti/bD0eRGEc1oSDhCpsqnkTGK1bwcD4AKpwY+c08Izh/2BOoY6McDoqh\ndQ89jzTuMtD4cNlnPiIrY9HbxoNjshK2ax1OaeXyYKZFG1TxqMFv5gA/G5+S3Cwr\nVG4fkAxYi5vdIK3b8jUXrTM/kpoTl+d3dlQ7rRZYf7KyT31/HtJ/GNzxFI6upzO7\niDNrrUOyeOPjWXdzUh9budv3j+6UfbYK7uZIoebHIQKBgQDykYX1L/akGaOC2tfS\njzCWUzPxGFYVZQ7zD1PM6UyouuS1KLURDEGk9RxqVzTPh/pYd8Ivnz3vOVski5Zt\n19ozLGxdtDhn122DcnVpfCdYzHBdAzPCzORenFohX+MhiX5fEotTlVi7wfOmzTP5\nhUCMSd/17bJrV4XMLhkdrMRBFQKBgQDH5fwV7o+Ej/ZfcdGIa3mAFazToPDzxhHU\nnwADxaxpNGKRU03XCaiYkykarLYdG6Rk+7dXUv8eLy+6Dcq1SWQtfCWKEor++IIp\n1RwWmFHfYriHGkmxSkkEkLFvL8old9xM5YWbEXc4QIXvnfR4BZxdyJHVzIDdbI2i\nFgcn17U3GwKBgDd1njMY7ENIuWHJt16k7m7wRwfwkH4DxQ89ieNn0+cgE/p3fC6R\nptCYWg7WMXThmhNwDi3lMrvnWTdZ0uL6XyEkHwKtmdfkIV3UZZPglv5uf6JEgSkg\nv3YCOXk3+y5HyWTjUIejtc334kVY1XFPThrFKTeJSSnRsP2l7IgkYBqhAoGAYGsr\nM3z1FrDF2nWw5odIfKJ30UAw2LRyB0eGH0uqhLgyzvwKcK2E98sLqYUi9llN6zOK\n1IEA8xM5hxl97AFxY4sdJEMbbi55wim7uZ5Q51nbvbbNUsmM/Lm6C/JWI8pzpVeU\nIR7EjYp50AE1WOsD6CyFQ0W35pWkn0jWvL4L938CgYEAlax4dLE5+YCG5DeLT1/N\nin6kQjl7JS4oCNlPEj6aRPLX2OJresI7oOn+uNatKvlVyfPm6rdxeHCmER1FpR7Q\nA/aNVjPeViKT/R4nK9ytsa+s/1IJVrwLFHJK3toGE660g5w3vKrCjWisMdP4yzzQ\nbv1KwcKoQbNVXwauH79JKc0=\n-----END PRIVATE KEY-----\n"),
			ProjectID:      ToPointer("project-id"),
			FormatVersion:  ToPointer(2),
			Format:         ToPointer("format"),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "pubsubs/cleanup", func(c *Client) {
			_ = c.DeletePubsub(&DeletePubsubInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-pubsub",
			})

			_ = c.DeletePubsub(&DeletePubsubInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-pubsub",
			})
		})
	}()

	if *pubsub.Name != "test-pubsub" {
		t.Errorf("bad name: %q", *pubsub.Name)
	}
	if *pubsub.Topic != "topic" {
		t.Errorf("bad topic: %q", *pubsub.Topic)
	}
	if *pubsub.User != "user" {
		t.Errorf("bad user: %q", *pubsub.User)
	}
	if *pubsub.AccountName != "service-account" {
		t.Errorf("bad account name: %q", *pubsub.AccountName)
	}
	if strings.TrimSpace(*pubsub.SecretKey) != `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC9aQoqdHVA86oq
TdRQ5HqwMfpiLBBMKNQcAJsO71RKNrDWwJJZiyYbvM4FOWRZFtRSdPIDgX0C0Wg1
NnqWYvHDyA5Ug+T8kowiQDn56dU6Km2FWO4wnqZeA8q5G7rQVXlqdibuiP7FglHA
eURUzFsqyymXMUGrqDPqrHsVWC2E3NTJEb4QlywtrwI13qbhlvTx6/9oRfUjytXJ
RuUIE5xL8yhRCagNr5ZW250aa+wBwu5DSCk5fDNr0eDuZjw84WHDll+mHxBFGV+X
KJ5jCOmGumGqjVWZesJpNN1My3M9bsY9layNJJ0eiDeHDEi/yXhhO/mNEXhvhq/R
fN0Jh2A3AgMBAAECggEAef+CEL5aF6/aVs0yh7fiXkKSp1ECXkud8ztgpEn63KJF
XM1EdnBt50fA2xSQUeGmeEXi6+cngf0nRb8FToAEgLoGoOEjSJuLrzP3I8U9Fe3m
BRG2uZI2Ti/bD0eRGEc1oSDhCpsqnkTGK1bwcD4AKpwY+c08Izh/2BOoY6McDoqh
dQ89jzTuMtD4cNlnPiIrY9HbxoNjshK2ax1OaeXyYKZFG1TxqMFv5gA/G5+S3Cwr
VG4fkAxYi5vdIK3b8jUXrTM/kpoTl+d3dlQ7rRZYf7KyT31/HtJ/GNzxFI6upzO7
iDNrrUOyeOPjWXdzUh9budv3j+6UfbYK7uZIoebHIQKBgQDykYX1L/akGaOC2tfS
jzCWUzPxGFYVZQ7zD1PM6UyouuS1KLURDEGk9RxqVzTPh/pYd8Ivnz3vOVski5Zt
19ozLGxdtDhn122DcnVpfCdYzHBdAzPCzORenFohX+MhiX5fEotTlVi7wfOmzTP5
hUCMSd/17bJrV4XMLhkdrMRBFQKBgQDH5fwV7o+Ej/ZfcdGIa3mAFazToPDzxhHU
nwADxaxpNGKRU03XCaiYkykarLYdG6Rk+7dXUv8eLy+6Dcq1SWQtfCWKEor++IIp
1RwWmFHfYriHGkmxSkkEkLFvL8old9xM5YWbEXc4QIXvnfR4BZxdyJHVzIDdbI2i
Fgcn17U3GwKBgDd1njMY7ENIuWHJt16k7m7wRwfwkH4DxQ89ieNn0+cgE/p3fC6R
ptCYWg7WMXThmhNwDi3lMrvnWTdZ0uL6XyEkHwKtmdfkIV3UZZPglv5uf6JEgSkg
v3YCOXk3+y5HyWTjUIejtc334kVY1XFPThrFKTeJSSnRsP2l7IgkYBqhAoGAYGsr
M3z1FrDF2nWw5odIfKJ30UAw2LRyB0eGH0uqhLgyzvwKcK2E98sLqYUi9llN6zOK
1IEA8xM5hxl97AFxY4sdJEMbbi55wim7uZ5Q51nbvbbNUsmM/Lm6C/JWI8pzpVeU
IR7EjYp50AE1WOsD6CyFQ0W35pWkn0jWvL4L938CgYEAlax4dLE5+YCG5DeLT1/N
in6kQjl7JS4oCNlPEj6aRPLX2OJresI7oOn+uNatKvlVyfPm6rdxeHCmER1FpR7Q
A/aNVjPeViKT/R4nK9ytsa+s/1IJVrwLFHJK3toGE660g5w3vKrCjWisMdP4yzzQ
bv1KwcKoQbNVXwauH79JKc0=
-----END PRIVATE KEY-----` {
		t.Errorf("bad secret_key: %q", *pubsub.SecretKey)
	}
	if *pubsub.ProjectID != "project-id" {
		t.Errorf("bad project_id: %q", *pubsub.ProjectID)
	}
	if *pubsub.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *pubsub.FormatVersion)
	}
	if *pubsub.Format != "format" {
		t.Errorf("bad format: %q", *pubsub.Format)
	}
	if *pubsub.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *pubsub.Placement)
	}

	// List
	var pubsubs []*Pubsub
	Record(t, "pubsubs/list", func(c *Client) {
		pubsubs, err = c.ListPubsubs(&ListPubsubsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(pubsubs) < 1 {
		t.Errorf("bad pubsubs: %v", pubsubs)
	}

	// Get
	var npubsub *Pubsub
	Record(t, "pubsubs/get", func(c *Client) {
		npubsub, err = c.GetPubsub(&GetPubsubInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-pubsub",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *pubsub.Name != *npubsub.Name {
		t.Errorf("bad name: %q", *pubsub.Name)
	}
	if *pubsub.Topic != *npubsub.Topic {
		t.Errorf("bad topic: %q", *pubsub.Topic)
	}
	if *pubsub.User != *npubsub.User {
		t.Errorf("bad user: %q", *pubsub.User)
	}
	if *pubsub.SecretKey != *npubsub.SecretKey {
		t.Errorf("bad secret_key: %q", *pubsub.SecretKey)
	}
	if *pubsub.ProjectID != *npubsub.ProjectID {
		t.Errorf("bad project_id: %q", *pubsub.ProjectID)
	}
	if *pubsub.FormatVersion != *npubsub.FormatVersion {
		t.Errorf("bad format_version: %q", *pubsub.FormatVersion)
	}
	if *pubsub.Format != *npubsub.Format {
		t.Errorf("bad format: %q", *pubsub.Format)
	}
	if *pubsub.Placement != *npubsub.Placement {
		t.Errorf("bad placement: %q", *pubsub.Placement)
	}

	// Update
	var upubsub *Pubsub
	Record(t, "pubsubs/update", func(c *Client) {
		upubsub, err = c.UpdatePubsub(&UpdatePubsubInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-pubsub",
			NewName:        ToPointer("new-test-pubsub"),
			Topic:          ToPointer("new-topic"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *upubsub.Name != "new-test-pubsub" {
		t.Errorf("bad name: %q", *upubsub.Name)
	}
	if *upubsub.Topic != "new-topic" {
		t.Errorf("bad topic: %q", *upubsub.Topic)
	}

	// Delete
	Record(t, "pubsubs/delete", func(c *Client) {
		err = c.DeletePubsub(&DeletePubsubInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-pubsub",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPubsubs_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListPubsubs(&ListPubsubsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListPubsubs(&ListPubsubsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreatePubsub_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreatePubsub(&CreatePubsubInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreatePubsub(&CreatePubsubInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetPubsub_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetPubsub(&GetPubsubInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetPubsub(&GetPubsubInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetPubsub(&GetPubsubInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePubsub_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdatePubsub(&UpdatePubsubInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdatePubsub(&UpdatePubsubInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdatePubsub(&UpdatePubsubInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeletePubsub_validation(t *testing.T) {
	var err error

	err = TestClient.DeletePubsub(&DeletePubsubInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeletePubsub(&DeletePubsubInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeletePubsub(&DeletePubsubInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
