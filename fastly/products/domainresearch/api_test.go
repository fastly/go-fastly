package domainresearch_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/domainresearch"
	"github.com/fastly/go-fastly/v19/internal/productcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          domainresearch.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domainresearch.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          domainresearch.Get,
		ProductID:     domainresearch.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[domainresearch.EnableOutput, products.NullInput]{
		OpNoInputFn: domainresearch.Enable,
		ProductID:   domainresearch.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domainresearch.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      domainresearch.Get,
		ProductID: domainresearch.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: domainresearch.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[domainresearch.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          domainresearch.Get,
		ProductID:     domainresearch.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
