package fanout_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/fanout"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          fanout.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[fanout.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          fanout.Get,
		ProductID:     fanout.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[fanout.EnableOutput, products.ServiceNullInput]{
		OpNoInputFn: fanout.Enable,
		ProductID:   fanout.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[fanout.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      fanout.Get,
		ProductID: fanout.ProductID,
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: fanout.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[fanout.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          fanout.Get,
		ProductID:     fanout.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestComputeServiceID)
}
