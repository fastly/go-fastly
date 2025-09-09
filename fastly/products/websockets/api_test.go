package websockets_test

import (
	"testing"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/products"
	"github.com/fastly/go-fastly/v11/fastly/products/websockets"
	"github.com/fastly/go-fastly/v11/internal/productcore"
	"github.com/fastly/go-fastly/v11/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          websockets.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[websockets.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          websockets.Get,
		ProductID:     websockets.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[websockets.EnableOutput, products.NullInput]{
		OpNoInputFn: websockets.Enable,
		ProductID:   websockets.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[websockets.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      websockets.Get,
		ProductID: websockets.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: websockets.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[websockets.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          websockets.Get,
		ProductID:     websockets.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}

func TestEnablementCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestComputeServiceID)
}
