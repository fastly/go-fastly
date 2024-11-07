package fastly

import (
	"testing"
)

func TestClient_ProductEnablement_fanout(t *testing.T) {
	t.Parallel()

	var err error

	// Enable Product - Bot Management
	var pe *ProductEnablement
	record(t, "product_enablement/enable_fanout", func(c *Client) {
		pe, err = c.EnableProduct(&ProductEnablementInput{
			ProductID: ProductFanout,
			ServiceID: testComputeServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *pe.Product.ProductID != ProductFanout.String() {
		t.Errorf("bad feature_revision: %s", *pe.Product.ProductID)
	}

	// Get Product status
	var gpe *ProductEnablement
	record(t, "product_enablement/get_fanout", func(c *Client) {
		gpe, err = c.GetProduct(&ProductEnablementInput{
			ProductID: ProductFanout,
			ServiceID: testComputeServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *gpe.Product.ProductID != ProductFanout.String() {
		t.Errorf("bad feature_revision: %s", *gpe.Product.ProductID)
	}

	// Disable Product
	record(t, "product_enablement/disable_fanout", func(c *Client) {
		err = c.DisableProduct(&ProductEnablementInput{
			ProductID: ProductFanout,
			ServiceID: testComputeServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get Product status again to check disabled
	record(t, "product_enablement/get-disabled_fanout", func(c *Client) {
		gpe, err = c.GetProduct(&ProductEnablementInput{
			ProductID: ProductFanout,
			ServiceID: testComputeServiceID,
		})
	})

	// The API returns a 400 if Product is not enabled.
	// The API client returns an error if a non-2xx is returned from the API.
	if err == nil {
		t.Fatal("expected a 400 from the API but got a 2xx")
	}
}

func TestClient_GetProduct_validation_fanout(t *testing.T) {
	var err error

	_, err = testClient.GetProduct(&ProductEnablementInput{
		ProductID: ProductFanout,
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

func TestClient_EnableProduct_validation_fanout(t *testing.T) {
	var err error
	_, err = testClient.EnableProduct(&ProductEnablementInput{
		ProductID: ProductFanout,
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

func TestClient_DisableProduct_validation_fanout(t *testing.T) {
	var err error

	err = testClient.DisableProduct(&ProductEnablementInput{
		ProductID: ProductFanout,
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
