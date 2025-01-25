package logexplorerinsights_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/logexplorerinsights"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          logexplorerinsights.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          logexplorerinsights.Get,
		ProductID:     logexplorerinsights.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[logexplorerinsights.EnableOutput, products.NullInput]{
		OpNoInputFn: logexplorerinsights.Enable,
		ProductID:   logexplorerinsights.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      logexplorerinsights.Get,
		ProductID: logexplorerinsights.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      logexplorerinsights.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[logexplorerinsights.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          logexplorerinsights.Get,
		ProductID:     logexplorerinsights.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
