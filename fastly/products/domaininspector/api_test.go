package domaininspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/domaininspector"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          domaininspector.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          domaininspector.Get,
		ProductID:     domaininspector.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[domaininspector.EnableOutput, products.ServiceNullInput]{
		OpNoInputFn: domaininspector.Enable,
		ProductID:   domaininspector.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      domaininspector.Get,
		ProductID: domaininspector.ProductID,
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: domaininspector.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          domaininspector.Get,
		ProductID:     domaininspector.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}
