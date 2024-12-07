package origin_inspector_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/origin_inspector"
	p "github.com/fastly/go-fastly/v9/internal/products"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTestInput{
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "before enablement",
		Executor:      origin_inspector.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	p.TestEnable(&p.TestEnableInput[p.NullInput, fastly.ProductEnablement]{
		ExecutorNoInput: origin_inspector.Enable,
		ServiceID:       serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:     "after enablement",
		Executor:  origin_inspector.Get,
		ServiceID: serviceID,
	}),
	p.TestDisable(&p.TestDisableInput{
		Executor:  origin_inspector.Disable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "after disablement",
		Executor:      origin_inspector.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
