// Code generated by 'service_linked_product' generator, DO NOT EDIT.

package brotli_compression_test

import (
	"testing"

	fastly "github.com/fastly/go-fastly/v9/fastly"
	brotlicompression "github.com/fastly/go-fastly/v9/fastly/products/brotli_compression"
)

func Test_Get_validation(t *testing.T) {
	if _, err := brotlicompression.Get(fastly.TestClient, ""); err != fastly.ErrMissingServiceID {
		t.Fatalf("expected '%s', got: '%s'", fastly.ErrMissingServiceID, err)
	}
}

func Test_Enable_validation(t *testing.T) {
	if _, err := brotlicompression.Enable(fastly.TestClient, ""); err != fastly.ErrMissingServiceID {
		t.Fatalf("expected '%s', got: '%s'", fastly.ErrMissingServiceID, err)
	}
}

func Test_Disable_validation(t *testing.T) {
	if err := brotlicompression.Disable(fastly.TestClient, ""); err != fastly.ErrMissingServiceID {
		t.Fatalf("expected '%s', got: '%s'", fastly.ErrMissingServiceID, err)
	}
}
