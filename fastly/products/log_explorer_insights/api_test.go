package log_explorer_insights_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/log_explorer_insights"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          log_explorer_insights.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*log_explorer_insights.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          log_explorer_insights.Get,
		ProductID:     log_explorer_insights.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[*log_explorer_insights.EnableOutput, *productcore.NullInput]{
		OpNoInputFn: log_explorer_insights.Enable,
		ProductID:   log_explorer_insights.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*log_explorer_insights.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      log_explorer_insights.Get,
		ProductID: log_explorer_insights.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      log_explorer_insights.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*log_explorer_insights.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          log_explorer_insights.Get,
		ProductID:     log_explorer_insights.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
