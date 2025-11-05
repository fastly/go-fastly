package brotlicompression_test

import (
	"testing"

	"github.com/fastly/go-fastly/v12/fastly"
	"github.com/fastly/go-fastly/v12/fastly/products"
	"github.com/fastly/go-fastly/v12/fastly/products/brotlicompression"
	"github.com/fastly/go-fastly/v12/internal/productcore"
	"github.com/fastly/go-fastly/v12/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase: "ensure disabled before testing",
		OpFn:  brotlicompression.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[brotlicompression.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          brotlicompression.Get,
		ProductID:     brotlicompression.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[brotlicompression.EnableOutput, products.NullInput]{
		OpNoInputFn: brotlicompression.Enable,
		ProductID:   brotlicompression.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[brotlicompression.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      brotlicompression.Get,
		ProductID: brotlicompression.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: brotlicompression.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[brotlicompression.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          brotlicompression.Get,
		ProductID:     brotlicompression.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}
