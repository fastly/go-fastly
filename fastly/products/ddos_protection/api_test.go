package ddos_protection_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/ddos_protection"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"

	"github.com/stretchr/testify/require"
)

var serviceID = fastly.TestDeliveryServiceID

func TestUpdateConfigurationMissingMode(t *testing.T) {
	t.Parallel()

	_, err := ddos_protection.UpdateConfiguration(nil, serviceID, &ddos_protection.ConfigureInput{Mode: ""})

	require.ErrorIs(t, err, ddos_protection.ErrMissingMode)
}

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          ddos_protection.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          ddos_protection.Get,
		ProductID:     ddos_protection.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[*products.EnableOutput, *productcore.NullInput]{
		OpNoInputFn: ddos_protection.Enable,
		ProductID:   ddos_protection.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      ddos_protection.Get,
		ProductID: ddos_protection.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[*ddos_protection.ConfigureOutput]{
		Phase:     "default",
		OpFn:      ddos_protection.GetConfiguration,
		ProductID: ddos_protection.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o *ddos_protection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "log", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewUpdateConfigurationTest(&productcore.UpdateConfigurationTestInput[*ddos_protection.ConfigureOutput, *ddos_protection.ConfigureInput]{
		OpFn:      ddos_protection.UpdateConfiguration,
		Input:     &ddos_protection.ConfigureInput{Mode: "block"},
		ProductID: ddos_protection.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o *ddos_protection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[*ddos_protection.ConfigureOutput]{
		Phase:     "after update",
		OpFn:      ddos_protection.GetConfiguration,
		ProductID: ddos_protection.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o *ddos_protection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      ddos_protection.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[*products.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          ddos_protection.Get,
		ProductID:     ddos_protection.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablementAndConfiguration(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
