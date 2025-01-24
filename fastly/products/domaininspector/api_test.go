package domaininspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/domaininspector"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          domaininspector.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          domaininspector.Get,
		ProductID:     domaininspector.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[domaininspector.EnableOutput, products.NullInput]{
		OpNoInputFn: domaininspector.Enable,
		ProductID:   domaininspector.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      domaininspector.Get,
		ProductID: domaininspector.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      domaininspector.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domaininspector.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          domaininspector.Get,
		ProductID:     domaininspector.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
