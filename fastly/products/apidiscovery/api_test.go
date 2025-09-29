package apidiscovery_test

import (
	"testing"

	"github.com/fastly/go-fastly/v12/fastly"
	"github.com/fastly/go-fastly/v12/fastly/products"
	"github.com/fastly/go-fastly/v12/fastly/products/apidiscovery"
	"github.com/fastly/go-fastly/v12/internal/productcore"
	"github.com/fastly/go-fastly/v12/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          apidiscovery.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[apidiscovery.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          apidiscovery.Get,
		ProductID:     apidiscovery.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[apidiscovery.EnableOutput, products.NullInput]{
		OpNoInputFn: apidiscovery.Enable,
		ProductID:   apidiscovery.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[apidiscovery.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      apidiscovery.Get,
		ProductID: apidiscovery.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: apidiscovery.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[apidiscovery.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          apidiscovery.Get,
		ProductID:     apidiscovery.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}
