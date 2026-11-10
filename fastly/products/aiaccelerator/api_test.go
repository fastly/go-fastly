package aiaccelerator_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/aiaccelerator"
	"github.com/fastly/go-fastly/v19/internal/productcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          aiaccelerator.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[aiaccelerator.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          aiaccelerator.Get,
		ProductID:     aiaccelerator.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[aiaccelerator.EnableOutput, products.NullInput]{
		OpNoInputFn: aiaccelerator.Enable,
		ProductID:   aiaccelerator.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[aiaccelerator.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      aiaccelerator.Get,
		ProductID: aiaccelerator.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: aiaccelerator.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[aiaccelerator.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          aiaccelerator.Get,
		ProductID:     aiaccelerator.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
