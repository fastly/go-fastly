package domain_inspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/domain_inspector"
	"github.com/fastly/go-fastly/v9/internal/productcore"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          domain_inspector.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          domain_inspector.Get,
		ProductID:     domain_inspector.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[*products.EnableOutput, *productcore.NullInput]{
		OpNoInputFn: domain_inspector.Enable,
		ProductID:   domain_inspector.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      domain_inspector.Get,
		ProductID: domain_inspector.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      domain_inspector.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          domain_inspector.Get,
		ProductID:     domain_inspector.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
