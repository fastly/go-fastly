package brotlicompression_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/brotlicompression"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          brotlicompression.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[brotlicompression.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          brotlicompression.Get,
		ProductID:     brotlicompression.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[brotlicompression.EnableOutput, products.NullInput]{
		OpNoInputFn: brotlicompression.Enable,
		ProductID:   brotlicompression.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[brotlicompression.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      brotlicompression.Get,
		ProductID: brotlicompression.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      brotlicompression.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[brotlicompression.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          brotlicompression.Get,
		ProductID:     brotlicompression.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
