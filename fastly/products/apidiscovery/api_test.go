package apidiscovery_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/apidiscovery"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          apidiscovery.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[apidiscovery.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          apidiscovery.Get,
		ProductID:     apidiscovery.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[apidiscovery.EnableOutput, products.ServiceNullInput]{
		OpNoInputFn: apidiscovery.Enable,
		ProductID:   apidiscovery.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[apidiscovery.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      apidiscovery.Get,
		ProductID: apidiscovery.ProductID,
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: apidiscovery.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[apidiscovery.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          apidiscovery.Get,
		ProductID:     apidiscovery.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}
