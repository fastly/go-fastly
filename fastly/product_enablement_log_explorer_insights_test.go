package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_ProductEnablement_log_explorer_insights(t *testing.T) {
	t.Parallel()

	var err error

	// Enable Product - Log Explorer & Insights
	var pe *ProductEnablement
	Record(t, "product_enablement/enable_log_explorer_insights", func(c *Client) {
		pe, err = c.EnableProduct(context.TODO(), &ProductEnablementInput{
			ProductID: ProductLogExplorerInsights,
			ServiceID: TestDeliveryServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *pe.Product.ProductID != ProductLogExplorerInsights.String() {
		t.Errorf("bad feature_revision: %s", *pe.Product.ProductID)
	}

	// Get Product status
	var gpe *ProductEnablement
	Record(t, "product_enablement/get_log_explorer_insights", func(c *Client) {
		gpe, err = c.GetProduct(context.TODO(), &ProductEnablementInput{
			ProductID: ProductLogExplorerInsights,
			ServiceID: TestDeliveryServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *gpe.Product.ProductID != ProductLogExplorerInsights.String() {
		t.Errorf("bad feature_revision: %s", *gpe.Product.ProductID)
	}

	// Disable Product
	Record(t, "product_enablement/disable_log_explorer_insights", func(c *Client) {
		err = c.DisableProduct(context.TODO(), &ProductEnablementInput{
			ProductID: ProductLogExplorerInsights,
			ServiceID: TestDeliveryServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get Product status again to check disabled
	Record(t, "product_enablement/get-disabled_log_explorer_insights", func(c *Client) {
		gpe, err = c.GetProduct(context.TODO(), &ProductEnablementInput{
			ProductID: ProductLogExplorerInsights,
			ServiceID: TestDeliveryServiceID,
		})
	})

	// The API returns a 400 if Product is not enabled.
	// The API client returns an error if a non-2xx is returned from the API.
	if err == nil {
		t.Fatal("expected a 400 from the API but got a 2xx")
	}
}

func TestClient_GetProduct_validation_log_explorer_insights(t *testing.T) {
	var err error

	_, err = TestClient.GetProduct(context.TODO(), &ProductEnablementInput{
		ProductID: ProductLogExplorerInsights,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetProduct(context.TODO(), &ProductEnablementInput{
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingProductID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_EnableProduct_validation_log_explorer_insights(t *testing.T) {
	var err error
	_, err = TestClient.EnableProduct(context.TODO(), &ProductEnablementInput{
		ProductID: ProductLogExplorerInsights,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.EnableProduct(context.TODO(), &ProductEnablementInput{
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingProductID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DisableProduct_validation_log_explorer_insights(t *testing.T) {
	var err error

	err = TestClient.DisableProduct(context.TODO(), &ProductEnablementInput{
		ProductID: ProductLogExplorerInsights,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DisableProduct(context.TODO(), &ProductEnablementInput{
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingProductID) {
		t.Errorf("bad error: %s", err)
	}
}
