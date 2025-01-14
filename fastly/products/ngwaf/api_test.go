package ngwaf_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/ngwaf"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"

	"github.com/stretchr/testify/require"
)

var serviceID = fastly.TestDeliveryServiceID

func TestEnableMissingWorkspaceID(t *testing.T) {
	t.Parallel()

	_, err := ngwaf.Enable(nil, serviceID, ngwaf.EnableInput{WorkspaceID: ""})

	require.ErrorIs(t, err, ngwaf.ErrMissingWorkspaceID)
}

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          ngwaf.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ngwaf.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          ngwaf.Get,
		ProductID:     ngwaf.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[ngwaf.EnableOutput, ngwaf.EnableInput]{
		OpWithInputFn: ngwaf.Enable,
		Input:         ngwaf.EnableInput{WorkspaceID: fastly.TestNGWAFWorkspaceID},
		ProductID:     ngwaf.ProductID,
		ServiceID:     serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ngwaf.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      ngwaf.Get,
		ProductID: ngwaf.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[ngwaf.ConfigureOutput]{
		Phase:     "default",
		OpFn:      ngwaf.GetConfiguration,
		ProductID: ngwaf.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, output ngwaf.ConfigureOutput) {
			require.NotNilf(t, output.Configuration.TrafficRamp, "test '%s'", tc.Name)
			require.Equalf(t, "100", *output.Configuration.TrafficRamp, "test '%s'", tc.Name)
		},
	}),
	productcore.NewUpdateConfigurationTest(&productcore.UpdateConfigurationTestInput[ngwaf.ConfigureOutput, ngwaf.ConfigureInput]{
		OpFn:      ngwaf.UpdateConfiguration,
		Input:     ngwaf.ConfigureInput{TrafficRamp: "45"},
		ProductID: ngwaf.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, output ngwaf.ConfigureOutput) {
			require.NotNilf(t, output.Configuration.TrafficRamp, "test '%s'", tc.Name)
			require.Equalf(t, "45", *output.Configuration.TrafficRamp, "test '%s'", tc.Name)
		},
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[ngwaf.ConfigureOutput]{
		Phase:     "after update",
		OpFn:      ngwaf.GetConfiguration,
		ProductID: ngwaf.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, output ngwaf.ConfigureOutput) {
			require.NotNilf(t, output.Configuration.TrafficRamp, "test '%s'", tc.Name)
			require.Equalf(t, "45", *output.Configuration.TrafficRamp, "test '%s'", tc.Name)
		},
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      ngwaf.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ngwaf.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          ngwaf.Get,
		ProductID:     ngwaf.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablementAndConfiguration(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
