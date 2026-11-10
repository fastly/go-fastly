package websockets_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/websockets"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          websockets.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[websockets.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          websockets.Get,
		ProductID:     websockets.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[websockets.EnableOutput, products.ServiceNullInput]{
		OpNoInputFn: websockets.Enable,
		ProductID:   websockets.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[websockets.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      websockets.Get,
		ProductID: websockets.ProductID,
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: websockets.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[websockets.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          websockets.Get,
		ProductID:     websockets.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}

func TestEnablementCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestComputeServiceID)
}
