package ddos_protection_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	// tp is 'this product' package
	tp "github.com/fastly/go-fastly/v9/fastly/products/ddos_protection"
	// fp is 'fastly products' package
	fp "github.com/fastly/go-fastly/v9/fastly/products"
	// ip is 'internal products' package
	ip "github.com/fastly/go-fastly/v9/internal/products"

	"github.com/stretchr/testify/require"
)

var serviceID = fastly.TestDeliveryServiceID

func TestUpdateConfigurationMissingMode(t *testing.T) {
	t.Parallel()

	_, err := tp.UpdateConfiguration(nil, serviceID, &tp.ConfigureInput{Mode: ""})

	require.ErrorIs(t, err, tp.ErrMissingMode)
}

var functionalTests = []*fastly.FunctionalTest{
	ip.NewDisableTest(&ip.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          tp.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	ip.NewGetTest(&ip.GetTestInput[*fp.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          tp.Get,
		ProductID:     tp.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	ip.NewEnableTest(&ip.EnableTestInput[*ip.NullInput, *fp.EnableOutput]{
		OpNoInputFn: tp.Enable,
		ProductID:   tp.ProductID,
		ServiceID:   serviceID,
	}),
	ip.NewGetTest(&ip.GetTestInput[*fp.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      tp.Get,
		ProductID: tp.ProductID,
		ServiceID: serviceID,
	}),
	ip.NewGetConfigurationTest(&ip.GetConfigurationTestInput[*tp.ConfigureOutput]{
		Phase:     "default",
		OpFn:      tp.GetConfiguration,
		ProductID: tp.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *fastly.FunctionalTest, output *tp.ConfigureOutput) {
			require.NotNilf(t, output.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "log", *output.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	ip.NewUpdateConfigurationTest(&ip.UpdateConfigurationTestInput[*tp.ConfigureInput, *tp.ConfigureOutput]{
		OpFn:      tp.UpdateConfiguration,
		Input:     &tp.ConfigureInput{Mode: "block"},
		ProductID: tp.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *fastly.FunctionalTest, output *tp.ConfigureOutput) {
			require.NotNilf(t, output.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *output.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	ip.NewGetConfigurationTest(&ip.GetConfigurationTestInput[*tp.ConfigureOutput]{
		Phase:     "after update",
		OpFn:      tp.GetConfiguration,
		ProductID: tp.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *fastly.FunctionalTest, output *tp.ConfigureOutput) {
			require.NotNilf(t, output.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *output.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	ip.NewDisableTest(&ip.DisableTestInput{
		OpFn:      tp.Disable,
		ServiceID: serviceID,
	}),
	ip.NewGetTest(&ip.GetTestInput[*fp.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          tp.Get,
		ProductID:     tp.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablementAndConfiguration(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
