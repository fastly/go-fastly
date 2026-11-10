package airuntimecontrol_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/airuntimecontrol"
	"github.com/fastly/go-fastly/v19/internal/productcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          airuntimecontrol.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[airuntimecontrol.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          airuntimecontrol.Get,
		ProductID:     airuntimecontrol.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[airuntimecontrol.EnableOutput, products.NullInput]{
		OpNoInputFn: airuntimecontrol.Enable,
		ProductID:   airuntimecontrol.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[airuntimecontrol.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      airuntimecontrol.Get,
		ProductID: airuntimecontrol.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: airuntimecontrol.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[airuntimecontrol.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          airuntimecontrol.Get,
		ProductID:     airuntimecontrol.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
