package logexplorerinsights_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/logexplorerinsights"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          logexplorerinsights.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          logexplorerinsights.Get,
		ProductID:     logexplorerinsights.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[logexplorerinsights.EnableOutput, products.ServiceNullInput]{
		OpNoInputFn: logexplorerinsights.Enable,
		ProductID:   logexplorerinsights.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      logexplorerinsights.Get,
		ProductID: logexplorerinsights.ProductID,
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: logexplorerinsights.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          logexplorerinsights.Get,
		ProductID:     logexplorerinsights.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}

func TestEnablementCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestComputeServiceID)
}
