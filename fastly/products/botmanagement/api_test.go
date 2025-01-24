package botmanagement_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/botmanagement"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          botmanagement.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          botmanagement.Get,
		ProductID:     botmanagement.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[botmanagement.EnableOutput, products.NullInput]{
		OpNoInputFn: botmanagement.Enable,
		ProductID:   botmanagement.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      botmanagement.Get,
		ProductID: botmanagement.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      botmanagement.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          botmanagement.Get,
		ProductID:     botmanagement.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
