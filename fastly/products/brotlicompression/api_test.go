package brotlicompression_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/brotlicompression"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase: "ensure disabled before testing",
		OpFn:  brotlicompression.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[brotlicompression.ServiceEnableOutput]{
		Phase:         "before enablement",
		OpFn:          brotlicompression.Get,
		ProductID:     brotlicompression.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[brotlicompression.ServiceEnableOutput, products.ServiceNullInput]{
		OpNoInputFn: brotlicompression.Enable,
		ProductID:   brotlicompression.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[brotlicompression.ServiceEnableOutput]{
		Phase:     "after enablement",
		OpFn:      brotlicompression.Get,
		ProductID: brotlicompression.ProductID,
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: brotlicompression.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[brotlicompression.ServiceEnableOutput]{
		Phase:         "after disablement",
		OpFn:          brotlicompression.Get,
		ProductID:     brotlicompression.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}
