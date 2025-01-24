package imageoptimizer_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/imageoptimizer"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          imageoptimizer.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          imageoptimizer.Get,
		ProductID:     imageoptimizer.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[imageoptimizer.EnableOutput, products.NullInput]{
		OpNoInputFn: imageoptimizer.Enable,
		ProductID:   imageoptimizer.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      imageoptimizer.Get,
		ProductID: imageoptimizer.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      imageoptimizer.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          imageoptimizer.Get,
		ProductID:     imageoptimizer.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
