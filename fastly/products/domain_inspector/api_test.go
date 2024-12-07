package domain_inspector_test

import (
	"testing"
	"github.com/fastly/go-fastly/v9/fastly"
	p "github.com/fastly/go-fastly/v9/internal/products"
	"github.com/fastly/go-fastly/v9/fastly/products/domain_inspector"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTestInput{
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "before enablement",
		Executor: domain_inspector.Get,
		ServiceID: serviceID,
		ExpectFailure: true,
	}),
	p.TestEnable(&p.TestEnableInput[p.NullInput, fastly.ProductEnablement]{
		ExecutorNoInput: domain_inspector.Enable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "after enablement",
		Executor: domain_inspector.Get,
		ServiceID: serviceID,
	}),
	p.TestDisable(&p.TestDisableInput{
		Executor: domain_inspector.Disable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "after disablement",
		Executor: domain_inspector.Get,
		ServiceID: serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
