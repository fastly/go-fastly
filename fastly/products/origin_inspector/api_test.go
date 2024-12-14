package origin_inspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/origin_inspector"
	"github.com/fastly/go-fastly/v9/internal/productcore"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          origin_inspector.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          origin_inspector.Get,
		ProductID:     origin_inspector.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[*products.EnableOutput, *productcore.NullInput]{
		OpNoInputFn: origin_inspector.Enable,
		ProductID:   origin_inspector.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      origin_inspector.Get,
		ProductID: origin_inspector.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      origin_inspector.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          origin_inspector.Get,
		ProductID:     origin_inspector.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
