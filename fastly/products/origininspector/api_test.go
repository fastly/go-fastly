package origininspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/origininspector"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          origininspector.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[origininspector.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          origininspector.Get,
		ProductID:     origininspector.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[origininspector.EnableOutput, products.ServiceNullInput]{
		OpNoInputFn: origininspector.Enable,
		ProductID:   origininspector.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[origininspector.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      origininspector.Get,
		ProductID: origininspector.ProductID,
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: origininspector.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[origininspector.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          origininspector.Get,
		ProductID:     origininspector.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}

func TestEnablementCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestComputeServiceID)
}
