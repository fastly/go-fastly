package ngwaf_test

import (
	"context"
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products/ngwaf"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"

	"github.com/stretchr/testify/require"
)

func TestEnableMissingWorkspaceIDDelivery(t *testing.T) {
	t.Parallel()

	_, err := ngwaf.Enable(context.TODO(), nil, fastly.TestDeliveryServiceID, ngwaf.EnableInput{WorkspaceID: ""})

	require.ErrorIs(t, err, ngwaf.ErrMissingWorkspaceID)
}

func TestEnableMissingWorkspaceIDCompute(t *testing.T) {
	t.Parallel()

	_, err := ngwaf.Enable(context.TODO(), nil, fastly.TestComputeServiceID, ngwaf.EnableInput{WorkspaceID: ""})

	require.ErrorIs(t, err, ngwaf.ErrMissingWorkspaceID)
}

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          ngwaf.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[ngwaf.ServiceEnableOutput]{
		Phase:         "before enablement",
		OpFn:          ngwaf.Get,
		ProductID:     ngwaf.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[ngwaf.ServiceEnableOutput, ngwaf.EnableInput]{
		OpWithInputFn: ngwaf.Enable,
		Input:         ngwaf.EnableInput{WorkspaceID: fastly.TestNGWAFWorkspaceID},
		ProductID:     ngwaf.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[ngwaf.ServiceEnableOutput]{
		Phase:     "after enablement",
		OpFn:      ngwaf.Get,
		ProductID: ngwaf.ProductID,
	}),
	serviceproductcore.NewGetConfigurationTest(&serviceproductcore.GetConfigurationTestInput[ngwaf.ServiceConfigureOutput]{
		Phase:     "default",
		OpFn:      ngwaf.GetConfiguration,
		ProductID: ngwaf.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, output ngwaf.ServiceConfigureOutput) {
			require.NotNilf(t, output.Configuration.TrafficRamp, "test '%s'", tc.Name)
			require.Equalf(t, "100", *output.Configuration.TrafficRamp, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewUpdateConfigurationTest(&serviceproductcore.UpdateConfigurationTestInput[ngwaf.ServiceConfigureOutput, ngwaf.ConfigureInput]{
		OpFn:      ngwaf.UpdateConfiguration,
		Input:     ngwaf.ConfigureInput{TrafficRamp: "45"},
		ProductID: ngwaf.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, output ngwaf.ServiceConfigureOutput) {
			require.NotNilf(t, output.Configuration.TrafficRamp, "test '%s'", tc.Name)
			require.Equalf(t, "45", *output.Configuration.TrafficRamp, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewGetConfigurationTest(&serviceproductcore.GetConfigurationTestInput[ngwaf.ServiceConfigureOutput]{
		Phase:     "after update",
		OpFn:      ngwaf.GetConfiguration,
		ProductID: ngwaf.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, output ngwaf.ServiceConfigureOutput) {
			require.NotNilf(t, output.Configuration.TrafficRamp, "test '%s'", tc.Name)
			require.Equalf(t, "45", *output.Configuration.TrafficRamp, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: ngwaf.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[ngwaf.ServiceEnableOutput]{
		Phase:         "after disablement",
		OpFn:          ngwaf.Get,
		ProductID:     ngwaf.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}
