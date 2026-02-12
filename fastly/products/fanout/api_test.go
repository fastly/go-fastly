package fanout_test

import (
	"testing"

	"github.com/fastly/go-fastly/v13/fastly"
	"github.com/fastly/go-fastly/v13/fastly/products"
	"github.com/fastly/go-fastly/v13/fastly/products/fanout"
	"github.com/fastly/go-fastly/v13/internal/productcore"
	"github.com/fastly/go-fastly/v13/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          fanout.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[fanout.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          fanout.Get,
		ProductID:     fanout.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[fanout.EnableOutput, products.NullInput]{
		OpNoInputFn: fanout.Enable,
		ProductID:   fanout.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[fanout.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      fanout.Get,
		ProductID: fanout.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: fanout.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[fanout.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          fanout.Get,
		ProductID:     fanout.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestComputeServiceID)
}
