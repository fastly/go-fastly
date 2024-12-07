package log_explorer_insights_test

import (
	"testing"
	"github.com/fastly/go-fastly/v9/fastly"
	p "github.com/fastly/go-fastly/v9/internal/products"
	"github.com/fastly/go-fastly/v9/fastly/products/log_explorer_insights"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTestInput{
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "before enablement",
		Executor: log_explorer_insights.Get,
		ServiceID: serviceID,
		ExpectFailure: true,
	}),
	p.TestEnable(&p.TestEnableInput[p.NullInput, fastly.ProductEnablement]{
		ExecutorNoInput: log_explorer_insights.Enable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "after enablement",
		Executor: log_explorer_insights.Get,
		ServiceID: serviceID,
	}),
	p.TestDisable(&p.TestDisableInput{
		Executor: log_explorer_insights.Disable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "after disablement",
		Executor: log_explorer_insights.Get,
		ServiceID: serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}