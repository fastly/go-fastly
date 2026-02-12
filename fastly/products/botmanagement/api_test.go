package botmanagement_test

import (
	"testing"

	"github.com/fastly/go-fastly/v13/fastly"
	"github.com/fastly/go-fastly/v13/fastly/products"
	"github.com/fastly/go-fastly/v13/fastly/products/botmanagement"
	"github.com/fastly/go-fastly/v13/internal/productcore"
	"github.com/fastly/go-fastly/v13/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          botmanagement.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          botmanagement.Get,
		ProductID:     botmanagement.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[botmanagement.EnableOutput, products.NullInput]{
		OpNoInputFn: botmanagement.Enable,
		ProductID:   botmanagement.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      botmanagement.Get,
		ProductID: botmanagement.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: botmanagement.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          botmanagement.Get,
		ProductID:     botmanagement.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}
