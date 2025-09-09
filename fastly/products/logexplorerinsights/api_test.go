package logexplorerinsights_test

import (
	"testing"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/products"
	"github.com/fastly/go-fastly/v11/fastly/products/logexplorerinsights"
	"github.com/fastly/go-fastly/v11/internal/productcore"
	"github.com/fastly/go-fastly/v11/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          logexplorerinsights.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          logexplorerinsights.Get,
		ProductID:     logexplorerinsights.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[logexplorerinsights.EnableOutput, products.NullInput]{
		OpNoInputFn: logexplorerinsights.Enable,
		ProductID:   logexplorerinsights.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      logexplorerinsights.Get,
		ProductID: logexplorerinsights.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: logexplorerinsights.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          logexplorerinsights.Get,
		ProductID:     logexplorerinsights.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}

func TestEnablementCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestComputeServiceID)
}
