package fastly

import (
	"errors"
	"testing"
)

func TestClient_ProductEnablement_domain_inspector(t *testing.T) {
	t.Parallel()

	var err error

	// Enable Product - Bot Management
	var pe *ProductEnablement
	Record(t, "product_enablement/enable_domain_inspector", func(c *Client) {
		pe, err = c.EnableProduct(&ProductEnablementInput{
			ProductID: ProductDomainInspector,
			ServiceID: TestDeliveryServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *pe.Product.ProductID != ProductDomainInspector.String() {
		t.Errorf("bad feature_revision: %s", *pe.Product.ProductID)
	}

	// Get Product status
	var gpe *ProductEnablement
	Record(t, "product_enablement/get_domain_inspector", func(c *Client) {
		gpe, err = c.GetProduct(&ProductEnablementInput{
			ProductID: ProductDomainInspector,
			ServiceID: TestDeliveryServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *gpe.Product.ProductID != ProductDomainInspector.String() {
		t.Errorf("bad feature_revision: %s", *gpe.Product.ProductID)
	}

	// Disable Product
	Record(t, "product_enablement/disable_domain_inspector", func(c *Client) {
		err = c.DisableProduct(&ProductEnablementInput{
			ProductID: ProductDomainInspector,
			ServiceID: TestDeliveryServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get Product status again to check disabled
	Record(t, "product_enablement/get-disabled_domain_inspector", func(c *Client) {
		gpe, err = c.GetProduct(&ProductEnablementInput{
			ProductID: ProductDomainInspector,
			ServiceID: TestDeliveryServiceID,
		})
	})

	// The API returns a 400 if Product is not enabled.
	// The API client returns an error if a non-2xx is returned from the API.
	if err == nil {
		t.Fatal("expected a 400 from the API but got a 2xx")
	}
}

func TestClient_GetProduct_validation_domain_inspector(t *testing.T) {
	var err error

	_, err = TestClient.GetProduct(&ProductEnablementInput{
		ProductID: ProductDomainInspector,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetProduct(&ProductEnablementInput{
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingProductID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_EnableProduct_validation_domain_inspector(t *testing.T) {
	var err error
	_, err = TestClient.EnableProduct(&ProductEnablementInput{
		ProductID: ProductDomainInspector,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.EnableProduct(&ProductEnablementInput{
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingProductID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DisableProduct_validation_domain_inspector(t *testing.T) {
	var err error

	err = TestClient.DisableProduct(&ProductEnablementInput{
		ProductID: ProductDomainInspector,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DisableProduct(&ProductEnablementInput{
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingProductID) {
		t.Errorf("bad error: %s", err)
	}
}
