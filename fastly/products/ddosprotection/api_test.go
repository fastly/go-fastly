package ddosprotection_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/fastly/products/ddosprotection"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/fastly/go-fastly/v9/internal/test_utils"

	"github.com/stretchr/testify/require"
)

var serviceID = fastly.TestDeliveryServiceID

func TestUpdateConfigurationMissingMode(t *testing.T) {
	t.Parallel()

	_, err := ddosprotection.UpdateConfiguration(nil, serviceID, ddosprotection.ConfigureInput{Mode: ""})

	require.ErrorIs(t, err, ddosprotection.ErrMissingMode)
}

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          ddosprotection.Disable,
		ServiceID:     serviceID,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ddosprotection.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          ddosprotection.Get,
		ProductID:     ddosprotection.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[ddosprotection.EnableOutput, products.NullInput]{
		OpNoInputFn: ddosprotection.Enable,
		ProductID:   ddosprotection.ProductID,
		ServiceID:   serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ddosprotection.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      ddosprotection.Get,
		ProductID: ddosprotection.ProductID,
		ServiceID: serviceID,
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[ddosprotection.ConfigureOutput]{
		Phase:     "default",
		OpFn:      ddosprotection.GetConfiguration,
		ProductID: ddosprotection.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o ddosprotection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "log", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewUpdateConfigurationTest(&productcore.UpdateConfigurationTestInput[ddosprotection.ConfigureOutput, ddosprotection.ConfigureInput]{
		OpFn:      ddosprotection.UpdateConfiguration,
		Input:     ddosprotection.ConfigureInput{Mode: "block"},
		ProductID: ddosprotection.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o ddosprotection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[ddosprotection.ConfigureOutput]{
		Phase:     "after update",
		OpFn:      ddosprotection.GetConfiguration,
		ProductID: ddosprotection.ProductID,
		ServiceID: serviceID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o ddosprotection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn:      ddosprotection.Disable,
		ServiceID: serviceID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ddosprotection.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          ddosprotection.Get,
		ProductID:     ddosprotection.ProductID,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablementAndConfiguration(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
