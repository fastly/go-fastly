package websockets_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/websockets"
	p "github.com/fastly/go-fastly/v9/internal/products"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTestInput{
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "before enablement",
		Executor:      websockets.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	p.TestEnable(&p.TestEnableInput[p.NullInput, fastly.ProductEnablement]{
		ExecutorNoInput: websockets.Enable,
		ServiceID:       serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:     "after enablement",
		Executor:  websockets.Get,
		ServiceID: serviceID,
	}),
	p.TestDisable(&p.TestDisableInput{
		Executor:  websockets.Disable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "after disablement",
		Executor:      websockets.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
