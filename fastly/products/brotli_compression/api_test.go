package brotli_compression_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/brotli_compression"
	"github.com/fastly/go-fastly/v9/internal/productcore"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          brotli_compression.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          brotli_compression.Get,
		ProductID:     brotli_compression.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[*products.EnableOutput, *productcore.NullInput]{
		OpNoInputFn: brotli_compression.Enable,
		ProductID:   brotli_compression.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      brotli_compression.Get,
		ProductID: brotli_compression.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      brotli_compression.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          brotli_compression.Get,
		ProductID:     brotli_compression.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
