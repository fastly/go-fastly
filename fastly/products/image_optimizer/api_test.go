package image_optimizer_test

import (
	"testing"
	"github.com/fastly/go-fastly/v9/fastly"
	p "github.com/fastly/go-fastly/v9/internal/products"
	"github.com/fastly/go-fastly/v9/fastly/products/image_optimizer"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTestInput{
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "before enablement",
		Executor: image_optimizer.Get,
		ServiceID: serviceID,
		ExpectFailure: true,
	}),
	p.TestEnable(&p.TestEnableInput[p.NullInput, fastly.ProductEnablement]{
		ExecutorNoInput: image_optimizer.Enable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "after enablement",
		Executor: image_optimizer.Get,
		ServiceID: serviceID,
	}),
	p.TestDisable(&p.TestDisableInput{
		Executor: image_optimizer.Disable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase: "after disablement",
		Executor: image_optimizer.Get,
		ServiceID: serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
