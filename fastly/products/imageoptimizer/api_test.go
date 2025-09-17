package imageoptimizer_test

import (
	"testing"

	"github.com/fastly/go-fastly/v12/fastly"
	"github.com/fastly/go-fastly/v12/fastly/products"
	"github.com/fastly/go-fastly/v12/fastly/products/imageoptimizer"
	"github.com/fastly/go-fastly/v12/internal/productcore"
	"github.com/fastly/go-fastly/v12/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          imageoptimizer.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          imageoptimizer.Get,
		ProductID:     imageoptimizer.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[imageoptimizer.EnableOutput, products.NullInput]{
		OpNoInputFn: imageoptimizer.Enable,
		ProductID:   imageoptimizer.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      imageoptimizer.Get,
		ProductID: imageoptimizer.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: imageoptimizer.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          imageoptimizer.Get,
		ProductID:     imageoptimizer.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}
