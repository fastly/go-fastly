package origininspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v13/fastly"
	"github.com/fastly/go-fastly/v13/fastly/products"
	"github.com/fastly/go-fastly/v13/fastly/products/origininspector"
	"github.com/fastly/go-fastly/v13/internal/productcore"
	"github.com/fastly/go-fastly/v13/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          origininspector.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[origininspector.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          origininspector.Get,
		ProductID:     origininspector.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[origininspector.EnableOutput, products.NullInput]{
		OpNoInputFn: origininspector.Enable,
		ProductID:   origininspector.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[origininspector.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      origininspector.Get,
		ProductID: origininspector.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: origininspector.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[origininspector.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          origininspector.Get,
		ProductID:     origininspector.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}

func TestEnablementCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestComputeServiceID)
}
