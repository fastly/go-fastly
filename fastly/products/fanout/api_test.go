package fanout_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/fanout"
	p "github.com/fastly/go-fastly/v9/internal/products"
)

var serviceID = fastly.TestComputeServiceID

var functionalTests = []*fastly.FunctionalTestInput{
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "before enablement",
		Executor:      fanout.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	p.TestEnable(&p.TestEnableInput[p.NullInput, fastly.ProductEnablement]{
		ExecutorNoInput: fanout.Enable,
		ServiceID:       serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:     "after enablement",
		Executor:  fanout.Get,
		ServiceID: serviceID,
	}),
	p.TestDisable(&p.TestDisableInput{
		Executor:  fanout.Disable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "after disablement",
		Executor:      fanout.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
