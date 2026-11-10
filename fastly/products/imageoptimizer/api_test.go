package imageoptimizer_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/imageoptimizer"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          imageoptimizer.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          imageoptimizer.Get,
		ProductID:     imageoptimizer.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[imageoptimizer.EnableOutput, products.ServiceNullInput]{
		OpNoInputFn: imageoptimizer.Enable,
		ProductID:   imageoptimizer.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      imageoptimizer.Get,
		ProductID: imageoptimizer.ProductID,
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: imageoptimizer.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[imageoptimizer.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          imageoptimizer.Get,
		ProductID:     imageoptimizer.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}
