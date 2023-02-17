package fastly

import (
	"testing"
)

func TestClient_ProductEnablement(t *testing.T) {
	t.Parallel()

	var err error

	// Enable Product
	var pe *ProductEnablement
	record(t, "product_enablement/enable", func(c *Client) {
		pe, err = c.EnableProduct(&ProductEnablementInput{
			ProductID: ProductBrotliCompression,
			ServiceID: testServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if pe.Product.ID != ProductBrotliCompression.String() {
		t.Errorf("bad feature_revision: %s", pe.Product.ID)
	}

	// Get Product status
	var gpe *ProductEnablement
	record(t, "product_enablement/get", func(c *Client) {
		gpe, err = c.GetProduct(&ProductEnablementInput{
			ProductID: ProductBrotliCompression,
			ServiceID: testServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if gpe.Product.ID != ProductBrotliCompression.String() {
		t.Errorf("bad feature_revision: %s", gpe.Product.ID)
	}

	// Disable Product
	record(t, "product_enablement/disable", func(c *Client) {
		err = c.DisableProduct(&ProductEnablementInput{
			ProductID: ProductBrotliCompression,
			ServiceID: testServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get Product status again to check disabled
	record(t, "product_enablement/get-disabled", func(c *Client) {
		gpe, err = c.GetProduct(&ProductEnablementInput{
			ProductID: ProductBrotliCompression,
			ServiceID: testServiceID,
		})
	})

	// The API returns a 400 if Product is not enabled.
	// The API client returns an error if a non-2xx is returned from the API.
	if err == nil {
		t.Fatal("expected a 400 from the API but got a 2xx")
	}
}

func TestClient_GetProduct_validation(t *testing.T) {
	var err error

	_, err = testClient.GetProduct(&ProductEnablementInput{
		ProductID: ProductBrotliCompression,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetProduct(&ProductEnablementInput{
		ServiceID: "foo",
	})
	if err != ErrMissingProductID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_EnableProduct_validation(t *testing.T) {
	var err error
	_, err = testClient.EnableProduct(&ProductEnablementInput{
		ProductID: ProductBrotliCompression,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.EnableProduct(&ProductEnablementInput{
		ServiceID: "foo",
	})
	if err != ErrMissingProductID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DisableProduct_validation(t *testing.T) {
	var err error

	err = testClient.DisableProduct(&ProductEnablementInput{
		ProductID: ProductBrotliCompression,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DisableProduct(&ProductEnablementInput{
		ServiceID: "foo",
	})
	if err != ErrMissingProductID {
		t.Errorf("bad error: %s", err)
	}
}
