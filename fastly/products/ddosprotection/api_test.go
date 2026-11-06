package ddosprotection_test

import (
	"context"
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products/ddosprotection"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"

	"github.com/stretchr/testify/require"
)

var serviceID = fastly.TestDeliveryServiceID

func TestUpdateConfigurationMissingMode(t *testing.T) {
	t.Parallel()

	_, err := ddosprotection.UpdateConfiguration(context.TODO(), nil, serviceID, ddosprotection.ConfigureInput{Mode: ""})

	require.ErrorIs(t, err, ddosprotection.ErrMissingMode)
}

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          ddosprotection.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[ddosprotection.ServiceEnableOutput]{
		Phase:         "before enablement",
		OpFn:          ddosprotection.Get,
		ProductID:     ddosprotection.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[ddosprotection.ServiceEnableOutput, ddosprotection.EnableInput]{
		OpWithInputFn: ddosprotection.Enable,
		Input:         ddosprotection.EnableInput{},
		ProductID:     ddosprotection.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[ddosprotection.ServiceEnableOutput]{
		Phase:     "after enablement",
		OpFn:      ddosprotection.Get,
		ProductID: ddosprotection.ProductID,
	}),
	serviceproductcore.NewGetConfigurationTest(&serviceproductcore.GetConfigurationTestInput[ddosprotection.ServiceConfigureOutput]{
		Phase:     "default",
		OpFn:      ddosprotection.GetConfiguration,
		ProductID: ddosprotection.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, o ddosprotection.ServiceConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "log", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewUpdateConfigurationTest(&serviceproductcore.UpdateConfigurationTestInput[ddosprotection.ServiceConfigureOutput, ddosprotection.ConfigureInput]{
		OpFn:      ddosprotection.UpdateConfiguration,
		Input:     ddosprotection.ConfigureInput{Mode: "block"},
		ProductID: ddosprotection.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, o ddosprotection.ServiceConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewGetConfigurationTest(&serviceproductcore.GetConfigurationTestInput[ddosprotection.ServiceConfigureOutput]{
		Phase:     "after update",
		OpFn:      ddosprotection.GetConfiguration,
		ProductID: ddosprotection.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, o ddosprotection.ServiceConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: ddosprotection.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[ddosprotection.ServiceEnableOutput]{
		Phase:         "after disablement",
		OpFn:          ddosprotection.Get,
		ProductID:     ddosprotection.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[ddosprotection.ServiceEnableOutput, ddosprotection.EnableInput]{
		Phase:         "enable with explicit mode",
		OpWithInputFn: ddosprotection.Enable,
		Input:         ddosprotection.EnableInput{Mode: "block"},
		ProductID:     ddosprotection.ProductID,
	}),
	serviceproductcore.NewGetConfigurationTest(&serviceproductcore.GetConfigurationTestInput[ddosprotection.ServiceConfigureOutput]{
		Phase:     "verify mode parameter",
		OpFn:      ddosprotection.GetConfiguration,
		ProductID: ddosprotection.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, o ddosprotection.ServiceConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase: "cleanup after mode test",
		OpFn:  ddosprotection.Disable,
	}),
}

func TestEnablementAndConfigurationDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}

func TestEnablementAndConfigurationCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestComputeServiceID)
}
