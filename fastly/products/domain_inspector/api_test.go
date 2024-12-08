package domain_inspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	// tp is 'this product' package
	tp "github.com/fastly/go-fastly/v9/fastly/products/domain_inspector"
	// fp is 'fastly products' package
	fp "github.com/fastly/go-fastly/v9/fastly/products"
	// ip is 'internal products' package
	ip "github.com/fastly/go-fastly/v9/internal/products"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTest{
	ip.NewDisableTest(&ip.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          tp.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	ip.NewGetTest(&ip.GetTestInput[*fp.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          tp.Get,
		ProductID:     tp.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	ip.NewEnableTest(&ip.EnableTestInput[*ip.NullInput, *fp.EnableOutput]{
		OpNoInputFn: tp.Enable,
		ProductID:   tp.ProductID,
		ServiceID:   serviceID,
	}),
	ip.NewGetTest(&ip.GetTestInput[*fp.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      tp.Get,
		ProductID: tp.ProductID,
		ServiceID: serviceID,
	}),
	ip.NewDisableTest(&ip.DisableTestInput{
		OpFn:      tp.Disable,
		ServiceID: serviceID,
	}),
	ip.NewGetTest(&ip.GetTestInput[*fp.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          tp.Get,
		ProductID:     tp.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
