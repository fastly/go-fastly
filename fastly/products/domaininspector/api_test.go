package domaininspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v13/fastly"
	"github.com/fastly/go-fastly/v13/fastly/products"
	"github.com/fastly/go-fastly/v13/fastly/products/domaininspector"
	"github.com/fastly/go-fastly/v13/internal/productcore"
	"github.com/fastly/go-fastly/v13/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          domaininspector.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          domaininspector.Get,
		ProductID:     domaininspector.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[domaininspector.EnableOutput, products.NullInput]{
		OpNoInputFn: domaininspector.Enable,
		ProductID:   domaininspector.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      domaininspector.Get,
		ProductID: domaininspector.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: domaininspector.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          domaininspector.Get,
		ProductID:     domaininspector.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}
